package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/milkyway-labs/milkyway/x/services/types"
)

// Simulation parameter constants
const (
	keyServiceRegistrationFees = "operators_count"
	keyServices                = "services"
	keyServicesParams          = "services_params"
)

func getServices(r *rand.Rand, simState *module.SimulationState) []types.Service {
	count := r.Intn(10) + 1
	var services []types.Service
	for i := 0; i < count; i++ {
		adminAccount, _ := simulation.RandomAcc(r, simState.Accounts)
		service := RandomService(r, uint32(i)+1, adminAccount.Address.String())
		services = append(services, service)
	}

	return services
}

func getServiceParams(r *rand.Rand, services []types.Service) []types.ServiceParamsRecord {
	var params []types.ServiceParamsRecord
	for _, service := range services {
		generate := (r.Uint64() % 2) == 0
		if !generate {
			continue
		}

		serviceParams := types.NewServiceParams([]string{"umilk"})
		params = append(params, types.NewServiceParamsRecord(service.ID, serviceParams))
	}
	return params
}

// RandomizedGenState generates a random GenesisState for the services module
func RandomizedGenState(simState *module.SimulationState) {
	var (
		services       []types.Service
		servicesParams []types.ServiceParamsRecord
	)

	simState.AppParams.GetOrGenerate(keyServices, &services, simState.Rand, func(r *rand.Rand) {
		services = getServices(r, simState)
	})

	simState.AppParams.GetOrGenerate(keyServicesParams, &servicesParams, simState.Rand, func(r *rand.Rand) {
		servicesParams = getServiceParams(r, services)
	})

	params := types.DefaultParams()
	nextServiceID := uint32(len(services)) + 1

	servicesGenesis := types.NewGenesisState(nextServiceID, services, servicesParams, params)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(servicesGenesis)
}