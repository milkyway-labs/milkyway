package v11

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/milkyway-labs/milkyway/v10/app/keepers"
	investorstypes "github.com/milkyway-labs/milkyway/v10/x/investors/types"
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

		// Create the module account if it doesn't exist
		keepers.AccountKeeper.GetModuleAccount(ctx, investorstypes.ModuleName)

		// TODO: specify vesting investors list
		//keepers.InvestorsKeeper.SetVestingInvestor(ctx, "...")

		// Set the default investors parameters. Note that it uses
		// UpdateInvestorsRewardRatio instead of SetInvestorsRewardRatio, just in case
		// the investors were already delegating.
		err = keepers.InvestorsKeeper.UpdateInvestorsRewardRatio(ctx, investorstypes.DefaultInvestorsRewardRatio)
		if err != nil {
			return nil, err
		}

		return vm, nil
	}
}
