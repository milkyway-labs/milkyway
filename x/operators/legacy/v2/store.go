package v2

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

// MigrateStore performs in-place store migrations from v1 to v2
// The things done here are the following:
// 1. setting up the module params
func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)

	// Set the module params
	params := types.DefaultParams()
	store.Set(types.ParamsKey, cdc.MustMarshal(&params))

	return nil
}
