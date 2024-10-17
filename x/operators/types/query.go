package types

import (
	"github.com/cosmos/cosmos-sdk/types/query"
)

// NewQueryOperatorRequest creates a new QueryOperatorRequest object
func NewQueryOperatorRequest(operatorID uint32) *QueryOperatorRequest {
	return &QueryOperatorRequest{OperatorId: operatorID}
}

// NewQueryOperatorParamsRequest creates a new QueryOperatorParamsRequest object
func NewQueryOperatorParamsRequest(operatorID uint32) *QueryOperatorParamsRequest {
	return &QueryOperatorParamsRequest{OperatorId: operatorID}
}

// NewQueryOperatorsRequest creates a new QueryOperatorsRequest object
func NewQueryOperatorsRequest(pagination *query.PageRequest) *QueryOperatorsRequest {
	return &QueryOperatorsRequest{Pagination: pagination}
}

// NewQueryParamsRequest creates a new QueryParamsRequest object
func NewQueryParamsRequest() *QueryParamsRequest {
	return &QueryParamsRequest{}
}
