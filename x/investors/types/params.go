package types

import (
	sdkmath "cosmossdk.io/math"
)

var (
	DefaultInvestorsRewardRatio = sdkmath.LegacyNewDecWithPrec(5, 1) // 50%
)

// ValidateInvestorsRewardRatio validates the investors reward ratio.
func ValidateInvestorsRewardRatio(ratio sdkmath.LegacyDec) error {
	if ratio.IsNegative() {
		return ErrInvalidInvestorsRewardRatio.Wrapf("ratio cannot be negative: %s", ratio)
	} else if ratio.GT(sdkmath.LegacyOneDec()) {
		return ErrInvalidInvestorsRewardRatio.Wrapf("ratio cannot be greater than one: %s", ratio)
	}
	return nil
}
