package keeper

import (
	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
)

var _ operatorstypes.OperatorsHooks = &OperatorsHooks{}

type OperatorsHooks struct {
	*Keeper
}

func (k *Keeper) OperatorsHooks() operatorstypes.OperatorsHooks {
	return &OperatorsHooks{k}
}

// ------------------------------------------------------------------------------

// AfterOperatorDeleted implements types.OperatorsHooks.
func (o *OperatorsHooks) AfterOperatorDeleted(ctx sdk.Context, operatorID uint32) error {
	// After the operator has completed its inactivation
	// we remove the data that we keep in the x/restaking module that are linked
	// to the operator.

	// Wipe the list of services that this operator has joined
	iter, err := o.operatorJoinedServices.Iterate(ctx, collections.NewPrefixedPairRange[uint32, uint32](operatorID))
	if err != nil {
		return err
	}
	defer iter.Close()

	toRemoveOperatorJoinedServices, err := iter.Keys()
	if err != nil {
		return err
	}
	for _, key := range toRemoveOperatorJoinedServices {
		err = o.operatorJoinedServices.Remove(ctx, key)
		if err != nil {
			return err
		}
	}
	return nil
}

// AfterOperatorInactivatingCompleted implements types.OperatorsHooks.
func (o *OperatorsHooks) AfterOperatorInactivatingCompleted(ctx sdk.Context, operatorID uint32) error {
	return nil
}

// AfterOperatorInactivatingStarted implements types.OperatorsHooks.
func (o *OperatorsHooks) AfterOperatorInactivatingStarted(ctx sdk.Context, operatorID uint32) error {
	return nil
}

// AfterOperatorRegistered implements types.OperatorsHooks.
func (o *OperatorsHooks) AfterOperatorRegistered(ctx sdk.Context, operatorID uint32) error {
	return nil
}
