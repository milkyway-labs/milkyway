package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var _ stakingtypes.StakingHooks = Hooks{}

type Hooks struct {
	k *Keeper
}

func (k *Keeper) Hooks() Hooks {
	return Hooks{k}
}

func (h Hooks) BeforeDelegationSharesModified(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	delegation, isVestingAcc, err := h.k.GetDelegation(ctx, delAddr, valAddr)
	if err != nil {
		return err
	}
	// If the account is not a vesting account, do nothing
	if !isVestingAcc {
		return nil
	}

	return h.k.DecrementValidatorVestingAccountsShares(ctx, valAddr, delegation.GetShares())
}

func (h Hooks) AfterDelegationModified(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	delegation, isVestingAcc, err := h.k.GetDelegation(ctx, delAddr, valAddr)
	if err != nil {
		return err
	}
	// If the account is not a vesting account, do nothing
	if !isVestingAcc {
		return nil
	}

	return h.k.IncrementValidatorVestingAccountsShares(ctx, valAddr, delegation.GetShares())
}

func (h Hooks) BeforeDelegationRemoved(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	// Just call BeforeDelegationSharesModified hook since all we want to do is
	// decrementing the shares
	return h.BeforeDelegationSharesModified(ctx, delAddr, valAddr)
}

func (h Hooks) BeforeValidatorSlashed(_ context.Context, _ sdk.ValAddress, _ sdkmath.LegacyDec) error {
	// TODO: implement it
	return nil
}

func (h Hooks) AfterValidatorCreated(_ context.Context, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) BeforeValidatorModified(_ context.Context, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) AfterValidatorRemoved(_ context.Context, _ sdk.ConsAddress, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) AfterValidatorBonded(_ context.Context, _ sdk.ConsAddress, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) AfterValidatorBeginUnbonding(_ context.Context, _ sdk.ConsAddress, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) BeforeDelegationCreated(_ context.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}

func (h Hooks) AfterUnbondingInitiated(_ context.Context, _ uint64) error {
	return nil
}
