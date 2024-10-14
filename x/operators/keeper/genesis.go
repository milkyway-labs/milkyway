package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

// ExportGenesis returns the GenesisState associated with the given context
func (k *Keeper) ExportGenesis(ctx sdk.Context) (*types.GenesisState, error) {
	inactiveOperators, err := k.GetInactivatingOperators(ctx)
	if err != nil {
		return nil, err
	}

	return types.NewGenesisState(
		k.exportNextOperatorID(ctx),
		k.GetOperators(ctx),
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
	}

	// Store the inactivating operators
	for _, entry := range state.UnbondingOperators {
		k.setOperatorAsInactivating(ctx, entry.OperatorID, entry.UnbondingCompletionTime)
	}

	// Store params
	k.SetParams(ctx, state.Params)

	return nil
}
