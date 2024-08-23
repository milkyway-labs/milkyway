package keeper

import (
	"context"

	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/assets/types"
)

type Keeper struct {
	cdc          codec.Codec
	storeService corestoretypes.KVStoreService

	Schema        collections.Schema
	Params        collections.Item[types.Params]
	Assets        collections.Map[string, types.Asset]                 // denom => types.Asset
	TickerIndexes collections.KeySet[collections.Pair[string, string]] // ticker + denom => nil

	authority string
}

func NewKeeper(
	cdc codec.Codec,
	storeService corestoretypes.KVStoreService,
	authority string,
) *Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	k := &Keeper{
		cdc:          cdc,
		storeService: storeService,

		Params: collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Assets: collections.NewMap(
			sb,
			types.AssetKeyPrefix,
			"assets",
			collections.StringKey,
			codec.CollValue[types.Asset](cdc),
		),
		TickerIndexes: collections.NewKeySet(
			sb,
			types.TickerIndexKeyPrefix,
			"ticker_indexes",
			collections.PairKeyCodec(collections.StringKey, collections.StringKey),
		),

		authority: authority,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k *Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", "x/"+types.ModuleName)
}
