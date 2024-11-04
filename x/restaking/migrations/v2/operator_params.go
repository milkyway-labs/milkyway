package v2

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	operatortypes "github.com/milkyway-labs/milkyway/x/operators/types"
	legacytypes "github.com/milkyway-labs/milkyway/x/restaking/legacy/types"
)

// migrateOperatorParams migrates all the operators commissions rates from the
// restaking module to the operators module
func migateOperatorParams(
	ctx sdk.Context,
	storeKey storetypes.StoreKey,
	cdc codec.Codec,
	restakingKeeper RestakingKeeper,
	operatorsKeeper OperatorsKeeper,
) error {
	store := ctx.KVStore(storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, legacytypes.OperatorParamsPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		// Get the operators params from the store
		var params legacytypes.LegacyOperatorParams
		cdc.MustUnmarshal(iterator.Value(), &params)

		// Parse the operator id from the iterator key
		operatorID, err := legacytypes.ParseOperatorParamsKey(iterator.Key())
		if err != nil {
			return err
		}

		// Get the operator from the operators keeper
		_, found := operatorsKeeper.GetOperator(ctx, operatorID)
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
		store.Delete(iterator.Key())
	}
	return nil
}
