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
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"
	"github.com/stretchr/testify/require"

	milkyway "github.com/milkyway-labs/milkyway/v6/app"
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

	// Generate the genesis file
	genesis := milkyway.NewDefaultGenesisState(app.AppCodec(), app.ModuleBasics)

	// Get the codes
	cdc, _ := milkyway.MakeCodecs()

	// Make sure there are some markets in the genesis state
	var marketMapGenesis marketmaptypes.GenesisState
	err := cdc.UnmarshalJSON(genesis[marketmaptypes.ModuleName], &marketMapGenesis)
	require.NoError(t, err)
	require.NotEmpty(t, marketMapGenesis.MarketMap.Markets)

	// Make sure the oracle genesis state is properly initialized
	var oracleGenesis oracletypes.GenesisState
	err = cdc.UnmarshalJSON(genesis[oracletypes.ModuleName], &oracleGenesis)
	require.NoError(t, err)
	require.Len(t, oracleGenesis.CurrencyPairGenesis, len(marketMapGenesis.MarketMap.Markets))
}
