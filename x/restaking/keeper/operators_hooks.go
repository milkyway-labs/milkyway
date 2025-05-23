package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"

	operatorstypes "github.com/milkyway-labs/milkyway/v12/x/operators/types"
	servicestypes "github.com/milkyway-labs/milkyway/v12/x/services/types"
)

var _ operatorstypes.OperatorsHooks = &OperatorsHooks{}

type OperatorsHooks struct {
	*Keeper
}

func (k *Keeper) OperatorsHooks() operatorstypes.OperatorsHooks {
	return &OperatorsHooks{k}
}

// ------------------------------------------------------------------------------

// BeforeOperatorDeleted implements types.OperatorsHooks.
func (o *OperatorsHooks) BeforeOperatorDeleted(ctx context.Context, operatorID uint32) error {
	// After the operator has been deleted
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

	// Remove the operator from the services allow list where has been
	// added
	err = o.removeOperatorFromServicesAllowList(ctx, operatorID)
	if err != nil {
		return err
	}

	return nil
}

func (o *OperatorsHooks) removeOperatorFromServicesAllowList(ctx context.Context, operatorID uint32) error {
	// Get all the keys to remove
	var toRemoveKeys []collections.Pair[uint32, uint32]
	err := o.IterateAllServicesAllowedOperators(ctx, func(serviceID uint32, oID uint32) (stop bool, err error) {
		if oID == operatorID {
			toRemoveKeys = append(toRemoveKeys, collections.Join(serviceID, oID))
		}
		return false, nil
	})
	if err != nil {
		return err
	}

	// Iterate over the keys and remove them from the service operators allow list
	for _, key := range toRemoveKeys {
		// Remove the operator from the service allow list
		err := o.serviceOperatorsAllowList.Remove(ctx, key)
		if err != nil {
			return err
		}

		// Since we may have removed the last operator from the service allow
		// list lets check if is now empty and in this case we have to disable
		// the service to prevent unwanted operators to join.
		serviceID := key.K1()
		isConfigured, err := o.IsServiceOperatorsAllowListConfigured(ctx, serviceID)
		if err != nil {
			return err
		}
		if !isConfigured {
			service, err := o.servicesKeeper.GetService(ctx, serviceID)
			if err != nil {
				if errors.Is(err, collections.ErrNotFound) {
					return errorsmod.Wrapf(servicestypes.ErrServiceNotFound, "service %d not found", serviceID)
				}
				return err
			}

			if !service.IsActive() {
				// The service is not active, nothing to do
				continue
			}

			// The service is active and its operators allow list has become
			// empty, deactivate the service.
			err = o.servicesKeeper.DeactivateService(ctx, serviceID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// AfterOperatorInactivatingCompleted implements types.OperatorsHooks.
func (o *OperatorsHooks) AfterOperatorInactivatingCompleted(ctx context.Context, operatorID uint32) error {
	return nil
}

// AfterOperatorInactivatingStarted implements types.OperatorsHooks.
func (o *OperatorsHooks) AfterOperatorInactivatingStarted(ctx context.Context, operatorID uint32) error {
	return nil
}

// AfterOperatorReactivated implements types.OperatorsHooks.
func (o *OperatorsHooks) AfterOperatorReactivated(ctx context.Context, operatorID uint32) error {
	return nil
}

// AfterOperatorRegistered implements types.OperatorsHooks.
func (o *OperatorsHooks) AfterOperatorRegistered(ctx context.Context, operatorID uint32) error {
	return nil
}
