package keeper

import (
	"bytes"
	"context"
	"fmt"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

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

	// Store the operator
	err = k.Keeper.RegisterOperator(ctx, operator)
	if err != nil {
		return nil, err
	}

	// Update the ID for the next operator
	k.SetNextOperatorID(ctx, operator.ID+1)

	// Emit the event
	ctx.EventManager().EmitEvents(sdk.Events{
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
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	// Store the updated operator
	k.SaveOperator(ctx, updated)

	// Emit the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateOperator,
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

	// Make sure the operator is not already inactive or inactivating
	if operator.Status == types.OPERATOR_STATUS_INACTIVATING || operator.Status == types.OPERATOR_STATUS_INACTIVE {
		return nil, errors.Wrapf(sdkerrors.ErrInvalidRequest, "operator is already inactive or inactivating")
	}

	// Start the operator inactivation
	k.StartOperatorInactivation(ctx, operator)

	// Emit the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeStartOperatorInactivation,
			sdk.NewAttribute(types.AttributeKeyOperatorID, fmt.Sprintf("%d", msg.OperatorID)),
		),
	})

	return &types.MsgDeactivateOperatorResponse{}, nil
}

// ExecuteMessages defines the rpc method for Msg/ExecuteMessages
func (k msgServer) ExecuteMessages(goCtx context.Context, msg *types.MsgExecuteMessages) (*types.MsgExecuteMessagesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the operator exists
	operator, found := k.GetOperator(ctx, msg.OperatorID)
	if !found {
		return nil, types.ErrOperatorNotFound
	}

	// Make sure only the admin can execute messages
	if operator.Admin != msg.Sender {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "only the admin can execute messages")
	}

	operatorAddr, err := k.accountKeeper.AddressCodec().StringToBytes(operator.Address)
	if err != nil {
		return nil, err
	}

	messages, err := msg.GetMsgs()
	if err != nil {
		return nil, err
	}

	events := sdk.EmptyEvents()
	for _, msg := range messages {
		// perform a basic validation of the message
		if m, ok := msg.(sdk.HasValidateBasic); ok {
			if err := m.ValidateBasic(); err != nil {
				return nil, errors.Wrap(types.ErrInvalidExecuteMsg, err.Error())
			}
		}

		signers, _, err := k.cdc.GetMsgV1Signers(msg)
		if err != nil {
			return nil, err
		}
		if len(signers) != 1 {
			return nil, types.ErrInvalidExecuteMessagesSigner
		}

		// assert that the operator is the only signer for ExecuteMessages message
		if !bytes.Equal(signers[0], operatorAddr) {
			return nil, errors.Wrapf(types.ErrInvalidExecuteMessagesSigner, sdk.AccAddress(signers[0]).String())
		}

		handler := k.Router().Handler(msg)
		if handler == nil {
			return nil, errors.Wrap(types.ErrUnroutableExecuteMsg, sdk.MsgTypeURL(msg))
		}

		var res *sdk.Result
		res, err = handler(ctx, msg)
		if err != nil {
			return nil, err
		}

		events = append(events, res.GetEvents()...)
	}

	// TODO - merge events of MsgExecuteMessages itself
	ctx.EventManager().EmitEvents(events)

	return &types.MsgExecuteMessagesResponse{}, nil
}

// TransferOperatorOwnership defines the rpc method for Msg/TransferOperatorOwnership
func (k msgServer) TransferOperatorOwnership(goCtx context.Context, msg *types.MsgTransferOperatorOwnership) (*types.MsgTransferOperatorOwnershipResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the operator exists
	operator, found := k.GetOperator(ctx, msg.OperatorID)
	if !found {
		return nil, types.ErrOperatorNotFound
	}

	// Make sure only the admin can transfer the operator ownership
	if operator.Admin != msg.Sender {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "only the admin can transfer the operator ownership")
	}

	// Update the operator admin
	operator.Admin = msg.NewAdmin
	k.SaveOperator(ctx, operator)

	// Emit the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferOperatorOwnership,
			sdk.NewAttribute(types.AttributeKeyOperatorID, fmt.Sprintf("%d", msg.OperatorID)),
			sdk.NewAttribute(types.AttributeKeyNewAdmin, msg.NewAdmin),
		),
	})

	return &types.MsgTransferOperatorOwnershipResponse{}, nil
}

// UpdateParams defines the rpc method for Msg/UpdateParams
func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	// Check the authority
	authority := k.authority
	if authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", authority, msg.Authority)
	}

	// Update the params
	ctx := sdk.UnwrapSDKContext(goCtx)
	k.SetParams(ctx, msg.Params)

	return &types.MsgUpdateParamsResponse{}, nil
}
