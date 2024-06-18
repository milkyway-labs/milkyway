package types

import (
	"fmt"
)

// NewGenesis creates a new genesis state
func NewGenesis(
	poolsDelegations []PoolDelegation,
	servicesDelegations []ServiceDelegation,
	operatorsDelegations []OperatorDelegation,
	params Params,
) *GenesisState {
	return &GenesisState{
		PoolsDelegations:     poolsDelegations,
		ServicesDelegations:  servicesDelegations,
		OperatorsDelegations: operatorsDelegations,
		Params:               params,
	}
}

// DefaultGenesis returns a default genesis state
func DefaultGenesis() *GenesisState {
	return NewGenesis(nil, nil, nil, DefaultParams())
}

// Validate performs basic validation of genesis data
func (g *GenesisState) Validate() error {
	// Validate pools delegations
	for _, entry := range g.PoolsDelegations {
		err := entry.Validate()
		if err != nil {
			return fmt.Errorf("invalid pool delegation: %w", err)
		}
	}

	// Validate services delegations
	for _, entry := range g.ServicesDelegations {
		err := entry.Validate()
		if err != nil {
			return fmt.Errorf("invalid service delegation: %w", err)
		}
	}

	// Validate operators delegations
	for _, entry := range g.OperatorsDelegations {
		err := entry.Validate()
		if err != nil {
			return fmt.Errorf("invalid operator delegation: %w", err)
		}
	}

	// Validate the params
	err := g.Params.Validate()
	if err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	return nil
}
