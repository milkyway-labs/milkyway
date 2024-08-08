package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

// SaveOperatorParams stored the given params for the given operator
func (k *Keeper) SaveOperatorParams(ctx sdk.Context, operatorID uint32, params types.OperatorParams) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.OperatorParamsStoreKey(operatorID), k.cdc.MustMarshal(&params))
}

// GetOperatorParams returns the params for the given operator, if any.
// If not params are found, false is returned instead.
func (k *Keeper) GetOperatorParams(ctx sdk.Context, operatorID uint32) (params types.OperatorParams, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.OperatorParamsStoreKey(operatorID))
	if bz == nil {
		return params, false
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params, true
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
	k.operatorsKeeper.SaveOperator(ctx, operator)
	return operator, addedShares, nil
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
func (k *Keeper) GetOperatorUnbondingDelegation(ctx sdk.Context, delegatorAddress string, operatorID uint32) (types.UnbondingDelegation, bool) {
	store := ctx.KVStore(k.storeKey)
	ubdBz := store.Get(types.UserOperatorUnbondingDelegationKey(delegatorAddress, operatorID))
	if ubdBz == nil {
		return types.UnbondingDelegation{}, false
	}

	return types.MustUnmarshalUnbondingDelegation(k.cdc, ubdBz), true
}
