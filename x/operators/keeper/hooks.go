package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

var _ types.OperatorsHooks = &Keeper{}

// AfterOperatorRegistered implements OperatorsHooks
func (k *Keeper) AfterOperatorRegistered(ctx sdk.Context, operatorID uint32) {
	if k.hooks != nil {
		k.hooks.AfterOperatorRegistered(ctx, operatorID)
	}
}

// AfterOperatorInactivatingStarted implements OperatorsHooks
func (k *Keeper) AfterOperatorInactivatingStarted(ctx sdk.Context, operatorID uint32) {
	if k.hooks != nil {
		k.hooks.AfterOperatorInactivatingStarted(ctx, operatorID)
	}
}

// AfterOperatorInactivatingCompleted implements OperatorsHooks
func (k *Keeper) AfterOperatorInactivatingCompleted(ctx sdk.Context, operatorID uint32) {
	if k.hooks != nil {
		k.hooks.AfterOperatorInactivatingCompleted(ctx, operatorID)
	}
}
