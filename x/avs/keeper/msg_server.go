package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/milkyway-labs/milkyway/x/avs/types"
)

var (
	_ types.MsgServer = msgServer{}
)

type msgServer struct {
	Keeper
}

// RegisterService defines the rpc method for Msg/RegisterService
func (k msgServer) RegisterService(goCtx context.Context, msg *types.MsgRegisterService) (*types.MsgRegisterServiceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the next reaction id
	avsID, err := k.GetNextAVSID(ctx)
	if err != nil {
		return nil, err
	}

	// Create the AVS and validate it
	avs := types.NewAVS(
		avsID,
		types.AVS_STATUS_CREATED,
		msg.Name,
		msg.Description,
		msg.Website,
		msg.PictureURL,
		msg.Sender,
	)
	if err := avs.Validate(); err != nil {
		return nil, errors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
	}

	// Store the AVS
	k.SaveAVS(ctx, avs)

	// Update the ID for the next AVS
	k.SetNextAVSID(ctx, avs.ID+1)

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

	// Check if the AVS exists
	avs, found := k.GetAVS(ctx, msg.ServiceID)
	if !found {
		return nil, errors.Wrapf(sdkerrors.ErrInvalidRequest, "avs with id %d not found", msg.ServiceID)
	}

	// Update the AVS and validate it
	updated := avs.Update(types.NewAVSUpdate(msg.Name, msg.Description, msg.Website, msg.PictureURL))
	if err := updated.Validate(); err != nil {
		return nil, errors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
	}

	// Save the AVS
	k.SaveAVS(ctx, updated)

	// Emit the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdatedService,
			sdk.NewAttribute(types.AttributeKeyServiceID, fmt.Sprintf("%d", msg.ServiceID)),
		),
	})

	return &types.MsgUpdateServiceResponse{}, nil
}
