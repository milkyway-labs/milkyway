package v8

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/milkyway-labs/milkyway/v7/app/keepers"
	"github.com/milkyway-labs/milkyway/v7/utils"
	poolstypes "github.com/milkyway-labs/milkyway/v7/x/pools/types"
	restakingtypes "github.com/milkyway-labs/milkyway/v7/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v7/x/services/types"
)

func BeginFork(ctx sdk.Context, mm *module.Manager, cfg module.Configurator, keepers *keepers.AppKeepers) {
	ctx.Logger().Info(`
===================================================================================================
==== Forking chain state
===================================================================================================
`)

	// Run the store migrations manually since we're not using software upgrade.
	fromVM, err := keepers.UpgradeKeeper.GetModuleVersionMap(ctx)
	if err != nil {
		panic(err)
	}
	vm, err := mm.RunMigrations(ctx, cfg, fromVM)
	if err != nil {
		panic(err)
	}
	err = keepers.UpgradeKeeper.SetModuleVersionMap(ctx, vm)
	if err != nil {
		panic(err)
	}

	err = updateServiceSecuringPools(ctx, keepers)
	if err != nil {
		panic(err)
	}

	err = updatePoolServiceTotalDelShares(ctx, keepers)
	if err != nil {
		panic(err)
	}
}

func updateServiceSecuringPools(ctx sdk.Context, keepers *keepers.AppKeepers) error {
	// Get all pools IDs first.
	pools, err := keepers.PoolsKeeper.GetPools(ctx)
	if err != nil {
		return err
	}
	poolsIDs := utils.Map(pools, func(pool poolstypes.Pool) uint32 {
		return pool.ID
	})

	// Update existing services that have no securing pools configured(which means
	// they are secured by all pools by default) to have all the pools as their
	// securing pools.
	var servicesIDsToUpdate []uint32
	err = keepers.ServicesKeeper.IterateServices(ctx, func(service servicestypes.Service) (stop bool, err error) {
		// We only need to update services that have not been configured to secure pools
		// since we changed the semantic of "empty securing pools"
		configured, err := keepers.RestakingKeeper.IsServiceSecuringPoolsConfigured(ctx, service.ID)
		if err != nil {
			return true, err
		}
		if configured {
			return false, nil
		}

		servicesIDsToUpdate = append(servicesIDsToUpdate, service.ID)
		return false, nil
	})
	if err != nil {
		return err
	}

	for _, serviceID := range servicesIDsToUpdate {
		for _, poolID := range poolsIDs {
			err = keepers.RestakingKeeper.AddPoolToServiceSecuringPools(ctx, serviceID, poolID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// updatePoolServiceTotalDelShares synchronizes the pool-service total delegator shares
// with the updated user preferences.
func updatePoolServiceTotalDelShares(ctx sdk.Context, keepers *keepers.AppKeepers) error {
	preferencesCache := map[string]restakingtypes.UserPreferences{}

	allServices, err := keepers.ServicesKeeper.GetServices(ctx)
	if err != nil {
		return err
	}
	allServicesIDs := utils.Map(allServices, func(service servicestypes.Service) uint32 {
		return service.ID
	})

	// First clear all the existing pool-service total delegator shares.
	err = keepers.RewardsKeeper.PoolServiceTotalDelegatorShares.Clear(ctx, nil)
	if err != nil {
		return err
	}

	// Iterating over all pool delegations, increment pool-service total delegator shares
	// only if the service is trusted with the pool.
	err = keepers.RestakingKeeper.IterateAllPoolDelegations(ctx, func(delegation restakingtypes.Delegation) (stop bool, err error) {
		preferences, ok := preferencesCache[delegation.UserAddress]
		if !ok {
			preferences, err = keepers.RestakingKeeper.GetUserPreferences(ctx, delegation.UserAddress)
			if err != nil {
				return true, err
			}
			preferencesCache[delegation.UserAddress] = preferences
		}

		for _, serviceID := range allServicesIDs {
			if !preferences.IsServiceTrustedWithPool(serviceID, delegation.TargetID) {
				continue
			}

			err = keepers.RewardsKeeper.IncrementPoolServiceTotalDelegatorShares(
				ctx,
				delegation.TargetID,
				serviceID,
				delegation.Shares,
			)
			if err != nil {
				return true, err
			}
		}
		return false, nil
	})
	if err != nil {
		return err
	}
	return nil
}
