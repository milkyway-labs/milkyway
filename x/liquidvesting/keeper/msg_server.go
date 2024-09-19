package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	*Keeper
}

func NewMsgServer(k *Keeper) types.MsgServer {
	return &msgServer{Keeper: k}
}

// MintVestedRepresentation implements types.MsgServer.
func (m msgServer) MintVestedRepresentation(
	goCtx context.Context,
	msg *types.MsgMintVestedRepresentation,
) (*types.MsgMintVestedRepresentationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return nil, err
	}

	isMinter, err := m.IsMinter(ctx, sender)
	if !isMinter {
		return nil, types.ErrNotMinter
	}

	err = m.Keeper.MintVestedRepresentation(ctx, receiver, msg.Amount)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintVestedRepresentation,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, msg.Receiver),
		),
	})

	return &types.MsgMintVestedRepresentationResponse{}, nil
}

// BurnVestedRepresentation implements types.MsgServer.
func (m msgServer) BurnVestedRepresentation(
	goCtx context.Context,
	msg *types.MsgBurnVestedRepresentation,
) (*types.MsgBurnVestedRepresentationResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	user, err := sdk.AccAddressFromBech32(msg.User)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	isBurner, err := m.IsBurner(ctx, sender)
	if !isBurner {
		return nil, types.ErrNotBurner
	}

	err = m.Keeper.BurnVestedRepresentation(ctx, user, msg.Amount)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnVestedRepresentation,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyUser, msg.User),
		),
	})

	return &types.MsgBurnVestedRepresentationResponse{}, nil
}

// UpdateParams implements types.MsgServer.
func (m msgServer) UpdateParams(
	goCtx context.Context,
	msg *types.MsgUpdateParams,
) (*types.MsgUpdateParamsResponse, error) {
	// Check the authority
	authority := m.authority
	if authority != msg.Authority {
		return nil, errors.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority; expected %s, got %s",
			authority, msg.Authority,
		)
	}

	// Update the params
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := m.SetParams(ctx, msg.Params)
	if err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
