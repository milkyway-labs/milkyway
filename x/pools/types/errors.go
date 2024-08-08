package types

import (
	"cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrPoolNotFound       = errors.Wrap(sdkerrors.ErrNotFound, "pool not found")
	ErrInvalidGenesis     = errors.Register(ModuleName, 1, "invalid genesis state")
	ErrInsufficientShares = errors.Register(ModuleName, 2, "insufficient delegation shares")
	ErrInvalidDenom       = errors.Register(ModuleName, 3, "invalid token denom")
	ErrMultipleTokens     = errors.Register(ModuleName, 4, "multiple tokens not allowed")
	ErrInvalidShares      = errors.Register(ModuleName, 5, "invalid shares")
)
