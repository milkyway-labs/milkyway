package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

// SaveOperatorDelegation stores the given operator delegation in the store
func (k *Keeper) SaveOperatorDelegation(ctx sdk.Context, delegation types.OperatorDelegation) {
	store := ctx.KVStore(k.storeKey)

	// Marshal and store the delegation
	delegationBz := types.MustMarshalOperatorDelegation(k.cdc, delegation)
	store.Set(types.UserOperatorDelegationStoreKey(delegation.UserAddress, delegation.OperatorID), delegationBz)

	// Store the delegation in the delegations by operator ID store
	store.Set(types.DelegationByOperatorIDStoreKey(delegation.OperatorID, delegation.UserAddress), []byte{})
}

// GetOperatorDelegation retrieves the delegation for the given user and operator
// If the delegation does not exist, false is returned instead
func (k *Keeper) GetOperatorDelegation(ctx sdk.Context, operatorID uint32, userAddress string) (types.OperatorDelegation, bool) {
	store := ctx.KVStore(k.storeKey)
	delegationBz := store.Get(types.UserOperatorDelegationStoreKey(userAddress, operatorID))
	if delegationBz == nil {
		return types.OperatorDelegation{}, false
	}

	return types.MustUnmarshalOperatorDelegation(k.cdc, delegationBz), true
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
		Amount:    amount,
		Delegator: delegator,
		Receiver:  &operator,
		GetDelegation: func(ctx sdk.Context, receiverID uint32, delegator string) (types.Delegation, bool) {
			return k.GetOperatorDelegation(ctx, receiverID, delegator)
		},
		BuildDelegation: func(receiverID uint32, delegator string) types.Delegation {
			return types.NewOperatorDelegation(receiverID, delegator, sdk.NewDecCoins())
		},
		UpdateDelegation: func(ctx sdk.Context, delegation types.Delegation) (newShares sdk.DecCoins, err error) {
			// Calculate the new shares and add the tokens to the operator
			_, newShares, err = k.AddOperatorTokensAndShares(ctx, operator, amount)
			if err != nil {
				return newShares, err
			}

			// Update the delegation shares
			operatorDelegation, ok := delegation.(types.OperatorDelegation)
			if !ok {
				return newShares, fmt.Errorf("invalid delegation type: %T", delegation)
			}
			operatorDelegation.Shares = operatorDelegation.Shares.Add(newShares...)

			// Store the updated delegation
			k.SaveOperatorDelegation(ctx, operatorDelegation)

			return newShares, err
		},
		Hooks: types.DelegationHooks{
			BeforeDelegationSharesModified: k.BeforeOperatorDelegationSharesModified,
			BeforeDelegationCreated:        k.BeforeOperatorDelegationCreated,
			AfterDelegationModified:        k.AfterOperatorDelegationModified,
		},
	})
}
