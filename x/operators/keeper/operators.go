package keeper

import (
	goerrors "errors"
	"time"

	"cosmossdk.io/collections"
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
func (k *Keeper) SaveOperator(ctx sdk.Context, operator types.Operator) error {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.OperatorStoreKey(operator.ID), k.cdc.MustMarshal(&operator))
	return k.operatorAddressSet.Set(ctx, operator.Address)
}

// StartOperatorInactivation starts the inactivation process for the operator with the given ID
func (k *Keeper) StartOperatorInactivation(ctx sdk.Context, operator types.Operator) error {
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
	k.insertIntoInactivatingQueue(ctx, operator)

	// Call the hook
	return k.AfterOperatorInactivatingStarted(ctx, operator.ID)
}

// CompleteOperatorInactivation completes the inactivation process for the operator with the given ID
func (k *Keeper) CompleteOperatorInactivation(ctx sdk.Context, operator types.Operator) error {
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
func (k *Keeper) SaveOperatorParams(ctx sdk.Context, operatorID uint32, params types.OperatorParams) error {
	return k.operatorParams.Set(ctx, operatorID, params)
}

// GetOperatorParams returns the operator params
func (k *Keeper) GetOperatorParams(ctx sdk.Context, operatorID uint32) (types.OperatorParams, error) {
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
func (k *Keeper) DeleteOperatorParams(ctx sdk.Context, operatorID uint32) error {
	return k.operatorParams.Remove(ctx, operatorID)
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
func (k *Keeper) removeFromInactivatingQueue(ctx sdk.Context, operatorID uint32) error {
	store := ctx.KVStore(k.storeKey)

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
	store.Delete(types.InactivatingOperatorQueueKey(operatorID, inactivatingTime))
	return nil
}
