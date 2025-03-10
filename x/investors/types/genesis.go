package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v9/utils"
)

// NewGenesisState creates a new GenesisState instance.
func NewGenesisState(investorsRewardRatio sdkmath.LegacyDec, vestingInvestorsAddrs []string,
) *GenesisState {
	return &GenesisState{
		InvestorsRewardRatio:      investorsRewardRatio,
		VestingInvestorsAddresses: vestingInvestorsAddrs,
	}
}

// DefaultGenesisState returns a default genesis state.
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultInvestorsRewardRatio, nil)
}

// Validate validates the genesis state.
func (gs *GenesisState) Validate() error {
	err := ValidateInvestorsRewardRatio(gs.InvestorsRewardRatio)
	if err != nil {
		return err
	}

	for _, investor := range gs.VestingInvestorsAddresses {
		_, err = sdk.AccAddressFromBech32(investor)
		if err != nil {
			return fmt.Errorf("invalid investor address: %w", err)
		}
	}

	duplicatedInvestor := utils.FindDuplicate(gs.VestingInvestorsAddresses)
	if duplicatedInvestor != nil {
		return fmt.Errorf("duplicated investor address: %s", *duplicatedInvestor)
	}
	return nil
}
