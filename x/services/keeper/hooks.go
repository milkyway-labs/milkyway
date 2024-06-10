package keeper

import (
	"context"

	"github.com/milkyway-labs/milkyway/x/services/types"
)

// Implement ServicesHooks interface
var _ types.ServicesHooks = &Keeper{}

// AfterServiceCreated implements ServicesHooks
func (k *Keeper) AfterServiceCreated(ctx context.Context, serviceID uint32) {
	if k.hooks != nil {
		k.hooks.AfterServiceCreated(ctx, serviceID)
	}
}

// AfterServiceActivated implements ServicesHooks
func (k *Keeper) AfterServiceActivated(ctx context.Context, serviceID uint32) {
	if k.hooks != nil {
		k.hooks.AfterServiceActivated(ctx, serviceID)
	}
}

// AfterServiceDeactivated implements ServicesHooks
func (k *Keeper) AfterServiceDeactivated(ctx context.Context, serviceID uint32) {
	if k.hooks != nil {
		k.hooks.AfterServiceDeactivated(ctx, serviceID)
	}
}
