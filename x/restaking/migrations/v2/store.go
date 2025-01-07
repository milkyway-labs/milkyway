package v2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MigrateStore performs in-place store migrations from v1 to v2. The migrations include:
// - Removing joined operators that are not allowed by the services they have joined
func MigrateStore(ctx sdk.Context, keeper Keeper) error {
	return removeNotAllowedJoinedServices(ctx, keeper)
}

type operatorJoinedService struct {
	operatorID uint32
	serviceID  uint32
}

func removeNotAllowedJoinedServices(ctx sdk.Context, keeper Keeper) error {
	// Get the list of deletable operators (i.e. operators that have joined services that they should not have joined)
	var operatorsToDelete []operatorJoinedService
	err := keeper.IterateAllOperatorsJoinedServices(ctx, func(operatorID uint32, serviceID uint32) (stop bool, err error) {
		canOperatorValidator, err := keeper.CanOperatorValidateService(ctx, serviceID, operatorID)
		if err != nil {
			return true, err
		}

		if canOperatorValidator {
			// Skip if the operator is allowed to join the service
			return false, nil
		}

		// Add the service to the list of deletable services if the operator should have not been allowed to join it
		operatorsToDelete = append(operatorsToDelete, operatorJoinedService{
			operatorID: operatorID,
			serviceID:  serviceID,
		})

		return false, nil
	})
	if err != nil {
		return err
	}

	for _, service := range operatorsToDelete {
		err = keeper.RemoveServiceFromOperatorJoinedServices(ctx, service.operatorID, service.serviceID)
		if err != nil {
			return err
		}
	}

	return nil
}
