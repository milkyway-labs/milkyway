package milkyway_test

import (
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	milkyway "github.com/milkyway-labs/milkyway/v4/app"
	milkywayhelpers "github.com/milkyway-labs/milkyway/v4/app/helpers"
	liquidvestingtypes "github.com/milkyway-labs/milkyway/v4/x/liquidvesting/types"
)

func TestMilkyWayApp_BlockedModuleAccountAddrs(t *testing.T) {
	moduleAccountAddresses := milkyway.ModuleAccountAddrs()
	blockedAddrs := milkyway.BlockedModuleAccountAddrs(moduleAccountAddresses)

	require.NotContains(t, blockedAddrs, authtypes.NewModuleAddress(liquidvestingtypes.ModuleName).String())
	require.NotContains(t, blockedAddrs, authtypes.NewModuleAddress(govtypes.ModuleName).String())
}

func TestMilkyWayApp_Export(t *testing.T) {
	app := milkywayhelpers.Setup(t)
	_, err := app.ExportAppStateAndValidators(true, []string{}, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}
