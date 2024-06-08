package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewPool creates a new Pool instance
func NewPool(id uint32, denom string) Pool {
	return Pool{
		ID:    id,
		Denom: denom,
	}
}

// Validate checks if the pool is valid
func (p *Pool) Validate() error {
	if p.ID == 0 {
		return fmt.Errorf("invalid pool id")
	}

	if sdk.ValidateDenom(p.Denom) != nil {
		return fmt.Errorf("invalid pool denom")
	}

	return nil
}
