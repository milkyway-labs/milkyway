package keeper

import (
	"context"
	"slices"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v4/utils"
	"github.com/milkyway-labs/milkyway/v4/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v4/x/services/types"
)

// ServiceAllowedOperatorsIterator returns an iterator that iterates over all
// operators allowed to secure a service
func (k *Keeper) ServiceAllowedOperatorsIterator(ctx context.Context, serviceID uint32) (collections.KeySetIterator[collections.Pair[uint32, uint32]], error) {
	return k.serviceOperatorsAllowList.Iterate(ctx, collections.NewPrefixedPairRange[uint32, uint32](serviceID))
}

// GetAllServiceAllowedOperators returns all operators that have been whitelisted
// by a service
func (k *Keeper) GetAllServiceAllowedOperators(ctx context.Context, serviceID uint32) ([]uint32, error) {
	iterator, err := k.ServiceAllowedOperatorsIterator(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var operators []uint32
	for ; iterator.Valid(); iterator.Next() {
		serviceOperatorPair, err := iterator.Key()
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
func (k *Keeper) AddOperatorToServiceAllowList(ctx context.Context, serviceID uint32, operatorID uint32) error {
	key := collections.Join(serviceID, operatorID)
	return k.serviceOperatorsAllowList.Set(ctx, key)
}

// RemoveOperatorFromServiceAllowList removes an operator from the list of operators
// allowed to secure a service.
// If the operator is not in the list, no action is taken.
func (k *Keeper) RemoveOperatorFromServiceAllowList(ctx context.Context, serviceID uint32, operatorID uint32) error {
	key := collections.Join(serviceID, operatorID)
	return k.serviceOperatorsAllowList.Remove(ctx, key)
}

// IsServiceOperatorsAllowListConfigured returns true if the operators allow list
func (k *Keeper) IsServiceOperatorsAllowListConfigured(ctx context.Context, serviceID uint32) (bool, error) {
	iterator, err := k.ServiceAllowedOperatorsIterator(ctx, serviceID)
	if err != nil {
		return false, err
	}
	defer iterator.Close()
	return iterator.Valid(), nil
}

// IsOperatorInServiceAllowList returns true if the given operator is in the
// service operators allow list
func (k *Keeper) IsOperatorInServiceAllowList(ctx context.Context, serviceID uint32, operatorID uint32) (bool, error) {
	key := collections.Join(serviceID, operatorID)
	return k.serviceOperatorsAllowList.Has(ctx, key)
}

// CanOperatorValidateService returns true if the given operator can secure
// the given service
func (k *Keeper) CanOperatorValidateService(ctx context.Context, serviceID uint32, operatorID uint32) (bool, error) {
	configured, err := k.IsServiceOperatorsAllowListConfigured(ctx, serviceID)
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
func (k *Keeper) ServiceSecuringPoolsIterator(ctx context.Context, serviceID uint32) (collections.KeySetIterator[collections.Pair[uint32, uint32]], error) {
	return k.serviceSecuringPools.Iterate(ctx, collections.NewPrefixedPairRange[uint32, uint32](serviceID))
}

// GetAllServiceSecuringPools returns all pools that have been allowed to
// secure the give service
func (k *Keeper) GetAllServiceSecuringPools(ctx context.Context, serviceID uint32) ([]uint32, error) {
	iterator, err := k.ServiceSecuringPoolsIterator(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var pools []uint32
	for ; iterator.Valid(); iterator.Next() {
		servicePoolPair, err := iterator.Key()
		if err != nil {
			return nil, err
		}
		pools = append(pools, servicePoolPair.K2())
	}

	return pools, nil
}

// AddPoolToServiceSecuringPools adds a pool to the list of pools
// permitted for securing the service
func (k *Keeper) AddPoolToServiceSecuringPools(ctx context.Context, serviceID uint32, poolID uint32) error {
	key := collections.Join(serviceID, poolID)
	return k.serviceSecuringPools.Set(ctx, key)
}

// RemovePoolFromServiceSecuringPools removes a pool from the list of pools from which the
// service is borrowing the security from
func (k *Keeper) RemovePoolFromServiceSecuringPools(ctx context.Context, serviceID uint32, poolID uint32) error {
	key := collections.Join(serviceID, poolID)
	return k.serviceSecuringPools.Remove(ctx, key)
}

// IsServiceSecuringPoolsConfigured returns true if the list of securing pools
// has been configured for the given service
func (k *Keeper) IsServiceSecuringPoolsConfigured(ctx context.Context, serviceID uint32) (bool, error) {
	iterator, err := k.ServiceSecuringPoolsIterator(ctx, serviceID)
	if err != nil {
		return false, err
	}
	defer iterator.Close()
	return iterator.Valid(), nil
}

// IsPoolInServiceSecuringPools returns true if the pool is in the list
// of pools from which the service can borrow security
func (k *Keeper) IsPoolInServiceSecuringPools(ctx context.Context, serviceID uint32, poolID uint32) (bool, error) {
	key := collections.Join(serviceID, poolID)
	return k.serviceSecuringPools.Has(ctx, key)
}

// IsServiceSecuredByPool returns true if the service is being secured
// by the given pool
func (k *Keeper) IsServiceSecuredByPool(ctx context.Context, serviceID uint32, poolID uint32) (bool, error) {
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
func (k *Keeper) GetServiceDelegation(ctx context.Context, serviceID uint32, userAddress string) (types.Delegation, bool, error) {
	store := k.storeService.OpenKVStore(ctx)

	delegationBz, err := store.Get(types.UserServiceDelegationStoreKey(userAddress, serviceID))
	if err != nil {
		return types.Delegation{}, false, err
	}

	if delegationBz == nil {
		return types.Delegation{}, false, nil
	}

	return types.MustUnmarshalDelegation(k.cdc, delegationBz), true, nil
}

// AddServiceTokensAndShares adds the given amount of tokens to the service and returns the added shares
func (k *Keeper) AddServiceTokensAndShares(
	ctx context.Context, service servicestypes.Service, tokensToAdd sdk.Coins,
) (serviceOut servicestypes.Service, addedShares sdk.DecCoins, err error) {
	// Update the service tokens and shares and get the added shares
	service, addedShares = service.AddTokensFromDelegation(tokensToAdd)

	// Save the service
	err = k.servicesKeeper.SaveService(ctx, service)
	return service, addedShares, err
}

// RemoveServiceDelegation removes the given service delegation from the store
func (k *Keeper) RemoveServiceDelegation(ctx context.Context, delegation types.Delegation) error {
	store := k.storeService.OpenKVStore(ctx)

	err := store.Delete(types.UserServiceDelegationStoreKey(delegation.UserAddress, delegation.TargetID))
	if err != nil {
		return err
	}

	return store.Delete(types.DelegationByServiceIDStoreKey(delegation.TargetID, delegation.UserAddress))
}

// --------------------------------------------------------------------------------------------------------------------

// DelegateToService sends the given amount to the service account and saves the delegation for the given user
func (k *Keeper) DelegateToService(ctx context.Context, serviceID uint32, amount sdk.Coins, delegator string) (sdk.DecCoins, error) {
	// Get the service
	service, err := k.servicesKeeper.GetService(ctx, serviceID)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return sdk.NewDecCoins(), servicestypes.ErrServiceNotFound
		}
		return nil, err
	}

	restakableDenoms, err := k.GetRestakableDenoms(ctx)
	if err != nil {
		return nil, err
	}

	// Get the service parameters
	serviceParams, err := k.servicesKeeper.GetServiceParams(ctx, serviceID)
	if err != nil {
		return sdk.NewDecCoins(), err
	}

	// Update restakable denoms with the ones allowed by the service
	if len(restakableDenoms) == 0 {
		// No restakable denoms configured, use the service restakable denoms
		restakableDenoms = serviceParams.AllowedDenoms
	} else if len(serviceParams.AllowedDenoms) > 0 {
		// We have both x/restaking restakable denoms and the service
		// restakable denoms, intersect them
		restakableDenoms = utils.Intersect(restakableDenoms, serviceParams.AllowedDenoms)
		if len(restakableDenoms) == 0 {
			// The intersection is empty, this service doesn't allow any delegation
			return sdk.NewDecCoins(), errors.Wrapf(types.ErrDenomNotRestakable, "tokens cannot be restaked")
		}
	}

	// Ensure the provided amount can be restaked
	if len(restakableDenoms) > 0 {
		for _, coin := range amount {
			isRestakable := slices.Contains(restakableDenoms, coin.Denom)
			if !isRestakable {
				return sdk.NewDecCoins(), errors.Wrapf(types.ErrDenomNotRestakable, "%s cannot be restaked", coin.Denom)
			}
		}
	}

	// Make sure the service is active
	if !service.IsActive() {
		return sdk.NewDecCoins(), servicestypes.ErrServiceNotActive
	}

	return k.PerformDelegation(ctx, types.DelegationData{
		Amount:          amount,
		Delegator:       delegator,
		Target:          service,
		BuildDelegation: types.NewServiceDelegation,
		UpdateDelegation: func(ctx context.Context, delegation types.Delegation) (newShares sdk.DecCoins, err error) {
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
func (k *Keeper) GetServiceUnbondingDelegation(ctx context.Context, serviceID uint32, delegator string) (types.UnbondingDelegation, bool, error) {
	store := k.storeService.OpenKVStore(ctx)

	ubdBz, err := store.Get(types.UserServiceUnbondingDelegationKey(delegator, serviceID))
	if err != nil {
		return types.UnbondingDelegation{}, false, err
	}

	if ubdBz == nil {
		return types.UnbondingDelegation{}, false, nil
	}

	return types.MustUnmarshalUnbondingDelegation(k.cdc, ubdBz), true, nil
}

// UndelegateFromService removes the given amount from the service account and saves the
// unbonding delegation for the given user
func (k *Keeper) UndelegateFromService(ctx context.Context, serviceID uint32, amount sdk.Coins, delegator string) (time.Time, error) {
	// Find the service
	service, err := k.servicesKeeper.GetService(ctx, serviceID)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return time.Time{}, servicestypes.ErrServiceNotFound
		}
		return time.Time{}, err
	}

	// Get the shares
	shares, err := k.ValidateUnbondAmount(ctx, delegator, service, amount)
	if err != nil {
		return time.Time{}, err
	}

	return k.PerformUndelegation(ctx, types.UndelegationData{
		Amount:                   amount,
		Delegator:                delegator,
		Target:                   service,
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
