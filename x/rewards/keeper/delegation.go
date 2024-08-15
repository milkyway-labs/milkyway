package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
)

func (k *Keeper) GetDelegation(ctx context.Context, target *types.DelegationTarget, del sdk.AccAddress) (restakingtypes.Delegation, bool) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	switch target.Type() {
	case restakingtypes.DELEGATION_TYPE_POOL:
		return k.restakingKeeper.GetPoolDelegation(sdkCtx, target.GetID(), del.String())
	case restakingtypes.DELEGATION_TYPE_OPERATOR:
		return k.restakingKeeper.GetOperatorDelegation(sdkCtx, target.GetID(), del.String())
	case restakingtypes.DELEGATION_TYPE_SERVICE:
		return k.restakingKeeper.GetServiceDelegation(sdkCtx, target.GetID(), del.String())
	default:
		panic("unknown delegation type")
	}
}

// initialize starting info for a new delegation
func (k *Keeper) initializeDelegation(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32, del sdk.AccAddress) error {
	target, err := k.GetDelegationTarget(ctx, delType, targetID)
	if err != nil {
		return err
	}

	// period has already been incremented - we want to store the period ended by this delegation action
	currentRewards, err := k.GetCurrentRewards(ctx, target)
	if err != nil {
		return err
	}
	previousPeriod := currentRewards.Period - 1

	// increment reference count for the period we're going to track
	err = k.incrementReferenceCount(ctx, target, previousPeriod)
	if err != nil {
		return err
	}

	delegation, found := k.GetDelegation(ctx, target, del)
	if !found {
		return sdkerrors.ErrNotFound.Wrapf("delegation not found: %d, %s", targetID, del.String())
	}

	// calculate delegation stake in tokens
	// we don't store directly, so multiply delegation shares * (tokens per share)
	// note: necessary to truncate so we don't allow withdrawing more rewards than owed
	stake := target.TokensFromSharesTruncated(delegation.Shares)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return k.SetDelegatorStartingInfo(
		ctx, target, del,
		types.NewDelegatorStartingInfo(previousPeriod, stake, uint64(sdkCtx.BlockHeight())))
}

// calculate the rewards accrued by a delegation between two periods
func (k *Keeper) calculateDelegationRewardsBetween(
	ctx context.Context, target *types.DelegationTarget, startingPeriod, endingPeriod uint64, stakes sdk.DecCoins,
) (rewards types.DecPools, err error) {
	// sanity check
	if startingPeriod > endingPeriod {
		panic("startingPeriod cannot be greater than endingPeriod")
	}

	// sanity check
	if stakes.IsAnyNegative() {
		panic("stake should not be negative")
	}

	// return staking * (ending - starting)
	starting, err := k.GetHistoricalRewards(ctx, target, startingPeriod)
	if err != nil {
		return nil, err
	}

	ending, err := k.GetHistoricalRewards(ctx, target, endingPeriod)
	if err != nil {
		return nil, err
	}

	differences := ending.CumulativeRewardRatios.Sub(starting.CumulativeRewardRatios)
	if differences.IsAnyNegative() {
		panic("negative rewards should not be possible")
	}

	for _, diff := range differences {
		rewards = append(rewards, types.NewDecPool(
			diff.Denom,
			diff.DecCoins.MulDecTruncate(stakes.AmountOf(diff.Denom)),
		))
	}
	return
}

// calculate the total rewards accrued by a delegation
func (k *Keeper) CalculateDelegationRewards(ctx context.Context, target *types.DelegationTarget, del restakingtypes.Delegation, endingPeriod uint64) (rewards types.DecPools, err error) {
	delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(del.UserAddress)
	if err != nil {
		return nil, err
	}

	// fetch starting info for delegation
	startingInfo, err := k.GetDelegatorStartingInfo(ctx, target, delAddr)
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
	//startingHeight := startingInfo.Height
	//// Slashes this block happened after reward allocation, but we have to account
	//// for them for the stake sanity check below.
	//endingHeight := uint64(sdkCtx.BlockHeight())
	//if endingHeight > startingHeight {
	//}

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

	// calculate rewards for final period
	delRewards, err := k.calculateDelegationRewardsBetween(ctx, target, startingPeriod, endingPeriod, stakes)
	if err != nil {
		return nil, err
	}

	rewards = rewards.Add(delRewards...)
	return rewards, nil
}

func (k *Keeper) withdrawDelegationRewards(
	ctx context.Context, target *types.DelegationTarget, del restakingtypes.Delegation,
) (types.Pools, error) {
	delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(del.UserAddress)
	if err != nil {
		return nil, err
	}

	// check existence of delegator starting info
	hasInfo, err := k.HasDelegatorStartingInfo(ctx, target, delAddr)
	if err != nil {
		return nil, err
	}
	if !hasInfo {
		return nil, types.ErrEmptyDelegationDistInfo
	}

	// end current period and calculate rewards
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

	// defensive edge case may happen on the very final digits
	// of the decCoins due to operation order of the distribution mechanism.
	rewards := rewardsRaw.Intersect(outstanding)
	if !rewards.IsEqual(rewardsRaw) {
		logger := k.Logger(ctx)
		logger.Info(
			"rounding error withdrawing rewards from delegation target",
			"delegator", del.UserAddress,
			"delegation_type", target.Type().String(),
			"delegation_target_id", target.GetID(),
			"got", rewards.String(),
			"expected", rewardsRaw.String(),
		)
	}

	// truncate reward dec coins, return remainder to community operator
	// TODO: return remainder to community operator
	pools, _ := rewards.TruncateDecimal()
	coins := pools.Sum()

	// add pools to user account
	if !pools.IsEmpty() {
		withdrawAddr, err := k.GetDelegatorWithdrawAddr(ctx, delAddr)
		if err != nil {
			return nil, err
		}

		err = k.bankKeeper.SendCoins(ctx, types.RewardsPoolAddress, withdrawAddr, coins)
		if err != nil {
			return nil, err
		}
	}

	// update the outstanding rewards and the community operator only if the
	// transaction was successful
	err = k.SetOutstandingRewards(ctx, target, types.OutstandingRewards{Rewards: outstanding.Sub(rewards)})
	if err != nil {
		return nil, err
	}

	// decrement reference count of starting period
	startingInfo, err := k.GetDelegatorStartingInfo(ctx, target, delAddr)
	if err != nil {
		return nil, err
	}

	startingPeriod := startingInfo.PreviousPeriod
	err = k.decrementReferenceCount(ctx, target, startingPeriod)
	if err != nil {
		return nil, err
	}

	// remove delegator starting info
	err = k.RemoveDelegatorStartingInfo(ctx, target, delAddr)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWithdrawRewards,
			sdk.NewAttribute(sdk.AttributeKeyAmount, coins.String()),
			sdk.NewAttribute(types.AttributeKeyDelegationType, target.Type().String()),
			sdk.NewAttribute(types.AttributeKeyDelegationTargetID, fmt.Sprint(target.GetID())),
			sdk.NewAttribute(restakingtypes.AttributeKeyDelegator, del.UserAddress),
		),
	)

	return pools, nil
}
