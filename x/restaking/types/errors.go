package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrInvalidGenesis              = errors.Register(ModuleName, 1, "invalid genesis state")
	ErrInvalidShares               = errors.Register(ModuleName, 2, "invalid shares amount")
	ErrDelegatorShareExRateInvalid = errors.Register(ModuleName, 3, "cannot delegate to pool/operator/service with invalid (zero) ex-rate")
)
