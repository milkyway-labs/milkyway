package testutils

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	marketmapkeeper "github.com/skip-mev/connect/v2/x/marketmap/keeper"
	marketmaptypes "github.com/skip-mev/connect/v2/x/marketmap/types"
	oraclekeeper "github.com/skip-mev/connect/v2/x/oracle/keeper"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/v9/testutils/storetesting"
	assetskeeper "github.com/milkyway-labs/milkyway/v9/x/assets/keeper"
	assetstypes "github.com/milkyway-labs/milkyway/v9/x/assets/types"
	"github.com/milkyway-labs/milkyway/v9/x/liquidvesting"
	"github.com/milkyway-labs/milkyway/v9/x/liquidvesting/keeper"
	"github.com/milkyway-labs/milkyway/v9/x/liquidvesting/types"
	operatorskeeper "github.com/milkyway-labs/milkyway/v9/x/operators/keeper"
	operatorstypes "github.com/milkyway-labs/milkyway/v9/x/operators/types"
	poolskeeper "github.com/milkyway-labs/milkyway/v9/x/pools/keeper"
	poolstypes "github.com/milkyway-labs/milkyway/v9/x/pools/types"
	restakingkeeper "github.com/milkyway-labs/milkyway/v9/x/restaking/keeper"
	restakingtypes "github.com/milkyway-labs/milkyway/v9/x/restaking/types"
	rewardskeeper "github.com/milkyway-labs/milkyway/v9/x/rewards/keeper"
	rewardstypes "github.com/milkyway-labs/milkyway/v9/x/rewards/types"
	serviceskeeper "github.com/milkyway-labs/milkyway/v9/x/services/keeper"
	servicestypes "github.com/milkyway-labs/milkyway/v9/x/services/types"
)

type KeeperTestData struct {
	storetesting.BaseKeeperTestData

	LiquidVestingModuleAddress sdk.AccAddress

	Keeper          *keeper.Keeper
	IBCMiddleware   porttypes.IBCModule
	OperatorsKeeper *operatorskeeper.Keeper
	PoolsKeeper     *poolskeeper.Keeper
	ServicesKeeper  *serviceskeeper.Keeper
	MarketMapKeeper *marketmapkeeper.Keeper
	OracleKeeper    *oraclekeeper.Keeper
	AssetsKeeper    *assetskeeper.Keeper
	RestakingKeeper *restakingkeeper.Keeper
	RewardsKeeper   *rewardskeeper.Keeper
}

// NewKeeperTestData returns a new KeeperTestData
func NewKeeperTestData(t *testing.T) KeeperTestData {
	t.Helper()

	var data = KeeperTestData{
		BaseKeeperTestData: storetesting.NewBaseKeeperTestData(t, []string{
			types.StoreKey,
			operatorstypes.StoreKey, poolstypes.StoreKey, servicestypes.StoreKey,
			restakingtypes.StoreKey, marketmaptypes.StoreKey, oracletypes.StoreKey,
			assetstypes.StoreKey, rewardstypes.StoreKey,
		}),
	}

	// Module addresses
	data.LiquidVestingModuleAddress = authtypes.NewModuleAddress(types.ModuleName)

	// Build keepers
	data.PoolsKeeper = poolskeeper.NewKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.Keys[poolstypes.StoreKey]),
		data.AccountKeeper,
	)
	data.OperatorsKeeper = operatorskeeper.NewKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.Keys[operatorstypes.StoreKey]),
		data.AccountKeeper,
		data.DistributionKeeper,
		data.AuthorityAddress,
	)
	data.ServicesKeeper = serviceskeeper.NewKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.Keys[servicestypes.StoreKey]),
		data.AccountKeeper,
		data.DistributionKeeper,
		data.AuthorityAddress,
	)
	data.MarketMapKeeper = marketmapkeeper.NewKeeper(
		runtime.NewKVStoreService(data.Keys[marketmaptypes.StoreKey]),
		data.Cdc,
		authtypes.NewModuleAddress(govtypes.ModuleName),
	)
	oracleKeeper := oraclekeeper.NewKeeper(
		runtime.NewKVStoreService(data.Keys[oracletypes.StoreKey]),
		data.Cdc,
		data.MarketMapKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName),
	)
	data.OracleKeeper = &oracleKeeper
	data.AssetsKeeper = assetskeeper.NewKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.Keys[assetstypes.StoreKey]),
		data.AuthorityAddress,
	)
	data.RestakingKeeper = restakingkeeper.NewKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.Keys[restakingtypes.StoreKey]),
		data.AccountKeeper,
		data.BankKeeper,
		data.PoolsKeeper,
		data.OperatorsKeeper,
		data.ServicesKeeper,
		data.OracleKeeper,
		data.AssetsKeeper,
		data.AuthorityAddress,
	)
	data.Keeper = &keeper.Keeper{}
	data.RewardsKeeper = rewardskeeper.NewKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.Keys[rewardstypes.StoreKey]),
		data.AccountKeeper,
		data.BankKeeper,
		data.DistributionKeeper,
		data.OracleKeeper,
		data.PoolsKeeper,
		data.OperatorsKeeper,
		data.Keeper.AdjustedServicesKeeper(data.ServicesKeeper),
		data.Keeper.AdjustedRestakingKeeper(data.RestakingKeeper),
		data.AssetsKeeper,
		data.AuthorityAddress,
	)
	*data.Keeper = *keeper.NewKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.Keys[types.StoreKey]),
		data.AccountKeeper,
		data.BankKeeper,
		data.OperatorsKeeper,
		data.PoolsKeeper,
		data.ServicesKeeper,
		data.RestakingKeeper,
		data.RewardsKeeper,
		data.LiquidVestingModuleAddress.String(),
		data.AuthorityAddress,
	)
	data.PoolsKeeper.SetHooks(
		data.RewardsKeeper.PoolsHooks(),
	)
	data.OperatorsKeeper.SetHooks(operatorstypes.NewMultiOperatorsHooks(
		data.RestakingKeeper.OperatorsHooks(),
		data.RewardsKeeper.OperatorsHooks(),
	))
	data.ServicesKeeper.SetHooks(servicestypes.NewMultiServicesHooks(
		data.RestakingKeeper.ServicesHooks(),
		data.RewardsKeeper.ServicesHooks(),
	))
	data.RestakingKeeper.SetHooks(restakingtypes.NewMultiRestakingHooks(
		data.Keeper.RestakingHooks(),
		data.RewardsKeeper.RestakingHooks(),
	))

	// Set bank hooks
	data.BankKeeper.AppendSendRestriction(data.Keeper.SendRestrictionFn)

	// Set ibc hooks
	var ibcStack porttypes.IBCModule = mockIBCMiddleware{}
	data.IBCMiddleware = liquidvesting.NewIBCMiddleware(ibcStack, data.Keeper)

	// Call InitGenesis to set default params of the modules
	require.NoError(t, data.PoolsKeeper.InitGenesis(data.Context, poolstypes.DefaultGenesis()))
	require.NoError(t, data.OperatorsKeeper.InitGenesis(data.Context, operatorstypes.DefaultGenesis()))
	require.NoError(t, data.ServicesKeeper.InitGenesis(data.Context, servicestypes.DefaultGenesis()))
	require.NoError(t, data.RewardsKeeper.InitGenesis(data.Context, rewardstypes.DefaultGenesis()))
	require.NoError(t, data.RestakingKeeper.InitGenesis(data.Context, restakingtypes.DefaultGenesis()))
	require.NoError(t, data.Keeper.InitGenesis(data.Context, types.DefaultGenesisState()))

	return data
}
