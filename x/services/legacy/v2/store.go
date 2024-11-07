package v2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MigrateStore performs the in-place migration from version 1 to version 2.
// It takes care of iterating through all the services that have been set as allowed to use pool-restaking,
// and setting them as accredited.
func MigrateStore(ctx sdk.Context, k ServicesKeeper, pk PoolsKeeper) error {
	return setAccreditedServices(ctx, k, pk)
}

// setAccreditedServices sets all the services that have been allowed to use pool-restaking as accredited.
func setAccreditedServices(ctx sdk.Context, k ServicesKeeper, pk PoolsKeeper) error {
	// Get the pool's params
	poolsParams := pk.GetParams(ctx)

	// Iterate over all the services that have been allowed to use pool-restaking
	for _, serviceID := range poolsParams.AllowedServicesIDs {
		// Get the service
		service, found := k.GetService(ctx, serviceID)
		if !found {
			continue
		}

		// Set the service as accredited
		service.Accredited = true

		// Save the service
		err := k.SaveService(ctx, service)
		if err != nil {
			return err
		}
	}

	return nil
}
