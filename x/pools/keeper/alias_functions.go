package keeper

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/pools/types"
)

// IteratePools iterates over the pools in the store and performs a callback function
func (k *Keeper) IteratePools(ctx sdk.Context, cb func(pool types.Pool) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.PoolPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var pool types.Pool
		k.cdc.MustUnmarshal(iterator.Value(), &pool)

		if cb(pool) {
			break
		}
	}
}

// GetPoolByDenom returns the pool for the given denom if it exists.
// If the pool does not exist, false is returned instead
func (k *Keeper) GetPoolByDenom(ctx sdk.Context, denom string) (types.Pool, bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.PoolPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var pool types.Pool
		k.cdc.MustUnmarshal(iterator.Value(), &pool)

		if pool.Denom == denom {
			return pool, true
		}
	}

	return types.Pool{}, false
}

// CreateOrGetPoolByDenom creates a new pool for the given denom if it does not exist.
// If the pool already exists, no action is taken.
// In both cases, the pool is returned.
func (k *Keeper) CreateOrGetPoolByDenom(ctx sdk.Context, denom string) (types.Pool, error) {
	// If the pool already exists, just return
	if pool, found := k.GetPoolByDenom(ctx, denom); found {
		return pool, nil
	}

	// Get the pool id
	poolID, err := k.GetNextPoolID(ctx)
	if err != nil {
		return types.Pool{}, err
	}

	// Create the pool and validate it
	pool := types.NewPool(poolID, denom)
	err = pool.Validate()
	if err != nil {
		return types.Pool{}, err
	}

	// Save the pool
	k.SavePool(ctx, pool)

	// Increment the pool id
	k.SetNextPoolID(ctx, poolID+1)

	// Log the event
	k.Logger(ctx).Debug("created pool", "id", poolID, "denom", denom)

	return pool, nil
}
