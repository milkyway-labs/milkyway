package v11

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/milkyway-labs/milkyway/v11/app/keepers"
	investorstypes "github.com/milkyway-labs/milkyway/v11/x/investors/types"
)

const foundationAddress = "milk108zdtldyt6r98rlg6la6nvwczzxnh2mjajlj4g"

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

		// Set the default investors parameters. Note that it uses
		// UpdateInvestorsRewardRatio instead of SetInvestorsRewardRatio, just in case
		// the investors were already delegating.
		err = keepers.InvestorsKeeper.SetInvestorsRewardRatio(ctx, investorstypes.DefaultInvestorsRewardRatio)
		if err != nil {
			return nil, err
		}

		// Mint MILK token to the foundation account, so that we can distribute it later
		foundationAddr, err := sdk.AccAddressFromBech32(foundationAddress)
		if err != nil {
			return nil, err
		}
		mintAmt := sdk.NewCoins(sdk.NewInt64Coin("umilk", 999_999_980_000_000)) // 1B - 20
		err = keepers.BankKeeper.MintCoins(ctx, minttypes.ModuleName, mintAmt)
		if err != nil {
			return nil, err
		}
		err = keepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, foundationAddr, mintAmt)
		if err != nil {
			return nil, err
		}

		return vm, nil
	}
}
