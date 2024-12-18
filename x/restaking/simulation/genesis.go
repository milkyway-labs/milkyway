package simulation

import (
	"github.com/cosmos/cosmos-sdk/types/module"

	operatorssimulation "github.com/milkyway-labs/milkyway/v6/x/operators/simulation"
	poolssimulation "github.com/milkyway-labs/milkyway/v6/x/pools/simulation"
	"github.com/milkyway-labs/milkyway/v6/x/restaking/types"
	servicessimulation "github.com/milkyway-labs/milkyway/v6/x/services/simulation"
)

// RandomizedGenState generates a random GenesisState for the restaking module
func RandomizedGenState(simState *module.SimulationState) {
	poolsGenesis := poolssimulation.GetGenesisState(simState)
	operatorsGenesis := operatorssimulation.GetGenesisState(simState)
	servicesGenesis := servicessimulation.GetGenesisState(simState)

	genesis := types.NewGenesis(
		RandomOperatorJoinedServices(simState.Rand, operatorsGenesis.Operators, servicesGenesis.Services),
		RandomServiceAllowedOperators(simState.Rand, servicesGenesis.Services, operatorsGenesis.Operators),
		RandomServiceSecuringPools(simState.Rand, poolsGenesis.Pools, servicesGenesis.Services),
		// empty delegations and undelegations since we need to also perform side effects to other
		// modules to keep the shares consistent.
		nil,
		nil,
		RandomUserPreferencesEntries(simState.Rand, servicesGenesis.Services),
		RandomParams(simState.Rand),
	)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(genesis)
}
