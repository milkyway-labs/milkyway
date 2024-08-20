package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DelegationTarget is an interface that represents the target of a delegation (operator, pool, service, etc).
type DelegationTarget interface {
	GetID() uint32
	GetAddress() string
	GetTokens() sdk.Coins
	GetDelegatorShares() sdk.DecCoins
	InvalidExRate() bool
	GetSharesDenom(tokenDenom string) string
	TokensFromShares(shares sdk.DecCoins) sdk.DecCoins
	TokensFromSharesTruncated(shares sdk.DecCoins) sdk.DecCoins
	SharesFromTokens(amt sdk.Coins) (sdk.DecCoins, error)
	SharesFromTokensTruncated(tokens sdk.Coins) (sdk.DecCoins, error)
}

// DelegationGetter represents a function that allows to retrieve an existing delegation
type DelegationGetter func(ctx sdk.Context, receiverID uint32, delegator string) (Delegation, bool)

// DelegationBuilder represents a function that allows to build a new delegation
type DelegationBuilder func(targetID uint32, delegator string, shares sdk.DecCoins) Delegation

// DelegationUpdater represents a function that allows to update an existing delegation
type DelegationUpdater func(ctx sdk.Context, delegation Delegation) (newShares sdk.DecCoins, err error)

// DelegationHooks contains the hooks that can be called before and after a delegation is modified.
type DelegationHooks struct {
	BeforeDelegationSharesModified func(ctx sdk.Context, receiverID uint32, delegator string) error
	BeforeDelegationCreated        func(ctx sdk.Context, receiverID uint32, delegator string) error
	AfterDelegationModified        func(ctx sdk.Context, receiverID uint32, delegator string) error
	BeforeDelegationRemoved        func(ctx sdk.Context, receiverID uint32, delegator string) error
}

// DelegationData contains the data required to perform a delegation.
type DelegationData struct {
	Amount           sdk.Coins
	Delegator        string
	Target           DelegationTarget
	BuildDelegation  DelegationBuilder
	UpdateDelegation DelegationUpdater
	Hooks            DelegationHooks
}

type UnbondingDelegationBuilder func(
	delegatorAddress string, targetID uint32,
	creationHeight int64, minTime time.Time, balance sdk.Coins, id uint64,
) UnbondingDelegation

type UndelegationData struct {
	Amount                   sdk.Coins
	Delegator                string
	Target                   DelegationTarget
	BuildUnbondingDelegation UnbondingDelegationBuilder
	Hooks                    DelegationHooks
	Shares                   sdk.DecCoins
}
