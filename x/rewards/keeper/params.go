package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"

	"github.com/milkyway-labs/milkyway/v4/x/rewards/types"
)

// SetParams sets module parameters
func (k *Keeper) SetParams(ctx context.Context, params types.Params) error {
	return k.Params.Set(ctx, params)
}

// GetParams returns the module parameters
func (k *Keeper) GetParams(ctx context.Context) (p types.Params, err error) {
	p, err = k.Params.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return types.DefaultParams(), nil
		}
		return types.Params{}, err
	}
	return p, nil
}
