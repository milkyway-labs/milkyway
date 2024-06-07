package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

var (
	_ types.MsgServer = msgServer{}
)

type msgServer struct {
	*Keeper
}

func NewMsgServer(k *Keeper) types.MsgServer {
	return &msgServer{Keeper: k}
}

// RegisterOperator defines the rpc method for Msg/RegisterOperator
func (k msgServer) RegisterOperator(goCtx context.Context, msg *types.MsgRegisterOperator) (*types.MsgRegisterOperatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the next operator id
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
		return nil, errors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
	}

	// Store the operator
	k.Keeper.RegisterOperator(ctx, operator)

	// Update the ID for the next operator
	k.SetNextOperatorID(ctx, operator.ID+1)

	// Emit the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRegisteredOperator,
			sdk.NewAttribute(types.AttributeKeyOperatorID, fmt.Sprintf("%d", operator.ID)),
		),
	})

	return &types.MsgRegisterOperatorResponse{}, nil

}

// UpdateOperator defines the rpc method for Msg/UpdateOperator
func (k msgServer) UpdateOperator(goCtx context.Context, msg *types.MsgUpdateOperator) (*types.MsgUpdateOperatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the operator exists
	operator, found := k.GetOperator(ctx, msg.OperatorID)
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
	err := updated.Validate()
	if err != nil {
		return nil, errors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
	}

	// Store the updated operator
	err = k.Keeper.UpdateOperator(ctx, updated)
	if err != nil {
		return nil, err
	}

	// Emit the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdatedOperator,
			sdk.NewAttribute(types.AttributeKeyOperatorID, fmt.Sprintf("%d", operator.ID)),
		),
	})

	return &types.MsgUpdateOperatorResponse{}, nil
}

// DeactivateOperator defines the rpc method for Msg/DeactivateOperator
func (k msgServer) DeactivateOperator(goCtx context.Context, msg *types.MsgDeactivateOperator) (*types.MsgDeactivateOperatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the operator exists
	operator, found := k.GetOperator(ctx, msg.OperatorID)
	if !found {
		return nil, types.ErrOperatorNotFound
	}

	// Make sure only the admin can deactivate the operator
	if operator.Admin != msg.Sender {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "only the admin can deactivate the operator")
	}

	// Start the operator inactivation
	err := k.Keeper.StartOperatorInactivation(ctx, msg.OperatorID)
	if err != nil {
		return nil, err
	}

	// Emit the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeStartedOperatorInactivation,
			sdk.NewAttribute(types.AttributeKeyOperatorID, fmt.Sprintf("%d", msg.OperatorID)),
		),
	})

	return &types.MsgDeactivateOperatorResponse{}, nil
}
