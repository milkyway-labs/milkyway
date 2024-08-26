package types

import (
	"github.com/cosmos/cosmos-sdk/types/query"
)

// NewQueryRewardsPlansRequest creates a new instance of QueryRewardsPlansRequest.
func NewQueryRewardsPlansRequest(pagination *query.PageRequest) *QueryRewardsPlansRequest {
	return &QueryRewardsPlansRequest{
		Pagination: pagination,
	}
}

// NewQueryRewardsPlanRequest creates a new instance of QueryRewardsPlanRequest.
func NewQueryRewardsPlanRequest(planID uint64) *QueryRewardsPlanRequest {
	return &QueryRewardsPlanRequest{
		PlanId: planID,
	}
}

// NewQueryPoolOutstandingRewardsRequest creates a new instance of QueryPoolOutstandingRewardsRequest.
func NewQueryPoolOutstandingRewardsRequest(poolID uint32) *QueryPoolOutstandingRewardsRequest {
	return &QueryPoolOutstandingRewardsRequest{
		PoolId: poolID,
	}
}

// NewQueryOperatorOutstandingRewardsRequest creates a new instance of QueryOperatorOutstandingRewardsRequest.
func NewQueryOperatorOutstandingRewardsRequest(operatorID uint32) *QueryOperatorOutstandingRewardsRequest {
	return &QueryOperatorOutstandingRewardsRequest{
		OperatorId: operatorID,
	}
}

// NewQueryServiceOutstandingRewardsRequest creates a new instance of QueryServiceOutstandingRewardsRequest.
func NewQueryServiceOutstandingRewardsRequest(serviceID uint32) *QueryServiceOutstandingRewardsRequest {
	return &QueryServiceOutstandingRewardsRequest{
		ServiceId: serviceID,
	}
}

// NewQueryOperatorCommissionRequest creates a new instance of QueryOperatorCommissionRequest.
func NewQueryOperatorCommissionRequest(operatorID uint32) *QueryOperatorCommissionRequest {
	return &QueryOperatorCommissionRequest{
		OperatorId: operatorID,
	}
}

// NewQueryPoolDelegationRewardsRequest creates a new instance of QueryPoolDelegationRewardsRequest.
func NewQueryPoolDelegationRewardsRequest(poolID uint32, delegator string) *QueryPoolDelegationRewardsRequest {
	return &QueryPoolDelegationRewardsRequest{
		PoolId:           poolID,
		DelegatorAddress: delegator,
	}
}

// NewQueryOperatorDelegationRewardsRequest creates a new instance of QueryOperatorDelegationRewardsRequest.
func NewQueryOperatorDelegationRewardsRequest(operatorID uint32, delegator string) *QueryOperatorDelegationRewardsRequest {
	return &QueryOperatorDelegationRewardsRequest{
		OperatorId:       operatorID,
		DelegatorAddress: delegator,
	}
}

// NewQueryServiceDelegationRewardsRequest creates a new instance of QueryServiceDelegationRewardsRequest.
func NewQueryServiceDelegationRewardsRequest(serviceID uint32, delegator string) *QueryServiceDelegationRewardsRequest {
	return &QueryServiceDelegationRewardsRequest{
		ServiceId:        serviceID,
		DelegatorAddress: delegator,
	}
}

func NewQueryDelegatorTotalRewardsRequest(delegator string) *QueryDelegatorTotalRewardsRequest {
	return &QueryDelegatorTotalRewardsRequest{
		DelegatorAddress: delegator,
	}
}

func NewQueryDelegatorWithdrawAddressRequest(delegator string) *QueryDelegatorWithdrawAddressRequest {
	return &QueryDelegatorWithdrawAddressRequest{
		DelegatorAddress: delegator,
	}
}
