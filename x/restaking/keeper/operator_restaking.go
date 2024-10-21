package keeper

import (
	"errors"
	"time"

	"cosmossdk.io/collections"
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

// GetOperatorJoinedServices gets the services joined by the operator with the given ID.
func (k *Keeper) GetOperatorJoinedServices(ctx sdk.Context, operatorID uint32) (types.OperatorJoinedServices, error) {
	joinedServices, err := k.operatorJoinedServices.Get(ctx, operatorID)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return types.NewEmptyOperatorJoinedServices(), nil
		} else {
			return types.OperatorJoinedServices{}, err
		}
	}
	return joinedServices, nil
}

// SaveOperatorJoinedServices sets the services joined by the operator with the
// given ID.
func (k *Keeper) SaveOperatorJoinedServices(
	ctx sdk.Context,
	operatorID uint32,
	joinedServices types.OperatorJoinedServices,
) error {
	return k.operatorJoinedServices.Set(ctx, operatorID, joinedServices)
}

// AddServiceToOperator adds the given service to the list of services joined by
// the operator with the given ID
func (k *Keeper) AddServiceToOperator(ctx sdk.Context, operatorID uint32, serviceID uint32) error {
	joinedServices, err := k.GetOperatorJoinedServices(ctx, operatorID)
	if err != nil {
		return err
	}

	err = joinedServices.Add(serviceID)
	if err != nil {
		return sdkerrors.Wrap(types.ErrServiceAlreadyJoinedByOperator, err.Error())
	}

	return k.SaveOperatorJoinedServices(ctx, operatorID, joinedServices)
}

// RemoveServiceFromOperator removes the given service from the list of services joined by
// the operator with the given ID
func (k *Keeper) RemoveServiceFromOperator(ctx sdk.Context, operatorID uint32, serviceID uint32) error {
	// Get the operator's joined services
	joinedServices, err := k.GetOperatorJoinedServices(ctx, operatorID)
	if err != nil {
		return err
	}

	// Try to remove the service
	removed := joinedServices.Remove(serviceID)
	if !removed {
		return types.ErrServiceNotJoinedByOperator
	}

	return k.SaveOperatorJoinedServices(ctx, operatorID, joinedServices)
}

// --------------------------------------------------------------------------------------------------------------------

// GetOperatorDelegation retrieves the delegation for the given user and operator
// If the delegation does not exist, false is returned instead
func (k *Keeper) GetOperatorDelegation(ctx sdk.Context, operatorID uint32, userAddress string) (types.Delegation, bool) {
	store := ctx.KVStore(k.storeKey)
	delegationBz := store.Get(types.UserOperatorDelegationStoreKey(userAddress, operatorID))
	if delegationBz == nil {
		return types.Delegation{}, false
	}

	return types.MustUnmarshalDelegation(k.cdc, delegationBz), true
}

// AddOperatorTokensAndShares adds the given amount of tokens to the operator and returns the added shares
func (k *Keeper) AddOperatorTokensAndShares(
	ctx sdk.Context, operator operatorstypes.Operator, tokensToAdd sdk.Coins,
) (operatorOut operatorstypes.Operator, addedShares sdk.DecCoins, err error) {
	// Update the operator tokens and shares and get the added shares
	operator, addedShares = operator.AddTokensFromDelegation(tokensToAdd)

	// Save the operator
	err = k.operatorsKeeper.SaveOperator(ctx, operator)
	return operator, addedShares, err
}

// RemoveOperatorDelegation removes the given operator delegation from the store
func (k *Keeper) RemoveOperatorDelegation(ctx sdk.Context, delegation types.Delegation) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.UserOperatorDelegationStoreKey(delegation.UserAddress, delegation.TargetID))
	store.Delete(types.DelegationByOperatorIDStoreKey(delegation.TargetID, delegation.UserAddress))
}

// --------------------------------------------------------------------------------------------------------------------

// DelegateToOperator sends the given amount to the operator account and saves the delegation for the given user
func (k *Keeper) DelegateToOperator(ctx sdk.Context, operatorID uint32, amount sdk.Coins, delegator string) (sdk.DecCoins, error) {
	// Get the operator
	operator, found := k.operatorsKeeper.GetOperator(ctx, operatorID)
	if !found {
		return sdk.NewDecCoins(), operatorstypes.ErrOperatorNotFound
	}

	// MAke sure the operator is active
	if !operator.IsActive() {
		return sdk.NewDecCoins(), operatorstypes.ErrOperatorNotActive
	}

	return k.PerformDelegation(ctx, types.DelegationData{
		Amount:          amount,
		Delegator:       delegator,
		Target:          &operator,
		BuildDelegation: types.NewOperatorDelegation,
		UpdateDelegation: func(ctx sdk.Context, delegation types.Delegation) (newShares sdk.DecCoins, err error) {
			// Calculate the new shares and add the tokens to the operator
			_, newShares, err = k.AddOperatorTokensAndShares(ctx, operator, amount)
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
			BeforeDelegationSharesModified: k.BeforeOperatorDelegationSharesModified,
			BeforeDelegationCreated:        k.BeforeOperatorDelegationCreated,
			AfterDelegationModified:        k.AfterOperatorDelegationModified,
		},
	})
}

// --------------------------------------------------------------------------------------------------------------------

// GetOperatorUnbondingDelegation returns the unbonding delegation for the given delegator address and operator id.
// If no unbonding delegation is found, false is returned instead.
func (k *Keeper) GetOperatorUnbondingDelegation(ctx sdk.Context, operatorID uint32, delegatorAddress string) (types.UnbondingDelegation, bool) {
	store := ctx.KVStore(k.storeKey)
	ubdBz := store.Get(types.UserOperatorUnbondingDelegationKey(delegatorAddress, operatorID))
	if ubdBz == nil {
		return types.UnbondingDelegation{}, false
	}

	return types.MustUnmarshalUnbondingDelegation(k.cdc, ubdBz), true
}

// UndelegateFromOperator removes the given amount from the operator account and saves the
// unbonding delegation for the given user
func (k *Keeper) UndelegateFromOperator(ctx sdk.Context, operatorID uint32, amount sdk.Coins, delegator string) (time.Time, error) {
	// Find the operator
	operator, found := k.operatorsKeeper.GetOperator(ctx, operatorID)
	if !found {
		return time.Time{}, operatorstypes.ErrOperatorNotFound
	}

	// Get the shares
	shares, err := k.ValidateUnbondAmount(ctx, delegator, &operator, amount)
	if err != nil {
		return time.Time{}, err
	}

	return k.PerformUndelegation(ctx, types.UndelegationData{
		Amount:                   amount,
		Delegator:                delegator,
		Target:                   &operator,
		BuildUnbondingDelegation: types.NewOperatorUnbondingDelegation,
		Hooks: types.DelegationHooks{
			BeforeDelegationSharesModified: k.BeforeOperatorDelegationSharesModified,
			BeforeDelegationCreated:        k.BeforeOperatorDelegationCreated,
			AfterDelegationModified:        k.AfterOperatorDelegationModified,
			BeforeDelegationRemoved:        k.BeforeOperatorDelegationRemoved,
		},
		Shares: shares,
	})
}
