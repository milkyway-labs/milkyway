package keeper

import (
	"errors"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
)

func (k *Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	err := params.Validate()
	if err != nil {
		return err
	}

	return k.params.Set(ctx, params)
}

func (k *Keeper) GetParams(ctx sdk.Context) (types.Params, error) {
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
