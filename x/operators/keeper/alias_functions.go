package keeper

import (
	"fmt"
	"time"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

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

// GetOperators returns the services stored in the KVStore
func (k *Keeper) GetOperators(ctx sdk.Context) []types.Operator {
	var operators []types.Operator
	k.IterateOperators(ctx, func(service types.Operator) (stop bool) {
		operators = append(operators, service)
		return false
	})
	return operators
}

// IterateInactivatingOperatorQueue iterates over all the operators that are set to be inactivated
// by the given time and calls the given function.
func (k *Keeper) IterateInactivatingOperatorQueue(ctx sdk.Context, endTime time.Time, fn func(operator types.Operator) (stop bool)) {
	k.iterateInactivatingOperatorsKeys(ctx, endTime, func(key, value []byte) (stop bool) {
		operatorID, _ := types.SplitInactivatingOperatorQueueKey(key)
		operator, found := k.GetOperator(ctx, operatorID)
		if !found {
			panic(fmt.Sprintf("operator %d does not exist", operatorID))
		}

		return fn(operator)
	})
}

// iterateInactivatingOperatorsKeys iterates over all the keys of the operators set to be inactivated
// by the given time, and calls the given function.
// If endTime is zero it iterates over all the keys.
func (k *Keeper) iterateInactivatingOperatorsKeys(ctx sdk.Context, endTime time.Time, fn func(key, value []byte) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	var iterator storetypes.Iterator
	if endTime.IsZero() {
		iterator = store.Iterator(types.InactivatingOperatorQueuePrefix, nil)
	} else {
		iterator = store.Iterator(types.InactivatingOperatorQueuePrefix, storetypes.PrefixEndBytes(types.InactivatingOperatorByTime(endTime)))
	}
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		stop := fn(iterator.Key(), iterator.Value())
		if stop {
			break
		}
	}
}

// GetInactivatingOperators returns the inactivating operators stored in the KVStore
func (k *Keeper) GetInactivatingOperators(ctx sdk.Context) []types.UnbondingOperator {
	var operators []types.UnbondingOperator
	k.iterateInactivatingOperatorsKeys(ctx, time.Time{}, func(key, value []byte) (stop bool) {
		operatorID, endTime := types.SplitInactivatingOperatorQueueKey(key)
		operators = append(operators, types.NewUnbondingOperator(operatorID, endTime))
		return false
	})
	return operators
}
