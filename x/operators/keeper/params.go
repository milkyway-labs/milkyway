package keeper

import (
	"context"

	"github.com/milkyway-labs/milkyway/v12/x/operators/types"
)

// SetParams sets module parameters
func (k *Keeper) SetParams(ctx context.Context, params types.Params) error {
	return k.params.Set(ctx, params)
}

// GetParams returns the module parameters
func (k *Keeper) GetParams(ctx context.Context) (types.Params, error) {
	return k.params.Get(ctx)
}
