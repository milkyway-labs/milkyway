package v3

import (
	"context"
	"slices"
	"time"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/milkyway-labs/milkyway/v3/app/keepers"
	"github.com/milkyway-labs/milkyway/v3/utils"
	operatorskeeper "github.com/milkyway-labs/milkyway/v3/x/operators/keeper"
	restakingkeeper "github.com/milkyway-labs/milkyway/v3/x/restaking/keeper"
	rewardskeeper "github.com/milkyway-labs/milkyway/v3/x/rewards/keeper"
	serviceskeeper "github.com/milkyway-labs/milkyway/v3/x/services/keeper"
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

		denomsToRemove := []string{
			"ibc/9ACD338BC3B488E0F50A54DE9A844C8326AF0739D917922A9CE04D42AD66017E", // miscalculated TIA
			"ibc/84FBEC4BBB48BD7CC534ED7518F339CCF6C45529DC00C7BFB8605C9EE7D68AFC", // miscalculated stTIA
		}

		err = removeDenomsFromOperatorsModule(ctx, keepers.OperatorsKeeper, denomsToRemove)
		if err != nil {
			return nil, err
		}

		err = removeDenomsFromServicesModule(ctx, keepers.ServicesKeeper, denomsToRemove)
		if err != nil {
			return nil, err
		}

		err = removeDenomsFromRestakingModule(ctx, keepers.RestakingKeeper, denomsToRemove)
		if err != nil {
			return nil, err
		}

		err = removeDenomsFromRewardsModule(ctx, keepers.RewardsKeeper, denomsToRemove)
		if err != nil {
			return nil, err
		}

		// Set trusted delegate for liquid vesting module.
		liquidVestingParams, err := keepers.LiquidVestingKeeper.GetParams(ctx)
		if err != nil {
			return nil, err
		}
		err = keepers.LiquidVestingKeeper.SetParams(ctx, liquidVestingParams)
		if err != nil {
			return nil, err
		}

		// Increase downtime jail duration from 1 minute to 1 hour.
		slashingParams, err := keepers.SlashingKeeper.GetParams(ctx)
		if err != nil {
			return nil, err
		}
		slashingParams.DowntimeJailDuration = time.Hour
		err = keepers.SlashingKeeper.SetParams(ctx, slashingParams)
		if err != nil {
			return nil, err
		}

		return vm, nil
	}
}

func removeDenomsFromOperatorsModule(ctx context.Context, ok *operatorskeeper.Keeper, denomsToRemove []string) error {
	params, err := ok.GetParams(ctx)
	if err != nil {
		return err
	}
	params.OperatorRegistrationFee = utils.Filter(params.OperatorRegistrationFee, func(coin sdk.Coin) bool {
		return !slices.Contains(denomsToRemove, coin.Denom)
	})
	return ok.SetParams(ctx, params)
}

func removeDenomsFromServicesModule(ctx context.Context, sk *serviceskeeper.Keeper, denomsToRemove []string) error {
	params, err := sk.GetParams(ctx)
	if err != nil {
		return err
	}
	params.ServiceRegistrationFee = utils.Filter(params.ServiceRegistrationFee, func(coin sdk.Coin) bool {
		return !slices.Contains(denomsToRemove, coin.Denom)
	})
	return sk.SetParams(ctx, params)
}

func removeDenomsFromRestakingModule(ctx context.Context, rk *restakingkeeper.Keeper, denomsToRemove []string) error {
	params, err := rk.GetParams(ctx)
	if err != nil {
		return err
	}
	params.AllowedDenoms = utils.Filter(params.AllowedDenoms, func(denom string) bool {
		return !slices.Contains(denomsToRemove, denom)
	})
	return rk.SetParams(ctx, params)
}

func removeDenomsFromRewardsModule(ctx context.Context, rk *rewardskeeper.Keeper, denomsToRemove []string) error {
	params, err := rk.GetParams(ctx)
	if err != nil {
		return err
	}
	params.RewardsPlanCreationFee = utils.Filter(params.RewardsPlanCreationFee, func(coin sdk.Coin) bool {
		return !slices.Contains(denomsToRemove, coin.Denom)
	})
	return rk.SetParams(ctx, params)
}
