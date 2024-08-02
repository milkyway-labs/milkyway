package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
)

// initialize starting info for a new delegation
func (k *Keeper) initializePoolDelegation(ctx context.Context, poolID uint32, del sdk.AccAddress) error {
	// period has already been incremented - we want to store the period ended by this delegation action
	currentRewards, err := k.PoolCurrentRewards.Get(ctx, poolID)
	if err != nil {
		return err
	}
	previousPeriod := currentRewards.Period - 1

	// increment reference count for the period we're going to track
	err = k.incrementPoolReferenceCount(ctx, poolID, previousPeriod)
	if err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	pool, found := k.poolsKeeper.GetPool(sdkCtx, poolID)
	if !found {
		return poolstypes.ErrPoolNotFound
	}

	delegation, found := k.restakingKeeper.GetPoolDelegation(sdkCtx, poolID, del.String())
	if !found {
		return sdkerrors.ErrNotFound.Wrapf("pool delegation not found: %d, %s", poolID, del.String())
	}

	// calculate delegation stake in tokens
	// we don't store directly, so multiply delegation shares * (tokens per share)
	// note: necessary to truncate so we don't allow withdrawing more rewards than owed

	stake := pool.TokensFromSharesTruncated(delegation.Shares.AmountOf(pool.Denom))
	return k.PoolDelegatorStartingInfos.Set(ctx, collections.Join(poolID, del), types.NewDelegatorStartingInfo(previousPeriod, stake, uint64(sdkCtx.BlockHeight())))
}

// calculate the rewards accrued by a delegation between two periods
func (k *Keeper) calculatePoolDelegationRewardsBetween(ctx context.Context, pool poolstypes.Pool,
	startingPeriod, endingPeriod uint64, stake math.LegacyDec,
) (sdk.DecCoins, error) {
	// sanity check
	if startingPeriod > endingPeriod {
		panic("startingPeriod cannot be greater than endingPeriod")
	}

	// sanity check
	if stake.IsNegative() {
		panic("stake should not be negative")
	}

	// return staking * (ending - starting)
	starting, err := k.PoolHistoricalRewards.Get(ctx, collections.Join(pool.ID, startingPeriod))
	if err != nil {
		return sdk.DecCoins{}, err
	}

	ending, err := k.PoolHistoricalRewards.Get(ctx, collections.Join(pool.ID, endingPeriod))
	if err != nil {
		return sdk.DecCoins{}, err
	}

	difference := ending.CumulativeRewardRatio.Sub(starting.CumulativeRewardRatio)
	if difference.IsAnyNegative() {
		panic("negative rewards should not be possible")
	}
	// note: necessary to truncate so we don't allow withdrawing more rewards than owed
	rewards := difference.MulDecTruncate(stake)
	return rewards, nil
}

// calculate the total rewards accrued by a delegation
func (k *Keeper) CalculatePoolDelegationRewards(ctx context.Context, pool poolstypes.Pool, del restakingtypes.Delegation, endingPeriod uint64) (rewards sdk.DecCoins, err error) {
	delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(del.UserAddress)
	if err != nil {
		return nil, err
	}

	// fetch starting info for delegation
	startingInfo, err := k.PoolDelegatorStartingInfos.Get(ctx, collections.Join(pool.ID, sdk.AccAddress(delAddr)))
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if startingInfo.Height == uint64(sdkCtx.BlockHeight()) {
		// started this height, no rewards yet
		return nil, nil
	}

	startingPeriod := startingInfo.PreviousPeriod
	stake := startingInfo.Stake

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
	currentStake := pool.TokensFromShares(del.Shares.AmountOf(pool.Denom))

	if stake.GT(currentStake) {
		// AccountI for rounding inconsistencies between:
		//
		//     currentStake: calculated as in staking with a single computation
		//     stake:        calculated as an accumulation of stake
		//                   calculations across pool's distribution periods
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
		if stake.LTE(currentStake.Add(marginOfErr)) {
			stake = currentStake
		} else {
			panic(fmt.Sprintf("calculated final stake for delegator %s greater than current stake"+
				"\n\tfinal stake:\t%s"+
				"\n\tcurrent stake:\t%s",
				del.UserAddress, stake, currentStake))
		}
	}

	// calculate rewards for final period
	delRewards, err := k.calculatePoolDelegationRewardsBetween(ctx, pool, startingPeriod, endingPeriod, stake)
	if err != nil {
		return sdk.DecCoins{}, err
	}

	rewards = rewards.Add(delRewards...)
	return rewards, nil
}

func (k *Keeper) withdrawPoolDelegationRewards(ctx context.Context, pool poolstypes.Pool, del restakingtypes.Delegation) (sdk.Coins, error) {
	delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(del.UserAddress)
	if err != nil {
		return nil, err
	}

	// check existence of delegator starting info
	hasInfo, err := k.PoolDelegatorStartingInfos.Has(ctx, collections.Join(pool.ID, sdk.AccAddress(delAddr)))
	if err != nil {
		return nil, err
	}

	if !hasInfo {
		return nil, types.ErrEmptyDelegationDistInfo
	}

	// end current period and calculate rewards
	endingPeriod, err := k.IncrementPoolPeriod(ctx, pool)
	if err != nil {
		return nil, err
	}

	rewardsRaw, err := k.CalculatePoolDelegationRewards(ctx, pool, del, endingPeriod)
	if err != nil {
		return nil, err
	}

	outstanding, err := k.GetPoolOutstandingRewardsCoins(ctx, pool.ID)
	if err != nil {
		return nil, err
	}

	// defensive edge case may happen on the very final digits
	// of the decCoins due to operation order of the distribution mechanism.
	rewards := rewardsRaw.Intersect(outstanding)
	if !rewards.Equal(rewardsRaw) {
		logger := k.Logger(ctx)
		logger.Info(
			"rounding error withdrawing rewards from pool",
			"delegator", del.UserAddress,
			"pool", pool.ID,
			"got", rewards.String(),
			"expected", rewardsRaw.String(),
		)
	}

	// truncate reward dec coins, return remainder to community pool
	finalRewards, _ := rewards.TruncateDecimal()

	// add coins to user account
	if !finalRewards.IsZero() {
		withdrawAddr, err := k.GetDelegatorWithdrawAddr(ctx, delAddr)
		if err != nil {
			return nil, err
		}

		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawAddr, finalRewards)
		if err != nil {
			return nil, err
		}
	}

	// update the outstanding rewards and the community pool only if the
	// transaction was successful
	err = k.PoolOutstandingRewards.Set(ctx, pool.ID, types.OutstandingRewards{Rewards: outstanding.Sub(rewards)})
	if err != nil {
		return nil, err
	}

	// We don't use truncation remainder

	// decrement reference count of starting period
	startingInfo, err := k.PoolDelegatorStartingInfos.Get(ctx, collections.Join(pool.ID, sdk.AccAddress(delAddr)))
	if err != nil {
		return nil, err
	}

	startingPeriod := startingInfo.PreviousPeriod
	err = k.decrementPoolReferenceCount(ctx, pool.ID, startingPeriod)
	if err != nil {
		return nil, err
	}

	// remove delegator starting info
	err = k.PoolDelegatorStartingInfos.Remove(ctx, collections.Join(pool.ID, sdk.AccAddress(delAddr)))
	if err != nil {
		return nil, err
	}

	if finalRewards.IsZero() {
		baseDenom, _ := sdk.GetBaseDenom()
		if baseDenom == "" {
			baseDenom = sdk.DefaultBondDenom
		}

		// Note, we do not call the NewCoins constructor as we do not want the zero
		// coin removed.
		finalRewards = sdk.Coins{sdk.NewCoin(baseDenom, math.ZeroInt())}
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWithdrawRewards,
			sdk.NewAttribute(sdk.AttributeKeyAmount, finalRewards.String()),
			sdk.NewAttribute(types.AttributeKeyPoolID, fmt.Sprint(pool.ID)),
			sdk.NewAttribute(types.AttributeKeyDelegator, del.UserAddress),
		),
	)

	return finalRewards, nil
}
