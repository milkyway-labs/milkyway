package keeper

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/collections/indexes"
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.Codec

	authority string

	accountKeeper   types.AccountKeeper
	bankKeeper      types.BankKeeper
	poolsKeeper     types.PoolsKeeper
	operatorsKeeper types.OperatorsKeeper
	servicesKeeper  types.ServicesKeeper

	// Keeper data
	schema collections.Schema

	// Here we use a IndexMap with NoValue instead of a KeySet because the cosmos-sdk don't
	// provide a KeySet with indexes that we need in order to get the list of operators
	// that have joined a service given a serviceID.
	operatorJoinedServices *collections.IndexedMap[collections.Pair[uint32, uint32], collections.NoValue, operatorServiceIndex]

	// The pair represents the service ID and the operator ID
	serviceOperatorsAllowList collections.KeySet[collections.Pair[uint32, uint32]]

	// The pair represents the service ID and the pool ID
	serviceSecuringPools collections.KeySet[collections.Pair[uint32, uint32]]

	hooks types.RestakingHooks
}

func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	storeService corestoretypes.KVStoreService,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	poolsKeeper types.PoolsKeeper,
	operatorsKeeper types.OperatorsKeeper,
	servicesKeeper types.ServicesKeeper,
	authority string,
) *Keeper {
	// Ensure that authority is a valid AccAddress
	if _, err := accountKeeper.AddressCodec().StringToBytes(authority); err != nil {
		panic("authority is not a valid account address")
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := &Keeper{
		storeKey: storeKey,
		cdc:      cdc,

		accountKeeper:   accountKeeper,
		bankKeeper:      bankKeeper,
		poolsKeeper:     poolsKeeper,
		operatorsKeeper: operatorsKeeper,
		servicesKeeper:  servicesKeeper,

		operatorJoinedServices: collections.NewIndexedMap(
			sb, types.OperatorJoinedServicesPrefix,
			"operator_joined_services",
			collections.PairKeyCodec(collections.Uint32Key, collections.Uint32Key),
			collections.NoValue{},
			newOperatorServiceIndex(sb),
		),
		serviceOperatorsAllowList: collections.NewKeySet(
			sb, types.ServiceOperatorsAllowListPrefix,
			"service_operators_allow_list",
			collections.PairKeyCodec(collections.Uint32Key, collections.Uint32Key),
		),
		serviceSecuringPools: collections.NewKeySet(
			sb, types.ServiceSecuringPoolsPrefix,
			"service_securing_pools",
			collections.PairKeyCodec(collections.Uint32Key, collections.Uint32Key),
		),
		authority: authority,
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
func (k *Keeper) SetHooks(rs types.RestakingHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set services hooks twice")
	}

	k.hooks = rs
	return k
}

// ------------------------------------------------------------------------------

type operatorServiceIndex struct {
	// Index that allows to perform a reverse lookup where given a service ID
	// we retrieve all the operators that have joined it
	Service *indexes.ReversePair[uint32, uint32, collections.NoValue]
}

func (i operatorServiceIndex) IndexesList() []collections.Index[collections.Pair[uint32, uint32], collections.NoValue] {
	return []collections.Index[collections.Pair[uint32, uint32], collections.NoValue]{i.Service}
}

func newOperatorServiceIndex(sb *collections.SchemaBuilder) operatorServiceIndex {
	return operatorServiceIndex{
		Service: indexes.NewReversePair[collections.NoValue](
			sb,
			types.ServiceJoinedByOperatorIndexPrefix,
			"service_joined_by_operator",
			collections.PairKeyCodec(collections.Uint32Key, collections.Uint32Key),
			indexes.WithReversePairUncheckedValue(),
		),
	}
}
