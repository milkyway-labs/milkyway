package types

import (
	"fmt"

	"github.com/milkyway-labs/milkyway/utils"
)

// NewGenesisState returns a new GenesisState instance
func NewGenesisState(nextAVSID uint32, services []Service, params Params) *GenesisState {
	return &GenesisState{
		NextAVSID: nextAVSID,
		Services:  services,
		Params:    params,
	}
}

// DefaultGenesis returns a default GenesisState
func DefaultGenesis() *GenesisState {
	return NewGenesisState(1, nil, DefaultParams())
}

// --------------------------------------------------------------------------------------------------------------------

// Validate validates the GenesisState and returns an error if it is invalid.
func (data *GenesisState) Validate() error {
	// Check for the next Service ID
	if data.NextAVSID == 0 {
		return fmt.Errorf("invalid next Service ID: %d", data.NextAVSID)
	}

	// Check for duplicated services
	if duplicate := findDuplicatedService(data.Services); duplicate != nil {
		return fmt.Errorf("duplicated service: %d", duplicate.ID)
	}

	// Validate services
	for _, service := range data.Services {
		if err := service.Validate(); err != nil {
			return fmt.Errorf("invalid service with id %d: %s", service.ID, err)
		}
	}

	// Validate params
	if err := data.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %s", err)
	}

	return nil
}

// findDuplicatedService returns the first duplicated service in the slice.
// If no duplicates are found, it returns nil instead.
func findDuplicatedService(services []Service) *Service {
	return utils.FindDuplicate(services, func(a, b Service) bool {
		return a.ID == b.ID
	})
}
