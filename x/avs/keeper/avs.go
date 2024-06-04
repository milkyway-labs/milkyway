package keeper

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/avs/types"
)

// SetNextAVSID sets the next AVS ID to be used when registering a new AVS
func (k Keeper) SetNextAVSID(ctx sdk.Context, avsID uint32) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.NextAVSIDKey(), types.GetAVSIDBytes(avsID))
}

// HasNextAVSID checks if the next AVS ID is set
func (k Keeper) HasNextAVSID(ctx sdk.Context) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.NextAVSIDKey())
}

// GetNextAVSID returns the next AVS ID to be used when registering a new AVS
func (k Keeper) GetNextAVSID(ctx sdk.Context) (avsID uint32, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextAVSIDKey())
	if bz == nil {
		return 0, errors.Wrapf(types.ErrInvalidGenesis, "initial avs id not set")
	}

	avsID = types.GetAVSIDFromBytes(bz)
	return avsID, nil
}

// --------------------------------------------------------------------------------------------------------------------

// SaveAVS stores a new AVS in the KVStore
func (k Keeper) SaveAVS(ctx sdk.Context, avs types.AVS) {
	previous, existed := k.GetAVS(ctx, avs.ID)

	// Save the AVS data
	store := ctx.KVStore(k.storeKey)
	store.Set(types.AVSStoreKey(avs.ID), k.cdc.MustMarshal(&avs))
	k.Logger(ctx).Debug("saved avs", "id", avs.ID)

	// Call the hook based on the AVS status change
	switch {
	case !existed:
		k.AfterAVSCreated(ctx, avs.ID)
	case previous.Status == types.AVS_STATUS_CREATED && avs.Status == types.AVS_STATUS_REGISTERED:
		k.AfterAVSRegistered(ctx, avs.ID)
	case previous.Status == types.AVS_STATUS_REGISTERED && avs.Status == types.AVS_STATUS_UNREGISTERED:
		k.AfterAVSDeregistered(ctx, avs.ID)
	}
}

// GetAVS returns an AVS from the KVStore
func (k Keeper) GetAVS(ctx sdk.Context, avsID uint32) (avs types.AVS, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.AVSStoreKey(avsID))
	if bz == nil {
		return avs, false
	}

	k.cdc.MustUnmarshal(bz, &avs)
	return avs, true
}
