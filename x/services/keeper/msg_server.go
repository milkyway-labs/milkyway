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
	serviceID, err := k.GetNextServiceID(ctx)
	if err != nil {
		return nil, err
	}

	// Create the Service and validate it
	service := types.NewService(
		serviceID,
		types.SERVICE_STATUS_CREATED,
		msg.Name,
		msg.Description,
		msg.Website,
		msg.PictureURL,
		msg.Sender,
	)

	// Validate the service before storing
	err = service.Validate()
	if err != nil {
		return nil, errors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
	}

	// Create the service
	err = k.Keeper.CreateService(ctx, service)
	if err != nil {
		return nil, err
	}

	// Update the ID for the next service
	k.SetNextServiceID(ctx, service.ID+1)

	// Emit the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreatedService,
			sdk.NewAttribute(types.AttributeKeyServiceID, fmt.Sprintf("%d", service.ID)),
		),
	})

	return &types.MsgCreateServiceResponse{
		NewServiceID: service.ID,
	}, nil
}

// UpdateService defines the rpc method for Msg/UpdateService
func (k msgServer) UpdateService(goCtx context.Context, msg *types.MsgUpdateService) (*types.MsgUpdateServiceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the service exists
	service, found := k.GetService(ctx, msg.ServiceID)
	if !found {
		return nil, errors.Wrapf(sdkerrors.ErrInvalidRequest, "service with id %d not found", msg.ServiceID)
	}

	// Make sure the user that is updating the service is the admin
	if service.Admin != msg.Sender {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "you are not the admin of the service")
	}

	// Update the service
	updated := service.Update(types.NewServiceUpdate(msg.Name, msg.Description, msg.Website, msg.PictureURL))

	// Validate the updated service
	err := updated.Validate()
	if err != nil {
		return nil, errors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
	}

	// Save the service
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
func (k msgServer) DeactivateService(goCtx context.Context, msg *types.MsgDeactivateService) (*types.MsgDeactivateServiceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the service exists
	service, found := k.GetService(ctx, msg.ServiceID)
	if !found {
		return nil, errors.Wrapf(sdkerrors.ErrInvalidRequest, "service with id %d not found", msg.ServiceID)
	}

	// Make sure the user that is deactivating the service is the admin
	if service.Admin != msg.Sender {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "you are not the admin of the service")
	}

	// Deactivate the service
	service.Status = types.SERVICE_STATUS_INACTIVE

	// Save the service
	k.SaveService(ctx, service)

	// Emit the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDeactivatedService,
			sdk.NewAttribute(types.AttributeKeyServiceID, fmt.Sprintf("%d", msg.ServiceID)),
		),
	})

	return &types.MsgDeactivateServiceResponse{}, nil
}
