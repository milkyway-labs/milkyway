package keeper

import (
	"context"
	"fmt"
	"sort"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v7/utils"
	operatorstypes "github.com/milkyway-labs/milkyway/v7/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/v7/x/pools/types"
	"github.com/milkyway-labs/milkyway/v7/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v7/x/services/types"
)

// --------------------------------------------------------------------------------------------------------------------

// IterateServiceValidatingOperators iterates over all the operators that have
// joined the given service, performing the given action. If the action returns
// true, the iteration will stop.
func (k *Keeper) IterateServiceValidatingOperators(ctx context.Context, serviceID uint32, action func(operatorID uint32) (stop bool, err error)) error {
	return k.operatorJoinedServices.Indexes.Service.Walk(ctx, collections.NewPrefixedPairRange[uint32, uint32](serviceID), func(_, operatorID uint32) (stop bool, err error) {
		return action(operatorID)
	})
}

// IterateAllOperatorsJoinedServices iterates over all the operators and their joined services,
// performing the given action. If the action returns true, the iteration will stop.
func (k *Keeper) IterateAllOperatorsJoinedServices(ctx context.Context, action func(operatorID uint32, serviceID uint32) (stop bool, err error)) error {
	err := k.operatorJoinedServices.Walk(ctx, nil, func(key collections.Pair[uint32, uint32], _ collections.NoValue) (stop bool, err error) {
		operatorID := key.K1()
		serviceID := key.K2()
		return action(operatorID, serviceID)
	})
	return err
}

// GetAllOperatorsJoinedServices returns all services that each operator has joined
func (k *Keeper) GetAllOperatorsJoinedServices(ctx context.Context) ([]types.OperatorJoinedServices, error) {
	items := make(map[uint32]types.OperatorJoinedServices)
	err := k.IterateAllOperatorsJoinedServices(ctx, func(operatorID uint32, serviceID uint32) (stop bool, err error) {
		joinedServicesRecord, ok := items[operatorID]
		if !ok {
			joinedServicesRecord = types.NewOperatorJoinedServices(operatorID, nil)
		}
		joinedServicesRecord.ServiceIDs = append(joinedServicesRecord.ServiceIDs, serviceID)
		items[operatorID] = joinedServicesRecord
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, nil
	}

	// Convert back to list
	itemsList := make([]types.OperatorJoinedServices, 0, len(items))
	for _, v := range items {
		itemsList = append(itemsList, v)
	}
	// Ensure that the items always maintain the same order,
	// as iterating over the map may result in different item orders.
	sort.Slice(itemsList, func(i, j int) bool {
		return itemsList[i].OperatorID < itemsList[j].OperatorID
	})
	return itemsList, nil
}

// --------------------------------------------------------------------------------------------------------------------

// IterateAllServicesAllowedOperators iterates over all the services and their allowed operators,
// performing the given action. If the action returns true, the iteration will stop.
func (k *Keeper) IterateAllServicesAllowedOperators(ctx context.Context, action func(serviceID uint32, operatorID uint32) (stop bool, err error)) error {
	err := k.serviceOperatorsAllowList.Walk(ctx, nil, func(key collections.Pair[uint32, uint32]) (stop bool, err error) {
		serviceID := key.K1()
		operatorID := key.K2()
		return action(serviceID, operatorID)
	})
	return err
}

// GetAllServicesAllowedOperators returns all the operators that are allowed to secure a service for all the services
func (k *Keeper) GetAllServicesAllowedOperators(ctx context.Context) ([]types.ServiceAllowedOperators, error) {
	items := make(map[uint32]types.ServiceAllowedOperators)
	err := k.IterateAllServicesAllowedOperators(ctx, func(serviceID uint32, operatorID uint32) (stop bool, err error) {
		allowedOperators, ok := items[serviceID]
		if !ok {
			allowedOperators = types.NewServiceAllowedOperators(serviceID, nil)
		}
		allowedOperators.OperatorIDs = append(allowedOperators.OperatorIDs, operatorID)
		items[serviceID] = allowedOperators

		return false, nil
	})
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, nil
	}

	// Convert back to list
	itemsList := make([]types.ServiceAllowedOperators, 0, len(items))
	for _, v := range items {
		itemsList = append(itemsList, v)
	}

	// Ensure that the items always maintain the same order,
	// as iterating over the map may result in different item orders.
	sort.Slice(itemsList, func(i, j int) bool {
		return itemsList[i].ServiceID < itemsList[j].ServiceID
	})
	return itemsList, nil
}

// GetAllServicesSecuringPools returns all the pools from which the services are allowed to borrow security
func (k *Keeper) GetAllServicesSecuringPools(ctx context.Context) ([]types.ServiceSecuringPools, error) {
	items := make(map[uint32]types.ServiceSecuringPools)
	err := k.serviceSecuringPools.Walk(ctx, nil, func(key collections.Pair[uint32, uint32]) (stop bool, err error) {
		serviceID := key.K1()
		poolID := key.K2()
		securingPools, ok := items[serviceID]
		if !ok {
			securingPools = types.NewServiceSecuringPools(serviceID, nil)
		}
		securingPools.PoolIDs = append(securingPools.PoolIDs, poolID)
		items[serviceID] = securingPools
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}

	// Convert back to list
	itemsList := make([]types.ServiceSecuringPools, 0, len(items))
	for _, v := range items {
		itemsList = append(itemsList, v)
	}

	// Ensure that the items always maintain the same order,
	// as iterating over the map may result in different item orders.
	sort.Slice(itemsList, func(i, j int) bool {
		return itemsList[i].ServiceID < itemsList[j].ServiceID
	})
	return itemsList, nil
}

// --------------------------------------------------------------------------------------------------------------------
// --- Delegation operations
// --------------------------------------------------------------------------------------------------------------------

// SetDelegation stores the given delegation in the store
func (k *Keeper) SetDelegation(ctx context.Context, delegation types.Delegation) error {
	store := k.storeService.OpenKVStore(ctx)

	// Get the keys builders
	getDelegationKey, getDelegationByTargetID, err := types.GetDelegationKeyBuilders(delegation)
	if err != nil {
		return err
	}

	// Marshal and store the delegation
	delegationBz := types.MustMarshalDelegation(k.cdc, delegation)
	err = store.Set(getDelegationKey(delegation.UserAddress, delegation.TargetID), delegationBz)
	if err != nil {
		return err
	}

	// Store the key in the delegation by target ID index (used for reversed lookup)
	err = store.Set(getDelegationByTargetID(delegation.TargetID, delegation.UserAddress), []byte{})
	if err != nil {
		return err
	}

	return nil
}

// GetDelegationForTarget returns the delegation for the given delegator and target.
func (k *Keeper) GetDelegationForTarget(ctx context.Context, target types.DelegationTarget, delegator string) (types.Delegation, bool, error) {
	switch target.(type) {
	case poolstypes.Pool:
		return k.GetPoolDelegation(ctx, target.GetID(), delegator)
	case operatorstypes.Operator:
		return k.GetOperatorDelegation(ctx, target.GetID(), delegator)
	case servicestypes.Service:
		return k.GetServiceDelegation(ctx, target.GetID(), delegator)
	default:
		return types.Delegation{}, false, fmt.Errorf("invalid target type %T", target)
	}
}

// GetDelegationTargetFromDelegation returns the target of the given delegation.
func (k *Keeper) GetDelegationTargetFromDelegation(ctx context.Context, delegation types.Delegation) (types.DelegationTarget, error) {
	switch delegation.Type {
	case types.DELEGATION_TYPE_POOL:
		return k.poolsKeeper.GetPool(ctx, delegation.TargetID)
	case types.DELEGATION_TYPE_SERVICE:
		return k.servicesKeeper.GetService(ctx, delegation.TargetID)
	case types.DELEGATION_TYPE_OPERATOR:
		return k.operatorsKeeper.GetOperator(ctx, delegation.TargetID)
	default:
		return nil, nil
	}
}

// RemoveDelegation removes the given delegation from the store
func (k *Keeper) RemoveDelegation(ctx context.Context, delegation types.Delegation) error {
	switch delegation.Type {
	case types.DELEGATION_TYPE_POOL:
		return k.RemovePoolDelegation(ctx, delegation)
	case types.DELEGATION_TYPE_OPERATOR:
		return k.RemoveOperatorDelegation(ctx, delegation)
	case types.DELEGATION_TYPE_SERVICE:
		return k.RemoveServiceDelegation(ctx, delegation)
	default:
		return errors.Wrapf(types.ErrInvalidDelegationType, "invalid delegation type %v", delegation.Type)
	}
}

// --------------------------------------------------------------------------------------------------------------------
// --- Delegations iterations operations
// --------------------------------------------------------------------------------------------------------------------

// IterateUserPoolDelegations iterates all the pool delegations of a user and performs the given callback function
func (k *Keeper) IterateUserPoolDelegations(ctx context.Context, userAddress string, cb func(del types.Delegation) (stop bool, err error)) error {
	store := k.storeService.OpenKVStore(ctx)

	iterator := storetypes.KVStorePrefixIterator(runtime.KVStoreAdapter(store), types.UserPoolDelegationsStorePrefix(userAddress))
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

// IterateAllPoolDelegations iterates all the pool delegations and performs the given callback function
func (k *Keeper) IterateAllPoolDelegations(ctx context.Context, cb func(del types.Delegation) (stop bool, err error)) error {
	store := k.storeService.OpenKVStore(ctx)

	iterator, err := store.Iterator(types.PoolDelegationPrefix, storetypes.PrefixEndBytes(types.PoolDelegationPrefix))
	if err != nil {
		return err
	}
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

// GetAllPoolDelegations returns all the pool delegations
func (k *Keeper) GetAllPoolDelegations(ctx context.Context) ([]types.Delegation, error) {
	var delegations []types.Delegation
	err := k.IterateAllPoolDelegations(ctx, func(del types.Delegation) (stop bool, err error) {
		delegations = append(delegations, del)
		return false, nil
	})
	return delegations, err
}

// IterateUserOperatorDelegations iterates all the operator delegations of a user and performs the given callback function
func (k *Keeper) IterateUserOperatorDelegations(ctx context.Context, userAddress string, cb func(del types.Delegation) (stop bool, err error)) error {
	store := k.storeService.OpenKVStore(ctx)

	iterator := storetypes.KVStorePrefixIterator(runtime.KVStoreAdapter(store), types.UserOperatorDelegationsStorePrefix(userAddress))
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

// IterateAllOperatorDelegations iterates all the operator delegations and performs the given callback function
func (k *Keeper) IterateAllOperatorDelegations(ctx context.Context, cb func(del types.Delegation) (stop bool, err error)) error {
	store := k.storeService.OpenKVStore(ctx)

	iterator, err := store.Iterator(types.OperatorDelegationPrefix, storetypes.PrefixEndBytes(types.OperatorDelegationPrefix))
	if err != nil {
		return err
	}
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

// GetAllOperatorDelegations returns all the operator delegations
func (k *Keeper) GetAllOperatorDelegations(ctx context.Context) ([]types.Delegation, error) {
	var delegations []types.Delegation
	err := k.IterateAllOperatorDelegations(ctx, func(del types.Delegation) (stop bool, err error) {
		delegations = append(delegations, del)
		return false, nil
	})
	return delegations, err
}

// IterateUserServiceDelegations iterates all the service delegations of a user and performs the given callback function
func (k *Keeper) IterateUserServiceDelegations(ctx context.Context, userAddress string, cb func(del types.Delegation) (stop bool, err error)) error {
	store := k.storeService.OpenKVStore(ctx)

	iterator := storetypes.KVStorePrefixIterator(runtime.KVStoreAdapter(store), types.UserServiceDelegationsStorePrefix(userAddress))
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

// IterateAllServiceDelegations iterates all the service delegations and performs the given callback function
func (k *Keeper) IterateAllServiceDelegations(ctx context.Context, cb func(del types.Delegation) (stop bool, err error)) error {
	store := k.storeService.OpenKVStore(ctx)

	iterator, err := store.Iterator(types.ServiceDelegationPrefix, storetypes.PrefixEndBytes(types.ServiceDelegationPrefix))
	if err != nil {
		return err
	}
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

// IterateServiceDelegations iterates all the delegations of a service and
// performs the given callback function
func (k *Keeper) IterateServiceDelegations(ctx context.Context, serviceID uint32, cb func(del types.Delegation) (stop bool, err error)) error {
	store := k.storeService.OpenKVStore(ctx)

	iterator := storetypes.KVStorePrefixIterator(runtime.KVStoreAdapter(store), types.DelegationsByServiceIDStorePrefix(serviceID))
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

// GetAllServiceDelegations returns all the service delegations
func (k *Keeper) GetAllServiceDelegations(ctx context.Context) ([]types.Delegation, error) {
	var delegations []types.Delegation
	err := k.IterateAllServiceDelegations(ctx, func(del types.Delegation) (stop bool, err error) {
		delegations = append(delegations, del)
		return false, nil
	})
	return delegations, err
}

// IterateUserDelegations iterates over the user's delegations.
// The delegations will be iterated in the following order:
// 1. Pool delegations
// 2. Service delegations
// 3. Operator delegations
func (k *Keeper) IterateUserDelegations(ctx context.Context, userAddress string, cb func(del types.Delegation) (stop bool, err error)) error {
	store := k.storeService.OpenKVStore(ctx)

	poolIterator := storetypes.KVStorePrefixIterator(runtime.KVStoreAdapter(store), types.UserPoolDelegationsStorePrefix(userAddress))
	defer poolIterator.Close()

	// Iterate pools delegations
	for ; poolIterator.Valid(); poolIterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, poolIterator.Value())
		stop, err := cb(delegation)
		if err != nil {
			return err
		}
		if stop {
			return nil
		}
	}

	// Iterate service delegations
	servicesIterator := storetypes.KVStorePrefixIterator(runtime.KVStoreAdapter(store), types.UserServiceDelegationsStorePrefix(userAddress))
	defer servicesIterator.Close()
	for ; servicesIterator.Valid(); servicesIterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, servicesIterator.Value())
		stop, err := cb(delegation)
		if err != nil {
			return err
		}
		if stop {
			return nil
		}
	}

	// Iterate operators delegations
	operatorsIterator := storetypes.KVStorePrefixIterator(runtime.KVStoreAdapter(store), types.UserOperatorDelegationsStorePrefix(userAddress))
	defer operatorsIterator.Close()
	for ; operatorsIterator.Valid(); operatorsIterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, operatorsIterator.Value())
		stop, err := cb(delegation)
		if err != nil {
			return err
		}
		if stop {
			return nil
		}
	}

	return nil
}

// GetAllDelegations returns all the delegations
func (k *Keeper) GetAllDelegations(ctx context.Context) ([]types.Delegation, error) {

	var delegations []types.Delegation

	poolsDelegations, err := k.GetAllPoolDelegations(ctx)
	if err != nil {
		return nil, err
	}
	delegations = append(delegations, poolsDelegations...)

	operatorDelegations, err := k.GetAllOperatorDelegations(ctx)
	if err != nil {
		return nil, err
	}
	delegations = append(delegations, operatorDelegations...)

	serviceDelegations, err := k.GetAllServiceDelegations(ctx)
	if err != nil {
		return nil, err
	}
	delegations = append(delegations, serviceDelegations...)

	return delegations, nil
}

// GetAllUserRestakedCoins returns all the user's restaked coins
func (k *Keeper) GetAllUserRestakedCoins(ctx context.Context, userAddress string) (sdk.DecCoins, error) {
	totalDelegatedCoins := sdk.NewDecCoins()
	err := k.IterateUserDelegations(ctx, userAddress, func(d types.Delegation) (bool, error) {
		target, err := k.GetDelegationTargetFromDelegation(ctx, d)
		if err != nil {
			return true, err
		}

		totalDelegatedCoins = totalDelegatedCoins.Add(target.TokensFromShares(d.Shares)...)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return totalDelegatedCoins, nil
}

// PerformDelegation performs a delegation of the given amount from the delegator to the receiver.
// It sends the coins to the receiver address and updates the delegation object and returns the new
// shares of the delegation.
// NOTE: This is done so that if we implement other delegation types in the future we can have a single
// function that performs common operations for all of them.
func (k *Keeper) PerformDelegation(ctx context.Context, data types.DelegationData) (sdk.DecCoins, error) {
	// Get the data
	receiver := data.Target
	delegator := data.Delegator
	hooks := data.Hooks

	// In some situations, the exchange rate becomes invalid, e.g. if
	// the receives loses all tokens due to slashing. In this case,
	// make all future delegations invalid.
	if receiver.InvalidExRate() {
		return nil, types.ErrDelegatorShareExRateInvalid
	}

	// Check if the restake operation is allowed
	err := k.ValidateRestake(ctx, delegator, data.Amount, data.Target)
	if err != nil {
		return nil, err
	}

	// Get or create the delegation object and call the appropriate hook if present
	delegation, found, err := k.GetDelegationForTarget(ctx, receiver, delegator)
	if err != nil {
		return nil, err
	}

	if found {
		// Delegation was found
		err = hooks.BeforeDelegationSharesModified(ctx, receiver.GetID(), delegator)
		if err != nil {
			return nil, err
		}
	} else {
		// Delegation was not found
		delegation = data.BuildDelegation(receiver.GetID(), delegator, sdk.NewDecCoins())
		err = hooks.BeforeDelegationCreated(ctx, receiver.GetID(), delegator)
		if err != nil {
			return nil, err
		}
	}

	// Convert the addresses to sdk.AccAddress
	delegatorAddress, err := k.accountKeeper.AddressCodec().StringToBytes(delegator)
	if err != nil {
		return nil, err
	}
	receiverAddress, err := k.accountKeeper.AddressCodec().StringToBytes(receiver.GetAddress())
	if err != nil {
		return nil, err
	}

	// Send the coins to the receiver address
	err = k.bankKeeper.SendCoins(ctx, delegatorAddress, receiverAddress, data.Amount)
	if err != nil {
		return nil, err
	}

	// Update the delegation
	newShares, err := data.UpdateDelegation(ctx, delegation)
	if err != nil {
		return nil, err
	}

	// Call the after-modification hook
	err = hooks.AfterDelegationModified(ctx, receiver.GetID(), delegator)
	if err != nil {
		return nil, err
	}

	return newShares, nil
}

// --------------------------------------------------------------------------------------------------------------------
// --- Unbonding operations
// --------------------------------------------------------------------------------------------------------------------

// getUnbondingDelegationTarget returns the target of the given unbonding delegation
func (k *Keeper) getUnbondingDelegationTarget(ctx context.Context, ubd types.UnbondingDelegation) (types.DelegationTarget, error) {
	switch ubd.Type {
	case types.DELEGATION_TYPE_POOL:
		pool, err := k.poolsKeeper.GetPool(ctx, ubd.TargetID)
		if err != nil {
			return nil, err
		}
		return pool, nil

	case types.DELEGATION_TYPE_OPERATOR:
		operator, err := k.operatorsKeeper.GetOperator(ctx, ubd.TargetID)
		if err != nil {
			return nil, err
		}
		return operator, nil

	case types.DELEGATION_TYPE_SERVICE:
		service, err := k.servicesKeeper.GetService(ctx, ubd.TargetID)
		if err != nil {
			return nil, err
		}
		return service, nil

	default:
		return nil, errors.Wrapf(types.ErrInvalidDelegationType, "invalid delegation type %v", ubd.Type)
	}
}

// getUnbondingDelegationKeyBuilder returns the key builder for the given unbonding delegation
func (k *Keeper) getUnbondingDelegationKeyBuilder(ud types.UnbondingDelegation) (types.UnbondingDelegationKeyBuilder, error) {
	switch ud.Type {
	case types.DELEGATION_TYPE_POOL:
		return types.UserPoolUnbondingDelegationKey, nil

	case types.DELEGATION_TYPE_OPERATOR:
		return types.UserOperatorUnbondingDelegationKey, nil

	case types.DELEGATION_TYPE_SERVICE:
		return types.UserServiceUnbondingDelegationKey, nil

	default:
		return nil, errors.Wrapf(types.ErrInvalidDelegationType, "invalid delegation type %v", ud.Type)
	}
}

// SetUnbondingDelegation stores the given unbonding delegation in the store
func (k *Keeper) SetUnbondingDelegation(ctx context.Context, ud types.UnbondingDelegation) ([]byte, error) {
	// Get the key to be used to store the unbonding delegation
	getUnbondingDelegation, err := k.getUnbondingDelegationKeyBuilder(ud)
	if err != nil {
		return nil, err
	}
	unbondingDelegationKey := getUnbondingDelegation(ud.DelegatorAddress, ud.TargetID)

	// Store the unbonding delegation
	store := k.storeService.OpenKVStore(ctx)
	err = store.Set(unbondingDelegationKey, types.MustMarshalUnbondingDelegation(k.cdc, ud))
	if err != nil {
		return nil, err
	}

	return unbondingDelegationKey, nil
}

// SetUnbondingDelegationByUnbondingID sets an index to look up an UnbondingDelegation
// by the unbondingID of an UnbondingDelegationEntry that it contains Note, it does not
// set the unbonding delegation itself, use SetUnbondingDelegation(ctx, ubd) for that
func (k *Keeper) SetUnbondingDelegationByUnbondingID(ctx context.Context, ubd types.UnbondingDelegation, ubdKey []byte, id uint64) error {
	// Set the index allowing to lookup the UnbondingDelegation by the unbondingID of an
	// UnbondingDelegationEntry that it contains
	store := k.storeService.OpenKVStore(ctx)
	err := store.Set(types.GetUnbondingIndexKey(id), ubdKey)
	if err != nil {
		return err
	}

	// Set the type of the unbonding delegation so that we know how to deserialize id
	return store.Set(types.GetUnbondingTypeKey(id), utils.Uint32ToBigEndian(ubd.TargetID))
}

// GetUnbondingDelegation returns the unbonding delegation for the given delegator and target.
func (k *Keeper) GetUnbondingDelegation(
	ctx context.Context, delegatorAddress string, ubdType types.DelegationType, targetID uint32,
) (types.UnbondingDelegation, bool, error) {
	switch ubdType {
	case types.DELEGATION_TYPE_POOL:
		return k.GetPoolUnbondingDelegation(ctx, targetID, delegatorAddress)
	case types.DELEGATION_TYPE_OPERATOR:
		return k.GetOperatorUnbondingDelegation(ctx, targetID, delegatorAddress)
	case types.DELEGATION_TYPE_SERVICE:
		return k.GetServiceUnbondingDelegation(ctx, targetID, delegatorAddress)
	default:
		return types.UnbondingDelegation{}, false, fmt.Errorf("invalid delegation type %v", ubdType)
	}
}

// RemoveUnbondingDelegation removes the unbonding delegation object and associated index.
func (k *Keeper) RemoveUnbondingDelegation(ctx context.Context, ubd types.UnbondingDelegation) error {
	// Get the key to be used to store the unbonding delegation
	getUnbondingDelegation, err := k.getUnbondingDelegationKeyBuilder(ubd)
	if err != nil {
		return err
	}
	unbondingDelegationKey := getUnbondingDelegation(ubd.DelegatorAddress, ubd.TargetID)

	store := k.storeService.OpenKVStore(ctx)
	return store.Delete(unbondingDelegationKey)
}

// PerformUndelegation unbonds an amount of delegator shares from a given validator. It
// will verify that the unbonding entries between the delegator and validator
// are not exceeded and unbond the staked tokens (based on shares) by creating
// an unbonding object and inserting it into the unbonding queue which will be
// processed during the staking EndBlocker.
func (k *Keeper) PerformUndelegation(ctx context.Context, data types.UndelegationData) (time.Time, error) {
	// TODO: Probably we should implement this as well
	// if k.HasMaxUnbondingDelegationEntries(ctx, delAddr, valAddr) {
	//	 return time.Time{}, types.ErrMaxUnbondingDelegationEntries
	// }

	// Unbond the tokens
	returnAmount, err := k.Unbond(ctx, data)
	if err != nil {
		return time.Time{}, err
	}

	// Compute the time at which the unbonding delegation should end
	unbondingTime, err := k.UnbondingTime(ctx)
	if err != nil {
		return time.Time{}, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	completionTime := sdkCtx.BlockHeader().Time.Add(unbondingTime)

	// Store the unbonding delegation entry inside the store
	ubd, err := k.SetUnbondingDelegationEntry(ctx, data, sdkCtx.BlockHeight(), completionTime, returnAmount)
	if err != nil {
		return time.Time{}, err
	}

	// Insert the unbonding delegation into the unbonding queue
	err = k.InsertUBDQueue(ctx, ubd, completionTime)
	if err != nil {
		return time.Time{}, err
	}

	return completionTime, nil
}

// UnbondRestakedAssets unbonds the provided amount from the user's delegations.
// The algorithm will go over the user's delegation in the following order: pools, services and operators
// until the token undelegated matches the provided amount.
func (k *Keeper) UnbondRestakedAssets(ctx context.Context, user sdk.AccAddress, amount sdk.Coins) (time.Time, error) {
	var undelegations []types.UndelegationData
	toUndelegateTokens := sdk.NewDecCoinsFromCoins(amount...)

	err := k.IterateUserDelegations(ctx, user.String(), func(delegation types.Delegation) (bool, error) {
		target, err := k.GetDelegationTargetFromDelegation(ctx, delegation)
		if err != nil {
			if errors.IsOf(err, collections.ErrNotFound) {
				return false, nil
			}
			return true, err
		}

		// Compute the shares that this delegation should have to undelegate
		// all the remaining tokens
		involvedShares, err := target.SharesFromDecCoins(toUndelegateTokens)
		if err != nil {
			return true, err
		}
		// Filter the shares to only the ones that can be removed from the current
		// delegation
		involvedShares = utils.FilterDecCoinsByDenom(involvedShares, delegation.Shares)
		// No shares, keep iterating
		if len(involvedShares) == 0 {
			return false, nil
		}

		toUndelegatedShares := sdk.NewDecCoins()
		for _, share := range involvedShares {
			delShareAmount := delegation.Shares.AmountOf(share.Denom)
			toUndelegatedShares = toUndelegatedShares.Add(
				sdk.NewDecCoinFromDec(share.Denom, math.LegacyMinDec(share.Amount, delShareAmount)))
		}

		// Update the coins to undelegate
		coins := target.TokensFromShares(toUndelegatedShares)
		toUndelegateTokens = toUndelegateTokens.Sub(coins)

		// Update the list of undelegations to perform
		var buildUnbondingDelegation types.UnbondingDelegationBuilder
		var hooks types.DelegationHooks
		switch delegation.Type {
		case types.DELEGATION_TYPE_POOL:
			buildUnbondingDelegation = types.NewPoolUnbondingDelegation
			hooks = types.DelegationHooks{
				BeforeDelegationSharesModified: k.BeforePoolDelegationSharesModified,
				BeforeDelegationCreated:        k.BeforePoolDelegationCreated,
				AfterDelegationModified:        k.AfterPoolDelegationModified,
				BeforeDelegationRemoved:        k.BeforePoolDelegationRemoved,
			}
		case types.DELEGATION_TYPE_SERVICE:
			buildUnbondingDelegation = types.NewServiceUnbondingDelegation
			hooks = types.DelegationHooks{
				BeforeDelegationSharesModified: k.BeforeServiceDelegationSharesModified,
				BeforeDelegationCreated:        k.BeforeServiceDelegationCreated,
				AfterDelegationModified:        k.AfterServiceDelegationModified,
				BeforeDelegationRemoved:        k.BeforeServiceDelegationRemoved,
			}
		case types.DELEGATION_TYPE_OPERATOR:
			buildUnbondingDelegation = types.NewOperatorUnbondingDelegation
			hooks = types.DelegationHooks{
				BeforeDelegationSharesModified: k.BeforeOperatorDelegationSharesModified,
				BeforeDelegationCreated:        k.BeforeOperatorDelegationCreated,
				AfterDelegationModified:        k.AfterOperatorDelegationModified,
				BeforeDelegationRemoved:        k.BeforeOperatorDelegationRemoved,
			}
		default:
			return true, fmt.Errorf("unsupported delegation type: %s", delegation.Type.String())
		}
		truncatedCoins, _ := coins.TruncateDecimal()
		undelegations = append(undelegations, types.UndelegationData{
			Amount:                   truncatedCoins,
			Delegator:                user.String(),
			Target:                   target,
			BuildUnbondingDelegation: buildUnbondingDelegation,
			Hooks:                    hooks,
			Shares:                   toUndelegatedShares,
		})

		// We have finished to undelegate the tokens, stop the iteration.
		if toUndelegateTokens.IsZero() {
			return true, nil
		}

		return false, nil
	})
	if err != nil {
		return time.Time{}, err
	}

	truncatedToUndelegateTokens, _ := toUndelegateTokens.TruncateDecimal()
	if !truncatedToUndelegateTokens.IsZero() {
		return time.Time{}, fmt.Errorf("user hasn't restaked the provided amount: %s", amount.String())
	}

	// Enqueue the undelegations
	var completionTime time.Time
	for _, u := range undelegations {
		ct, err := k.PerformUndelegation(ctx, u)
		if err != nil {
			return time.Time{}, err
		}
		if completionTime.IsZero() {
			completionTime = ct
		}
	}

	return completionTime, nil
}

// --------------------------------------------------------------------------------------------------------------------
// --- Unbonding iterations operations
// --------------------------------------------------------------------------------------------------------------------

// GetAllPoolUnbondingDelegations returns all the pool unbonding delegations
func (k *Keeper) GetAllPoolUnbondingDelegations(ctx context.Context) ([]types.UnbondingDelegation, error) {
	store := k.storeService.OpenKVStore(ctx)

	iterator, err := store.Iterator(types.PoolUnbondingDelegationPrefix, storetypes.PrefixEndBytes(types.PoolUnbondingDelegationPrefix))
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var unbondingDelegations []types.UnbondingDelegation
	for ; iterator.Valid(); iterator.Next() {
		unbondingDelegation := types.MustUnmarshalUnbondingDelegation(k.cdc, iterator.Value())
		unbondingDelegations = append(unbondingDelegations, unbondingDelegation)
	}

	return unbondingDelegations, nil
}

// GetAllUserPoolUnbondingDelegations returns all the user's unbonding delegations
// from a pool
func (k *Keeper) GetAllUserPoolUnbondingDelegations(ctx context.Context, userAddress string) []types.UnbondingDelegation {
	store := k.storeService.OpenKVStore(ctx)

	iterator := storetypes.KVStorePrefixIterator(runtime.KVStoreAdapter(store), types.PoolUnbondingDelegationsStorePrefix(userAddress))
	defer iterator.Close()

	var unbondingDelegations []types.UnbondingDelegation
	for ; iterator.Valid(); iterator.Next() {
		unbondingDelegation := types.MustUnmarshalUnbondingDelegation(k.cdc, iterator.Value())
		unbondingDelegations = append(unbondingDelegations, unbondingDelegation)
	}

	return unbondingDelegations
}

// GetAllOperatorUnbondingDelegations returns all the operator unbonding delegations
func (k *Keeper) GetAllOperatorUnbondingDelegations(ctx context.Context) ([]types.UnbondingDelegation, error) {
	store := k.storeService.OpenKVStore(ctx)

	iterator, err := store.Iterator(types.OperatorUnbondingDelegationPrefix, storetypes.PrefixEndBytes(types.OperatorUnbondingDelegationPrefix))
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var unbondingDelegations []types.UnbondingDelegation
	for ; iterator.Valid(); iterator.Next() {
		unbondingDelegation := types.MustUnmarshalUnbondingDelegation(k.cdc, iterator.Value())
		unbondingDelegations = append(unbondingDelegations, unbondingDelegation)
	}

	return unbondingDelegations, nil
}

// GetAllUserOperatorUnbondingDelegations returns all the user's unbonding delegations
// from an operator
func (k *Keeper) GetAllUserOperatorUnbondingDelegations(ctx context.Context, userAddress string) []types.UnbondingDelegation {
	store := k.storeService.OpenKVStore(ctx)

	iterator := storetypes.KVStorePrefixIterator(runtime.KVStoreAdapter(store), types.OperatorUnbondingDelegationsStorePrefix(userAddress))
	defer iterator.Close()

	var unbondingDelegations []types.UnbondingDelegation
	for ; iterator.Valid(); iterator.Next() {
		unbondingDelegation := types.MustUnmarshalUnbondingDelegation(k.cdc, iterator.Value())
		unbondingDelegations = append(unbondingDelegations, unbondingDelegation)
	}

	return unbondingDelegations
}

// GetAllServiceUnbondingDelegations returns all the service unbonding delegations
func (k *Keeper) GetAllServiceUnbondingDelegations(ctx context.Context) ([]types.UnbondingDelegation, error) {
	store := k.storeService.OpenKVStore(ctx)

	iterator, err := store.Iterator(types.ServiceUnbondingDelegationPrefix, storetypes.PrefixEndBytes(types.ServiceUnbondingDelegationPrefix))
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var unbondingDelegations []types.UnbondingDelegation
	for ; iterator.Valid(); iterator.Next() {
		unbondingDelegation := types.MustUnmarshalUnbondingDelegation(k.cdc, iterator.Value())
		unbondingDelegations = append(unbondingDelegations, unbondingDelegation)
	}

	return unbondingDelegations, nil
}

// GetAllUserServiceUnbondingDelegations returns all the user's unbonding delegations
// from a service
func (k *Keeper) GetAllUserServiceUnbondingDelegations(ctx context.Context, userAddress string) []types.UnbondingDelegation {
	store := k.storeService.OpenKVStore(ctx)

	iterator := storetypes.KVStorePrefixIterator(runtime.KVStoreAdapter(store), types.ServiceUnbondingDelegationsStorePrefix(userAddress))
	defer iterator.Close()

	var unbondingDelegations []types.UnbondingDelegation
	for ; iterator.Valid(); iterator.Next() {
		unbondingDelegation := types.MustUnmarshalUnbondingDelegation(k.cdc, iterator.Value())
		unbondingDelegations = append(unbondingDelegations, unbondingDelegation)
	}

	return unbondingDelegations
}

// GetAllUnbondingDelegations returns all the unbonding delegations
func (k *Keeper) GetAllUnbondingDelegations(ctx context.Context) ([]types.UnbondingDelegation, error) {
	var unbondingDelegations []types.UnbondingDelegation

	unbondingPoolDelegations, err := k.GetAllPoolUnbondingDelegations(ctx)
	if err != nil {
		return nil, err
	}
	unbondingDelegations = append(unbondingDelegations, unbondingPoolDelegations...)

	unbondingOperatorDelegations, err := k.GetAllOperatorUnbondingDelegations(ctx)
	if err != nil {
		return nil, err
	}
	unbondingDelegations = append(unbondingDelegations, unbondingOperatorDelegations...)

	unbondingServiceDelegations, err := k.GetAllServiceUnbondingDelegations(ctx)
	if err != nil {
		return nil, err
	}
	unbondingDelegations = append(unbondingDelegations, unbondingServiceDelegations...)

	return unbondingDelegations, nil
}

// GetAllUserUnbondingDelegations returns all the user's unbonding delegations
func (k *Keeper) GetAllUserUnbondingDelegations(ctx context.Context, userAddress string) []types.UnbondingDelegation {
	var unbondingDelegations []types.UnbondingDelegation

	unbondingDelegations = append(unbondingDelegations, k.GetAllUserPoolUnbondingDelegations(ctx, userAddress)...)
	unbondingDelegations = append(unbondingDelegations, k.GetAllUserOperatorUnbondingDelegations(ctx, userAddress)...)
	unbondingDelegations = append(unbondingDelegations, k.GetAllUserServiceUnbondingDelegations(ctx, userAddress)...)

	return unbondingDelegations
}

// --------------------------------------------------------------------------------------------------------------------

// GetUserPreferencesEntries returns all the user preferences entries
func (k *Keeper) GetUserPreferencesEntries(ctx context.Context) ([]types.UserPreferencesEntry, error) {
	var entries []types.UserPreferencesEntry
	err := k.usersPreferences.Walk(ctx, nil, func(userAddress string, preferences types.UserPreferences) (stop bool, err error) {
		entries = append(entries, types.NewUserPreferencesEntry(userAddress, preferences))
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return entries, nil
}

// --------------------------------------------------------------------------------------------------------------------

// GetTotalRestakedAssets returns the total amount of restaked assets
func (k *Keeper) GetTotalRestakedAssets(ctx context.Context) (sdk.Coins, error) {
	totalRestakedAssets := sdk.NewCoins()

	err := k.poolsKeeper.IteratePools(ctx, func(pool poolstypes.Pool) (bool, error) {
		totalRestakedAssets = totalRestakedAssets.Add(pool.GetTokens()...)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	err = k.operatorsKeeper.IterateOperators(ctx, func(operator operatorstypes.Operator) (bool, error) {
		totalRestakedAssets = totalRestakedAssets.Add(operator.GetTokens()...)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	err = k.servicesKeeper.IterateServices(ctx, func(service servicestypes.Service) (bool, error) {
		totalRestakedAssets = totalRestakedAssets.Add(service.GetTokens()...)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return totalRestakedAssets, nil
}
