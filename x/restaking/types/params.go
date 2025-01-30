package types

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultUnbondingTime = 3 * 24 * time.Hour
	DefaultMaxEntries    = 7
)

var (
	DefaultRestakingCap = math.LegacyZeroDec() // no cap
)

// NewParams returns a new Params instance
func NewParams(
	unbondingTime time.Duration,
	allowedDenoms []string,
	restakingCap math.LegacyDec,
	maxEntries uint32,
) Params {
	return Params{
		UnbondingTime: unbondingTime,
		AllowedDenoms: allowedDenoms,
		RestakingCap:  restakingCap,
		MaxEntries:    maxEntries,
	}
}

// DefaultParams return a Params instance with default values set
func DefaultParams() Params {
	return NewParams(DefaultUnbondingTime, nil, DefaultRestakingCap, DefaultMaxEntries)
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

	if p.MaxEntries == 0 {
		return fmt.Errorf("max entries must be positive: %d", p.MaxEntries)
	}

	return nil
}
