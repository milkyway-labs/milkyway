package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	restakingtypes "github.com/milkyway-labs/milkyway/v9/x/restaking/types"
)

// NewGenesisState creates a new GenesisState instance.
func NewGenesisState(
	params Params,
	burnCoins []BurnCoins,
	usersInsuranceFunds []UserInsuranceFundEntry,
	lockedRepresentationDelegators []string,
	targetsCoveredLockedShares []TargetCoveredLockedSharesRecord,
) *GenesisState {
	return &GenesisState{
		Params:                         params,
		BurnCoins:                      burnCoins,
		UserInsuranceFunds:             usersInsuranceFunds,
		LockedRepresentationDelegators: lockedRepresentationDelegators,
		TargetsCoveredLockedShares:     targetsCoveredLockedShares,
	}
}

// DefaultGenesisState returns a default GenesisState.
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams(), nil, nil, nil, nil)
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

	for _, delegator := range g.LockedRepresentationDelegators {
		_, err := sdk.AccAddressFromBech32(delegator)
		if err != nil {
			return err
		}
	}

	for _, targetCoveredLockedShares := range g.TargetsCoveredLockedShares {
		if err := targetCoveredLockedShares.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Validate validates TargetCoveredLockedSharesRecord.
func (r *TargetCoveredLockedSharesRecord) Validate() error {
	if !(r.DelegationType == restakingtypes.DELEGATION_TYPE_POOL ||
		r.DelegationType == restakingtypes.DELEGATION_TYPE_OPERATOR ||
		r.DelegationType == restakingtypes.DELEGATION_TYPE_SERVICE) {
		return fmt.Errorf("invalid delegation type: %v", r.DelegationType)
	}

	if r.DelegationTargetID == 0 {
		return fmt.Errorf("invalid delegation target id")
	}

	if r.Shares.IsAnyNegative() {
		return restakingtypes.ErrInvalidShares
	}

	return nil
}
