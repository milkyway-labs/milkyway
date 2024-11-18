package types

import (
	"context"
)

// DONTCOVER

// Event Hooks
// These can be utilized to communicate between a services keeper
// and another keeper which must take particular actions when
// services change state. The second keeper must implement this
// interface, which then the services keeper can call.

// ServicesHooks event hooks for services objects (noalias)
type ServicesHooks interface {
	AfterServiceCreated(ctx context.Context, serviceID uint32) error               // Must be called after a service is created
	AfterServiceActivated(ctx context.Context, serviceID uint32) error             // Must be called after a service is registered
	AfterServiceDeactivated(ctx context.Context, serviceID uint32) error           // Must be called after a service is deregistered
	AfterServiceDeleted(ctx context.Context, serviceID uint32) error               // Must be called after a service is deleted
	AfterServiceAccreditationModified(ctx context.Context, serviceID uint32) error // Must be called after a service accreditation is changed
}

// --------------------------------------------------------------------------------------------------------------------

var _ ServicesHooks = MultiServicesHooks{}

// MultiServicesHooks combines multiple services hooks, all hook functions are run in array sequence
type MultiServicesHooks []ServicesHooks

// NewMultiServicesHooks creates a new MultiServicesHooks object
func NewMultiServicesHooks(hooks ...ServicesHooks) MultiServicesHooks {
	return hooks
}

// AfterServiceCreated implements ServicesHooks
func (m MultiServicesHooks) AfterServiceCreated(ctx context.Context, serviceID uint32) error {
	for _, hook := range m {
		if err := hook.AfterServiceCreated(ctx, serviceID); err != nil {
			return err
		}
	}
	return nil
}

// AfterServiceActivated implements ServicesHooks
func (m MultiServicesHooks) AfterServiceActivated(ctx context.Context, serviceID uint32) error {
	for _, hook := range m {
		if err := hook.AfterServiceActivated(ctx, serviceID); err != nil {
			return err
		}
	}
	return nil
}

// AfterServiceDeactivated implements ServicesHooks
func (m MultiServicesHooks) AfterServiceDeactivated(ctx context.Context, serviceID uint32) error {
	for _, hook := range m {
		if err := hook.AfterServiceDeactivated(ctx, serviceID); err != nil {
			return err
		}
	}
	return nil
}

// AfterServiceDeleted implements ServicesHooks
func (m MultiServicesHooks) AfterServiceDeleted(ctx context.Context, serviceID uint32) error {
	for _, hook := range m {
		if err := hook.AfterServiceDeleted(ctx, serviceID); err != nil {
			return err
		}
	}
	return nil
}

// AfterServiceAccreditationModified implements ServicesHooks
func (m MultiServicesHooks) AfterServiceAccreditationModified(ctx context.Context, serviceID uint32) error {
	for _, hook := range m {
		if err := hook.AfterServiceAccreditationModified(ctx, serviceID); err != nil {
			return err
		}
	}
	return nil
}
