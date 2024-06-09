package types

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ParsePoolID parses a pool id from a string
func ParsePoolID(value string) (uint32, error) {
	parsed, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid pool id: %s", value)
	}
	return uint32(parsed), nil
}

// --------------------------------------------------------------------------------------------------------------------

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
