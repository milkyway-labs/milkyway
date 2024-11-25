package testutils

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/runtime"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	marketmapkeeper "github.com/skip-mev/connect/v2/x/marketmap/keeper"
	marketmaptypes "github.com/skip-mev/connect/v2/x/marketmap/types"
	oraclekeeper "github.com/skip-mev/connect/v2/x/oracle/keeper"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/testutils/storetesting"
	assetskeeper "github.com/milkyway-labs/milkyway/x/assets/keeper"
	assetstypes "github.com/milkyway-labs/milkyway/x/assets/types"
	operatorskeeper "github.com/milkyway-labs/milkyway/x/operators/keeper"
	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolskeeper "github.com/milkyway-labs/milkyway/x/pools/keeper"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	restakingkeeper "github.com/milkyway-labs/milkyway/x/restaking/keeper"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	rewardskeeper "github.com/milkyway-labs/milkyway/x/rewards/keeper"
	rewardstypes "github.com/milkyway-labs/milkyway/x/rewards/types"
	serviceskeeper "github.com/milkyway-labs/milkyway/x/services/keeper"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

type KeeperTestData struct {
	storetesting.BaseKeeperTestData

	MarketMapKeeper *marketmapkeeper.Keeper
	OracleKeeper    oraclekeeper.Keeper
	PoolsKeeper     *poolskeeper.Keeper
	OperatorsKeeper *operatorskeeper.Keeper
	ServicesKeeper  *serviceskeeper.Keeper
	RestakingKeeper *restakingkeeper.Keeper
	AssetsKeeper    *assetskeeper.Keeper

	Keeper *rewardskeeper.Keeper
}

func NewKeeperTestData(t *testing.T) KeeperTestData {
	var data = KeeperTestData{
		BaseKeeperTestData: storetesting.NewBaseKeeperTestData(t, []string{
			marketmaptypes.StoreKey,
			oracletypes.StoreKey,
			poolstypes.StoreKey,
			operatorstypes.StoreKey,
			servicestypes.StoreKey,
			restakingtypes.StoreKey,
			assetstypes.StoreKey,
			rewardstypes.StoreKey,
		}),
	}

	data.Context = data.Context.
		WithBlockHeight(1).
		WithBlockTime(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))

	data.MarketMapKeeper = marketmapkeeper.NewKeeper(
		runtime.NewKVStoreService(data.Keys[marketmaptypes.StoreKey]),
		data.Cdc,
		authtypes.NewModuleAddress(govtypes.ModuleName),
	)

	data.OracleKeeper = oraclekeeper.NewKeeper(
		runtime.NewKVStoreService(data.Keys[oracletypes.StoreKey]),
		data.Cdc,
		data.MarketMapKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName),
	)

	data.PoolsKeeper = poolskeeper.NewKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.Keys[poolstypes.StoreKey]),
		data.AccountKeeper,
	)
	data.OperatorsKeeper = operatorskeeper.NewKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.Keys[operatorstypes.StoreKey]),
		data.AccountKeeper,
		data.BankKeeper,
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
	data.RestakingKeeper = restakingkeeper.NewKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.Keys[restakingtypes.StoreKey]),
		data.AccountKeeper,
		data.BankKeeper,
		data.PoolsKeeper,
		data.OperatorsKeeper,
		data.ServicesKeeper,
		data.AuthorityAddress,
	)
	data.AssetsKeeper = assetskeeper.NewKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.Keys[assetstypes.StoreKey]),
		data.AuthorityAddress,
	)

	data.Keeper = rewardskeeper.NewKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.Keys[rewardstypes.StoreKey]),
		data.AccountKeeper,
		data.BankKeeper,
		data.DistributionKeeper,
		&data.OracleKeeper,
		data.PoolsKeeper,
		data.OperatorsKeeper,
		data.ServicesKeeper,
		data.RestakingKeeper,
		data.AssetsKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Set the hooks
	data.PoolsKeeper.SetHooks(data.Keeper.PoolsHooks())
	data.OperatorsKeeper.SetHooks(operatorstypes.NewMultiOperatorsHooks(
		data.RestakingKeeper.OperatorsHooks(),
		data.Keeper.OperatorsHooks(),
	))
	data.ServicesKeeper.SetHooks(servicestypes.NewMultiServicesHooks(
		data.RestakingKeeper.ServicesHooks(),
		data.Keeper.ServicesHooks(),
	))
	data.RestakingKeeper.SetHooks(data.Keeper.RestakingHooks())

	// Set the base params
	require.NoError(t, data.PoolsKeeper.InitGenesis(data.Context, poolstypes.DefaultGenesis()))
	require.NoError(t, data.ServicesKeeper.InitGenesis(data.Context, servicestypes.DefaultGenesis()))
	require.NoError(t, data.OperatorsKeeper.InitGenesis(data.Context, operatorstypes.DefaultGenesis()))
	require.NoError(t, data.Keeper.InitGenesis(data.Context, rewardstypes.DefaultGenesis()))

	return data
}
