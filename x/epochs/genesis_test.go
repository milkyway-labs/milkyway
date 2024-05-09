package epochs_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/milkyway-labs/milk/testutil/keeper"
	"github.com/milkyway-labs/milk/testutil/nullify"

	"github.com/milkyway-labs/milk/x/epochs"
	"github.com/milkyway-labs/milk/x/epochs/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.EpochsKeeper(t)
	epochs.InitGenesis(ctx, *k, genesisState)
	got := epochs.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
