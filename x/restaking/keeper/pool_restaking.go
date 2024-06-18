package keeper

import (
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
func (k *Keeper) DelegateToPool(ctx sdk.Context, amount sdk.Coin, delegator string) (sdkmath.LegacyDec, error) {
	// Get or create the pool for the given amount denom
	pool, err := k.poolsKeeper.CreateOrGetPoolByDenom(ctx, amount.Denom)
	if err != nil {
		return sdkmath.LegacyZeroDec(), err
	}

	// In some situations, the exchange rate becomes invalid, e.g. if
	// Pool loses all tokens due to slashing. In this case,
	// make all future delegations invalid.
	if pool.InvalidExRate() {
		return sdkmath.LegacyZeroDec(), types.ErrDelegatorShareExRateInvalid
	}

	// Get or create the delegation object and call the appropriate hook if present
	delegation, found := k.GetPoolDelegation(ctx, pool.ID, delegator)

	if found {
		// Delegation was found
		err = k.BeforePoolDelegationSharesModified(ctx, pool.ID, delegator)
		if err != nil {
			return sdkmath.LegacyZeroDec(), err
		}
	} else {
		// Delegation was not found
		delegation = types.NewPoolDelegation(pool.ID, delegator, sdkmath.LegacyZeroDec())
		err = k.BeforePoolDelegationCreated(ctx, pool.ID, delegator)
		if err != nil {
			return sdkmath.LegacyZeroDec(), err
		}
	}

	// Convert the addresses to sdk.AccAddress
	delegatorAddress, err := k.accountKeeper.AddressCodec().StringToBytes(delegator)
	if err != nil {
		return sdkmath.LegacyZeroDec(), err
	}
	poolAddress, err := k.accountKeeper.AddressCodec().StringToBytes(pool.Address)
	if err != nil {
		return sdkmath.LegacyZeroDec(), err
	}

	// Get the bond amount
	bondAmount := amount.Amount

	// Send the funds to the pool account
	coins := sdk.NewCoins(sdk.NewCoin(pool.Denom, bondAmount))
	err = k.bankKeeper.SendCoins(ctx, delegatorAddress, poolAddress, coins)
	if err != nil {
		return sdkmath.LegacyDec{}, err
	}

	// Calculate the new shares and add the tokens to the pool
	_, newShares, err := k.AddPoolTokensAndShares(ctx, pool, bondAmount)
	if err != nil {
		return newShares, err
	}

	// Update delegation
	delegation.Shares = delegation.Shares.Add(newShares)
	k.SavePoolDelegation(ctx, delegation)

	// Call the after-modification hook
	err = k.AfterPoolDelegationModified(ctx, pool.ID, delegator)
	if err != nil {
		return newShares, err
	}

	return newShares, nil
}
