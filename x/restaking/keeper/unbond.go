package keeper

import (
	"context"
	"encoding/binary"
	"time"

	"cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v8/utils"
	operatorstypes "github.com/milkyway-labs/milkyway/v8/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/v8/x/pools/types"
	"github.com/milkyway-labs/milkyway/v8/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v8/x/services/types"
)

// IncrementUnbondingID increments and returns a unique ID for an unbonding operation
func (k *Keeper) IncrementUnbondingID(ctx context.Context) (unbondingID uint64, err error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.UnbondingIDKey)
	if err != nil {
		return 0, err
	}

	if bz != nil {
		unbondingID = binary.BigEndian.Uint64(bz)
	}

	// Increment the unbonding id
	unbondingID++

	// Convert back into bytes for storage
	bz = make([]byte, 8)
	binary.BigEndian.PutUint64(bz, unbondingID)

	// Store the new unbonding id
	err = store.Set(types.UnbondingIDKey, bz)
	if err != nil {
		return 0, err
	}

	return unbondingID, nil
}

// DeleteUnbondingIndex removes a mapping from UnbondingId to unbonding operation
func (k *Keeper) DeleteUnbondingIndex(ctx context.Context, id uint64) error {
	store := k.storeService.OpenKVStore(ctx)
	return store.Delete(types.GetUnbondingIndexKey(id))
}

// ValidateUnbondAmount validates that a given unbond or redelegation amount is valid based on upon the
// converted shares. If the amount is valid, the total amount of respective shares is returned,
// otherwise an error is returned.
func (k *Keeper) ValidateUnbondAmount(
	ctx context.Context, delAddr string, target types.DelegationTarget, amt sdk.Coins,
) (shares sdk.DecCoins, err error) {
	// Get the delegation
	delegation, found, err := k.GetDelegationForTarget(ctx, target, delAddr)
	if err != nil {
		return nil, err
	}

	if !found {
		return shares, types.ErrDelegationNotFound
	}
	delShares := delegation.Shares

	shares, err = target.SharesFromTokens(amt)
	if err != nil {
		return shares, err
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
func (k *Keeper) Unbond(ctx context.Context, data types.UndelegationData) (amount sdk.Coins, err error) {
	// Check if a delegation object exists in the store
	delegation, found, err := k.GetDelegationForTarget(ctx, data.Target, data.Delegator)
	if err != nil {
		return nil, err
	}

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

	// Charge gas cost based on the number of denoms inside the unbonding delegation
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.GasMeter().ConsumeGas(types.BaseDelegationDenomCost*uint64(len(delegation.Shares)), "undelegation gas cost")

	// Subtract shares from delegation
	delegation.Shares = delegation.Shares.Sub(data.Shares)

	if delegation.Shares.IsZero() {
		// Call the before delegation removed hook
		err = data.Hooks.BeforeDelegationRemoved(ctx, data.Target.GetID(), data.Delegator)
		if err != nil {
			return amount, err
		}

		// Remove the delegation
		err = k.RemoveDelegation(ctx, delegation)
		if err != nil {
			return amount, err
		}
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
func (k *Keeper) RemoveTargetTokensAndShares(ctx context.Context, data types.UndelegationData) (sdk.Coins, error) {
	var issuedTokensAmount sdk.Coins

	switch target := data.Target.(type) {
	case operatorstypes.Operator:
		operator, amount := target.RemoveDelShares(data.Shares)
		if err := k.operatorsKeeper.SaveOperator(ctx, operator); err != nil {
			return nil, err
		}
		issuedTokensAmount = amount
	case servicestypes.Service:
		service, amount := target.RemoveDelShares(data.Shares)
		if err := k.servicesKeeper.SaveService(ctx, service); err != nil {
			return nil, err
		}
		issuedTokensAmount = amount
	case poolstypes.Pool:
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
	ctx context.Context, data types.UndelegationData, creationHeight int64, minTime time.Time, balance sdk.Coins,
) (types.UnbondingDelegation, error) {
	// Get the ID of the next unbonding delegation entry
	id, err := k.IncrementUnbondingID(ctx)
	if err != nil {
		return types.UnbondingDelegation{}, err
	}

	// Either get the existing unbonding delegation, or create a new one
	ubdType, err := types.GetDelegationTypeFromTarget(data.Target)
	if err != nil {
		return types.UnbondingDelegation{}, err
	}

	isNewUbdEntry := true
	ubd, found, err := k.GetUnbondingDelegation(ctx, data.Delegator, ubdType, data.Target.GetID())
	if err != nil {
		return types.UnbondingDelegation{}, err
	}

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
		err = k.SetUnbondingDelegationByUnbondingID(ctx, ubd, unbondingDelegationKey, id)
		if err != nil {
			return types.UnbondingDelegation{}, err
		}

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
func (k *Keeper) CompleteUnbonding(ctx context.Context, data types.DTData) (sdk.Coins, error) {
	// Get the unbonding delegation entry
	ubd, found, err := k.GetUnbondingDelegation(ctx, data.DelegatorAddress, data.UnbondingDelegationType, data.TargetID)
	if err != nil {
		return nil, err
	}

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

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	ctxTime := sdkCtx.BlockHeader().Time

	// Loop through all the entries and complete unbonding mature entries
	balances := sdk.NewCoins()
	for i := 0; i < len(ubd.Entries); i++ {
		entry := ubd.Entries[i]
		if entry.IsMature(ctxTime) {
			// Remove the entry
			ubd.RemoveEntry(int64(i))
			i--

			// Delete the index
			err = k.DeleteUnbondingIndex(ctx, entry.UnbondingID)
			if err != nil {
				return nil, err
			}

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
func (k *Keeper) GetUBDQueueTimeSlice(ctx context.Context, timestamp time.Time) (dvPairs []types.DTData, err error) {
	store := k.storeService.OpenKVStore(ctx)

	bz, err := store.Get(types.GetUnbondingDelegationTimeKey(timestamp))
	if err != nil {
		return nil, err
	}

	if bz == nil {
		return []types.DTData{}, nil
	}

	pairs := types.DTDataList{}
	k.cdc.MustUnmarshal(bz, &pairs)

	return pairs.Data, nil
}

// SetUBDQueueTimeSlice sets a specific unbonding queue timeslice.
func (k *Keeper) SetUBDQueueTimeSlice(ctx context.Context, timestamp time.Time, keys []types.DTData) error {
	store := k.storeService.OpenKVStore(ctx)
	bz := k.cdc.MustMarshal(&types.DTDataList{Data: keys})
	return store.Set(types.GetUnbondingDelegationTimeKey(timestamp), bz)
}

// InsertUBDQueue inserts an unbonding delegation to the appropriate timeslice
// in the unbonding queue.
func (k *Keeper) InsertUBDQueue(ctx context.Context, ubd types.UnbondingDelegation, completionTime time.Time) error {
	dvPair := types.DTData{
		UnbondingDelegationType: ubd.Type,
		DelegatorAddress:        ubd.DelegatorAddress,
		TargetID:                ubd.TargetID,
	}

	timeSlice, err := k.GetUBDQueueTimeSlice(ctx, completionTime)
	if err != nil {
		return err
	}

	if len(timeSlice) == 0 {
		return k.SetUBDQueueTimeSlice(ctx, completionTime, []types.DTData{dvPair})

	}

	timeSlice = append(timeSlice, dvPair)
	return k.SetUBDQueueTimeSlice(ctx, completionTime, timeSlice)
}

// UBDQueueIterator returns all the unbonding queue timeslices from time 0 until endTime.
func (k *Keeper) UBDQueueIterator(ctx context.Context, endTime time.Time) (storetypes.Iterator, error) {
	store := k.storeService.OpenKVStore(ctx)
	return store.Iterator(types.UnbondingQueueKey, storetypes.InclusiveEndBytes(types.GetUnbondingDelegationTimeKey(endTime)))
}

// DequeueAllMatureUBDQueue returns a concatenated list of all the timeslices inclusively previous to
// currTime, and deletes the timeslices from the queue.
func (k *Keeper) DequeueAllMatureUBDQueue(ctx context.Context, currTime time.Time) (matureUnbonds []types.DTData, err error) {
	store := k.storeService.OpenKVStore(ctx)

	// Get an iterator for all timeslices from time 0 until the current BlockHeader time
	unbondingTimesliceIterator, err := k.UBDQueueIterator(ctx, currTime)
	if err != nil {
		return nil, err
	}
	defer unbondingTimesliceIterator.Close()

	for ; unbondingTimesliceIterator.Valid(); unbondingTimesliceIterator.Next() {
		timeslice := types.DTDataList{}
		value := unbondingTimesliceIterator.Value()
		k.cdc.MustUnmarshal(value, &timeslice)

		matureUnbonds = append(matureUnbonds, timeslice.Data...)

		err = store.Delete(unbondingTimesliceIterator.Key())
		if err != nil {
			return nil, err
		}
	}

	return matureUnbonds, nil
}
