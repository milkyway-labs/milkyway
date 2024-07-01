package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrInvalidGenesis          = errors.Register(ModuleName, 1, "invalid genesis state")
	ErrInvalidDeactivationTime = errors.Register(ModuleName, 2, "invalid deactivation time")
	ErrOperatorNotFound        = errors.Register(ModuleName, 3, "operator not found")
	ErrOperatorNotActive       = errors.Register(ModuleName, 4, "operator not active")
	ErrInsufficientShares      = errors.Register(ModuleName, 5, "insufficient delegation shares")
)
