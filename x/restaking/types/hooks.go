package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type RestakingHooks interface {
	PoolRestakingHooks
	OperatorRestakingHooks
	ServiceRestakingHooks
}

type PoolRestakingHooks interface {
	BeforePoolDelegationCreated(ctx sdk.Context, poolID uint32, delegator string) error
	BeforePoolDelegationSharesModified(ctx sdk.Context, poolID uint32, delegator string) error
	AfterPoolDelegationModified(ctx sdk.Context, poolID uint32, delegator string) error
	BeforePoolDelegationRemoved(ctx sdk.Context, poolID uint32, delegator string) error
}

type OperatorRestakingHooks interface {
	BeforeOperatorDelegationCreated(ctx sdk.Context, operatorID uint32, delegator string) error
	BeforeOperatorDelegationSharesModified(ctx sdk.Context, operatorID uint32, delegator string) error
	AfterOperatorDelegationModified(ctx sdk.Context, operatorID uint32, delegator string) error
	BeforeOperatorDelegationRemoved(ctx sdk.Context, operatorID uint32, delegator string) error
}

type ServiceRestakingHooks interface {
	BeforeServiceDelegationCreated(ctx sdk.Context, serviceID uint32, delegator string) error
	BeforeServiceDelegationSharesModified(ctx sdk.Context, serviceID uint32, delegator string) error
	AfterServiceDelegationModified(ctx sdk.Context, serviceID uint32, delegator string) error
	BeforeServiceDelegationRemoved(ctx sdk.Context, serviceID uint32, delegator string) error
}
