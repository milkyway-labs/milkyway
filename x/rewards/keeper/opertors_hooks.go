package keeper

import (
	"context"

	operatorstypes "github.com/milkyway-labs/milkyway/v6/x/operators/types"
	restakingtypes "github.com/milkyway-labs/milkyway/v6/x/restaking/types"
)

var _ operatorstypes.OperatorsHooks = OperatorsHooks{}

type OperatorsHooks struct {
	k *Keeper
}

func (k *Keeper) OperatorsHooks() OperatorsHooks {
	return OperatorsHooks{k}
}

// AfterOperatorRegistered implements operatorstypes.OperatorsHooks
func (h OperatorsHooks) AfterOperatorRegistered(ctx context.Context, operatorID uint32) error {
	return h.k.AfterDelegationTargetCreated(ctx, restakingtypes.DELEGATION_TYPE_OPERATOR, operatorID)
}

// AfterOperatorInactivatingStarted implements operatorstypes.OperatorsHooks
func (h OperatorsHooks) AfterOperatorInactivatingStarted(context.Context, uint32) error {
	return nil
}

// AfterOperatorInactivatingCompleted implements operatorstypes.OperatorsHooks
func (h OperatorsHooks) AfterOperatorInactivatingCompleted(context.Context, uint32) error {
	return nil
}

// AfterOperatorReactivated implements operatorstypes.OperatorsHooks
func (h OperatorsHooks) AfterOperatorReactivated(context.Context, uint32) error {
	return nil
}

// BeforeOperatorDeleted implements operatorstypes.OperatorsHooks
func (h OperatorsHooks) BeforeOperatorDeleted(ctx context.Context, operatorID uint32) error {
	return h.k.BeforeDelegationTargetRemoved(ctx, restakingtypes.DELEGATION_TYPE_OPERATOR, operatorID)
}
