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

// NewQueryPoolUnbondingDelegationsRequest creates a new QueryPoolUnbondingDelegationsRequest instance
func NewQueryPoolUnbondingDelegationsRequest(poolID uint32, pagination *query.PageRequest) *QueryPoolUnbondingDelegationsRequest {
	return &QueryPoolUnbondingDelegationsRequest{PoolId: poolID, Pagination: pagination}
}

// NewQueryPoolUnbondingDelegationRequest creates a new QueryPoolUnbondingDelegationRequest instance
func NewQueryPoolUnbondingDelegationRequest(poolID uint32, delegatorAddress string) *QueryPoolUnbondingDelegationRequest {
	return &QueryPoolUnbondingDelegationRequest{PoolId: poolID, DelegatorAddress: delegatorAddress}
}

// NewQueryOperatorJoinedServicesRequest creates a new QueryOperatorJoinedServicesRequest instance
func NewQueryOperatorJoinedServicesRequest(operatorID uint32, pagination *query.PageRequest) *QueryOperatorJoinedServicesRequest {
	return &QueryOperatorJoinedServicesRequest{OperatorId: operatorID, Pagination: pagination}
}

// NewQueryOperatorDelegationsRequest creates a new QueryOperatorDelegationsRequest instance
func NewQueryOperatorDelegationsRequest(operatorID uint32, pagination *query.PageRequest) *QueryOperatorDelegationsRequest {
	return &QueryOperatorDelegationsRequest{OperatorId: operatorID, Pagination: pagination}
}

// NewQueryOperatorDelegationRequest creates a new QueryOperatorDelegationRequest instance
func NewQueryOperatorDelegationRequest(operatorID uint32, delegatorAddress string) *QueryOperatorDelegationRequest {
	return &QueryOperatorDelegationRequest{OperatorId: operatorID, DelegatorAddress: delegatorAddress}
}

// NewQueryOperatorUnbondingDelegationsRequest creates a new QueryOperatorUnbondingDelegationsRequest instance
func NewQueryOperatorUnbondingDelegationsRequest(operatorID uint32, pagination *query.PageRequest) *QueryOperatorUnbondingDelegationsRequest {
	return &QueryOperatorUnbondingDelegationsRequest{OperatorId: operatorID, Pagination: pagination}
}

// NewQueryOperatorUnbondingDelegationRequest creates a new QueryOperatorUnbondingDelegationRequest instance
func NewQueryOperatorUnbondingDelegationRequest(operatorID uint32, delegatorAddress string) *QueryOperatorUnbondingDelegationRequest {
	return &QueryOperatorUnbondingDelegationRequest{OperatorId: operatorID, DelegatorAddress: delegatorAddress}
}

// NewQueryServiceAllowedOperatorsRequest creates a new QueryServiceAllowedOperatorsRequest instance
func NewQueryServiceAllowedOperatorsRequest(serviceID uint32, pagination *query.PageRequest) *QueryServiceAllowedOperatorsRequest {
	return &QueryServiceAllowedOperatorsRequest{ServiceId: serviceID, Pagination: pagination}
}

// NewQueryServiceSecuringPoolsRequest creates a new QueryServiceSecuringPoolsRequest instance
func NewQueryServiceSecuringPoolsRequest(serviceID uint32, pagination *query.PageRequest) *QueryServiceSecuringPoolsRequest {
	return &QueryServiceSecuringPoolsRequest{ServiceId: serviceID, Pagination: pagination}
}

// NewQueryServiceOperatorsRequest creates a new QueryServiceOperatorsRequest instance
func NewQueryServiceOperatorsRequest(serviceID uint32, pagination *query.PageRequest) *QueryServiceOperatorsRequest {
	return &QueryServiceOperatorsRequest{ServiceId: serviceID, Pagination: pagination}
}

// NewQueryServiceDelegationsRequest creates a new QueryServiceDelegationsRequest instance
func NewQueryServiceDelegationsRequest(serviceID uint32, pagination *query.PageRequest) *QueryServiceDelegationsRequest {
	return &QueryServiceDelegationsRequest{ServiceId: serviceID, Pagination: pagination}
}

// NewQueryServiceDelegationRequest creates a new QueryServiceDelegationRequest instance
func NewQueryServiceDelegationRequest(serviceID uint32, delegatorAddress string) *QueryServiceDelegationRequest {
	return &QueryServiceDelegationRequest{ServiceId: serviceID, DelegatorAddress: delegatorAddress}
}

// NewQueryServiceUnbondingDelegationsRequest creates a new QueryServiceUnbondingDelegationsRequest instance
func NewQueryServiceUnbondingDelegationsRequest(serviceID uint32, pagination *query.PageRequest) *QueryServiceUnbondingDelegationsRequest {
	return &QueryServiceUnbondingDelegationsRequest{ServiceId: serviceID, Pagination: pagination}
}

// NewQueryServiceUnbondingDelegationRequest creates a new QueryServiceUnbondingDelegationRequest instance
func NewQueryServiceUnbondingDelegationRequest(serviceID uint32, delegatorAddress string) *QueryServiceUnbondingDelegationRequest {
	return &QueryServiceUnbondingDelegationRequest{ServiceId: serviceID, DelegatorAddress: delegatorAddress}
}

// NewQueryDelegatorPoolDelegationsRequest creates a new QueryDelegatorPoolDelegationsRequest instance
func NewQueryDelegatorPoolDelegationsRequest(delegatorAddress string, pagination *query.PageRequest) *QueryDelegatorPoolDelegationsRequest {
	return &QueryDelegatorPoolDelegationsRequest{DelegatorAddress: delegatorAddress, Pagination: pagination}
}

// NewQueryDelegatorPoolUnbondingDelegationsRequest creates a new QueryDelegatorPoolUnbondingDelegationsRequest instance
func NewQueryDelegatorPoolUnbondingDelegationsRequest(delegatorAddress string, pagination *query.PageRequest) *QueryDelegatorPoolUnbondingDelegationsRequest {
	return &QueryDelegatorPoolUnbondingDelegationsRequest{DelegatorAddress: delegatorAddress, Pagination: pagination}
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

// NewQueryDelegatorOperatorUnbondingDelegationsRequest creates a new QueryDelegatorOperatorUnbondingDelegationsRequest instance
func NewQueryDelegatorOperatorUnbondingDelegationsRequest(delegatorAddress string, pagination *query.PageRequest) *QueryDelegatorOperatorUnbondingDelegationsRequest {
	return &QueryDelegatorOperatorUnbondingDelegationsRequest{DelegatorAddress: delegatorAddress, Pagination: pagination}
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

// NewQueryDelegatorServiceUnbondingDelegationsRequest creates a new QueryDelegatorServiceUnbondingDelegationsRequest instance
func NewQueryDelegatorServiceUnbondingDelegationsRequest(delegatorAddress string, pagination *query.PageRequest) *QueryDelegatorServiceUnbondingDelegationsRequest {
	return &QueryDelegatorServiceUnbondingDelegationsRequest{DelegatorAddress: delegatorAddress, Pagination: pagination}
}

// NewQueryDelegatorServicesRequest creates a new QueryDelegatorServicesRequest instance
func NewQueryDelegatorServicesRequest(delegatorAddress string, pagination *query.PageRequest) *QueryDelegatorServicesRequest {
	return &QueryDelegatorServicesRequest{DelegatorAddress: delegatorAddress, Pagination: pagination}
}

// NewQueryDelegatorServiceRequest creates a new QueryDelegatorServiceRequest instance
func NewQueryDelegatorServiceRequest(delegatorAddress string, serviceID uint32) *QueryDelegatorServiceRequest {
	return &QueryDelegatorServiceRequest{DelegatorAddress: delegatorAddress, ServiceId: serviceID}
}

// NewQueryUserPreferencesRequest creates a new QueryUserPreferencesRequest instance
func NewQueryUserPreferencesRequest(userAddress string) *QueryUserPreferencesRequest {
	return &QueryUserPreferencesRequest{UserAddress: userAddress}
}

// NewQueryParamsRequest creates a new QueryParamsRequest instance
func NewQueryParamsRequest() *QueryParamsRequest {
	return &QueryParamsRequest{}
}
