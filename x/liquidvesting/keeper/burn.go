package keeper

import (
	"context"
	"fmt"
	"slices"
	"time"

	"cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v4/x/liquidvesting/types"
)

// IsBurner tells if a user have the permissions to burn tokens
// from a user's balance.
func (k *Keeper) IsBurner(ctx context.Context, user sdk.AccAddress) (bool, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return false, err
	}

	stringAddr, err := k.accountKeeper.AddressCodec().BytesToString(user)
	if err != nil {
		return false, err
	}

	return slices.Contains(params.Burners, stringAddr), nil
}

// BurnLockedRepresentation burns the locked staking representation
// from the user's balance.
// NOTE: If the coins are restaked they will be unstaked first.
func (k *Keeper) BurnLockedRepresentation(
	ctx context.Context,
	accAddress sdk.AccAddress,
	amount sdk.Coins,
) error {
	// Ensure that we are burning locked representations tokens
	for _, c := range amount {
		if !types.IsLockedRepresentationDenom(c.Denom) {
			return fmt.Errorf("invalid denom %s", c.Denom)
		}
	}

	// Get the user balance
	userBalance := k.bankKeeper.GetAllBalances(ctx, accAddress)

	liquidCoins := sdk.NewCoins()
	toUnbondCoins := sdk.NewCoins()
	for _, c := range amount {
		userBalanceOfC := userBalance.AmountOf(c.Denom)
		if userBalanceOfC.GTE(c.Amount) {
			liquidCoins = liquidCoins.Add(c)
		} else {
			// The user's balance of the coin c is lower than the amount to burn,
			// consider it as to unbond
			liquidCoins = liquidCoins.Add(sdk.NewCoin(c.Denom, userBalanceOfC))
			toUnbondCoins = toUnbondCoins.Add(sdk.NewCoin(c.Denom, c.Amount.Sub(userBalanceOfC)))
		}
	}

	liquidCoinsIsZero := liquidCoins.IsZero()
	toUnbondCoinsIsZero := toUnbondCoins.IsZero()

	if liquidCoinsIsZero && toUnbondCoinsIsZero {
		return errors.Wrap(types.ErrInsufficientBalance, amount.String())
	}

	// The amount to burn is not in the user balance, check if we can remove that
	// amount from the user's delegations.
	if !toUnbondCoinsIsZero {
		completionTime, err := k.restakingKeeper.UnbondRestakedAssets(ctx, accAddress, toUnbondCoins)
		if err != nil {
			return err
		}

		stringAddr, err := k.accountKeeper.AddressCodec().BytesToString(accAddress)
		if err != nil {
			return err
		}

		// Store in the burn coins queue that we have to burn those coins once
		// they are undelegated.
		err = k.InsertBurnCoinsToUnbondingQueue(ctx, types.NewBurnCoins(stringAddr, completionTime, toUnbondCoins))
		if err != nil {
			return err
		}
	}

	if !liquidCoinsIsZero {
		// Burn the liquid coins
		err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, accAddress, types.ModuleName, liquidCoins)
		if err != nil {
			return err
		}

		err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, liquidCoins)
		if err != nil {
			return err
		}
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------
// --- BurnCoins queue operations
// --------------------------------------------------------------------------------------------------------------------

// GetBurnCoinsQueueTimeSlice gets a specific burn coins queue timeslice. A timeslice
// is a slice of BurnCoins corresponding to the BurnCoins that needs to be burned at a certain time.
func (k *Keeper) GetBurnCoinsQueueTimeSlice(ctx context.Context, timestamp time.Time) (dvPairs []types.BurnCoins, err error) {
	store := k.storeService.OpenKVStore(ctx)

	bz, err := store.Get(types.GetBurnCoinsQueueTimeKey(timestamp))
	if err != nil {
		return nil, err
	}

	if bz == nil {
		return []types.BurnCoins{}, nil
	}

	pairs := types.BurnCoinsList{}
	k.cdc.MustUnmarshal(bz, &pairs)

	return pairs.Data, nil
}

// SetBurnCoinsQueueTimeSlice sets a specific burn coins queue timeslice.
func (k *Keeper) SetBurnCoinsQueueTimeSlice(ctx context.Context, timestamp time.Time, keys []types.BurnCoins) error {
	store := k.storeService.OpenKVStore(ctx)
	bz := k.cdc.MustMarshal(&types.BurnCoinsList{Data: keys})
	return store.Set(types.GetBurnCoinsQueueTimeKey(timestamp), bz)
}

// InsertBurnCoinsToUnbondingQueue inserts an BurnCoin to the appropriate timeslice
// in the burn coins queue.
func (k *Keeper) InsertBurnCoinsToUnbondingQueue(ctx context.Context, burnCoins types.BurnCoins) error {
	// Get the existing list of coins to be burned
	timeSlice, err := k.GetBurnCoinsQueueTimeSlice(ctx, burnCoins.CompletionTime)
	if err != nil {
		return err
	}

	// Add the new coin
	timeSlice = append(timeSlice, burnCoins)
	return k.SetBurnCoinsQueueTimeSlice(ctx, burnCoins.CompletionTime, timeSlice)
}

// BurnCoinsUnbondingQueueIterator returns all the BurnCoins from time 0 until endTime.
func (k *Keeper) BurnCoinsUnbondingQueueIterator(ctx context.Context, endTime time.Time) (storetypes.Iterator, error) {
	store := k.storeService.OpenKVStore(ctx)
	return store.Iterator(types.BurnCoinsQueueKey, storetypes.InclusiveEndBytes(types.GetBurnCoinsQueueTimeKey(endTime)))
}

// IterateBurnCoinsUnbondingQueue iterates all the BurnCoins from time 0 until endTime.
func (k *Keeper) IterateBurnCoinsUnbondingQueue(ctx context.Context, endTime time.Time, iterF func(burnCoin types.BurnCoins) (bool, error)) error {
	iter, err := k.BurnCoinsUnbondingQueueIterator(ctx, endTime)
	if err != nil {
		return err
	}

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		timeslice := types.BurnCoinsList{}
		value := iter.Value()
		k.cdc.MustUnmarshal(value, &timeslice)
		for _, burnCoin := range timeslice.Data {
			stop, err := iterF(burnCoin)
			if stop || err != nil {
				return err
			}
		}
	}

	return nil
}

// GetUnbondedCoinsFromQueue returns a concatenated list of all the timeslices inclusively previous to
// endTime.
func (k *Keeper) GetUnbondedCoinsFromQueue(ctx context.Context, endTime time.Time) ([]types.BurnCoins, error) {
	var burnCoins []types.BurnCoins
	err := k.IterateBurnCoinsUnbondingQueue(ctx, endTime, func(burnCoin types.BurnCoins) (bool, error) {
		burnCoins = append(burnCoins, burnCoin)
		return false, nil
	})
	return burnCoins, err
}

// DequeueAllBurnCoinsFromUnbondingQueue returns a concatenated list of all the timeslices inclusively previous to
// currTime, and deletes the timeslices from the queue.
func (k *Keeper) DequeueAllBurnCoinsFromUnbondingQueue(ctx context.Context, currTime time.Time) (burnCoins []types.BurnCoins, err error) {
	store := k.storeService.OpenKVStore(ctx)

	// Get an iterator for all timeslices from time 0 until the current BlockHeader time
	iter, err := k.BurnCoinsUnbondingQueueIterator(ctx, currTime)
	if err != nil {
		return nil, err
	}

	var toDeleteKeys [][]byte
	for ; iter.Valid(); iter.Next() {
		timeslice := types.BurnCoinsList{}
		value := iter.Value()
		k.cdc.MustUnmarshal(value, &timeslice)

		burnCoins = append(burnCoins, timeslice.Data...)

		toDeleteKeys = append(toDeleteKeys, iter.Key())
	}

	// Close the iterator
	err = iter.Close()
	if err != nil {
		return nil, err
	}

	// Delete all the keys
	for _, key := range toDeleteKeys {
		err = store.Delete(key)
		if err != nil {
			return nil, err
		}
	}

	return burnCoins, nil
}

// GetAllBurnCoins returns all the coins that are scheduled to be burned.
func (k *Keeper) GetAllBurnCoins(ctx context.Context) []types.BurnCoins {
	store := k.storeService.OpenKVStore(ctx)

	iterator := storetypes.KVStorePrefixIterator(runtime.KVStoreAdapter(store), types.BurnCoinsQueueKey)
	defer iterator.Close()

	var burnCoins []types.BurnCoins
	for ; iterator.Valid(); iterator.Next() {
		var burnCoinsList types.BurnCoinsList
		val := iterator.Value()
		k.cdc.MustUnmarshal(val, &burnCoinsList)
		burnCoins = append(burnCoins, burnCoinsList.Data...)
	}

	return burnCoins
}
