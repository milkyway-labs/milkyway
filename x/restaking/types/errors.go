package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrInvalidGenesis                 = errors.Register(ModuleName, 1, "invalid genesis state")
	ErrInvalidShares                  = errors.Register(ModuleName, 2, "invalid shares amount")
	ErrDelegatorShareExRateInvalid    = errors.Register(ModuleName, 3, "cannot delegate to pool/operator/service with invalid (zero) ex-rate")
	ErrDelegationNotFound             = errors.Register(ModuleName, 4, "delegation not found")
	ErrNotEnoughDelegationShares      = errors.Register(ModuleName, 5, "not enough delegation shares")
	ErrInvalidDelegationType          = errors.Register(ModuleName, 6, "invalid delegation type")
	ErrNoUnbondingDelegation          = errors.Register(ModuleName, 7, "no unbonding delegation found")
	ErrServiceAlreadyJoinedByOperator = errors.Register(ModuleName, 8, "service already joined by operator")
)
