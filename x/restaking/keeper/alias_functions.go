package keeper

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

// IterateAllPoolDelegations iterates over all the pool delegations and performs a callback function
func (k *Keeper) IterateAllPoolDelegations(ctx sdk.Context, cb func(delegation types.PoolDelegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := store.Iterator(types.PoolDelegationPrefix, storetypes.PrefixEndBytes(types.PoolDelegationPrefix))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalPoolDelegation(k.cdc, iterator.Value())
		if cb(delegation) {
			break
		}
	}
}

// GetAllPoolDelegations returns all the pool delegations
func (k *Keeper) GetAllPoolDelegations(ctx sdk.Context) []types.PoolDelegation {
	var delegations []types.PoolDelegation
	k.IterateAllPoolDelegations(ctx, func(delegation types.PoolDelegation) bool {
		delegations = append(delegations, delegation)
		return false
	})
	return delegations
}

// GetAllDelegatorPoolDelegations returns all the pool delegations of a given delegator
func (k *Keeper) GetAllDelegatorPoolDelegations(ctx sdk.Context, delegator string) []types.PoolDelegation {
	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.UserPoolDelegationsStorePrefix(delegator)

	iterator := store.Iterator(delegatorPrefixKey, storetypes.PrefixEndBytes(delegatorPrefixKey)) // Smallest to largest
	defer iterator.Close()

	var delegations []types.PoolDelegation
	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalPoolDelegation(k.cdc, iterator.Value())
		delegations = append(delegations, delegation)

	}

	return delegations
}

// GetPoolDelegations returns all the delegations to a given pool
func (k *Keeper) GetPoolDelegations(ctx sdk.Context, poolID uint32) ([]types.PoolDelegation, error) {
	store := ctx.KVStore(k.storeKey)
	prefix := types.DelegationsByPoolIDStorePrefix(poolID)
	iterator := store.Iterator(prefix, storetypes.PrefixEndBytes(prefix))
	defer iterator.Close()

	var delegations []types.PoolDelegation
	for ; iterator.Valid(); iterator.Next() {
		_, delegatorAddress, err := types.ParseDelegationsByPoolIDKey(iterator.Key())
		if err != nil {
			return nil, err
		}

		delegation, found := k.GetPoolDelegation(ctx, poolID, delegatorAddress)
		if !found {
			return nil, fmt.Errorf("delegation not found for pool %d and delegator %s", poolID, delegatorAddress)
		}

		delegations = append(delegations, delegation)
	}

	return delegations, nil
}
