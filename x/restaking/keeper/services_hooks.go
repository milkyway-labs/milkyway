package keeper

import (
	"context"

	"cosmossdk.io/collections"

	servicestypes "github.com/milkyway-labs/milkyway/v8/x/services/types"
)

type ServicesHooks struct {
	*Keeper
}

func (k *Keeper) ServicesHooks() servicestypes.ServicesHooks {
	return &ServicesHooks{k}
}

// ------------------------------------------------------------------------------

// BeforeServiceDeleted implements types.ServicesHooks.
func (h *ServicesHooks) BeforeServiceDeleted(ctx context.Context, serviceID uint32) error {
	// After the service has been deleted
	// we remove the data that we keep in the x/restaking
	// associated to this service.

	// Get the iterator to iterate over the operators that have joined this
	// service
	serviceValidatingOperatorsIter, err := h.operatorJoinedServices.Indexes.Service.Iterate(ctx, collections.NewPrefixedPairRange[uint32, uint32](serviceID))
	if err != nil {
		return err
	}
	defer serviceValidatingOperatorsIter.Close()

	// Get all the keys to remove
	toRemoveOperatorJoinedServices, err := serviceValidatingOperatorsIter.PrimaryKeys()
	if err != nil {
		return err
	}
	for _, key := range toRemoveOperatorJoinedServices {
		err = h.operatorJoinedServices.Remove(ctx, key)
		if err != nil {
			return err
		}
	}

	// Get the iterator to iterate over the operators that are
	// allowed to secure this service
	serviceOperatorsAllowListIter, err := h.ServiceAllowedOperatorsIterator(ctx, serviceID)
	if err != nil {
		return err
	}
	defer serviceOperatorsAllowListIter.Close()

	// Get all the keys to remove
	toRemoveServiceAllowedOperators, err := serviceOperatorsAllowListIter.Keys()
	if err != nil {
		return err
	}
	for _, key := range toRemoveServiceAllowedOperators {
		err = h.serviceOperatorsAllowList.Remove(ctx, key)
		if err != nil {
			return err
		}
	}

	// Get the iterator to iterate over the list of pools from
	// which the service is allowed to borrow security
	serviceSecuringPoolsIter, err := h.ServiceSecuringPoolsIterator(ctx, serviceID)
	if err != nil {
		return err
	}
	defer serviceSecuringPoolsIter.Close()

	// Get all the keys to remove
	toRemoveServiceSecuringPools, err := serviceSecuringPoolsIter.Keys()
	if err != nil {
		return err
	}
	for _, key := range toRemoveServiceSecuringPools {
		err = h.serviceSecuringPools.Remove(ctx, key)
		if err != nil {
			return err
		}
	}

	return nil
}

// AfterServiceDeactivated implements types.ServicesHooks.
func (h *ServicesHooks) AfterServiceDeactivated(ctx context.Context, serviceID uint32) error {
	return nil
}

// AfterServiceActivated implements types.ServicesHooks.
func (h *ServicesHooks) AfterServiceActivated(ctx context.Context, serviceID uint32) error {
	return nil
}

// AfterServiceCreated implements types.ServicesHooks.
func (h *ServicesHooks) AfterServiceCreated(ctx context.Context, serviceID uint32) error {
	return nil
}

// AfterServiceAccreditationModified implements types.ServicesHooks.
func (h *ServicesHooks) AfterServiceAccreditationModified(ctx context.Context, serviceID uint32) error {
	return nil
}
