package keeper

import (
	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v9/x/liquidvesting/types"
)

type Keeper struct {
	cdc          codec.Codec
	storeService corestoretypes.KVStoreService

	// Keepers
	accountKeeper   types.AccountKeeper
	bankKeeper      types.BankKeeper
	operatorsKeeper types.OperatorsKeeper
	poolsKeeper     types.PoolsKeeper
	servicesKeeper  types.ServicesKeeper
	restakingKeeper types.RestakingKeeper
	rewardsKeeper   types.RewardsKeeper

	// Keeper data
	schema         collections.Schema
	params         collections.Item[types.Params]
	insuranceFunds collections.Map[string, types.UserInsuranceFund] // User address -> UserInsuranceFund
	// (delegationType, targetID) -> types.CoveredLockedShares
	TargetsCoveredLockedShares collections.Map[collections.Pair[int32, uint32], types.CoveredLockedShares]

	// Addresses
	ModuleAddress string
	authority     string

	restakingOverrider restakingOverrider
}

func NewKeeper(
	cdc codec.Codec,
	storeService corestoretypes.KVStoreService,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	operatorsKeeper types.OperatorsKeeper,
	poolsKeeper types.PoolsKeeper,
	servicesKeeper types.ServicesKeeper,
	restakingKeeper types.RestakingKeeper,
	rewardsKeeper types.RewardsKeeper,
	moduleAddress string,
	authority string,
) *Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := &Keeper{
		cdc:          cdc,
		storeService: storeService,

		accountKeeper:   accountKeeper,
		bankKeeper:      bankKeeper,
		operatorsKeeper: operatorsKeeper,
		poolsKeeper:     poolsKeeper,
		servicesKeeper:  servicesKeeper,
		restakingKeeper: restakingKeeper,
		rewardsKeeper:   rewardsKeeper,

		params: collections.NewItem(
			sb,
			types.ParamsKey,
			"params",
			codec.CollValue[types.Params](cdc),
		),
		insuranceFunds: collections.NewMap(
			sb,
			types.InsuranceFundKey,
			"insurance_fund",
			collections.StringKey,
			codec.CollValue[types.UserInsuranceFund](cdc),
		),
		TargetsCoveredLockedShares: collections.NewMap(
			sb,
			types.CoveredLockedSharesKeyPrefix,
			"covered_locked_shares",
			collections.PairKeyCodec(collections.Int32Key, collections.Uint32Key),
			codec.CollValue[types.CoveredLockedShares](cdc),
		),

		ModuleAddress: moduleAddress,
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
