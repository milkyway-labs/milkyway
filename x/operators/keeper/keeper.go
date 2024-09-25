package keeper

import (
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.Codec
	hooks    types.OperatorsHooks

	accountKeeper types.AccountKeeper
	poolKeeper    types.CommunityPoolKeeper

	// Msg server router
	router *baseapp.MsgServiceRouter

	authority string
}

func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	accountKeeper types.AccountKeeper,
	poolKeeper types.CommunityPoolKeeper,
	router *baseapp.MsgServiceRouter,
	authority string,
) *Keeper {
	return &Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		authority:     authority,
		accountKeeper: accountKeeper,
		poolKeeper:    poolKeeper,
		router:        router,
	}
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// Router returns the gov keeper's router
func (k Keeper) Router() *baseapp.MsgServiceRouter {
	return k.router
}

// SetHooks allows to set the operators hooks
func (k *Keeper) SetHooks(rs types.OperatorsHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set avs hooks twice")
	}

	k.hooks = rs
	return k
}
