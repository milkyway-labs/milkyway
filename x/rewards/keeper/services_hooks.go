package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

var _ servicestypes.ServicesHooks = ServicesHooks{}

type ServicesHooks struct {
	k *Keeper
}

func (k *Keeper) ServicesHooks() ServicesHooks {
	return ServicesHooks{k}
}

// AfterServiceCreated implements servicestypes.ServicesHooks
func (h ServicesHooks) AfterServiceCreated(ctx sdk.Context, serviceID uint32) error {
	return h.k.AfterDelegationTargetCreated(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, serviceID)
}

// AfterServiceActivated implements servicestypes.ServicesHooks
func (h ServicesHooks) AfterServiceActivated(sdk.Context, uint32) error {
	return nil
}

// AfterServiceDeactivated implements servicestypes.ServicesHooks
func (h ServicesHooks) AfterServiceDeactivated(sdk.Context, uint32) error {
	return nil
}

// AfterServiceDeleted implements servicestypes.ServicesHooks
func (h ServicesHooks) AfterServiceDeleted(ctx sdk.Context, serviceID uint32) error {
	return h.k.AfterDelegationTargetRemoved(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, serviceID)
}

// AfterServiceAccreditationModified implements servicestypes.ServicesHooks
func (h ServicesHooks) AfterServiceAccreditationModified(ctx sdk.Context, serviceID uint32, accredited bool) error {
	err := h.k.restakingKeeper.IterateServiceDelegations(ctx, serviceID, func(del restakingtypes.Delegation) (stop bool, err error) {
		preferences, err := h.k.restakingKeeper.GetUserPreferences(ctx, del.UserAddress)
		if err != nil {
			return true, err
		}
		trustedBefore := preferences.IsServiceTrusted(serviceID, !accredited)
		trustedAfter := preferences.IsServiceTrusted(serviceID, accredited)
		if trustedBefore != trustedAfter {
			err = h.k.AfterUserTrustedServiceUpdated(ctx, del.UserAddress, serviceID, trustedAfter)
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
