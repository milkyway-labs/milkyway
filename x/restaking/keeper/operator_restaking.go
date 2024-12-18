package keeper

import (
	"context"
	"slices"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	operatorstypes "github.com/milkyway-labs/milkyway/v6/x/operators/types"
	"github.com/milkyway-labs/milkyway/v6/x/restaking/types"
)

// AddServiceToOperatorJoinedServices adds the given service to the list of services joined by
// the operator with the given ID
func (k *Keeper) AddServiceToOperatorJoinedServices(ctx context.Context, operatorID uint32, serviceID uint32) error {
	operatorServicePair := collections.Join(operatorID, serviceID)
	return k.operatorJoinedServices.Set(ctx, operatorServicePair, collections.NoValue{})
}

// RemoveServiceFromOperatorJoinedServices removes the given service from the list of services joined by
// the operator with the given ID
func (k *Keeper) RemoveServiceFromOperatorJoinedServices(ctx context.Context, operatorID uint32, serviceID uint32) error {
	operatorServicePair := collections.Join(operatorID, serviceID)
	return k.operatorJoinedServices.Remove(ctx, operatorServicePair)
}

// HasOperatorJoinedService returns whether the operator with the given ID has
// joined the provided service
func (k *Keeper) HasOperatorJoinedService(ctx context.Context, operatorID uint32, serviceID uint32) (bool, error) {
	operatorServicePair := collections.Join(operatorID, serviceID)
	return k.operatorJoinedServices.Has(ctx, operatorServicePair)
}

// --------------------------------------------------------------------------------------------------------------------

// GetOperatorDelegation retrieves the delegation for the given user and operator
// If the delegation does not exist, false is returned instead
func (k *Keeper) GetOperatorDelegation(ctx context.Context, operatorID uint32, userAddress string) (types.Delegation, bool, error) {
	store := k.storeService.OpenKVStore(ctx)
	delegationBz, err := store.Get(types.UserOperatorDelegationStoreKey(userAddress, operatorID))
	if err != nil {
		return types.Delegation{}, false, err
	}

	if delegationBz == nil {
		return types.Delegation{}, false, nil
	}

	return types.MustUnmarshalDelegation(k.cdc, delegationBz), true, nil
}

// AddOperatorTokensAndShares adds the given amount of tokens to the operator and returns the added shares
func (k *Keeper) AddOperatorTokensAndShares(
	ctx context.Context, operator operatorstypes.Operator, tokensToAdd sdk.Coins,
) (operatorOut operatorstypes.Operator, addedShares sdk.DecCoins, err error) {
	// Update the operator tokens and shares and get the added shares
	operator, addedShares = operator.AddTokensFromDelegation(tokensToAdd)

	// Save the operator
	err = k.operatorsKeeper.SaveOperator(ctx, operator)
	return operator, addedShares, err
}

// RemoveOperatorDelegation removes the given operator delegation from the store
func (k *Keeper) RemoveOperatorDelegation(ctx context.Context, delegation types.Delegation) error {
	store := k.storeService.OpenKVStore(ctx)

	err := store.Delete(types.UserOperatorDelegationStoreKey(delegation.UserAddress, delegation.TargetID))
	if err != nil {
		return err
	}

	return store.Delete(types.DelegationByOperatorIDStoreKey(delegation.TargetID, delegation.UserAddress))
}

// --------------------------------------------------------------------------------------------------------------------

// DelegateToOperator sends the given amount to the operator account and saves the delegation for the given user
func (k *Keeper) DelegateToOperator(ctx context.Context, operatorID uint32, amount sdk.Coins, delegator string) (sdk.DecCoins, error) {
	// Get the operator
	operator, err := k.operatorsKeeper.GetOperator(ctx, operatorID)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return sdk.NewDecCoins(), operatorstypes.ErrOperatorNotFound
		}
		return nil, err
	}

	restakableDenoms, err := k.GetRestakableDenoms(ctx)
	if err != nil {
		return nil, err
	}

	if len(restakableDenoms) > 0 {
		// Ensure the provided amount can be restaked
		for _, coin := range amount {
			isRestakable := slices.Contains(restakableDenoms, coin.Denom)
			if !isRestakable {
				return sdk.NewDecCoins(), errors.Wrapf(types.ErrDenomNotRestakable, "%s cannot be restaked", coin.Denom)
			}
		}
	}

	// Make sure the operator is active
	if !operator.IsActive() {
		return sdk.NewDecCoins(), operatorstypes.ErrOperatorNotActive
	}

	return k.PerformDelegation(ctx, types.DelegationData{
		Amount:          amount,
		Delegator:       delegator,
		Target:          operator,
		BuildDelegation: types.NewOperatorDelegation,
		UpdateDelegation: func(ctx context.Context, delegation types.Delegation) (newShares sdk.DecCoins, err error) {
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
func (k *Keeper) GetOperatorUnbondingDelegation(ctx context.Context, operatorID uint32, delegatorAddress string) (types.UnbondingDelegation, bool, error) {
	store := k.storeService.OpenKVStore(ctx)

	ubdBz, err := store.Get(types.UserOperatorUnbondingDelegationKey(delegatorAddress, operatorID))
	if err != nil {
		return types.UnbondingDelegation{}, false, err
	}

	if ubdBz == nil {
		return types.UnbondingDelegation{}, false, nil
	}

	return types.MustUnmarshalUnbondingDelegation(k.cdc, ubdBz), true, nil
}

// UndelegateFromOperator removes the given amount from the operator account and saves the
// unbonding delegation for the given user
func (k *Keeper) UndelegateFromOperator(ctx context.Context, operatorID uint32, amount sdk.Coins, delegator string) (time.Time, error) {
	// Find the operator
	operator, err := k.operatorsKeeper.GetOperator(ctx, operatorID)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return time.Time{}, operatorstypes.ErrOperatorNotFound
		}
		return time.Time{}, err
	}

	// Get the shares
	shares, err := k.ValidateUnbondAmount(ctx, delegator, operator, amount)
	if err != nil {
		return time.Time{}, err
	}

	return k.PerformUndelegation(ctx, types.UndelegationData{
		Amount:                   amount,
		Delegator:                delegator,
		Target:                   operator,
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
