package testutils

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"go.uber.org/mock/gomock"

	"github.com/milkyway-labs/milkyway/app/keepers"
	"github.com/milkyway-labs/milkyway/testutils/storetesting"
	poolskeeper "github.com/milkyway-labs/milkyway/x/pools/keeper"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	"github.com/milkyway-labs/milkyway/x/services/keeper"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

type KeeperTestData struct {
	storetesting.BaseKeeperTestData

	StoreKey storetypes.StoreKey

	MockCtrl *gomock.Controller

	PoolKeeper  *MockCommunityPoolKeeper
	PoolsKeeper *poolskeeper.Keeper

	Keeper *keeper.Keeper
	Hooks  *MockHooks
}

// NewKeeperTestData returns a new KeeperTestData
func NewKeeperTestData(t *testing.T) KeeperTestData {
	t.Helper()

	// Initialize the base test data
	var data = KeeperTestData{
		BaseKeeperTestData: storetesting.NewBaseKeeperTestData(t, []string{
			authtypes.StoreKey, banktypes.StoreKey, servicestypes.StoreKey, poolstypes.StoreKey,
		}),
	}

	// Define store keys
	data.StoreKey = data.Keys[servicestypes.StoreKey]

	// Mocks initializations
	data.MockCtrl = gomock.NewController(t)
	data.PoolKeeper = NewMockCommunityPoolKeeper(data.MockCtrl)

	// Build keepers
	data.PoolsKeeper = poolskeeper.NewKeeper(
		data.Cdc,
		data.Keys[poolstypes.StoreKey],
		runtime.NewKVStoreService(data.Keys[poolstypes.StoreKey]),
		data.AccountKeeper,
	)
	data.Keeper = keeper.NewKeeper(
		data.Cdc,
		data.StoreKey,
		runtime.NewKVStoreService(data.Keys[servicestypes.StoreKey]),
		data.AccountKeeper,
		keepers.NewCommunityPoolKeeper(data.BankKeeper, authtypes.FeeCollectorName),
		data.AuthorityAddress,
	)

	// Set hooks
	data.Hooks = NewMockHooks()
	data.Keeper = data.Keeper.SetHooks(data.Hooks)

	return data
}
