package simulation

import (
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/milkyway-labs/milkyway/v3/x/rewards/types"
)

// RandomizedGenState generates a random GenesisState for the operators module
func RandomizedGenState(simState *module.SimulationState) {
	rewardsPlans := RandomRewardsPlans(simState.Rand, []string{simState.BondDenom})
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
		RandomPoolServiceTotalDelegatorShares(simState.Rand, []string{simState.BondDenom}),
	)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(genesis)
}
