package services

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/services/keeper"
	"github.com/milkyway-labs/milkyway/x/services/types"
)

// ExportGenesis returns the GenesisState associated with the given context
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState(
		exportNextAVSID(ctx, k),
		exportServices(ctx, k),
		k.GetParams(ctx),
	)
}

// exportNextAVSID returns the next AVS ID stored in the KVStore
func exportNextAVSID(ctx sdk.Context, k keeper.Keeper) uint32 {
	nextAVSID, err := k.GetNextAVSID(ctx)
	if err != nil {
		panic(err)
	}
	return nextAVSID
}

// exportServices returns the services stored in the KVStore
func exportServices(ctx sdk.Context, k keeper.Keeper) []types.AVS {
	var services []types.AVS
	k.IterateServices(ctx, func(service types.AVS) (stop bool) {
		services = append(services, service)
		return false
	})
	return services
}

// InitGenesis initializes the state from a GenesisState
func InitGenesis(ctx sdk.Context, k keeper.Keeper, state types.GenesisState) {
	// Set the next AVS ID
	k.SetNextAVSID(ctx, state.NextAVSID)

	// Store the services
	for _, service := range state.Services {
		k.SaveAVS(ctx, service)
	}

	// Store params
	k.SetParams(ctx, state.Params)
}
