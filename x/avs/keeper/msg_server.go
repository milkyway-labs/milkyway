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

// RegisterAVS defines the rpc method for Msg/RegisterAVS
func (k msgServer) RegisterAVS(goCtx context.Context, msg *types.MsgRegisterAVS) (*types.MsgRegisterAVSResponse, error) {
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
			types.EventTypeRegisteredAVS,
			sdk.NewAttribute(types.AttributeKeyAVSID, fmt.Sprintf("%d", avs.ID)),
		),
	})

	return &types.MsgRegisterAVSResponse{
		NewAVSID: avs.ID,
	}, nil
}

// UpdateAVS defines the rpc method for Msg/UpdateAVS
func (k msgServer) UpdateAVS(goCtx context.Context, msg *types.MsgUpdateAVS) (*types.MsgUpdateAVSResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the AVS exists
	avs, found := k.GetAVS(ctx, msg.AVSID)
	if !found {
		return nil, errors.Wrapf(sdkerrors.ErrInvalidRequest, "avs with id %d not found", msg.AVSID)
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
			types.EventTypeUpdatedAVS,
			sdk.NewAttribute(types.AttributeKeyAVSID, fmt.Sprintf("%d", msg.AVSID)),
		),
	})

	return &types.MsgUpdateAVSResponse{}, nil
}
