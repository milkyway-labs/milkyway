package keeper

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

// IterateOperators iterates over the operators in the store and performs a callback function
func (k *Keeper) IterateOperators(ctx sdk.Context, cb func(operator types.Operator) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.OperatorPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var operator types.Operator
		k.cdc.MustUnmarshal(iterator.Value(), &operator)

		if cb(operator) {
			break
		}
	}
}
