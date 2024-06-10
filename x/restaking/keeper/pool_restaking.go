package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

// PerformPoolDelegation sends the given amount to the pool account and saves the delegation for the given user
func (k *Keeper) PerformPoolDelegation(ctx sdk.Context, amount sdk.Coin, delegator string) error {
	// Get or create the pool for the given amount denom
	pool, err := k.poolsKeeper.CreateOrGetPoolByDenom(ctx, amount.Denom)
	if err != nil {
		return err
	}

	// Send the funds to the pool account
	delegatorAddr, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		return err
	}
	poolAddr, err := sdk.AccAddressFromBech32(pool.AccountAddress)
	if err != nil {
		return err
	}
	err = k.bankKeeper.SendCoins(ctx, delegatorAddr, poolAddr, sdk.NewCoins(amount))
	if err != nil {
		return err
	}

	// Get the current delegation for the user
	delegationAmount, found, err := k.GetUserPoolDelegationAmount(ctx, pool.ID, delegator)
	if err != nil {
		return err
	}

	// If a delegation already exists, add the new amount to the existing one
	if found {
		amount = amount.AddAmount(delegationAmount.Amount)
	}

	// Save the new delegation amount
	k.SavePoolDelegation(ctx, pool.ID, amount, delegator)

	return nil
}

// SavePoolDelegation stores the given amount as the delegation for the given user in the given pool
func (k *Keeper) SavePoolDelegation(ctx sdk.Context, poolID uint32, amount sdk.Coin, delegator string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.UserPoolDelegationStoreKey(poolID, delegator), []byte(amount.String()))
}

// GetUserPoolDelegationAmount returns the delegation amount for the given user in the given pool
func (k *Keeper) GetUserPoolDelegationAmount(ctx sdk.Context, poolID uint32, userAddress string) (sdk.Coin, bool, error) {
	// Get the delegation amount from the store
	store := ctx.KVStore(k.storeKey)
	delegationAmountBz := store.Get(types.UserPoolDelegationStoreKey(poolID, userAddress))
	if delegationAmountBz == nil {
		return sdk.Coin{}, false, nil
	}

	// Parse the delegation amount
	amount, err := sdk.ParseCoinNormalized(string(delegationAmountBz))
	if err != nil {
		return sdk.Coin{}, false, err
	}

	return amount, true, nil
}
