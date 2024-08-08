package keeper

import (
	"encoding/binary"
	"time"

	"cosmossdk.io/errors"
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

// ValidateUnbondAmount validates that a given unbond or redelegation amount is valid based on upon the
// converted shares. If the amount is valid, the total amount of respective shares is returned,
// otherwise an error is returned.
func (k *Keeper) ValidateUnbondAmount(
	ctx sdk.Context, delAddr sdk.AccAddress, target types.DelegationTarget, amt sdk.Coins,
) (shares sdk.DecCoins, err error) {
	// Get the delegation
	delegation, found := k.GetDelegationForTarget(ctx, target, delAddr.String())
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

// Undelegate unbonds an amount of delegator shares from a given validator. It
// will verify that the unbonding entries between the delegator and validator
// are not exceeded and unbond the staked tokens (based on shares) by creating
// an unbonding object and inserting it into the unbonding queue which will be
// processed during the staking EndBlocker.
func (k *Keeper) Undelegate(ctx sdk.Context, data types.UndelegationData) (time.Time, error) {
	// TODO: Probably we should implement this as well
	//if k.HasMaxUnbondingDelegationEntries(ctx, delAddr, valAddr) {
	//	return time.Time{}, types.ErrMaxUnbondingDelegationEntries
	//}

	// Unbond the tokens
	returnAmount, err := k.Unbond(ctx, data)
	if err != nil {
		return time.Time{}, err
	}

	// Compute the time at which the unbonding delegation should end
	completionTime := ctx.BlockHeader().Time.Add(k.UnbondingTime(ctx))

	// Store the unbonding delegation entry inside the store
	ubd, err := k.SetUnbondingDelegationEntry(ctx, data, ctx.BlockHeight(), completionTime, returnAmount)
	if err != nil {
		return time.Time{}, err
	}

	// Insert the unbonding delegation into the unbonding queue
	k.InsertUBDQueue(ctx, ubd, completionTime)

	return completionTime, nil
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
		k.operatorsKeeper.SaveOperator(ctx, operator)
		issuedTokensAmount = amount
	case *servicestypes.Service:
		service, amount := target.RemoveDelShares(data.Shares)
		k.servicesKeeper.SaveService(ctx, service)
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
	ubd, found := k.GetUnbondingDelegation(ctx, data.Delegator, data.Target)
	if found {
		ubd.AddEntry(creationHeight, minTime, balance, id)
	} else {
		ubd = data.BuildUnbondingDelegation(data.Delegator, data.Target.GetID(), creationHeight, minTime, balance, id)
	}

	err := k.SetUnbondingDelegation(ctx, ubd, id)
	if err != nil {
		return types.UnbondingDelegation{}, err
	}

	// Call the hook after the unbonding has been initiated
	err = k.hooks.AfterUnbondingInitiated(ctx, id)
	if err != nil {
		k.Logger(ctx).Error("failed to call after unbonding initiated hook", "error", err)
	}

	return ubd, nil
}
