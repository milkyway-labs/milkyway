package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewAsset returns a new Asset instance
func NewAsset(denom, ticker string, exponent uint32) Asset {
	return Asset{
		Denom:    denom,
		Ticker:   ticker,
		Exponent: exponent,
	}
}

// Validate validates the Asset instance
func (asset *Asset) Validate() error {
	err := sdk.ValidateDenom(asset.Denom)
	if err != nil {
		return fmt.Errorf("invalid denom: %w", err)
	}

	err = ValidateTicker(asset.Ticker)
	if err != nil {
		return fmt.Errorf("invalid ticker: %w", err)
	}

	return nil
}
