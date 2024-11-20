package types

import (
	"fmt"
)

// NewGenesis creates a new GenesisState instance
func NewGenesis(nextPoolID uint32, pools []Pool, params Params) *GenesisState {
	return &GenesisState{
		Params:     params,
		NextPoolID: nextPoolID,
		Pools:      pools,
	}
}

// DefaultGenesis returns the default GenesisState
func DefaultGenesis() *GenesisState {
	return NewGenesis(1, nil, DefaultParams())
}

// Validate checks if the GenesisState is valid
func (data *GenesisState) Validate() error {
	err := data.Params.Validate()
	if err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	// Validate the next pool ID
	if data.NextPoolID == 0 {
		return fmt.Errorf("invalid next pool id")
	}

	// Validate the pools
	for _, pool := range data.Pools {
		err = pool.Validate()
		if err != nil {
			return fmt.Errorf("invalid pool with id %d: %w", pool.ID, err)
		}
	}

	return nil
}
