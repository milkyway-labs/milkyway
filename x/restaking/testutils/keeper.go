package testutils

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	appkeepers "github.com/milkyway-labs/milkyway/app/keepers"
	"github.com/milkyway-labs/milkyway/testutils/storetesting"
	operatorskeeper "github.com/milkyway-labs/milkyway/x/operators/keeper"
	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolskeeper "github.com/milkyway-labs/milkyway/x/pools/keeper"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	"github.com/milkyway-labs/milkyway/x/restaking/keeper"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
	serviceskeeper "github.com/milkyway-labs/milkyway/x/services/keeper"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

type KeeperTestData struct {
	storetesting.BaseKeeperTestData

	StoreKey storetypes.StoreKey

	PoolsKeeper     *poolskeeper.Keeper
	OperatorsKeeper *operatorskeeper.Keeper
	ServicesKeeper  *serviceskeeper.Keeper
	Keeper          *keeper.Keeper
}

func NewKeeperTestData(t *testing.T) KeeperTestData {
	// Build the base data
	data := KeeperTestData{
		BaseKeeperTestData: storetesting.NewBaseKeeperTestData(t, []string{
			types.StoreKey,
			authtypes.StoreKey, banktypes.StoreKey,
			poolstypes.StoreKey, operatorstypes.StoreKey, servicestypes.StoreKey,
		}),
	}

	// Setup the keys
	data.StoreKey = data.Keys[types.StoreKey]

	// Build the keepers
	communityPoolKeeper := appkeepers.NewCommunityPoolKeeper(data.BankKeeper, authtypes.FeeCollectorName)

	data.PoolsKeeper = poolskeeper.NewKeeper(
		data.Cdc,
		data.Keys[poolstypes.StoreKey],
		runtime.NewKVStoreService(data.Keys[poolstypes.StoreKey]),
		data.AccountKeeper,
	)
	data.OperatorsKeeper = operatorskeeper.NewKeeper(
		data.Cdc,
		data.Keys[operatorstypes.StoreKey],
		runtime.NewKVStoreService(data.Keys[operatorstypes.StoreKey]),
		data.AccountKeeper,
		communityPoolKeeper,
		data.AuthorityAddress,
	)
	data.ServicesKeeper = serviceskeeper.NewKeeper(
		data.Cdc,
		data.Keys[servicestypes.StoreKey],
		runtime.NewKVStoreService(data.Keys[servicestypes.StoreKey]),
		data.AccountKeeper,
		communityPoolKeeper,
		data.AuthorityAddress,
	)
	data.Keeper = keeper.NewKeeper(
		data.Cdc,
		data.StoreKey,
		runtime.NewKVStoreService(data.Keys[types.StoreKey]),
		data.AccountKeeper,
		data.BankKeeper,
		data.PoolsKeeper,
		data.OperatorsKeeper,
		data.ServicesKeeper,
		data.AuthorityAddress,
	).SetHooks(NewMockHooks())

	return data
}