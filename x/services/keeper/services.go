package keeper

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/services/types"
)

// SetNextServiceID sets the next service ID to be used when registering a new Service
func (k *Keeper) SetNextServiceID(ctx sdk.Context, avsID uint32) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.NextServiceIDKey, types.GetServiceIDBytes(avsID))
}

// HasNextServiceID checks if the next service ID is set
func (k *Keeper) HasNextServiceID(ctx sdk.Context) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.NextServiceIDKey)
}

// GetNextServiceID returns the next service ID to be used when registering a new Service
func (k *Keeper) GetNextServiceID(ctx sdk.Context) (avsID uint32, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextServiceIDKey)
	if bz == nil {
		return 0, errors.Wrapf(types.ErrInvalidGenesis, "initial service id not set")
	}

	avsID = types.GetServiceIDFromBytes(bz)
	return avsID, nil
}

// --------------------------------------------------------------------------------------------------------------------

func (k *Keeper) storeService(ctx sdk.Context, service types.Service) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ServiceStoreKey(service.ID), k.cdc.MustMarshal(&service))
}

// CreateService creates a new Service and stores it in the KVStore
func (k *Keeper) CreateService(ctx sdk.Context, service types.Service) error {
	// Store the service
	k.storeService(ctx, service)

	// Charge for the creation
	registrationFees := k.GetParams(ctx).ServiceRegistrationFee
	if registrationFees != nil && registrationFees.IsZero() {
		userAddress, err := sdk.AccAddressFromBech32(service.Admin)
		if err != nil {
			return err
		}

		err = k.poolKeeper.FundCommunityPool(ctx, registrationFees, userAddress)
		if err != nil {
			return err
		}
	}

	// Log and call the hooks
	k.Logger(ctx).Debug("created service", "id", service.ID)
	k.AfterServiceCreated(ctx, service.ID)

	return nil
}

// UpdateService updates an existing Service
func (k *Keeper) UpdateService(ctx sdk.Context, service types.Service) error {
	previous, existed := k.GetService(ctx, service.ID)
	if !existed {
		return types.ErrServiceNotFound
	}

	// Update the service
	k.storeService(ctx, service)
	k.Logger(ctx).Debug("updated service", "id", service.ID)

	// Call the hook based on the Service status change
	switch {
	case previous.Status == types.SERVICE_STATUS_CREATED && service.Status == types.SERVICE_STATUS_ACTIVE:
		k.AfterServiceActivated(ctx, service.ID)
	case previous.Status == types.SERVICE_STATUS_ACTIVE && service.Status == types.SERVICE_STATUS_INACTIVE:
		k.AfterServiceDeactivated(ctx, service.ID)
	}

	return nil
}

// GetService returns an Service from the KVStore
func (k *Keeper) GetService(ctx sdk.Context, serviceID uint32) (service types.Service, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.ServiceStoreKey(serviceID))
	if bz == nil {
		return service, false
	}

	k.cdc.MustUnmarshal(bz, &service)
	return service, true
}
