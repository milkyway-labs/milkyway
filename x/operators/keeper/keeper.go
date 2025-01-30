package keeper

import (
	"context"

	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v8/x/operators/types"
)

type Keeper struct {
	storeService corestoretypes.KVStoreService
	cdc          codec.BinaryCodec
	hooks        types.OperatorsHooks

	accountKeeper types.AccountKeeper
	poolKeeper    types.CommunityPoolKeeper
	Schema        collections.Schema

	nextOperatorID     collections.Sequence                          // Next operator ID
	operators          collections.Map[uint32, types.Operator]       // operator ID -> operator
	operatorAddressSet collections.KeySet[string]                    // Set of operator addresses
	operatorParams     collections.Map[uint32, types.OperatorParams] // operator ID -> parameters
	params             collections.Item[types.Params]                // global parameters

	// authority represents the address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority string
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService corestoretypes.KVStoreService,
	accountKeeper types.AccountKeeper,
	poolKeeper types.CommunityPoolKeeper,
	authority string,
) *Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := &Keeper{
		storeService:  storeService,
		cdc:           cdc,
		authority:     authority,
		accountKeeper: accountKeeper,
		poolKeeper:    poolKeeper,

		nextOperatorID: collections.NewSequence(
			sb,
			types.NextOperatorIDKey,
			"next_operator_id",
		),
		operators: collections.NewMap(
			sb,
			types.OperatorPrefix,
			"operators",
			collections.Uint32Key,
			codec.CollValue[types.Operator](cdc),
		),
		operatorAddressSet: collections.NewKeySet(
			sb,
			types.OperatorAddressSetPrefix,
			"operators_addresses",
			collections.StringKey,
		),
		operatorParams: collections.NewMap(
			sb,
			types.OperatorParamsMapPrefix,
			"operators_params",
			collections.Uint32Key,
			codec.CollValue[types.OperatorParams](cdc),
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

// SetHooks allows to set the operators hooks
func (k *Keeper) SetHooks(rs types.OperatorsHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set operators hooks twice")
	}

	k.hooks = rs
	return k
}
