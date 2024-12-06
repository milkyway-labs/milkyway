package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/milkyway-labs/milkyway/v3/x/operators/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	*Keeper
}

func NewMsgServer(k *Keeper) types.MsgServer {
	return &msgServer{Keeper: k}
}

// RegisterOperator defines the rpc method for Msg/CreateOperator
func (k msgServer) RegisterOperator(ctx context.Context, msg *types.MsgRegisterOperator) (*types.MsgRegisterOperatorResponse, error) {
	// Get the operator id
	operatorID, err := k.GetNextOperatorID(ctx)
	if err != nil {
		return nil, err
	}

	// Create the new operator
	operator := types.NewOperator(
		operatorID,
		types.OPERATOR_STATUS_ACTIVE,
		msg.Moniker,
		msg.Website,
		msg.PictureURL,
		msg.Sender,
	)

	// Validate the operator before storing
	err = operator.Validate()
	if err != nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	// Charge for the creation
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	if !params.OperatorRegistrationFee.IsZero() {
		// Make sure the specified fees are enough
		if !msg.FeeAmount.IsAnyGTE(params.OperatorRegistrationFee) {
			return nil, errors.Wrapf(sdkerrors.ErrInsufficientFunds, "insufficient funds: %s < %s", msg.FeeAmount, params.OperatorRegistrationFee)
		}

		userAddress, err := sdk.AccAddressFromBech32(operator.Admin)
		if err != nil {
			return nil, errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid operator admin address: %s", operator.Admin)
		}

		err = k.poolKeeper.FundCommunityPool(ctx, msg.FeeAmount, userAddress)
		if err != nil {
			return nil, err
		}
	}

	// Store the operator
	err = k.Keeper.CreateOperator(ctx, operator)
	if err != nil {
		return nil, err
	}

	// Update the ID for the next operator
	err = k.SetNextOperatorID(ctx, operator.ID+1)
	if err != nil {
		return nil, err
	}

	// Emit the event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRegisterOperator,
			sdk.NewAttribute(types.AttributeKeyOperatorID, fmt.Sprintf("%d", operator.ID)),
		),
	})

	return &types.MsgRegisterOperatorResponse{
		NewOperatorID: operatorID,
	}, nil
}

// UpdateOperator defines the rpc method for Msg/UpdateOperator
func (k msgServer) UpdateOperator(ctx context.Context, msg *types.MsgUpdateOperator) (*types.MsgUpdateOperatorResponse, error) {
	// Check if the operator exists
	operator, found, err := k.GetOperator(ctx, msg.OperatorID)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, types.ErrOperatorNotFound
	}

	// Make sure only the admin can update the operator
	if operator.Admin != msg.Sender {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "only the admin can update the operator")
	}

	// Update the operator
	updated := operator.Update(types.NewOperatorUpdate(msg.Moniker, msg.Website, msg.PictureURL))

	// Validate the updated operator before storing
	err = updated.Validate()
	if err != nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	// Store the updated operator
	if err := k.SaveOperator(ctx, updated); err != nil {
		return nil, err
	}

	// Emit the event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateOperator,
			sdk.NewAttribute(types.AttributeKeyOperatorID, fmt.Sprintf("%d", operator.ID)),
		),
	})

	return &types.MsgUpdateOperatorResponse{}, nil
}

// DeactivateOperator defines the rpc method for Msg/DeactivateOperator
func (k msgServer) DeactivateOperator(ctx context.Context, msg *types.MsgDeactivateOperator) (*types.MsgDeactivateOperatorResponse, error) {
	// Check if the operator exists
	operator, found, err := k.GetOperator(ctx, msg.OperatorID)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, types.ErrOperatorNotFound
	}

	// Make sure only the admin can deactivate the operator
	if operator.Admin != msg.Sender {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "only the admin can deactivate the operator")
	}

	// Start the operator inactivation
	if err := k.StartOperatorInactivation(ctx, operator); err != nil {
		return nil, err
	}

	// Emit the event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeStartOperatorInactivation,
			sdk.NewAttribute(types.AttributeKeyOperatorID, fmt.Sprintf("%d", msg.OperatorID)),
		),
	})

	return &types.MsgDeactivateOperatorResponse{}, nil
}

// ReactivateOperator defines the rpc method for Msg/ReactivateOperator
func (k msgServer) ReactivateOperator(ctx context.Context, msg *types.MsgReactivateOperator) (*types.MsgReactivateOperatorResponse, error) {
	// Check if the operator exists
	operator, found, err := k.GetOperator(ctx, msg.OperatorID)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, types.ErrOperatorNotFound
	}

	// Make sure only the admin can reactivate the operator
	if operator.Admin != msg.Sender {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "only the admin can deactivate the operator")
	}

	// Reactivate the operator
	if err := k.ReactivateInactiveOperator(ctx, operator); err != nil {
		return nil, err
	}

	// Emit the event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeReactivateOperator,
			sdk.NewAttribute(types.AttributeKeyOperatorID, fmt.Sprintf("%d", msg.OperatorID)),
		),
	})

	return &types.MsgReactivateOperatorResponse{}, nil
}

// DeleteOperator defines the rpc method for Msg/DeleteOperator
func (k msgServer) DeleteOperator(ctx context.Context, msg *types.MsgDeleteOperator) (*types.MsgDeleteOperatorResponse, error) {
	// Check if the operator exists
	operator, found, err := k.GetOperator(ctx, msg.OperatorID)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, types.ErrOperatorNotFound
	}

	// Make sure only the admin can delete the operator
	if operator.Admin != msg.Sender {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "only the admin can delete the operator")
	}

	// Delete the operator
	if err := k.Keeper.DeleteOperator(ctx, operator); err != nil {
		return nil, err
	}

	// Emit the event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDeleteOperator,
			sdk.NewAttribute(types.AttributeKeyOperatorID, fmt.Sprintf("%d", msg.OperatorID)),
		),
	})

	return &types.MsgDeleteOperatorResponse{}, nil
}

// TransferOperatorOwnership defines the rpc method for Msg/TransferOperatorOwnership
func (k msgServer) TransferOperatorOwnership(ctx context.Context, msg *types.MsgTransferOperatorOwnership) (*types.MsgTransferOperatorOwnershipResponse, error) {
	// Check if the operator exists
	operator, found, err := k.GetOperator(ctx, msg.OperatorID)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, types.ErrOperatorNotFound
	}

	// Make sure only the admin can transfer the operator ownership
	if operator.Admin != msg.Sender {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "only the admin can transfer the operator ownership")
	}

	// Update the operator admin
	operator.Admin = msg.NewAdmin
	if err := k.SaveOperator(ctx, operator); err != nil {
		return nil, err
	}

	// Emit the event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferOperatorOwnership,
			sdk.NewAttribute(types.AttributeKeyOperatorID, fmt.Sprintf("%d", msg.OperatorID)),
			sdk.NewAttribute(types.AttributeKeyNewAdmin, msg.NewAdmin),
		),
	})

	return &types.MsgTransferOperatorOwnershipResponse{}, nil
}

// SetOperatorParams defines the rpc method for Msg/SetOperatorParams
func (k msgServer) SetOperatorParams(ctx context.Context, msg *types.MsgSetOperatorParams) (*types.MsgSetOperatorParamsResponse, error) {
	operator, found, err := k.GetOperator(ctx, msg.OperatorID)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, types.ErrOperatorNotFound
	}

	// Make sure only the admin can update the operator
	if operator.Admin != msg.Sender {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "only the admin can update the operator params")
	}

	// Make sure that the received params are valid
	if err := msg.Params.Validate(); err != nil {
		return nil, errors.Wrap(types.ErrInvalidOperatorParams, err.Error())
	}

	// Update the operator params
	err = k.SaveOperatorParams(ctx, msg.OperatorID, msg.Params)
	if err != nil {
		return nil, err
	}

	// Emit the event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeSetOperatorParams),
	})

	return &types.MsgSetOperatorParamsResponse{}, nil
}

// UpdateParams defines the rpc method for Msg/UpdateParams
func (k msgServer) UpdateParams(ctx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	// Check the authority
	authority := k.authority
	if authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", authority, msg.Authority)
	}

	// Update the params
	err := k.SetParams(ctx, msg.Params)
	if err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
