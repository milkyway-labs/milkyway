package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrInvalidGenesis          = errors.Register(ModuleName, 1, "invalid genesis state")
	ErrInvalidDeactivationTime = errors.Register(ModuleName, 2, "invalid deactivation time")
	ErrOperatorNotFound        = errors.Register(ModuleName, 3, "operator not found")
	ErrInsufficientShares      = errors.Register(ModuleName, 4, "insufficient delegation shares")
)
