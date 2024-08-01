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

func (k msgServer) WithdrawPoolDelReward(ctx context.Context, msg *types.MsgWithdrawPoolDelReward) (*types.MsgWithdrawPoolDelRewardResponse, error) {
	amount, err := k.WithdrawPoolDelegationRewards(ctx, msg.DelegatorAddress, msg.PoolID)
	if err != nil {
		return nil, err
	}

	// TODO: telemetry?

	return &types.MsgWithdrawPoolDelRewardResponse{Amount: amount}, nil
}

func (k msgServer) WithdrawOperatorDelReward(ctx context.Context, msg *types.MsgWithdrawOperatorDelReward) (*types.MsgWithdrawOperatorDelRewardResponse, error) {
	rewards, err := k.WithdrawOperatorDelegationRewards(ctx, msg.DelegatorAddress, msg.OperatorID)
	if err != nil {
		return nil, err
	}

	amount := rewards.Sum()
	// TODO: telemetry?

	return &types.MsgWithdrawOperatorDelRewardResponse{Amount: amount}, nil
}

func (k msgServer) WithdrawServiceDelReward(ctx context.Context, msg *types.MsgWithdrawServiceDelReward) (*types.MsgWithdrawServiceDelRewardResponse, error) {
	rewards, err := k.WithdrawServiceDelegationRewards(ctx, msg.DelegatorAddress, msg.ServiceID)
	if err != nil {
		return nil, err
	}

	amount := rewards.Sum()
	// TODO: telemetry?

	return &types.MsgWithdrawServiceDelRewardResponse{Amount: amount}, nil
}

func (k msgServer) WithdrawOperatorCommission(ctx context.Context, msg *types.MsgWithdrawOperatorCommission) (*types.MsgWithdrawOperatorCommissionResponse, error) {
	commissions, err := k.Keeper.WithdrawOperatorCommission(ctx, msg.OperatorID)
	if err != nil {
		return nil, err
	}

	amount := commissions.Sum()
	// TODO: telemetry?

	return &types.MsgWithdrawOperatorCommissionResponse{Amount: amount}, nil
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
