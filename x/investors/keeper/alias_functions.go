package keeper

import (
	"context"
	"errors"
	"slices"

	"cosmossdk.io/collections"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	vestingexported "github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/milkyway-labs/milkyway/v9/x/investors/types"
)

// GetAllVestingInvestorsAddresses returns all the vesting investors addresses.
func (k *Keeper) GetAllVestingInvestorsAddresses(ctx context.Context) ([]string, error) {
	var investors []string
	err := k.VestingInvestors.Walk(ctx, nil, func(investorAddr sdk.AccAddress) (stop bool, err error) {
		investor, err := k.accountKeeper.AddressCodec().BytesToString(investorAddr)
		if err != nil {
			return true, err
		}
		investors = append(investors, investor)
		return false, nil
	})
	return investors, err
}

// SetVestingInvestor sets an account as a vesting investor. It returns an error
// if the account does not exist or is not a vesting account. It also adds the
// account to the vesting investors queue.
func (k *Keeper) SetVestingInvestor(ctx context.Context, addr sdk.AccAddress) error {
	acc := k.accountKeeper.GetAccount(ctx, addr)
	if acc == nil {
		return sdkerrors.ErrUnknownAddress.Wrapf("account %s does not exist", addr)
	}
	vacc, isVestingAcc := acc.(vestingexported.VestingAccount)
	if !isVestingAcc {
		return sdkerrors.ErrInvalidRequest.Wrapf("account %s is not a vesting account", addr)
	}

	err := k.InvestorsVestingQueue.Set(ctx, collections.Join(vacc.GetEndTime(), addr))
	if err != nil {
		return err
	}
	return k.VestingInvestors.Set(ctx, addr)
}

// --------------------------------------------------------------------------------------------------------------------

func (k *Keeper) GetValidatorInvestorsShares(ctx context.Context, valAddr sdk.ValAddress) (sdkmath.LegacyDec, error) {
	shares, err := k.ValidatorsInvestorsShares.Get(ctx, valAddr)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return sdkmath.LegacyZeroDec(), nil
		}
		return sdkmath.LegacyDec{}, err
	}
	return shares, nil
}

func (k *Keeper) IncrementValidatorInvestorsShares(ctx context.Context, valAddr sdk.ValAddress, shares sdkmath.LegacyDec) error {
	prevShares, err := k.GetValidatorInvestorsShares(ctx, valAddr)
	if err != nil {
		return err
	}
	return k.ValidatorsInvestorsShares.Set(ctx, valAddr, prevShares.Add(shares))
}

func (k *Keeper) DecrementValidatorInvestorsShares(ctx context.Context, valAddr sdk.ValAddress, shares sdkmath.LegacyDec) error {
	prevShares, err := k.GetValidatorInvestorsShares(ctx, valAddr)
	if err != nil {
		return err
	}
	newShares := prevShares.Sub(shares)
	if newShares.IsNegative() {
		panic("cannot set negative shares")
	} else if newShares.IsZero() {
		return k.ValidatorsInvestorsShares.Remove(ctx, valAddr)
	}
	return k.ValidatorsInvestorsShares.Set(ctx, valAddr, newShares)
}

// --------------------------------------------------------------------------------------------------------------------

// IncrementValidatorPeriod increments the period of a validator using
// GetAdjustedValidator.
func (k *Keeper) IncrementValidatorPeriod(ctx context.Context, valAddr sdk.ValAddress) error {
	validator, err := k.GetAdjustedValidator(ctx, valAddr)
	if err != nil {
		return err
	}
	_, err = k.distrKeeper.IncrementValidatorPeriod(ctx, validator)
	return err
}

// UpdateInvestorsRewardRatio updates the investors reward ratio. It also
// increments the period of all validators to finalize the current epoch's
// cumulated rewards before updating the ratio.
func (k *Keeper) UpdateInvestorsRewardRatio(ctx context.Context, ratio sdkmath.LegacyDec) error {
	err := types.ValidateInvestorsRewardRatio(ratio)
	if err != nil {
		return err
	}

	// Get all validators that the vesting investors are delegating to.
	var valAddrs []sdk.ValAddress
	investors, err := k.GetAllVestingInvestorsAddresses(ctx)
	if err != nil {
		return err
	}
	for _, investor := range investors {
		investorAddr, err := k.accountKeeper.AddressCodec().StringToBytes(investor)
		if err != nil {
			return err
		}
		delegations, err := k.stakingKeeper.GetAllDelegatorDelegations(ctx, investorAddr)
		for _, delegation := range delegations {
			valAddr, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(delegation.ValidatorAddress)
			if err != nil {
				return err
			}
			valAddrs = append(valAddrs, valAddr)
		}
	}

	// Remove duplicated addresses
	valAddrs = slices.CompactFunc(valAddrs, func(a, b sdk.ValAddress) bool {
		return a.Equals(b)
	})

	for _, valAddr := range valAddrs {
		err = k.IncrementValidatorPeriod(ctx, valAddr)
		if err != nil {
			return err
		}
	}

	return k.InvestorsRewardRatio.Set(ctx, ratio)
}

// --------------------------------------------------------------------------------------------------------------------

func (k *Keeper) GetAdjustedValidator(ctx context.Context, valAddr sdk.ValAddress) (stakingtypes.Validator, error) {
	validator, err := k.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		return stakingtypes.Validator{}, err
	}
	investorsShares, err := k.GetValidatorInvestorsShares(ctx, valAddr)
	if err != nil {
		return stakingtypes.Validator{}, err
	}
	if investorsShares.IsPositive() {
		rewardRatio, err := k.InvestorsRewardRatio.Get(ctx)
		if err != nil {
			return stakingtypes.Validator{}, err
		}
		oneMinusRewardRatio := sdkmath.LegacyOneDec().Sub(rewardRatio)
		validator, _ = validator.RemoveDelShares(investorsShares.MulRoundUp(oneMinusRewardRatio))
	}
	return validator, nil
}

// --------------------------------------------------------------------------------------------------------------------

// RemoveVestingInvestor removes an investor from the vesting investors list and
// withdraws rewards from all the validators that the investor was delegating to.
func (k *Keeper) RemoveVestingInvestor(ctx context.Context, investorAddr sdk.AccAddress) error {
	// Remove from VestingInvestors first so that inside WithdrawDelegationRewards
	// the account's new starting info can be set correctly using non-deducted
	// delegation shares. Note that WithdrawDelegationRewards uses the previous
	// starting info's stake(=tokens) to calculate the rewards so the change to the
	// delegation shares prior to calling it doesn't affect the rewards amount.
	err := k.VestingInvestors.Remove(ctx, investorAddr)
	if err != nil {
		return err
	}

	// Withdraw rewards from all validators that the investor was delegating to.
	delegations, err := k.stakingKeeper.GetAllDelegatorDelegations(ctx, investorAddr)
	if err != nil {
		return err
	}
	for _, delegation := range delegations {
		valAddr, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(delegation.ValidatorAddress)
		if err != nil {
			return err
		}
		// Calling WithdrawDelegationRewards increments the validator's period so no need
		// to call IncrementValidatorPeriod here.
		_, err = k.distrKeeper.WithdrawDelegationRewards(ctx, investorAddr, valAddr)
		if err != nil {
			return err
		}

		err = k.DecrementValidatorInvestorsShares(ctx, valAddr, delegation.Shares)
		if err != nil {
			return err
		}
	}

	return nil
}
