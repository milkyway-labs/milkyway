package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

func (k msgServer) SetWithdrawAddress(ctx context.Context, msg *types.MsgSetWithdrawAddress) (*types.MsgSetWithdrawAddressResponse, error) {
	delegatorAddress, err := k.accountKeeper.AddressCodec().StringToBytes(msg.DelegatorAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}

	withdrawAddress, err := k.accountKeeper.AddressCodec().StringToBytes(msg.WithdrawAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid withdraw address: %s", err)
	}

	err = k.SetWithdrawAddr(ctx, delegatorAddress, withdrawAddress)
	if err != nil {
		return nil, err
	}

	return &types.MsgSetWithdrawAddressResponse{}, nil
}

func (k msgServer) WithdrawDelegatorReward(ctx context.Context, msg *types.MsgWithdrawDelegatorReward) (*types.MsgWithdrawDelegatorRewardResponse, error) {
	delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(msg.DelegatorAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}

	target, err := k.GetDelegationTarget(ctx, msg.DelegationType, msg.TargetID)
	if err != nil {
		return nil, err
	}

	rewards, err := k.WithdrawDelegationRewards(ctx, delAddr, target)
	if err != nil {
		return nil, err
	}
	amount := rewards.Sum()

	// TODO: telemetry?

	return &types.MsgWithdrawDelegatorRewardResponse{Amount: amount}, nil
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
