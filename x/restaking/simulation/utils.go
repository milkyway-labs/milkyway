package simulation

import (
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simulation "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/milkyway-labs/milkyway/v2/utils"
	operatorstypes "github.com/milkyway-labs/milkyway/v2/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/v2/x/pools/types"
	"github.com/milkyway-labs/milkyway/v2/x/restaking/keeper"
	"github.com/milkyway-labs/milkyway/v2/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v2/x/services/types"
)

func getOperatorsGenState(simState *module.SimulationState) operatorstypes.GenesisState {
	operatorsGenesisJSON, found := simState.GenState[operatorstypes.ModuleName]
	var operatorsGenesis operatorstypes.GenesisState
	if found {
		simState.Cdc.MustUnmarshalJSON(operatorsGenesisJSON, &operatorsGenesis)
	} else {
		operatorsGenesis = *operatorstypes.DefaultGenesis()
	}

	return operatorsGenesis
}

func getServicesGenState(simState *module.SimulationState) servicestypes.GenesisState {
	servicesGenesisJSON, found := simState.GenState[servicestypes.ModuleName]
	var servicesGenesis servicestypes.GenesisState
	if found {
		simState.Cdc.MustUnmarshalJSON(servicesGenesisJSON, &servicesGenesis)
	} else {
		servicesGenesis = *servicestypes.DefaultGenesis()
	}

	return servicesGenesis
}

func getPoolsGenState(simState *module.SimulationState) poolstypes.GenesisState {
	poolsGenesisJSON, found := simState.GenState[poolstypes.ModuleName]
	var poolsGenesis poolstypes.GenesisState
	if found {
		simState.Cdc.MustUnmarshalJSON(poolsGenesisJSON, &poolsGenesis)
	} else {
		poolsGenesis = *poolstypes.DefaultGenesis()
	}

	return poolsGenesis
}

func RandomOperatorJoinedServices(simState *module.SimulationState) []types.OperatorJoinedServices {
	operatorsGenesis := getOperatorsGenState(simState)
	servicesGenesis := getServicesGenState(simState)

	// Randomly join an operator to a service
	var operatorJoinedServices []types.OperatorJoinedServices
	if len(operatorsGenesis.Operators) > 0 && len(servicesGenesis.Services) > 0 {
		for _, operator := range operatorsGenesis.Operators {
			// 50% of creating a record for this operator
			if simState.Rand.Intn(2) == 0 {
				continue
			}

			var serviceIDs []uint32
			for _, service := range servicesGenesis.Services {
				// 50% of adding this service to the operator
				if simState.Rand.Intn(2) == 0 {
					continue
				}
				serviceIDs = append(serviceIDs, service.ID)
			}

			// Don't add if there's no service
			if len(serviceIDs) == 0 {
				continue
			}

			operatorJoinedServices = append(operatorJoinedServices, types.NewOperatorJoinedServices(operator.ID, serviceIDs))
		}
	}

	return operatorJoinedServices
}

func RandomServiceAllowedOperators(simState *module.SimulationState) []types.ServiceAllowedOperators {
	operatorsGenesis := getOperatorsGenState(simState)
	servicesGenesis := getServicesGenState(simState)

	var serviceAllowedOperators []types.ServiceAllowedOperators
	if len(operatorsGenesis.Operators) > 0 && len(servicesGenesis.Services) > 0 {
		for _, service := range servicesGenesis.Services {
			// 50% of creating an operator allow list for this service
			if simState.Rand.Intn(2) == 0 {
				continue
			}

			var allowedOperatorIDs []uint32
			for _, operator := range operatorsGenesis.Operators {
				// 50% of adding the operator to the allow list
				if simState.Rand.Intn(2) == 0 {
					continue
				}

				allowedOperatorIDs = append(allowedOperatorIDs, operator.ID)
			}
			// Ignore if the allow list is empty
			if len(allowedOperatorIDs) == 0 {
				continue
			}

			serviceAllowedOperators = append(
				serviceAllowedOperators,
				types.NewServiceAllowedOperators(service.ID, allowedOperatorIDs),
			)
		}
	}

	return serviceAllowedOperators
}

func RandomServiceSecuringPools(simState *module.SimulationState) []types.ServiceSecuringPools {
	servicesGenesis := getServicesGenState(simState)
	poolsGenesis := getPoolsGenState(simState)

	var serviceSecuringPools []types.ServiceSecuringPools
	if len(poolsGenesis.Pools) > 0 && len(servicesGenesis.Services) > 0 {
		for _, service := range servicesGenesis.Services {
			// 50% of defining which pools are allowed to secure this service
			if simState.Rand.Intn(2) == 0 {
				continue
			}

			var allowedPoolIDs []uint32
			for _, pool := range poolsGenesis.Pools {
				// 50% of adding the operator to the allow list
				if simState.Rand.Intn(2) == 0 {
					continue
				}

				allowedPoolIDs = append(allowedPoolIDs, pool.ID)
			}
			// Ignore if the allow list is empty
			if len(allowedPoolIDs) == 0 {
				continue
			}

			serviceSecuringPools = append(
				serviceSecuringPools,
				types.NewServiceSecuringPools(service.ID, allowedPoolIDs),
			)
		}
	}

	return serviceSecuringPools
}

func RandomUserPreferencesEntries(simState *module.SimulationState) []types.UserPreferencesEntry {
	servicesGenesis := getServicesGenState(simState)

	var usersPreferences []types.UserPreferencesEntry
	if len(servicesGenesis.Services) > 0 {
		accounts := simulation.RandomAccounts(simState.Rand, simState.Rand.Intn(10))
		for _, account := range accounts {
			// Create the user preferences
			userPreferences := RandomUserPreferences(simState.Rand, servicesGenesis.Services)
			usersPreferences = append(
				usersPreferences,
				types.NewUserPreferencesEntry(account.Address.String(), userPreferences),
			)
		}
	}

	return usersPreferences
}

func RandomParams(simState *module.SimulationState) types.Params {
	unbondingDays := time.Duration(simState.Rand.Intn(7) + 1)
	return types.NewParams(time.Hour*24*unbondingDays, nil)
}

func RandomUserPreferences(r *rand.Rand, services []servicestypes.Service) types.UserPreferences {
	// Add some services to the user's trusted services
	var userTrustedServiceIDs []uint32
	for _, service := range services {
		// 50% of adding the service to the user's trusted services
		if r.Intn(2) == 0 {
			continue
		}
		userTrustedServiceIDs = append(userTrustedServiceIDs, service.ID)
	}

	// Create the user preferences
	userPreferences := types.NewUserPreferences(
		// 50% of trusting non accredited service
		r.Intn(2) == 0,
		// 50% of trusting accredited service
		r.Intn(2) == 0,
		userTrustedServiceIDs,
	)
	return userPreferences
}

func GetRandomExistingDelegation(r *rand.Rand, ctx sdk.Context, k *keeper.Keeper, filter func(delegation types.Delegation) bool) (types.Delegation, bool) {
	delegations, err := k.GetAllDelegations(ctx)
	if err != nil {
		panic(err)
	}

	if filter != nil {
		delegations = utils.Filter(delegations, filter)
	}

	if len(delegations) == 0 {
		return types.Delegation{}, false
	}

	randomIndex := r.Intn(len(delegations))
	return delegations[randomIndex], true
}
