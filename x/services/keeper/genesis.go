package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/services/types"
)

// ExportGenesis returns the GenesisState associated with the given context
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return types.NewGenesisState(
		k.exportNextServiceID(ctx),
		k.exportServices(ctx),
		k.GetParams(ctx),
	)
}

// exportNextServiceID returns the next Service ID stored in the KVStore
func (k *Keeper) exportNextServiceID(ctx sdk.Context) uint32 {
	nextAVSID, err := k.GetNextServiceID(ctx)
	if err != nil {
		panic(err)
	}
	return nextAVSID
}

// exportServices returns the services stored in the KVStore
func (k *Keeper) exportServices(ctx sdk.Context) []types.Service {
	var services []types.Service
	k.IterateServices(ctx, func(service types.Service) (stop bool) {
		services = append(services, service)
		return false
	})
	return services
}

// --------------------------------------------------------------------------------------------------------------------

// InitGenesis initializes the state from a GenesisState
func (k *Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) {
	// Set the next service ID
	k.SetNextServiceID(ctx, state.NextServiceID)

	// Store the services
	for _, service := range state.Services {
		k.SaveService(ctx, service)
	}

	// Store params
	k.SetParams(ctx, state.Params)
}
