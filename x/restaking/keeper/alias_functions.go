package keeper

import (
	"fmt"
	"sort"
	"time"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/utils"
	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

// --------------------------------------------------------------------------------------------------------------------
// --- Params operations
// --------------------------------------------------------------------------------------------------------------------

// IterateAllOperatorsJoinedServices iterates over all the operators and their joined services,
// performing the given action. If the action returns true, the iteration will stop.
func (k *Keeper) IterateAllOperatorsJoinedServices(ctx sdk.Context, action func(operatorID uint32, serviceID uint32) (stop bool, err error)) error {
	iterator, err := k.operatorJoinedServices.Iterate(ctx, nil)
	if err != nil {
		return err
	}
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		operatorServicePair, err := iterator.Key()
		if err != nil {
			return err
		}

		stop, err := action(operatorServicePair.K1(), operatorServicePair.K2())
		if err != nil {
			return err
		}

		if stop {
			break
		}
	}

	return nil
}

// GetAllOperatorsJoinedServices returns all the operators joined services
func (k *Keeper) GetAllOperatorsJoinedServices(ctx sdk.Context) ([]types.OperatorJoinedServices, error) {
	iterator, err := k.operatorJoinedServices.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	items := make(map[uint32]types.OperatorJoinedServices)
	k.IterateAllOperatorsJoinedServices(ctx, func(operatorID uint32, serviceID uint32) (stop bool, err error) {
		joinedServicesRecord, ok := items[operatorID]
		if !ok {
			joinedServicesRecord = types.NewOperatorJoinedServices(operatorID, nil)
		}
		joinedServicesRecord.ServiceIDs = append(joinedServicesRecord.ServiceIDs, serviceID)
		items[operatorID] = joinedServicesRecord
		return false, nil
	})

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
func (k *Keeper) IterateAllServicesAllowedOperators(ctx sdk.Context, action func(serviceID uint32, operatorID uint32) (stop bool, err error)) error {
	iterator, err := k.serviceOperatorsAllowList.Iterate(ctx, nil)
	if err != nil {
		return err
	}
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		serviceOperatorPair, err := iterator.Key()
		if err != nil {
			return err
		}
		serviceID := serviceOperatorPair.K1()
		operatorID := serviceOperatorPair.K2()

		stop, err := action(serviceID, operatorID)
		if err != nil {
			return err
		}

		if stop {
			break
		}
	}

	return nil
}

// GetAllServicesAllowedOperators returns all the operators that are allowed to secure a service for all the services
func (k *Keeper) GetAllServicesAllowedOperators(ctx sdk.Context) ([]types.ServiceAllowedOperators, error) {
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
func (k *Keeper) GetAllServicesSecuringPools(ctx sdk.Context) ([]types.ServiceSecuringPools, error) {
	iterator, err := k.serviceSecuringPools.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	items := make(map[uint32]types.ServiceSecuringPools)
	for ; iterator.Valid(); iterator.Next() {
		servicePoolPair, err := iterator.Key()
		if err != nil {
			return nil, err
		}
		serviceID := servicePoolPair.K1()
		poolID := servicePoolPair.K2()

		securingPools, ok := items[serviceID]
		if !ok {
			securingPools = types.NewServiceSecuringPools(serviceID, nil)
		}
		securingPools.PoolIDs = append(securingPools.PoolIDs, poolID)
		items[serviceID] = securingPools
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

// getDelegationKeyBuilders returns the key builders for the given delegation
func (k *Keeper) getDelegationKeyBuilders(delegation types.Delegation) (types.DelegationKeyBuilder, types.DelegationByTargetIDBuilder, error) {
	switch delegation.Type {
	case types.DELEGATION_TYPE_POOL:
		return types.UserPoolDelegationStoreKey, types.DelegationByPoolIDStoreKey, nil

	case types.DELEGATION_TYPE_OPERATOR:
		return types.UserOperatorDelegationStoreKey, types.DelegationByOperatorIDStoreKey, nil

	case types.DELEGATION_TYPE_SERVICE:
		return types.UserServiceDelegationStoreKey, types.DelegationByServiceIDStoreKey, nil

	default:
		return nil, nil, errors.Wrapf(types.ErrInvalidDelegationType, "invalid delegation type: %v", delegation.Type)
	}
}

// SetDelegation stores the given delegation in the store
func (k *Keeper) SetDelegation(ctx sdk.Context, delegation types.Delegation) error {
	store := ctx.KVStore(k.storeKey)

	// Get the keys builders
	getDelegationKey, getDelegationByTargetID, err := k.getDelegationKeyBuilders(delegation)
	if err != nil {
		return err
	}

	// Marshal and store the delegation
	delegationBz := types.MustMarshalDelegation(k.cdc, delegation)
	store.Set(getDelegationKey(delegation.UserAddress, delegation.TargetID), delegationBz)

	// Store the delegation in the delegations by pool ID store
	store.Set(getDelegationByTargetID(delegation.TargetID, delegation.UserAddress), []byte{})

	return nil
}

// GetDelegationForTarget returns the delegation for the given delegator and target.
func (k *Keeper) GetDelegationForTarget(
	ctx sdk.Context, target types.DelegationTarget, delegator string,
) (types.Delegation, bool) {
	switch target.(type) {
	case *poolstypes.Pool:
		return k.GetPoolDelegation(ctx, target.GetID(), delegator)
	case *operatorstypes.Operator:
		return k.GetOperatorDelegation(ctx, target.GetID(), delegator)
	case *servicestypes.Service:
		return k.GetServiceDelegation(ctx, target.GetID(), delegator)
	default:
		return types.Delegation{}, false
	}
}

// GetDelegationTargetFromDelegation returns the target of the given delegation.
func (k *Keeper) GetDelegationTargetFromDelegation(
	ctx sdk.Context, delegation types.Delegation,
) (types.DelegationTarget, bool) {
	switch delegation.Type {
	case types.DELEGATION_TYPE_POOL:
		if t, found := k.poolsKeeper.GetPool(ctx, delegation.TargetID); found {
			return &t, true
		} else {
			return nil, false
		}
	case types.DELEGATION_TYPE_SERVICE:
		if t, found := k.servicesKeeper.GetService(ctx, delegation.TargetID); found {
			return &t, true
		} else {
			return nil, false
		}
	case types.DELEGATION_TYPE_OPERATOR:
		if t, found := k.operatorsKeeper.GetOperator(ctx, delegation.TargetID); found {
			return &t, true
		} else {
			return nil, false
		}
	default:
		return nil, false
	}
}

// RemoveDelegation removes the given delegation from the store
func (k *Keeper) RemoveDelegation(ctx sdk.Context, delegation types.Delegation) {
	switch delegation.Type {
	case types.DELEGATION_TYPE_POOL:
		k.RemovePoolDelegation(ctx, delegation)
	case types.DELEGATION_TYPE_OPERATOR:
		k.RemoveOperatorDelegation(ctx, delegation)
	case types.DELEGATION_TYPE_SERVICE:
		k.RemoveServiceDelegation(ctx, delegation)
	}
}

// --------------------------------------------------------------------------------------------------------------------
// --- Delegations iterations operations
// --------------------------------------------------------------------------------------------------------------------

// IterateUserPoolDelegations iterates all the pool delegations of a user and performs the given callback function
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

// IterateAllPoolDelegations iterates all the pool delegations and performs the given callback function
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

// IterateUserOperatorDelegations iterates all the operator delegations of a user and performs the given callback function
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

// IterateAllOperatorDelegations iterates all the operator delegations and performs the given callback function
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

// IterateUserServiceDelegations iterates all the service delegations of a user and performs the given callback function
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

// IterateAllServiceDelegations iterates all the service delegations and performs the given callback function
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

// IterateUserDelegations iterates over the user's delegations.
// The delegations will be iterated in the following order:
// 1. Pool delegations
// 2. Service delegations
// 3. Operator delegations
func (k *Keeper) IterateUserDelegations(
	ctx sdk.Context, userAddress string, cb func(del types.Delegation) (stop bool, err error),
) error {
	store := ctx.KVStore(k.storeKey)
	poolIterator := storetypes.KVStorePrefixIterator(store, types.UserPoolDelegationsStorePrefix(userAddress))
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
	servicesIterator := storetypes.KVStorePrefixIterator(store, types.UserServiceDelegationsStorePrefix(userAddress))
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
	operatorsIterator := storetypes.KVStorePrefixIterator(store, types.UserOperatorDelegationsStorePrefix(userAddress))
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
func (k *Keeper) GetAllDelegations(ctx sdk.Context) []types.Delegation {
	var delegations []types.Delegation

	delegations = append(delegations, k.GetAllPoolDelegations(ctx)...)
	delegations = append(delegations, k.GetAllOperatorDelegations(ctx)...)
	delegations = append(delegations, k.GetAllServiceDelegations(ctx)...)

	return delegations
}

// GetAllUserRestakedCoins returns all the user's restaked coins
func (k *Keeper) GetAllUserRestakedCoins(ctx sdk.Context, userAddress string) (sdk.DecCoins, error) {
	totalDelegatedCoins := sdk.NewDecCoins()
	k.IterateUserDelegations(ctx, userAddress, func(d types.Delegation) (bool, error) {
		target, found := k.GetDelegationTargetFromDelegation(ctx, d)
		if !found {
			return true, fmt.Errorf("can't find target for delegation %d, target id: %d", d.Type, d.TargetID)
		}
		totalDelegatedCoins = totalDelegatedCoins.Add(target.TokensFromShares(d.Shares)...)
		return false, nil
	})

	return totalDelegatedCoins, nil
}

// PerformDelegation performs a delegation of the given amount from the delegator to the receiver.
// It sends the coins to the receiver address and updates the delegation object and returns the new
// shares of the delegation.
// NOTE: This is done so that if we implement other delegation types in the future we can have a single
// function that performs common operations for all of them.
func (k *Keeper) PerformDelegation(ctx sdk.Context, data types.DelegationData) (sdk.DecCoins, error) {
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

	// Get or create the delegation object and call the appropriate hook if present
	delegation, found := k.GetDelegationForTarget(ctx, receiver, delegator)

	if found {
		// Delegation was found
		err := hooks.BeforeDelegationSharesModified(ctx, receiver.GetID(), delegator)
		if err != nil {
			return nil, err
		}
	} else {
		// Delegation was not found
		delegation = data.BuildDelegation(receiver.GetID(), delegator, sdk.NewDecCoins())
		err := hooks.BeforeDelegationCreated(ctx, receiver.GetID(), delegator)
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
func (k *Keeper) getUnbondingDelegationTarget(ctx sdk.Context, ubd types.UnbondingDelegation) (types.DelegationTarget, error) {
	switch ubd.Type {
	case types.DELEGATION_TYPE_POOL:
		pool, found := k.poolsKeeper.GetPool(ctx, ubd.TargetID)
		if !found {
			return nil, poolstypes.ErrPoolNotFound
		}
		return &pool, nil

	case types.DELEGATION_TYPE_OPERATOR:
		operator, found := k.operatorsKeeper.GetOperator(ctx, ubd.TargetID)
		if !found {
			return nil, operatorstypes.ErrOperatorNotFound
		}
		return &operator, nil

	case types.DELEGATION_TYPE_SERVICE:
		service, found := k.servicesKeeper.GetService(ctx, ubd.TargetID)
		if !found {
			return nil, servicestypes.ErrServiceNotFound
		}
		return &service, nil

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
func (k *Keeper) SetUnbondingDelegation(ctx sdk.Context, ud types.UnbondingDelegation) ([]byte, error) {
	// Get the key to be used to store the unbonding delegation
	getUnbondingDelegation, err := k.getUnbondingDelegationKeyBuilder(ud)
	if err != nil {
		return nil, err
	}
	unbondingDelegationKey := getUnbondingDelegation(ud.DelegatorAddress, ud.TargetID)

	// Store the unbonding delegation
	store := ctx.KVStore(k.storeKey)
	store.Set(unbondingDelegationKey, types.MustMarshalUnbondingDelegation(k.cdc, ud))

	return unbondingDelegationKey, nil
}

// SetUnbondingDelegationByUnbondingID sets an index to look up an UnbondingDelegation
// by the unbondingID of an UnbondingDelegationEntry that it contains Note, it does not
// set the unbonding delegation itself, use SetUnbondingDelegation(ctx, ubd) for that
func (k *Keeper) SetUnbondingDelegationByUnbondingID(ctx sdk.Context, ubd types.UnbondingDelegation, ubdKey []byte, id uint64) {
	// Set the index allowing to lookup the UnbondingDelegation by the unbondingID of an
	// UnbondingDelegationEntry that it contains
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetUnbondingIndexKey(id), ubdKey)

	// Set the type of the unbonding delegation so that we know how to deserialize id
	store.Set(types.GetUnbondingTypeKey(id), utils.Uint32ToBigEndian(ubd.TargetID))
}

// GetUnbondingDelegation returns the unbonding delegation for the given delegator and target.
func (k *Keeper) GetUnbondingDelegation(
	ctx sdk.Context, delegatorAddress string, ubdType types.DelegationType, targetID uint32,
) (types.UnbondingDelegation, bool) {
	switch ubdType {
	case types.DELEGATION_TYPE_POOL:
		return k.GetPoolUnbondingDelegation(ctx, targetID, delegatorAddress)
	case types.DELEGATION_TYPE_OPERATOR:
		return k.GetOperatorUnbondingDelegation(ctx, targetID, delegatorAddress)
	case types.DELEGATION_TYPE_SERVICE:
		return k.GetServiceUnbondingDelegation(ctx, targetID, delegatorAddress)
	default:
		return types.UnbondingDelegation{}, false
	}
}

// RemoveUnbondingDelegation removes the unbonding delegation object and associated index.
func (k *Keeper) RemoveUnbondingDelegation(ctx sdk.Context, ubd types.UnbondingDelegation) error {
	// Get the key to be used to store the unbonding delegation
	getUnbondingDelegation, err := k.getUnbondingDelegationKeyBuilder(ubd)
	if err != nil {
		return err
	}
	unbondingDelegationKey := getUnbondingDelegation(ubd.DelegatorAddress, ubd.TargetID)

	store := ctx.KVStore(k.storeKey)
	store.Delete(unbondingDelegationKey)
	return nil
}

// PerformUndelegation unbonds an amount of delegator shares from a given validator. It
// will verify that the unbonding entries between the delegator and validator
// are not exceeded and unbond the staked tokens (based on shares) by creating
// an unbonding object and inserting it into the unbonding queue which will be
// processed during the staking EndBlocker.
func (k *Keeper) PerformUndelegation(ctx sdk.Context, data types.UndelegationData) (time.Time, error) {
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
	completionTime := ctx.BlockHeader().Time.Add(k.UnbondingTime(ctx))

	// Store the unbonding delegation entry inside the store
	ubd, err := k.SetUnbondingDelegationEntry(ctx, data, ctx.BlockHeight(), completionTime, returnAmount)
	if err != nil {
		return time.Time{}, err
	}

	// Insert the unbonding delegation into the unbonding queue
	k.InsertUBDQueue(ctx, ubd, completionTime)

	return completionTime, nil
}

// UnbondRestakedAssets unbonds the provided amount from the user's delegations.
// The algorithm will go over the user's delegation in the following order: pools, services and operators
// until the token undelegated matches the provided amount.
func (k *Keeper) UnbondRestakedAssets(ctx sdk.Context, user sdk.AccAddress, amount sdk.Coins) (time.Time, error) {
	var undelegations []types.UndelegationData
	toUndelegateTokens := sdk.NewDecCoinsFromCoins(amount...)

	err := k.IterateUserDelegations(ctx, user.String(), func(delegation types.Delegation) (bool, error) {
		target, found := k.GetDelegationTargetFromDelegation(ctx, delegation)
		if !found {
			return false, nil
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

	truncatedToUndalegateTokens, _ := toUndelegateTokens.TruncateDecimal()
	if !truncatedToUndalegateTokens.IsZero() {
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
func (k *Keeper) GetAllPoolUnbondingDelegations(ctx sdk.Context) []types.UnbondingDelegation {
	store := ctx.KVStore(k.storeKey)
	iterator := store.Iterator(types.PoolUnbondingDelegationPrefix, storetypes.PrefixEndBytes(types.PoolUnbondingDelegationPrefix))
	defer iterator.Close()

	var unbondingDelegations []types.UnbondingDelegation
	for ; iterator.Valid(); iterator.Next() {
		unbondingDelegation := types.MustUnmarshalUnbondingDelegation(k.cdc, iterator.Value())
		unbondingDelegations = append(unbondingDelegations, unbondingDelegation)
	}

	return unbondingDelegations
}

// GetAllUserPoolUnbondingDelegations returns all the user's unbonding delegations
// from a pool
func (k *Keeper) GetAllUserPoolUnbondingDelegations(ctx sdk.Context, userAddress string) []types.UnbondingDelegation {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.PoolUnbondingDelegationsStorePrefix(userAddress))
	defer iterator.Close()

	var unbondingDelegations []types.UnbondingDelegation
	for ; iterator.Valid(); iterator.Next() {
		unbondingDelegation := types.MustUnmarshalUnbondingDelegation(k.cdc, iterator.Value())
		unbondingDelegations = append(unbondingDelegations, unbondingDelegation)
	}

	return unbondingDelegations
}

// GetAllOperatorUnbondingDelegations returns all the operator unbonding delegations
func (k *Keeper) GetAllOperatorUnbondingDelegations(ctx sdk.Context) []types.UnbondingDelegation {
	store := ctx.KVStore(k.storeKey)
	iterator := store.Iterator(types.OperatorUnbondingDelegationPrefix, storetypes.PrefixEndBytes(types.OperatorUnbondingDelegationPrefix))
	defer iterator.Close()

	var unbondingDelegations []types.UnbondingDelegation
	for ; iterator.Valid(); iterator.Next() {
		unbondingDelegation := types.MustUnmarshalUnbondingDelegation(k.cdc, iterator.Value())
		unbondingDelegations = append(unbondingDelegations, unbondingDelegation)
	}

	return unbondingDelegations
}

// GetAllUserOperatorUnbondingDelegations returns all the user's unbonding delegations
// from a operator
func (k *Keeper) GetAllUserOperatorUnbondingDelegations(ctx sdk.Context, userAddress string) []types.UnbondingDelegation {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.OperatorUnbondingDelegationsStorePrefix(userAddress))
	defer iterator.Close()

	var unbondingDelegations []types.UnbondingDelegation
	for ; iterator.Valid(); iterator.Next() {
		unbondingDelegation := types.MustUnmarshalUnbondingDelegation(k.cdc, iterator.Value())
		unbondingDelegations = append(unbondingDelegations, unbondingDelegation)
	}

	return unbondingDelegations
}

// GetAllServiceUnbondingDelegations returns all the service unbonding delegations
func (k *Keeper) GetAllServiceUnbondingDelegations(ctx sdk.Context) []types.UnbondingDelegation {
	store := ctx.KVStore(k.storeKey)
	iterator := store.Iterator(types.ServiceUnbondingDelegationPrefix, storetypes.PrefixEndBytes(types.ServiceUnbondingDelegationPrefix))
	defer iterator.Close()

	var unbondingDelegations []types.UnbondingDelegation
	for ; iterator.Valid(); iterator.Next() {
		unbondingDelegation := types.MustUnmarshalUnbondingDelegation(k.cdc, iterator.Value())
		unbondingDelegations = append(unbondingDelegations, unbondingDelegation)
	}

	return unbondingDelegations
}

// GetAllUserServiceUnbondingDelegations returns all the user's unbonding delegations
// from a service
func (k *Keeper) GetAllUserServiceUnbondingDelegations(ctx sdk.Context, userAddress string) []types.UnbondingDelegation {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.ServiceUnbondingDelegationsStorePrefix(userAddress))
	defer iterator.Close()

	var unbondingDelegations []types.UnbondingDelegation
	for ; iterator.Valid(); iterator.Next() {
		unbondingDelegation := types.MustUnmarshalUnbondingDelegation(k.cdc, iterator.Value())
		unbondingDelegations = append(unbondingDelegations, unbondingDelegation)
	}

	return unbondingDelegations
}

// GetAllUnbondingDelegations returns all the unbonding delegations
func (k *Keeper) GetAllUnbondingDelegations(ctx sdk.Context) []types.UnbondingDelegation {
	var unbondingDelegations []types.UnbondingDelegation

	unbondingDelegations = append(unbondingDelegations, k.GetAllPoolUnbondingDelegations(ctx)...)
	unbondingDelegations = append(unbondingDelegations, k.GetAllOperatorUnbondingDelegations(ctx)...)
	unbondingDelegations = append(unbondingDelegations, k.GetAllServiceUnbondingDelegations(ctx)...)

	return unbondingDelegations
}

// GetAllUserUnbondingDelegations returns all the user's unbonding delegations
func (k *Keeper) GetAllUserUnbondingDelegations(ctx sdk.Context, userAddress string) []types.UnbondingDelegation {
	var unbondingDelegations []types.UnbondingDelegation

	unbondingDelegations = append(unbondingDelegations, k.GetAllUserPoolUnbondingDelegations(ctx, userAddress)...)
	unbondingDelegations = append(unbondingDelegations, k.GetAllUserOperatorUnbondingDelegations(ctx, userAddress)...)
	unbondingDelegations = append(unbondingDelegations, k.GetAllUserServiceUnbondingDelegations(ctx, userAddress)...)

	return unbondingDelegations
}
