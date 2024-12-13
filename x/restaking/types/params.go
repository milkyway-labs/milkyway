package types

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultUnbondingTime = 3 * 24 * time.Hour
)

var (
	DefaultRestakingCap = math.LegacyZeroDec() // no cap
)

// NewParams returns a new Params instance
func NewParams(unbondingTime time.Duration, allowedDenoms []string, restakingCap math.LegacyDec) Params {
	return Params{
		UnbondingTime: unbondingTime,
		AllowedDenoms: allowedDenoms,
		RestakingCap:  restakingCap,
	}
}

// DefaultParams return a Params instance with default values set
func DefaultParams() Params {
	return NewParams(DefaultUnbondingTime, nil, DefaultRestakingCap)
}

// Validate performs basic validation of params
func (p *Params) Validate() error {
	if p.UnbondingTime == 0 {
		return fmt.Errorf("invalid unbonding time: %s", p.UnbondingTime)
	}

	for _, denom := range p.AllowedDenoms {
		err := sdk.ValidateDenom(denom)
		if err != nil {
			return fmt.Errorf("invalid allowed denom: %s, %w", denom, err)
		}
	}

	if p.RestakingCap.IsNegative() {
		return fmt.Errorf("restaking cap cannot be negative: %s", p.RestakingCap)
	}

	return nil
}
