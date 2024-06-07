package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

// ExportGenesis returns the GenesisState associated with the given context
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return types.NewGenesisState(
		k.exportNextOperatorID(ctx),
		k.exportOperators(ctx),
		k.exportInactivatingOperators(ctx),
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

// exportInactivatingOperators returns the inactivating operators stored in the KVStore
func (k *Keeper) exportInactivatingOperators(ctx sdk.Context) []types.UnbondingOperator {
	var operators []types.UnbondingOperator
	k.iterateInactivatingOperatorsKeys(ctx, time.Time{}, func(key, value []byte) (stop bool) {
		operatorID, endTime := types.SplitInactivatingOperatorQueueKey(key)
		operators = append(operators, types.NewUnbondingOperator(operatorID, endTime))
		return false
	})
	return operators
}

// --------------------------------------------------------------------------------------------------------------------

// InitGenesis initializes the state from a GenesisState
func (k *Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) {
	// Set the next operator ID
	k.SetNextOperatorID(ctx, state.NextOperatorID)

	// Store the operators
	for _, operator := range state.Operators {
		k.storeOperator(ctx, operator)
	}

	// Store the inactivating operators
	for _, entry := range state.UnbondingOperators {
		k.setOperatorAsInactivating(ctx, entry.OperatorID, entry.UnbondCompletionTime)
	}

	// Store params
	k.SetParams(ctx, state.Params)
}
