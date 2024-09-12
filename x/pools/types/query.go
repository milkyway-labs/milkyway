package types

import (
	"github.com/cosmos/cosmos-sdk/types/query"
)

// NewQueryPoolByIDRequest creates a new instance of QueryPoolByIdRequest
func NewQueryPoolByIDRequest(poolID uint32) *QueryPoolByIdRequest {
	return &QueryPoolByIdRequest{
		PoolId: poolID,
	}
}

// NewQueryPoolByDenomRequest creates a new instance of QueryPoolByDenomRequest
func NewQueryPoolByDenomRequest(denom string) *QueryPoolByDenomRequest {
	return &QueryPoolByDenomRequest{
		Denom: denom,
	}
}

// NewQueryPoolsRequest creates a new instance of QueryPoolsRequest
func NewQueryPoolsRequest(pagination *query.PageRequest) *QueryPoolsRequest {
	return &QueryPoolsRequest{
		Pagination: pagination,
	}
}
