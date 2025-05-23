package keeper

import (
	"context"

	"github.com/milkyway-labs/milkyway/v12/x/operators/types"
)

var _ types.OperatorsHooks = &Keeper{}

// AfterOperatorRegistered implements OperatorsHooks
func (k *Keeper) AfterOperatorRegistered(ctx context.Context, operatorID uint32) error {
	if k.hooks != nil {
		return k.hooks.AfterOperatorRegistered(ctx, operatorID)
	}
	return nil
}

// AfterOperatorInactivatingStarted implements OperatorsHooks
func (k *Keeper) AfterOperatorInactivatingStarted(ctx context.Context, operatorID uint32) error {
	if k.hooks != nil {
		return k.hooks.AfterOperatorInactivatingStarted(ctx, operatorID)
	}
	return nil
}

// AfterOperatorInactivatingCompleted implements OperatorsHooks
func (k *Keeper) AfterOperatorInactivatingCompleted(ctx context.Context, operatorID uint32) error {
	if k.hooks != nil {
		return k.hooks.AfterOperatorInactivatingCompleted(ctx, operatorID)
	}
	return nil
}

// AfterOperatorReactivated implements OperatorsHooks
func (k *Keeper) AfterOperatorReactivated(ctx context.Context, operatorID uint32) error {
	if k.hooks != nil {
		return k.hooks.AfterOperatorReactivated(ctx, operatorID)
	}
	return nil
}

// BeforeOperatorDeleted implements OperatorsHooks
func (k *Keeper) BeforeOperatorDeleted(ctx context.Context, operatorID uint32) error {
	if k.hooks != nil {
		return k.hooks.BeforeOperatorDeleted(ctx, operatorID)
	}
	return nil
}
