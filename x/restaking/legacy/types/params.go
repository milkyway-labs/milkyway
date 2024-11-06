package types

import (
	"fmt"
	"time"
)

func NewLegacyParams(unbondingTime time.Duration) Params {
	return Params{UnbondingTime: unbondingTime}
}

// Validate performs basic validation of params
func (p *Params) Validate() error {
	if p.UnbondingTime == 0 {
		return fmt.Errorf("invalid unbonding time: %s", p.UnbondingTime)
	}

	return nil
}
