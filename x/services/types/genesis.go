package types

import (
	"fmt"

	"github.com/milkyway-labs/milkyway/v9/utils"
)

// NewGenesisState returns a new GenesisState instance
func NewGenesisState(nextServiceID uint32, services []Service, servicesParams []ServiceParamsRecord, params Params) *GenesisState {
	return &GenesisState{
		NextServiceID:  nextServiceID,
		Services:       services,
		ServicesParams: servicesParams,
		Params:         params,
	}
}

// DefaultGenesis returns a default GenesisState
func DefaultGenesis() *GenesisState {
	return NewGenesisState(1, nil, nil, DefaultParams())
}

// --------------------------------------------------------------------------------------------------------------------

// Validate validates the GenesisState and returns an error if it is invalid.
func (data *GenesisState) Validate() error {
	// Check for the next service ID
	if data.NextServiceID == 0 {
		return fmt.Errorf("invalid next service ID: %d", data.NextServiceID)
	}

	// Check for duplicated services
	if duplicate := findDuplicatedService(data.Services); duplicate != nil {
		return fmt.Errorf("duplicated service: %d", duplicate.ID)
	}

	// Validate services
	for _, service := range data.Services {
		err := service.Validate()
		if err != nil {
			return fmt.Errorf("invalid service with id %d: %s", service.ID, err)
		}
	}

	// Check for duplicated service params
	if duplicate := findDuplicatedServiceParamsRecord(data.ServicesParams); duplicate != nil {
		return fmt.Errorf("duplicated service params record: %d", duplicate.ServiceID)
	}

	// Validate the service params
	for _, serviceParams := range data.ServicesParams {
		err := serviceParams.Validate()
		if err != nil {
			return fmt.Errorf("invalid service params record for service %d: %w", serviceParams.ServiceID, err)
		}
	}

	// Validate params
	err := data.Params.Validate()
	if err != nil {
		return fmt.Errorf("invalid params: %s", err)
	}

	return nil
}

// findDuplicatedService returns the first duplicated service in the slice.
// If no duplicates are found, it returns nil instead.
func findDuplicatedService(services []Service) *Service {
	return utils.FindDuplicateFunc(services, func(a, b Service) bool {
		return a.ID == b.ID
	})
}

// findDuplicatedServiceParamsRecord returns the first duplicated ServiceParamsRecord in the slice.
// If no duplicates are found, it returns nil instead.
func findDuplicatedServiceParamsRecord(records []ServiceParamsRecord) *ServiceParamsRecord {
	return utils.FindDuplicateFunc(records, func(a, b ServiceParamsRecord) bool {
		return a.ServiceID == b.ServiceID
	})
}

// --------------------------------------------------------------------------------------------------------------------

// NewServiceParamsRecord returns a new ServiceParamsRecord instance.
func NewServiceParamsRecord(serviceID uint32, params ServiceParams) ServiceParamsRecord {
	return ServiceParamsRecord{
		ServiceID: serviceID,
		Params:    params,
	}
}

func (r *ServiceParamsRecord) Validate() error {
	if r.ServiceID == 0 {
		return fmt.Errorf("invalid service ID: %d", r.ServiceID)
	}

	return r.Params.Validate()
}
