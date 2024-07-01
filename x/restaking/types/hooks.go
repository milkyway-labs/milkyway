package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type RestakingHooks interface {
	PoolRestakingHooks
	OperatorRestakingHooks
}

type PoolRestakingHooks interface {
	BeforePoolDelegationCreated(ctx sdk.Context, poolID uint32, delegator string) error
	BeforePoolDelegationSharesModified(ctx sdk.Context, poolID uint32, delegator string) error
	AfterPoolDelegationModified(ctx sdk.Context, poolID uint32, delegator string) error
}

type OperatorRestakingHooks interface {
	BeforeOperatorDelegationCreated(ctx sdk.Context, operatorID uint32, delegator string) error
	BeforeOperatorDelegationSharesModified(ctx sdk.Context, operatorID uint32, delegator string) error
	AfterOperatorDelegationModified(ctx sdk.Context, operatorID uint32, delegator string) error
}
