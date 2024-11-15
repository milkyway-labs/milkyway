package types

import (
	"context"
)

type RestakingHooks interface {
	PoolRestakingHooks
	OperatorRestakingHooks
	ServiceRestakingHooks
	AfterUnbondingInitiated(ctx context.Context, unbondingDelegationID uint64) error
}

type PoolRestakingHooks interface {
	BeforePoolDelegationCreated(ctx context.Context, poolID uint32, delegator string) error
	BeforePoolDelegationSharesModified(ctx context.Context, poolID uint32, delegator string) error
	AfterPoolDelegationModified(ctx context.Context, poolID uint32, delegator string) error
	BeforePoolDelegationRemoved(ctx context.Context, poolID uint32, delegator string) error
}

type OperatorRestakingHooks interface {
	BeforeOperatorDelegationCreated(ctx context.Context, operatorID uint32, delegator string) error
	BeforeOperatorDelegationSharesModified(ctx context.Context, operatorID uint32, delegator string) error
	AfterOperatorDelegationModified(ctx context.Context, operatorID uint32, delegator string) error
	BeforeOperatorDelegationRemoved(ctx context.Context, operatorID uint32, delegator string) error
}

type ServiceRestakingHooks interface {
	BeforeServiceDelegationCreated(ctx context.Context, serviceID uint32, delegator string) error
	BeforeServiceDelegationSharesModified(ctx context.Context, serviceID uint32, delegator string) error
	AfterServiceDelegationModified(ctx context.Context, serviceID uint32, delegator string) error
	BeforeServiceDelegationRemoved(ctx context.Context, serviceID uint32, delegator string) error
}
