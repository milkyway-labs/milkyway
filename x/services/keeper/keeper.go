package keeper

import (
	"context"

	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v2/x/services/types"
)

type Keeper struct {
	cdc   codec.BinaryCodec
	hooks types.ServicesHooks

	accountKeeper types.AccountKeeper
	poolKeeper    types.CommunityPoolKeeper

	storeService corestoretypes.KVStoreService
	Schema       collections.Schema

	nextServiceID     collections.Sequence                         // Next service ID
	services          collections.Map[uint32, types.Service]       // service ID -> service
	serviceAddressSet collections.KeySet[string]                   // Set of service addresses
	serviceParams     collections.Map[uint32, types.ServiceParams] // service ID -> parameters
	params            collections.Item[types.Params]

	// authority represents the address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority string
}

// NewKeeper creates a new keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService corestoretypes.KVStoreService,
	accountKeeper types.AccountKeeper,
	poolKeeper types.CommunityPoolKeeper,
	authority string,
) *Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := &Keeper{
		cdc:           cdc,
		accountKeeper: accountKeeper,
		poolKeeper:    poolKeeper,
		authority:     authority,
		storeService:  storeService,

		nextServiceID: collections.NewSequence(
			sb,
			types.NextServiceIDKey,
			"next_service_id",
		),
		services: collections.NewMap(
			sb,
			types.ServicePrefix,
			"services",
			collections.Uint32Key,
			codec.CollValue[types.Service](cdc),
		),
		serviceAddressSet: collections.NewKeySet(
			sb,
			types.ServiceAddressSetPrefix,
			"services_address_set",
			collections.StringKey,
		),
		serviceParams: collections.NewMap(
			sb,
			types.ServiceParamsPrefix,
			"services_params",
			collections.Uint32Key,
			codec.CollValue[types.ServiceParams](cdc),
		),
		params: collections.NewItem(
			sb,
			types.ParamsKey,
			"params",
			codec.CollValue[types.Params](cdc),
		),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", "x/"+types.ModuleName)
}

// SetHooks allows to set the reactions hooks
func (k *Keeper) SetHooks(rs types.ServicesHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set services hooks twice")
	}

	k.hooks = rs
	return k
}
