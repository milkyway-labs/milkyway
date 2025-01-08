package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/milkyway-labs/milkyway/v7/x/rewards/types"
	servicestypes "github.com/milkyway-labs/milkyway/v7/x/services/types"
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
	// Make sure the creator is the admin of the service
	service, err := k.servicesKeeper.GetService(ctx, msg.ServiceID)
	if err != nil {
		return nil, err
	}

	if msg.Sender != service.Admin {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "only service admin can create rewards plan")
	}

	// Charge a scaling gas consumption fee
	rewardsPlans, err := k.GetRewardsPlans(ctx)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.GasMeter().ConsumeGas(types.BaseGasFeeForNewPlan*uint64(len(rewardsPlans)), "create rewards plan gas cost")

	// Create the plan
	plan, err := k.Keeper.CreateRewardsPlan(
		ctx,
		msg.Description,
		msg.ServiceID,
		msg.Amount,
		msg.StartTime,
		msg.EndTime,
		msg.PoolsDistribution,
		msg.OperatorsDistribution,
		msg.UsersDistribution,
	)
	if err != nil {
		return nil, err
	}

	// Charge fee for rewards plan creation. Fee is charged only in msg server and
	// not when calling the keeper's method directly. This gives freedom to other
	// modules to call the keeper's method directly without charging the fee.
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	if !params.RewardsPlanCreationFee.IsZero() {
		// Make sure the specified fees are enough
		if !msg.FeeAmount.IsAnyGTE(params.RewardsPlanCreationFee) {
			return nil, errors.Wrapf(sdkerrors.ErrInsufficientFunds, "insufficient funds: %s < %s", msg.FeeAmount, params.RewardsPlanCreationFee)
		}

		userAddress, err := sdk.AccAddressFromBech32(service.Admin)
		if err != nil {
			return nil, errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %s", service.Admin)
		}

		err = k.communityPoolKeeper.FundCommunityPool(ctx, msg.FeeAmount, userAddress)
		if err != nil {
			return nil, err
		}
	}

	// Emit the event
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateRewardsPlan,
			sdk.NewAttribute(types.AttributeKeyRewardsPlanID, fmt.Sprint(plan.ID)),
			sdk.NewAttribute(servicestypes.AttributeKeyServiceID, fmt.Sprint(msg.ServiceID)),
		),
	})

	return &types.MsgCreateRewardsPlanResponse{NewRewardsPlanID: plan.ID}, nil
}

// EditRewardsPlan defines the rpc method for Msg/EditRewardsPlan
func (k msgServer) EditRewardsPlan(ctx context.Context, msg *types.MsgEditRewardsPlan) (*types.MsgEditRewardsPlanResponse, error) {
	// Get the rewards plan to edit
	rewardsPlan, err := k.GetRewardsPlan(ctx, msg.ID)
	if err != nil {
		return nil, err
	}

	// Get the service to which the rewards is associated
	service, err := k.servicesKeeper.GetService(ctx, rewardsPlan.ServiceID)
	if err != nil {
		return nil, err
	}

	// Make sure the editor is the admin of the service
	if msg.Sender != service.Admin {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "only service admin can create rewards plan")
	}

	// Edit the rewards plan
	err = k.Keeper.EditRewardsPlan(
		ctx,
		msg.ID,
		msg.Description,
		msg.Amount,
		msg.StartTime,
		msg.EndTime,
		msg.PoolsDistribution,
		msg.OperatorsDistribution,
		msg.UsersDistribution,
	)
	if err != nil {
		return nil, err
	}

	// Emit the event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEditRewardsPlan,
			sdk.NewAttribute(types.AttributeKeyRewardsPlanID, fmt.Sprint(rewardsPlan.ID)),
			sdk.NewAttribute(servicestypes.AttributeKeyServiceID, fmt.Sprint(rewardsPlan.ServiceID)),
		),
	})

	return &types.MsgEditRewardsPlanResponse{}, nil
}

// SetWithdrawAddress sets the withdraw address for a delegator(or an operator
// when withdrawing commission). The default withdraw address if not set
// specified is the delegator(or an operator) address.
func (k msgServer) SetWithdrawAddress(ctx context.Context, msg *types.MsgSetWithdrawAddress) (*types.MsgSetWithdrawAddressResponse, error) {
	// Parse the addresses
	senderAddr, err := k.accountKeeper.AddressCodec().StringToBytes(msg.Sender)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}

	withdrawAddress, err := k.accountKeeper.AddressCodec().StringToBytes(msg.WithdrawAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid withdraw address: %s", err)
	}

	// Set the withdraw address
	err = k.Keeper.SetWithdrawAddress(ctx, senderAddr, withdrawAddress)
	if err != nil {
		return nil, err
	}

	// Emit an event
	// Emit the event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSetWithdrawAddress,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyWithdrawAddress, msg.WithdrawAddress),
		),
	)

	return &types.MsgSetWithdrawAddressResponse{}, nil
}

// WithdrawDelegatorReward defines the rpc method Msg/WithdrawDelegatorReward
func (k msgServer) WithdrawDelegatorReward(ctx context.Context, msg *types.MsgWithdrawDelegatorReward) (*types.MsgWithdrawDelegatorRewardResponse, error) {
	delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(msg.DelegatorAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}

	if msg.DelegationTargetID == 0 {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("invalid delegation target ID: %d", msg.DelegationTargetID)
	}

	target, err := k.GetDelegationTarget(ctx, msg.DelegationType, msg.DelegationTargetID)
	if err != nil {
		return nil, err
	}

	rewards, err := k.WithdrawDelegationRewards(ctx, delAddr, target)
	if err != nil {
		return nil, err
	}

	return &types.MsgWithdrawDelegatorRewardResponse{Amount: rewards.Sum()}, nil
}

// WithdrawOperatorCommission defines the rpc method Msg/WithdrawOperatorCommission
func (k msgServer) WithdrawOperatorCommission(ctx context.Context, msg *types.MsgWithdrawOperatorCommission) (*types.MsgWithdrawOperatorCommissionResponse, error) {
	_, err := k.accountKeeper.AddressCodec().StringToBytes(msg.Sender)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}

	operator, err := k.operatorsKeeper.GetOperator(ctx, msg.OperatorID)
	if err != nil {
		return nil, err
	}

	if msg.Sender != operator.Admin {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "only operator admin can withdraw operator commission")
	}

	commissions, err := k.Keeper.WithdrawOperatorCommission(ctx, msg.OperatorID)
	if err != nil {
		return nil, err
	}

	return &types.MsgWithdrawOperatorCommissionResponse{Amount: commissions.Sum()}, nil
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
