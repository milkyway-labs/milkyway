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
	Keeper
}

func NewMsgServer(k Keeper) types.MsgServer {
	return &msgServer{Keeper: k}
}

// RegisterService defines the rpc method for Msg/RegisterService
func (k msgServer) RegisterService(goCtx context.Context, msg *types.MsgRegisterService) (*types.MsgRegisterServiceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the next reaction id
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

	// Create the Service
	err = k.CreateService(ctx, avs)
	if err != nil {
		return nil, err
	}

	// Update the ID for the next Service
	k.SetNextServiceID(ctx, avs.ID+1)

	// Emit the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRegisteredService,
			sdk.NewAttribute(types.AttributeKeyServiceID, fmt.Sprintf("%d", avs.ID)),
		),
	})

	return &types.MsgRegisterServiceResponse{
		NewServiceID: avs.ID,
	}, nil
}

// UpdateService defines the rpc method for Msg/UpdateService
func (k msgServer) UpdateService(goCtx context.Context, msg *types.MsgUpdateService) (*types.MsgUpdateServiceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the Service exists
	avs, found := k.GetService(ctx, msg.ServiceID)
	if !found {
		return nil, errors.Wrapf(sdkerrors.ErrInvalidRequest, "avs with id %d not found", msg.ServiceID)
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
