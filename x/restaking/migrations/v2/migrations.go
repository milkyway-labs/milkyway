package v2

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	operatortypes "github.com/milkyway-labs/milkyway/x/operators/types"
	legacytypes "github.com/milkyway-labs/milkyway/x/restaking/legacy/types"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
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
	iterator := storetypes.KVStorePrefixIterator(store, types.OperatorParamsPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		// Get the operators params from the store
		var params legacytypes.LegacyOperatorParams
		cdc.MustUnmarshal(iterator.Value(), &params)

		// Parse the operator id from the iterator key
		operatorID, err := types.ParseOperatorParamsKey(iterator.Key())
		if err != nil {
			return err
		}

		// Get the operator from the operators keeper
		_, found := operatorsKeeper.GetOperator(ctx, operatorID)
		if !found {
			return fmt.Errorf("operator %d not found", operatorID)
		}

		// Update the operator params with the params retrieved from the
		// restaking module
		err = operatorsKeeper.SaveOperatorParams(ctx, operatorID,
			operatortypes.NewOperatorParams(params.CommissionRate),
		)
		if err != nil {
			return err
		}

		// Get the operator's joined services.
		joinedServices, err := restakingKeeper.GetOperatorJoinedServices(ctx, operatorID)
		if err != nil {
			return err
		}

		// Update the operator joined services with the ones retrieved from
		// the old params structure.
		for _, serviceID := range params.JoinedServicesIDs {
			err := joinedServices.Add(serviceID)
			if err != nil {
				return err
			}
		}

		// Store the services joined by the operator
		err = restakingKeeper.SaveOperatorJoinedServices(ctx, operatorID, joinedServices)
		if err != nil {
			return err
		}

		// Delete the params from the store
		store.Delete(iterator.Key())
	}
	return nil
}

func Migrate1To2(
	ctx sdk.Context,
	storeKey storetypes.StoreKey,
	cdc codec.Codec,
	restakingKeeper RestakingKeeper,
	operatorsKeeper OperatorsKeeper,
) error {
	return migateOperatorParams(ctx, storeKey, cdc, restakingKeeper, operatorsKeeper)
}
