package v2

import (
	"bytes"
	"fmt"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

// ParseLegacyServiceParamsKey parses the service ID from the given key
func parseLegacyServiceParamsKey(bz []byte) (serviceID uint32, err error) {
	//nolint:staticcheck // SA1004
	// We disable the deprecated lint error
	// since we need to use this key to perform the migration
	bz = bytes.TrimPrefix(bz, types.LegacyServiceParamsPrefix)
	if len(bz) != 4 {
		return 0, fmt.Errorf("invalid key length; expected: 4, got: %d", len(bz))
	}

	return servicestypes.GetServiceIDFromBytes(bz), nil
}

// migrateServiceParams migrates the LegacyServiceParams to the new ServiceParams
// moving some of the parameters contained in the LegacyServiceParams to the
// services module
func migrateServiceParams(
	ctx sdk.Context,
	storeKey storetypes.StoreKey,
	cdc codec.Codec,
	restakingKeeper RestakingKeeper,
	servicesKeeper ServicesKeeper,
) error {
	store := ctx.KVStore(storeKey)

	//nolint:staticcheck // SA1004
	// We disable the deprecated lint error
	// since we need to use this key to perform the migration
	iterator := storetypes.KVStorePrefixIterator(store, types.LegacyServiceParamsPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		//nolint:staticcheck // SA1004
		// We disable the deprecated lint error
		// since we need to use this struct to perform the migration
		var params types.LegacyServiceParamsRecord
		cdc.MustUnmarshal(iterator.Value(), &params)

		serviceID, err := parseLegacyServiceParamsKey(iterator.Key())
		if err != nil {
			return err
		}

		legacyParams := params.Params

		_, found := servicesKeeper.GetService(ctx, serviceID)
		if !found {
			return fmt.Errorf("service %d not found", serviceID)
		}

		//nolint:staticcheck // SA1004
		// We disable the deprecated lint error
		// since we need to use this fields to perform the migration
		// Store the new services params in the restaking module
		newRestakinParams := types.NewServiceParams(legacyParams.WhitelistedPoolsIDs, legacyParams.WhitelistedOperatorsIDs)
		err = restakingKeeper.SaveServiceParams(ctx, serviceID, newRestakinParams)
		if err != nil {
			return err
		}

		//nolint:staticcheck // SA1004
		// We disable the deprecated lint error
		// since we need to use this field to perform the migration
		// Store the new services params in the restaking module
		// Store the service params to the services module
		newServicesParams := servicestypes.NewServiceParams(params.Params.SlashFraction)
		err = servicesKeeper.SaveServiceParams(ctx, serviceID, newServicesParams)
		if err != nil {
			return err
		}

		// Delete the data after migration
		store.Delete(iterator.Key())
	}

	return nil
}
