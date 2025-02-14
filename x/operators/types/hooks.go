package types

// DONTCOVER

import (
	"context"
)

// Event Hooks
// These can be utilized to communicate between an operators keeper
// and another keeper which must take particular actions when
// operators change state. The second keeper must implement this
// interface, which then the operators keeper can call.

// OperatorsHooks event hooks for operators objects (noalias)
type OperatorsHooks interface {
	AfterOperatorRegistered(ctx context.Context, operatorID uint32) error            // Must be called after an operator is registered
	AfterOperatorInactivatingStarted(ctx context.Context, operatorID uint32) error   // Must be called after an operator has started inactivating
	AfterOperatorInactivatingCompleted(ctx context.Context, operatorID uint32) error // Must be called after an operator has completed inactivating
	AfterOperatorReactivated(ctx context.Context, operatorID uint32) error           // Must be called after an operator has been reactivated
	BeforeOperatorDeleted(ctx context.Context, operatorID uint32) error              // Must be called before an operator is deleted
}

// --------------------------------------------------------------------------------------------------------------------

// MultiOperatorsHooks combines multiple operators hooks, all hook functions are run in array sequence
type MultiOperatorsHooks []OperatorsHooks

var _ OperatorsHooks = &MultiOperatorsHooks{}

// NewMultiOperatorsHooks creates a new MultiOperatorsHooks object
func NewMultiOperatorsHooks(hooks ...OperatorsHooks) MultiOperatorsHooks {
	return hooks
}

// AfterOperatorRegistered implements OperatorsHooks
func (h MultiOperatorsHooks) AfterOperatorRegistered(ctx context.Context, operatorID uint32) error {
	for _, hook := range h {
		if err := hook.AfterOperatorRegistered(ctx, operatorID); err != nil {
			return err
		}
	}
	return nil
}

// AfterOperatorInactivatingStarted implements OperatorsHooks
func (h MultiOperatorsHooks) AfterOperatorInactivatingStarted(ctx context.Context, operatorID uint32) error {
	for _, hook := range h {
		if err := hook.AfterOperatorInactivatingStarted(ctx, operatorID); err != nil {
			return err
		}
	}
	return nil
}

// AfterOperatorInactivatingCompleted implements OperatorsHooks
func (h MultiOperatorsHooks) AfterOperatorInactivatingCompleted(ctx context.Context, operatorID uint32) error {
	for _, hook := range h {
		if err := hook.AfterOperatorInactivatingCompleted(ctx, operatorID); err != nil {
			return err
		}
	}
	return nil
}

// AfterOperatorReactivated implements OperatorsHooks
func (h MultiOperatorsHooks) AfterOperatorReactivated(ctx context.Context, operatorID uint32) error {
	for _, hook := range h {
		if err := hook.AfterOperatorReactivated(ctx, operatorID); err != nil {
			return err
		}
	}
	return nil
}

// BeforeOperatorDeleted implements OperatorsHooks
func (h MultiOperatorsHooks) BeforeOperatorDeleted(ctx context.Context, operatorID uint32) error {
	for _, hook := range h {
		if err := hook.BeforeOperatorDeleted(ctx, operatorID); err != nil {
			return err
		}
	}
	return nil
}
