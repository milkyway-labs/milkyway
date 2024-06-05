package types

// DONTCOVER

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Event Hooks
// These can be utilized to communicate between a services keeper
// and another keeper which must take particular actions when
// services change state. The second keeper must implement this
// interface, which then the services keeper can call.

// ServicesHooks event hooks for services objects (noalias)
type ServicesHooks interface {
	AfterServiceCreated(ctx sdk.Context, serviceID uint32)     // Must be called after a service is created
	AfterServiceActivated(ctx sdk.Context, serviceID uint32)   // Must be called after a service is registered
	AfterServiceDeactivated(ctx sdk.Context, serviceID uint32) // Must be called after a service is deregistered
}

// --------------------------------------------------------------------------------------------------------------------

// MultiServicesHooks combines multiple services hooks, all hook functions are run in array sequence
type MultiServicesHooks []ServicesHooks

// NewMultiServicesHooks creates a new MultiServicesHooks object
func NewMultiServicesHooks(hooks ...ServicesHooks) MultiServicesHooks {
	return hooks
}

// AfterServiceCreated implements ServicesHooks
func (m MultiServicesHooks) AfterServiceCreated(ctx sdk.Context, serviceID uint32) {
	for _, hook := range m {
		hook.AfterServiceCreated(ctx, serviceID)
	}
}

// AfterServiceActivated implements ServicesHooks
func (m MultiServicesHooks) AfterServiceActivated(ctx sdk.Context, serviceID uint32) {
	for _, hook := range m {
		hook.AfterServiceActivated(ctx, serviceID)
	}
}

// AfterServiceDeactivated implements ServicesHooks
func (m MultiServicesHooks) AfterServiceDeactivated(ctx sdk.Context, serviceID uint32) {
	for _, hook := range m {
		hook.AfterServiceDeactivated(ctx, serviceID)
	}
}
