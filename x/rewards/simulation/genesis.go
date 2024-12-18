package simulation

import (
	"github.com/cosmos/cosmos-sdk/types/module"

	operatorssimulation "github.com/milkyway-labs/milkyway/v6/x/operators/simulation"
	poolssimulation "github.com/milkyway-labs/milkyway/v6/x/pools/simulation"
	"github.com/milkyway-labs/milkyway/v6/x/rewards/types"
	servicessimulation "github.com/milkyway-labs/milkyway/v6/x/services/simulation"
)

// RandomizedGenState generates a random GenesisState for the operators module
func RandomizedGenState(simState *module.SimulationState) {
	servicesGenesis := servicessimulation.GetGenesisState(simState)
	poolsGenesis := poolssimulation.GetGenesisState(simState)
	operatorsGenesis := operatorssimulation.GetGenesisState(simState)

	rewardsPlans := RandomRewardsPlans(
		simState.Rand,
		poolsGenesis.Pools,
		operatorsGenesis.Operators,
		servicesGenesis.Services,
		[]string{simState.BondDenom},
	)
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

		// Empty accumulated commissions since we need to perform side effects on
		// other modules to have valid commissions
		nil,

		RandomPoolServiceTotalDelegatorShares(simState.Rand, poolsGenesis, servicesGenesis, []string{simState.BondDenom}),
	)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(genesis)
}
