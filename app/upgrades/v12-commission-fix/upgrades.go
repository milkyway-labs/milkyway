package v12commissionfix

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/milkyway-labs/milkyway/v12/app/keepers"
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

		// Get the global minimum commission rate.
		stakingParams, err := keepers.StakingKeeper.GetParams(ctx)
		if err != nil {
			return nil, err
		}
		minCommissionRate := stakingParams.MinCommissionRate

		// Update all validators' commssion rate and max commission rate to match the
		// global minimum commission rate.
		validators, err := keepers.StakingKeeper.GetAllValidators(ctx)
		if err != nil {
			return nil, err
		}
		for _, validator := range validators {
			// If the validator's commission rate and maximum commission rate are less than
			// the global minimum commission rate, increase them.
			if validator.Commission.MaxRate.LT(minCommissionRate) {
				validator.Commission.MaxRate = minCommissionRate
			}
			if validator.Commission.Rate.LT(minCommissionRate) {
				validator.Commission.Rate = minCommissionRate
			}
			err = keepers.StakingKeeper.SetValidator(ctx, validator)
			if err != nil {
				return nil, err
			}
		}

		return vm, nil
	}
}
