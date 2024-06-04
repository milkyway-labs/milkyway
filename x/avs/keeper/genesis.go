package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/avs/types"
)

// ExportGenesis returns the GenesisState associated with the given context
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return types.NewGenesisState(
		k.exportNextAVSID(ctx),
		k.exportServices(ctx),
		k.GetParams(ctx),
	)
}

// exportNextAVSID returns the next AVS ID stored in the KVStore
func (k *Keeper) exportNextAVSID(ctx sdk.Context) uint32 {
	nextAVSID, err := k.GetNextAVSID(ctx)
	if err != nil {
		panic(err)
	}
	return nextAVSID
}

// exportServices returns the services stored in the KVStore
func (k *Keeper) exportServices(ctx sdk.Context) []types.AVS {
	var services []types.AVS
	k.IterateServices(ctx, func(service types.AVS) (stop bool) {
		services = append(services, service)
		return false
	})
	return services
}

// --------------------------------------------------------------------------------------------------------------------

// InitGenesis initializes the state from a GenesisState
func (k *Keeper) InitGenesis(ctx sdk.Context, state *types.GenesisState) {
	// Set the next AVS ID
	k.SetNextAVSID(ctx, state.NextAVSID)

	// Store the services
	for _, service := range state.Services {
		k.SaveAVS(ctx, service)
	}

	// Store params
	k.SetParams(ctx, state.Params)
}
