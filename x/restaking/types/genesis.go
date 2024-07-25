package types

import (
	"fmt"
)

// NewGenesis creates a new genesis state
func NewGenesis(delegations []Delegation, params Params) *GenesisState {
	return &GenesisState{
		Delegations: delegations,
		Params:      params,
	}
}

// DefaultGenesis returns a default genesis state
func DefaultGenesis() *GenesisState {
	return NewGenesis(nil, DefaultParams())
}

// Validate performs basic validation of genesis data
func (g *GenesisState) Validate() error {
	// Validate delegations
	for _, entry := range g.Delegations {
		err := entry.Validate()
		if err != nil {
			return fmt.Errorf("invalid delegation: %w", err)
		}
	}

	// Validate the params
	err := g.Params.Validate()
	if err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	return nil
}
