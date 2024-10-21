package types

import (
	"fmt"

	"github.com/milkyway-labs/milkyway/utils"
)

// NewGenesis creates a new genesis state
func NewGenesis(
	operatorsJoinedServices []OperatorJoinedServicesRecord,
	servicesWhitelistedOperators []ServiceWhitelistedOperators,
	servicesWhitelistedPools []ServiceWhitelistedPools,
	delegations []Delegation,
	unbondingDelegations []UnbondingDelegation,
	params Params,
) *GenesisState {
	return &GenesisState{
		OperatorsJoinedServices:      operatorsJoinedServices,
		ServicesWhitelistedOperators: servicesWhitelistedOperators,
		ServicesWhitelistedPools:     servicesWhitelistedPools,
		Delegations:                  delegations,
		UnbondingDelegations:         unbondingDelegations,
		Params:                       params,
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
		DefaultParams(),
	)
}

// Validate performs basic validation of genesis data
func (g *GenesisState) Validate() error {
	if duplicate := findDuplicateOperatorJoinedServiceRecords(g.OperatorsJoinedServices); duplicate != nil {
		return fmt.Errorf("duplicated operator joined services in: %d", duplicate.OperatorID)
	}

	if duplicate := findDuplicateServiceWhitelistedOperators(g.ServicesWhitelistedOperators); duplicate != nil {
		return fmt.Errorf("duplicated service whitelisted operators, service id: %d", duplicate.ServiceID)
	}

	if duplicate := findDuplicateServiceWhitelistedPools(g.ServicesWhitelistedPools); duplicate != nil {
		return fmt.Errorf("duplicated service whitelisted pools, service id: %d", duplicate.ServiceID)
	}

	for _, record := range g.OperatorsJoinedServices {
		err := record.Validate()
		if err != nil {
			return fmt.Errorf("invalid operator joined service record, operator id %d %s", record.OperatorID, err)
		}
	}

	// Validate the whitelisted services operators
	for _, entry := range g.ServicesWhitelistedOperators {
		err := entry.Validate()
		if err != nil {
			return fmt.Errorf("invalid service whitelisted operators with id %d: %s", entry.ServiceID, err)
		}
	}

	// Validate the whitelisted services pools
	for _, entry := range g.ServicesWhitelistedPools {
		err := entry.Validate()
		if err != nil {
			return fmt.Errorf("invalid service whitelisted pools with id %d: %s", entry.ServiceID, err)
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

	// Validate the params
	err := g.Params.Validate()
	if err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	return nil
}

func findDuplicateOperatorJoinedServiceRecords(records []OperatorJoinedServicesRecord) *OperatorJoinedServicesRecord {
	return utils.FindDuplicateFunc(records, func(a, b OperatorJoinedServicesRecord) bool {
		return a.OperatorID == b.OperatorID
	})
}

func findDuplicateServiceWhitelistedOperators(records []ServiceWhitelistedOperators) *ServiceWhitelistedOperators {
	return utils.FindDuplicateFunc(records, func(a, b ServiceWhitelistedOperators) bool {
		return a.ServiceID == b.ServiceID
	})
}

func findDuplicateServiceWhitelistedPools(records []ServiceWhitelistedPools) *ServiceWhitelistedPools {
	return utils.FindDuplicateFunc(records, func(a, b ServiceWhitelistedPools) bool {
		return a.ServiceID == b.ServiceID
	})
}

// --------------------------------------------------------------------------------------------------------------------

// NewOperatorJoinedServicesRecord creates a new instance of OperatorServiceIdRecord.
func NewOperatorJoinedServicesRecord(operatorID uint32, joinedServices OperatorJoinedServices) OperatorJoinedServicesRecord {
	return OperatorJoinedServicesRecord{
		OperatorID:     operatorID,
		JoinedServices: joinedServices,
	}
}

func (o *OperatorJoinedServicesRecord) Validate() error {
	if o.OperatorID == 0 {
		return fmt.Errorf("the operator id must be greater than 0")
	}
	return o.JoinedServices.Validate()
}

// --------------------------------------------------------------------------------------------------------------------

// NewServiceWhitelistedOperators creates a new instance of ServiceWhitelistedOperators.
func NewServiceWhitelistedOperators(serviceID uint32, operatorIDs []uint32) ServiceWhitelistedOperators {
	return ServiceWhitelistedOperators{
		ServiceID:   serviceID,
		OperatorIDs: operatorIDs,
	}
}

func (r *ServiceWhitelistedOperators) Validate() error {
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

// NewServiceWhitelistedPools creates a new instance of ServiceWhitelistedPools.
func NewServiceWhitelistedPools(serviceID uint32, poolIDs []uint32) ServiceWhitelistedPools {
	return ServiceWhitelistedPools{
		ServiceID: serviceID,
		PoolIDs:   poolIDs,
	}
}

func (r *ServiceWhitelistedPools) Validate() error {
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
