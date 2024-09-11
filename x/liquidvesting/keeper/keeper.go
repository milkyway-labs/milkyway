package keeper

import (
	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
)

type Keeper struct {
	cdc          codec.Codec
	storeService corestoretypes.KVStoreService

	// Keepers
	BankKeeper      types.BankKeeper
	OperatorsKeeper types.OperatorsKeeper
	PoolsKeeper     types.PoolsKeeper
	ServicesKeeper  types.ServicesKeeper
	RestakingKeeper types.RestakingKeeper

	// Keeper data
	Schema         collections.Schema
	Params         collections.Item[types.Params]
	InsuranceFunds collections.Map[sdk.AccAddress, types.UserInsuranceFund]

	// Addresses
	ModuleAddress string
	authority     string
}

func NewKeeper(
	cdc codec.Codec,
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
		storeService: storeService,

		BankKeeper:      bankKeeper,
		OperatorsKeeper: operatorsKeeper,
		PoolsKeeper:     poolsKeeper,
		ServicesKeeper:  servicesKeeper,
		RestakingKeeper: restakingKeeper,

		Params: collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		InsuranceFunds: collections.NewMap[sdk.AccAddress, types.UserInsuranceFund](
			sb,
			types.InsuranceFundKey,
			"insurance_fund",
			sdk.AccAddressKey,
			codec.CollValue[types.UserInsuranceFund](cdc),
		),

		ModuleAddress: moduleAddress,
		authority:     authority,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}
