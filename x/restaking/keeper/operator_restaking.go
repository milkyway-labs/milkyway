package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

// SaveOperatorParams stored the given params for the given operator
func (k *Keeper) SaveOperatorParams(ctx sdk.Context, operatorID uint32, params types.OperatorParams) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.OperatorParamsStoreKey(operatorID), k.cdc.MustMarshal(&params))

	// Store the operator params in the x/opeators module.
	// TODO: Once we have moved also the operator's joined services in a dedicated
	// collection the whole SaveOperatorParams method should be removed.
	err := k.operatorsKeeper.SaveOperatorParams(ctx, operatorID, operatorstypes.NewOperatorParams(params.CommissionRate))
	if err != nil {
		panic(err)
	}
}

// GetOperatorParams returns the params for the given operator, if any.
// If not params are found, false is returned instead.
func (k *Keeper) GetOperatorParams(ctx sdk.Context, operatorID uint32) (params types.OperatorParams) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.OperatorParamsStoreKey(operatorID))
	if bz == nil {
		params = types.DefaultOperatorParams()
	} else {
		k.cdc.MustUnmarshal(bz, &params)
	}

	// Get the commission rate from the x/opeators module.
	// TODO: Once we have moved also the operator's joined services in a dedicated
	// collection the whole GetOperatorParams method should be removed.
	operatorParams, err := k.operatorsKeeper.GetOperatorParams(ctx, operatorID)
	if err != nil {
		panic(err)
	}
	params.CommissionRate = operatorParams.CommissionRate

	return params
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
