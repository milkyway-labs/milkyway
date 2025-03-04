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

// BeforeDelegationSharesModified is called before the existing delegation is
// being modified. If the delegator was a vesting investor, we decrement the
// validator's vesting investors shares so that we can re-increment it with the
// updated delegation shares after the delegation is modified.
func (h Hooks) BeforeDelegationSharesModified(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	isVestingInvestor, err := h.k.VestingInvestors.Has(ctx, delAddr)
	if err != nil {
		return err
	}
	if !isVestingInvestor {
		return nil
	}

	delegation, err := h.k.stakingKeeper.GetDelegation(ctx, delAddr, valAddr)
	if err != nil {
		return err
	}

	return h.k.DecrementValidatorInvestorsShares(ctx, valAddr, delegation.GetShares())
}

// AfterDelegationModified is called after a new delegation is created or the
// existing delegation is modified. If the delegator was a vesting investor, we
// increment the validator's vesting investors shares because we decremented it
// in BeforeDelegationSharesModified.
func (h Hooks) AfterDelegationModified(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	isVestingInvestor, err := h.k.VestingInvestors.Has(ctx, delAddr)
	if err != nil {
		return err
	}
	if !isVestingInvestor {
		return nil
	}

	delegation, err := h.k.stakingKeeper.GetDelegation(ctx, delAddr, valAddr)
	if err != nil {
		return err
	}

	return h.k.IncrementValidatorInvestorsShares(ctx, valAddr, delegation.GetShares())
}

func (h Hooks) BeforeDelegationRemoved(context.Context, sdk.AccAddress, sdk.ValAddress) error {
	return nil
}

func (h Hooks) BeforeValidatorSlashed(context.Context, sdk.ValAddress, sdkmath.LegacyDec) error {
	// TODO: implement it
	return nil
}

func (h Hooks) AfterValidatorCreated(context.Context, sdk.ValAddress) error {
	return nil
}

func (h Hooks) BeforeValidatorModified(context.Context, sdk.ValAddress) error {
	return nil
}

func (h Hooks) AfterValidatorRemoved(context.Context, sdk.ConsAddress, sdk.ValAddress) error {
	return nil
}

func (h Hooks) AfterValidatorBonded(context.Context, sdk.ConsAddress, sdk.ValAddress) error {
	return nil
}

func (h Hooks) AfterValidatorBeginUnbonding(context.Context, sdk.ConsAddress, sdk.ValAddress) error {
	return nil
}

func (h Hooks) BeforeDelegationCreated(context.Context, sdk.AccAddress, sdk.ValAddress) error {
	return nil
}

func (h Hooks) AfterUnbondingInitiated(_ context.Context, _ uint64) error {
	return nil
}
