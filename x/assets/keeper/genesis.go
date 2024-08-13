package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/assets/types"
)

// ExportGenesis returns the GenesisState associated with the given context
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params, err := k.Params.Get(ctx)
	if err != nil {
		panic(err)
	}

	assets := []types.Asset{}
	_ = k.Assets.Walk(ctx, nil, func(_ string, asset types.Asset) (stop bool, err error) {
		assets = append(assets, asset)
		return false, nil
	})

	return types.NewGenesisState(params, assets)
}

// --------------------------------------------------------------------------------------------------------------------

// InitGenesis initializes the state from a GenesisState
func (k *Keeper) InitGenesis(ctx sdk.Context, state *types.GenesisState) {
	// Store params
	if err := k.Params.Set(ctx, state.Params); err != nil {
		panic(err)
	}

	for _, asset := range state.Assets {
		if err := k.SetAsset(ctx, asset); err != nil {
			panic(err)
		}
	}
}
