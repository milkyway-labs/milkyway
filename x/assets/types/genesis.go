package types

// NewGenesisState returns a new GenesisState instance
func NewGenesisState(assets []Asset) *GenesisState {
	return &GenesisState{
		Assets: assets,
	}
}

// DefaultGenesis returns a default GenesisState
func DefaultGenesis() *GenesisState {
	return NewGenesisState(nil)
}

// --------------------------------------------------------------------------------------------------------------------

// Validate validates the GenesisState and returns an error if it is invalid.
func (data *GenesisState) Validate() error {
	// Validate the assets
	for _, asset := range data.Assets {
		err := asset.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
