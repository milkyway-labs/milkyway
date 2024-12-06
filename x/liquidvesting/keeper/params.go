package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"

	"github.com/milkyway-labs/milkyway/v3/x/liquidvesting/types"
)

func (k *Keeper) SetParams(ctx context.Context, params types.Params) error {
	err := params.Validate()
	if err != nil {
		return err
	}

	return k.params.Set(ctx, params)
}

func (k *Keeper) GetParams(ctx context.Context) (types.Params, error) {
	params, err := k.params.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return types.DefaultParams(), nil
		} else {
			return types.Params{}, err
		}
	}
	return params, nil
}
