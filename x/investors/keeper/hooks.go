package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	distrkeeper "github.com/milkyway-labs/milkyway/v10/x/distribution/keeper"
)

var (
	_ distrkeeper.DistrHooks    = DistrHooks{}
	_ stakingtypes.StakingHooks = StakingHooks{}
)

type DistrHooks struct {
	*Keeper
}

func (k *Keeper) DistrHooks() DistrHooks {
	return DistrHooks{k}
}

// BeforeWithdrawDelegationRewards is a hook called before a delegator's
// delegation rewards are withdrawn. It sets the current delegator in the
// context.
func (h DistrHooks) BeforeWithdrawDelegationRewards(ctx context.Context, delAddr sdk.AccAddress, _ sdk.ValAddress) error {
	delegator, err := h.accountKeeper.AddressCodec().BytesToString(delAddr)
	if err != nil {
		return err
	}
	err = h.SetCurrentDelegator(ctx, delegator)
	if err != nil {
		return err
	}
	return nil
}

// AfterWithdrawDelegationRewards is a hook called after a delegator's
// delegation rewards are withdrawn. It removes the current delegator from the
// context.
func (h DistrHooks) AfterWithdrawDelegationRewards(ctx context.Context, _ sdk.AccAddress, _ sdk.ValAddress, _ sdk.Coins) error {
	return h.RemoveCurrentDelegator(ctx)
}

type StakingHooks struct {
	*Keeper
}

func (k *Keeper) StakingHooks() StakingHooks {
	return StakingHooks{k}
}

// BeforeDelegationSharesModified is a hook called before a delegator's delegation
// shares are modified. It sets the current delegator in the context.
func (h StakingHooks) BeforeDelegationSharesModified(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	delegator, err := h.accountKeeper.AddressCodec().BytesToString(delAddr)
	if err != nil {
		return err
	}
	err = h.SetCurrentDelegator(ctx, delegator)
	if err != nil {
		return err
	}
	return nil
}

// AfterDelegationModified is a hook called after a delegator's delegation shares
// are modified. It removes the current delegator from the context.
func (h StakingHooks) AfterDelegationModified(ctx context.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return h.RemoveCurrentDelegator(ctx)
}

// Below are no-op implementations of the staking hooks

func (h StakingHooks) AfterValidatorCreated(context.Context, sdk.ValAddress) error {
	return nil
}

func (h StakingHooks) BeforeValidatorModified(context.Context, sdk.ValAddress) error {
	return nil
}

func (h StakingHooks) AfterValidatorRemoved(context.Context, sdk.ConsAddress, sdk.ValAddress) error {
	return nil
}

func (h StakingHooks) AfterValidatorBonded(context.Context, sdk.ConsAddress, sdk.ValAddress) error {
	return nil
}

func (h StakingHooks) AfterValidatorBeginUnbonding(context.Context, sdk.ConsAddress, sdk.ValAddress) error {
	return nil
}

func (h StakingHooks) BeforeDelegationCreated(context.Context, sdk.AccAddress, sdk.ValAddress) error {
	return nil
}

func (h StakingHooks) BeforeDelegationRemoved(context.Context, sdk.AccAddress, sdk.ValAddress) error {
	return nil
}

func (h StakingHooks) BeforeValidatorSlashed(context.Context, sdk.ValAddress, sdkmath.LegacyDec) error {
	return nil
}

func (h StakingHooks) AfterUnbondingInitiated(context.Context, uint64) error {
	return nil
}
