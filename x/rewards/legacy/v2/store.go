package v2

import (
	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/rewards/types"
)

// MigrateStore performs in-place store migrations from v1 to v2
// The things done here are the following:
// 1. setting up the module params
func MigrateStore(ctx sdk.Context, params collections.Item[types.Params]) error {
	// Set the module params
	return params.Set(ctx, types.DefaultParams())
}
