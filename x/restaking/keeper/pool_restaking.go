package keeper

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

// SavePoolDelegation stores the given pool delegation in the store
func (k *Keeper) SavePoolDelegation(ctx sdk.Context, delegation types.PoolDelegation) {
	store := ctx.KVStore(k.storeKey)

	// Marshal and store the delegation
	delegationBz := types.MustMarshalPoolDelegation(k.cdc, delegation)
	store.Set(types.UserPoolDelegationStoreKey(delegation.UserAddress, delegation.PoolID), delegationBz)

	// Store the delegation in the delegations by pool ID store
	store.Set(types.DelegationByPoolIDStoreKey(delegation.PoolID, delegation.UserAddress), []byte{})
}

// GetPoolDelegation retrieves the delegation for the given user and pool
// If the delegation does not exist, false is returned instead
func (k *Keeper) GetPoolDelegation(ctx sdk.Context, poolID uint32, userAddress string) (types.PoolDelegation, bool) {
	// Get the delegation amount from the store
	store := ctx.KVStore(k.storeKey)
	delegationAmountBz := store.Get(types.UserPoolDelegationStoreKey(userAddress, poolID))
	if delegationAmountBz == nil {
		return types.PoolDelegation{}, false
	}

	// Parse the delegation amount
	return types.MustUnmarshalPoolDelegation(k.cdc, delegationAmountBz), true
}

// AddPoolTokensAndShares adds the given amount of tokens to the pool and returns the added shares
func (k *Keeper) AddPoolTokensAndShares(
	ctx sdk.Context, pool poolstypes.Pool, tokensToAdd sdkmath.Int,
) (poolOut poolstypes.Pool, addedShares sdkmath.LegacyDec, err error) {

	// Update the pool tokens and shares and get the added shares
	pool, addedShares = pool.AddTokensFromDelegation(tokensToAdd)

	// Save the pool
	err = k.poolsKeeper.SavePool(ctx, pool)
	return pool, addedShares, err
}

// --------------------------------------------------------------------------------------------------------------------

// DelegateToPool sends the given amount to the pool account and saves the delegation for the given user
func (k *Keeper) DelegateToPool(ctx sdk.Context, amount sdk.Coin, delegator string) (sdk.DecCoins, error) {
	// Get or create the pool for the given amount denom
	pool, err := k.poolsKeeper.CreateOrGetPoolByDenom(ctx, amount.Denom)
	if err != nil {
		return sdk.NewDecCoins(), err
	}

	// Get the amount to be bonded
	coins := sdk.NewCoins(sdk.NewCoin(pool.Denom, amount.Amount))

	return k.PerformDelegation(ctx, types.DelegationData{
		Amount:    coins,
		Delegator: delegator,
		Receiver:  &pool,
		GetDelegation: func(ctx sdk.Context, receiverID uint32, delegator string) (types.Delegation, bool) {
			return k.GetPoolDelegation(ctx, receiverID, delegator)
		},
		BuildDelegation: func(receiverID uint32, delegator string) types.Delegation {
			return types.NewPoolDelegation(receiverID, delegator, sdkmath.LegacyZeroDec())
		},
		UpdateDelegation: func(ctx sdk.Context, delegation types.Delegation) (sdk.DecCoins, error) {
			// Calculate the new shares and add the tokens to the pool
			_, newShares, err := k.AddPoolTokensAndShares(ctx, pool, amount.Amount)
			if err != nil {
				return nil, err
			}

			// Update the delegation shares
			poolDelegation, ok := delegation.(types.PoolDelegation)
			if !ok {
				return nil, fmt.Errorf("invalid delegation type: %T", delegation)
			}
			poolDelegation.Shares = poolDelegation.Shares.Add(newShares)

			// Store the updated delegation
			k.SavePoolDelegation(ctx, poolDelegation)

			return sdk.NewDecCoins(sdk.NewDecCoinFromDec(amount.Denom, newShares)), err
		},
		Hooks: types.DelegationHooks{
			BeforeDelegationSharesModified: k.BeforePoolDelegationSharesModified,
			BeforeDelegationCreated:        k.BeforePoolDelegationCreated,
			AfterDelegationModified:        k.AfterPoolDelegationModified,
		},
	})
}
