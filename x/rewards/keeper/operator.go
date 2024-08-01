package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
)

// initialize rewards for a new operator
func (k *Keeper) InitializeOperator(ctx context.Context, operator operatorstypes.Operator) error {
	// set initial historical rewards (period 0) with reference count of 1
	err := k.OperatorHistoricalRewards.Set(ctx, collections.Join(operator.ID, uint64(0)), types.NewMultiHistoricalRewards(types.DecPools{}, 1))
	if err != nil {
		return err
	}

	// set current rewards (starting at period 1)
	err = k.OperatorCurrentRewards.Set(ctx, operator.ID, types.NewMultiCurrentRewards(types.DecPools{}, 1))
	if err != nil {
		return err
	}

	// set accumulated commission
	err = k.OperatorAccumulatedCommissions.Set(ctx, operator.ID, types.InitialMultiAccumulatedCommission())
	if err != nil {
		return err
	}

	// set outstanding rewards
	err = k.OperatorOutstandingRewards.Set(ctx, operator.ID, types.MultiOutstandingRewards{Rewards: types.DecPools{}})
	return err
}

// increment operator period, returning the period just ended
func (k *Keeper) IncrementOperatorPeriod(ctx context.Context, operator operatorstypes.Operator) (uint64, error) {
	// fetch current rewards
	rewards, err := k.OperatorCurrentRewards.Get(ctx, operator.ID)
	if err != nil {
		return 0, err
	}

	// calculate current ratio
	var current types.DecPools

	tokens := operator.Tokens
	communityFunding := types.DecPools{}
	for _, token := range tokens {
		rewardCoins := rewards.Rewards.CoinsOf(token.Denom)
		if token.IsZero() {
			// can't calculate ratio for zero-token operators
			// ergo we instead add to the community pool
			communityFunding = communityFunding.Add(types.NewDecPool(token.Denom, rewardCoins))
			current = append(current, types.NewDecPool(token.Denom, sdk.DecCoins{}))
		} else {
			current = append(current,
				types.NewDecPool(token.Denom, rewardCoins.QuoDecTruncate(math.LegacyNewDecFromInt(token.Amount))))
		}
	}

	outstanding, err := k.OperatorOutstandingRewards.Get(ctx, operator.ID)
	if err != nil {
		return 0, err
	}

	communityFundingCoins, _ := communityFunding.TruncateDecimal()
	moduleAcc := k.accountKeeper.GetModuleAddress(types.ModuleName)

	if err := k.communityPoolKeeper.FundCommunityPool(ctx, communityFundingCoins.Sum(), moduleAcc); err != nil {
		return 0, err
	}
	// Since we sent only truncated coins, subtract that amount from outstanding
	// rewards, too.
	outstanding.Rewards = outstanding.Rewards.Sub(types.NewDecPoolsFromPools(communityFundingCoins))

	err = k.OperatorOutstandingRewards.Set(ctx, operator.ID, outstanding)
	if err != nil {
		return 0, err
	}

	// fetch historical rewards for last period
	historical, err := k.OperatorHistoricalRewards.Get(ctx, collections.Join(operator.ID, rewards.Period-1))
	if err != nil {
		return 0, err
	}

	// decrement reference count
	err = k.decrementOperatorReferenceCount(ctx, operator.ID, rewards.Period-1)
	if err != nil {
		return 0, err
	}

	// set new historical rewards with reference count of 1
	err = k.OperatorHistoricalRewards.Set(
		ctx, collections.Join(operator.ID, rewards.Period),
		types.NewMultiHistoricalRewards(historical.CumulativeRewardRatios.Add(current...), 1))
	if err != nil {
		return 0, err
	}

	// set current rewards, incrementing period by 1
	err = k.OperatorCurrentRewards.Set(ctx, operator.ID, types.NewMultiCurrentRewards(types.DecPools{}, rewards.Period+1))
	if err != nil {
		return 0, err
	}

	return rewards.Period, nil
}

// increment the reference count for a historical rewards value
func (k *Keeper) incrementOperatorReferenceCount(ctx context.Context, operatorID uint32, period uint64) error {
	historical, err := k.OperatorHistoricalRewards.Get(ctx, collections.Join(operatorID, period))
	if err != nil {
		return err
	}
	if historical.ReferenceCount > 2 {
		panic("reference count should never exceed 2")
	}
	historical.ReferenceCount++
	return k.OperatorHistoricalRewards.Set(ctx, collections.Join(operatorID, period), historical)
}

// decrement the reference count for a historical rewards value, and delete if zero references remain
func (k *Keeper) decrementOperatorReferenceCount(ctx context.Context, operatorID uint32, period uint64) error {
	historical, err := k.OperatorHistoricalRewards.Get(ctx, collections.Join(operatorID, period))
	if err != nil {
		return err
	}

	if historical.ReferenceCount == 0 {
		panic("cannot set negative reference count")
	}
	historical.ReferenceCount--
	if historical.ReferenceCount == 0 {
		return k.OperatorHistoricalRewards.Remove(ctx, collections.Join(operatorID, period))
	}

	return k.OperatorHistoricalRewards.Set(ctx, collections.Join(operatorID, period), historical)
}
