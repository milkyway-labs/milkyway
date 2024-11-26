package keeper

import (
	"context"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	restakingtypes "github.com/milkyway-labs/milkyway/v2/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v2/x/services/types"
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

// AfterServiceAccreditationModified implements servicestypes.ServicesHooks
func (k *Keeper) AfterServiceAccreditationModified(ctx context.Context, serviceID uint32) error {
	service, found, err := k.servicesKeeper.GetService(ctx, serviceID)
	if err != nil {
		return err
	}

	if !found {
		return servicestypes.ErrServiceNotFound
	}

	err = k.restakingKeeper.IterateServiceDelegations(ctx, serviceID, func(del restakingtypes.Delegation) (stop bool, err error) {
		preferences, err := k.restakingKeeper.GetUserPreferences(ctx, del.UserAddress)
		if err != nil {
			return true, err
		}

		// Clone the service and invert the accreditation status to get the
		// previous state
		serviceBefore := service
		serviceBefore.Accredited = !serviceBefore.Accredited

		trustedBefore := preferences.IsServiceTrusted(serviceBefore)
		trustedAfter := preferences.IsServiceTrusted(service)
		if trustedBefore != trustedAfter {
			err = k.AfterUserTrustedServiceUpdated(ctx, del.UserAddress, service.ID, trustedAfter)
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

// AfterUserPreferencesModified is called after a user's trust in a service is
// updated. It updates the total delegator shares for the service in all pools
// where the user has a delegation.
func (k *Keeper) AfterUserPreferencesModified(
	ctx context.Context,
	userAddress string,
	oldPreferences, newPreferences restakingtypes.UserPreferences,
) error {
	return k.servicesKeeper.IterateServices(ctx, func(service servicestypes.Service) (bool, error) {
		trustedBefore := oldPreferences.IsServiceTrusted(service)
		trustedAfter := newPreferences.IsServiceTrusted(service)
		if trustedBefore != trustedAfter {
			err := k.AfterUserTrustedServiceUpdated(ctx, userAddress, service.ID, trustedAfter)
			if err != nil {
				return true, err
			}
		}
		return false, nil
	})
}
