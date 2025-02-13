package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v9/x/pools/types"
)

// ExportGenesis returns a new GenesisState instance containing the information currently present inside the store
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	nextPoolID, err := k.GetNextPoolID(ctx)
	if err != nil {
		panic(err)
	}

	pools, err := k.GetPools(ctx)
	if err != nil {
		panic(err)
	}

	return types.NewGenesis(
		nextPoolID,
		pools,
	)
}

// --------------------------------------------------------------------------------------------------------------------

// InitGenesis initializes the genesis store using the provided data
func (k *Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) error {
	// Set the next pool id
	err := k.SetNextPoolID(ctx, data.NextPoolID)
	if err != nil {
		return err
	}

	// Store the pools
	for _, pool := range data.Pools {
		err = k.SavePool(ctx, pool)
		if err != nil {
			return err
		}

		// Call the hook
		err = k.AfterPoolCreated(ctx, pool.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
