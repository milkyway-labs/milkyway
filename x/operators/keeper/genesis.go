package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

// ExportGenesis returns the GenesisState associated with the given context
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	nextOperatorID, err := k.GetNextOperatorID(ctx)
	if err != nil {
		panic(err)
	}

	operators, err := k.GetOperators(ctx)
	if err != nil {
		panic(err)
	}

	inactiveOperators, err := k.GetInactivatingOperators(ctx)
	if err != nil {
		panic(err)
	}

	operatorParamsRecords, err := k.GetAllOperatorParamsRecords(ctx)
	if err != nil {
		panic(err)
	}

	params, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}

	return types.NewGenesisState(
		nextOperatorID,
		operators,
		operatorParamsRecords,
		inactiveOperators,
		params,
	)
}

// --------------------------------------------------------------------------------------------------------------------

// InitGenesis initializes the state from a GenesisState
func (k *Keeper) InitGenesis(ctx sdk.Context, state *types.GenesisState) error {
	// Set the next operator ID
	err := k.SetNextOperatorID(ctx, state.NextOperatorID)
	if err != nil {
		return err
	}

	// Store the operators
	for _, operator := range state.Operators {
		err = k.CreateOperator(ctx, operator)
		if err != nil {
			return err
		}
	}

	// Store the operator params
	for _, operatorParams := range state.OperatorsParams {
		// Ensure that the operator is present
		_, found, err := k.GetOperator(ctx, operatorParams.OperatorID)
		if err != nil {
			return err
		}

		if !found {
			return fmt.Errorf("can't set operator params for %d, operator not found", operatorParams.OperatorID)
		}

		err = k.SaveOperatorParams(ctx, operatorParams.OperatorID, operatorParams.Params)
		if err != nil {
			return err
		}
	}

	// Store the inactivating operators
	for _, entry := range state.UnbondingOperators {
		k.setOperatorAsInactivating(ctx, entry.OperatorID, entry.UnbondingCompletionTime)
	}

	// Store params
	err = k.SetParams(ctx, state.Params)
	if err != nil {
		return err
	}

	return nil
}
