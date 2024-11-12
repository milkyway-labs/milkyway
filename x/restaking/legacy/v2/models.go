package v2

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

// --------------------------------------------------------------------------------------------------------------------

// NewLegacyServiceParams creates a new LegacyServiceParams instance
func NewLegacyServiceParams(
	slashFraction math.LegacyDec,
	whitelistedPoolsIDs []uint32,
	whitelistedOperatorsIDs []uint32,
) LegacyServiceParams {
	return LegacyServiceParams{
		SlashFraction:           slashFraction,
		WhitelistedPoolsIDs:     whitelistedPoolsIDs,
		WhitelistedOperatorsIDs: whitelistedOperatorsIDs,
	}
}

// DefaultServiceParams returns the default service params
func DefaultLegacyServiceParams() LegacyServiceParams {
	return NewLegacyServiceParams(math.LegacyZeroDec(), nil, nil)
}

// Validate validates the service params
func (p *LegacyServiceParams) Validate() error {
	if p.SlashFraction.IsNegative() || p.SlashFraction.GT(math.LegacyOneDec()) {
		return fmt.Errorf("invalid slash fraction %s", p.SlashFraction)
	}

	if duplicate := utils.FindDuplicate(p.WhitelistedPoolsIDs); duplicate != nil {
		return fmt.Errorf("duplicated whitelisted pool id: %v", duplicate)
	}

	if duplicate := utils.FindDuplicate(p.WhitelistedOperatorsIDs); duplicate != nil {
		return fmt.Errorf("duplicated whitelisted operator id: %v", duplicate)
	}

	for _, poolID := range p.WhitelistedPoolsIDs {
		if poolID == 0 {
			return fmt.Errorf("invalid whitelisted pool id: %d", poolID)
		}
	}

	for _, operatorID := range p.WhitelistedOperatorsIDs {
		if operatorID == 0 {
			return fmt.Errorf("invalid whitelisted operator id: %d", operatorID)
		}
	}
	return nil
}
