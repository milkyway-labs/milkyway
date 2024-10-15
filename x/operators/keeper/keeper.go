package keeper

import (
	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

type Keeper struct {
	storeKey     storetypes.StoreKey
	storeService corestoretypes.KVStoreService
	cdc          codec.BinaryCodec
	hooks        types.OperatorsHooks

	accountKeeper types.AccountKeeper
	poolKeeper    types.CommunityPoolKeeper
	schema        collections.Schema
	// Index to check if an address is an operator
	operatorAddressSet collections.KeySet[string]
	operatorParams     collections.Map[uint32, types.OperatorParams]

	authority string
}

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
		storeService:  storeService,
		cdc:           cdc,
		authority:     authority,
		accountKeeper: accountKeeper,
		poolKeeper:    poolKeeper,
		operatorAddressSet: collections.NewKeySet(
			sb,
			types.OperatorAddressSetPrefix,
			"operators_address",
			collections.StringKey,
		),
		operatorParams: collections.NewMap(
			sb,
			types.OperatorParamsMapPrefix,
			"operator_params",
			collections.Uint32Key,
			codec.CollValue[types.OperatorParams](cdc),
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

// SetHooks allows to set the operators hooks
func (k *Keeper) SetHooks(rs types.OperatorsHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set avs hooks twice")
	}

	k.hooks = rs
	return k
}
