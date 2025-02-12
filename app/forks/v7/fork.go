package v7

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/milkyway-labs/milkyway/v9/app/keepers"
)

func BeginFork(ctx sdk.Context, _ *module.Manager, _ module.Configurator, keepers *keepers.AppKeepers) {
	ctx.Logger().Info(`
===================================================================================================
==== Forking chain state
===================================================================================================
`)

	// Set the restaking cap to zero
	params, err := keepers.RestakingKeeper.GetParams(ctx)
	if err != nil {
		panic(err)
	}
	params.RestakingCap = sdkmath.LegacyNewDec(0)

	err = keepers.RestakingKeeper.SetParams(ctx, params)
	if err != nil {
		panic(err)
	}
}
