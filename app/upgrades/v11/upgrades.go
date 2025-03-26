package v11

import (
	"context"
	"encoding/json"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/milkyway-labs/milkyway/v10/app/keepers"
	investorstypes "github.com/milkyway-labs/milkyway/v10/x/investors/types"
)

type UpgradeData struct {
	VestingInvestors []string `json:"vesting_investors"`
}

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

		// Load the embedded upgrade data
		var upgradeData UpgradeData
		err = json.Unmarshal(dataBz, &upgradeData)
		if err != nil {
			return nil, fmt.Errorf("unmarshal upgrade data: %w", err)
		}

		// Set the vesting investors
		for _, investor := range upgradeData.VestingInvestors {
			err = keepers.InvestorsKeeper.SetVestingInvestor(ctx, investor)
			if err != nil {
				return nil, err
			}
		}

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
