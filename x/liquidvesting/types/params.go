package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewParams creates a new Params instance.
func NewParams(
	insurancePercentage math.LegacyDec,
	burners []string,
	minters []string,
	trustedDelegates []string,
) Params {
	return Params{
		InsurancePercentage: insurancePercentage,
		Burners:             burners,
		Minters:             minters,
		TrustedDelegates:    trustedDelegates,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(math.LegacyNewDec(2), nil, nil, nil)
}

// Validate ensure that the Prams structure is correct
func (p *Params) Validate() error {
	if p.InsurancePercentage.LTE(math.LegacyNewDec(0)) || p.InsurancePercentage.GT(math.LegacyNewDec(100)) {
		return ErrInvalidInsurancePercentage
	}
	for _, address := range p.Minters {
		_, err := sdk.AccAddressFromBech32(address)
		if err != nil {
			return err
		}
	}
	for _, address := range p.Burners {
		_, err := sdk.AccAddressFromBech32(address)
		if err != nil {
			return err
		}
	}
	for _, address := range p.TrustedDelegates {
		_, err := sdk.AccAddressFromBech32(address)
		if err != nil {
			return err
		}
	}
	return nil
}
