package types

import (
	"context"
)

type RestakingHooks interface {
	PoolRestakingHooks
	OperatorRestakingHooks
	ServiceRestakingHooks
	AfterUnbondingInitiated(ctx context.Context, unbondingDelegationID uint64) error
	AfterUserPreferencesModified(ctx context.Context, userAddress string, oldPreferences, newPreferences UserPreferences) error
}

// MultiRestakingHooks combines multiple restaking hooks, all hook functions are
// run in array sequence
type MultiRestakingHooks []RestakingHooks

// NewMultiRestakingHooks creates a new MultiRestakingHooks object
func NewMultiRestakingHooks(hooks ...RestakingHooks) MultiRestakingHooks {
	return hooks
}

// BeforePoolDelegationCreated implements RestakingHooks
func (m MultiRestakingHooks) BeforePoolDelegationCreated(ctx context.Context, poolID uint32, delegator string) error {
	for _, hook := range m {
		if err := hook.BeforePoolDelegationCreated(ctx, poolID, delegator); err != nil {
			return err
		}
	}
	return nil
}

// BeforePoolDelegationSharesModified implements RestakingHooks
func (m MultiRestakingHooks) BeforePoolDelegationSharesModified(ctx context.Context, poolID uint32, delegator string) error {
	for _, hook := range m {
		if err := hook.BeforePoolDelegationSharesModified(ctx, poolID, delegator); err != nil {
			return err
		}
	}
	return nil
}

// AfterPoolDelegationModified implements RestakingHooks
func (m MultiRestakingHooks) AfterPoolDelegationModified(ctx context.Context, poolID uint32, delegator string) error {
	for _, hook := range m {
		if err := hook.AfterPoolDelegationModified(ctx, poolID, delegator); err != nil {
			return err
		}
	}
	return nil
}

// BeforePoolDelegationRemoved implements RestakingHooks
func (m MultiRestakingHooks) BeforePoolDelegationRemoved(ctx context.Context, poolID uint32, delegator string) error {
	for _, hook := range m {
		if err := hook.BeforePoolDelegationRemoved(ctx, poolID, delegator); err != nil {
			return err
		}
	}
	return nil
}

// BeforeOperatorDelegationCreated implements RestakingHooks
func (m MultiRestakingHooks) BeforeOperatorDelegationCreated(ctx context.Context, operatorID uint32, delegator string) error {
	for _, hook := range m {
		if err := hook.BeforeOperatorDelegationCreated(ctx, operatorID, delegator); err != nil {
			return err
		}
	}
	return nil
}

// BeforeOperatorDelegationSharesModified implements RestakingHooks
func (m MultiRestakingHooks) BeforeOperatorDelegationSharesModified(ctx context.Context, operatorID uint32, delegator string) error {
	for _, hook := range m {
		if err := hook.BeforeOperatorDelegationSharesModified(ctx, operatorID, delegator); err != nil {
			return err
		}
	}
	return nil
}

// AfterOperatorDelegationModified implements RestakingHooks
func (m MultiRestakingHooks) AfterOperatorDelegationModified(ctx context.Context, operatorID uint32, delegator string) error {
	for _, hook := range m {
		if err := hook.AfterOperatorDelegationModified(ctx, operatorID, delegator); err != nil {
			return err
		}
	}
	return nil
}

// BeforeOperatorDelegationRemoved implements RestakingHooks
func (m MultiRestakingHooks) BeforeOperatorDelegationRemoved(ctx context.Context, operatorID uint32, delegator string) error {
	for _, hook := range m {
		if err := hook.BeforeOperatorDelegationRemoved(ctx, operatorID, delegator); err != nil {
			return err
		}
	}
	return nil
}

// BeforeServiceDelegationCreated implements RestakingHooks
func (m MultiRestakingHooks) BeforeServiceDelegationCreated(ctx context.Context, serviceID uint32, delegator string) error {
	for _, hook := range m {
		if err := hook.BeforeServiceDelegationCreated(ctx, serviceID, delegator); err != nil {
			return err
		}
	}
	return nil
}

// BeforeServiceDelegationSharesModified implements RestakingHooks
func (m MultiRestakingHooks) BeforeServiceDelegationSharesModified(ctx context.Context, serviceID uint32, delegator string) error {
	for _, hook := range m {
		if err := hook.BeforeServiceDelegationSharesModified(ctx, serviceID, delegator); err != nil {
			return err
		}
	}
	return nil
}

// AfterServiceDelegationModified implements RestakingHooks
func (m MultiRestakingHooks) AfterServiceDelegationModified(ctx context.Context, serviceID uint32, delegator string) error {
	for _, hook := range m {
		if err := hook.AfterServiceDelegationModified(ctx, serviceID, delegator); err != nil {
			return err
		}
	}
	return nil
}

// BeforeServiceDelegationRemoved implements RestakingHooks
func (m MultiRestakingHooks) BeforeServiceDelegationRemoved(ctx context.Context, serviceID uint32, delegator string) error {
	for _, hook := range m {
		if err := hook.BeforeServiceDelegationRemoved(ctx, serviceID, delegator); err != nil {
			return err
		}
	}
	return nil
}

// AfterUnbondingInitiated implements RestakingHooks
func (m MultiRestakingHooks) AfterUnbondingInitiated(ctx context.Context, unbondingDelegationID uint64) error {
	for _, hook := range m {
		if err := hook.AfterUnbondingInitiated(ctx, unbondingDelegationID); err != nil {
			return err
		}
	}
	return nil
}

// AfterUserPreferencesModified implements RestakingHooks
func (m MultiRestakingHooks) AfterUserPreferencesModified(ctx context.Context, userAddress string, oldPreferences, newPreferences UserPreferences) error {
	for _, hook := range m {
		if err := hook.AfterUserPreferencesModified(ctx, userAddress, oldPreferences, newPreferences); err != nil {
			return err
		}
	}
	return nil
}

type PoolRestakingHooks interface {
	BeforePoolDelegationCreated(ctx context.Context, poolID uint32, delegator string) error
	BeforePoolDelegationSharesModified(ctx context.Context, poolID uint32, delegator string) error
	AfterPoolDelegationModified(ctx context.Context, poolID uint32, delegator string) error
	BeforePoolDelegationRemoved(ctx context.Context, poolID uint32, delegator string) error
}

type OperatorRestakingHooks interface {
	BeforeOperatorDelegationCreated(ctx context.Context, operatorID uint32, delegator string) error
	BeforeOperatorDelegationSharesModified(ctx context.Context, operatorID uint32, delegator string) error
	AfterOperatorDelegationModified(ctx context.Context, operatorID uint32, delegator string) error
	BeforeOperatorDelegationRemoved(ctx context.Context, operatorID uint32, delegator string) error
}

type ServiceRestakingHooks interface {
	BeforeServiceDelegationCreated(ctx context.Context, serviceID uint32, delegator string) error
	BeforeServiceDelegationSharesModified(ctx context.Context, serviceID uint32, delegator string) error
	AfterServiceDelegationModified(ctx context.Context, serviceID uint32, delegator string) error
	BeforeServiceDelegationRemoved(ctx context.Context, serviceID uint32, delegator string) error
}
