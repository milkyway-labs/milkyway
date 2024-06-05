package types

// DONTCOVER

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Event Hooks
// These can be utilized to communicate between an operators keeper
// and another keeper which must take particular actions when
// services change state. The second keeper must implement this
// interface, which then the operators keeper can call.

// OperatorsHooks event hooks for operators objects (noalias)
type OperatorsHooks interface {
	AfterOperatorRegistered(ctx sdk.Context, operatorID uint32)            // Must be called after an operator is registered
	AfterOperatorInactivatingStarted(ctx sdk.Context, operatorID uint32)   // Must be called after an operator has started inactivating
	AfterOperatorInactivatingCompleted(ctx sdk.Context, operatorID uint32) // Must be called after an operator has completed inactivating
}

// --------------------------------------------------------------------------------------------------------------------

// MultiOperatorsHooks combines multiple operators hooks, all hook functions are run in array sequence
type MultiOperatorsHooks []OperatorsHooks

// NewMultiOperatorsHooks creates a new MultiOperatorsHooks object
func NewMultiOperatorsHooks(hooks ...OperatorsHooks) MultiOperatorsHooks {
	return hooks
}

// AfterOperatorRegistered implements OperatorsHooks
func (h MultiOperatorsHooks) AfterOperatorRegistered(ctx sdk.Context, operatorID uint32) {
	for _, hook := range h {
		hook.AfterOperatorRegistered(ctx, operatorID)
	}
}

// AfterOperatorInactivatingStarted implements OperatorsHooks
func (h MultiOperatorsHooks) AfterOperatorInactivatingStarted(ctx sdk.Context, operatorID uint32) {
	for _, hook := range h {
		hook.AfterOperatorInactivatingStarted(ctx, operatorID)
	}
}

// AfterOperatorInactivatingCompleted implements OperatorsHooks
func (h MultiOperatorsHooks) AfterOperatorInactivatingCompleted(ctx sdk.Context, operatorID uint32) {
	for _, hook := range h {
		hook.AfterOperatorInactivatingCompleted(ctx, operatorID)
	}
}
