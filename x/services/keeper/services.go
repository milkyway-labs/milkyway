package keeper

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/milkyway-labs/milkyway/x/services/types"
)

// SetNextServiceID sets the next service ID to be used when registering a new Service
func (k *Keeper) SetNextServiceID(ctx sdk.Context, serviceID uint32) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.NextServiceIDKey, types.GetServiceIDBytes(serviceID))
}

// GetNextServiceID returns the next service ID to be used when registering a new Service
func (k *Keeper) GetNextServiceID(ctx sdk.Context) (serviceID uint32, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextServiceIDKey)
	if bz == nil {
		return 0, errors.Wrapf(types.ErrInvalidGenesis, "initial service id not set")
	}

	serviceID = types.GetServiceIDFromBytes(bz)
	return serviceID, nil
}

// --------------------------------------------------------------------------------------------------------------------

// SaveService stores a Service in the KVStore
func (k *Keeper) SaveService(ctx sdk.Context, service types.Service) error {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ServiceStoreKey(service.ID), k.cdc.MustMarshal(&service))
	return k.serviceAddressSet.Set(ctx, service.Address)
}

// CreateService creates a new Service and stores it in the KVStore
func (k *Keeper) CreateService(ctx sdk.Context, service types.Service) error {
	// Charge for the creation
	registrationFees := k.GetParams(ctx).ServiceRegistrationFee
	if !registrationFees.IsZero() {
		userAddress, err := sdk.AccAddressFromBech32(service.Admin)
		if err != nil {
			return err
		}

		err = k.poolKeeper.FundCommunityPool(ctx, registrationFees, userAddress)
		if err != nil {
			return err
		}
	}

	// Create the service account
	serviceAddress, err := sdk.AccAddressFromBech32(service.Address)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid service address: %s", service.Address)
	}
	k.createAccountIfNotExists(ctx, serviceAddress)

	// Store the service
	if err := k.SaveService(ctx, service); err != nil {
		return err
	}

	// Log and call the hooks
	k.Logger(ctx).Debug("created service", "id", service.ID)
	k.AfterServiceCreated(ctx, service.ID)

	return nil
}

// ActivateService activates the service with the given ID
func (k *Keeper) ActivateService(ctx sdk.Context, serviceID uint32) error {
	service, found := k.GetService(ctx, serviceID)
	if !found {
		return types.ErrServiceNotFound
	}

	// Check if the service is already active
	if service.Status == types.SERVICE_STATUS_ACTIVE {
		return types.ErrServiceAlreadyActive
	}

	service.Status = types.SERVICE_STATUS_ACTIVE
	if err := k.SaveService(ctx, service); err != nil {
		return err
	}

	// Call the hook
	k.AfterServiceActivated(ctx, serviceID)

	return nil
}

// DeactivateService deactivates the service with the given ID
func (k *Keeper) DeactivateService(ctx sdk.Context, serviceID uint32) error {
	service, existed := k.GetService(ctx, serviceID)
	if !existed {
		return types.ErrServiceNotFound
	}

	// Make sure the service is active
	if service.Status != types.SERVICE_STATUS_ACTIVE {
		return types.ErrServiceNotActive
	}

	// Update the status
	service.Status = types.SERVICE_STATUS_INACTIVE

	// Update the service
	if err := k.SaveService(ctx, service); err != nil {
		return err
	}

	// Call the hook
	k.AfterServiceDeactivated(ctx, service.ID)

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
