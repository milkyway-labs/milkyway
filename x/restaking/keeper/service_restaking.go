package keeper

import (
	"time"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

// ServiceAllowedOperatorsIterator returns an iterator that iterates over all
// operators allowed to secure a service
func (k *Keeper) ServiceAllowedOperatorsIterator(ctx sdk.Context, serviceID uint32) (collections.KeySetIterator[collections.Pair[uint32, uint32]], error) {
	return k.serviceOperatorsAllowList.Iterate(ctx, collections.NewPrefixedPairRange[uint32, uint32](serviceID))
}

// GetAllServiceAllowedOperators returns all operators that have been whitelisted
// by a service
func (k *Keeper) GetAllServiceAllowedOperators(ctx sdk.Context, serviceID uint32) ([]uint32, error) {
	iteretor, err := k.ServiceAllowedOperatorsIterator(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	defer iteretor.Close()
	var operators []uint32
	for ; iteretor.Valid(); iteretor.Next() {
		serviceOperatorPair, err := iteretor.Key()
		if err != nil {
			return nil, err
		}
		operators = append(operators, serviceOperatorPair.K2())
	}

	return operators, nil
}

// AddOperatorToServiceAllowList adds an operator to the list of operators
// allowed to secure a service.
// If the operator is already in the list, no action is taken.
func (k *Keeper) AddOperatorToServiceAllowList(ctx sdk.Context, serviceID uint32, operatorID uint32) error {
	key := collections.Join(serviceID, operatorID)
	return k.serviceOperatorsAllowList.Set(ctx, key)
}

// RemoveOperatorFromServiceAllowList removes an operator from the list of operators
// allowed to secure a service.
// If the operator is not in the list, no action is taken.
func (k *Keeper) RemoveOperatorFromServiceAllowList(ctx sdk.Context, serviceID uint32, operatorID uint32) error {
	key := collections.Join(serviceID, operatorID)
	return k.serviceOperatorsAllowList.Remove(ctx, key)
}

// IsServiceOpertorsAllowListConfigured returns true if the operators allow list
func (k *Keeper) IsServiceOpertorsAllowListConfigured(ctx sdk.Context, serviceID uint32) (bool, error) {
	iteretor, err := k.ServiceAllowedOperatorsIterator(ctx, serviceID)
	if err != nil {
		return false, err
	}
	defer iteretor.Close()

	for ; iteretor.Valid(); iteretor.Next() {
		return true, nil
	}

	return false, nil
}

// IsOperatorInServiceAllowList returns true if the given operator is in the
// service operators allow list
func (k *Keeper) IsOperatorInServiceAllowList(ctx sdk.Context, serviceID uint32, operatorID uint32) (bool, error) {
	key := collections.Join(serviceID, operatorID)
	return k.serviceOperatorsAllowList.Has(ctx, key)
}

// CanOpeartorValidateService returns true if the given operator can secure
// the given service
func (k *Keeper) CanOperatorValidateService(ctx sdk.Context, serviceID uint32, operatorID uint32) (bool, error) {
	configured, err := k.IsServiceOpertorsAllowListConfigured(ctx, serviceID)
	if err != nil {
		return false, err
	}
	// Allow all when the list is empty
	if !configured {
		return true, nil
	}

	return k.IsOperatorInServiceAllowList(ctx, serviceID, operatorID)
}

// --------------------------------------------------------------------------------------------------------------------

// ServiceSecuringPoolsIterator returns an iterator that iterates over all
// pools allowed to secure the given service.
func (k *Keeper) ServiceSecuringPoolsIterator(ctx sdk.Context, serviceID uint32) (collections.KeySetIterator[collections.Pair[uint32, uint32]], error) {
	return k.serviceSecuringPools.Iterate(ctx, collections.NewPrefixedPairRange[uint32, uint32](serviceID))
}

// GetAllServiceSecuringPools returns all pools that have been allowed to
// to secure the give service
func (k *Keeper) GetAllServiceSecuringPools(ctx sdk.Context, serviceID uint32) ([]uint32, error) {
	iteretor, err := k.ServiceSecuringPoolsIterator(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	defer iteretor.Close()
	var pools []uint32
	for ; iteretor.Valid(); iteretor.Next() {
		servicePoolPair, err := iteretor.Key()
		if err != nil {
			return nil, err
		}
		pools = append(pools, servicePoolPair.K2())
	}

	return pools, nil
}

// AddPoolToServiceSecuringPools adds a pool to the list of pools
// permitted for securing the service
func (k *Keeper) AddPoolToServiceSecuringPools(ctx sdk.Context, serviceID uint32, poolID uint32) error {
	key := collections.Join(serviceID, poolID)
	return k.serviceSecuringPools.Set(ctx, key)
}

// IsServiceSecuringPoolsConfigured returns true if the list of securing pools
// has been configured for the given service
func (k *Keeper) IsServiceSecuringPoolsConfigured(ctx sdk.Context, serviceID uint32) (bool, error) {
	iteretor, err := k.ServiceSecuringPoolsIterator(ctx, serviceID)
	if err != nil {
		return false, err
	}
	defer iteretor.Close()

	for ; iteretor.Valid(); iteretor.Next() {
		return true, nil
	}

	return false, nil
}

// IsServiceSecuredByPool returns true if the service is being secured
// by the given pool
func (k *Keeper) IsServiceSecuredByPool(ctx sdk.Context, serviceID uint32, poolID uint32) (bool, error) {
	configured, err := k.IsServiceSecuringPoolsConfigured(ctx, serviceID)
	if err != nil {
		return false, err
	}
	// Allow all when the list is empty
	if !configured {
		return true, nil
	}

	key := collections.Join(serviceID, poolID)
	return k.serviceSecuringPools.Has(ctx, key)
}

// --------------------------------------------------------------------------------------------------------------------

// GetServiceDelegation retrieves the delegation for the given user and service
// If the delegation does not exist, false is returned instead
func (k *Keeper) GetServiceDelegation(ctx sdk.Context, serviceID uint32, userAddress string) (types.Delegation, bool) {
	store := ctx.KVStore(k.storeKey)
	delegationBz := store.Get(types.UserServiceDelegationStoreKey(userAddress, serviceID))
	if delegationBz == nil {
		return types.Delegation{}, false
	}

	return types.MustUnmarshalDelegation(k.cdc, delegationBz), true
}

// AddServiceTokensAndShares adds the given amount of tokens to the service and returns the added shares
func (k *Keeper) AddServiceTokensAndShares(
	ctx sdk.Context, service servicestypes.Service, tokensToAdd sdk.Coins,
) (serviceOut servicestypes.Service, addedShares sdk.DecCoins, err error) {
	// Update the service tokens and shares and get the added shares
	service, addedShares = service.AddTokensFromDelegation(tokensToAdd)

	// Save the service
	err = k.servicesKeeper.SaveService(ctx, service)
	return service, addedShares, err
}

// RemoveServiceDelegation removes the given service delegation from the store
func (k *Keeper) RemoveServiceDelegation(ctx sdk.Context, delegation types.Delegation) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.UserServiceDelegationStoreKey(delegation.UserAddress, delegation.TargetID))
	store.Delete(types.DelegationByServiceIDStoreKey(delegation.TargetID, delegation.UserAddress))
}

// --------------------------------------------------------------------------------------------------------------------

// DelegateToService sends the given amount to the service account and saves the delegation for the given user
func (k *Keeper) DelegateToService(ctx sdk.Context, serviceID uint32, amount sdk.Coins, delegator string) (sdk.DecCoins, error) {
	// Get the service
	service, found := k.servicesKeeper.GetService(ctx, serviceID)
	if !found {
		return sdk.NewDecCoins(), servicestypes.ErrServiceNotFound
	}

	// Make sure the service is active
	if !service.IsActive() {
		return sdk.NewDecCoins(), servicestypes.ErrServiceNotActive
	}

	return k.PerformDelegation(ctx, types.DelegationData{
		Amount:          amount,
		Delegator:       delegator,
		Target:          &service,
		BuildDelegation: types.NewServiceDelegation,
		UpdateDelegation: func(ctx sdk.Context, delegation types.Delegation) (newShares sdk.DecCoins, err error) {
			// Calculate the new shares and add the tokens to the service
			_, newShares, err = k.AddServiceTokensAndShares(ctx, service, amount)
			if err != nil {
				return newShares, err
			}

			// Update the delegation shares
			delegation.Shares = delegation.Shares.Add(newShares...)

			// Store the updated delegation
			err = k.SetDelegation(ctx, delegation)
			if err != nil {
				return nil, err
			}

			return newShares, err
		},
		Hooks: types.DelegationHooks{
			BeforeDelegationSharesModified: k.BeforeServiceDelegationSharesModified,
			BeforeDelegationCreated:        k.BeforeServiceDelegationCreated,
			AfterDelegationModified:        k.AfterServiceDelegationModified,
		},
	})
}

// --------------------------------------------------------------------------------------------------------------------

// GetServiceUnbondingDelegation returns the unbonding delegation for the given delegator address and service id.
// If no unbonding delegation is found, false is returned instead.
func (k *Keeper) GetServiceUnbondingDelegation(ctx sdk.Context, serviceID uint32, delegator string) (types.UnbondingDelegation, bool) {
	store := ctx.KVStore(k.storeKey)
	ubdBz := store.Get(types.UserServiceUnbondingDelegationKey(delegator, serviceID))
	if ubdBz == nil {
		return types.UnbondingDelegation{}, false
	}

	return types.MustUnmarshalUnbondingDelegation(k.cdc, ubdBz), true
}

// UndelegateFromService removes the given amount from the service account and saves the
// unbonding delegation for the given user
func (k *Keeper) UndelegateFromService(ctx sdk.Context, serviceID uint32, amount sdk.Coins, delegator string) (time.Time, error) {
	// Find the service
	service, found := k.servicesKeeper.GetService(ctx, serviceID)
	if !found {
		return time.Time{}, servicestypes.ErrServiceNotFound
	}

	// Get the shares
	shares, err := k.ValidateUnbondAmount(ctx, delegator, &service, amount)
	if err != nil {
		return time.Time{}, err
	}

	return k.PerformUndelegation(ctx, types.UndelegationData{
		Amount:                   amount,
		Delegator:                delegator,
		Target:                   &service,
		BuildUnbondingDelegation: types.NewServiceUnbondingDelegation,
		Hooks: types.DelegationHooks{
			BeforeDelegationSharesModified: k.BeforeServiceDelegationSharesModified,
			BeforeDelegationCreated:        k.BeforeServiceDelegationCreated,
			AfterDelegationModified:        k.AfterServiceDelegationModified,
			BeforeDelegationRemoved:        k.BeforeServiceDelegationRemoved,
		},
		Shares: shares,
	})
}
