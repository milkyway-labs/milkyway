package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v12/x/pools/types"
)

// createAccountIfNotExists creates an account if it does not exist
func (k *Keeper) createAccountIfNotExists(ctx context.Context, address sdk.AccAddress) {
	if !k.accountKeeper.HasAccount(ctx, address) {
		defer telemetry.IncrCounter(1, "new", "account")
		k.accountKeeper.SetAccount(ctx, k.accountKeeper.NewAccountWithAddress(ctx, address))
	}
}

// IteratePools iterates over the pools in the store and performs a callback function
func (k *Keeper) IteratePools(ctx context.Context, cb func(pool types.Pool) (stop bool, err error)) error {
	err := k.pools.Walk(ctx, nil, func(_ uint32, pool types.Pool) (stop bool, err error) {
		return cb(pool)
	})
	return err
}

// GetPools returns the list of stored pools
func (k *Keeper) GetPools(ctx context.Context) ([]types.Pool, error) {
	var pools []types.Pool
	err := k.IteratePools(ctx, func(pool types.Pool) (stop bool, err error) {
		pools = append(pools, pool)
		return false, nil
	})
	return pools, err
}

// GetPoolByDenom returns the pool for the given denom if it exists.
// If the pool does not exist, false is returned instead
func (k *Keeper) GetPoolByDenom(ctx context.Context, denom string) (types.Pool, bool, error) {
	var poolFound types.Pool
	err := k.pools.Walk(ctx, nil, func(_ uint32, pool types.Pool) (stop bool, err error) {
		if pool.Denom == denom {
			poolFound = pool
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return types.Pool{}, false, err
	}
	if poolFound != (types.Pool{}) {
		return poolFound, true, nil
	}
	return types.Pool{}, false, nil
}

// CreateOrGetPoolByDenom creates a new pool for the given denom if it does not exist.
// If the pool already exists, no action is taken.
// In both cases, the pool is returned.
func (k *Keeper) CreateOrGetPoolByDenom(ctx context.Context, denom string) (types.Pool, error) {
	// If the pool already exists, just return
	pool, found, err := k.GetPoolByDenom(ctx, denom)
	if err != nil {
		return types.Pool{}, err
	}

	if found {
		return pool, nil
	}

	// Get the pool id
	poolID, err := k.GetNextPoolID(ctx)
	if err != nil {
		return types.Pool{}, err
	}

	// Create the pool and validate it
	pool = types.NewPool(poolID, denom)
	err = pool.Validate()
	if err != nil {
		return types.Pool{}, err
	}

	// Save the pool
	err = k.SavePool(ctx, pool)
	if err != nil {
		return types.Pool{}, err
	}

	// Increment the pool id
	err = k.SetNextPoolID(ctx, poolID+1)
	if err != nil {
		return types.Pool{}, err
	}

	// Log the event
	k.Logger(ctx).Debug("created pool", "id", poolID, "denom", denom)

	// Call the hook
	err = k.AfterPoolCreated(ctx, pool.ID)
	if err != nil {
		return pool, err
	}

	return pool, nil
}

// IsPoolAddress returns true if the provided address is the address
// where the users' asset are kept when they perform a pool restaking.
func (k *Keeper) IsPoolAddress(ctx context.Context, address string) (bool, error) {
	return k.poolAddressSet.Has(ctx, address)
}
