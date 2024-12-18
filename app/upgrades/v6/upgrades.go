package v6

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"

	"github.com/milkyway-labs/milkyway/v6/app/keepers"
	liquidvestingtypes "github.com/milkyway-labs/milkyway/v6/x/liquidvesting/types"
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

		// Unwrap the context
		sdkCtx := sdk.UnwrapSDKContext(ctx)

		// Overwrite liquidvesting module account.
		liquidVestingModuleAddr := keepers.AccountKeeper.GetModuleAddress(liquidvestingtypes.ModuleName)
		acc := keepers.AccountKeeper.GetAccount(ctx, liquidVestingModuleAddr)
		if acc != nil {
			_, ok := acc.(sdk.ModuleAccountI)
			if !ok {
				keepers.AccountKeeper.RemoveAccount(ctx, acc)
				keepers.AccountKeeper.GetModuleAccount(ctx, liquidvestingtypes.ModuleName)
			}
		}

		// Delete all currency pairs. We use IterateCurrencyPairs instead of
		// GetAllCurrencyPairs to check the returned error.
		var currencyPairs []connecttypes.CurrencyPair
		err = keepers.OracleKeeper.IterateCurrencyPairs(sdkCtx, func(cp connecttypes.CurrencyPair, _ oracletypes.CurrencyPairState) {
			currencyPairs = append(currencyPairs, cp)
		})
		if err != nil {
			return nil, err
		}
		for _, cp := range currencyPairs {
			err = keepers.OracleKeeper.RemoveCurrencyPair(sdkCtx, cp)
			if err != nil {
				return nil, err
			}
		}

		return vm, nil
	}
}
