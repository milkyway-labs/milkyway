package testutils

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/milkyway-labs/milkyway/testutils/storetesting"
	"github.com/milkyway-labs/milkyway/x/operators/keeper"
	"github.com/milkyway-labs/milkyway/x/operators/types"
)

type KeeperTestData struct {
	storetesting.BaseKeeperTestData

	StoreKey storetypes.StoreKey
	Keeper   *keeper.Keeper

	Hooks *MockHooks
}

func NewKeeperTestData(t *testing.T) KeeperTestData {
	t.Helper()

	var data = KeeperTestData{
		BaseKeeperTestData: storetesting.NewBaseKeeperTestData(t, []string{
			types.StoreKey, authtypes.StoreKey, banktypes.StoreKey,
		}),
	}

	// Set the store key
	data.StoreKey = data.Keys[types.StoreKey]

	// Build keepers
	data.Keeper = keeper.NewKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.Keys[types.StoreKey]),
		data.AccountKeeper,
		data.DistributionKeeper,
		data.AuthorityAddress,
	)

	// Setup the hooks
	data.Hooks = NewMockHooks()
	data.Keeper = data.Keeper.SetHooks(data.Hooks)

	return data
}
