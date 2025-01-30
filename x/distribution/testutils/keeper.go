package testutils

import (
	"testing"

	corestoretypes "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	marketmapkeeper "github.com/skip-mev/connect/v2/x/marketmap/keeper"
	marketmaptypes "github.com/skip-mev/connect/v2/x/marketmap/types"
	oraclekeeper "github.com/skip-mev/connect/v2/x/oracle/keeper"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"

	"github.com/milkyway-labs/milkyway/v7/testutils/storetesting"
	assetskeeper "github.com/milkyway-labs/milkyway/v7/x/assets/keeper"
	assetstypes "github.com/milkyway-labs/milkyway/v7/x/assets/types"
	"github.com/milkyway-labs/milkyway/v7/x/distribution/keeper"
	operatorskeeper "github.com/milkyway-labs/milkyway/v7/x/operators/keeper"
	operatorstypes "github.com/milkyway-labs/milkyway/v7/x/operators/types"
	poolskeeper "github.com/milkyway-labs/milkyway/v7/x/pools/keeper"
	poolstypes "github.com/milkyway-labs/milkyway/v7/x/pools/types"
	serviceskeeper "github.com/milkyway-labs/milkyway/v7/x/services/keeper"
	servicestypes "github.com/milkyway-labs/milkyway/v7/x/services/types"
)

type KeeperTestData struct {
	storetesting.BaseKeeperTestData

	StoreService corestoretypes.KVStoreService

	StakingKeeper   *stakingkeeper.Keeper
	PoolsKeeper     *poolskeeper.Keeper
	OperatorsKeeper *operatorskeeper.Keeper
	ServicesKeeper  *serviceskeeper.Keeper
	MarketMapKeeper *marketmapkeeper.Keeper
	OracleKeeper    *oraclekeeper.Keeper
	AssetsKeeper    *assetskeeper.Keeper
	Keeper          keeper.Keeper
}

func NewKeeperTestData(t *testing.T) KeeperTestData {
	// Build the base data
	data := KeeperTestData{
		BaseKeeperTestData: storetesting.NewBaseKeeperTestData(t, []string{
			distrtypes.StoreKey,
			authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
			poolstypes.StoreKey, operatorstypes.StoreKey, servicestypes.StoreKey,
			marketmaptypes.StoreKey, oracletypes.StoreKey, assetstypes.StoreKey,
		}),
	}

	// Setup the keys
	data.StoreService = runtime.NewKVStoreService(data.Keys[distrtypes.StoreKey])

	// Build the keepers
	data.StakingKeeper = stakingkeeper.NewKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.Keys[stakingtypes.StoreKey]),
		data.AccountKeeper,
		data.BankKeeper,
		data.AuthorityAddress,
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
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
	data.Keeper = keeper.NewKeeper(
		data.Cdc,
		data.StoreService,
		data.AccountKeeper,
		data.BankKeeper,
		data.StakingKeeper,
		authtypes.FeeCollectorName,
		data.AuthorityAddress,
	)
	data.StakingKeeper.SetHooks(data.Keeper.Hooks())

	return data
}
