package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewParams creates a new Params object
func NewParams(rewardsPlanCreationFee sdk.Coins) Params {
	return Params{
		RewardsPlanCreationFee: rewardsPlanCreationFee,
	}
}

// DefaultParams returns default Params
func DefaultParams() Params {
	return NewParams(nil)
}

// Validate checks that the parameters have valid values.
func (p *Params) Validate() error {
	err := p.RewardsPlanCreationFee.Validate()
	if err != nil {
		return fmt.Errorf("invalid rewards plan creation fee: %w", err)
	}
	return nil
}
