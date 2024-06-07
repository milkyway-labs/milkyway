package keeper

import (
	"time"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

// SetNextOperatorID sets the next operator ID to be used when registering a new Operator
func (k *Keeper) SetNextOperatorID(ctx sdk.Context, operatorID uint32) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.NextOperatorIDKey, types.GetOperatorIDBytes(operatorID))
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

// RegisterOperator creates a new Operator and stores it in the KVStore
func (k *Keeper) RegisterOperator(ctx sdk.Context, operator types.Operator) error {
	// Charge for the creation
	registrationFees := k.GetParams(ctx).OperatorRegistrationFee
	if !registrationFees.IsZero() {
		userAddress, err := sdk.AccAddressFromBech32(operator.Admin)
		if err != nil {
			return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid operator admin address: %s", operator.Admin)
		}

		err = k.poolKeeper.FundCommunityPool(ctx, registrationFees, userAddress)
		if err != nil {
			return err
		}
	}

	// Store the operator
	k.SaveOperator(ctx, operator)

	// Log and call the hooks
	k.Logger(ctx).Info("operator created", "id", operator.ID)
	k.AfterOperatorRegistered(ctx, operator.ID)

	return nil
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

// SaveOperator stores the given operator in the KVStore
func (k *Keeper) SaveOperator(ctx sdk.Context, operator types.Operator) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.OperatorStoreKey(operator.ID), k.cdc.MustMarshal(&operator))
}

// StartOperatorInactivation starts the inactivation process for the operator with the given ID
func (k *Keeper) StartOperatorInactivation(ctx sdk.Context, operator types.Operator) {
	// Update the operator status
	operator.Status = types.OPERATOR_STATUS_INACTIVATING
	k.SaveOperator(ctx, operator)

	// Insert the operator into the inactivating queue
	k.insertIntoInactivatingQueue(ctx, operator)

	// Call the hook
	k.AfterOperatorInactivatingStarted(ctx, operator.ID)
}

// CompleteOperatorInactivation completes the inactivation process for the operator with the given ID
func (k *Keeper) CompleteOperatorInactivation(ctx sdk.Context, operator types.Operator) {
	// Update the operator status
	operator.Status = types.OPERATOR_STATUS_INACTIVE
	k.SaveOperator(ctx, operator)

	// Remove the operator from the inactivating queue
	k.removeFromInactivatingQueue(ctx, operator.ID)

	// Call the hook
	k.AfterOperatorInactivatingCompleted(ctx, operator.ID)
}

// --------------------------------------------------------------------------------------------------------------------

// setOperatorAsInactivating sets the operator as inactivating in the KVStore
func (k *Keeper) setOperatorAsInactivating(ctx sdk.Context, operatorID uint32, endTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.InactivatingOperatorQueueKey(operatorID, endTime), types.GetOperatorIDBytes(operatorID))
}

// insertIntoInactivatingQueue inserts the operator into the inactivating queue
func (k *Keeper) insertIntoInactivatingQueue(ctx sdk.Context, operator types.Operator) {
	endTime := ctx.BlockTime().Add(k.GetParams(ctx).DeactivationTime)
	k.setOperatorAsInactivating(ctx, operator.ID, endTime)
}

// RemoveFromInactivatingQueue removes the operator from the inactivating queue
func (k *Keeper) removeFromInactivatingQueue(ctx sdk.Context, operatorID uint32) {
	store := ctx.KVStore(k.storeKey)

	// Find the inactivating time for the operator
	var inactivatingTime time.Time
	k.iterateInactivatingOperatorsKeys(ctx, time.Time{}, func(key, _ []byte) (stop bool) {
		inactivatingOperatorID, endTime := types.SplitInactivatingOperatorQueueKey(key)
		if inactivatingOperatorID == operatorID {
			inactivatingTime = endTime
			return true
		}

		return false
	})

	// Remove the operator from the inactivating queue
	store.Delete(types.InactivatingOperatorQueueKey(operatorID, inactivatingTime))
}
