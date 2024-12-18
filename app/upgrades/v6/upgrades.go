package v6

import (
	"context"
	"maps"
	"slices"
	"strings"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

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

		// Remove all the markets except from the TIA/USD market
		markets, err := keepers.MarketMapKeeper.GetAllMarkets(sdkCtx)
		if err != nil {
			return nil, err
		}

		for _, ticker := range slices.Sorted(maps.Keys(markets)) {
			market := markets[ticker]

			if strings.Contains(ticker, "TIA") {
				continue
			}

			err = keepers.MarketMapKeeper.DeleteMarket(sdkCtx, ticker)
			if err != nil {
				return nil, err
			}

			err = keepers.OracleKeeper.RemoveCurrencyPair(sdkCtx, market.Ticker.CurrencyPair)
			if err != nil {
				return nil, err
			}
		}

		keepers.AccountKeeper.GetModuleAccount(ctx, liquidvestingtypes.ModuleName)

		return vm, nil
	}
}
