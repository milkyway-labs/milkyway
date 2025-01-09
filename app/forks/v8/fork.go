package v8

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v7/app/keepers"
	"github.com/milkyway-labs/milkyway/v7/utils"
	poolstypes "github.com/milkyway-labs/milkyway/v7/x/pools/types"
	servicestypes "github.com/milkyway-labs/milkyway/v7/x/services/types"
)

func BeginFork(ctx sdk.Context, keepers *keepers.AppKeepers) {
	ctx.Logger().Info(`
===================================================================================================
==== Forking chain state
===================================================================================================
`)

	// Get all pools IDs first.
	pools, err := keepers.PoolsKeeper.GetPools(ctx)
	if err != nil {
		panic(err)
	}
	poolIDs := utils.Map(pools, func(pool poolstypes.Pool) uint32 {
		return pool.ID
	})

	// Update existing services that have no securing pools configured(which means
	// they are secured by all pools by default) to have all the pools as their
	// securing pools.
	var serviceIDsToUpdate []uint32
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

		serviceIDsToUpdate = append(serviceIDsToUpdate, service.ID)
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	for _, serviceID := range serviceIDsToUpdate {
		for _, poolID := range poolIDs {
			err = keepers.RestakingKeeper.AddPoolToServiceSecuringPools(ctx, serviceID, poolID)
			if err != nil {
				panic(err)
			}
		}
	}
}
