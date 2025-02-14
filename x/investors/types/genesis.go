package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/milkyway-labs/milkyway/v9/utils"
)

var (
	DefaultInvestorsRewardRatio = sdkmath.LegacyNewDecWithPrec(5, 1) // 50%
)

// NewGenesisState returns a new GenesisState instance
func NewGenesisState(
	investorsRewardRatio sdkmath.LegacyDec,
	vestingInvestorsAddrs []string,
	validatorsInvestorsShares []ValidatorInvestorsShares,
) *GenesisState {
	return &GenesisState{
		InvestorsRewardRatio:      investorsRewardRatio,
		VestingInvestorsAddresses: vestingInvestorsAddrs,
		ValidatorsInvestorsShares: validatorsInvestorsShares,
	}
}

// DefaultGenesis returns a default GenesisState
func DefaultGenesis() *GenesisState {
	return NewGenesisState(DefaultInvestorsRewardRatio, nil, nil)
}

// Validate validates the GenesisState and returns an error if it is invalid.
func (data *GenesisState) Validate() error {
	err := ValidateInvestorsRewardRatio(data.InvestorsRewardRatio)
	if err != nil {
		return err
	}

	for _, investor := range data.VestingInvestorsAddresses {
		_, err = sdk.AccAddressFromBech32(investor)
		if err != nil {
			return fmt.Errorf("invalid investor address: %w", err)
		}
	}

	duplicatedInvestor := utils.FindDuplicate(data.VestingInvestorsAddresses)
	if duplicatedInvestor != nil {
		return fmt.Errorf("duplicated investor address: %s", *duplicatedInvestor)
	}

	for _, shares := range data.ValidatorsInvestorsShares {
		err := shares.Validate()
		if err != nil {
			return fmt.Errorf("invalid validator investors shares for %s: %w", shares.ValidatorAddress, err)
		}
	}

	duplicatedShares := utils.FindDuplicateFunc(data.ValidatorsInvestorsShares, func(a, b ValidatorInvestorsShares) bool {
		return a.ValidatorAddress == b.ValidatorAddress
	})
	if duplicatedShares != nil {
		return fmt.Errorf("duplicated validator address: %s", duplicatedShares.ValidatorAddress)
	}

	return nil
}

// ValidateInvestorsRewardRatio validates the investors reward ratio.
func ValidateInvestorsRewardRatio(ratio sdkmath.LegacyDec) error {
	if ratio.IsNegative() {
		return sdkerrors.ErrInvalidRequest.Wrapf(
			"investors reward ratio cannot be negative: %s",
			ratio,
		)
	} else if ratio.GT(sdkmath.LegacyOneDec()) {
		return sdkerrors.ErrInvalidRequest.Wrapf(
			"investors reward ratio cannot be greater than one: %s",
			ratio,
		)
	}
	return nil
}

// Validate validates the ValidatorInvestorsShares.
func (shares *ValidatorInvestorsShares) Validate() error {
	_, err := sdk.AccAddressFromBech32(shares.ValidatorAddress)
	if err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}

	if shares.InvestorsShares.IsNegative() {
		return sdkerrors.ErrInvalidRequest.Wrapf(
			"investors shares cannot be negative: %s",
			shares.InvestorsShares,
		)
	}

	return nil
}
