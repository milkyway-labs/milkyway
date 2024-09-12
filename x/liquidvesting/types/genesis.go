package types

// NewGenesisState creates a new GenesisState instance.
func NewGenesisState(params Params) *GenesisState {
	return &GenesisState{
		Params: params,
	}
}

// DefaultGenesisState returns a default GenesisState.
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams())
}

func (g *GenesisState) Validate() error {
	return g.Params.Validate()
}
