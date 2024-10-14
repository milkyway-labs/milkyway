package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

// SaveServiceParams stored the given params for the given service
func (k *Keeper) SaveServiceParams(ctx sdk.Context, serviceID uint32, params types.ServiceParams) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ServiceParamsStoreKey(serviceID), k.cdc.MustMarshal(&params))
}

// GetServiceParams returns the params for the given service, if any.
// If not params are found, false is returned instead.
func (k *Keeper) GetServiceParams(ctx sdk.Context, serviceID uint32) (params types.ServiceParams) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ServiceParamsStoreKey(serviceID))
	if bz == nil {
		return types.DefaultServiceParams()
	}
	k.cdc.MustUnmarshal(bz, &params)
	return params
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
