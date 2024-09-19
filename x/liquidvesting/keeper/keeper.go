package keeper

import (
	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
)

type Keeper struct {
	cdc          codec.Codec
	storeKey     storetypes.StoreKey
	storeService corestoretypes.KVStoreService

	// Keepers
	bankKeeper      types.BankKeeper
	operatorsKeeper types.OperatorsKeeper
	poolsKeeper     types.PoolsKeeper
	servicesKeeper  types.ServicesKeeper
	restakingKeeper types.RestakingKeeper

	// Keeper data
	schema         collections.Schema
	params         collections.Item[types.Params]
	insuranceFunds collections.Map[sdk.AccAddress, types.UserInsuranceFund]

	// Addresses
	moduleAddress string
	authority     string
}

func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	storeService corestoretypes.KVStoreService,
	bankKeeper types.BankKeeper,
	operatorsKeeper types.OperatorsKeeper,
	poolsKeeper types.PoolsKeeper,
	servicesKeeper types.ServicesKeeper,
	restakingKeeper types.RestakingKeeper,
	moduleAddress string,
	authority string,
) *Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := &Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		storeService: storeService,

		bankKeeper:      bankKeeper,
		operatorsKeeper: operatorsKeeper,
		poolsKeeper:     poolsKeeper,
		servicesKeeper:  servicesKeeper,
		restakingKeeper: restakingKeeper,

		params: collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		insuranceFunds: collections.NewMap[sdk.AccAddress, types.UserInsuranceFund](
			sb,
			types.InsuranceFundKey,
			"insurance_fund",
			sdk.AccAddressKey,
			codec.CollValue[types.UserInsuranceFund](cdc),
		),

		moduleAddress: moduleAddress,
		authority:     authority,
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
