package v2

import (
	corestoretypes "cosmossdk.io/core/store"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	operatortypes "github.com/milkyway-labs/milkyway/x/operators/types"
)

func Migrate1To2(
	ctx sdk.Context,
	storeService corestoretypes.KVStoreService,
	cdc codec.Codec,
	restakingKeeper RestakingKeeper,
	operatorsKeeper OperatorsKeeper,
	servicesKeeper ServicesKeeper,
) error {
	err := migrateOperatorParams(ctx, storeService, cdc, restakingKeeper, operatorsKeeper)
	if err != nil {
		return err
	}

	err = migrateServiceParams(ctx, storeService, cdc, restakingKeeper, servicesKeeper)
	if err != nil {
		return err
	}

	return nil
}

// migrateOperatorParams migrates all the operators commissions rates from the
// restaking module to the operators module
func migrateOperatorParams(
	ctx sdk.Context,
	storeService corestoretypes.KVStoreService,
	cdc codec.Codec,
	restakingKeeper RestakingKeeper,
	operatorsKeeper OperatorsKeeper,
) error {
	store := storeService.OpenKVStore(ctx)

	iterator := storetypes.KVStorePrefixIterator(runtime.KVStoreAdapter(store), OperatorParamsPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		// Get the operators params from the store
		var params LegacyOperatorParams
		cdc.MustUnmarshal(iterator.Value(), &params)

		// Parse the operator id from the iterator key
		operatorID, err := ParseOperatorParamsKey(iterator.Key())
		if err != nil {
			return err
		}

		// Get the operator from the operators keeper
		_, found, err := operatorsKeeper.GetOperator(ctx, operatorID)
		if err != nil {
			return err
		}

		if found {
			// Update the operator params with the params retrieved from the
			// restaking module
			err = operatorsKeeper.SaveOperatorParams(ctx, operatorID, operatortypes.NewOperatorParams(params.CommissionRate))
			if err != nil {
				return err
			}

			// Get the operator's joined services.
			// Update the operator joined services with the ones retrieved from
			// the old params structure.
			for _, serviceID := range params.JoinedServicesIDs {
				err = restakingKeeper.AddServiceToOperatorJoinedServices(ctx, operatorID, serviceID)
				if err != nil {
					return err
				}
			}
		}

		// Delete the params from the store
		err = store.Delete(iterator.Key())
		if err != nil {
			return err
		}
	}
	return nil
}

// migrateServiceParams migrates the LegacyServiceParams to the new ServiceParams
// moving the data contained inside the LegacyServiceParams to the restaking module
func migrateServiceParams(
	ctx sdk.Context,
	storeService corestoretypes.KVStoreService,
	cdc codec.Codec,
	restakingKeeper RestakingKeeper,
	servicesKeeper ServicesKeeper,
) error {
	store := storeService.OpenKVStore(ctx)

	iterator := storetypes.KVStorePrefixIterator(runtime.KVStoreAdapter(store), ServiceParamsPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var legacyParams LegacyServiceParams
		cdc.MustUnmarshal(iterator.Value(), &legacyParams)

		serviceID, err := ParseServiceParamsKey(iterator.Key())
		if err != nil {
			return err
		}

		_, found, err := servicesKeeper.GetService(ctx, serviceID)
		if err != nil {
			return err
		}

		if found {
			// Store the service's whitelisted operators in the restaking module
			for _, operatorID := range legacyParams.WhitelistedOperatorsIDs {
				err = restakingKeeper.AddOperatorToServiceAllowList(ctx, serviceID, operatorID)
				if err != nil {
					return err
				}
			}

			// Store the service's whitelisted pools in the restaking module
			for _, poolID := range legacyParams.WhitelistedPoolsIDs {
				err = restakingKeeper.AddPoolToServiceSecuringPools(ctx, serviceID, poolID)
				if err != nil {
					return err
				}
			}
		}

		// Delete the data after migration
		err = store.Delete(iterator.Key())
		if err != nil {
			return err
		}
	}

	return nil
}
