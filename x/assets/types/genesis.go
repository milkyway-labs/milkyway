package types

import (
	"fmt"
)

// NewGenesisState returns a new GenesisState instance
func NewGenesisState(params Params, assets []Asset) *GenesisState {
	return &GenesisState{
		Params: params,
		Assets: assets,
	}
}

// DefaultGenesis returns a default GenesisState
func DefaultGenesis() *GenesisState {
	return NewGenesisState(DefaultParams(), nil)
}

// --------------------------------------------------------------------------------------------------------------------

// Validate validates the GenesisState and returns an error if it is invalid.
func (data *GenesisState) Validate() error {
	// Validate params
	err := data.Params.Validate()
	if err != nil {
		return fmt.Errorf("invalid params: %s", err)
	}

	// Validate the assets
	for _, asset := range data.Assets {
		err = asset.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
