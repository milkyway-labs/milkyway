package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

// ExportGenesis returns the GenesisState associated with the given context
func (k *Keeper) ExportGenesis(ctx sdk.Context) (*types.GenesisState, error) {
	inactiveOperators, err := k.GetInactivatingOperators(ctx)
	if err != nil {
		return nil, err
	}

	operatorParamsRecords, err := k.GetAllOperatorParamsRecords(ctx)
	if err != nil {
		return nil, err
	}

	return types.NewGenesisState(
		k.exportNextOperatorID(ctx),
		k.GetOperators(ctx),
		operatorParamsRecords,
		inactiveOperators,
		k.GetParams(ctx),
	), nil
}

// exportNextOperatorID returns the next operator ID stored in the KVStore
func (k *Keeper) exportNextOperatorID(ctx sdk.Context) uint32 {
	nextAVSID, err := k.GetNextOperatorID(ctx)
	if err != nil {
		panic(err)
	}
	return nextAVSID
}

// --------------------------------------------------------------------------------------------------------------------

// InitGenesis initializes the state from a GenesisState
func (k *Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) error {
	// Set the next operator ID
	k.SetNextOperatorID(ctx, state.NextOperatorID)

	// Store the operators
	for _, operator := range state.Operators {
		if err := k.SaveOperator(ctx, operator); err != nil {
			return err
		}
		// Init the operator with the default params
		if err := k.SaveOperatorParams(ctx, operator.ID, types.DefaultOperatorParams()); err != nil {
			return err
		}
	}

	// Store the operator params
	for _, operatorParams := range state.OperatorsParams {
		// Ensure that the operator is present
		_, found := k.GetOperator(ctx, operatorParams.OperatorID)
		if !found {
			return fmt.Errorf("can't set operator params for %d, operator not found", operatorParams.OperatorID)
		}

		err := k.SaveOperatorParams(ctx, operatorParams.OperatorID, operatorParams.Params)
		if err != nil {
			return err
		}
	}

	// Store the inactivating operators
	for _, entry := range state.UnbondingOperators {
		k.setOperatorAsInactivating(ctx, entry.OperatorID, entry.UnbondingCompletionTime)
	}

	// Store params
	k.SetParams(ctx, state.Params)

	return nil
}
