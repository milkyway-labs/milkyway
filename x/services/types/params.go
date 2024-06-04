package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewParams creates a new Params object
func NewParams(serviceRegistrationFee sdk.Coins) Params {
	return Params{
		ServiceRegistrationFee: serviceRegistrationFee,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		ServiceRegistrationFee: sdk.NewCoins(),
	}
}

// Validate checks that the parameters have valid values.
func (p *Params) Validate() error {
	if !p.ServiceRegistrationFee.IsValid() {
		return errors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid service registration fee: %s", p.ServiceRegistrationFee)
	}

	return nil
}
