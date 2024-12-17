package keeper

import (
	"context"
	goerrors "errors"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/milkyway-labs/milkyway/v4/x/operators/types"
)

// SetNextOperatorID sets the next operator ID to be used when registering a new Operator
func (k *Keeper) SetNextOperatorID(ctx context.Context, operatorID uint32) error {
	return k.nextOperatorID.Set(ctx, uint64(operatorID))
}

// GetNextOperatorID returns the next operator ID to be used when registering a new Operator
func (k *Keeper) GetNextOperatorID(ctx context.Context) (operatorID uint32, err error) {
	nextOperatorID, err := k.nextOperatorID.Next(ctx)
	if err != nil {
		return 0, err
	}

	// If the next operator ID is 0, we need to increment it
	if nextOperatorID == 0 {
		return k.GetNextOperatorID(ctx)
	}

	return uint32(nextOperatorID), nil
}

// --------------------------------------------------------------------------------------------------------------------

// CreateOperator creates a new Operator and stores it in the KVStore
func (k *Keeper) CreateOperator(ctx context.Context, operator types.Operator) error {
	// Create the operator account if it does not exist
	operatorAddress, err := sdk.AccAddressFromBech32(operator.Address)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid operator address: %s", operator.Address)
	}
	k.createAccountIfNotExists(ctx, operatorAddress)

	// Store the operator
	if err := k.SaveOperator(ctx, operator); err != nil {
		return err
	}

	// Log and call the hooks
	k.Logger(ctx).Info("operator created", "id", operator.ID)
	return k.AfterOperatorRegistered(ctx, operator.ID)
}

// GetOperator returns the operator with the given ID.
// If the operator does not exist, false is returned.
func (k *Keeper) GetOperator(ctx context.Context, operatorID uint32) (operator types.Operator, err error) {
	return k.operators.Get(ctx, operatorID)
}

// SaveOperator stores the given operator in the KVStore
func (k *Keeper) SaveOperator(ctx context.Context, operator types.Operator) error {
	err := k.operators.Set(ctx, operator.ID, operator)
	if err != nil {
		return err
	}

	return k.operatorAddressSet.Set(ctx, operator.Address)
}

// StartOperatorInactivation starts the inactivation process for the operator with the given ID
func (k *Keeper) StartOperatorInactivation(ctx context.Context, operator types.Operator) error {
	// Make sure the operator is not already inactive or inactivating
	if operator.Status == types.OPERATOR_STATUS_INACTIVATING || operator.Status == types.OPERATOR_STATUS_INACTIVE {
		return types.ErrOperatorNotActive
	}

	// Update the operator status
	operator.Status = types.OPERATOR_STATUS_INACTIVATING
	if err := k.SaveOperator(ctx, operator); err != nil {
		return err
	}

	// Insert the operator into the inactivating queue
	err := k.insertIntoInactivatingQueue(ctx, operator)
	if err != nil {
		return err
	}

	// Call the hook
	return k.AfterOperatorInactivatingStarted(ctx, operator.ID)
}

// ReactivateInactiveOperator reactivates an inactive operator
func (k *Keeper) ReactivateInactiveOperator(ctx context.Context, operator types.Operator) error {
	// Make sure the operator is inactive
	if operator.Status != types.OPERATOR_STATUS_INACTIVE {
		return types.ErrOperatorNotInactive
	}

	// Update the operator status
	operator.Status = types.OPERATOR_STATUS_ACTIVE
	if err := k.SaveOperator(ctx, operator); err != nil {
		return err
	}

	// Call the hook
	return k.AfterOperatorReactivated(ctx, operator.ID)
}

// DeleteOperator deletes the operator with the given ID
func (k *Keeper) DeleteOperator(ctx context.Context, operator types.Operator) error {
	// Make sure the operator is inactive
	if operator.Status != types.OPERATOR_STATUS_INACTIVE {
		return types.ErrOperatorNotInactive
	}

	// Call the hook
	err := k.BeforeOperatorDeleted(ctx, operator.ID)
	if err != nil {
		return err
	}

	// Delete the operator
	err = k.operators.Remove(ctx, operator.ID)
	if err != nil {
		return err
	}

	return k.operatorAddressSet.Remove(ctx, operator.Address)
}

// CompleteOperatorInactivation completes the inactivation process for the operator with the given ID
func (k *Keeper) CompleteOperatorInactivation(ctx context.Context, operator types.Operator) error {
	// Update the operator status
	operator.Status = types.OPERATOR_STATUS_INACTIVE
	if err := k.SaveOperator(ctx, operator); err != nil {
		return err
	}

	// Remove the operator from the inactivating queue
	if err := k.removeFromInactivatingQueue(ctx, operator.ID); err != nil {
		return err
	}

	// Remove the operator params when completed to avoid an invariant breaking.
	err := k.DeleteOperatorParams(ctx, operator.ID)
	if err != nil {
		return err
	}

	// Call the hook
	return k.AfterOperatorInactivatingCompleted(ctx, operator.ID)
}

// --------------------------------------------------------------------------------------------------------------------

// SaveOperatorParams stores the given operator params
func (k *Keeper) SaveOperatorParams(ctx context.Context, operatorID uint32, params types.OperatorParams) error {
	return k.operatorParams.Set(ctx, operatorID, params)
}

// GetOperatorParams returns the operator params
func (k *Keeper) GetOperatorParams(ctx context.Context, operatorID uint32) (types.OperatorParams, error) {
	params, err := k.operatorParams.Get(ctx, operatorID)
	if err != nil {
		if goerrors.Is(err, collections.ErrNotFound) {
			return types.DefaultOperatorParams(), nil
		} else {
			return types.OperatorParams{}, err
		}
	}
	return params, nil
}

// DeleteOperatorParams the operator params associated to the operator with the provided ID.
// If we don't have params associated to the provided operator ID no action will be performed.
func (k *Keeper) DeleteOperatorParams(ctx context.Context, operatorID uint32) error {
	return k.operatorParams.Remove(ctx, operatorID)
}

// --------------------------------------------------------------------------------------------------------------------

// setOperatorAsInactivating sets the operator as inactivating in the KVStore
func (k *Keeper) setOperatorAsInactivating(ctx context.Context, operatorID uint32, endTime time.Time) error {
	store := k.storeService.OpenKVStore(ctx)
	return store.Set(types.InactivatingOperatorQueueKey(operatorID, endTime), types.GetOperatorIDBytes(operatorID))
}

// insertIntoInactivatingQueue inserts the operator into the inactivating queue
func (k *Keeper) insertIntoInactivatingQueue(ctx context.Context, operator types.Operator) error {
	params, err := k.GetParams(ctx)
	if err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	endTime := sdkCtx.BlockTime().Add(params.DeactivationTime)

	return k.setOperatorAsInactivating(ctx, operator.ID, endTime)
}

// RemoveFromInactivatingQueue removes the operator from the inactivating queue
func (k *Keeper) removeFromInactivatingQueue(ctx context.Context, operatorID uint32) error {
	store := k.storeService.OpenKVStore(ctx)

	// Find the inactivating time for the operator
	var inactivatingTime time.Time
	err := k.iterateInactivatingOperatorsKeys(ctx, time.Time{}, func(key, _ []byte) (stop bool, err error) {
		inactivatingOperatorID, endTime := types.SplitInactivatingOperatorQueueKey(key)
		if inactivatingOperatorID == operatorID {
			inactivatingTime = endTime
			return true, nil
		}

		return false, nil
	})
	if err != nil {
		return err
	}

	// Remove the operator from the inactivating queue
	return store.Delete(types.InactivatingOperatorQueueKey(operatorID, inactivatingTime))
}
