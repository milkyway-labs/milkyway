package keeper

import (
	"context"
	"errors"
	"time"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v12/x/operators/types"
)

// createAccountIfNotExists creates an account if it does not exist
func (k *Keeper) createAccountIfNotExists(ctx context.Context, address sdk.AccAddress) {
	if !k.accountKeeper.HasAccount(ctx, address) {
		defer telemetry.IncrCounter(1, "new", "account")
		k.accountKeeper.SetAccount(ctx, k.accountKeeper.NewAccountWithAddress(ctx, address))
	}
}

// IterateOperators iterates over the operators in the store and performs a callback function
func (k *Keeper) IterateOperators(ctx context.Context, cb func(operator types.Operator) (stop bool, err error)) error {
	err := k.operators.Walk(ctx, nil, func(_ uint32, operator types.Operator) (stop bool, err error) {
		return cb(operator)
	})
	return err
}

// GetOperators returns the operators stored in the KVStore
func (k *Keeper) GetOperators(ctx context.Context) ([]types.Operator, error) {
	var operators []types.Operator
	err := k.IterateOperators(ctx, func(operator types.Operator) (stop bool, err error) {
		operators = append(operators, operator)
		return false, nil
	})
	return operators, err
}

// IterateInactivatingOperatorQueue iterates over all the operators that are set to be inactivated
// by the given time and calls the given function.
func (k *Keeper) IterateInactivatingOperatorQueue(ctx context.Context, endTime time.Time, fn func(operator types.Operator) (stop bool, err error)) error {
	return k.iterateInactivatingOperatorsKeys(ctx, endTime, func(key, value []byte) (stop bool, err error) {
		operatorID, _ := types.SplitInactivatingOperatorQueueKey(key)
		operator, err := k.GetOperator(ctx, operatorID)
		if err != nil {
			if errors.Is(err, collections.ErrNotFound) {
				return true, types.ErrOperatorNotFound
			}
			return true, err
		}
		return fn(operator)
	})
}

// iterateInactivatingOperatorsKeys iterates over all the keys of the operators set to be inactivated
// by the given time, and calls the given function.
// If endTime is zero it iterates over all the keys.
func (k *Keeper) iterateInactivatingOperatorsKeys(ctx context.Context, endTime time.Time, fn func(key, value []byte) (stop bool, err error)) error {
	store := k.storeService.OpenKVStore(ctx)

	var iterator storetypes.Iterator
	if endTime.IsZero() {
		iterator = storetypes.KVStorePrefixIterator(runtime.KVStoreAdapter(store), types.InactivatingOperatorQueuePrefix)
	} else {
		storeIterator, err := store.Iterator(types.InactivatingOperatorQueuePrefix, storetypes.PrefixEndBytes(types.InactivatingOperatorByTime(endTime)))
		if err != nil {
			return err
		}
		iterator = storeIterator
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
func (k *Keeper) GetInactivatingOperators(ctx context.Context) ([]types.UnbondingOperator, error) {
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
func (k *Keeper) IsOperatorAddress(ctx context.Context, address string) (bool, error) {
	return k.operatorAddressSet.Has(ctx, address)
}

// GetAllOperatorParamsRecords returns all the operator params records
func (k *Keeper) GetAllOperatorParamsRecords(ctx context.Context) ([]types.OperatorParamsRecord, error) {
	var records []types.OperatorParamsRecord
	err := k.operatorParams.Walk(ctx, nil, func(operatorID uint32, params types.OperatorParams) (stop bool, err error) {
		records = append(records, types.NewOperatorParamsRecord(operatorID, params))
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	return records, nil
}
