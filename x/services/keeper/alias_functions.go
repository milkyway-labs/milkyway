package keeper

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/services/types"
)

// IterateServices iterates over the services in the store and performs a callback function
func (k Keeper) IterateServices(ctx sdk.Context, cb func(service types.AVS) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.AVSPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var service types.AVS
		k.cdc.MustUnmarshal(iterator.Value(), &service)

		if cb(service) {
			break
		}
	}
}
