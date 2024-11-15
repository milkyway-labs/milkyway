package keeper

import (
	"context"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
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

// AfterOperatorDeleted implements operatorstypes.OperatorsHooks
func (h OperatorsHooks) AfterOperatorDeleted(ctx context.Context, operatorID uint32) error {
	return h.k.AfterDelegationTargetRemoved(ctx, restakingtypes.DELEGATION_TYPE_OPERATOR, operatorID)
}
