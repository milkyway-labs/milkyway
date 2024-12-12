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
	ErrServiceNotJoinedByOperator     = errors.Register(ModuleName, 9, "the service has not been joined by the operator")
	ErrOperatorAlreadyAllowed         = errors.Register(ModuleName, 10, "operator already allowed")
	ErrOperatorNotAllowed             = errors.Register(ModuleName, 11, "operator not allowed")
	ErrPoolAlreadySecuringService     = errors.Register(ModuleName, 12, "pool already securing the service")
	ErrPoolNotSecuringService         = errors.Register(ModuleName, 13, "pool not securing the service")
	ErrDenomNotRestakable             = errors.Register(ModuleName, 14, "denom not restakable")
	ErrRestakingCapExceeded           = errors.Register(ModuleName, 15, "restaking cap exceeded")
)
