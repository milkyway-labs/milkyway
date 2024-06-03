package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewParams creates a new Params object
func NewParams(avsRegistrationFee sdk.Coins) Params {
	return Params{
		AvsRegistrationFee: avsRegistrationFee,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		AvsRegistrationFee: sdk.NewCoins(),
	}
}

// Validate checks that the parameters have valid values.
func (p *Params) Validate() error {
	if !p.AvsRegistrationFee.IsValid() {
		return errors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid AVS registration fee: %s", p.AvsRegistrationFee)
	}

	return nil
}
