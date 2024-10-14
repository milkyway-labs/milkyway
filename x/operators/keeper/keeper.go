package keeper

import (
	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

type Keeper struct {
	storeKey     storetypes.StoreKey
	storeService corestoretypes.KVStoreService
	cdc          codec.Codec
	hooks        types.OperatorsHooks

	accountKeeper types.AccountKeeper
	poolKeeper    types.CommunityPoolKeeper
	schema        collections.Schema
	// Index to check if an address is an operator
	operatorAddressSet collections.KeySet[string]

	// Msg server router
	router *baseapp.MsgServiceRouter

	authority string
}

func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	storeService corestoretypes.KVStoreService,
	accountKeeper types.AccountKeeper,
	poolKeeper types.CommunityPoolKeeper,
	router *baseapp.MsgServiceRouter,
	authority string,
) *Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := &Keeper{
		storeKey:      storeKey,
		storeService:  storeService,
		cdc:           cdc,
		authority:     authority,
		accountKeeper: accountKeeper,
		poolKeeper:    poolKeeper,
		router:        router,
		operatorAddressSet: collections.NewKeySet(
			sb,
			types.OperatorAddressSetPrefix,
			"operators_address",
			collections.StringKey,
		),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.schema = schema

	return k
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
