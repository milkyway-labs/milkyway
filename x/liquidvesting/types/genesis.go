package types

// NewGenesisState creates a new GenesisState instance.
func NewGenesisState(
	params Params,
	burnCoins []BurnCoins,
	userInsuUserInsuranceFundItems []UserInsuranceFundState,
) *GenesisState {
	return &GenesisState{
		Params:             params,
		BurnCoins:          burnCoins,
		UserInsuranceFunds: userInsuUserInsuranceFundItems,
	}
}

// DefaultGenesisState returns a default GenesisState.
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams(), nil, nil)
}

func (g *GenesisState) Validate() error {
	if err := g.Params.Validate(); err != nil {
		return err
	}

	for _, burnCoin := range g.BurnCoins {
		if err := burnCoin.Validate(); err != nil {
			return err
		}
	}

	for _, userInsuranceFund := range g.UserInsuranceFunds {
		if err := userInsuranceFund.Validate(); err != nil {
			return err
		}
	}

	return nil
}
