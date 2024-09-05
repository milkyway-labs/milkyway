package types

// NewQueryUserInsuranceFundRequest creates a new QueryUserInsuranceFundRequest object.
func NewQueryUserInsuranceFundRequest(user string) *QueryUserInsuranceFundRequest {
	return &QueryUserInsuranceFundRequest{
		UserAddress: user,
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
