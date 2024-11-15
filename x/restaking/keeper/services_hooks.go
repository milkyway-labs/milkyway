package keeper

import (
	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

type ServicesHooks struct {
	*Keeper
}

func (k *Keeper) ServicesHooks() servicestypes.ServicesHooks {
	return &ServicesHooks{k}
}

// ------------------------------------------------------------------------------

// AfterServiceDeactivated implements types.ServicesHooks.
func (h *ServicesHooks) AfterServiceDeleted(ctx sdk.Context, serviceID uint32) error {
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
	serviceOperatorsAllowListIter, err := h.serviceOperatorsAllowList.Iterate(ctx, collections.NewPrefixedPairRange[uint32, uint32](serviceID))
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
	serviceSecuringPoolsIter, err := h.serviceSecuringPools.Iterate(ctx, collections.NewPrefixedPairRange[uint32, uint32](serviceID))
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
func (h *ServicesHooks) AfterServiceDeactivated(ctx sdk.Context, serviceID uint32) error {
	return nil
}

// AfterServiceActivated implements types.ServicesHooks.
func (h *ServicesHooks) AfterServiceActivated(ctx sdk.Context, serviceID uint32) error {
	return nil
}

// AfterServiceCreated implements types.ServicesHooks.
func (h *ServicesHooks) AfterServiceCreated(ctx sdk.Context, serviceID uint32) error {
	return nil
}

// AfterServiceAccreditationModified implements types.ServicesHooks.
func (h *ServicesHooks) AfterServiceAccreditationModified(ctx sdk.Context, serviceID uint32, accredited bool) error {
	store := ctx.KVStore(h.storeKey)
	iter := storetypes.KVStorePrefixIterator(store, types.DelegationsByServiceIDStorePrefix(serviceID))
	defer iter.Close()
	for iter.Valid() {
		_, delegator, err := types.ParseDelegationsByServiceIDKey(iter.Key())
		if err != nil {
			return err
		}
		preferences, err := h.GetUserPreferences(ctx, delegator)
		if err != nil {
			return err
		}
		trustedBefore := preferences.IsServiceTrusted(serviceID, !accredited)
		trustedAfter := preferences.IsServiceTrusted(serviceID, accredited)
		if trustedBefore != trustedAfter {
			err = h.AfterUserTrustedServiceUpdated(ctx, delegator, serviceID, trustedAfter)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
