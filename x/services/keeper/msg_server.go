package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/milkyway-labs/milkyway/x/services/types"
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

// CreateService defines the rpc method for Msg/CreateService
func (k msgServer) CreateService(goCtx context.Context, msg *types.MsgCreateService) (*types.MsgCreateServiceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the next service id
	avsID, err := k.GetNextServiceID(ctx)
	if err != nil {
		return nil, err
	}

	// Create the Service and validate it
	avs := types.NewService(
		avsID,
		types.SERVICE_STATUS_CREATED,
		msg.Name,
		msg.Description,
		msg.Website,
		msg.PictureURL,
		msg.Sender,
	)

	// Validate the service before storing
	err = avs.Validate()
	if err != nil {
		return nil, errors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
	}

	// Create the service
	err = k.Keeper.CreateService(ctx, avs)
	if err != nil {
		return nil, err
	}

	// Update the ID for the next service
	k.SetNextServiceID(ctx, avs.ID+1)

	// Emit the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreatedService,
			sdk.NewAttribute(types.AttributeKeyServiceID, fmt.Sprintf("%d", avs.ID)),
		),
	})

	return &types.MsgCreateServiceResponse{
		NewServiceID: avs.ID,
	}, nil
}

// UpdateService defines the rpc method for Msg/UpdateService
func (k msgServer) UpdateService(goCtx context.Context, msg *types.MsgUpdateService) (*types.MsgUpdateServiceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the service exists
	avs, found := k.GetService(ctx, msg.ServiceID)
	if !found {
		return nil, errors.Wrapf(sdkerrors.ErrInvalidRequest, "avs with id %d not found", msg.ServiceID)
	}

	// Make sure the user that is updating the service is the admin
	if avs.Admin != msg.Sender {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "you are not the admin of the service")
	}

	// Update the service
	updated := avs.Update(types.NewServiceUpdate(msg.Name, msg.Description, msg.Website, msg.PictureURL))

	// Validate the updated service
	err := updated.Validate()
	if err != nil {
		return nil, errors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
	}

	// Save the Service
	k.SaveService(ctx, updated)

	// Emit the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdatedService,
			sdk.NewAttribute(types.AttributeKeyServiceID, fmt.Sprintf("%d", msg.ServiceID)),
		),
	})

	return &types.MsgUpdateServiceResponse{}, nil
}

// DeactivateService defines the rpc method for Msg/DeactivateService
func (k msgServer) DeactivateService(goCtx context.Context, service *types.MsgDeactivateService) (*types.MsgDeactivateServiceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the service exists
	avs, found := k.GetService(ctx, service.ServiceID)
	if !found {
		return nil, errors.Wrapf(sdkerrors.ErrInvalidRequest, "service with id %d not found", service.ServiceID)
	}

	// Make sure the user that is deactivating the service is the admin
	if avs.Admin != service.Sender {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "you are not the admin of the service")
	}

	// Deactivate the service
	avs.Status = types.SERVICE_STATUS_INACTIVE

	// Save the Service
	k.SaveService(ctx, avs)

	// Emit the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDeactivatedService,
			sdk.NewAttribute(types.AttributeKeyServiceID, fmt.Sprintf("%d", service.ServiceID)),
		),
	})

	return &types.MsgDeactivateServiceResponse{}, nil
}
