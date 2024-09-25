package types

import (
	"cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrOperatorNotFound             = errors.Wrap(sdkerrors.ErrNotFound, "operator not found")
	ErrInvalidGenesis               = errors.Register(ModuleName, 1, "invalid genesis state")
	ErrInsufficientShares           = errors.Register(ModuleName, 2, "insufficient delegation shares")
	ErrInvalidDeactivationTime      = errors.Register(ModuleName, 3, "invalid deactivation time")
	ErrOperatorNotActive            = errors.Register(ModuleName, 4, "operator not active")
	ErrInvalidExecuteMessagesSigner = errors.Register(ModuleName, 5, "expected operator as only signer for executing messages")
	ErrInvalidExecuteMsg            = errors.Register(ModuleName, 6, "invalid execute message")
	ErrUnroutableExecuteMsg         = errors.Register(ModuleName, 7, "unroutable execute message")
)
