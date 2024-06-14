package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type RestakingHooks interface {
	BeforePoolDelegationCreated(ctx sdk.Context, poolID uint32, delegator string) error
	BeforePoolDelegationSharesModified(ctx sdk.Context, poolID uint32, delegator string) error
	AfterPoolDelegationModified(ctx sdk.Context, poolID uint32, delegator string) error
}
