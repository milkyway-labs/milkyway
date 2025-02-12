package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/milkyway-labs/milkyway/v9/x/liquidvesting/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	*Keeper
}

func NewMsgServer(k *Keeper) types.MsgServer {
	return &msgServer{Keeper: k}
}

// MintLockedRepresentation implements types.MsgServer.
func (m msgServer) MintLockedRepresentation(ctx context.Context, msg *types.MsgMintLockedRepresentation) (*types.MsgMintLockedRepresentationResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return nil, err
	}

	isMinter, err := m.IsMinter(ctx, sender)
	if err != nil {
		return nil, err
	}

	if !isMinter {
		return nil, types.ErrNotMinter
	}

	mintedAmount, err := m.Keeper.MintLockedRepresentation(ctx, receiver, msg.Amount)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintLockedRepresentation,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(sdk.AttributeKeyAmount, mintedAmount.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, msg.Receiver),
		),
	})

	return &types.MsgMintLockedRepresentationResponse{}, nil
}

// BurnLockedRepresentation implements types.MsgServer.
func (m msgServer) BurnLockedRepresentation(ctx context.Context, msg *types.MsgBurnLockedRepresentation) (*types.MsgBurnLockedRepresentationResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	user, err := sdk.AccAddressFromBech32(msg.User)
	if err != nil {
		return nil, err
	}

	isBurner, err := m.IsBurner(ctx, sender)
	if err != nil {
		return nil, err
	}

	if !isBurner {
		return nil, types.ErrNotBurner
	}

	err = m.Keeper.BurnLockedRepresentation(ctx, user, msg.Amount)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnLockedRepresentation,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyUser, msg.User),
		),
	})

	return &types.MsgBurnLockedRepresentationResponse{}, nil
}

// WithdrawInsuranceFund implements types.MsgServer.
func (m msgServer) WithdrawInsuranceFund(ctx context.Context, msg *types.MsgWithdrawInsuranceFund) (*types.MsgWithdrawInsuranceFundResponse, error) {
	canWithdraw, err := m.CanWithdrawFromInsuranceFund(ctx, msg.Sender, msg.Amount)
	if err != nil {
		return nil, err
	}

	if !canWithdraw {
		return nil, types.ErrInsufficientBalance
	}

	// Send the tokens back to the user
	err = m.WithdrawFromUserInsuranceFund(ctx, msg.Sender, msg.Amount)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawInsuranceFund,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
	})

	return &types.MsgWithdrawInsuranceFundResponse{}, nil
}

// UpdateParams implements types.MsgServer.
func (m msgServer) UpdateParams(ctx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	// Check the authority
	authority := m.authority
	if authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", authority, msg.Authority)
	}

	// Update the params
	err := m.SetParams(ctx, msg.Params)
	if err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
