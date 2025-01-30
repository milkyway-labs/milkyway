package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	restakingtypes "github.com/milkyway-labs/milkyway/v8/x/restaking/types"
	"github.com/milkyway-labs/milkyway/v8/x/rewards/types"
)

// AfterDelegationTargetCreated is called after a delegation target is created
func (k *Keeper) AfterDelegationTargetCreated(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32) error {
	target, err := k.GetDelegationTarget(ctx, delType, targetID)
	if err != nil {
		return err
	}

	return k.initializeDelegationTarget(ctx, target)
}

// BeforeDelegationTargetRemoved is called before a delegation target is removed
func (k *Keeper) BeforeDelegationTargetRemoved(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32) error {
	target, err := k.GetDelegationTarget(ctx, delType, targetID)
	if err != nil {
		return err
	}

	return k.clearDelegationTarget(ctx, target)
}

// BeforeDelegationCreated is called before a delegation to a target is created
func (k *Keeper) BeforeDelegationCreated(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32) error {
	target, err := k.GetDelegationTarget(ctx, delType, targetID)
	if err != nil {
		return err
	}

	_, err = k.IncrementDelegationTargetPeriod(ctx, target)
	return err
}

// BeforeDelegationSharesModified is called before a delegation to a target is modified
func (k *Keeper) BeforeDelegationSharesModified(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32, delegator string) error {
	target, err := k.GetDelegationTarget(ctx, delType, targetID)
	if err != nil {
		return err
	}

	// We don't have to initialize target here because we can assume BeforeDelegationCreated
	// has already been called when delegation shares are being modified.
	del, found, err := k.restakingKeeper.GetDelegationForTarget(ctx, target.DelegationTarget, delegator)
	if err != nil {
		return err
	}

	if !found {
		return sdkerrors.ErrNotFound.Wrapf("delegation not found: %d, %s", target.GetID(), delegator)
	}

	_, err = k.withdrawDelegationRewards(ctx, target, del)
	if err != nil {
		return err
	}

	if delType == restakingtypes.DELEGATION_TYPE_POOL {
		preferences, err := k.restakingKeeper.GetUserPreferences(ctx, delegator)
		if err != nil {
			return err
		}

		for _, entry := range preferences.TrustedServices {
			if preferences.IsServiceTrustedWithPool(entry.ServiceID, targetID) {
				// We decrement the amount of shares within the pool-service pair here so that we
				// can later increment those shares again within the AfterDelegationModified
				// hook. This is due in order to keep consistency if the shares change due to a
				// new delegation or an undelegation
				err = k.DecrementPoolServiceTotalDelegatorShares(ctx, targetID, entry.ServiceID, del.Shares)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// AfterDelegationModified is called after a delegation to a target is modified
func (k *Keeper) AfterDelegationModified(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32, delegator string) error {
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
		delegation, found, err := k.restakingKeeper.GetPoolDelegation(ctx, targetID, delegator)
		if err != nil {
			return err
		}

		if !found {
			return sdkerrors.ErrNotFound.Wrapf("pool delegation not found: %d, %s", targetID, delegator)
		}

		preferences, err := k.restakingKeeper.GetUserPreferences(ctx, delegator)
		if err != nil {
			return err
		}

		for _, entry := range preferences.TrustedServices {
			if preferences.IsServiceTrustedWithPool(entry.ServiceID, targetID) {
				// We decremented the amount of shares within the pool-service pair in the
				// BeforeDelegationSharesModified hook. We increment the shares here again
				// to keep consistency if the shares change due to a new delegation or an
				// undelegation
				err = k.IncrementPoolServiceTotalDelegatorShares(ctx, targetID, entry.ServiceID, delegation.Shares)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// BeforeServiceDeleted is called before a service is deleted
func (k *Keeper) BeforeServiceDeleted(ctx context.Context, serviceID uint32) error {
	err := k.BeforeDelegationTargetRemoved(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, serviceID)
	if err != nil {
		return err
	}

	// Clean up pool-service total delegator shares for the service
	var keysToDelete []collections.Pair[uint32, uint32]
	err = k.PoolServiceTotalDelegatorShares.Walk(ctx, nil, func(key collections.Pair[uint32, uint32], _ types.PoolServiceTotalDelegatorShares) (stop bool, err error) {
		if key.K2() != serviceID {
			return false, nil
		}
		keysToDelete = append(keysToDelete, key)
		return false, nil
	})
	if err != nil {
		return err
	}

	for _, key := range keysToDelete {
		err = k.PoolServiceTotalDelegatorShares.Remove(ctx, key)
		if err != nil {
			return err
		}
	}

	return nil
}

// AfterUserPreferencesModified is called after a user's trust in a service is
// updated. It updates the total delegator shares for the service in all pools
// where the user has a delegation.
func (k *Keeper) AfterUserPreferencesModified(
	ctx context.Context,
	userAddress string,
	oldPreferences restakingtypes.UserPreferences,
	newPreferences restakingtypes.UserPreferences,
) error {
	delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(userAddress)
	if err != nil {
		return err
	}

	changedServicesIDs := restakingtypes.ComputeChangedServicesIDs(oldPreferences, newPreferences)

	err = k.restakingKeeper.IterateUserPoolDelegations(ctx, userAddress, func(delegation restakingtypes.Delegation) (stop bool, err error) {
		poolID := delegation.TargetID

		pool, err := k.GetDelegationTarget(ctx, restakingtypes.DELEGATION_TYPE_POOL, poolID)
		if err != nil {
			return true, err
		}

		for _, serviceID := range changedServicesIDs {
			trustedBefore := oldPreferences.IsServiceTrustedWithPool(serviceID, poolID)
			trustedAfter := newPreferences.IsServiceTrustedWithPool(serviceID, poolID)

			if trustedBefore == trustedAfter {
				continue
			}

			// Calling these two methods has same effect as calling
			// BeforePoolDelegationSharesModified and then AfterPoolDelegationModified.
			_, err = k.withdrawDelegationRewards(ctx, pool, delegation)
			if err != nil {
				return true, err
			}

			err = k.initializeDelegation(ctx, pool, delAddr)
			if err != nil {
				return true, err
			}

			if trustedAfter {
				err = k.IncrementPoolServiceTotalDelegatorShares(ctx, poolID, serviceID, delegation.Shares)
			} else {
				err = k.DecrementPoolServiceTotalDelegatorShares(ctx, poolID, serviceID, delegation.Shares)
			}

			if err != nil {
				return true, err
			}
		}

		return false, nil
	})
	if err != nil {
		return err
	}

	return nil
}
