package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrInvalidGenesis                           = errors.Register(ModuleName, 1, "invalid genesis state")
	ErrNotMinter                                = errors.Register(ModuleName, 2, "sender don't have permission to mint tokens")
	ErrNotBurner                                = errors.Register(ModuleName, 3, "sender don't have permission to burn tokens")
	ErrInvalidDenom                             = errors.Register(ModuleName, 4, "invalid denom")
	ErrDenomMetadataNotFound                    = errors.Register(ModuleName, 5, "denom metadata not found")
	ErrVestedRepresentationCannoteBeTransferred = errors.Register(ModuleName, 6, "vested representation can't be transferred")
	ErrInsufficentInsuranceFundBalance          = errors.Register(ModuleName, 7, "insufficient insurance fund balance")
)
