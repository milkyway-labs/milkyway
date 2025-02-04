package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/milkyway-labs/milkyway/v7/utils"
)

var (
	DefaultVestingAccountsRewardRatio = sdkmath.LegacyNewDecWithPrec(5, 1) // 50%
)

// NewGenesisState returns a new GenesisState instance
func NewGenesisState(
	vestingAccountsRewardRatio sdkmath.LegacyDec,
	validatorsVestingAccountsShares []ValidatorVestingAccountsShares,
) *GenesisState {
	return &GenesisState{
		VestingAccountsRewardRatio:      vestingAccountsRewardRatio,
		ValidatorsVestingAccountsShares: validatorsVestingAccountsShares,
	}
}

// DefaultGenesis returns a default GenesisState
func DefaultGenesis() *GenesisState {
	return NewGenesisState(DefaultVestingAccountsRewardRatio, nil)
}

// Validate validates the GenesisState and returns an error if it is invalid.
func (data *GenesisState) Validate() error {
	err := ValidateVestingAccountsRewardRatio(data.VestingAccountsRewardRatio)
	if err != nil {
		return err
	}

	duplicate := utils.FindDuplicateFunc(data.ValidatorsVestingAccountsShares, func(a, b ValidatorVestingAccountsShares) bool {
		return a.ValidatorAddress == b.ValidatorAddress
	})
	if duplicate != nil {
		return fmt.Errorf("duplicated validator address: %s", duplicate.ValidatorAddress)
	}

	for _, shares := range data.ValidatorsVestingAccountsShares {
		err := shares.Validate()
		if err != nil {
			return fmt.Errorf("invalid validator vesting accounts shares for %s: %w", shares.ValidatorAddress, err)
		}
	}

	return nil
}

// ValidateVestingAccountsRewardRatio validates the vesting accounts reward ratio.
func ValidateVestingAccountsRewardRatio(ratio sdkmath.LegacyDec) error {
	if ratio.IsNegative() {
		return sdkerrors.ErrInvalidRequest.Wrapf(
			"vesting accounts reward ratio cannot be negative: %s",
			ratio,
		)
	} else if ratio.GT(sdkmath.LegacyOneDec()) {
		return sdkerrors.ErrInvalidRequest.Wrapf(
			"vesting accounts reward ratio cannot be greater than one: %s",
			ratio,
		)
	}
	return nil
}

// Validate validates the ValidatorVestingAccountsShares.
func (shares *ValidatorVestingAccountsShares) Validate() error {
	_, err := sdk.AccAddressFromBech32(shares.ValidatorAddress)
	if err != nil {
		return fmt.Errorf("invalid validator address: %w", err)
	}

	if shares.VestingAccountsShares.IsNegative() {
		return sdkerrors.ErrInvalidRequest.Wrapf(
			"vesting accounts shares cannot be negative: %s",
			shares.VestingAccountsShares,
		)
	}

	return nil
}
