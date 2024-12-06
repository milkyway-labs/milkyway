package testutils

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/milkyway-labs/milkyway/v3/testutils/storetesting"
	"github.com/milkyway-labs/milkyway/v3/x/assets/keeper"
	"github.com/milkyway-labs/milkyway/v3/x/assets/types"
)

type KeeperTestData struct {
	storetesting.BaseKeeperTestData

	AuthorityAddress string

	StoreKey *storetypes.KVStoreKey

	Keeper *keeper.Keeper
}

func NewKeeperTestData(t *testing.T) KeeperTestData {
	var data = KeeperTestData{
		BaseKeeperTestData: storetesting.NewBaseKeeperTestData(t, []string{
			types.StoreKey,
		}),
	}

	data.StoreKey = data.Keys[types.StoreKey]

	// Setup the addresses
	data.AuthorityAddress = authtypes.NewModuleAddress(govtypes.ModuleName).String()

	// Build the keepers
	data.Keeper = keeper.NewKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.StoreKey),
		data.AuthorityAddress,
	)

	return data
}
