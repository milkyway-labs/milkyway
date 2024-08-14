package keeper

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/utils"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

func (k *Keeper) IterateAllOperatorParams(
	ctx sdk.Context, cb func(operatorID uint32, params types.OperatorParams) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.OperatorParamsPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var params types.OperatorParams
		k.cdc.MustUnmarshal(iterator.Value(), &params)

		operatorID := utils.BigEndianToUint32(iterator.Key())
		if cb(operatorID, params) {
			break
		}
	}
}

func (k *Keeper) IterateAllServiceParams(ctx sdk.Context, cb func(serviceID uint32, params types.ServiceParams) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.ServiceParamsPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var params types.ServiceParams
		k.cdc.MustUnmarshal(iterator.Value(), &params)

		serviceID := utils.BigEndianToUint32(iterator.Key())
		if cb(serviceID, params) {
			break
		}
	}
}

// --------------------------------------------------------------------------------------------------------------------

func (k *Keeper) IterateUserPoolDelegations(ctx sdk.Context, userAddress string, cb func(del types.Delegation) (stop bool, err error)) error {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.UserPoolDelegationsStorePrefix(userAddress))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		stop, err := cb(delegation)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return nil
}

func (k *Keeper) IterateAllPoolDelegations(ctx sdk.Context, cb func(del types.Delegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := store.Iterator(types.PoolDelegationPrefix, storetypes.PrefixEndBytes(types.PoolDelegationPrefix))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		if cb(delegation) {
			break
		}
	}
}

// GetAllPoolDelegations returns all the pool delegations
func (k *Keeper) GetAllPoolDelegations(ctx sdk.Context) []types.Delegation {
	var delegations []types.Delegation
	k.IterateAllPoolDelegations(ctx, func(del types.Delegation) (stop bool) {
		delegations = append(delegations, del)
		return false
	})

	return delegations
}

func (k *Keeper) IterateUserOperatorDelegations(ctx sdk.Context, userAddress string, cb func(del types.Delegation) (stop bool, err error)) error {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.UserOperatorDelegationsStorePrefix(userAddress))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		stop, err := cb(delegation)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return nil
}

func (k *Keeper) IterateAllOperatorDelegations(ctx sdk.Context, cb func(del types.Delegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := store.Iterator(types.OperatorDelegationPrefix, storetypes.PrefixEndBytes(types.OperatorDelegationPrefix))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		if cb(delegation) {
			break
		}
	}
}

// GetAllOperatorDelegations returns all the operator delegations
func (k *Keeper) GetAllOperatorDelegations(ctx sdk.Context) []types.Delegation {
	var delegations []types.Delegation
	k.IterateAllOperatorDelegations(ctx, func(del types.Delegation) (stop bool) {
		delegations = append(delegations, del)
		return false
	})

	return delegations
}

func (k *Keeper) IterateUserServiceDelegations(ctx sdk.Context, userAddress string, cb func(del types.Delegation) (stop bool, err error)) error {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.UserServiceDelegationsStorePrefix(userAddress))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		stop, err := cb(delegation)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return nil
}

func (k *Keeper) IterateAllServiceDelegations(ctx sdk.Context, cb func(del types.Delegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := store.Iterator(types.ServiceDelegationPrefix, storetypes.PrefixEndBytes(types.ServiceDelegationPrefix))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		if cb(delegation) {
			break
		}
	}
}

// GetAllServiceDelegations returns all the service delegations
func (k *Keeper) GetAllServiceDelegations(ctx sdk.Context) []types.Delegation {
	var delegations []types.Delegation
	k.IterateAllServiceDelegations(ctx, func(del types.Delegation) (stop bool) {
		delegations = append(delegations, del)
		return false
	})

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
