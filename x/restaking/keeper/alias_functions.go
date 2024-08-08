package keeper

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

// IterateAllOperatorParams iterates all operators params and performs the given callback function
func (k *Keeper) IterateAllOperatorParams(ctx sdk.Context, cb func(operatorID uint32, params types.OperatorParams) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.OperatorParamsPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var params types.OperatorParams
		k.cdc.MustUnmarshal(iterator.Value(), &params)

		operatorID := operatorstypes.GetOperatorIDFromBytes(iterator.Key())
		if cb(operatorID, params) {
			break
		}
	}
}

// --------------------------------------------------------------------------------------------------------------------

// IterateAllServiceParams iterates all services params and performs the given callback function
func (k *Keeper) IterateAllServiceParams(ctx sdk.Context, cb func(serviceID uint32, params types.ServiceParams) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.ServiceParamsPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var params types.ServiceParams
		k.cdc.MustUnmarshal(iterator.Value(), &params)

		serviceID := servicestypes.GetServiceIDFromBytes(iterator.Key())
		if cb(serviceID, params) {
			break
		}
	}
}

// --------------------------------------------------------------------------------------------------------------------

// GetAllPoolDelegations returns all the pool delegations
func (k *Keeper) GetAllPoolDelegations(ctx sdk.Context) []types.Delegation {
	store := ctx.KVStore(k.storeKey)
	iterator := store.Iterator(types.PoolDelegationPrefix, storetypes.PrefixEndBytes(types.PoolDelegationPrefix))
	defer iterator.Close()

	var delegations []types.Delegation
	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		delegations = append(delegations, delegation)
	}

	return delegations
}

// GetAllOperatorDelegations returns all the operator delegations
func (k *Keeper) GetAllOperatorDelegations(ctx sdk.Context) []types.Delegation {
	store := ctx.KVStore(k.storeKey)
	iterator := store.Iterator(types.OperatorDelegationPrefix, storetypes.PrefixEndBytes(types.OperatorDelegationPrefix))
	defer iterator.Close()

	var delegations []types.Delegation
	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		delegations = append(delegations, delegation)
	}

	return delegations
}

// GetAllServiceDelegations returns all the service delegations
func (k *Keeper) GetAllServiceDelegations(ctx sdk.Context) []types.Delegation {
	store := ctx.KVStore(k.storeKey)
	iterator := store.Iterator(types.ServiceDelegationPrefix, storetypes.PrefixEndBytes(types.ServiceDelegationPrefix))
	defer iterator.Close()

	var delegations []types.Delegation
	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		delegations = append(delegations, delegation)
	}

	return delegations
}

// GetAllDelegations returns all the delegations
func (k *Keeper) GetAllDelegations(ctx sdk.Context) []types.Delegation {
	var delegations []types.Delegation

	delegations = append(delegations, k.GetAllPoolDelegations(ctx)...)
	delegations = append(delegations, k.GetAllOperatorDelegations(ctx)...)
	delegations = append(delegations, k.GetAllServiceDelegations(ctx)...)

	return delegations
}
