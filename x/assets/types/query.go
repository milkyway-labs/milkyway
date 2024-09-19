package types

import (
	"github.com/cosmos/cosmos-sdk/types/query"
)

// NewQueryAssetsRequest creates a new instance of QueryAssetsRequest
func NewQueryAssetsRequest(ticker string, pagination *query.PageRequest) *QueryAssetsRequest {
	return &QueryAssetsRequest{Ticker: ticker, Pagination: pagination}
}

// NewQueryAssetRequest creates a new instance of QueryAssetRequest
func NewQueryAssetRequest(denom string) *QueryAssetRequest {
	return &QueryAssetRequest{Denom: denom}
}
