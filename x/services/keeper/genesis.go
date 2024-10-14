package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/services/types"
)

// ExportGenesis returns the GenesisState associated with the given context
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return types.NewGenesisState(
		k.exportNextServiceID(ctx),
		k.GetServices(ctx),
		k.GetParams(ctx),
	)
}

// exportNextServiceID returns the next Service ID stored in the KVStore
func (k *Keeper) exportNextServiceID(ctx sdk.Context) uint32 {
	nextServiceID, err := k.GetNextServiceID(ctx)
	if err != nil {
		panic(err)
	}
	return nextServiceID
}

// --------------------------------------------------------------------------------------------------------------------

// InitGenesis initializes the state from a GenesisState
func (k *Keeper) InitGenesis(ctx sdk.Context, state *types.GenesisState) error {
	// Set the next service ID
	k.SetNextServiceID(ctx, state.NextServiceID)

	// Store the services
	for _, service := range state.Services {
		if err := k.SaveService(ctx, service); err != nil {
			return err
		}
	}

	// Store params
	k.SetParams(ctx, state.Params)

	return nil
}
