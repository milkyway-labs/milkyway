package keeper

import (
	"context"
	"slices"
	"time"

	"github.com/milkyway-labs/milkyway/v3/x/restaking/types"
)

// UnbondingTime returns the unbonding time.
func (k *Keeper) UnbondingTime(ctx context.Context) (time.Duration, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return 0, err
	}

	return params.UnbondingTime, nil
}

// SetRestakableDenoms sets the denoms that are allowed to be restaked.
func (k *Keeper) SetRestakableDenoms(ctx context.Context, denoms []string) error {
	params, err := k.GetParams(ctx)
	if err != nil {
		return err
	}

	params.AllowedDenoms = denoms
	return k.SetParams(ctx, params)
}

// GetRestakableDenoms gets the list of denoms that are allowed to be restaked.
// If the list is empty, all denoms are allowed.
func (k *Keeper) GetRestakableDenoms(ctx context.Context) ([]string, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	return params.AllowedDenoms, nil
}

// IsDenomRestakable checks if the asset with the provided denom is allowed
// to be restaked.
func (k *Keeper) IsDenomRestakable(ctx context.Context, denom string) (bool, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return false, err
	}

	return len(params.AllowedDenoms) == 0 || slices.Contains(params.AllowedDenoms, denom), nil
}

// SetParams sets module parameters
func (k *Keeper) SetParams(ctx context.Context, params types.Params) error {
	store := k.storeService.OpenKVStore(ctx)
	bz := k.cdc.MustMarshal(&params)
	return store.Set(types.ParamsKey, bz)
}

// GetParams returns the module parameters
func (k *Keeper) GetParams(ctx context.Context) (p types.Params, err error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.ParamsKey)
	if err != nil {
		return p, err
	}

	if bz == nil {
		return types.DefaultParams(), nil
	}

	k.cdc.MustUnmarshal(bz, &p)
	return p, nil
}
