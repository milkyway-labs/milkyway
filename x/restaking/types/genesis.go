package types

import (
	"fmt"

	"github.com/milkyway-labs/milkyway/utils"
)

// NewGenesis creates a new genesis state
func NewGenesis(
	operatorParamsRecords []OperatorParamsRecord, serviceParamsRecords []ServiceParamsRecord,
	delegations []Delegation, params Params) *GenesisState {
	return &GenesisState{
		OperatorParams: operatorParamsRecords,
		ServiceParams:  serviceParamsRecords,
		Delegations:    delegations,
		Params:         params,
	}
}

// DefaultGenesis returns a default genesis state
func DefaultGenesis() *GenesisState {
	return NewGenesis(nil, nil, nil, DefaultParams())
}

// Validate performs basic validation of genesis data
func (g *GenesisState) Validate() error {
	if duplicate := findDuplicateOperatorParamsRecords(g.OperatorParams); duplicate != nil {
		return fmt.Errorf("duplicated operator params: %d", duplicate.OperatorID)
	}

	if duplicate := findDuplicateServiceParamsRecords(g.ServiceParams); duplicate != nil {
		return fmt.Errorf("duplicated service params: %d", duplicate.ServiceID)
	}

	for _, record := range g.OperatorParams {
		if record.OperatorID == 0 {
			return fmt.Errorf("invalid operator id: %d", record.OperatorID)
		}
		err := record.Params.Validate()
		if err != nil {
			return fmt.Errorf("invalid operator params with id %d: %s", record.OperatorID, err)
		}
	}

	for _, record := range g.ServiceParams {
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

	// Validate the params
	err := g.Params.Validate()
	if err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	return nil
}

func findDuplicateOperatorParamsRecords(records []OperatorParamsRecord) *OperatorParamsRecord {
	return utils.FindDuplicateFunc(records, func(a, b OperatorParamsRecord) bool {
		return a.OperatorID == b.OperatorID
	})
}

func findDuplicateServiceParamsRecords(records []ServiceParamsRecord) *ServiceParamsRecord {
	return utils.FindDuplicateFunc(records, func(a, b ServiceParamsRecord) bool {
		return a.ServiceID == b.ServiceID
	})
}
