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

// ------------------------------------------------------------------------------

func (k *Keeper) OperatorsHooks() operatorstypes.OperatorsHooks {
	return &OperatorsHooks{k}
}

// ------------------------------------------------------------------------------

// AfterOperatorInactivatingCompleted implements types.OperatorsHooks.
func (o *OperatorsHooks) AfterOperatorInactivatingCompleted(ctx sdk.Context, operatorID uint32) error {
	// After the operator has completed its inactivation
	// we remove the data that we keep in the x/restaking module that is linked
	// to the operator.

	iter, err := o.operatorJoinedServices.Iterate(ctx, collections.NewPrefixedPairRange[uint32, uint32](operatorID))
	if err != nil {
		return err
	}
	defer iter.Close()

	// Iterate over the operator's joined service and remove all records
	for ; iter.Valid(); iter.Next() {
		operatorServicePair, err := iter.Key()
		if err != nil {
			return err
		}
		err = o.operatorJoinedServices.Remove(ctx, operatorServicePair)
		if err != nil {
			return err
		}
	}
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
