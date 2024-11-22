package simulation

import (
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/milkyway-labs/milkyway/testutils/simtesting"
	"github.com/milkyway-labs/milkyway/x/operators/types"
)

// RandomizedGenState generates a random GenesisState for the services module
func RandomizedGenState(simState *module.SimulationState) {
	// Generate a random list of operators
	var operators []types.Operator
	for i := 0; i < simState.Rand.Intn(100); i++ {
		operators = append(operators, RandomOperator(simState.Rand, simState.Accounts))
	}

	// Get the next operator ID
	var nextOperatorID uint32 = 1
	for _, operator := range operators {
		if operator.ID >= nextOperatorID {
			nextOperatorID = operator.ID + 1
		}
	}

	// Generate the operator params
	var operatorParams []types.OperatorParamsRecord
	for _, operator := range operators {
		// 50% chance of having default params
		if simState.Rand.Intn(2) == 0 {
			continue
		}

		operatorParams = append(operatorParams, types.NewOperatorParamsRecord(
			operator.ID,
			RandomOperatorParams(simState.Rand),
		))
	}

	// Generate the unbonding operators
	var unbondingOperators []types.UnbondingOperator
	for _, operator := range operators {
		// 85% chance of not being unbonded
		if simState.Rand.Intn(100) < 85 {
			continue
		}

		unbondingOperators = append(unbondingOperators, types.NewUnbondingOperator(
			operator.ID,
			simtesting.RandomFutureTime(simState.Rand, simState.GenTimestamp),
		))
	}

	// Generate the params
	params := RandomParams(simState.Rand, simState.BondDenom)

	// Set the genesis state inside the simulation
	genesis := types.NewGenesisState(nextOperatorID, operators, operatorParams, unbondingOperators, params)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(genesis)
}
