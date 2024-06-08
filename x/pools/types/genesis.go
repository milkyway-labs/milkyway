package types

import (
	"fmt"
)

// NewGenesis creates a new GenesisState instance
func NewGenesis(nextPoolID uint32, pools []Pool) GenesisState {
	return GenesisState{
		NextPoolID: nextPoolID,
		Pools:      pools,
	}
}

// Validate checks if the GenesisState is valid
func (data *GenesisState) Validate() error {
	// Validate the next pool ID
	if data.NextPoolID == 0 {
		return fmt.Errorf("invalid next pool id")
	}

	// Validate the pools
	for _, pool := range data.Pools {
		err := pool.Validate()
		if err != nil {
			return fmt.Errorf("invalid pool with id %d: %w", pool.ID, err)
		}
	}

	return nil
}
