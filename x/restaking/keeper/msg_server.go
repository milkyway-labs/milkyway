package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/hashicorp/go-metrics"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

var (
	_ types.MsgServer = msgServer{}
)

type msgServer struct {
	*Keeper
}

func NewMsgServer(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// DelegatePool defines the rpc method for Msg/DelegatePool
func (k msgServer) DelegatePool(goCtx context.Context, msg *types.MsgDelegatePool) (*types.MsgDelegatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return nil, errors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid delegation amount",
		)
	}

	newShares, err := k.Keeper.DelegateToPool(ctx, msg.Amount, msg.Delegator)
	if err != nil {
		return nil, err
	}

	if msg.Amount.Amount.IsInt64() {
		defer func() {
			telemetry.IncrCounter(1, types.ModuleName, "pool restake")
			telemetry.SetGaugeWithLabels(
				[]string{"tx", "msg", sdk.MsgTypeURL(msg)},
				float32(msg.Amount.Amount.Int64()),
				[]metrics.Label{telemetry.NewLabel("denom", msg.Amount.Denom)},
			)
		}()
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDelegatePool,
			sdk.NewAttribute(types.AttributeKeyDelegator, msg.Delegator),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyNewShares, newShares.String()),
		),
	})

	return &types.MsgDelegatePoolResponse{}, nil
}

// DelegateOperator defines the rpc method for Msg/DelegateOperator
func (k msgServer) DelegateOperator(goCtx context.Context, msg *types.MsgDelegateOperator) (*types.MsgDelegateOperatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !msg.Amount.IsValid() || !msg.Amount.IsAllPositive() {
		return nil, errors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid delegation amount",
		)
	}

	newShares, err := k.Keeper.DelegateToOperator(ctx, msg.OperatorID, msg.Amount, msg.Delegator)
	if err != nil {
		return nil, err
	}

	for _, token := range msg.Amount {
		if token.Amount.IsInt64() {
			defer func() {
				telemetry.IncrCounter(1, types.ModuleName, "operator restake")
				telemetry.SetGaugeWithLabels(
					[]string{"tx", "msg", sdk.MsgTypeURL(msg)},
					float32(token.Amount.Int64()),
					[]metrics.Label{telemetry.NewLabel("denom", token.Denom)},
				)
			}()
		}
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDelegateOperator,
			sdk.NewAttribute(types.AttributeKeyDelegator, msg.Delegator),
			sdk.NewAttribute(types.AttributeKeyOperatorID, fmt.Sprintf("%d", msg.OperatorID)),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyNewShares, newShares.String()),
		),
	})

	return &types.MsgDelegateOperatorResponse{}, nil
}

// DelegateService defines the rpc method for Msg/DelegateService
func (k msgServer) DelegateService(goCtx context.Context, msg *types.MsgDelegateService) (*types.MsgDelegateServiceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !msg.Amount.IsValid() || !msg.Amount.IsAllPositive() {
		return nil, errors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid delegation amount",
		)
	}

	newShares, err := k.Keeper.DelegateToService(ctx, msg.ServiceID, msg.Amount, msg.Delegator)
	if err != nil {
		return nil, err
	}

	for _, token := range msg.Amount {
		if token.Amount.IsInt64() {
			defer func() {
				telemetry.IncrCounter(1, types.ModuleName, "service restake")
				telemetry.SetGaugeWithLabels(
					[]string{"tx", "msg", sdk.MsgTypeURL(msg)},
					float32(token.Amount.Int64()),
					[]metrics.Label{telemetry.NewLabel("denom", token.Denom)},
				)
			}()
		}
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDelegateService,
			sdk.NewAttribute(types.AttributeKeyDelegator, msg.Delegator),
			sdk.NewAttribute(types.AttributeKeyServiceID, fmt.Sprintf("%d", msg.ServiceID)),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyNewShares, newShares.String()),
		),
	})

	return &types.MsgDelegateServiceResponse{}, nil
}

// UpdateParams defines the rpc method for Msg/UpdateParams
func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	// Check the authority
	authority := k.authority
	if authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", authority, msg.Authority)
	}

	// Update the params
	ctx := sdk.UnwrapSDKContext(goCtx)
	k.SetParams(ctx, msg.Params)

	return &types.MsgUpdateParamsResponse{}, nil
}