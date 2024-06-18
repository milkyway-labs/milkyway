package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/pools/types"
)

// ExportGenesis returns a new GenesisState instance containing the information currently present inside the store
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return types.NewGenesis(
		k.exportNextPoolID(ctx),
		k.GetPools(ctx),
	)
}

// exportNextPoolID exports the next pool id stored inside the store
func (k *Keeper) exportNextPoolID(ctx sdk.Context) uint32 {
	nextPoolID, err := k.GetNextPoolID(ctx)
	if err != nil {
		panic(err)
	}
	return nextPoolID
}

// --------------------------------------------------------------------------------------------------------------------

// InitGenesis initializes the genesis store using the provided data
func (k *Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) {
	// Set the next pool id
	k.SetNextPoolID(ctx, data.NextPoolID)

	// Store the pools
	for _, pool := range data.Pools {
		err := k.SavePool(ctx, pool)
		if err != nil {
			panic(err)
		}
	}
}
