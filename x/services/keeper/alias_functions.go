package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v9/x/services/types"
)

// createAccountIfNotExists creates an account if it does not exist
func (k *Keeper) createAccountIfNotExists(ctx context.Context, address sdk.AccAddress) {
	if !k.accountKeeper.HasAccount(ctx, address) {
		defer telemetry.IncrCounter(1, "new", "account")
		k.accountKeeper.SetAccount(ctx, k.accountKeeper.NewAccountWithAddress(ctx, address))
	}
}

// IterateServices iterates over the services in the store and performs a callback function
func (k *Keeper) IterateServices(ctx context.Context, cb func(service types.Service) (stop bool, err error)) error {
	err := k.services.Walk(ctx, nil, func(_ uint32, service types.Service) (stop bool, err error) {
		return cb(service)
	})
	return err
}

// GetServices returns the services stored in the KVStore
func (k *Keeper) GetServices(ctx context.Context) ([]types.Service, error) {
	var services []types.Service
	err := k.IterateServices(ctx, func(service types.Service) (stop bool, err error) {
		services = append(services, service)
		return false, nil
	})
	return services, err
}

// IsServiceAddress returns true if the provided address is the address
// where the users' asset are kept when they restake toward a service.
func (k *Keeper) IsServiceAddress(ctx context.Context, address string) (bool, error) {
	return k.serviceAddressSet.Has(ctx, address)
}

// GetAllServicesParams returns the parameters that have been configured for all
// services.
func (k *Keeper) GetAllServicesParams(ctx context.Context) ([]types.ServiceParamsRecord, error) {
	var records []types.ServiceParamsRecord
	err := k.serviceParams.Walk(ctx, nil, func(serviceID uint32, params types.ServiceParams) (bool, error) {
		records = append(records, types.NewServiceParamsRecord(serviceID, params))
		return false, nil
	})

	return records, err
}
