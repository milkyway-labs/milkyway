package keeper

import (
	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
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

	accountKeeper types.AccountKeeper
	poolKeeper    types.CommunityPoolKeeper

	// Data
	storeService      corestoretypes.KVStoreService
	schema            collections.Schema
	serviceAddressSet collections.KeySet[string]
	serviceParams     collections.Map[uint32, types.ServiceParams]

	// authority represents the address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority string
}

// NewKeeper creates a new keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	storeService corestoretypes.KVStoreService,
	accountKeeper types.AccountKeeper,
	poolKeeper types.CommunityPoolKeeper,
	authority string,
) *Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := &Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		accountKeeper: accountKeeper,
		poolKeeper:    poolKeeper,
		authority:     authority,
		storeService:  storeService,
		serviceAddressSet: collections.NewKeySet(
			sb,
			types.ServiceAddressSetPrefix,
			"service_address_set",
			collections.StringKey,
		),
		serviceParams: collections.NewMap(
			sb,
			types.ServiceParamsPrefix,
			"service_params",
			collections.Uint32Key,
			codec.CollValue[types.ServiceParams](cdc),
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

// SetHooks allows to set the reactions hooks
func (k *Keeper) SetHooks(rs types.ServicesHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set services hooks twice")
	}

	k.hooks = rs
	return k
}
