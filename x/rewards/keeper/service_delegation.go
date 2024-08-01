package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

// initialize starting info for a new delegation
func (k *Keeper) initializeServiceDelegation(ctx context.Context, serviceID uint32, del string) error {
	// period has already been incremented - we want to store the period ended by this delegation action
	currentRewards, err := k.ServiceCurrentRewards.Get(ctx, serviceID)
	if err != nil {
		return err
	}
	previousPeriod := currentRewards.Period - 1

	// increment reference count for the period we're going to track
	err = k.incrementServiceReferenceCount(ctx, serviceID, previousPeriod)
	if err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	service, found := k.servicesKeeper.GetService(sdkCtx, serviceID)
	if !found {
		return servicestypes.ErrServiceNotFound
	}

	delegation, found := k.restakingKeeper.GetServiceDelegation(sdkCtx, serviceID, del)
	if !found {
		return sdkerrors.ErrNotFound.Wrapf("service delegation not found: %d, %s", serviceID, del)
	}

	// calculate delegation stake in tokens
	// we don't store directly, so multiply delegation shares * (tokens per share)
	// note: necessary to truncate so we don't allow withdrawing more rewards than owed
	stake := service.TokensFromShares(delegation.Shares)
	return k.ServiceDelegatorStartingInfos.Set(
		ctx, collections.Join(serviceID, del),
		types.NewMultiDelegatorStartingInfo(previousPeriod, stake, uint64(sdkCtx.BlockHeight())))
}

// calculate the rewards accrued by a delegation between two periods
func (k *Keeper) calculateServiceDelegationRewardsBetween(ctx context.Context, service servicestypes.Service,
	startingPeriod, endingPeriod uint64, stakes sdk.DecCoins,
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
	starting, err := k.ServiceHistoricalRewards.Get(ctx, collections.Join(service.ID, startingPeriod))
	if err != nil {
		return nil, err
	}

	ending, err := k.ServiceHistoricalRewards.Get(ctx, collections.Join(service.ID, endingPeriod))
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
func (k *Keeper) CalculateServiceDelegationRewards(ctx context.Context, service servicestypes.Service, del restakingtypes.Delegation, endingPeriod uint64) (rewards types.DecPools, err error) {
	// fetch starting info for delegation
	startingInfo, err := k.ServiceDelegatorStartingInfos.Get(ctx, collections.Join(service.ID, del.UserAddress))
	if err != nil {
		return
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if startingInfo.Height == uint64(sdkCtx.BlockHeight()) {
		// started this height, no rewards yet
		return
	}

	startingPeriod := startingInfo.PreviousPeriod
	stakes := startingInfo.Stakes

	startingHeight := startingInfo.Height
	// Slashes this block happened after reward allocation, but we have to account
	// for them for the stake sanity check below.
	endingHeight := uint64(sdkCtx.BlockHeight())
	if endingHeight > startingHeight {
		// TODO: handle slash events
	}

	// A total stake sanity check; Recalculated final stake should be less than or
	// equal to current stake here. We cannot use Equals because stake is truncated
	// when multiplied by slash fractions (see above). We could only use equals if
	// we had arbitrary-precision rationals.
	currentStakes := service.TokensFromShares(del.Shares)

	for i, stake := range stakes {
		currentStake := currentStakes.AmountOf(stake.Denom)
		if stake.Amount.GT(currentStake) {
			// AccountI for rounding inconsistencies between:
			//
			//     currentStake: calculated as in staking with a single computation
			//     stake:        calculated as an accumulation of stake
			//                   calculations across service's distribution periods
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
	delRewards, err := k.calculateServiceDelegationRewardsBetween(ctx, service, startingPeriod, endingPeriod, stakes)
	if err != nil {
		return nil, err
	}

	rewards = rewards.Add(delRewards...)
	return rewards, nil
}

func (k *Keeper) withdrawServiceDelegationRewards(ctx context.Context, service servicestypes.Service, del restakingtypes.Delegation) (types.Pools, error) {
	// check existence of delegator starting info
	hasInfo, err := k.ServiceDelegatorStartingInfos.Has(ctx, collections.Join(service.ID, del.UserAddress))
	if err != nil {
		return nil, err
	}
	if !hasInfo {
		return nil, types.ErrEmptyDelegationDistInfo
	}

	// end current period and calculate rewards
	endingPeriod, err := k.IncrementServicePeriod(ctx, service)
	if err != nil {
		return nil, err
	}

	rewardsRaw, err := k.CalculateServiceDelegationRewards(ctx, service, del, endingPeriod)
	if err != nil {
		return nil, err
	}

	outstanding, err := k.GetServiceOutstandingRewardsCoins(ctx, service.ID)
	if err != nil {
		return nil, err
	}

	// defensive edge case may happen on the very final digits
	// of the decCoins due to operation order of the distribution mechanism.
	rewards := rewardsRaw.Intersect(outstanding)
	if !rewards.IsEqual(rewardsRaw) {
		logger := k.Logger(ctx)
		logger.Info(
			"rounding error withdrawing rewards from service",
			"delegator", del.UserAddress,
			"service", service.ID,
			"got", rewards.String(),
			"expected", rewardsRaw.String(),
		)
	}

	// truncate reward dec coins, return remainder to community service
	// TODO: return remainder to community service
	pools, _ := rewards.TruncateDecimal()
	coins := pools.Sum()

	// add pools to user account
	if !pools.IsEmpty() {
		withdrawAddr, err := k.GetDelegatorWithdrawAddr(ctx, del.UserAddress)
		if err != nil {
			return nil, err
		}

		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawAddr, coins)
		if err != nil {
			return nil, err
		}
	}

	// update the outstanding rewards and the community service only if the
	// transaction was successful
	err = k.ServiceOutstandingRewards.Set(ctx, service.ID, types.MultiOutstandingRewards{Rewards: outstanding.Sub(rewards)})
	if err != nil {
		return nil, err
	}

	// decrement reference count of starting period
	startingInfo, err := k.ServiceDelegatorStartingInfos.Get(ctx, collections.Join(service.ID, del.UserAddress))
	if err != nil {
		return nil, err
	}

	startingPeriod := startingInfo.PreviousPeriod
	err = k.decrementServiceReferenceCount(ctx, service.ID, startingPeriod)
	if err != nil {
		return nil, err
	}

	// remove delegator starting info
	err = k.ServiceDelegatorStartingInfos.Remove(ctx, collections.Join(service.ID, del.UserAddress))
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWithdrawRewards,
			sdk.NewAttribute(sdk.AttributeKeyAmount, coins.String()),
			sdk.NewAttribute(types.AttributeKeyServiceID, fmt.Sprint(service.ID)),
			sdk.NewAttribute(types.AttributeKeyDelegator, del.UserAddress),
		),
	)

	return pools, nil
}
