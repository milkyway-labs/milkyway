package keeper

import (
	"fmt"
	"slices"
	"time"

	"cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
)

// IsBurner tells if a user have the permissions to burn tokens
// from a user's balance.
func (k *Keeper) IsBurner(ctx sdk.Context, user sdk.AccAddress) (bool, error) {
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

// BurnVestedRepresentation burns the vested staking representation
// from the user's balance.
// NOTE: If the coins are restaked they will be unstaked first.
func (k *Keeper) BurnVestedRepresentation(
	ctx sdk.Context,
	accAddress sdk.AccAddress,
	amount sdk.Coins,
) error {
	// Ensure that we are burning vested representations tokens
	for _, c := range amount {
		if !types.IsVestedRepresentationDenom(c.Denom) {
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
			// The user's balance of the coin c is lower then the amount to burn,
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
		k.InsertBurnCoinsQueue(ctx, types.NewBurnCoins(stringAddr, completionTime, toUnbondCoins))
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
func (k *Keeper) GetBurnCoinsQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (dvPairs []types.BurnCoins) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetBurnCoinsQueueTimeKey(timestamp))
	if bz == nil {
		return []types.BurnCoins{}
	}

	pairs := types.BurnCoinsList{}
	k.cdc.MustUnmarshal(bz, &pairs)

	return pairs.Data
}

// SetBurnCoinsQueueTimeSlice sets a specific burn coins queue timeslice.
func (k *Keeper) SetBurnCoinsQueueTimeSlice(ctx sdk.Context, timestamp time.Time, keys []types.BurnCoins) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&types.BurnCoinsList{Data: keys})
	store.Set(types.GetBurnCoinsQueueTimeKey(timestamp), bz)
}

// InsertBurnCoinsQueue inserts an BurnCoin to the appropriate timeslice
// in the burn coins queue.
func (k *Keeper) InsertBurnCoinsQueue(ctx sdk.Context, burnCoins types.BurnCoins) {
	// Get the existing list of coins to be burned
	timeSlice := k.GetBurnCoinsQueueTimeSlice(ctx, burnCoins.CompletionTime)

	// Add the new coin
	timeSlice = append(timeSlice, burnCoins)
	k.SetBurnCoinsQueueTimeSlice(ctx, burnCoins.CompletionTime, timeSlice)
}

// BurnCoinsQueueIterator returns all the BurnCoins from time 0 until endTime.
func (k *Keeper) BurnCoinsQueueIterator(ctx sdk.Context, endTime time.Time) storetypes.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.BurnCoinsQueueKey, storetypes.InclusiveEndBytes(types.GetBurnCoinsQueueTimeKey(endTime)))
}

// DequeueAllBurnCoinsQueue returns a concatenated list of all the timeslices inclusively previous to
// currTime, and deletes the timeslices from the queue.
func (k *Keeper) DequeueAllBurnCoinsQueue(ctx sdk.Context, currTime time.Time) (burnCoins []types.BurnCoins) {
	store := ctx.KVStore(k.storeKey)

	// Get an iterator for all timeslices from time 0 until the current BlockHeader time
	bcTimesliceIterator := k.BurnCoinsQueueIterator(ctx, currTime)
	defer bcTimesliceIterator.Close()

	for ; bcTimesliceIterator.Valid(); bcTimesliceIterator.Next() {
		timeslice := types.BurnCoinsList{}
		value := bcTimesliceIterator.Value()
		k.cdc.MustUnmarshal(value, &timeslice)

		burnCoins = append(burnCoins, timeslice.Data...)

		store.Delete(bcTimesliceIterator.Key())
	}

	return burnCoins
}

// GetAllBurnCoins returns all the coins that are scheduled to be burned.
func (k *Keeper) GetAllBurnCoins(ctx sdk.Context) []types.BurnCoins {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.BurnCoinsQueueKey)
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
