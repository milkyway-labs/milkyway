package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/services/types"
)

// Implement ServicesHooks interface
var _ types.ServicesHooks = &Keeper{}

// AfterServiceCreated implements ServicesHooks
func (k *Keeper) AfterServiceCreated(ctx sdk.Context, avsID uint32) {
	if k.hooks != nil {
		k.hooks.AfterServiceCreated(ctx, avsID)
	}
}

// AfterServiceRegistered implements ServicesHooks
func (k *Keeper) AfterServiceActivated(ctx sdk.Context, avsID uint32) {
	if k.hooks != nil {
		k.hooks.AfterServiceActivated(ctx, avsID)
	}
}

// AfterServiceDeregistered implements ServicesHooks
func (k *Keeper) AfterServiceDeactivated(ctx sdk.Context, avsID uint32) {
	if k.hooks != nil {
		k.hooks.AfterServiceDeactivated(ctx, avsID)
	}
}
