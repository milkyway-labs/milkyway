package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
)

// AfterDelegationTargetCreated is called after a delegation target is created
func (k *Keeper) AfterDelegationTargetCreated(ctx sdk.Context, delType restakingtypes.DelegationType, targetID uint32) error {
	target, err := k.GetDelegationTarget(ctx, delType, targetID)
	if err != nil {
		return err
	}

	return k.initializeDelegationTarget(ctx, target)
}

// AfterDelegationTargetRemoved is called after a delegation target is removed
func (k *Keeper) AfterDelegationTargetRemoved(ctx sdk.Context, delType restakingtypes.DelegationType, targetID uint32) error {
	target, err := k.GetDelegationTarget(ctx, delType, targetID)
	if err != nil {
		return err
	}

	return k.clearDelegationTarget(ctx, target)
}

// BeforeDelegationCreated is called before a delegation to a target is created
func (k *Keeper) BeforeDelegationCreated(ctx sdk.Context, delType restakingtypes.DelegationType, targetID uint32) error {
	target, err := k.GetDelegationTarget(ctx, delType, targetID)
	if err != nil {
		return err
	}

	_, err = k.IncrementDelegationTargetPeriod(ctx, target)
	return err
}

// BeforeDelegationSharesModified is called before a delegation to a target is modified
func (k *Keeper) BeforeDelegationSharesModified(ctx sdk.Context, delType restakingtypes.DelegationType, targetID uint32, delegator string) error {
	target, err := k.GetDelegationTarget(ctx, delType, targetID)
	if err != nil {
		return err
	}

	// We don't have to initialize target here because we can assume BeforeDelegationCreated
	// has already been called when delegation shares are being modified.

	del, found := k.restakingKeeper.GetDelegationForTarget(ctx, target, delegator)
	if !found {
		return sdkerrors.ErrNotFound.Wrapf("delegation not found: %d, %s", target.GetID(), delegator)
	}

	_, err = k.withdrawDelegationRewards(ctx, target, del)
	return err
}

// AfterDelegationModified is called after a delegation to a target is modified
func (k *Keeper) AfterDelegationModified(ctx sdk.Context, delType restakingtypes.DelegationType, targetID uint32, delegator string) error {
	target, err := k.GetDelegationTarget(ctx, delType, targetID)
	if err != nil {
		return err
	}

	delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(delegator)
	if err != nil {
		return err
	}

	return k.initializeDelegation(ctx, target, delAddr)
}
