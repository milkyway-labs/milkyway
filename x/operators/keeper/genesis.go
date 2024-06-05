package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

// ExportGenesis returns the GenesisState associated with the given context
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return types.NewGenesisState(
		k.exportNextOperatorID(ctx),
		k.exportOperators(ctx),
		k.GetParams(ctx),
	)
}

// exportNextOperatorID returns the next operator ID stored in the KVStore
func (k *Keeper) exportNextOperatorID(ctx sdk.Context) uint32 {
	nextAVSID, err := k.GetNextOperatorID(ctx)
	if err != nil {
		panic(err)
	}
	return nextAVSID
}

// exportOperators returns the services stored in the KVStore
func (k *Keeper) exportOperators(ctx sdk.Context) []types.Operator {
	var operators []types.Operator
	k.IterateOperators(ctx, func(service types.Operator) (stop bool) {
		operators = append(operators, service)
		return false
	})
	return operators
}

// --------------------------------------------------------------------------------------------------------------------

// InitGenesis initializes the state from a GenesisState
func (k *Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) {
	// Set the next operator ID
	k.SetNextOperatorID(ctx, state.NextOperatorID)

	// Store the services
	for _, service := range state.Operators {
		k.storeOperator(ctx, service)
	}

	// Store params
	k.SetParams(ctx, state.Params)
}
