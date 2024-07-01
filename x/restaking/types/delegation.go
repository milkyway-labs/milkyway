package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DelegationReceiver is an interface that represents the receiver of a delegation (operator, pool, service, etc).
type DelegationReceiver interface {
	GetID() uint32
	GetAddress() string
	InvalidExRate() bool
}

// DelegationGetter represents a function that allows to retrieve an existing delegation
type DelegationGetter func(ctx sdk.Context, receiverID uint32, delegator string) (Delegation, bool)

// DelegationBuilder represents a function that allows to build a new delegation
type DelegationBuilder func(receiverID uint32, delegator string) Delegation

// DelegationUpdater represents a function that allows to update an existing delegation
type DelegationUpdater func(ctx sdk.Context, delegation Delegation) (newShares sdk.DecCoins, err error)

// Delegation is an interface that represents a delegation object.
type Delegation interface {
	isDelegation()
}

// DelegationHooks contains the hooks that can be called before and after a delegation is modified.
type DelegationHooks struct {
	BeforeDelegationSharesModified func(ctx sdk.Context, receiverID uint32, delegator string) error
	BeforeDelegationCreated        func(ctx sdk.Context, receiverID uint32, delegator string) error
	AfterDelegationModified        func(ctx sdk.Context, receiverID uint32, delegator string) error
}

// DelegationData contains the data required to perform a delegation.
type DelegationData struct {
	Amount           sdk.Coins
	Delegator        string
	Receiver         DelegationReceiver
	GetDelegation    DelegationGetter
	BuildDelegation  DelegationBuilder
	UpdateDelegation DelegationUpdater
	Hooks            DelegationHooks
}
