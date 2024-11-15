package keeper

import (
	"context"

	"github.com/milkyway-labs/milkyway/x/services/types"
)

// Implement ServicesHooks interface
var _ types.ServicesHooks = &Keeper{}

// AfterServiceCreated implements ServicesHooks
func (k *Keeper) AfterServiceCreated(ctx context.Context, serviceID uint32) error {
	if k.hooks != nil {
		return k.hooks.AfterServiceCreated(ctx, serviceID)
	}
	return nil
}

// AfterServiceActivated implements ServicesHooks
func (k *Keeper) AfterServiceActivated(ctx context.Context, serviceID uint32) error {
	if k.hooks != nil {
		return k.hooks.AfterServiceActivated(ctx, serviceID)
	}
	return nil
}

// AfterServiceDeactivated implements ServicesHooks
func (k *Keeper) AfterServiceDeactivated(ctx context.Context, serviceID uint32) error {
	if k.hooks != nil {
		return k.hooks.AfterServiceDeactivated(ctx, serviceID)
	}
	return nil
}

// AfterServiceDeleted implements ServicesHooks
func (k *Keeper) AfterServiceDeleted(ctx context.Context, serviceID uint32) error {
	if k.hooks != nil {
		return k.hooks.AfterServiceDeleted(ctx, serviceID)
	}
	return nil
}

// AfterServiceAccreditationModified implements ServicesHooks
func (k *Keeper) AfterServiceAccreditationModified(ctx sdk.Context, serviceID uint32) error {
	if k.hooks != nil {
		return k.hooks.AfterServiceAccreditationModified(ctx, serviceID)
	}
	return nil
}
