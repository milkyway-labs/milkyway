package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

	icacallbackskeeper "github.com/milkyway-labs/milk/x/icacallbacks/keeper"

	"github.com/milkyway-labs/milk/x/records/types"
)

type (
	Keeper struct {
		// *cosmosibckeeper.Keeper
		Cdc                codec.BinaryCodec
		storeKey           storetypes.StoreKey
		memKey             storetypes.StoreKey
		AccountKeeper      types.AccountKeeper
		TransferKeeper     ibctransferkeeper.Keeper
		IBCKeeper          ibckeeper.Keeper
		ICACallbacksKeeper icacallbackskeeper.Keeper
		params             collections.Item[types.Params]
	}
)

func NewKeeper(
	Cdc codec.BinaryCodec,
	storeService corestoretypes.KVStoreService,
	storeKey,
	memKey storetypes.StoreKey,
	AccountKeeper types.AccountKeeper,
	TransferKeeper ibctransferkeeper.Keeper,
	ibcKeeper ibckeeper.Keeper,
	ICACallbacksKeeper icacallbackskeeper.Keeper,
) *Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	return &Keeper{
		Cdc:                Cdc,
		storeKey:           storeKey,
		memKey:             memKey,
		AccountKeeper:      AccountKeeper,
		TransferKeeper:     TransferKeeper,
		IBCKeeper:          ibcKeeper,
		ICACallbacksKeeper: ICACallbacksKeeper,
		params:             collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](Cdc)),
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
