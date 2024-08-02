package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
)

var _ restakingtypes.RestakingHooks = Hooks{}

type Hooks struct {
	k *Keeper
}

func (k *Keeper) Hooks() Hooks {
	return Hooks{k}
}

func (h Hooks) BeforePoolDelegationCreated(ctx sdk.Context, poolID uint32, delegator string) error {
	return h.k.BeforeDelegationCreated(ctx, restakingtypes.DELEGATION_TYPE_POOL, poolID)
}

func (h Hooks) BeforePoolDelegationSharesModified(ctx sdk.Context, poolID uint32, delegator string) error {
	return h.k.BeforeDelegationSharesModified(ctx, restakingtypes.DELEGATION_TYPE_POOL, poolID, delegator)
}

func (h Hooks) AfterPoolDelegationModified(ctx sdk.Context, poolID uint32, delegator string) error {
	return h.k.AfterDelegationModified(ctx, restakingtypes.DELEGATION_TYPE_POOL, poolID, delegator)
}

func (h Hooks) BeforeOperatorDelegationCreated(ctx sdk.Context, operatorID uint32, delegator string) error {
	return h.k.BeforeDelegationCreated(ctx, restakingtypes.DELEGATION_TYPE_OPERATOR, operatorID)
}

func (h Hooks) BeforeOperatorDelegationSharesModified(ctx sdk.Context, operatorID uint32, delegator string) error {
	return h.k.BeforeDelegationSharesModified(ctx, restakingtypes.DELEGATION_TYPE_OPERATOR, operatorID, delegator)
}

func (h Hooks) AfterOperatorDelegationModified(ctx sdk.Context, operatorID uint32, delegator string) error {
	return h.k.AfterDelegationModified(ctx, restakingtypes.DELEGATION_TYPE_OPERATOR, operatorID, delegator)
}

func (h Hooks) BeforeServiceDelegationCreated(ctx sdk.Context, serviceID uint32, delegator string) error {
	return h.k.BeforeDelegationCreated(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, serviceID)
}

func (h Hooks) BeforeServiceDelegationSharesModified(ctx sdk.Context, serviceID uint32, delegator string) error {
	return h.k.BeforeDelegationSharesModified(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, serviceID, delegator)
}

func (h Hooks) AfterServiceDelegationModified(ctx sdk.Context, serviceID uint32, delegator string) error {
	return h.k.AfterDelegationModified(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, serviceID, delegator)
}

func (k *Keeper) BeforeDelegationCreated(ctx sdk.Context, delType restakingtypes.DelegationType, targetID uint32) error {
	target, err := k.GetDelegationTarget(ctx, delType, targetID)
	if err != nil {
		return err
	}

	// Initialize target if it doesn't exist yet.
	exists, err := k.HasCurrentRewards(ctx, target)
	if err != nil {
		return err
	}
	if !exists {
		if err := k.initializeDelegationTarget(ctx, target); err != nil {
			return err
		}
	}

	_, err = k.IncrementDelegationTargetPeriod(ctx, target)
	return err
}

func (k *Keeper) BeforeDelegationSharesModified(ctx sdk.Context, delType restakingtypes.DelegationType, targetID uint32, delegator string) error {
	target, err := k.GetDelegationTarget(ctx, delType, targetID)
	if err != nil {
		return err
	}

	// We don't have to initialize target here because we can assume BeforeDelegationCreated
	// has already been called when delegation shares are being modified.

	delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(delegator)
	if err != nil {
		return err
	}
	del, found := k.GetDelegation(ctx, target, delAddr)
	if !found {
		return sdkerrors.ErrNotFound.Wrapf("delegation not found: %d, %s", target.GetID(), delegator)
	}

	if _, err := k.withdrawDelegationRewards(ctx, target, del); err != nil {
		return err
	}

	return nil
}

func (k *Keeper) AfterDelegationModified(ctx sdk.Context, delType restakingtypes.DelegationType, targetID uint32, delegator string) error {
	delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(delegator)
	if err != nil {
		return err
	}
	return k.initializeDelegation(ctx, delType, targetID, delAddr)
}
