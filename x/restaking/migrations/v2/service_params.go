package v2

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	legacytypes "github.com/milkyway-labs/milkyway/x/restaking/legacy/types"
)

// migrateServiceParams migrates the LegacyServiceParams to the new ServiceParams
// moving the data contained inside the LegacyServiceParams to the restaking module
func migrateServiceParams(
	ctx sdk.Context,
	storeKey storetypes.StoreKey,
	cdc codec.Codec,
	restakingKeeper RestakingKeeper,
	servicesKeeper ServicesKeeper,
) error {
	store := ctx.KVStore(storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, legacytypes.ServiceParamsPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var legacyParams legacytypes.LegacyServiceParams
		cdc.MustUnmarshal(iterator.Value(), &legacyParams)

		serviceID, err := legacytypes.ParseServiceParamsKey(iterator.Key())
		if err != nil {
			return err
		}

		_, found := servicesKeeper.GetService(ctx, serviceID)
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
		store.Delete(iterator.Key())
	}

	return nil
}
