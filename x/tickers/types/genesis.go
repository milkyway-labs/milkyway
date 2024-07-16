package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Tickers = map[string]string

// NewGenesisState returns a new GenesisState instance
func NewGenesisState(params Params, tickers Tickers) *GenesisState {
	return &GenesisState{
		Params:  params,
		Tickers: tickers,
	}
}

// DefaultGenesis returns a default GenesisState
func DefaultGenesis() *GenesisState {
	return NewGenesisState(DefaultParams(), Tickers{})
}

// --------------------------------------------------------------------------------------------------------------------

// Validate validates the GenesisState and returns an error if it is invalid.
func (data *GenesisState) Validate() error {
	// Validate params
	err := data.Params.Validate()
	if err != nil {
		return fmt.Errorf("invalid params: %s", err)
	}

	for denom, ticker := range data.Tickers {
		if err := sdk.ValidateDenom(denom); err != nil {
			return err
		}
		if err := ValidateTicker(ticker); err != nil {
			return err
		}
	}
	return nil
}
