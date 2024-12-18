package testutils

import (
	"testing"

	corestoretypes "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/runtime"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/milkyway-labs/milkyway/v6/testutils/storetesting"
	"github.com/milkyway-labs/milkyway/v6/x/operators/keeper"
	"github.com/milkyway-labs/milkyway/v6/x/operators/types"
)

type KeeperTestData struct {
	storetesting.BaseKeeperTestData

	StoreService corestoretypes.KVStoreService
	Keeper       *keeper.Keeper

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
	data.StoreService = runtime.NewKVStoreService(data.Keys[types.StoreKey])

	// Build keepers
	data.Keeper = keeper.NewKeeper(
		data.Cdc,
		data.StoreService,
		data.AccountKeeper,
		data.DistributionKeeper,
		data.AuthorityAddress,
	)

	// Setup the hooks
	data.Hooks = NewMockHooks()
	data.Keeper = data.Keeper.SetHooks(data.Hooks)

	return data
}
