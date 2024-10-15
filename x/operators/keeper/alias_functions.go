package keeper

import (
	"fmt"
	"time"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

// createAccountIfNotExists creates an account if it does not exist
func (k *Keeper) createAccountIfNotExists(ctx sdk.Context, address sdk.AccAddress) {
	if !k.accountKeeper.HasAccount(ctx, address) {
		defer telemetry.IncrCounter(1, "new", "account")
		k.accountKeeper.SetAccount(ctx, k.accountKeeper.NewAccountWithAddress(ctx, address))
	}
}

// IterateOperators iterates over the operators in the store and performs a callback function
func (k *Keeper) IterateOperators(ctx sdk.Context, cb func(operator types.Operator) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.OperatorPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var operator types.Operator
		k.cdc.MustUnmarshal(iterator.Value(), &operator)

		if cb(operator) {
			break
		}
	}
}

// GetOperators returns the operators stored in the KVStore
func (k *Keeper) GetOperators(ctx sdk.Context) []types.Operator {
	var operators []types.Operator
	k.IterateOperators(ctx, func(operator types.Operator) (stop bool) {
		operators = append(operators, operator)
		return false
	})
	return operators
}

// IterateInactivatingOperatorQueue iterates over all the operators that are set to be inactivated
// by the given time and calls the given function.
func (k *Keeper) IterateInactivatingOperatorQueue(ctx sdk.Context, endTime time.Time, fn func(operator types.Operator) (stop bool, err error)) error {
	return k.iterateInactivatingOperatorsKeys(ctx, endTime, func(key, value []byte) (stop bool, err error) {
		operatorID, _ := types.SplitInactivatingOperatorQueueKey(key)
		operator, found := k.GetOperator(ctx, operatorID)
		if !found {
			return true, fmt.Errorf("operator %d does not exist", operatorID)
		}

		return fn(operator)
	})
}

// iterateInactivatingOperatorsKeys iterates over all the keys of the operators set to be inactivated
// by the given time, and calls the given function.
// If endTime is zero it iterates over all the keys.
func (k *Keeper) iterateInactivatingOperatorsKeys(ctx sdk.Context, endTime time.Time, fn func(key, value []byte) (stop bool, err error)) error {
	store := ctx.KVStore(k.storeKey)

	var iterator storetypes.Iterator
	if endTime.IsZero() {
		iterator = storetypes.KVStorePrefixIterator(store, types.InactivatingOperatorQueuePrefix)
	} else {
		iterator = store.Iterator(types.InactivatingOperatorQueuePrefix, storetypes.PrefixEndBytes(types.InactivatingOperatorByTime(endTime)))
	}
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		stop, err := fn(iterator.Key(), iterator.Value())
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return nil
}

// GetInactivatingOperators returns the inactivating operators stored in the KVStore
func (k *Keeper) GetInactivatingOperators(ctx sdk.Context) ([]types.UnbondingOperator, error) {
	var operators []types.UnbondingOperator

	err := k.iterateInactivatingOperatorsKeys(ctx, time.Time{}, func(key, value []byte) (stop bool, err error) {
		operatorID, endTime := types.SplitInactivatingOperatorQueueKey(key)
		operators = append(operators, types.NewUnbondingOperator(operatorID, endTime))
		return false, nil
	})
	return operators, err
}

// IsOperatorAddress returns true if the provided address is the address
// where the users' asset are kept when they restake toward an operator.
func (k *Keeper) IsOperatorAddress(ctx sdk.Context, address string) (bool, error) {
	return k.operatorAddressSet.Has(ctx, address)
}

// GetAllOperatorParamsRecords returns all the operator params records
func (k *Keeper) GetAllOperatorParamsRecords(ctx sdk.Context) ([]types.OperatorParamsRecord, error) {
	iterator, err := k.operatorParams.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var records []types.OperatorParamsRecord
	for ; iterator.Valid(); iterator.Next() {
		// Get the operator params
		params, err := iterator.Value()
		if err != nil {
			return nil, err
		}
		// Get the operator id from the map key
		operatorId, err := iterator.Key()
		if err != nil {
			return nil, err
		}
		records = append(records, types.NewOperatorParamsRecord(operatorId, params))
	}

	return records, nil
}
