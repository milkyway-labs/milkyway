package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/assets/types"
)

// ExportGenesis returns the GenesisState associated with the given context
func (k *Keeper) ExportGenesis(ctx sdk.Context) (*types.GenesisState, error) {
	// Get the params
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	// Get all the assets
	var assets []types.Asset
	_ = k.Assets.Walk(ctx, nil, func(_ string, asset types.Asset) (stop bool, err error) {
		assets = append(assets, asset)
		return false, nil
	})

	return types.NewGenesisState(params, assets), nil
}

// --------------------------------------------------------------------------------------------------------------------

// InitGenesis initializes the state from a GenesisState
func (k *Keeper) InitGenesis(ctx sdk.Context, state *types.GenesisState) error {
	// Store the params
	err := k.Params.Set(ctx, state.Params)
	if err != nil {
		return err
	}

	// Store the assets
	for _, asset := range state.Assets {
		err = k.SetAsset(ctx, asset)
		if err != nil {
			return err
		}
	}

	return nil
}
