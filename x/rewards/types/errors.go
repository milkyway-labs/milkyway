package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrEmptyDelegationDistInfo = errors.Register(ModuleName, 2, "no delegation distribution info")
	ErrNoOperatorCommission    = errors.Register(ModuleName, 3, "no operator commission to withdraw")
)
