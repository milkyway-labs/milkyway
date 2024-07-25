package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

// SavePoolDelegation stores the given pool delegation in the store
func (k *Keeper) SavePoolDelegation(ctx sdk.Context, delegation types.Delegation) {
	store := ctx.KVStore(k.storeKey)

	// Marshal and store the delegation
	delegationBz := types.MustMarshalDelegation(k.cdc, delegation)
	store.Set(types.UserPoolDelegationStoreKey(delegation.UserAddress, delegation.TargetID), delegationBz)

	// Store the delegation in the delegations by pool ID store
	store.Set(types.DelegationByPoolIDStoreKey(delegation.TargetID, delegation.UserAddress), []byte{})
}

// GetPoolDelegation retrieves the delegation for the given user and pool
// If the delegation does not exist, false is returned instead
func (k *Keeper) GetPoolDelegation(ctx sdk.Context, poolID uint32, userAddress string) (types.Delegation, bool) {
	store := ctx.KVStore(k.storeKey)
	delegationBz := store.Get(types.UserPoolDelegationStoreKey(userAddress, poolID))
	if delegationBz == nil {
		return types.Delegation{}, false
	}

	return types.MustUnmarshalDelegation(k.cdc, delegationBz), true
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
		Target:    &pool,
		GetDelegation: func(ctx sdk.Context, receiverID uint32, delegator string) (types.Delegation, bool) {
			return k.GetPoolDelegation(ctx, receiverID, delegator)
		},
		BuildDelegation: types.NewPoolDelegation,
		UpdateDelegation: func(ctx sdk.Context, delegation types.Delegation) (sdk.DecCoins, error) {
			// Calculate the new shares and add the tokens to the pool
			_, newShares, err := k.AddPoolTokensAndShares(ctx, pool, amount.Amount)
			if err != nil {
				return nil, err
			}

			// Update the delegation shares
			newDecShares := sdk.NewDecCoinFromDec(pool.GetSharesDenom(amount.Denom), newShares)
			delegation.Shares = delegation.Shares.Add(newDecShares)

			// Store the updated delegation
			k.SavePoolDelegation(ctx, delegation)

			sharesDenom := pool.GetSharesDenom(amount.Denom)
			return sdk.NewDecCoins(sdk.NewDecCoinFromDec(sharesDenom, newShares)), err
		},
		Hooks: types.DelegationHooks{
			BeforeDelegationSharesModified: k.BeforePoolDelegationSharesModified,
			BeforeDelegationCreated:        k.BeforePoolDelegationCreated,
			AfterDelegationModified:        k.AfterPoolDelegationModified,
		},
	})
}
