package types

// DONTCOVER

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Event Hooks
// These can be utilized to communicate between an avs keeper and another
// keeper which must take particular actions when AVSs change
// state. The second keeper must implement this interface, which then the
// avs keeper can call.

// AVSHooks event hooks for avs objects (noalias)
type AVSHooks interface {
	AfterAVSRegistered(ctx sdk.Context, avsID uint32)   // Must be called after an AVS is registered
	AfterAVSDeregistered(ctx sdk.Context, avsID uint32) // Must be called after an AVS is deregistered
}

// --------------------------------------------------------------------------------------------------------------------

// MultiAVSHooks combines multiple avs hooks, all hook functions are run in array sequence
type MultiAVSHooks []AVSHooks

// NewMultiAVSHooks creates a new MultiAVSHooks object
func NewMultiAVSHooks(hooks ...AVSHooks) MultiAVSHooks {
	return hooks
}

// AfterAVSRegistered implements AVSHooks
func (m MultiAVSHooks) AfterAVSRegistered(ctx sdk.Context, avsID uint32) {
	for _, hook := range m {
		hook.AfterAVSRegistered(ctx, avsID)
	}
}

// AfterAVSDeregistered implements AVSHooks
func (m MultiAVSHooks) AfterAVSDeregistered(ctx sdk.Context, avsID uint32) {
	for _, hook := range m {
		hook.AfterAVSDeregistered(ctx, avsID)
	}
}
