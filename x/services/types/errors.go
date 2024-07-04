package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrInvalidGenesis       = errors.Register(ModuleName, 1, "invalid genesis state")
	ErrServiceNotFound      = errors.Register(ModuleName, 2, "service not found")
	ErrServiceAlreadyActive = errors.Register(ModuleName, 3, "service is already active")
	ErrServiceNotActive     = errors.Register(ModuleName, 4, "service is not active")
	ErrInsufficientShares   = errors.Register(ModuleName, 5, "insufficient delegation shares")
)
