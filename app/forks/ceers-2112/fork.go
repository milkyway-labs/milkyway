package ceers2112

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v3/app/keepers"
)

func BeginFork(ctx sdk.Context, keepers *keepers.AppKeepers) {
	ctx.Logger().Debug(`
===================================================================================================
==== Forking chain state
===================================================================================================
`)

	// Update the gov params to make sure the min deposit amount is set to stake
	params, err := keepers.GovKeeper.Params.Get(ctx)
	if err != nil {
		panic(err)
	}

	params.MinDeposit = sdk.NewCoins(sdk.NewInt64Coin("stake", 1000000))
	params.ExpeditedMinDeposit = sdk.NewCoins(sdk.NewInt64Coin("stake", 5000000))

	err = keepers.GovKeeper.Params.Set(ctx, params)
	if err != nil {
		panic(err)
	}
}
