package simulation

import (
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/milkyway-labs/milkyway/v5/x/services/types"
)

// RandomizedGenState generates a random GenesisState for the services module
func RandomizedGenState(simState *module.SimulationState) {
	// Generate a list of random services
	var services []types.Service
	for i := 0; i < simState.Rand.Intn(100); i++ {
		services = append(services, RandomService(simState.Rand, simState.Accounts))
	}

	// Get the next service ID
	var nextServiceID uint32
	for _, service := range services {
		if service.ID >= nextServiceID {
			nextServiceID = service.ID + 1
		}
	}

	// Generate a list of random service params
	var servicesParams []types.ServiceParamsRecord
	for _, service := range services {
		// 50% of chance of not having custom params
		if simState.Rand.Intn(2) == 0 {
			continue
		}

		servicesParams = append(servicesParams, types.NewServiceParamsRecord(
			service.ID,
			RandomServiceParams(simState.Rand),
		))
	}

	// Generate random params
	params := RandomParams(simState.Rand, simState.BondDenom)

	servicesGenesis := types.NewGenesisState(nextServiceID, services, servicesParams, params)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(servicesGenesis)
}

// GetGenesisState returns the services genesis state from the SimulationState
func GetGenesisState(simState *module.SimulationState) types.GenesisState {
	operatorsGenesisJSON, found := simState.GenState[types.ModuleName]
	var operatorsGenesis types.GenesisState
	if found {
		simState.Cdc.MustUnmarshalJSON(operatorsGenesisJSON, &operatorsGenesis)
	} else {
		operatorsGenesis = *types.DefaultGenesis()
	}

	return operatorsGenesis
}
