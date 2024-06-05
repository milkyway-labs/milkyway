package keeper

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

// SetNextOperatorID sets the next operator ID to be used when registering a new Operator
func (k *Keeper) SetNextOperatorID(ctx sdk.Context, operatorID uint32) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.NextOperatorIDKey, types.GetOperatorIDBytes(operatorID))
}

// HasNextOperatorID checks if the next operator ID is set
func (k *Keeper) HasNextOperatorID(ctx sdk.Context) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.NextOperatorIDKey)
}

// GetNextOperatorID returns the next operator ID to be used when registering a new Operator
func (k *Keeper) GetNextOperatorID(ctx sdk.Context) (operatorID uint32, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextOperatorIDKey)
	if bz == nil {
		return 0, errors.Wrapf(types.ErrInvalidGenesis, "initial operator id not set")
	}

	operatorID = types.GetOperatorIDFromBytes(bz)
	return operatorID, nil
}

// --------------------------------------------------------------------------------------------------------------------

// storeOperator stores the given operator in the KVStore
func (k *Keeper) storeOperator(ctx sdk.Context, operator types.Operator) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.OperatorStoreKey(operator.ID), k.cdc.MustMarshal(&operator))
}

// RegisterOperator creates a new Operator and stores it in the KVStore
func (k *Keeper) RegisterOperator(ctx sdk.Context, operator types.Operator) {
	k.storeOperator(ctx, operator)

	k.Logger(ctx).Info("operator created", "id", operator.ID)
	k.AfterOperatorRegistered(ctx, operator.ID)
}

// GetOperator returns the operator with the given ID.
// If the operator does not exist, false is returned.
func (k *Keeper) GetOperator(ctx sdk.Context, operatorID uint32) (operator types.Operator, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.OperatorStoreKey(operatorID))
	if bz == nil {
		return operator, false
	}

	k.cdc.MustUnmarshal(bz, &operator)
	return operator, true
}

// UpdateOperator updates an existing operator in the KVStore
func (k *Keeper) UpdateOperator(ctx sdk.Context, operator types.Operator) error {
	// Make sure the operator exist
	existing, found := k.GetOperator(ctx, operator.ID)
	if !found {
		return types.ErrOperatorNotFound
	}

	// Store the updated operator
	k.storeOperator(ctx, operator)

	// Log the event
	k.Logger(ctx).Info("operator updated", "id", operator.ID)

	// Call the hook based on the operator status change
	switch {
	case existing.Status == types.OPERATOR_STATUS_ACTIVE && operator.Status == types.OPERATOR_STATUS_INACTIVATING:
		k.AfterOperatorInactivatingStarted(ctx, operator.ID)
	case existing.Status == types.OPERATOR_STATUS_INACTIVATING && operator.Status == types.OPERATOR_STATUS_INACTIVE:
		k.AfterOperatorInactivatingCompleted(ctx, operator.ID)
	}

	return nil
}

// StartOperatorInactivation starts the inactivation process for the operator with the given ID
func (k *Keeper) StartOperatorInactivation(ctx sdk.Context, operatorID uint32) error {
	operator, found := k.GetOperator(ctx, operatorID)
	if !found {
		return types.ErrOperatorNotFound
	}

	// Update the operator status
	operator.Status = types.OPERATOR_STATUS_INACTIVATING
	err := k.UpdateOperator(ctx, operator)
	if err != nil {
		return err
	}

	// TODO: Add the operator to the inactivating queue

	return nil
}
