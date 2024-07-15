package types

import (
	"cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrServiceNotFound      = errors.Wrap(sdkerrors.ErrNotFound, "service not found")
	ErrInvalidGenesis       = errors.Register(ModuleName, 1, "invalid genesis state")
	ErrInsufficientShares   = errors.Register(ModuleName, 2, "insufficient delegation shares")
	ErrServiceAlreadyActive = errors.Register(ModuleName, 3, "service is already active")
	ErrServiceNotActive     = errors.Register(ModuleName, 4, "service is not active")
)
