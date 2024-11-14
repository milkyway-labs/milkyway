package keeper

import (
	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/pools/types"
)

type Keeper struct {
	storeKey     storetypes.StoreKey
	cdc          codec.Codec
	storeService corestoretypes.KVStoreService
	hooks        types.PoolsHooks

	accountKeeper types.AccountKeeper

	// Data
	schema         collections.Schema
	poolAddressSet collections.KeySet[string]
}

func NewKeeper(cdc codec.Codec,
	storeKey storetypes.StoreKey,
	storeService corestoretypes.KVStoreService,
	accountKeeper types.AccountKeeper,
) *Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := &Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		storeService:  storeService,
		accountKeeper: accountKeeper,
		poolAddressSet: collections.NewKeySet(
			sb,
			types.PoolAddressSetPrefix,
			"pool_address_set",
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

// SetHooks allows to set the pools hooks
func (k *Keeper) SetHooks(rs types.PoolsHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set pools hooks twice")
	}

	k.hooks = rs
	return k
}
