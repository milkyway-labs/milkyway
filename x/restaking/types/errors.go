package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrInvalidGenesis              = errors.Register(ModuleName, 1, "invalid genesis state")
	ErrDelegationNotFound          = errors.Register(ModuleName, 2, "delegation not found")
	ErrDelegatorShareExRateInvalid = errors.Register(ModuleName, 34, "cannot delegate to pool/operator/service with invalid (zero) ex-rate")
)
