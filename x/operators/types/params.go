package types

import (
	"time"
)

// NewParams creates a new Params object
func NewParams(deactivationTime time.Duration) Params {
	return Params{
		DeactivationTime: deactivationTime,
	}
}

// DefaultParams returns default Params
func DefaultParams() Params {
	return Params{
		DeactivationTime: 3 * 24 * time.Hour, // 3 days
	}
}

// Validate checks that the parameters have valid values.
func (p *Params) Validate() error {
	if p.DeactivationTime <= 0 {
		return ErrInvalidDeactivationTime
	}

	return nil
}
