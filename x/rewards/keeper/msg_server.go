package keeper

import (
	"context"

	"cosmossdk.io/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/milkyway-labs/milkyway/x/rewards/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	*Keeper
}

func NewMsgServer(k *Keeper) types.MsgServer {
	return &msgServer{Keeper: k}
}

// CreateRewardsPlan defines the rpc method for Msg/CreateRewardsPlan
func (k msgServer) CreateRewardsPlan(ctx context.Context, msg *types.MsgCreateRewardsPlan) (*types.MsgCreateRewardsPlanResponse, error) {
	// TODO: need to charge fee?

	plan, err := k.Keeper.CreateRewardsPlan(
		ctx, msg.Description, msg.ServiceID, msg.Amount, msg.StartTime, msg.EndTime, msg.PoolsDistribution, msg.OperatorsDistribution,
		msg.UsersDistribution)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateRewardsPlanResponse{NewRewardsPlanID: plan.ID}, nil
}

// UpdateParams defines the rpc method for Msg/UpdateParams
func (k msgServer) UpdateParams(ctx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	if err := msg.Params.Validate(); err != nil {
		return nil, err
	}

	// store params
	if err := k.Params.Set(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
