package types

import (
	query "github.com/cosmos/cosmos-sdk/types/query"
)

// NewQueryUserInsuranceFundRequest creates a new QueryUserInsuranceFundRequest object.
func NewQueryUserInsuranceFundRequest(user string) *QueryUserInsuranceFundRequest {
	return &QueryUserInsuranceFundRequest{
		UserAddress: user,
	}
}

// NewQueryUserInsuranceFundsRequest creates a new QueryUsersInsuranceFundRequest object.
func NewQueryUserInsuranceFundsRequest(pagination *query.PageRequest) *QueryUserInsuranceFundsRequest {
	return &QueryUserInsuranceFundsRequest{
		Pagination: pagination,
	}
}

func NewUserInsuranceFundData(userAddress string, insuranceFund UserInsuranceFund) UserInsuranceFundData {
	return UserInsuranceFundData{
		UserAddress: userAddress,
		Balance:     insuranceFund.Balance,
		Used:        insuranceFund.Used,
	}
}

// NewQueryUserRestakableAssetsRequest creates a new QueryUserRestakableAssetsRequest object.
func NewQueryUserRestakableAssetsRequest(user string) *QueryUserRestakableAssetsRequest {
	return &QueryUserRestakableAssetsRequest{
		UserAddress: user,
	}
}

// NewQueryInsuranceFundRequest creates a new QueryInsuranceFundRequest object.
func NewQueryInsuranceFundRequest() *QueryInsuranceFundRequest {
	return &QueryInsuranceFundRequest{}
}

// NewQueryParamsRequest creates a new QueryParamsRequest object.
func NewQueryParamsRequest() *QueryParamsRequest {
	return &QueryParamsRequest{}
}
