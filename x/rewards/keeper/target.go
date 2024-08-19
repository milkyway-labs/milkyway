package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

func (k *Keeper) GetDelegationTarget(
	ctx context.Context, delType restakingtypes.DelegationType, targetID uint32,
) (restakingtypes.DelegationTarget, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	switch delType {
	case restakingtypes.DELEGATION_TYPE_POOL:
		pool, found := k.poolsKeeper.GetPool(sdkCtx, targetID)
		if !found {
			return nil, poolstypes.ErrPoolNotFound
		}
		return &pool, nil
	case restakingtypes.DELEGATION_TYPE_OPERATOR:
		operator, found := k.operatorsKeeper.GetOperator(sdkCtx, targetID)
		if !found {
			return nil, operatorstypes.ErrOperatorNotFound
		}
		return &operator, nil
	case restakingtypes.DELEGATION_TYPE_SERVICE:
		service, found := k.servicesKeeper.GetService(sdkCtx, targetID)
		if !found {
			return nil, servicestypes.ErrServiceNotFound
		}
		return &service, nil
	default:
		return nil, errors.Wrapf(restakingtypes.ErrInvalidDelegationType, "invalid delegation type: %s", delType)
	}
}

// initialize rewards for a new delegation target
func (k *Keeper) initializeDelegationTarget(ctx context.Context, target restakingtypes.DelegationTarget) error {
	// set initial historical rewards (period 0) with reference count of 1
	err := k.SetHistoricalRewards(ctx, target, uint64(0), types.NewHistoricalRewards(types.DecPools{}, 1))
	if err != nil {
		return err
	}

	// set current rewards (starting at period 1)
	err = k.SetCurrentRewards(ctx, target, types.NewCurrentRewards(types.DecPools{}, 1))
	if err != nil {
		return err
	}

	// set accumulated commission only if target is an operator
	if _, ok := target.(*operatorstypes.Operator); ok {
		err = k.OperatorAccumulatedCommissions.Set(ctx, target.GetID(), types.InitialAccumulatedCommission())
		if err != nil {
			return err
		}
	}

	// set outstanding rewards
	err = k.SetOutstandingRewards(ctx, target, types.OutstandingRewards{Rewards: types.DecPools{}})
	return err
}

// increment period, returning the period just ended
func (k *Keeper) IncrementDelegationTargetPeriod(ctx context.Context, target restakingtypes.DelegationTarget) (uint64, error) {
	// fetch current rewards
	rewards, err := k.GetCurrentRewards(ctx, target)
	if err != nil {
		return 0, err
	}

	// calculate current ratio
	var current types.DecPools

	tokens := target.GetTokens()
	communityFunding := types.DecPools{}
	for _, token := range tokens {
		rewardCoins := rewards.Rewards.CoinsOf(token.Denom)
		if token.IsZero() {
			// can't calculate ratio for zero-token targets
			// ergo we instead add to the community pool
			communityFunding = communityFunding.Add(types.NewDecPool(token.Denom, rewardCoins))
			current = append(current, types.NewDecPool(token.Denom, sdk.DecCoins{}))
		} else {
			current = append(current,
				types.NewDecPool(token.Denom, rewardCoins.QuoDecTruncate(math.LegacyNewDecFromInt(token.Amount))))
		}
	}

	outstanding, err := k.GetOutstandingRewards(ctx, target)
	if err != nil {
		return 0, err
	}

	communityFundingCoins, _ := communityFunding.TruncateDecimal()

	rewardsPoolAddr := k.accountKeeper.GetModuleAddress(types.RewardsPoolName)
	err = k.communityPoolKeeper.FundCommunityPool(ctx, communityFundingCoins.Sum(), rewardsPoolAddr)
	if err != nil {
		return 0, err
	}
	// Since we sent only truncated coins, subtract that amount from outstanding
	// rewards, too.
	outstanding.Rewards = outstanding.Rewards.Sub(types.NewDecPoolsFromPools(communityFundingCoins))

	err = k.SetOutstandingRewards(ctx, target, outstanding)
	if err != nil {
		return 0, err
	}

	// fetch historical rewards for last period
	historical, err := k.GetHistoricalRewards(ctx, target, rewards.Period-1)
	if err != nil {
		return 0, err
	}

	// decrement reference count
	err = k.decrementReferenceCount(ctx, target, rewards.Period-1)
	if err != nil {
		return 0, err
	}

	// set new historical rewards with reference count of 1
	err = k.SetHistoricalRewards(
		ctx, target, rewards.Period,
		types.NewHistoricalRewards(historical.CumulativeRewardRatios.Add(current...), 1))
	if err != nil {
		return 0, err
	}

	// set current rewards, incrementing period by 1
	err = k.SetCurrentRewards(ctx, target, types.NewCurrentRewards(types.DecPools{}, rewards.Period+1))
	if err != nil {
		return 0, err
	}

	return rewards.Period, nil
}

// increment the reference count for a historical rewards value
func (k *Keeper) incrementReferenceCount(ctx context.Context, target restakingtypes.DelegationTarget, period uint64) error {
	historical, err := k.GetHistoricalRewards(ctx, target, period)
	if err != nil {
		return err
	}
	if historical.ReferenceCount > 2 {
		panic("reference count should never exceed 2")
	}
	historical.ReferenceCount++
	return k.SetHistoricalRewards(ctx, target, period, historical)
}

// decrement the reference count for a historical rewards value, and delete if zero references remain
func (k *Keeper) decrementReferenceCount(ctx context.Context, target restakingtypes.DelegationTarget, period uint64) error {
	historical, err := k.GetHistoricalRewards(ctx, target, period)
	if err != nil {
		return err
	}

	if historical.ReferenceCount == 0 {
		panic("cannot set negative reference count")
	}
	historical.ReferenceCount--
	if historical.ReferenceCount == 0 {
		return k.RemoveHistoricalRewards(ctx, target, period)
	}

	return k.SetHistoricalRewards(ctx, target, period, historical)
}
