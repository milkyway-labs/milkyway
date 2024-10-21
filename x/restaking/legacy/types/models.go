package types

import (
	"fmt"

	"cosmossdk.io/math"

	"github.com/milkyway-labs/milkyway/utils"
)

// NewLegacyOperatorParams creates a new OperatorParams instance
func NewLegacyOperatorParams(commissionRate math.LegacyDec, joinedServicesIDs []uint32) LegacyOperatorParams {
	return LegacyOperatorParams{
		CommissionRate:    commissionRate,
		JoinedServicesIDs: joinedServicesIDs,
	}
}

// DefaultLegacyOperatorParams returns the default operator params
func DefaultLegacyOperatorParams() LegacyOperatorParams {
	return NewLegacyOperatorParams(math.LegacyZeroDec(), nil)
}

// Validate validates the operator params
func (p *LegacyOperatorParams) Validate() error {
	if p.CommissionRate.IsNegative() || p.CommissionRate.GT(math.LegacyOneDec()) {
		return fmt.Errorf("invalid commission rate: %s", p.CommissionRate.String())
	}

	if duplicate := utils.FindDuplicate(p.JoinedServicesIDs); duplicate != nil {
		return fmt.Errorf("duplicated joined service id: %v", duplicate)
	}

	for _, serviceID := range p.JoinedServicesIDs {
		if serviceID == 0 {
			return fmt.Errorf("invalid joined service id: %d", serviceID)
		}
	}
	return nil
}
