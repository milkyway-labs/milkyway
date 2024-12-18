package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v5/utils"
)

// NewGenesis creates a new genesis state
func NewGenesis(
	operatorsJoinedServices []OperatorJoinedServices,
	servicesAllowedOperators []ServiceAllowedOperators,
	servicesSecuringPools []ServiceSecuringPools,
	delegations []Delegation,
	unbondingDelegations []UnbondingDelegation,
	usersPreferencesEntries []UserPreferencesEntry,
	params Params,
) *GenesisState {
	return &GenesisState{
		OperatorsJoinedServices:  operatorsJoinedServices,
		ServicesAllowedOperators: servicesAllowedOperators,
		ServicesSecuringPools:    servicesSecuringPools,
		Delegations:              delegations,
		UnbondingDelegations:     unbondingDelegations,
		UsersPreferences:         usersPreferencesEntries,
		Params:                   params,
	}
}

// DefaultGenesis returns a default genesis state
func DefaultGenesis() *GenesisState {
	return NewGenesis(
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		DefaultParams(),
	)
}

// Validate performs basic validation of genesis data
func (g *GenesisState) Validate() error {
	if duplicate := findDuplicateOperatorJoinedServiceRecords(g.OperatorsJoinedServices); duplicate != nil {
		return fmt.Errorf("duplicated operator joined services in: %d", duplicate.OperatorID)
	}

	if duplicate := findDuplicateServiceAllowedOperators(g.ServicesAllowedOperators); duplicate != nil {
		return fmt.Errorf("duplicated service allowed operators, service id: %d", duplicate.ServiceID)
	}

	if duplicate := findDuplicateServiceSecuringPools(g.ServicesSecuringPools); duplicate != nil {
		return fmt.Errorf("duplicated service securing pools, service id: %d", duplicate.ServiceID)
	}

	for _, record := range g.OperatorsJoinedServices {
		err := record.Validate()
		if err != nil {
			return fmt.Errorf("invalid operator joined service record, operator id %d %s", record.OperatorID, err)
		}
	}

	// Validate the services allowed operators
	for _, entry := range g.ServicesAllowedOperators {
		err := entry.Validate()
		if err != nil {
			return fmt.Errorf("invalid service allowed operators with id %d: %s", entry.ServiceID, err)
		}
	}

	// Validate the services securing pools
	for _, entry := range g.ServicesSecuringPools {
		err := entry.Validate()
		if err != nil {
			return fmt.Errorf("invalid service securing pools with id %d: %s", entry.ServiceID, err)
		}
	}

	// Validate delegations
	for _, entry := range g.Delegations {
		err := entry.Validate()
		if err != nil {
			return fmt.Errorf("invalid delegation: %w", err)
		}
	}

	// Validate unbonding delegations
	for _, entry := range g.UnbondingDelegations {
		err := entry.Validate()
		if err != nil {
			return fmt.Errorf("invalid unbonding delegation: %w", err)
		}
	}

	for _, entry := range g.UsersPreferences {
		err := entry.Validate()
		if err != nil {
			return fmt.Errorf("invalid user preferences entry: %w", err)
		}
	}

	// Validate the params
	err := g.Params.Validate()
	if err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	return nil
}

func findDuplicateOperatorJoinedServiceRecords(records []OperatorJoinedServices) *OperatorJoinedServices {
	return utils.FindDuplicateFunc(records, func(a, b OperatorJoinedServices) bool {
		return a.OperatorID == b.OperatorID
	})
}

func findDuplicateServiceAllowedOperators(records []ServiceAllowedOperators) *ServiceAllowedOperators {
	return utils.FindDuplicateFunc(records, func(a, b ServiceAllowedOperators) bool {
		return a.ServiceID == b.ServiceID
	})
}

func findDuplicateServiceSecuringPools(records []ServiceSecuringPools) *ServiceSecuringPools {
	return utils.FindDuplicateFunc(records, func(a, b ServiceSecuringPools) bool {
		return a.ServiceID == b.ServiceID
	})
}

// --------------------------------------------------------------------------------------------------------------------

// NewOperatorJoinedServices creates a new instance of OperatorServiceIdRecord.
func NewOperatorJoinedServices(operatorID uint32, serviceIDs []uint32) OperatorJoinedServices {
	return OperatorJoinedServices{
		OperatorID: operatorID,
		ServiceIDs: serviceIDs,
	}
}

func (o *OperatorJoinedServices) Validate() error {
	if o.OperatorID == 0 {
		return fmt.Errorf("the operator id must be greater than 0")
	}

	for _, serviceID := range o.ServiceIDs {
		if serviceID == 0 {
			return fmt.Errorf("the service id must be greater than 0")
		}
	}
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// NewServiceAllowedOperators creates a new instance of ServiceAllowedOperators.
func NewServiceAllowedOperators(serviceID uint32, operatorIDs []uint32) ServiceAllowedOperators {
	return ServiceAllowedOperators{
		ServiceID:   serviceID,
		OperatorIDs: operatorIDs,
	}
}

func (r *ServiceAllowedOperators) Validate() error {
	if r.ServiceID == 0 {
		return fmt.Errorf("the service id must be greater than 0")
	}

	for _, operatorID := range r.OperatorIDs {
		if operatorID == 0 {
			return fmt.Errorf("the operator id must be greater than 0")
		}
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// NewServiceSecuringPools creates a new instance of ServiceSecuringPools.
func NewServiceSecuringPools(serviceID uint32, poolIDs []uint32) ServiceSecuringPools {
	return ServiceSecuringPools{
		ServiceID: serviceID,
		PoolIDs:   poolIDs,
	}
}

func (r *ServiceSecuringPools) Validate() error {
	if r.ServiceID == 0 {
		return fmt.Errorf("the service id must be greater than 0")
	}

	for _, poolID := range r.PoolIDs {
		if poolID == 0 {
			return fmt.Errorf("the pool id must be greater than 0")
		}
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// NewUserPreferencesEntry creates a new instance of UserPreferenceEntry.
func NewUserPreferencesEntry(userAddress string, preferences UserPreferences) UserPreferencesEntry {
	return UserPreferencesEntry{
		UserAddress: userAddress,
		Preferences: preferences,
	}
}

// Validate checks if the UserPreferenceEntry is valid
func (u *UserPreferencesEntry) Validate() error {
	_, err := sdk.AccAddressFromBech32(u.UserAddress)
	if err != nil {
		return fmt.Errorf("invalid user address: %w", err)
	}

	return u.Preferences.Validate()
}
