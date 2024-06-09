package keeper

import (
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.Codec
}

func NewKeeper(storeKey storetypes.StoreKey) *Keeper {
	return &Keeper{
		storeKey: storeKey,
	}
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}
