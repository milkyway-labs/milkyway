package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
)

// initialize rewards for a new pool
func (k *Keeper) InitializePool(ctx context.Context, pool poolstypes.Pool) error {
	// set initial historical rewards (period 0) with reference count of 1
	err := k.PoolHistoricalRewards.Set(ctx, collections.Join(pool.ID, uint64(0)), types.NewHistoricalRewards(sdk.DecCoins{}, 1))
	if err != nil {
		return err
	}

	// set current rewards (starting at period 1)
	err = k.PoolCurrentRewards.Set(ctx, pool.ID, types.NewCurrentRewards(sdk.DecCoins{}, 1))
	if err != nil {
		return err
	}

	// set outstanding rewards
	err = k.PoolOutstandingRewards.Set(ctx, pool.ID, types.OutstandingRewards{Rewards: sdk.DecCoins{}})
	return err
}

// increment pool period, returning the period just ended
func (k *Keeper) IncrementPoolPeriod(ctx context.Context, pool poolstypes.Pool) (uint64, error) {
	// fetch current rewards
	rewards, err := k.PoolCurrentRewards.Get(ctx, pool.ID)
	if err != nil {
		return 0, err
	}

	// calculate current ratio
	var current sdk.DecCoins
	if pool.Tokens.IsZero() {
		// can't calculate ratio for zero-token pool
		// ergo we instead add to the community pool
		truncatedRewards, _ := rewards.Rewards.TruncateDecimal()
		moduleAcc := k.accountKeeper.GetModuleAddress(types.ModuleName)
		if err := k.communityPoolKeeper.FundCommunityPool(ctx, truncatedRewards, moduleAcc); err != nil {
			return 0, err
		}

		outstanding, err := k.PoolOutstandingRewards.Get(ctx, pool.ID)
		if err != nil {
			return 0, err
		}

		outstanding.Rewards = outstanding.Rewards.Sub(rewards.Rewards)

		err = k.PoolOutstandingRewards.Set(ctx, pool.ID, outstanding)
		if err != nil {
			return 0, err
		}

		current = sdk.DecCoins{}
	} else {
		// note: necessary to truncate so we don't allow withdrawing more rewards than owed
		current = rewards.Rewards.QuoDecTruncate(math.LegacyNewDecFromInt(pool.Tokens))
	}

	// fetch historical rewards for last period
	historical, err := k.PoolHistoricalRewards.Get(ctx, collections.Join(pool.ID, rewards.Period-1))
	if err != nil {
		return 0, err
	}

	cumRewardRatio := historical.CumulativeRewardRatio

	// decrement reference count
	err = k.decrementPoolReferenceCount(ctx, pool.ID, rewards.Period-1)
	if err != nil {
		return 0, err
	}

	// set new historical rewards with reference count of 1
	err = k.PoolHistoricalRewards.Set(
		ctx, collections.Join(pool.ID, rewards.Period),
		types.NewHistoricalRewards(cumRewardRatio.Add(current...), 1))
	if err != nil {
		return 0, err
	}

	// set current rewards, incrementing period by 1
	err = k.PoolCurrentRewards.Set(ctx, pool.ID, types.NewCurrentRewards(sdk.DecCoins{}, rewards.Period+1))
	if err != nil {
		return 0, err
	}

	return rewards.Period, nil
}

// increment the reference count for a historical rewards value
func (k *Keeper) incrementPoolReferenceCount(ctx context.Context, poolID uint32, period uint64) error {
	historical, err := k.PoolHistoricalRewards.Get(ctx, collections.Join(poolID, period))
	if err != nil {
		return err
	}
	if historical.ReferenceCount > 2 {
		panic("reference count should never exceed 2")
	}
	historical.ReferenceCount++
	return k.PoolHistoricalRewards.Set(ctx, collections.Join(poolID, period), historical)
}

// decrement the reference count for a historical rewards value, and delete if zero references remain
func (k *Keeper) decrementPoolReferenceCount(ctx context.Context, poolID uint32, period uint64) error {
	historical, err := k.PoolHistoricalRewards.Get(ctx, collections.Join(poolID, period))
	if err != nil {
		return err
	}

	if historical.ReferenceCount == 0 {
		panic("cannot set negative reference count")
	}
	historical.ReferenceCount--
	if historical.ReferenceCount == 0 {
		return k.PoolHistoricalRewards.Remove(ctx, collections.Join(poolID, period))
	}

	return k.PoolHistoricalRewards.Set(ctx, collections.Join(poolID, period), historical)
}
