package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

	restakingtypes "github.com/milkyway-labs/milkyway/v4/x/restaking/types"
	"github.com/milkyway-labs/milkyway/v4/x/rewards/types"
)

var _ types.QueryServer = queryServer{}

type queryServer struct {
	k *Keeper
}

func NewQueryServer(k *Keeper) types.QueryServer {
	return queryServer{k: k}
}

func (q queryServer) Params(ctx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params, err := q.k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	return &types.QueryParamsResponse{Params: params}, nil
}

func (q queryServer) RewardsPlans(ctx context.Context, req *types.QueryRewardsPlansRequest) (*types.QueryRewardsPlansResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	plans, pageRes, err := query.CollectionPaginate(ctx, q.k.RewardsPlans, req.Pagination, func(_ uint64, plan types.RewardsPlan) (types.RewardsPlan, error) {
		return plan, nil
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryRewardsPlansResponse{RewardsPlans: plans, Pagination: pageRes}, nil
}

func (q queryServer) RewardsPlan(ctx context.Context, req *types.QueryRewardsPlanRequest) (*types.QueryRewardsPlanResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.PlanId == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid plan id")
	}

	plan, err := q.k.RewardsPlans.Get(ctx, req.PlanId)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "plan not found")
		}
		return nil, err
	}

	return &types.QueryRewardsPlanResponse{RewardsPlan: plan}, nil
}

func (q queryServer) PoolOutstandingRewards(ctx context.Context, req *types.QueryPoolOutstandingRewardsRequest) (*types.QueryPoolOutstandingRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.PoolId == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid pool id")
	}

	target, err := q.k.GetDelegationTarget(ctx, restakingtypes.DELEGATION_TYPE_POOL, req.PoolId)
	if err != nil {
		return nil, err
	}

	rewards, err := target.OutstandingRewards.Get(ctx, target.GetID())
	if err != nil && !errors.IsOf(err, collections.ErrNotFound) {
		return nil, err
	}

	return &types.QueryPoolOutstandingRewardsResponse{Rewards: rewards}, nil
}

func (q queryServer) OperatorOutstandingRewards(ctx context.Context, req *types.QueryOperatorOutstandingRewardsRequest) (*types.QueryOperatorOutstandingRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.OperatorId == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid operator id")
	}

	target, err := q.k.GetDelegationTarget(ctx, restakingtypes.DELEGATION_TYPE_OPERATOR, req.OperatorId)
	if err != nil {
		return nil, err
	}

	rewards, err := target.OutstandingRewards.Get(ctx, target.GetID())
	if err != nil && !errors.IsOf(err, collections.ErrNotFound) {
		return nil, err
	}

	return &types.QueryOperatorOutstandingRewardsResponse{Rewards: rewards}, nil
}

func (q queryServer) ServiceOutstandingRewards(ctx context.Context, req *types.QueryServiceOutstandingRewardsRequest) (*types.QueryServiceOutstandingRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.ServiceId == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid service id")
	}

	target, err := q.k.GetDelegationTarget(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, req.ServiceId)
	if err != nil {
		return nil, err
	}

	rewards, err := target.OutstandingRewards.Get(ctx, target.GetID())
	if err != nil && !errors.IsOf(err, collections.ErrNotFound) {
		return nil, err
	}

	return &types.QueryServiceOutstandingRewardsResponse{Rewards: rewards}, nil
}

func (q queryServer) OperatorCommission(ctx context.Context, req *types.QueryOperatorCommissionRequest) (*types.QueryOperatorCommissionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.OperatorId == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid operator id")
	}

	commission, err := q.k.GetOperatorAccumulatedCommission(ctx, req.OperatorId)
	if err != nil {
		return nil, err
	}

	return &types.QueryOperatorCommissionResponse{Commission: commission}, nil
}

func (q queryServer) PoolDelegationRewards(ctx context.Context, req *types.QueryPoolDelegationRewardsRequest) (*types.QueryPoolDelegationRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	delAddr, err := q.k.accountKeeper.AddressCodec().StringToBytes(req.DelegatorAddress)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid delegator address: %s", err)
	}

	if req.PoolId == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid pool id")
	}

	rewards, err := q.k.GetPoolDelegationRewards(ctx, delAddr, req.PoolId)
	if err != nil {
		return nil, err
	}

	return &types.QueryPoolDelegationRewardsResponse{Rewards: rewards}, nil
}

func (q queryServer) OperatorDelegationRewards(ctx context.Context, req *types.QueryOperatorDelegationRewardsRequest) (*types.QueryOperatorDelegationRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	delAddr, err := q.k.accountKeeper.AddressCodec().StringToBytes(req.DelegatorAddress)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid delegator address: %s", err)
	}

	if req.OperatorId == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid operator id")
	}

	rewards, err := q.k.GetOperatorDelegationRewards(ctx, delAddr, req.OperatorId)
	if err != nil {
		return nil, err
	}

	return &types.QueryOperatorDelegationRewardsResponse{Rewards: rewards}, nil
}

func (q queryServer) ServiceDelegationRewards(ctx context.Context, req *types.QueryServiceDelegationRewardsRequest) (*types.QueryServiceDelegationRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	delAddr, err := q.k.accountKeeper.AddressCodec().StringToBytes(req.DelegatorAddress)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid delegator address: %s", err)
	}

	if req.ServiceId == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid service id")
	}

	rewards, err := q.k.GetServiceDelegationRewards(ctx, delAddr, req.ServiceId)
	if err != nil {
		return nil, err
	}

	return &types.QueryServiceDelegationRewardsResponse{Rewards: rewards}, nil
}

func (q queryServer) DelegatorTotalRewards(ctx context.Context, req *types.QueryDelegatorTotalRewardsRequest) (*types.QueryDelegatorTotalRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	delAddr, err := q.k.accountKeeper.AddressCodec().StringToBytes(req.DelegatorAddress)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid delegator address: %s", err)
	}

	var total = types.DecPools{}
	var delRewards []types.DelegationDelegatorReward

	err = q.k.restakingKeeper.IterateUserPoolDelegations(ctx, req.DelegatorAddress, func(del restakingtypes.Delegation) (stop bool, err error) {
		delReward, err := q.k.GetPoolDelegationRewards(ctx, delAddr, del.TargetID)
		if err != nil {
			return false, err
		}
		delRewards = append(delRewards, types.NewDelegationDelegatorReward(del.Type, del.TargetID, delReward))
		total = total.Add(delReward...)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	err = q.k.restakingKeeper.IterateUserOperatorDelegations(ctx, req.DelegatorAddress, func(del restakingtypes.Delegation) (stop bool, err error) {
		delReward, err := q.k.GetOperatorDelegationRewards(ctx, delAddr, del.TargetID)
		if err != nil {
			return false, err
		}
		delRewards = append(delRewards, types.NewDelegationDelegatorReward(del.Type, del.TargetID, delReward))
		total = total.Add(delReward...)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	err = q.k.restakingKeeper.IterateUserServiceDelegations(ctx, req.DelegatorAddress, func(del restakingtypes.Delegation) (stop bool, err error) {
		delReward, err := q.k.GetServiceDelegationRewards(ctx, delAddr, del.TargetID)
		if err != nil {
			return false, err
		}
		delRewards = append(delRewards, types.NewDelegationDelegatorReward(del.Type, del.TargetID, delReward))
		total = total.Add(delReward...)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryDelegatorTotalRewardsResponse{Rewards: delRewards, Total: total}, nil
}

func (q queryServer) DelegatorWithdrawAddress(ctx context.Context, req *types.QueryDelegatorWithdrawAddressRequest) (*types.QueryDelegatorWithdrawAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	delAddr, err := q.k.accountKeeper.AddressCodec().StringToBytes(req.DelegatorAddress)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid delegator address: %s", err)
	}

	withdrawAddr, err := q.k.GetDelegatorWithdrawAddr(ctx, delAddr)
	if err != nil {
		return nil, err
	}

	withdrawAddrStr, err := q.k.accountKeeper.AddressCodec().BytesToString(withdrawAddr)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid withdraw address: %s", err)
	}

	return &types.QueryDelegatorWithdrawAddressResponse{WithdrawAddress: withdrawAddrStr}, nil
}
