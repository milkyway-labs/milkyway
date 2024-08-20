package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrTickerNotFound = errors.Register(ModuleName, 2, "ticker not found")
)
