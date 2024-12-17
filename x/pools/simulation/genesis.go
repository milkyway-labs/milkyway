package simulation

import (
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/milkyway-labs/milkyway/v4/x/pools/types"
)

// GetGenesisState returns the pools genesis state from the SimulationState
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
