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
// - Upgrading the user preferences to the new format
func MigrateStore(ctx sdk.Context, keeper Keeper, storeService corestoretypes.KVStoreService, cdc codec.BinaryCodec) error {
	err := removeNotAllowedJoinedServices(ctx, keeper)
	if err != nil {
		return err
	}

	return upgradeUserPreferences(ctx, storeService, cdc)
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

// --------------------------------------------------------------------------------------------------------------------

type userPreferencesEntry struct {
	Key   []byte
	Value UserPreferences
}

// upgradeUserPreferences upgrades the user preferences to the new format
func upgradeUserPreferences(ctx sdk.Context, storeService corestoretypes.KVStoreService, cdc codec.BinaryCodec) error {
	store := storeService.OpenKVStore(ctx)

	iterator, err := store.Iterator(types.UserPreferencesPrefix, nil)
	if err != nil {
		return err
	}

	// Get all the preferences
	var entries []userPreferencesEntry
	for ; iterator.Valid(); iterator.Next() {
		var preferences UserPreferences
		if err = cdc.Unmarshal(iterator.Value(), &preferences); err != nil {
			return err
		}

		entries = append(entries, userPreferencesEntry{
			Key:   iterator.Key(),
			Value: preferences,
		})
	}

	// Close the iterator
	if err = iterator.Close(); err != nil {
		return err
	}

	// Upgrade the preferences
	for _, entry := range entries {
		// Create the new preferences
		var trustedServices []types.TrustedServiceEntry
		for _, serviceID := range entry.Value.TrustedServicesIDs {
			trustedServices = append(trustedServices, types.NewTrustedServiceEntry(serviceID, nil))
		}

		newPreferences := types.NewUserPreferences(trustedServices)

		// Store the preferences
		err = store.Set(entry.Key, cdc.MustMarshal(&newPreferences))
		if err != nil {
			return err
		}
	}

	return nil
}
