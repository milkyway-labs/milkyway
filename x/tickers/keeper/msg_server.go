package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/milkyway-labs/milkyway/x/tickers/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	*Keeper
}

func NewMsgServer(k *Keeper) types.MsgServer {
	return &msgServer{Keeper: k}
}

func (k msgServer) RegisterTicker(ctx context.Context, msg *types.MsgRegisterTicker) (*types.MsgRegisterTickerResponse, error) {
	if err := msg.Validate(); err != nil {
		return nil, err
	}

	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	if err := k.SetTicker(ctx, msg.Denom, msg.Ticker); err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRegisterTicker,
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyTicker, msg.Ticker),
		),
	})

	return &types.MsgRegisterTickerResponse{}, nil
}

func (k msgServer) DeregisterTicker(ctx context.Context, msg *types.MsgDeregisterTicker) (*types.MsgDeregisterTickerResponse, error) {
	if err := msg.Validate(); err != nil {
		return nil, err
	}

	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	if err := k.RemoveTicker(ctx, msg.Denom); err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDeregisterTicker,
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
		),
	})

	return &types.MsgDeregisterTickerResponse{}, nil
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
