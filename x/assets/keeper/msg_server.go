package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/milkyway-labs/milkyway/v7/x/assets/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	*Keeper
}

func NewMsgServer(k *Keeper) types.MsgServer {
	return &msgServer{Keeper: k}
}

// RegisterAsset defines the rpc method for Msg/RegisterAsset
func (k msgServer) RegisterAsset(ctx context.Context, msg *types.MsgRegisterAsset) (*types.MsgRegisterAssetResponse, error) {
	err := msg.Validate()
	if err != nil {
		return nil, err
	}

	// Check if the authority is correct
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	// Store the asset
	err = k.SetAsset(ctx, msg.Asset)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRegisterAsset,
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Asset.Denom),
			sdk.NewAttribute(types.AttributeKeyTicker, msg.Asset.Ticker),
			sdk.NewAttribute(types.AttributeKeyExponent, fmt.Sprint(msg.Asset.Exponent)),
		),
	})

	return &types.MsgRegisterAssetResponse{}, nil
}

// DeregisterAsset defines the rpc method for Msg/DeregisterAsset
func (k msgServer) DeregisterAsset(ctx context.Context, msg *types.MsgDeregisterAsset) (*types.MsgDeregisterAssetResponse, error) {
	// Validate the message
	err := msg.Validate()
	if err != nil {
		return nil, err
	}

	// Check if the authority is correct
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	// Remove the asset
	err = k.RemoveAsset(ctx, msg.Denom)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDeregisterAsset,
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
		),
	})

	return &types.MsgDeregisterAssetResponse{}, nil
}
