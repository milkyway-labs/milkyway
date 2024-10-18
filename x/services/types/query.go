package types

import (
	"github.com/cosmos/cosmos-sdk/types/query"
)

// NewQueryParamsRequest creates a new QueryParamsRequest instance
func NewQueryParamsRequest() *QueryParamsRequest {
	return &QueryParamsRequest{}
}

// NewQueryServicesRequest creates a new QueryServiceRequest instance
func NewQueryServicesRequest(pagination *query.PageRequest) *QueryServicesRequest {
	return &QueryServicesRequest{
		Pagination: pagination,
	}
}

// NewQueryServiceRequest creates a new QueryServiceRequest instance
func NewQueryServiceRequest(serviceID uint32) *QueryServiceRequest {
	return &QueryServiceRequest{
		ServiceId: serviceID,
	}
}

// NewQueryServiceParamsRequest creates a new QueryServiceParamsRequest instance
func NewQueryServiceParamsRequest(serviceID uint32) *QueryServiceParamsRequest {
	return &QueryServiceParamsRequest{
		ServiceId: serviceID,
	}
}
