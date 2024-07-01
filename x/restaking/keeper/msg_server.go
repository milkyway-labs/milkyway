package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// PoolRestake defines the rpc method for Msg/PoolRestake
func (k msgServer) PoolRestake(goCtx context.Context, msg *types.MsgJoinRestakingPool) (*types.MsgJoinRestakingPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

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
			types.EventTypePoolRestake,
			sdk.NewAttribute(types.AttributeKeyDelegator, msg.Delegator),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyNewShares, newShares.String()),
		),
	})

	return &types.MsgJoinRestakingPoolResponse{}, nil
}

// OperatorRestake defines the rpc method for Msg/OperatorRestake
func (k msgServer) OperatorRestake(goCtx context.Context, msg *types.MsgDelegateOperator) (*types.MsgDelegateOperatorResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ServiceRestake defines the rpc method for Msg/ServiceRestake
func (k msgServer) ServiceRestake(goCtx context.Context, msg *types.MsgDelegateService) (*types.MsgDelegateServiceResponse, error) {
	//TODO implement me
	panic("implement me")
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
