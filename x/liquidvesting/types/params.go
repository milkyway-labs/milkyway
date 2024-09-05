package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
)

// NewParams creates a new Params instance.
func NewParams(
	insurancePercentage cosmossdk_io_math.LegacyDec,
	burners []string,
	minters []string,
) Params {
	return Params{
		InsurancePercentage: insurancePercentage,
		Burners:             burners,
		Minters:             minters,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(cosmossdk_io_math.LegacyNewDec(2), []string{}, []string{})
}
