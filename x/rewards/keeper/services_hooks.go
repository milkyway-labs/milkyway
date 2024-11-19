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
func (h ServicesHooks) AfterServiceAccreditationModified(ctx sdk.Context, serviceID uint32) error {
	return h.k.AfterServiceAccreditationModified(ctx, serviceID)
}
