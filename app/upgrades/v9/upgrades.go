package v9

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/milkyway-labs/milkyway/v7/app/keepers"
	"github.com/milkyway-labs/milkyway/v7/x/restaking/types"
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

		// Set the default max entries parameter
		params, err := keepers.RestakingKeeper.GetParams(ctx)
		if err != nil {
			return nil, err
		}
		params.MaxEntries = types.DefaultMaxEntries
		err = keepers.RestakingKeeper.SetParams(ctx, params)
		if err != nil {
			return nil, err
		}

		return vm, nil
	}
}
