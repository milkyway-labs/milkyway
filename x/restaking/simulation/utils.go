package simulation

import (
	"math/rand"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/milkyway-labs/milkyway/v7/testutils/simtesting"
	"github.com/milkyway-labs/milkyway/v7/utils"
	operatorstypes "github.com/milkyway-labs/milkyway/v7/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/v7/x/pools/types"
	"github.com/milkyway-labs/milkyway/v7/x/restaking/keeper"
	"github.com/milkyway-labs/milkyway/v7/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v7/x/services/types"
)

func RandomOperatorJoinedServices(
	r *rand.Rand,
	operators []operatorstypes.Operator,
	services []servicestypes.Service,
) []types.OperatorJoinedServices {
	// Randomly join an operator to a service
	var operatorJoinedServices []types.OperatorJoinedServices
	if len(operators) == 0 || len(services) == 0 {
		return operatorJoinedServices
	}

	for _, operator := range simtesting.RandomSubSlice(r, operators) {
		serviceIDs := utils.Map(simtesting.RandomSubSlice(r, services), func(s servicestypes.Service) uint32 {
			return s.ID
		})

		// Ignore if the joined services list is empty
		if len(serviceIDs) == 0 {
			continue
		}

		operatorJoinedServices = append(operatorJoinedServices, types.NewOperatorJoinedServices(operator.ID, serviceIDs))
	}

	return operatorJoinedServices
}

func RandomServiceAllowedOperators(
	r *rand.Rand,
	services []servicestypes.Service,
	operators []operatorstypes.Operator,
) []types.ServiceAllowedOperators {
	var serviceAllowedOperators []types.ServiceAllowedOperators
	if len(operators) == 0 || len(services) == 0 {
		return serviceAllowedOperators
	}

	for _, service := range simtesting.RandomSubSlice(r, services) {
		allowedOperatorIDs := utils.Map(simtesting.RandomSubSlice(r, operators), func(o operatorstypes.Operator) uint32 {
			return o.ID
		})

		// Ignore if the allow list is empty
		if len(allowedOperatorIDs) == 0 {
			continue
		}

		serviceAllowedOperators = append(
			serviceAllowedOperators,
			types.NewServiceAllowedOperators(service.ID, allowedOperatorIDs),
		)
	}

	return serviceAllowedOperators
}

func RandomServiceSecuringPools(
	r *rand.Rand,
	pools []poolstypes.Pool,
	services []servicestypes.Service,
) []types.ServiceSecuringPools {
	var serviceSecuringPools []types.ServiceSecuringPools
	if len(pools) == 0 || len(services) == 0 {
		return serviceSecuringPools
	}

	for _, service := range simtesting.RandomSubSlice(r, services) {
		allowedPoolIDs := utils.Map(simtesting.RandomSubSlice(r, pools), func(o poolstypes.Pool) uint32 {
			return o.ID
		})
		// Ignore if the allow list is empty
		if len(allowedPoolIDs) == 0 {
			continue
		}

		serviceSecuringPools = append(
			serviceSecuringPools,
			types.NewServiceSecuringPools(service.ID, allowedPoolIDs))
	}

	return serviceSecuringPools
}

func RandomUserPreferencesEntries(
	r *rand.Rand,
	services []servicestypes.Service,
) []types.UserPreferencesEntry {
	var usersPreferences []types.UserPreferencesEntry
	if len(services) == 0 {
		return usersPreferences
	}

	accounts := simulation.RandomAccounts(r, r.Intn(10))
	for _, account := range accounts {
		// Create the user preferences
		userPreferences := RandomUserPreferences(r, services)
		usersPreferences = append(
			usersPreferences,
			types.NewUserPreferencesEntry(account.Address.String(), userPreferences),
		)
	}

	return usersPreferences
}

func RandomParams(r *rand.Rand) types.Params {
	unbondingDays := time.Duration(r.Intn(7) + 1)
	return types.NewParams(time.Hour*24*unbondingDays, nil, simulation.RandomDecAmount(r, math.LegacyNewDec(10000)), uint32(r.Intn(20)+1))
}

func RandomUserPreferences(r *rand.Rand, services []servicestypes.Service) types.UserPreferences {
	// Add some services to the user's trusted services
	trustedServices := simtesting.RandomSubSlice(r, services)
	trustedServiceEntries := utils.Map(trustedServices, func(s servicestypes.Service) types.TrustedServiceEntry {
		return types.NewTrustedServiceEntry(s.ID, nil)
	})

	return types.NewUserPreferences(trustedServiceEntries)
}

func GetRandomExistingDelegation(
	r *rand.Rand,
	ctx sdk.Context,
	k *keeper.Keeper,
	filter func(delegation types.Delegation) bool,
) (types.Delegation, bool) {
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
