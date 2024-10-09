package keeper

import (
	"encoding/binary"
	"time"

	"cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/utils"
	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

// IncrementUnbondingID increments and returns a unique ID for an unbonding operation
func (k *Keeper) IncrementUnbondingID(ctx sdk.Context) (unbondingID uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.UnbondingIDKey)
	if bz != nil {
		unbondingID = binary.BigEndian.Uint64(bz)
	}

	// Increment the unbonding id
	unbondingID++

	// Convert back into bytes for storage
	bz = make([]byte, 8)
	binary.BigEndian.PutUint64(bz, unbondingID)

	// Store the new unbonding id
	store.Set(types.UnbondingIDKey, bz)

	return unbondingID
}

// DeleteUnbondingIndex removes a mapping from UnbondingId to unbonding operation
func (k *Keeper) DeleteUnbondingIndex(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetUnbondingIndexKey(id))
}

// ValidateUnbondAmount validates that a given unbond or redelegation amount is valid based on upon the
// converted shares. If the amount is valid, the total amount of respective shares is returned,
// otherwise an error is returned.
func (k *Keeper) ValidateUnbondAmount(
	ctx sdk.Context, delAddr string, target types.DelegationTarget, amt sdk.Coins,
) (shares sdk.DecCoins, err error) {
	// Get the delegation
	delegation, found := k.GetDelegationForTarget(ctx, target, delAddr)
	if !found {
		return shares, types.ErrDelegationNotFound
	}
	delShares := delegation.Shares

	shares, err = target.SharesFromTokens(amt)
	if err != nil {
		return shares, err
	}

	sharesTruncated, err := target.SharesFromTokensTruncated(amt)
	if err != nil {
		return shares, err
	}

	// Ensure that the shares to be unbonded are not greater than the shares that the delegator has
	if utils.IsAnyGT(sharesTruncated, delShares) {
		return shares, types.ErrInvalidShares
	}

	// Cap the shares at the delegation's shares. Shares being greater could occur
	// due to rounding, however we don't want to truncate the shares or take the
	// minimum because we want to allow for the full withdraw of shares from a
	// delegation.
	if utils.IsAnyGT(shares, delShares) {
		shares = delShares
	}

	return shares, nil
}

// Unbond unbonds a particular delegation and perform associated store operations.
func (k *Keeper) Unbond(ctx sdk.Context, data types.UndelegationData) (amount sdk.Coins, err error) {
	// Check if a delegation object exists in the store
	delegation, found := k.GetDelegationForTarget(ctx, data.Target, data.Delegator)
	if !found {
		return sdk.NewCoins(), types.ErrDelegationNotFound
	}

	// Call the before-delegation-modified hook
	err = data.Hooks.BeforeDelegationSharesModified(ctx, data.Target.GetID(), data.Delegator)
	if err != nil {
		return amount, err
	}

	// Ensure that we have enough shares to remove
	if utils.IsAnyLT(delegation.Shares, data.Shares) {
		return amount, errors.Wrap(types.ErrNotEnoughDelegationShares, delegation.Shares.String())
	}

	// Subtract shares from delegation
	delegation.Shares = delegation.Shares.Sub(data.Shares)

	if delegation.Shares.IsZero() {
		// Call the before delegation removed hook
		err = data.Hooks.BeforeDelegationRemoved(ctx, data.Target.GetID(), data.Delegator)
		if err != nil {
			return amount, err
		}

		// Remove the delegation
		k.RemoveDelegation(ctx, delegation)
	} else {
		// Store the updated delegation
		err = k.SetDelegation(ctx, delegation)
		if err != nil {
			return nil, err
		}

		// Call the after delegation modification hook
		err = data.Hooks.AfterDelegationModified(ctx, data.Target.GetID(), data.Delegator)
		if err != nil {
			return amount, err
		}
	}

	// Remove the shares and coins from the validator
	return k.RemoveTargetTokensAndShares(ctx, data)
}

// RemoveTargetTokensAndShares removes the given amount of tokens and shares from the target.
func (k *Keeper) RemoveTargetTokensAndShares(ctx sdk.Context, data types.UndelegationData) (sdk.Coins, error) {
	var issuedTokensAmount sdk.Coins

	switch target := data.Target.(type) {
	case *operatorstypes.Operator:
		operator, amount := target.RemoveDelShares(data.Shares)
		if err := k.operatorsKeeper.SaveOperator(ctx, operator); err != nil {
			return nil, err
		}
		issuedTokensAmount = amount
	case *servicestypes.Service:
		service, amount := target.RemoveDelShares(data.Shares)
		if err := k.servicesKeeper.SaveService(ctx, service); err != nil {
			return nil, err
		}
		issuedTokensAmount = amount
	case *poolstypes.Pool:
		pool, amount, err := target.RemoveDelShares(data.Shares)
		if err != nil {
			return nil, err
		}

		err = k.poolsKeeper.SavePool(ctx, pool)
		if err != nil {
			return nil, err
		}
		issuedTokensAmount = amount
	}

	return issuedTokensAmount, nil
}

// SetUnbondingDelegationEntry adds an entry to the unbonding delegation at
// the given addresses. It creates the unbonding delegation if it does not exist.
func (k *Keeper) SetUnbondingDelegationEntry(
	ctx sdk.Context, data types.UndelegationData, creationHeight int64, minTime time.Time, balance sdk.Coins,
) (types.UnbondingDelegation, error) {
	// Get the ID of the next unbonding delegation entry
	id := k.IncrementUnbondingID(ctx)

	// Either get the existing unbonding delegation, or create a new one
	ubdType, err := types.GetDelegationTypeFromTarget(data.Target)
	if err != nil {
		return types.UnbondingDelegation{}, err
	}

	isNewUbdEntry := true
	ubd, found := k.GetUnbondingDelegation(ctx, data.Delegator, ubdType, data.Target.GetID())
	if found {
		isNewUbdEntry = ubd.AddEntry(creationHeight, minTime, balance, id)
	} else {
		ubd = data.BuildUnbondingDelegation(data.Delegator, data.Target.GetID(), creationHeight, minTime, balance, id)
	}

	unbondingDelegationKey, err := k.SetUnbondingDelegation(ctx, ubd)
	if err != nil {
		return types.UnbondingDelegation{}, err
	}

	// Only call the hook for new entries since calls to AfterUnbondingInitiated are not idempotent
	if isNewUbdEntry {
		// Add to the UBDByUnbondingOp index to look up the UBD by the UBDE ID
		k.SetUnbondingDelegationByUnbondingID(ctx, ubd, unbondingDelegationKey, id)

		// Call the hook after the unbonding has been initiated
		err = k.AfterUnbondingInitiated(ctx, id)
		if err != nil {
			k.Logger(ctx).Error("failed to call after unbonding initiated hook", "error", err)
		}
	}

	return ubd, nil
}

// CompleteUnbonding completes the unbonding of all mature entries in the
// retrieved unbonding delegation object and returns the total unbonding balance
// or an error upon failure.
func (k *Keeper) CompleteUnbonding(ctx sdk.Context, data types.DTData) (sdk.Coins, error) {
	// Get the unbonding delegation entry
	ubd, found := k.GetUnbondingDelegation(ctx, data.DelegatorAddress, data.UnbondingDelegationType, data.TargetID)
	if !found {
		return nil, types.ErrNoUnbondingDelegation
	}

	// Get the target of the unbonding delegation
	target, err := k.getUnbondingDelegationTarget(ctx, ubd)
	if err != nil {
		return nil, err
	}

	// Get the address of the target
	targetAddress, err := sdk.AccAddressFromBech32(target.GetAddress())
	if err != nil {
		return nil, err
	}

	// Get the address of the delegator
	delegatorAddress, err := sdk.AccAddressFromBech32(ubd.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	balances := sdk.NewCoins()
	ctxTime := ctx.BlockHeader().Time

	// Loop through all the entries and complete unbonding mature entries
	for i := 0; i < len(ubd.Entries); i++ {
		entry := ubd.Entries[i]
		if entry.IsMature(ctxTime) {
			// Remove the entry
			ubd.RemoveEntry(int64(i))
			i--

			// Delete the index
			k.DeleteUnbondingIndex(ctx, entry.UnbondingID)

			// Track undelegation only when remaining or truncated shares are non-zero
			if !entry.Balance.IsZero() {
				amount := entry.Balance

				// Send the coins back to the delegator
				err = k.bankKeeper.SendCoins(ctx, targetAddress, delegatorAddress, amount)
				if err != nil {
					return nil, err
				}

				balances = balances.Add(amount...)
			}
		}
	}

	// Set the unbonding delegation or remove it if there are no more entries
	if len(ubd.Entries) == 0 {
		err = k.RemoveUnbondingDelegation(ctx, ubd)
		if err != nil {
			return nil, err
		}
	} else {
		_, err = k.SetUnbondingDelegation(ctx, ubd)
		if err != nil {
			return nil, err
		}
	}

	return balances, nil
}

// --------------------------------------------------------------------------------------------------------------------
// --- Unbonding queue operations
// --------------------------------------------------------------------------------------------------------------------

// GetUBDQueueTimeSlice gets a specific unbonding queue timeslice. A timeslice
// is a slice of DVPairs corresponding to unbonding delegations that expire at a
// certain time.
func (k *Keeper) GetUBDQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (dvPairs []types.DTData) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetUnbondingDelegationTimeKey(timestamp))
	if bz == nil {
		return []types.DTData{}
	}

	pairs := types.DTDataList{}
	k.cdc.MustUnmarshal(bz, &pairs)

	return pairs.Data
}

// SetUBDQueueTimeSlice sets a specific unbonding queue timeslice.
func (k *Keeper) SetUBDQueueTimeSlice(ctx sdk.Context, timestamp time.Time, keys []types.DTData) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&types.DTDataList{Data: keys})
	store.Set(types.GetUnbondingDelegationTimeKey(timestamp), bz)
}

// InsertUBDQueue inserts an unbonding delegation to the appropriate timeslice
// in the unbonding queue.
func (k *Keeper) InsertUBDQueue(ctx sdk.Context, ubd types.UnbondingDelegation, completionTime time.Time) {
	dvPair := types.DTData{
		UnbondingDelegationType: ubd.Type,
		DelegatorAddress:        ubd.DelegatorAddress,
		TargetID:                ubd.TargetID,
	}

	timeSlice := k.GetUBDQueueTimeSlice(ctx, completionTime)
	if len(timeSlice) == 0 {
		k.SetUBDQueueTimeSlice(ctx, completionTime, []types.DTData{dvPair})
	} else {
		timeSlice = append(timeSlice, dvPair)
		k.SetUBDQueueTimeSlice(ctx, completionTime, timeSlice)
	}
}

// UBDQueueIterator returns all the unbonding queue timeslices from time 0 until endTime.
func (k *Keeper) UBDQueueIterator(ctx sdk.Context, endTime time.Time) storetypes.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.UnbondingQueueKey, storetypes.InclusiveEndBytes(types.GetUnbondingDelegationTimeKey(endTime)))
}

// DequeueAllMatureUBDQueue returns a concatenated list of all the timeslices inclusively previous to
// currTime, and deletes the timeslices from the queue.
func (k *Keeper) DequeueAllMatureUBDQueue(ctx sdk.Context, currTime time.Time) (matureUnbonds []types.DTData) {
	store := ctx.KVStore(k.storeKey)

	// Get an iterator for all timeslices from time 0 until the current BlockHeader time
	unbondingTimesliceIterator := k.UBDQueueIterator(ctx, currTime)
	defer unbondingTimesliceIterator.Close()

	for ; unbondingTimesliceIterator.Valid(); unbondingTimesliceIterator.Next() {
		timeslice := types.DTDataList{}
		value := unbondingTimesliceIterator.Value()
		k.cdc.MustUnmarshal(value, &timeslice)

		matureUnbonds = append(matureUnbonds, timeslice.Data...)

		store.Delete(unbondingTimesliceIterator.Key())
	}

	return matureUnbonds
}
