package v2

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	operatortypes "github.com/milkyway-labs/milkyway/x/operators/types"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

// migrateOperatorParams migrates all the operators commissions rates from the
// restaking module to the operators module
func migateOperatorParams(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.Codec, operatorsKeeper types.OperatorsKeeper) error {
	store := ctx.KVStore(storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.OperatorParamsPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		// Get the operators params from the store
		var params types.OperatorParams
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

		// Update the operator params with the params readed from the
		// restaking modeule
		err = operatorsKeeper.SaveOperatorParams(ctx, operatorID, operatortypes.NewOperatorParams(params.CommissionRate))
		if err != nil {
			return err
		}
	}
	return nil
}

func Migrate1To2(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.Codec, operatorsKeeper types.OperatorsKeeper) error {
	return migateOperatorParams(ctx, storeKey, cdc, operatorsKeeper)
}
