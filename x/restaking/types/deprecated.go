package types

import (
	fmt "fmt"

	"cosmossdk.io/math"

	"github.com/milkyway-labs/milkyway/utils"
)

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

// --------------------------------------------------------------------------------------------------------------------

// NewLegacyServiceParamsRecord creates a new instance of LegacyServiceParamsRecord.
func NewLegacyServiceParamsRecord(serviceID uint32, params LegacyServiceParams) LegacyServiceParamsRecord {
	return LegacyServiceParamsRecord{
		ServiceID: serviceID,
		Params:    params,
	}
}

func (r *LegacyServiceParamsRecord) Validate() error {
	if r.ServiceID == 0 {
		return fmt.Errorf("the service id must be greater than 0")
	}
	return r.Params.Validate()
}
