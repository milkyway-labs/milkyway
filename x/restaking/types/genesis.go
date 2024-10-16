package types

import (
	"fmt"

	"github.com/milkyway-labs/milkyway/utils"
)

// NewGenesis creates a new genesis state
func NewGenesis(
	operatorsParamsRecords []OperatorSecuredServicesRecord,
	servicesParamsRecords []ServiceParamsRecord,
	delegations []Delegation,
	unbondingDelegations []UnbondingDelegation,
	params Params,
) *GenesisState {
	return &GenesisState{
		OperatorsSecuredServices: operatorsParamsRecords,
		ServicesParams:           servicesParamsRecords,
		Delegations:              delegations,
		UnbondingDelegations:     unbondingDelegations,
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
		DefaultParams(),
	)
}

// Validate performs basic validation of genesis data
func (g *GenesisState) Validate() error {
	if duplicate := findDuplicateOperatorSecuredServiceRecords(g.OperatorsSecuredServices); duplicate != nil {
		return fmt.Errorf("duplicated operator params: %d", duplicate.OperatorID)
	}

	if duplicate := findDuplicateServiceParamsRecords(g.ServicesParams); duplicate != nil {
		return fmt.Errorf("duplicated service params: %d", duplicate.ServiceID)
	}

	for _, record := range g.OperatorsSecuredServices {
		err := record.Validate()
		if err != nil {
			return fmt.Errorf("invalid operator secured service record, operator id %d %s", record.OperatorID, err)
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

func findDuplicateOperatorSecuredServiceRecords(records []OperatorSecuredServicesRecord) *OperatorSecuredServicesRecord {
	return utils.FindDuplicateFunc(records, func(a, b OperatorSecuredServicesRecord) bool {
		return a.OperatorID == b.OperatorID
	})
}

func findDuplicateServiceParamsRecords(records []ServiceParamsRecord) *ServiceParamsRecord {
	return utils.FindDuplicateFunc(records, func(a, b ServiceParamsRecord) bool {
		return a.ServiceID == b.ServiceID
	})
}

// --------------------------------------------------------------------------------------------------------------------

// NewOperatorServiceIdRecord creates a new instance of OperatorServiceIdRecord.
func NewOperatorSecuredServicesRecord(operatorID uint32, securedServices OperatorSecuredServices) OperatorSecuredServicesRecord {
	return OperatorSecuredServicesRecord{
		OperatorID:      operatorID,
		SecuredServices: securedServices,
	}
}

func (o *OperatorSecuredServicesRecord) Validate() error {
	if o.OperatorID == 0 {
		return fmt.Errorf("the operator id must be greater than 0")
	}
	return o.SecuredServices.Validate()
}
