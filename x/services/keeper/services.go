package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/milkyway-labs/milkyway/v5/x/services/types"
)

// SetNextServiceID sets the next service ID to be used when registering a new Service
func (k *Keeper) SetNextServiceID(ctx context.Context, serviceID uint32) error {
	return k.nextServiceID.Set(ctx, uint64(serviceID))
}

// GetNextServiceID returns the next service ID to be used when registering a new Service
func (k *Keeper) GetNextServiceID(ctx context.Context) (serviceID uint32, err error) {
	nextServiceID, err := k.nextServiceID.Next(ctx)
	if err != nil {
		return 0, err
	}

	// If the next service ID is 0, we need to increment it
	if nextServiceID == 0 {
		return k.GetNextServiceID(ctx)
	}

	return uint32(nextServiceID), nil
}

// --------------------------------------------------------------------------------------------------------------------

// SaveService stores a Service in the KVStore
func (k *Keeper) SaveService(ctx context.Context, service types.Service) error {
	err := k.services.Set(ctx, service.ID, service)
	if err != nil {
		return err
	}

	return k.serviceAddressSet.Set(ctx, service.Address)
}

// CreateService creates a new Service and stores it in the KVStore
func (k *Keeper) CreateService(ctx context.Context, service types.Service) error {
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
	return k.AfterServiceCreated(ctx, service.ID)
}

// ActivateService activates the service with the given ID
func (k *Keeper) ActivateService(ctx context.Context, serviceID uint32) error {
	service, err := k.GetService(ctx, serviceID)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return types.ErrServiceNotFound
		}
		return err
	}

	// Check if the service is already active
	if service.Status == types.SERVICE_STATUS_ACTIVE {
		return types.ErrServiceAlreadyActive
	}

	// Update the status
	service.Status = types.SERVICE_STATUS_ACTIVE

	// Update the service
	err = k.SaveService(ctx, service)
	if err != nil {
		return err
	}

	// Call the hook
	return k.AfterServiceActivated(ctx, serviceID)
}

// DeactivateService deactivates the service with the given ID
func (k *Keeper) DeactivateService(ctx context.Context, serviceID uint32) error {
	service, err := k.GetService(ctx, serviceID)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return types.ErrServiceNotFound
		}
		return err
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
	return k.AfterServiceDeactivated(ctx, service.ID)
}

// DeleteService deletes the service with the given ID
func (k *Keeper) DeleteService(ctx context.Context, serviceID uint32) error {
	service, err := k.GetService(ctx, serviceID)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return types.ErrServiceNotFound
		}
		return err
	}

	// Make sure the service is not active
	if service.Status == types.SERVICE_STATUS_ACTIVE {
		return types.ErrServiceIsActive
	}

	// Call the hook
	err = k.BeforeServiceDeleted(ctx, service.ID)
	if err != nil {
		return err
	}

	// Remove the service from the store
	err = k.services.Remove(ctx, serviceID)
	if err != nil {
		return err
	}

	return k.serviceAddressSet.Remove(ctx, service.Address)
}

// SetServiceAccredited sets the accreditation of the service with the given ID
func (k *Keeper) SetServiceAccredited(ctx context.Context, serviceID uint32, accredited bool) error {
	// Check if the service exists
	service, err := k.GetService(ctx, serviceID)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return types.ErrServiceNotFound
		}
		return err
	}

	// Skip any operation if the service accreditation status does not change
	if service.Accredited == accredited {
		return nil
	}

	// Update the service accreditation status
	service.Accredited = accredited
	err = k.SaveService(ctx, service)
	if err != nil {
		return err
	}

	// Call the hook
	return k.AfterServiceAccreditationModified(ctx, service.ID)
}

// HasService checks if a Service with the given ID exists in the KVStore
func (k *Keeper) HasService(ctx context.Context, serviceID uint32) (bool, error) {
	return k.services.Has(ctx, serviceID)
}

// GetService returns an Service from the KVStore
func (k *Keeper) GetService(ctx context.Context, serviceID uint32) (service types.Service, err error) {
	return k.services.Get(ctx, serviceID)
}

// GetServiceParams returns the params for the service with the given ID
func (k *Keeper) GetServiceParams(ctx context.Context, serviceID uint32) (types.ServiceParams, error) {
	params, err := k.serviceParams.Get(ctx, serviceID)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return types.DefaultServiceParams(), nil
		}
		return types.ServiceParams{}, err
	}

	return params, nil
}

// SetServiceParams sets the params for the service with the given ID
func (k *Keeper) SetServiceParams(ctx context.Context, serviceID uint32, params types.ServiceParams) error {
	err := params.Validate()
	if err != nil {
		return err
	}

	return k.serviceParams.Set(ctx, serviceID, params)
}
