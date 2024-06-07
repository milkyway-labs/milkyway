package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewParams creates a new Params object
func NewParams(registrationFee sdk.Coins, deactivationTime time.Duration) Params {
	return Params{
		OperatorRegistrationFee: registrationFee,
		DeactivationTime:        deactivationTime,
	}
}

// DefaultParams returns default Params
func DefaultParams() Params {
	return Params{
		OperatorRegistrationFee: sdk.NewCoins(),
		DeactivationTime:        3 * 24 * time.Hour, // 3 days
	}
}

// Validate checks that the parameters have valid values.
func (p *Params) Validate() error {
	if !p.OperatorRegistrationFee.IsValid() {
		return fmt.Errorf("invalid operator registration fee: %s", p.OperatorRegistrationFee)
	}

	if p.DeactivationTime <= 0 {
		return ErrInvalidDeactivationTime
	}

	return nil
}
