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
	setAccreditedServices(ctx, k, pk)
	return nil
}

// setAccreditedServices sets all the services that have been allowed to use pool-restaking as accredited.
func setAccreditedServices(ctx sdk.Context, k ServicesKeeper, pk PoolsKeeper) {
	// Get the pool's params
	poolsParams := pk.GetParams(ctx)

	// Iterate over all the services to update their accreditation status
	k.IterateServices(ctx, func(service servicestypes.Service) (stop bool) {
		service.Accredited = slices.Contains(poolsParams.AllowedServicesIDs, service.ID)

		// Save the service
		err := k.SaveService(ctx, service)
		if err != nil {
			panic(err)
		}

		return false
	})
}
