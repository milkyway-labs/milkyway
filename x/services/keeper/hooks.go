package keeper

import (
	"context"

	"github.com/milkyway-labs/milkyway/x/services/types"
)

// Implement ServicesHooks interface
var _ types.ServicesHooks = &Keeper{}

// AfterServiceCreated implements ServicesHooks
func (k *Keeper) AfterServiceCreated(ctx context.Context, avsID uint32) {
	if k.hooks != nil {
		k.hooks.AfterServiceCreated(ctx, avsID)
	}
}

// AfterServiceActivated implements ServicesHooks
func (k *Keeper) AfterServiceActivated(ctx context.Context, avsID uint32) {
	if k.hooks != nil {
		k.hooks.AfterServiceActivated(ctx, avsID)
	}
}

// AfterServiceDeactivated implements ServicesHooks
func (k *Keeper) AfterServiceDeactivated(ctx context.Context, avsID uint32) {
	if k.hooks != nil {
		k.hooks.AfterServiceDeactivated(ctx, avsID)
	}
}
