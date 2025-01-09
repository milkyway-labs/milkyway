package keeper

import (
	"context"

	restakingtypes "github.com/milkyway-labs/milkyway/v7/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v7/x/services/types"
)

var _ servicestypes.ServicesHooks = ServicesHooks{}

type ServicesHooks struct {
	k *Keeper
}

func (k *Keeper) ServicesHooks() ServicesHooks {
	return ServicesHooks{k}
}

// AfterServiceCreated implements servicestypes.ServicesHooks
func (h ServicesHooks) AfterServiceCreated(ctx context.Context, serviceID uint32) error {
	return h.k.AfterDelegationTargetCreated(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, serviceID)
}

// AfterServiceActivated implements servicestypes.ServicesHooks
func (h ServicesHooks) AfterServiceActivated(context.Context, uint32) error {
	return nil
}

// AfterServiceDeactivated implements servicestypes.ServicesHooks
func (h ServicesHooks) AfterServiceDeactivated(context.Context, uint32) error {
	return nil
}

// BeforeServiceDeleted implements servicestypes.ServicesHooks
func (h ServicesHooks) BeforeServiceDeleted(ctx context.Context, serviceID uint32) error {
	return h.k.BeforeDelegationTargetRemoved(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, serviceID)
}

// AfterServiceAccreditationModified implements servicestypes.ServicesHooks
func (h ServicesHooks) AfterServiceAccreditationModified(context.Context, uint32) error {
	return nil
}
