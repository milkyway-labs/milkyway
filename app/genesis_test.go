package milkyway_test

import (
	"testing"

	"cosmossdk.io/log"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	oracleconfig "github.com/skip-mev/connect/v2/oracle/config"
	marketmaptypes "github.com/skip-mev/connect/v2/x/marketmap/types"

	milkyway "github.com/milkyway-labs/milkyway/app"
)

func TestNewDefaultGenesisState(t *testing.T) {
	app := milkyway.NewMilkyWayApp(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		map[int64]bool{},
		t.TempDir(),
		oracleconfig.NewDefaultAppConfig(),
		simtestutil.NewAppOptionsWithFlagHome(t.TempDir()),
		[]wasmkeeper.Option{},
		baseapp.SetChainID("milkyway-app"),
	)

	genesis := milkyway.NewDefaultGenesisState(app.AppCodec(), app.ModuleBasics)
	println(genesis[marketmaptypes.ModuleName])
}
