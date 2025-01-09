package keeper

import (
	"context"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/milkyway-labs/milkyway/v7/utils"
	restakingtypes "github.com/milkyway-labs/milkyway/v7/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v7/x/services/types"
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

		servicesIDs, err := k.restakingKeeper.GetUserTrustedServicesIDs(ctx, delegator)
		if err != nil {
			return err
		}

		for _, serviceID := range servicesIDs {
			// We decremented the amount of shares within the pool-service pair in the
			// BeforeDelegationSharesModified hook. We increment the shares here again
			// to keep consistency if the shares change due to a new delegation or an
			// undelegation
			err = k.IncrementPoolServiceTotalDelegatorShares(ctx, targetID, serviceID, delegation.Shares)
			if err != nil {
				return err
			}
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

	allServices, err := k.servicesKeeper.GetServices(ctx)
	if err != nil {
		return err
	}
	allServicesIDs := utils.Map(allServices, func(service servicestypes.Service) uint32 {
		return service.ID
	})

	err = k.restakingKeeper.IterateUserPoolDelegations(ctx, userAddress, func(delegation restakingtypes.Delegation) (stop bool, err error) {
		poolID := delegation.TargetID

		pool, err := k.GetDelegationTarget(ctx, restakingtypes.DELEGATION_TYPE_POOL, poolID)
		if err != nil {
			return true, err
		}

		for _, serviceID := range allServicesIDs {
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
