package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/milkyway-labs/milkyway/v9/x/pools/types"
)

// SetNextPoolID sets the next pool ID to be used when registering a new Pool
func (k *Keeper) SetNextPoolID(ctx context.Context, poolID uint32) error {
	return k.nextPoolID.Set(ctx, uint64(poolID))
}

// GetNextPoolID returns the next pool ID to be used when registering a new Pool
func (k *Keeper) GetNextPoolID(ctx context.Context) (poolID uint32, err error) {
	nextPoolID, err := k.nextPoolID.Next(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get next pool ID")
	}

	// If the next pool ID is 0, we need to increment it to 1
	if nextPoolID == 0 {
		return k.GetNextPoolID(ctx)
	}

	return uint32(nextPoolID), nil
}

// --------------------------------------------------------------------------------------------------------------------

// SavePool stores the given pool inside the store
func (k *Keeper) SavePool(ctx context.Context, pool types.Pool) error {
	err := k.pools.Set(ctx, pool.ID, pool)
	if err != nil {
		return err
	}

	// Create the pool account if it does not exist
	poolAddress, err := sdk.AccAddressFromBech32(pool.Address)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid pool address: %s", pool.Address)
	}
	k.createAccountIfNotExists(ctx, poolAddress)

	return k.poolAddressSet.Set(ctx, pool.Address)
}

// HasPool checks if a pool with the given ID exists in the store
func (k *Keeper) HasPool(ctx context.Context, id uint32) (bool, error) {
	return k.pools.Has(ctx, id)
}

// GetPool retrieves the pool with the given ID from the store.
// If the pool does not exist, false is returned instead
func (k *Keeper) GetPool(ctx context.Context, id uint32) (types.Pool, error) {
	return k.pools.Get(ctx, id)
}
