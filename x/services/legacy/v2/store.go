package v2

import (
	"slices"

	sdk "github.com/cosmos/cosmos-sdk/types"

	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
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
	poolsParams, err := pk.GetParams(ctx)
	if err != nil {
		return err
	}

	// Iterate over all the services to update their accreditation status
	return k.IterateServices(ctx, func(service servicestypes.Service) (stop bool, err error) {
		service.Accredited = slices.Contains(poolsParams.AllowedServicesIDs, service.ID)

		// Save the service
		err = k.SaveService(ctx, service)
		if err != nil {
			return true, err
		}

		return false, nil
	})
}
