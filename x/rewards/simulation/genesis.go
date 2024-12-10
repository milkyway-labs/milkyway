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

	var operatorAccumalatedCommisionRecords []types.OperatorAccumulatedCommissionRecord
	var poolServiceTotalDelShares []types.PoolServiceTotalDelegatorShares

	genesis := types.NewGenesisState(
		RandomParams(simState.Rand, []string{simState.BondDenom}),
		nextRewardsPlan,
		rewardsPlans,
		nil,
		RandomDelegatorWithdrawInfos(simState.Rand, simState.Accounts),
		RandomDelegationTypeRecords(simState.Rand, []string{simState.BondDenom}),
		RandomDelegationTypeRecords(simState.Rand, []string{simState.BondDenom}),
		RandomDelegationTypeRecords(simState.Rand, []string{simState.BondDenom}),
		operatorAccumalatedCommisionRecords,
		poolServiceTotalDelShares,
	)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(genesis)
}
