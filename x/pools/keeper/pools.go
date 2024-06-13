package keeper

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/milkyway-labs/milkyway/x/pools/types"
)

// SetNextPoolID sets the next service ID to be used when registering a new Pool
func (k *Keeper) SetNextPoolID(ctx sdk.Context, serviceID uint32) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.NextPoolIDKey, types.GetPoolIDBytes(serviceID))
}

// GetNextPoolID returns the next service ID to be used when registering a new Pool
func (k *Keeper) GetNextPoolID(ctx sdk.Context) (serviceID uint32, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextPoolIDKey)
	if bz == nil {
		return 0, errors.Wrapf(types.ErrInvalidGenesis, "initial service id not set")
	}

	serviceID = types.GetPoolIDFromBytes(bz)
	return serviceID, nil
}

// --------------------------------------------------------------------------------------------------------------------

// SavePool stores the given pool inside the store
func (k *Keeper) SavePool(ctx sdk.Context, pool types.Pool) error {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetPoolStoreKey(pool.ID), k.cdc.MustMarshal(&pool))

	// Create the pool account if it does not exist
	poolAddress, err := sdk.AccAddressFromBech32(pool.Address)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid pool address: %s", pool.Address)
	}
	k.createAccountIfNotExists(ctx, poolAddress)

	return nil
}

// GetPool retrieves the pool with the given ID from the store.
// If the pool does not exist, false is returned instead
func (k *Keeper) GetPool(ctx sdk.Context, id uint32) (types.Pool, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPoolStoreKey(id))
	if bz == nil {
		return types.Pool{}, false
	}

	var pool types.Pool
	k.cdc.MustUnmarshal(bz, &pool)
	return pool, true
}
