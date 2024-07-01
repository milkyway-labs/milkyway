package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DelegationReceiver interface {
	GetID() uint32
	GetAddress() string
	InvalidExRate() bool
}

type DelegationGetter func(ctx sdk.Context, receiverID uint32, delegator string) (Delegation, bool)

type DelegationBuilder func(receiverID uint32, delegator string, shares sdkmath.LegacyDec) Delegation

type DelegationUpdater func(ctx sdk.Context, delegation Delegation) (newShares sdkmath.LegacyDec, err error)

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
