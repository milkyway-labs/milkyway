package simulation

import (
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/milkyway-labs/milkyway/v2/x/restaking/types"
)

// RandomizedGenState generates a random GenesisState for the restaking module
func RandomizedGenState(simState *module.SimulationState) {
	genesis := types.NewGenesis(
		RandomOperatorJoinedServices(simState),
		RandomServiceAllowedOperators(simState),
		RandomServiceSecuringPools(simState),
		// empty delegations and undelegations since we need to also perform side effects to other
		// modules to keep the shares consistent.
		nil,
		nil,
		RandomUserPreferencesEntries(simState),
		RandomParams(simState),
	)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(genesis)
}
