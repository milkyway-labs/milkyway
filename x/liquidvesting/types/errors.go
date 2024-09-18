package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrInvalidInsurancePercentage               = errors.Register(ModuleName, 1, "invalid insurance percentage value")
	ErrNotMinter                                = errors.Register(ModuleName, 2, "sender don't have permission to mint tokens")
	ErrNotBurner                                = errors.Register(ModuleName, 3, "sender don't have permission to burn tokens")
	ErrInvalidDenom                             = errors.Register(ModuleName, 4, "invalid denom")
	ErrVestedRepresentationCannoteBeTransferred = errors.Register(ModuleName, 5, "vested representation can't be transferred")
	ErrInsufficientInsuranceFundBalance         = errors.Register(ModuleName, 6, "insufficient insurance fund balance")
	ErrInsufficientBalance                      = errors.Register(ModuleName, 7, "insufficient balance")
)
