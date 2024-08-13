package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/milkyway-labs/milkyway/x/assets/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	*Keeper
}

func NewMsgServer(k *Keeper) types.MsgServer {
	return &msgServer{Keeper: k}
}

func (k msgServer) RegisterAsset(ctx context.Context, msg *types.MsgRegisterAsset) (*types.MsgRegisterAssetResponse, error) {
	if err := msg.Validate(); err != nil {
		return nil, err
	}

	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	if err := k.SetAsset(ctx, msg.Asset); err != nil {
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

func (k msgServer) DeregisterAsset(ctx context.Context, msg *types.MsgDeregisterAsset) (*types.MsgDeregisterAssetResponse, error) {
	if err := msg.Validate(); err != nil {
		return nil, err
	}

	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	if err := k.RemoveAsset(ctx, msg.Denom); err != nil {
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

// UpdateParams defines the rpc method for Msg/UpdateParams
func (k msgServer) UpdateParams(ctx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if err := msg.Validate(); err != nil {
		return nil, err
	}

	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	// store params
	if err := k.Params.Set(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
