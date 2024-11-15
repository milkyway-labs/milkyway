package keeper

import (
	"time"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

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
	ctx sdk.Context, pool poolstypes.Pool, tokensToAdd sdk.Coin,
) (poolOut poolstypes.Pool, addedShares sdk.DecCoin, err error) {
	// Update the pool tokens and shares and get the added shares
	pool, addedShares, err = pool.AddTokensFromDelegation(tokensToAdd)
	if err != nil {
		return pool, sdk.DecCoin{}, err
	}

	// Save the pool
	err = k.poolsKeeper.SavePool(ctx, pool)
	return pool, addedShares, err
}

// RemovePoolDelegation removes the given pool delegation from the store
func (k *Keeper) RemovePoolDelegation(ctx sdk.Context, delegation types.Delegation) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.UserPoolDelegationStoreKey(delegation.UserAddress, delegation.TargetID))
	store.Delete(types.DelegationByPoolIDStoreKey(delegation.TargetID, delegation.UserAddress))
}

// DelegateToPool sends the given amount to the pool account and saves the delegation for the given user
func (k *Keeper) DelegateToPool(ctx sdk.Context, amount sdk.Coin, delegator string) (sdk.DecCoins, error) {
	// Ensure the provided amount can be restaked
	isRestakable := k.IsDenomRestakable(ctx, amount.Denom)
	if !isRestakable {
		return sdk.NewDecCoins(), errors.Wrapf(types.ErrDenomNotRestakable, "%s cannot be restaked", amount.Denom)
	}

	// Get or create the pool for the given amount denom
	pool, err := k.poolsKeeper.CreateOrGetPoolByDenom(ctx, amount.Denom)
	if err != nil {
		return sdk.NewDecCoins(), err
	}

	// Get the amount to be bonded
	coins := sdk.NewCoins(sdk.NewCoin(pool.Denom, amount.Amount))

	return k.PerformDelegation(ctx, types.DelegationData{
		Amount:          coins,
		Delegator:       delegator,
		Target:          &pool,
		BuildDelegation: types.NewPoolDelegation,
		UpdateDelegation: func(ctx sdk.Context, delegation types.Delegation) (sdk.DecCoins, error) {
			// Calculate the new shares and add the tokens to the pool
			_, newShares, err := k.AddPoolTokensAndShares(ctx, pool, amount)
			if err != nil {
				return nil, err
			}

			// Update the delegation shares
			delegation.Shares = delegation.Shares.Add(newShares)

			// Store the updated delegation
			err = k.SetDelegation(ctx, delegation)
			if err != nil {
				return nil, err
			}

			return sdk.NewDecCoins(newShares), err
		},
		Hooks: types.DelegationHooks{
			BeforeDelegationSharesModified: k.BeforePoolDelegationSharesModified,
			BeforeDelegationCreated:        k.BeforePoolDelegationCreated,
			AfterDelegationModified:        k.AfterPoolDelegationModified,
		},
	})
}

// --------------------------------------------------------------------------------------------------------------------

// GetPoolUnbondingDelegation returns the unbonding delegation for the given delegator address and pool id.
// If no unbonding delegation is found, false is returned instead.
func (k *Keeper) GetPoolUnbondingDelegation(ctx sdk.Context, poolID uint32, delegator string) (types.UnbondingDelegation, bool) {
	store := ctx.KVStore(k.storeKey)
	ubdBz := store.Get(types.UserPoolUnbondingDelegationKey(delegator, poolID))
	if ubdBz == nil {
		return types.UnbondingDelegation{}, false
	}

	return types.MustUnmarshalUnbondingDelegation(k.cdc, ubdBz), true
}

// UndelegateFromPool removes the given amount from the pool account and saves the
// unbonding delegation for the given user
func (k *Keeper) UndelegateFromPool(ctx sdk.Context, amount sdk.Coin, delegator string) (time.Time, error) {
	// Find the pool
	pool, found, err := k.poolsKeeper.GetPoolByDenom(ctx, amount.Denom)
	if err != nil {
		return time.Time{}, err
	}

	if !found {
		return time.Time{}, poolstypes.ErrPoolNotFound
	}

	// Get the undelegation amount
	undelegationAmount := sdk.NewCoins(amount)

	// Get the shares
	shares, err := k.ValidateUnbondAmount(ctx, delegator, &pool, undelegationAmount)
	if err != nil {
		return time.Time{}, err
	}

	return k.PerformUndelegation(ctx, types.UndelegationData{
		Amount:                   undelegationAmount,
		Delegator:                delegator,
		Target:                   &pool,
		BuildUnbondingDelegation: types.NewPoolUnbondingDelegation,
		Hooks: types.DelegationHooks{
			BeforeDelegationSharesModified: k.BeforePoolDelegationSharesModified,
			BeforeDelegationCreated:        k.BeforePoolDelegationCreated,
			AfterDelegationModified:        k.AfterPoolDelegationModified,
			BeforeDelegationRemoved:        k.BeforePoolDelegationRemoved,
		},
		Shares: shares,
	})
}
