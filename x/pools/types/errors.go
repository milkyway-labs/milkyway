package types

import (
	"cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrPoolNotFound       = errors.Wrap(sdkerrors.ErrNotFound, "pool not found")
	ErrInvalidGenesis     = errors.Register(ModuleName, 1, "invalid genesis state")
	ErrInsufficientShares = errors.Register(ModuleName, 2, "insufficient delegation shares")
)
