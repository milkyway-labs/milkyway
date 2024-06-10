package types

import (
	"fmt"
	"time"
)

// NewParams returns a new Params instance
func NewParams(unbondingTime time.Duration) Params {
	return Params{
		UnbondingTime: unbondingTime,
	}
}

// DefaultParams return a Params instance with default values set
func DefaultParams() Params {
	return NewParams(3 * 24 * time.Hour)
}

// Validate performs basic validation of params
func (p *Params) Validate() error {
	if p.UnbondingTime == 0 {
		return fmt.Errorf("invalid unbonding time: %s", p.UnbondingTime)
	}

	return nil
}
