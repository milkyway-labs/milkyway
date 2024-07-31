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

// NewQueryOperatorParamsRequest creates a new QueryOperatorParamsRequest instance
func NewQueryOperatorParamsRequest(operatorID uint32) *QueryOperatorParamsRequest {
	return &QueryOperatorParamsRequest{OperatorId: operatorID}
}

// NewQueryOperatorDelegationsRequest creates a new QueryOperatorDelegationsRequest instance
func NewQueryOperatorDelegationsRequest(operatorID uint32, pagination *query.PageRequest) *QueryOperatorDelegationsRequest {
	return &QueryOperatorDelegationsRequest{OperatorId: operatorID, Pagination: pagination}
}

// NewQueryOperatorDelegationRequest creates a new QueryOperatorDelegationRequest instance
func NewQueryOperatorDelegationRequest(operatorID uint32, delegatorAddress string) *QueryOperatorDelegationRequest {
	return &QueryOperatorDelegationRequest{OperatorId: operatorID, DelegatorAddress: delegatorAddress}
}

// NewQueryServiceParamsRequest creates a new QueryServiceParamsRequest instance
func NewQueryServiceParamsRequest(serviceID uint32) *QueryServiceParamsRequest {
	return &QueryServiceParamsRequest{ServiceId: serviceID}
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

// NewQueryDelegatorPoolsRequest creates a new QueryDelegatorPoolsRequest instance
func NewQueryDelegatorPoolsRequest(delegatorAddress string, pagination *query.PageRequest) *QueryDelegatorPoolsRequest {
	return &QueryDelegatorPoolsRequest{DelegatorAddress: delegatorAddress, Pagination: pagination}
}

// NewQueryDelegatorPoolRequest creates a new QueryDelegatorPoolRequest instance
func NewQueryDelegatorPoolRequest(delegatorAddress string, poolID uint32) *QueryDelegatorPoolRequest {
	return &QueryDelegatorPoolRequest{DelegatorAddress: delegatorAddress, PoolId: poolID}
}

// NewQueryDelegatorOperatorDelegationsRequest creates a new QueryDelegatorOperatorDelegationsRequest instance
func NewQueryDelegatorOperatorDelegationsRequest(delegatorAddress string, pagination *query.PageRequest) *QueryDelegatorOperatorDelegationsRequest {
	return &QueryDelegatorOperatorDelegationsRequest{DelegatorAddress: delegatorAddress, Pagination: pagination}
}

// NewQueryDelegatorOperatorsRequest creates a new QueryDelegatorOperatorsRequest instance
func NewQueryDelegatorOperatorsRequest(delegatorAddress string, pagination *query.PageRequest) *QueryDelegatorOperatorsRequest {
	return &QueryDelegatorOperatorsRequest{DelegatorAddress: delegatorAddress, Pagination: pagination}
}

// NewQueryDelegatorOperatorRequest creates a new QueryDelegatorOperatorRequest instance
func NewQueryDelegatorOperatorRequest(delegatorAddress string, operatorID uint32) *QueryDelegatorOperatorRequest {
	return &QueryDelegatorOperatorRequest{DelegatorAddress: delegatorAddress, OperatorId: operatorID}
}

// NewQueryDelegatorServiceDelegationsRequest creates a new QueryDelegatorServiceDelegationsRequest instance
func NewQueryDelegatorServiceDelegationsRequest(delegatorAddress string, pagination *query.PageRequest) *QueryDelegatorServiceDelegationsRequest {
	return &QueryDelegatorServiceDelegationsRequest{DelegatorAddress: delegatorAddress, Pagination: pagination}
}

// NewQueryDelegatorServicesRequest creates a new QueryDelegatorServicesRequest instance
func NewQueryDelegatorServicesRequest(delegatorAddress string, pagination *query.PageRequest) *QueryDelegatorServicesRequest {
	return &QueryDelegatorServicesRequest{DelegatorAddress: delegatorAddress, Pagination: pagination}
}

// NewQueryDelegatorServiceRequest creates a new QueryDelegatorServiceRequest instance
func NewQueryDelegatorServiceRequest(delegatorAddress string, serviceID uint32) *QueryDelegatorServiceRequest {
	return &QueryDelegatorServiceRequest{DelegatorAddress: delegatorAddress, ServiceId: serviceID}
}

// NewQueryParamsRequest creates a new QueryParamsRequest instance
func NewQueryParamsRequest() *QueryParamsRequest {
	return &QueryParamsRequest{}
}
