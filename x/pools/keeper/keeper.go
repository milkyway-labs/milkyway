package keeper

import (
	"context"

	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v3/x/pools/types"
)

type Keeper struct {
	cdc          codec.Codec
	storeService corestoretypes.KVStoreService
	hooks        types.PoolsHooks

	accountKeeper types.AccountKeeper

	schema         collections.Schema
	nextPoolID     collections.Sequence                // Sequence for pool IDs
	pools          collections.Map[uint32, types.Pool] // Map of pool ID to pool
	poolAddressSet collections.KeySet[string]          // Set of pool addresses
}

func NewKeeper(cdc codec.Codec,
	storeService corestoretypes.KVStoreService,
	accountKeeper types.AccountKeeper,
) *Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := &Keeper{
		cdc:           cdc,
		storeService:  storeService,
		accountKeeper: accountKeeper,

		nextPoolID: collections.NewSequence(
			sb,
			types.NextPoolIDKey,
			"next_pool_id",
		),
		pools: collections.NewMap(
			sb,
			types.PoolPrefix,
			"pools",
			collections.Uint32Key,
			codec.CollValue[types.Pool](cdc),
		),
		poolAddressSet: collections.NewKeySet(
			sb,
			types.PoolAddressSetPrefix,
			"pools_addresses_set",
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
func (k *Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", "x/"+types.ModuleName)
}

// SetHooks allows to set the pools hooks
func (k *Keeper) SetHooks(rs types.PoolsHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set pools hooks twice")
	}

	k.hooks = rs
	return k
}
