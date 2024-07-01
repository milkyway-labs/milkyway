package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

// SaveServiceDelegation stores the given service delegation in the store
func (k *Keeper) SaveServiceDelegation(ctx sdk.Context, delegation types.ServiceDelegation) {
	store := ctx.KVStore(k.storeKey)

	// Marshal and store the delegation
	delegationBz := types.MustMarshalServiceDelegation(k.cdc, delegation)
	store.Set(types.UserServiceDelegationStoreKey(delegation.UserAddress, delegation.ServiceID), delegationBz)

	// Store the delegation in the delegations by service ID store
	store.Set(types.DelegationByServiceIDStoreKey(delegation.ServiceID, delegation.UserAddress), []byte{})
}

// GetServiceDelegation retrieves the delegation for the given user and service
// If the delegation does not exist, false is returned instead
func (k *Keeper) GetServiceDelegation(ctx sdk.Context, serviceID uint32, userAddress string) (types.ServiceDelegation, bool) {
	// Get the delegation amount from the store
	store := ctx.KVStore(k.storeKey)
	delegationAmountBz := store.Get(types.UserServiceDelegationStoreKey(userAddress, serviceID))
	if delegationAmountBz == nil {
		return types.ServiceDelegation{}, false
	}

	// Parse the delegation amount
	return types.MustUnmarshalServiceDelegation(k.cdc, delegationAmountBz), true
}

// AddServiceTokensAndShares adds the given amount of tokens to the service and returns the added shares
func (k *Keeper) AddServiceTokensAndShares(
	ctx sdk.Context, service servicestypes.Service, tokensToAdd sdk.Coins,
) (serviceOut servicestypes.Service, addedShares sdk.DecCoins, err error) {

	// Update the service tokens and shares and get the added shares
	service, addedShares = service.AddTokensFromDelegation(tokensToAdd)

	// Save the service
	k.servicesKeeper.SaveService(ctx, service)
	return service, addedShares, nil
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
		Amount:    amount,
		Delegator: delegator,
		Receiver:  &service,
		GetDelegation: func(ctx sdk.Context, receiverID uint32, delegator string) (types.Delegation, bool) {
			return k.GetServiceDelegation(ctx, receiverID, delegator)
		},
		BuildDelegation: func(receiverID uint32, delegator string) types.Delegation {
			return types.NewServiceDelegation(receiverID, delegator, sdk.NewDecCoins())
		},
		UpdateDelegation: func(ctx sdk.Context, delegation types.Delegation) (newShares sdk.DecCoins, err error) {
			// Calculate the new shares and add the tokens to the service
			_, newShares, err = k.AddServiceTokensAndShares(ctx, service, amount)
			if err != nil {
				return newShares, err
			}

			// Update the delegation shares
			serviceDelegation, ok := delegation.(types.ServiceDelegation)
			if !ok {
				return newShares, fmt.Errorf("invalid delegation type: %T", delegation)
			}
			serviceDelegation.Shares = serviceDelegation.Shares.Add(newShares...)

			// Store the updated delegation
			k.SaveServiceDelegation(ctx, serviceDelegation)

			return newShares, err
		},
		Hooks: types.DelegationHooks{
			BeforeDelegationSharesModified: k.BeforeServiceDelegationSharesModified,
			BeforeDelegationCreated:        k.BeforeServiceDelegationCreated,
			AfterDelegationModified:        k.AfterServiceDelegationModified,
		},
	})
}
