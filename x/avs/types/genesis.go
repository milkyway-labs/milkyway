package types

import (
	"fmt"

	"github.com/milkyway-labs/milkyway/utils"
)

// NewGenesisState returns a new GenesisState instance
func NewGenesisState(services []AVS, params Params) *GenesisState {
	return &GenesisState{
		Services: services,
		Params:   params,
	}
}

// DefaultGenesisState returns a default GenesisState
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(nil, DefaultParams())
}

// --------------------------------------------------------------------------------------------------------------------

// ValidateGenesis validates the given genesis state and returns an error if something is invalid
func ValidateGenesis(data *GenesisState) error {
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
func findDuplicatedService(services []AVS) *AVS {
	return utils.FindDuplicate(services, func(a, b AVS) bool {
		return a.ID == b.ID
	})
}
