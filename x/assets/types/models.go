package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewAsset(denom, ticker string, exponent uint32) Asset {
	return Asset{
		Denom:    denom,
		Ticker:   ticker,
		Exponent: exponent,
	}
}

func (asset Asset) Validate() error {
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
