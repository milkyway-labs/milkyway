package types

import (
	"github.com/cosmos/cosmos-sdk/types/query"
)

// NewQueryPoolDelegationsRequest creates a new QueryPoolDelegationsRequest instance
func NewQueryPoolDelegationsRequest(poolID uint32, pagination *query.PageRequest) *QueryPoolDelegationsRequest {
	return &QueryPoolDelegationsRequest{PoolId: poolID, Pagination: pagination}
}

// NewQueryPoolDelegationRequest creates a new QueryPoolDelegationRequest instance
func NewQueryPoolDelegationRequest(poolID uint32, delegatorAddress string) *QueryPoolDelegationRequest {
	return &QueryPoolDelegationRequest{PoolId: poolID, DelegatorAddress: delegatorAddress}
}

// NewQueryOperatorDelegationsRequest creates a new QueryOperatorDelegationsRequest instance
func NewQueryOperatorDelegationsRequest(operatorID uint32, pagination *query.PageRequest) *QueryOperatorDelegationsRequest {
	return &QueryOperatorDelegationsRequest{OperatorId: operatorID, Pagination: pagination}
}

// NewQueryOperatorDelegationRequest creates a new QueryOperatorDelegationRequest instance
func NewQueryOperatorDelegationRequest(operatorID uint32, delegatorAddress string) *QueryOperatorDelegationRequest {
	return &QueryOperatorDelegationRequest{OperatorId: operatorID, DelegatorAddress: delegatorAddress}
}

// NewQueryServiceDelegationsRequest creates a new QueryServiceDelegationsRequest instance
func NewQueryServiceDelegationsRequest(serviceID uint32, pagination *query.PageRequest) *QueryServiceDelegationsRequest {
	return &QueryServiceDelegationsRequest{ServiceId: serviceID, Pagination: pagination}
}

// NewQueryServiceDelegationRequest creates a new QueryServiceDelegationRequest instance
func NewQueryServiceDelegationRequest(serviceID uint32, delegatorAddress string) *QueryServiceDelegationRequest {
	return &QueryServiceDelegationRequest{ServiceId: serviceID, DelegatorAddress: delegatorAddress}
}

// NewQueryDelegatorPoolDelegationsRequest creates a new QueryDelegatorPoolDelegationsRequest instance
func NewQueryDelegatorPoolDelegationsRequest(delegatorAddress string, pagination *query.PageRequest) *QueryDelegatorPoolDelegationsRequest {
	return &QueryDelegatorPoolDelegationsRequest{DelegatorAddress: delegatorAddress, Pagination: pagination}
}

// NewQueryDelegatorOperatorDelegationsRequest creates a new QueryDelegatorOperatorDelegationsRequest instance
func NewQueryDelegatorOperatorDelegationsRequest(delegatorAddress string, pagination *query.PageRequest) *QueryDelegatorOperatorDelegationsRequest {
	return &QueryDelegatorOperatorDelegationsRequest{DelegatorAddress: delegatorAddress, Pagination: pagination}
}

// NewQueryDelegatorServiceDelegationsRequest creates a new QueryDelegatorServiceDelegationsRequest instance
func NewQueryDelegatorServiceDelegationsRequest(delegatorAddress string, pagination *query.PageRequest) *QueryDelegatorServiceDelegationsRequest {
	return &QueryDelegatorServiceDelegationsRequest{DelegatorAddress: delegatorAddress, Pagination: pagination}
}
