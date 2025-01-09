package v2

import (
	corestoretypes "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v7/x/restaking/types"
)

// MigrateStore performs in-place store migrations from v1 to v2. The migrations include:
// - Properly setting the delegation-by-target-id values
// - Removing joined operators that are not allowed by the services they have joined
func MigrateStore(ctx sdk.Context, keeper Keeper, storeService corestoretypes.KVStoreService, cdc codec.BinaryCodec) error {
	err := removeNotAllowedJoinedServices(ctx, keeper)
	if err != nil {
		return err
	}

	return setDelegationByTargetIDValues(ctx, keeper, storeService, cdc)
}

// setDelegationByTargetIDValues sets the delegation-by-target-id values for all delegations in the store
func setDelegationByTargetIDValues(ctx sdk.Context, keeper Keeper, storeService corestoretypes.KVStoreService, cdc codec.BinaryCodec) error {
	// Get all the delegations
	poolsDelegations, err := keeper.GetAllPoolDelegations(ctx)
	if err != nil {
		return err
	}

	operatorsDelegations, err := keeper.GetAllOperatorDelegations(ctx)
	if err != nil {
		return err
	}

	servicesDelegations, err := keeper.GetAllServiceDelegations(ctx)
	if err != nil {
		return err
	}

	// Join the delegations together
	allDelegations := append(append(poolsDelegations, operatorsDelegations...), servicesDelegations...)

	// For each delegation, set the delegation-by-target-id value
	store := storeService.OpenKVStore(ctx)
	for _, delegation := range allDelegations {
		_, delegationByTargetIDKeyBuilder, err := types.GetDelegationKeyBuilders(delegation)
		if err != nil {
			return err
		}

		// Set the key
		err = store.Set(delegationByTargetIDKeyBuilder(delegation.TargetID, delegation.UserAddress), cdc.MustMarshal(&delegation))
		if err != nil {
			return err
		}
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

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
