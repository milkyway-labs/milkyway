package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	operatorstypes "github.com/milkyway-labs/milkyway/v12/x/operators/types"
	restakingtypes "github.com/milkyway-labs/milkyway/v12/x/restaking/types"
	"github.com/milkyway-labs/milkyway/v12/x/rewards/types"
)

// DeleteHistoricalRewards deletes all historical rewards for a target
func (k *Keeper) DeleteHistoricalRewards(ctx context.Context, target DelegationTarget) error {
	collection := target.HistoricalRewards

	// Walk over the collection and get the list of keys to be deleted
	var keys []collections.Pair[uint32, uint64]
	err := collection.Walk(
		ctx,
		collections.NewPrefixedPairRange[uint32, uint64](target.GetID()),
		func(key collections.Pair[uint32, uint64], value types.HistoricalRewards) (stop bool, err error) {
			keys = append(keys, key)
			return false, nil
		},
	)
	if err != nil {
		return err
	}

	// Delete all the keys from the collection
	for _, key := range keys {
		if err := collection.Remove(ctx, key); err != nil {
			return err
		}
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// GetRewardsPlans returns all rewards plans
func (k *Keeper) GetRewardsPlans(ctx context.Context) ([]types.RewardsPlan, error) {
	var rewardsPlans []types.RewardsPlan
	err := k.RewardsPlans.Walk(ctx, nil, func(id uint64, plan types.RewardsPlan) (stop bool, err error) {
		rewardsPlans = append(rewardsPlans, plan)
		return false, nil
	})
	return rewardsPlans, err
}

// --------------------------------------------------------------------------------------------------------------------

// GetOutstandingRewardsCoins returns the outstanding rewards coins for a target
func (k *Keeper) GetOutstandingRewardsCoins(ctx context.Context, target DelegationTarget) (types.DecPools, error) {
	rewards, err := target.OutstandingRewards.Get(ctx, target.GetID())
	if err != nil && !errors.IsOf(err, collections.ErrNotFound) {
		return nil, err
	}
	return rewards.Rewards, nil
}

// GetOperatorAccumulatedCommission returns the accumulated commission for an operator.
func (k *Keeper) GetOperatorAccumulatedCommission(ctx context.Context, operatorID uint32) (commission types.AccumulatedCommission, err error) {
	commission, err = k.OperatorAccumulatedCommissions.Get(ctx, operatorID)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return types.AccumulatedCommission{}, nil
		}
		return types.AccumulatedCommission{}, err
	}
	return
}

// DeleteOperatorAccumulatedCommission deletes the accumulated commission for an operator.
func (k *Keeper) DeleteOperatorAccumulatedCommission(ctx context.Context, operatorID uint32) error {
	return k.OperatorAccumulatedCommissions.Remove(ctx, operatorID)
}

// GetOperatorWithdrawAddr returns the outstanding rewards coins for an operator
func (k *Keeper) GetOperatorWithdrawAddr(ctx context.Context, operator operatorstypes.Operator) (sdk.AccAddress, error) {
	// Try getting a custom withdraw address
	operatorAddr, err := sdk.AccAddressFromBech32(operator.Address)
	if err != nil {
		return nil, err
	}

	withdrawAddr, err := k.GetDelegatorWithdrawAddr(ctx, operatorAddr)
	if err != nil {
		return nil, err
	}

	if withdrawAddr != nil {
		return withdrawAddr, nil
	}

	// By default, use the operator admin address as the withdraw address
	adminAddress, err := sdk.AccAddressFromBech32(operator.Admin)
	if err != nil {
		return nil, err
	}

	return adminAddress, nil
}

// GetDelegatorWithdrawAddr returns the delegator's withdraw address if set, otherwise the delegator address is returned.
func (k *Keeper) GetDelegatorWithdrawAddr(ctx context.Context, delegator sdk.AccAddress) (sdk.AccAddress, error) {
	addr, err := k.DelegatorWithdrawAddrs.Get(ctx, delegator)
	if err != nil && errors.IsOf(err, collections.ErrNotFound) {
		return delegator, nil
	}
	return addr, err
}

// --------------------------------------------------------------------------------------------------------------------

// GetDelegationRewards returns the rewards for a delegation
func (k *Keeper) GetDelegationRewards(
	ctx context.Context, delAddr sdk.AccAddress, delType restakingtypes.DelegationType, targetID uint32,
) (types.DecPools, error) {
	target, err := k.GetDelegationTarget(ctx, delType, targetID)
	if err != nil {
		return nil, err
	}

	delegator, err := k.accountKeeper.AddressCodec().BytesToString(delAddr)
	if err != nil {
		return nil, err
	}

	delegation, found, err := k.restakingKeeper.GetDelegationForTarget(ctx, target.DelegationTarget, delegator)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, errors.Wrap(sdkerrors.ErrNotFound, "delegation not found")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cacheCtx, _ := sdkCtx.CacheContext()
	endingPeriod, err := k.IncrementDelegationTargetPeriod(cacheCtx, target)
	if err != nil {
		return nil, err
	}

	rewards, err := k.CalculateDelegationRewards(cacheCtx, target, delegation, endingPeriod)
	if err != nil {
		return nil, err
	}

	return rewards, nil
}

// GetPoolDelegationRewards returns the rewards for a pool delegation
func (k *Keeper) GetPoolDelegationRewards(ctx context.Context, delAddr sdk.AccAddress, poolID uint32) (types.DecPools, error) {
	return k.GetDelegationRewards(ctx, delAddr, restakingtypes.DELEGATION_TYPE_POOL, poolID)
}

// GetOperatorDelegationRewards returns the rewards for an operator delegation
func (k *Keeper) GetOperatorDelegationRewards(ctx context.Context, delAddr sdk.AccAddress, operatorID uint32) (types.DecPools, error) {
	return k.GetDelegationRewards(ctx, delAddr, restakingtypes.DELEGATION_TYPE_OPERATOR, operatorID)
}

// GetServiceDelegationRewards returns the rewards for a service delegation
func (k *Keeper) GetServiceDelegationRewards(
	ctx context.Context, delAddr sdk.AccAddress, serviceID uint32,
) (types.DecPools, error) {
	return k.GetDelegationRewards(ctx, delAddr, restakingtypes.DELEGATION_TYPE_SERVICE, serviceID)
}

// --------------------------------------------------------------------------------------------------------------------

// GetPoolServiceTotalDelegatorShares returns the total delegator shares for a
// pool-service pair.
func (k *Keeper) GetPoolServiceTotalDelegatorShares(ctx context.Context, poolID, serviceID uint32) (sdk.DecCoins, error) {
	shares, err := k.PoolServiceTotalDelegatorShares.Get(ctx, collections.Join(poolID, serviceID))
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return shares.Shares, nil
}

// SetPoolServiceTotalDelegatorShares sets the total delegator shares for a
// pool-service pair.
func (k *Keeper) SetPoolServiceTotalDelegatorShares(ctx context.Context, poolID, serviceID uint32, shares sdk.DecCoins) error {
	return k.PoolServiceTotalDelegatorShares.Set(
		ctx,
		collections.Join(poolID, serviceID),
		types.NewPoolServiceTotalDelegatorShares(poolID, serviceID, shares),
	)
}

// DeletePoolServiceTotalDelegatorShares deletes the total delegator shares for a
// pool-service pair.
func (k *Keeper) DeletePoolServiceTotalDelegatorShares(ctx context.Context, poolID, serviceID uint32) error {
	return k.PoolServiceTotalDelegatorShares.Remove(ctx, collections.Join(poolID, serviceID))
}

// DeleteAllPoolServiceTotalDelegatorSharesByService deletes all total delegator
// shares for a pool.
func (k *Keeper) DeleteAllPoolServiceTotalDelegatorSharesByService(ctx context.Context, serviceID uint32) error {
	// Walk over the collection and get the list of keys to be deleted
	var keys []collections.Pair[uint32, uint32]
	err := k.PoolServiceTotalDelegatorShares.Walk(
		ctx,
		nil, // TODO: is there a better way to do this?
		func(key collections.Pair[uint32, uint32], value types.PoolServiceTotalDelegatorShares) (stop bool, err error) {
			if key.K2() == serviceID {
				keys = append(keys, key)
			}
			return false, nil
		},
	)
	if err != nil {
		return err
	}
	for _, key := range keys {
		err = k.PoolServiceTotalDelegatorShares.Remove(ctx, key)
		if err != nil {
			return err
		}
	}
	return nil
}

// IncrementPoolServiceTotalDelegatorShares increments the total delegator shares
// for a pool-service pair.
func (k *Keeper) IncrementPoolServiceTotalDelegatorShares(
	ctx context.Context, poolID, serviceID uint32, shares sdk.DecCoins,
) error {
	prevShares, err := k.GetPoolServiceTotalDelegatorShares(ctx, poolID, serviceID)
	if err != nil {
		return err
	}
	return k.SetPoolServiceTotalDelegatorShares(ctx, poolID, serviceID, prevShares.Add(shares...))
}

// DecrementPoolServiceTotalDelegatorShares decrements the total delegator shares
// for a pool-service pair.
func (k *Keeper) DecrementPoolServiceTotalDelegatorShares(
	ctx context.Context, poolID, serviceID uint32, shares sdk.DecCoins,
) error {
	prevShares, err := k.GetPoolServiceTotalDelegatorShares(ctx, poolID, serviceID)
	if err != nil {
		return err
	}
	newShares := prevShares.Sub(shares)
	// Delete the pool-service total delegator shares record if it becomes zero
	if newShares.IsZero() {
		return k.DeletePoolServiceTotalDelegatorShares(ctx, poolID, serviceID)
	}
	return k.SetPoolServiceTotalDelegatorShares(ctx, poolID, serviceID, newShares)
}

// GetDelegationTargetTrustedTokens returns the amount of tokens of a delegation target
// considering trusted services(it only applies to pools).
func (k *Keeper) GetDelegationTargetTrustedTokens(ctx context.Context, serviceID uint32, target DelegationTarget) (sdk.Coins, error) {
	tokens := target.GetTokens()
	if target.DelegationType == restakingtypes.DELEGATION_TYPE_POOL {
		totalDelShares, err := k.GetPoolServiceTotalDelegatorShares(ctx, target.GetID(), serviceID)
		if err != nil {
			return nil, err
		}
		tokens, _ = target.TokensFromSharesTruncated(totalDelShares).TruncateDecimal()
	}
	return tokens, nil
}
