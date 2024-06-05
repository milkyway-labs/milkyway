package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrInvalidGenesis  = errors.Register(ModuleName, 1, "invalid genesis state")
	ErrServiceNotFound = errors.Register(ModuleName, 2, "service not found")
)
