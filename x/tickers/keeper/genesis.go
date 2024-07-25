package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/tickers/types"
)

// ExportGenesis returns the GenesisState associated with the given context
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params, err := k.Params.Get(ctx)
	if err != nil {
		panic(err)
	}

	tickers := types.Tickers{}
	_ = k.Tickers.Walk(ctx, nil, func(denom, ticker string) (stop bool, err error) {
		tickers[denom] = ticker
		return false, nil
	})

	return types.NewGenesisState(params, tickers)
}

// --------------------------------------------------------------------------------------------------------------------

// InitGenesis initializes the state from a GenesisState
func (k *Keeper) InitGenesis(ctx sdk.Context, state *types.GenesisState) {
	// Store params
	if err := k.Params.Set(ctx, state.Params); err != nil {
		panic(err)
	}

	for denom, ticker := range state.Tickers {
		if err := k.SetTicker(ctx, denom, ticker); err != nil {
			panic(err)
		}
	}
}
