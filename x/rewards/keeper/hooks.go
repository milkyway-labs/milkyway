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
	if err != nil {
		return err
	}

	if delType == restakingtypes.DELEGATION_TYPE_POOL {
		servicesIDs, err := k.restakingKeeper.GetUserTrustedServicesIDs(ctx, delegator)
		if err != nil {
			return err
		}
		for _, serviceID := range servicesIDs {
			// We decrement the amount of shares within the pool-service pair here so that we
			// can later increment those shares again within the AfterDelegationModified
			// hook. This is due in order to keep consistency if the shares change due to a
			// new delegation or an undelegation
			err = k.DecrementPoolServiceTotalDelegatorShares(ctx, targetID, serviceID, del.Shares)
			if err != nil {
				return err
			}
		}
	}
	return nil
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

	err = k.initializeDelegation(ctx, target, delAddr)
	if err != nil {
		return err
	}
	if delType == restakingtypes.DELEGATION_TYPE_POOL {
		del, found := k.restakingKeeper.GetPoolDelegation(ctx, targetID, delegator)
		if !found {
			return sdkerrors.ErrNotFound.Wrapf("pool delegation not found: %d, %s", targetID, delegator)
		}

		servicesIDs, err := k.restakingKeeper.GetUserTrustedServicesIDs(ctx, delegator)
		if err != nil {
			return err
		}
		for _, serviceID := range servicesIDs {
			// We decremented the amount of shares within the pool-service pair in the
			// BeforeDelegationSharesModified hook. We increment the shares here again
			// to keep consistency if the shares change due to a new delegation or an
			// undelegation
			err = k.IncrementPoolServiceTotalDelegatorShares(ctx, targetID, serviceID, del.Shares)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// AfterUserTrustedServiceUpdated is called after a user's trust in a service is
// updated. It updates the total delegator shares for the service in all pools
// where the user has a delegation.
func (k *Keeper) AfterUserTrustedServiceUpdated(ctx sdk.Context, userAddress string, serviceID uint32, trusted bool) error {
	delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(userAddress)
	if err != nil {
		return err
	}

	err = k.restakingKeeper.IterateUserPoolDelegations(ctx, userAddress, func(del restakingtypes.Delegation) (stop bool, err error) {
		isSecured, err := k.restakingKeeper.IsServiceSecuredByPool(ctx, serviceID, del.TargetID)
		if err != nil {
			return true, err
		}
		if !isSecured {
			return false, nil
		}

		pool, err := k.GetDelegationTarget(ctx, restakingtypes.DELEGATION_TYPE_POOL, del.TargetID)
		if err != nil {
			return true, err
		}

		// Calling these two methods has same effect as calling
		// BeforePoolDelegationSharesModified and then AfterPoolDelegationModified.
		_, err = k.withdrawDelegationRewards(ctx, pool, del)
		if err != nil {
			return true, err
		}
		err = k.initializeDelegation(ctx, pool, delAddr)
		if err != nil {
			return true, err
		}

		if trusted {
			err = k.IncrementPoolServiceTotalDelegatorShares(ctx, del.TargetID, serviceID, del.Shares)
		} else {
			err = k.DecrementPoolServiceTotalDelegatorShares(ctx, del.TargetID, serviceID, del.Shares)
		}
		if err != nil {
			return true, err
		}
		return false, nil
	})
	if err != nil {
		return err
	}
	return nil
}
