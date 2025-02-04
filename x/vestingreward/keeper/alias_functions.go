package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	vestingexported "github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/milkyway-labs/milkyway/v7/x/vestingreward/types"
)

func (k *Keeper) GetValidatorVestingAccountsShares(ctx context.Context, valAddr sdk.ValAddress) (sdkmath.LegacyDec, error) {
	shares, err := k.ValidatorsVestingAccountsShares.Get(ctx, valAddr)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return sdkmath.LegacyZeroDec(), nil
		}
		return sdkmath.LegacyDec{}, err
	}
	return shares, nil
}

func (k *Keeper) IncrementValidatorVestingAccountsShares(ctx context.Context, valAddr sdk.ValAddress, shares sdkmath.LegacyDec) error {
	prevShares, err := k.GetValidatorVestingAccountsShares(ctx, valAddr)
	if err != nil {
		return err
	}
	return k.ValidatorsVestingAccountsShares.Set(ctx, valAddr, prevShares.Add(shares))
}

func (k *Keeper) DecrementValidatorVestingAccountsShares(ctx context.Context, valAddr sdk.ValAddress, shares sdkmath.LegacyDec) error {
	prevShares, err := k.GetValidatorVestingAccountsShares(ctx, valAddr)
	if err != nil {
		return err
	}
	newShares := prevShares.Sub(shares)
	if newShares.IsNegative() {
		panic("cannot set negative shares")
	} else if newShares.IsZero() {
		return k.ValidatorsVestingAccountsShares.Remove(ctx, valAddr)
	}
	return k.ValidatorsVestingAccountsShares.Set(ctx, valAddr, newShares)
}

// UpdateVestingAccountsRewardRatio updates the vesting accounts reward ratio.
// It also increments the period of all validators, since each validator's total
// tokens need to be adjusted after updating the ratio.
func (k *Keeper) UpdateVestingAccountsRewardRatio(ctx context.Context, ratio sdkmath.LegacyDec) error {
	err := types.ValidateVestingAccountsRewardRatio(ratio)
	if err != nil {
		return err
	}

	var innerErr error
	err = k.stakingKeeper.IterateValidators(ctx, func(_ int64, validator stakingtypes.ValidatorI) (stop bool) {
		var valAddr sdk.ValAddress
		valAddr, innerErr = k.stakingKeeper.ValidatorAddressCodec().StringToBytes(validator.GetOperator())
		if innerErr != nil {
			return true
		}
		validator, innerErr = k.GetAdjustedValidator(ctx, valAddr)
		if innerErr != nil {
			return true
		}
		_, innerErr = k.distrKeeper.IncrementValidatorPeriod(ctx, validator)
		return innerErr != nil
	})
	if err != nil {
		return err
	}
	if innerErr != nil {
		return innerErr
	}

	return k.VestingAccountsRewardRatio.Set(ctx, ratio)
}

// GetDelegation returns a delegation and a boolean flag indicating if the
// account is a vesting account.
func (k *Keeper) GetDelegation(
	ctx context.Context,
	delAddr sdk.AccAddress,
	valAddr sdk.ValAddress,
) (delegation stakingtypes.Delegation, isVestingAcc bool, err error) {
	delegation, err = k.stakingKeeper.GetDelegation(ctx, delAddr, valAddr)
	if err != nil {
		return
	}

	acc := k.accountKeeper.GetAccount(ctx, delAddr)
	_, isVestingAcc = acc.(vestingexported.VestingAccount)
	return delegation, isVestingAcc, nil
}

func (k *Keeper) GetAdjustedValidator(ctx context.Context, valAddr sdk.ValAddress) (stakingtypes.Validator, error) {
	validator, err := k.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		return stakingtypes.Validator{}, err
	}
	vestingAccountsShares, err := k.GetValidatorVestingAccountsShares(ctx, valAddr)
	if err != nil {
		return stakingtypes.Validator{}, err
	}
	if vestingAccountsShares.IsPositive() {
		rewardRatio, err := k.VestingAccountsRewardRatio.Get(ctx)
		if err != nil {
			return stakingtypes.Validator{}, err
		}
		oneMinusRewardRatio := sdkmath.LegacyOneDec().Sub(rewardRatio)
		validator, _ = validator.RemoveDelShares(vestingAccountsShares.MulRoundUp(oneMinusRewardRatio))
	}
	return validator, nil
}
