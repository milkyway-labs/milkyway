package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrInvalidInsurancePercentage              = errors.Register(ModuleName, 1, "invalid insurance percentage value")
	ErrNotMinter                               = errors.Register(ModuleName, 2, "sender don't have permission to mint tokens")
	ErrNotBurner                               = errors.Register(ModuleName, 3, "sender don't have permission to burn tokens")
	ErrInvalidAmount                           = errors.Register(ModuleName, 4, "invalid amount")
	ErrInvalidDenom                            = errors.Register(ModuleName, 5, "invalid denom")
	ErrVestedRepresentationCannotBeTransferred = errors.Register(ModuleName, 6, "vested representation can't be transferred")
	ErrInsufficientInsuranceFundBalance        = errors.Register(ModuleName, 7, "insufficient insurance fund balance")
	ErrInsufficientBalance                     = errors.Register(ModuleName, 8, "insufficient balance")
	ErrTransferBetweenTargetsNotAllowed        = errors.Register(ModuleName, 9, "transfer between restaking targets is not allowed")
)
