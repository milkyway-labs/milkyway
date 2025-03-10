package types

import (
	sdkmath "cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	DefaultInvestorsRewardRatio = sdkmath.LegacyNewDecWithPrec(5, 1) // 50%
)

// ValidateInvestorsRewardRatio validates the investors reward ratio.
func ValidateInvestorsRewardRatio(ratio sdkmath.LegacyDec) error {
	if ratio.IsNegative() {
		return sdkerrors.ErrInvalidRequest.Wrapf(
			"investors reward ratio cannot be negative: %s",
			ratio,
		)
	} else if ratio.GT(sdkmath.LegacyOneDec()) {
		return sdkerrors.ErrInvalidRequest.Wrapf(
			"investors reward ratio cannot be greater than one: %s",
			ratio,
		)
	}
	return nil
}
