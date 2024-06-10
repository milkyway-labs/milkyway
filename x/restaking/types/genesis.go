package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesis creates a new genesis state
func NewGenesis(
	poolsDelegations []PoolDelegationEntry,
	servicesDelegations []ServiceDelegationEntry,
	operatorsDelegations []OperatorDelegationEntry,
	params Params,
) *GenesisState {
	return &GenesisState{
		PoolsDelegations:     poolsDelegations,
		ServicesDelegations:  servicesDelegations,
		OperatorsDelegations: operatorsDelegations,
		Params:               params,
	}
}

// DefaultGenesis returns a default genesis state
func DefaultGenesis() *GenesisState {
	return NewGenesis(nil, nil, nil, DefaultParams())
}

// Validate performs basic validation of genesis data
func (g *GenesisState) Validate() error {
	// Validate pools delegations
	for _, entry := range g.PoolsDelegations {
		err := entry.Validate()
		if err != nil {
			return fmt.Errorf("invalid pool delegation entry: %w", err)
		}
	}

	// Validate services delegations
	for _, entry := range g.ServicesDelegations {
		err := entry.Validate()
		if err != nil {
			return fmt.Errorf("invalid service delegation entry: %w", err)
		}
	}

	// Validate operators delegations
	for _, entry := range g.OperatorsDelegations {
		err := entry.Validate()
		if err != nil {
			return fmt.Errorf("invalid operator delegation entry: %w", err)
		}
	}

	// Validate the params
	err := g.Params.Validate()
	if err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// NewPoolDelegationEntry returns a new PoolDelegationEntry
func NewPoolDelegationEntry(poolID uint32, userAddress string, amount sdk.Coin) PoolDelegationEntry {
	return PoolDelegationEntry{
		PoolID:      poolID,
		UserAddress: userAddress,
		Amount:      amount,
	}
}

// Validate performs basic validation of a pool delegation entry
func (e *PoolDelegationEntry) Validate() error {
	if e.PoolID == 0 {
		return fmt.Errorf("invalid pool id: %d", e.PoolID)
	}

	_, err := sdk.AccAddressFromBech32(e.UserAddress)
	if err != nil {
		return fmt.Errorf("invalid user address: %s", e.UserAddress)
	}

	if !e.Amount.IsValid() || e.Amount.IsZero() {
		return fmt.Errorf("invalid amount: %s", e.Amount)
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// NewServiceDelegationEntry returns a new ServiceDelegationEntry
func NewServiceDelegationEntry(serviceID uint32, userAddress string, amount sdk.Coin) ServiceDelegationEntry {
	return ServiceDelegationEntry{
		ServiceID:   serviceID,
		UserAddress: userAddress,
		Amount:      amount,
	}
}

// Validate performs basic validation of a service delegation entry
func (e *ServiceDelegationEntry) Validate() error {
	if e.ServiceID == 0 {
		return fmt.Errorf("invalid service id: %d", e.ServiceID)
	}

	_, err := sdk.AccAddressFromBech32(e.UserAddress)
	if err != nil {
		return fmt.Errorf("invalid user address: %s", e.UserAddress)
	}

	if !e.Amount.IsValid() || e.Amount.IsZero() {
		return fmt.Errorf("invalid amount: %s", e.Amount)
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// NewOperatorDelegationEntry returns a new OperatorDelegationEntry
func NewOperatorDelegationEntry(operatorID uint32, userAddress string, amount sdk.Coin) OperatorDelegationEntry {
	return OperatorDelegationEntry{
		OperatorID:  operatorID,
		UserAddress: userAddress,
		Amount:      amount,
	}
}

// Validate performs basic validation of an operator delegation entry
func (e *OperatorDelegationEntry) Validate() error {
	if e.OperatorID == 0 {
		return fmt.Errorf("invalid operator id: %d", e.OperatorID)
	}

	_, err := sdk.AccAddressFromBech32(e.UserAddress)
	if err != nil {
		return fmt.Errorf("invalid user address: %s", e.UserAddress)
	}

	if !e.Amount.IsValid() || e.Amount.IsZero() {
		return fmt.Errorf("invalid amount: %s", e.Amount)
	}

	return nil
}
