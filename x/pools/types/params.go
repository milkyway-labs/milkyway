package types

import (
	"fmt"

	"github.com/milkyway-labs/milkyway/utils"
)

// NewParams returns a new Params instance
func NewParams(allowedServiceIDs []uint32) Params {
	return Params{
		AllowedServiceIDs: allowedServiceIDs,
	}
}

// DefaultParams return a Params instance with default values set
func DefaultParams() Params {
	return NewParams(nil)
}

// Validate performs basic validation of params
func (p *Params) Validate() error {
	if duplicate := utils.FindDuplicate(p.AllowedServiceIDs); duplicate != nil {
		return fmt.Errorf("duplicated allowed service id: %v", duplicate)
	}

	for _, serviceID := range p.AllowedServiceIDs {
		if serviceID == 0 {
			return fmt.Errorf("invalid service id: %d", serviceID)
		}
	}
	return nil
}