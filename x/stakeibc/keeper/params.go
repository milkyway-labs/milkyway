package keeper

import (
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milk/x/stakeibc/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	params, err := k.params.Get(ctx)
	if err != nil {
		panic(err) // XXX
	}
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	if err := k.params.Set(ctx, params); err != nil {
		panic(err) // XXX
	}
}

func (k *Keeper) GetParam(ctx sdk.Context, key []byte) uint64 {
	params, err := k.params.Get(ctx)
	if err != nil {
		panic(err) // XXX
	}
	// XXX
	f := reflect.ValueOf(params).FieldByName(string(key))
	return f.Uint()
}
