package keeper

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

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
