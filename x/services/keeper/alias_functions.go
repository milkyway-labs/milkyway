package keeper

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/services/types"
)

// createAccountIfNotExists creates an account if it does not exist
func (k *Keeper) createAccountIfNotExists(ctx sdk.Context, address sdk.AccAddress) {
	if !k.accountKeeper.HasAccount(ctx, address) {
		defer telemetry.IncrCounter(1, "new", "account")
		k.accountKeeper.SetAccount(ctx, k.accountKeeper.NewAccountWithAddress(ctx, address))
	}
}

// IterateServices iterates over the services in the store and performs a callback function
func (k *Keeper) IterateServices(ctx sdk.Context, cb func(service types.Service) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.ServicePrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var service types.Service
		k.cdc.MustUnmarshal(iterator.Value(), &service)

		if cb(service) {
			break
		}
	}
}

// GetServices returns the services stored in the KVStore
func (k *Keeper) GetServices(ctx sdk.Context) []types.Service {
	var services []types.Service
	k.IterateServices(ctx, func(service types.Service) (stop bool) {
		services = append(services, service)
		return false
	})
	return services
}

// IsServiceAddress returns true if the provided address is the address
// where the users' asset are kept when they restake toward a service.
func (k *Keeper) IsServiceAddress(ctx sdk.Context, address string) (bool, error) {
	return k.serviceAddressSet.Has(ctx, address)
}

// GetAllServicesParams returns the parameters that have been configured for all
// services.
func (k *Keeper) GetAllServicesParams(ctx sdk.Context) ([]types.ServiceParamsRecord, error) {
	var records []types.ServiceParamsRecord
	err := k.serviceParams.Walk(ctx, nil, func(serviceID uint32, params types.ServiceParams) (bool, error) {
		records = append(records, types.NewServiceParamsRecord(serviceID, params))
		return false, nil
	})

	return records, err
}
