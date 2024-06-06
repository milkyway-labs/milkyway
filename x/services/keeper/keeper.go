package keeper

import (
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/services/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec
	hooks    types.ServicesHooks

	poolKeeper CommunityPoolKeeper
}

// NewKeeper creates a new keeper
func NewKeeper(cdc codec.BinaryCodec, storeKey storetypes.StoreKey, poolKeeper CommunityPoolKeeper) *Keeper {
	return &Keeper{
		storeKey:   storeKey,
		cdc:        cdc,
		poolKeeper: poolKeeper,
	}
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// SetHooks allows to set the reactions hooks
func (k *Keeper) SetHooks(rs types.ServicesHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set avs hooks twice")
	}

	k.hooks = rs
	return k
}
