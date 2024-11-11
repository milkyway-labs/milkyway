package milkyway_test

import (
	"testing"

	"cosmossdk.io/log"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	db "github.com/cosmos/cosmos-db"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	milkyway "github.com/milkyway-labs/milkyway/app"
	milkywayhelpers "github.com/milkyway-labs/milkyway/app/helpers"
)

type EmptyAppOptions struct{}

var emptyWasmOption []wasmkeeper.Option

func (ao EmptyAppOptions) Get(_ string) interface{} {
	return nil
}

func TestMilkyWayApp_BlockedModuleAccountAddrs(t *testing.T) {
	app := milkyway.NewMilkyWayApp(
		log.NewNopLogger(),
		db.NewMemDB(),
		nil,
		true,
		map[int64]bool{},
		milkyway.DefaultNodeHome,
		EmptyAppOptions{},
		emptyWasmOption,
	)

	moduleAccountAddresses := app.ModuleAccountAddrs()
	blockedAddrs := app.BlockedModuleAccountAddrs(moduleAccountAddresses)

	require.NotContains(t, blockedAddrs, authtypes.NewModuleAddress(govtypes.ModuleName).String())
}

func TestMilkyWayApp_Export(t *testing.T) {
	app := milkywayhelpers.Setup(t)
	_, err := app.ExportAppStateAndValidators(true, []string{}, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}
