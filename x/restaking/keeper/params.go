package keeper

import (
	"slices"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

// UnbondingTime returns the unbonding time.
func (k *Keeper) UnbondingTime(ctx sdk.Context) time.Duration {
	return k.GetParams(ctx).UnbondingTime
}

// SetRestakableDenoms sets the denoms that are allowed to be restaked.
func (k *Keeper) SetRestakableDenoms(ctx sdk.Context, denoms []string) {
	params := k.GetParams(ctx)
	params.AllowedDenoms = denoms
	k.SetParams(ctx, params)
}

// GetAllowedDenoms gets the list of denoms that are allowed to be restaked.
// If the list is empty, all denoms are allowed.
func (k *Keeper) GetAllowedDenoms(ctx sdk.Context) []string {
	return k.GetParams(ctx).AllowedDenoms
}

// IsDenomRestakable checks if the asset with the provided denom is allowed
// to be restaked.
func (k *Keeper) IsDenomRestakable(ctx sdk.Context, denom string) bool {
	allowedDenoms := k.GetParams(ctx).AllowedDenoms
	if len(allowedDenoms) == 0 {
		return true
	}

	return slices.Contains(allowedDenoms, denom)
}

// SetParams sets module parameters
func (k *Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)
}

// GetParams returns the module parameters
func (k *Keeper) GetParams(ctx sdk.Context) (p types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return types.DefaultParams()
	}
	k.cdc.MustUnmarshal(bz, &p)
	return p
}
