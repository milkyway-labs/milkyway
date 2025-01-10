package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	restakingtypes "github.com/milkyway-labs/milkyway/v7/x/restaking/types"
	"github.com/milkyway-labs/milkyway/v7/x/rewards/types"
)

// initializeDelegation initializes a delegation for a target
func (k *Keeper) initializeDelegation(ctx context.Context, target DelegationTarget, delAddr sdk.AccAddress) error {
	// Period has already been incremented - we want to store the period ended by this delegation action
	currentRewards, err := target.CurrentRewards.Get(ctx, target.GetID())
	if err != nil {
		return err
	}
	previousPeriod := currentRewards.Period - 1

	// Increment reference count for the period we're going to track
	err = k.incrementReferenceCount(ctx, target, previousPeriod)
	if err != nil {
		return err
	}

	delegator, err := k.accountKeeper.AddressCodec().BytesToString(delAddr)
	if err != nil {
		return err
	}

	delegation, found, err := k.restakingKeeper.GetDelegationForTarget(ctx, target.DelegationTarget, delegator)
	if err != nil {
		return err
	}

	if !found {
		return sdkerrors.ErrNotFound.Wrapf("delegation not found: %d, %s", target.GetID(), delAddr.String())
	}

	// Calculate delegation stake in tokens.
	// We don't store directly, so multiply delegation shares * (tokens per share)
	// NOTE: it's necessary to truncate so we don't allow withdrawing more rewards than owed
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	stake := target.TokensFromSharesTruncated(delegation.Shares)
	return target.DelegatorStartingInfos.Set(
		ctx, collections.Join(target.GetID(), delAddr),
		types.NewDelegatorStartingInfo(previousPeriod, stake, uint64(sdkCtx.BlockHeight())),
	)
}

// calculateDelegationRewardsBetween calculates the rewards accrued by a delegation between two periods
func (k *Keeper) calculateDelegationRewardsBetween(
	ctx context.Context,
	target DelegationTarget,
	delegator string,
	startingPeriod, endingPeriod uint64,
	stakes sdk.DecCoins,
) (rewards types.DecPools, err error) {
	// Sanity check
	if startingPeriod > endingPeriod {
		panic("startingPeriod cannot be greater than endingPeriod")
	}

	// Sanity check
	if stakes.IsAnyNegative() {
		panic("stake should not be negative")
	}

	// Return staking * (ending - starting)
	starting, err := target.HistoricalRewards.Get(ctx, collections.Join(target.GetID(), startingPeriod))
	if err != nil {
		return nil, err
	}

	ending, err := target.HistoricalRewards.Get(ctx, collections.Join(target.GetID(), endingPeriod))
	if err != nil {
		return nil, err
	}

	differences := ending.CumulativeRewardRatios.Sub(starting.CumulativeRewardRatios)
	var decPools types.DecPools

	if target.DelegationType == restakingtypes.DELEGATION_TYPE_POOL {
		preferences, err := k.restakingKeeper.GetUserPreferences(ctx, delegator)
		if err != nil {
			return nil, err
		}

		for _, diff := range differences {
			if preferences.IsServiceTrustedWithPool(diff.ServiceID, target.GetID()) {
				decPools = decPools.Add(diff.DecPools...)
			}
		}
	} else {
		for _, diff := range differences {
			decPools = decPools.Add(diff.DecPools...)
		}
	}

	for _, decPool := range decPools {
		rewards = rewards.Add(types.NewDecPool(
			decPool.Denom,
			decPool.DecCoins.MulDecTruncate(stakes.AmountOf(decPool.Denom)),
		))
	}

	return rewards, nil
}

// CalculateDelegationRewards calculates the total rewards accrued by a delegation
// between the starting period and the current block height
func (k *Keeper) CalculateDelegationRewards(
	ctx context.Context, target DelegationTarget, del restakingtypes.Delegation, endingPeriod uint64,
) (rewards types.DecPools, err error) {
	delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(del.UserAddress)
	if err != nil {
		return nil, err
	}

	// Fetch starting info for delegation
	startingInfo, err := target.DelegatorStartingInfos.Get(ctx, collections.Join(target.GetID(), sdk.AccAddress(delAddr)))
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if startingInfo.Height == uint64(sdkCtx.BlockHeight()) {
		// started this height, no rewards yet
		return nil, nil
	}

	startingPeriod := startingInfo.PreviousPeriod
	stakes := startingInfo.Stakes

	// TODO: handle slash events
	// startingHeight := startingInfo.Height
	// // Slashes this block happened after reward allocation, but we have to account
	// // for them for the stake sanity check below.
	// endingHeight := uint64(sdkCtx.BlockHeight())
	// if endingHeight > startingHeight {
	// }

	// A total stake sanity check; Recalculated final stake should be less than or
	// equal to current stake here. We cannot use Equals because stake is truncated
	// when multiplied by slash fractions (see above). We could only use equals if
	// we had arbitrary-precision rationals.
	currentStakes := target.TokensFromShares(del.Shares)

	for i, stake := range stakes {
		currentStake := currentStakes.AmountOf(stake.Denom)
		if stake.Amount.GT(currentStake) {
			// AccountI for rounding inconsistencies between:
			//
			//     currentStake: calculated as in staking with a single computation
			//     stake:        calculated as an accumulation of stake
			//                   calculations across target's distribution periods
			//
			// These inconsistencies are due to differing order of operations which
			// will inevitably have different accumulated rounding and may lead to
			// the smallest decimal place being one greater in stake than
			// currentStake. When we calculated slashing by period, even if we
			// round down for each slash fraction, it's possible due to how much is
			// being rounded that we slash less when slashing by period instead of
			// for when we slash without periods. In other words, the single slash,
			// and the slashing by period could both be rounding down but the
			// slashing by period is simply rounding down less, thus making stake >
			// currentStake
			//
			// A small amount of this error is tolerated and corrected for,
			// however any greater amount should be considered a breach in expected
			// behavior.
			marginOfErr := math.LegacySmallestDec().MulInt64(3)
			if stake.Amount.LTE(currentStake.Add(marginOfErr)) {
				stakes[i].Amount = currentStake
			} else {
				panic(fmt.Sprintf("calculated final stake for delegator %s greater than current stake"+
					"\n\tstake denom:\t%s"+
					"\n\tfinal stake:\t%s"+
					"\n\tcurrent stake:\t%s",
					del.UserAddress, stake.Denom, stake.Amount, currentStake))
			}
		}
	}

	// Calculate the rewards for the final period
	delRewards, err := k.calculateDelegationRewardsBetween(
		ctx,
		target,
		del.UserAddress,
		startingPeriod,
		endingPeriod,
		stakes,
	)
	if err != nil {
		return nil, err
	}

	return delRewards, nil
}

// withdrawDelegationRewards withdraws the rewards from the delegation and reinitializes it
func (k *Keeper) withdrawDelegationRewards(
	ctx context.Context, target DelegationTarget, del restakingtypes.Delegation,
) (types.Pools, error) {
	delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(del.UserAddress)
	if err != nil {
		return nil, err
	}

	// Check the existence of delegator starting info
	hasInfo, err := target.DelegatorStartingInfos.Has(ctx, collections.Join(target.GetID(), sdk.AccAddress(delAddr)))
	if err != nil {
		return nil, err
	}
	if !hasInfo {
		return nil, types.ErrEmptyDelegationDistInfo
	}

	// End current period and calculate rewards
	endingPeriod, err := k.IncrementDelegationTargetPeriod(ctx, target)
	if err != nil {
		return nil, err
	}

	rewardsRaw, err := k.CalculateDelegationRewards(ctx, target, del, endingPeriod)
	if err != nil {
		return nil, err
	}

	outstanding, err := k.GetOutstandingRewardsCoins(ctx, target)
	if err != nil {
		return nil, err
	}

	// Defensive edge case may happen on the very final digits
	// of the decCoins due to operation order of the distribution mechanism.
	rewards := rewardsRaw.Intersect(outstanding)
	if !rewards.IsEqual(rewardsRaw) {
		logger := k.Logger(ctx)
		logger.Info(
			"rounding error withdrawing rewards from delegation target",
			"delegator", del.UserAddress,
			"delegation_type", del.Type.String(),
			"delegation_target_id", target.GetID(),
			"got", rewards.String(),
			"expected", rewardsRaw.String(),
		)
	}

	// Truncate reward dec coins, return remainder to community operator
	// TODO: return remainder to community operator
	pools, _ := rewards.TruncateDecimal()
	coins := pools.Sum()

	// Add pools to user account
	if !pools.IsEmpty() {
		withdrawAddr, err := k.GetDelegatorWithdrawAddr(ctx, delAddr)
		if err != nil {
			return nil, err
		}

		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.RewardsPoolName, withdrawAddr, coins)
		if err != nil {
			return nil, err
		}
	}

	// Update the outstanding rewards and the community operator only if the
	// transaction was successful
	err = target.OutstandingRewards.Set(ctx, target.GetID(), types.OutstandingRewards{Rewards: outstanding.Sub(rewards)})
	if err != nil {
		return nil, err
	}

	// Decrement reference count of starting period
	startingInfo, err := target.DelegatorStartingInfos.Get(ctx, collections.Join(target.GetID(), sdk.AccAddress(delAddr)))
	if err != nil {
		return nil, err
	}

	startingPeriod := startingInfo.PreviousPeriod
	err = k.decrementReferenceCount(ctx, target, startingPeriod)
	if err != nil {
		return nil, err
	}

	// Remove delegator starting info
	err = target.DelegatorStartingInfos.Remove(ctx, collections.Join(target.GetID(), sdk.AccAddress(delAddr)))
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWithdrawRewards,
			sdk.NewAttribute(types.AttributeKeyDelegationType, del.Type.String()),
			sdk.NewAttribute(types.AttributeKeyDelegationTargetID, fmt.Sprint(target.GetID())),
			sdk.NewAttribute(restakingtypes.AttributeKeyDelegator, del.UserAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, coins.String()),
			sdk.NewAttribute(types.AttributeKeyAmountPerPool, pools.String()),
		),
	)

	return pools, nil
}
