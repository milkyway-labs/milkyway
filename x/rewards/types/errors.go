package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrEmptyDelegationDistInfo = errors.Register(ModuleName, 2, "no delegation distribution info")
)
