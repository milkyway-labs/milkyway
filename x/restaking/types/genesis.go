package types

import (
	"fmt"

	"github.com/milkyway-labs/milkyway/utils"
)

// NewGenesis creates a new genesis state
func NewGenesis(
	operatorsJoinedServices []OperatorJoinedServicesRecord,
	servicesParamsRecords []ServiceParamsRecord,
	delegations []Delegation,
	unbondingDelegations []UnbondingDelegation,
	params Params,
) *GenesisState {
	return &GenesisState{
		OperatorsJoinedServices: operatorsJoinedServices,
		ServicesParams:          servicesParamsRecords,
		Delegations:             delegations,
		UnbondingDelegations:    unbondingDelegations,
		Params:                  params,
	}
}

// DefaultGenesis returns a default genesis state
func DefaultGenesis() *GenesisState {
	return NewGenesis(
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

	if duplicate := findDuplicateServiceParamsRecords(g.ServicesParams); duplicate != nil {
		return fmt.Errorf("duplicated service params: %d", duplicate.ServiceID)
	}

	for _, record := range g.OperatorsJoinedServices {
		err := record.Validate()
		if err != nil {
			return fmt.Errorf("invalid operator joined service record, operator id %d %s", record.OperatorID, err)
		}
	}

	for _, record := range g.ServicesParams {
		if record.ServiceID == 0 {
			return fmt.Errorf("invalid service id: %d", record.ServiceID)
		}
		err := record.Params.Validate()
		if err != nil {
			return fmt.Errorf("invalid service params with id %d: %s", record.ServiceID, err)
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

func findDuplicateServiceParamsRecords(records []ServiceParamsRecord) *ServiceParamsRecord {
	return utils.FindDuplicateFunc(records, func(a, b ServiceParamsRecord) bool {
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

// NewServiceParamsRecord creates a new instance of ServiceParamsRecord.
func NewServiceParamsRecord(serviceID uint32, params ServiceParams) ServiceParamsRecord {
	return ServiceParamsRecord{
		ServiceID: serviceID,
		Params:    params,
	}
}

func (r *ServiceParamsRecord) Validate() error {
	if r.ServiceID == 0 {
		return fmt.Errorf("the service id must be greater than 0")
	}
	return r.Params.Validate()
}
