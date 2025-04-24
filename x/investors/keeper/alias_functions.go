package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	vestingexported "github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/milkyway-labs/milkyway/v11/x/investors/types"
)

// SetInvestorsRewardRatio sets the investors reward ratio.
func (k *Keeper) SetInvestorsRewardRatio(ctx context.Context, ratio sdkmath.LegacyDec) error {
	err := types.ValidateInvestorsRewardRatio(ratio)
	if err != nil {
		return err
	}
	return k.InvestorsRewardRatio.Set(ctx, ratio)
}

// GetInvestorsRewardRatio returns the investors reward ratio.
func (k *Keeper) GetInvestorsRewardRatio(ctx context.Context) (sdkmath.LegacyDec, error) {
	return k.InvestorsRewardRatio.Get(ctx)
}

// UpdateInvestorsRewardRatio updates the investors reward ratio. It withdraws
// rewards from all the validators that the investors were delegating to.
func (k *Keeper) UpdateInvestorsRewardRatio(ctx context.Context, ratio sdkmath.LegacyDec) error {
	// Forcefully withdraw all vesting investors' staking rewards
	investors, err := k.GetAllVestingInvestorsAddresses(ctx)
	if err != nil {
		return err
	}
	for _, investor := range investors {
		err = k.WithdrawAllDelegationRewards(ctx, investor)
		if err != nil {
			return err
		}
	}

	return k.SetInvestorsRewardRatio(ctx, ratio)
}

// --------------------------------------------------------------------------------------------------------------------

// GetAllVestingInvestorsAddresses returns all the vesting investors addresses.
func (k *Keeper) GetAllVestingInvestorsAddresses(ctx context.Context) ([]string, error) {
	var investors []string
	err := k.VestingInvestors.Walk(ctx, nil, func(addr string) (stop bool, err error) {
		investors = append(investors, addr)
		return false, nil
	})
	return investors, err
}

// SetVestingInvestor sets an account as a vesting investor. It returns an error
// if the account does not exist or is not a vesting account. It also adds the
// account to the vesting investors queue.
func (k *Keeper) SetVestingInvestor(ctx context.Context, addr string) error {
	accAddr, err := k.accountKeeper.AddressCodec().StringToBytes(addr)
	if err != nil {
		return err
	}
	acc := k.accountKeeper.GetAccount(ctx, accAddr)
	if acc == nil {
		return sdkerrors.ErrUnknownAddress.Wrapf("account %s does not exist", addr)
	}
	vacc, isVestingAcc := acc.(vestingexported.VestingAccount)
	if !isVestingAcc {
		return sdkerrors.ErrInvalidRequest.Wrapf("account %s is not a vesting account", addr)
	}
	// If the vesting account's end time is in the past, do nothing and return early
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if vacc.GetEndTime() <= sdkCtx.BlockTime().Unix() {
		return nil
	}

	// Withdraw all delegation rewards accrued so far
	err = k.WithdrawAllDelegationRewards(ctx, addr)
	if err != nil {
		return err
	}

	err = k.InvestorsVestingQueue.Set(ctx, collections.Join(vacc.GetEndTime(), addr))
	if err != nil {
		return err
	}
	return k.VestingInvestors.Set(ctx, addr)
}

// IsVestingInvestor returns true if the account is a vesting investor.
func (k *Keeper) IsVestingInvestor(ctx context.Context, addr string) (bool, error) {
	return k.VestingInvestors.Has(ctx, addr)
}

// RemoveVestingInvestor removes an account from the vesting investors list.
// Rewards accrued so far will be withdrawn automatically.
func (k *Keeper) RemoveVestingInvestor(ctx context.Context, addr string) error {
	// Forcefully withdraw all delegation rewards
	err := k.WithdrawAllDelegationRewards(ctx, addr)
	if err != nil {
		return err
	}

	return k.VestingInvestors.Remove(ctx, addr)
}

// --------------------------------------------------------------------------------------------------------------------

// WithdrawAllDelegationRewards withdraws all the staking rewards allocated to
// the delegator.
func (k *Keeper) WithdrawAllDelegationRewards(ctx context.Context, delegator string) error {
	delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(delegator)
	if err != nil {
		return err
	}

	var innerErr error
	err = k.stakingKeeper.IterateDelegatorDelegations(ctx, delAddr, func(delegation stakingtypes.Delegation) (stop bool) {
		var valAddr sdk.ValAddress
		valAddr, innerErr = k.stakingKeeper.ValidatorAddressCodec().StringToBytes(delegation.ValidatorAddress)
		if innerErr != nil {
			return true
		}

		_, innerErr = k.distrKeeper.WithdrawDelegationRewards(ctx, delAddr, valAddr)
		if innerErr != nil {
			return true
		}
		return innerErr != nil
	})
	if err != nil {
		return err
	}
	return innerErr
}
