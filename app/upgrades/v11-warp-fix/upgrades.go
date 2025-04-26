package v11warpfix

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	hyperlanetypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	warptypes "github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/milkyway-labs/milkyway/v11/app/keepers"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configuration module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		vm, err := mm.RunMigrations(ctx, configuration, fromVM)
		if err != nil {
			return nil, err
		}

		// Create the warp module account
		keepers.AccountKeeper.GetModuleAccount(ctx, warptypes.ModuleName)
		keepers.AccountKeeper.GetModuleAccount(ctx, hyperlanetypes.ModuleName)

		return vm, nil
	}
}
