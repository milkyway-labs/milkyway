package simulation

import (
	"github.com/cosmos/cosmos-sdk/types/module"

	poolstypes "github.com/milkyway-labs/milkyway/v3/x/pools/types"
	"github.com/milkyway-labs/milkyway/v3/x/rewards/types"
	servicestypes "github.com/milkyway-labs/milkyway/v3/x/services/types"
)

// RandomizedGenState generates a random GenesisState for the operators module
func RandomizedGenState(simState *module.SimulationState) {
	// Get the services genesis state
	rawServicesGenesis := simState.GenState[servicestypes.ModuleName]
	var servicesGenesis servicestypes.GenesisState
	simState.Cdc.MustUnmarshalJSON(rawServicesGenesis, &servicesGenesis)

	// Get the pools genesis state
	rawPoolsGenesis := simState.GenState[poolstypes.ModuleName]
	var poolsGenesis poolstypes.GenesisState
	simState.Cdc.MustUnmarshalJSON(rawPoolsGenesis, &poolsGenesis)

	rewardsPlans := RandomRewardsPlans(simState.Rand, servicesGenesis.Services, []string{simState.BondDenom})
	nextRewardsPlan := uint64(1)
	for _, plan := range rewardsPlans {
		if plan.ID >= nextRewardsPlan {
			nextRewardsPlan = plan.ID + 1
		}
	}

	genesis := types.NewGenesisState(
		RandomParams(simState.Rand, []string{simState.BondDenom}),
		nextRewardsPlan,
		rewardsPlans,
		nil,
		RandomDelegatorWithdrawInfos(simState.Rand, simState.Accounts),
		// Empty delegation type records since we need to perform side effects on
		// other modules to have valid delegations
		types.NewDelegationTypeRecords(nil, nil, nil, nil),
		types.NewDelegationTypeRecords(nil, nil, nil, nil),
		types.NewDelegationTypeRecords(nil, nil, nil, nil),
		RandomOperatorAccumulatedCommissionRecords(simState.Rand, []string{simState.BondDenom}),
		RandomPoolServiceTotalDelegatorShares(simState.Rand, poolsGenesis, servicesGenesis, []string{simState.BondDenom}),
	)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(genesis)
}
