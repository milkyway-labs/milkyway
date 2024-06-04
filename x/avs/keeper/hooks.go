package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/avs/types"
)

// Implement AVSHooks interface
var _ types.AVSHooks = &Keeper{}

// AfterAVSCreated implements AVSHooks
func (k Keeper) AfterAVSCreated(ctx sdk.Context, avsID uint32) {
	if k.hooks != nil {
		k.hooks.AfterAVSCreated(ctx, avsID)
	}
}

// AfterAVSRegistered implements AVSHooks
func (k Keeper) AfterAVSRegistered(ctx sdk.Context, avsID uint32) {
	if k.hooks != nil {
		k.hooks.AfterAVSRegistered(ctx, avsID)
	}
}

// AfterAVSDeregistered implements AVSHooks
func (k Keeper) AfterAVSDeregistered(ctx sdk.Context, avsID uint32) {
	if k.hooks != nil {
		k.hooks.AfterAVSDeregistered(ctx, avsID)
	}
}
