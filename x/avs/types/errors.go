package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrInvalidGenesis = errors.Register(ModuleName, 1, "invalid genesis state")
)
