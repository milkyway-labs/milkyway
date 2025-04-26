package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	operatorstypes "github.com/milkyway-labs/milkyway/v12/x/operators/types"
	restakingtypes "github.com/milkyway-labs/milkyway/v12/x/restaking/types"
	"github.com/milkyway-labs/milkyway/v12/x/rewards/types"
)

// DelegationTarget is a wrapper around the delegation target that holds the
// delegation type and collection references for the target's type.
type DelegationTarget struct {
	restakingtypes.DelegationTarget
	DelegationType         restakingtypes.DelegationType
	DelegatorStartingInfos collections.Map[collections.Pair[uint32, sdk.AccAddress], types.DelegatorStartingInfo]
	HistoricalRewards      collections.Map[collections.Pair[uint32, uint64], types.HistoricalRewards]
	CurrentRewards         collections.Map[uint32, types.CurrentRewards]
	OutstandingRewards     collections.Map[uint32, types.OutstandingRewards]
}

// delegationTargetKey is used as a key of the delegation target cache map.
type delegationTargetKey struct {
	typ restakingtypes.DelegationType
	id  uint32
}

// delegationTargetCache is a cache of delegation targets.
type delegationTargetCache map[delegationTargetKey]DelegationTarget

// GetDelegationTarget returns the wrapped delegation target for the given
// delegation type and target ID.
func (k *Keeper) GetDelegationTarget(
	ctx context.Context,
	delType restakingtypes.DelegationType,
	targetID uint32,
) (DelegationTarget, error) {
	switch delType {
	case restakingtypes.DELEGATION_TYPE_POOL:
		pool, err := k.poolsKeeper.GetPool(ctx, targetID)
		if err != nil {
			return DelegationTarget{}, err
		}
		return DelegationTarget{
			DelegationTarget:       pool,
			DelegationType:         delType,
			DelegatorStartingInfos: k.PoolDelegatorStartingInfos,
			HistoricalRewards:      k.PoolHistoricalRewards,
			CurrentRewards:         k.PoolCurrentRewards,
			OutstandingRewards:     k.PoolOutstandingRewards,
		}, nil
	case restakingtypes.DELEGATION_TYPE_OPERATOR:
		operator, err := k.operatorsKeeper.GetOperator(ctx, targetID)
		if err != nil {
			return DelegationTarget{}, err
		}
		return DelegationTarget{
			DelegationTarget:       operator,
			DelegationType:         delType,
			DelegatorStartingInfos: k.OperatorDelegatorStartingInfos,
			HistoricalRewards:      k.OperatorHistoricalRewards,
			CurrentRewards:         k.OperatorCurrentRewards,
			OutstandingRewards:     k.OperatorOutstandingRewards,
		}, nil
	case restakingtypes.DELEGATION_TYPE_SERVICE:
		service, err := k.servicesKeeper.GetService(ctx, targetID)
		if err != nil {
			return DelegationTarget{}, err
		}
		return DelegationTarget{
			DelegationTarget:       service,
			DelegationType:         delType,
			DelegatorStartingInfos: k.ServiceDelegatorStartingInfos,
			HistoricalRewards:      k.ServiceHistoricalRewards,
			CurrentRewards:         k.ServiceCurrentRewards,
			OutstandingRewards:     k.ServiceOutstandingRewards,
		}, nil
	default:
		return DelegationTarget{}, errors.Wrapf(restakingtypes.ErrInvalidDelegationType, "invalid delegation type: %s", delType)
	}
}

// --------------------------------------------------------------------------------------------------------------------

// initialize rewards for a new delegation target
func (k *Keeper) initializeDelegationTarget(ctx context.Context, target DelegationTarget) error {
	// set initial historical rewards (period 0) with reference count of 1
	err := target.HistoricalRewards.Set(
		ctx,
		collections.Join(target.GetID(), uint64(0)),
		types.NewHistoricalRewards(types.ServicePools{}, 1),
	)
	if err != nil {
		return err
	}

	// set current rewards (starting at period 1)
	err = target.CurrentRewards.Set(ctx, target.GetID(), types.NewCurrentRewards(types.ServicePools{}, 1))
	if err != nil {
		return err
	}

	// set accumulated commission only if target is an operator
	if target.DelegationType == restakingtypes.DELEGATION_TYPE_OPERATOR {
		err = k.OperatorAccumulatedCommissions.Set(ctx, target.GetID(), types.InitialAccumulatedCommission())
		if err != nil {
			return err
		}
	}

	// set outstanding rewards
	return target.OutstandingRewards.Set(ctx, target.GetID(), types.OutstandingRewards{Rewards: types.DecPools{}})
}

// IncrementDelegationTargetPeriod increments the period, returning the period that just ended
func (k *Keeper) IncrementDelegationTargetPeriod(ctx context.Context, target DelegationTarget) (uint64, error) {
	// fetch current rewards
	rewards, err := target.CurrentRewards.Get(ctx, target.GetID())
	if err != nil {
		return 0, err
	}

	// calculate current ratio
	var current types.ServicePools

	communityFunding := types.DecPools{}
	for _, reward := range rewards.Rewards {
		var tokens sdk.DecCoins
		if target.DelegationType == restakingtypes.DELEGATION_TYPE_POOL {
			totalShares, err := k.GetPoolServiceTotalDelegatorShares(ctx, target.GetID(), reward.ServiceID)
			if err != nil {
				return 0, err
			}
			tokens = target.TokensFromSharesTruncated(totalShares)
		} else {
			tokens = sdk.NewDecCoinsFromCoins(target.GetTokens()...)
		}

		for _, token := range tokens {
			rewardCoins := reward.DecPools.CoinsOf(token.Denom)
			if token.IsZero() {
				// can't calculate ratio for zero-token targets
				// ergo we instead add to the community pool
				communityFunding = communityFunding.Add(types.NewDecPool(token.Denom, rewardCoins))
				current = current.Add(
					types.NewServicePool(reward.ServiceID, types.DecPools{types.NewDecPool(token.Denom, sdk.DecCoins{})}),
				)
			} else {
				current = current.Add(
					types.NewServicePool(
						reward.ServiceID,
						types.DecPools{types.NewDecPool(token.Denom, rewardCoins.QuoDecTruncate(token.Amount))},
					),
				)
			}
		}
	}

	outstanding, err := target.OutstandingRewards.Get(ctx, target.GetID())
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

	err = target.OutstandingRewards.Set(ctx, target.GetID(), outstanding)
	if err != nil {
		return 0, err
	}

	// fetch historical rewards for last period
	historical, err := target.HistoricalRewards.Get(ctx, collections.Join(target.GetID(), rewards.Period-1))
	if err != nil {
		return 0, err
	}

	// decrement reference count
	err = k.decrementReferenceCount(ctx, target, rewards.Period-1)
	if err != nil {
		return 0, err
	}

	// set new historical rewards with reference count of 1
	err = target.HistoricalRewards.Set(
		ctx, collections.Join(target.GetID(), rewards.Period),
		types.NewHistoricalRewards(historical.CumulativeRewardRatios.Add(current...), 1))
	if err != nil {
		return 0, err
	}

	// set current rewards, incrementing period by 1
	err = target.CurrentRewards.Set(ctx, target.GetID(), types.NewCurrentRewards(types.ServicePools{}, rewards.Period+1))
	if err != nil {
		return 0, err
	}

	return rewards.Period, nil
}

// increment the reference count for a historical rewards value
func (k *Keeper) incrementReferenceCount(ctx context.Context, target DelegationTarget, period uint64) error {
	historical, err := target.HistoricalRewards.Get(ctx, collections.Join(target.GetID(), period))
	if err != nil {
		return err
	}

	if historical.ReferenceCount > 2 {
		panic("reference count should never exceed 2")
	}

	historical.ReferenceCount++
	return target.HistoricalRewards.Set(ctx, collections.Join(target.GetID(), period), historical)
}

// decrement the reference count for a historical rewards value, and delete if zero references remain
func (k *Keeper) decrementReferenceCount(ctx context.Context, target DelegationTarget, period uint64) error {
	historical, err := target.HistoricalRewards.Get(ctx, collections.Join(target.GetID(), period))
	if err != nil {
		return err
	}

	if historical.ReferenceCount == 0 {
		panic("cannot set negative reference count")
	}

	historical.ReferenceCount--

	if historical.ReferenceCount == 0 {
		return target.HistoricalRewards.Remove(ctx, collections.Join(target.GetID(), period))
	}

	return target.HistoricalRewards.Set(ctx, collections.Join(target.GetID(), period), historical)
}

// clearDelegateTarget clears all rewards for a delegation target
func (k *Keeper) clearDelegationTarget(ctx context.Context, target DelegationTarget) error {
	// fetch outstanding
	outstandingCoins, err := k.GetOutstandingRewardsCoins(ctx, target)
	if err != nil {
		return err
	}

	outstanding := outstandingCoins.CoinsAmount()

	// Clear data related to an operator or service
	switch target.DelegationType {
	case restakingtypes.DELEGATION_TYPE_OPERATOR:
		// Clear data related to an operator
		operator, ok := target.DelegationTarget.(operatorstypes.Operator)
		if !ok {
			return fmt.Errorf("invalid delegation target type %T", target.DelegationTarget)
		}
		outstanding, err = k.clearOperator(ctx, outstanding, operator)
		if err != nil {
			return err
		}
	case restakingtypes.DELEGATION_TYPE_SERVICE:
		// Clear data related to a service
		err = k.DeleteAllPoolServiceTotalDelegatorSharesByService(ctx, target.GetID())
		if err != nil {
			return err
		}
	}

	// Add outstanding to community pool
	// The target is removed only after it has no more delegations.
	// This operation sends only the remaining dust to the community pool.
	rewardsPoolAddr := k.accountKeeper.GetModuleAddress(types.RewardsPoolName)

	// We truncate the outstanding to be able to send it to the community pool
	// The remainder will be just be removed
	outstandingTruncated, _ := outstanding.TruncateDecimal()
	err = k.communityPoolKeeper.FundCommunityPool(ctx, outstandingTruncated, rewardsPoolAddr)
	if err != nil {
		return err
	}

	// Delete outstanding rewards
	err = target.OutstandingRewards.Remove(ctx, target.GetID())
	if err != nil {
		return err
	}

	// Remove the commission record
	if target.DelegationType == restakingtypes.DELEGATION_TYPE_OPERATOR {
		err = k.DeleteOperatorAccumulatedCommission(ctx, target.GetID())
		if err != nil {
			return err
		}
	}

	// TODO: Clear slash events when we introduce slashing

	// Clear historical rewards
	err = k.DeleteHistoricalRewards(ctx, target)
	if err != nil {
		return err
	}

	// Clear current rewards
	return target.CurrentRewards.Remove(ctx, target.GetID())
}

func (k *Keeper) clearOperator(ctx context.Context, outstanding sdk.DecCoins, operator operatorstypes.Operator) (outstandingLeftOver sdk.DecCoins, err error) {
	// Force-withdraw commission
	valCommission, err := k.GetOperatorAccumulatedCommission(ctx, operator.ID)
	if err != nil {
		return outstanding, err
	}

	commission := valCommission.Commissions.CoinsAmount()

	if !commission.IsZero() {
		// Subtract from outstanding
		outstanding = outstanding.Sub(commission)

		// Split into integral & remainder
		coins, remainder := commission.TruncateDecimal()

		// We truncate the remainder to be able to send it to the community pool
		// The remainder will be just be removed
		remainderTruncated, _ := remainder.TruncateDecimal()

		// Send remainder to community pool
		rewardsPoolAddr := k.accountKeeper.GetModuleAddress(types.RewardsPoolName)
		err = k.communityPoolKeeper.FundCommunityPool(ctx, remainderTruncated, rewardsPoolAddr)
		if err != nil {
			return outstanding, err
		}

		// Add to operator account
		if !coins.IsZero() {
			withdrawAddr, err := k.GetOperatorWithdrawAddr(ctx, operator)
			if err != nil {
				return outstanding, err
			}

			err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawAddr, coins)
			if err != nil {
				return outstanding, err
			}
		}
	}

	return outstanding, nil
}
