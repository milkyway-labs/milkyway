package types

import (
	"github.com/cosmos/cosmos-sdk/types/address"
)

// RewardsPoolAddress is the address of global rewards pool where rewards are
// moved from each rewards plan's rewards pool and distributed to delegators
// later.
var RewardsPoolAddress = address.Module(ModuleName, []byte("RewardsPool"))

// NewParams creates a new Params object
func NewParams() Params {
	return Params{}
}

// DefaultParams returns default Params
func DefaultParams() Params {
	return Params{}
}

// Validate checks that the parameters have valid values.
func (p *Params) Validate() error {
	return nil
}
