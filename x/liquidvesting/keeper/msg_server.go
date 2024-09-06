package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
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

// MintStakingRepresentation implements types.MsgServer.
func (m msgServer) MintStakingRepresentation(
	goCtx context.Context,
	msg *types.MsgMintStakingRepresentation,
) (*types.MsgMintStakingRepresentationResponse, error) {
	accAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return nil, err
	}

	isMinter, err := m.IsMinter(goCtx, accAddr)
	if !isMinter {
		return nil, types.ErrNotMinter
	}

	err = m.Keeper.MintStakingRepresentation(goCtx, receiver, msg.Amount)
	if err != nil {
		return nil, err
	}

	return &types.MsgMintStakingRepresentationResponse{}, nil
}

// BurnStakingRepresentation implements types.MsgServer.
func (m msgServer) BurnStakingRepresentation(
	goCtx context.Context,
	msg *types.MsgBurnStakingRepresentation,
) (*types.MsgBurnStakingRepresentationResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	user, err := sdk.AccAddressFromBech32(msg.User)
	if err != nil {
		return nil, err
	}

	isBurner, err := m.IsBurner(goCtx, sender)
	if !isBurner {
		return nil, types.ErrNotBurner
	}

	err = m.Keeper.BurnStakingRepresentation(user, msg.Amount)
	if err != nil {
		return nil, err
	}

	return &types.MsgBurnStakingRepresentationResponse{}, nil
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
	m.SetParams(ctx, msg.Params)

	return &types.MsgUpdateParamsResponse{}, nil
}
