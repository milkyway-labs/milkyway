package v10

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/milkyway-labs/milkyway/v9/app/keepers"
	investorstypes "github.com/milkyway-labs/milkyway/v9/x/investors/types"
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

		// Set the default investors parameters
		err = keepers.InvestorsKeeper.UpdateInvestorsRewardRatio(ctx, investorstypes.DefaultInvestorsRewardRatio)
		if err != nil {
			return nil, err
		}
		// Create the module account if it doesn't exist
		keepers.AccountKeeper.GetModuleAccount(ctx, investorstypes.ModuleName)

		return vm, nil
	}
}
