package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrInvalidInvestorsRewardRatio = errorsmod.Register(ModuleName, 2, "invalid investors reward ratio")
)
