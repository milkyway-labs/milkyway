package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/milkyway-labs/milkyway/v11/x/investors/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	*Keeper
}

func NewMsgServer(k *Keeper) types.MsgServer {
	return &msgServer{Keeper: k}
}

// UpdateInvestorsRewardRatio defines the rpc method for Msg/UpdateInvestorsRewardRatio
func (k msgServer) UpdateInvestorsRewardRatio(ctx context.Context, msg *types.MsgUpdateInvestorsRewardRatio) (*types.MsgUpdateInvestorsRewardRatioResponse, error) {
	// Check the authority
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	err := k.Keeper.UpdateInvestorsRewardRatio(ctx, msg.InvestorsRewardRatio)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateInvestorsRewardRatio,
			sdk.NewAttribute(types.AttributeKeyInvestorsRewardRatio, msg.InvestorsRewardRatio.String()),
		),
	})

	return &types.MsgUpdateInvestorsRewardRatioResponse{}, nil
}
