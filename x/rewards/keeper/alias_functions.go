package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

// SetDelegatorStartingInfo sets the delegator starting info for a delegator.
func (k *Keeper) SetDelegatorStartingInfo(
	ctx context.Context, target restakingtypes.DelegationTarget, del sdk.AccAddress, info types.DelegatorStartingInfo,
) error {
	switch target.(type) {
	case poolstypes.Pool:
		return k.PoolDelegatorStartingInfos.Set(ctx, collections.Join(target.GetID(), del), info)
	case operatorstypes.Operator:
		return k.OperatorDelegatorStartingInfos.Set(ctx, collections.Join(target.GetID(), del), info)
	case servicestypes.Service:
		return k.ServiceDelegatorStartingInfos.Set(ctx, collections.Join(target.GetID(), del), info)
	default:
		return errors.Wrapf(restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target)
	}
}

// HasDelegatorStartingInfo returns true if the delegator starting info exists for a delegator and target
func (k *Keeper) HasDelegatorStartingInfo(
	ctx context.Context, target restakingtypes.DelegationTarget, delegator sdk.AccAddress,
) (bool, error) {
	switch target.(type) {
	case poolstypes.Pool:
		return k.PoolDelegatorStartingInfos.Has(ctx, collections.Join(target.GetID(), delegator))
	case operatorstypes.Operator:
		return k.OperatorDelegatorStartingInfos.Has(ctx, collections.Join(target.GetID(), delegator))
	case servicestypes.Service:
		return k.ServiceDelegatorStartingInfos.Has(ctx, collections.Join(target.GetID(), delegator))
	default:
		return false, errors.Wrapf(restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target)
	}
}

// GetDelegatorStartingInfo returns the delegator starting info for a delegator and target
func (k *Keeper) GetDelegatorStartingInfo(
	ctx context.Context, target restakingtypes.DelegationTarget, delegator sdk.AccAddress,
) (types.DelegatorStartingInfo, error) {
	switch target.(type) {
	case poolstypes.Pool:
		return k.PoolDelegatorStartingInfos.Get(ctx, collections.Join(target.GetID(), delegator))
	case operatorstypes.Operator:
		return k.OperatorDelegatorStartingInfos.Get(ctx, collections.Join(target.GetID(), delegator))
	case servicestypes.Service:
		return k.ServiceDelegatorStartingInfos.Get(ctx, collections.Join(target.GetID(), delegator))
	default:
		return types.DelegatorStartingInfo{}, errors.Wrapf(
			restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target,
		)
	}
}

// RemoveDelegatorStartingInfo removes the delegator starting info for a delegator and target
func (k *Keeper) RemoveDelegatorStartingInfo(
	ctx context.Context, target restakingtypes.DelegationTarget, delegator sdk.AccAddress,
) error {
	switch target.(type) {
	case poolstypes.Pool:
		return k.PoolDelegatorStartingInfos.Remove(ctx, collections.Join(target.GetID(), delegator))
	case operatorstypes.Operator:
		return k.OperatorDelegatorStartingInfos.Remove(ctx, collections.Join(target.GetID(), delegator))
	case servicestypes.Service:
		return k.ServiceDelegatorStartingInfos.Remove(ctx, collections.Join(target.GetID(), delegator))
	default:
		return errors.Wrapf(restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target)
	}
}

// --------------------------------------------------------------------------------------------------------------------

// SetOutstandingRewards sets the outstanding rewards for a target.
func (k *Keeper) SetOutstandingRewards(
	ctx context.Context, target restakingtypes.DelegationTarget, rewards types.OutstandingRewards,
) error {
	switch target.(type) {
	case poolstypes.Pool:
		return k.PoolOutstandingRewards.Set(ctx, target.GetID(), rewards)
	case operatorstypes.Operator:
		return k.OperatorOutstandingRewards.Set(ctx, target.GetID(), rewards)
	case servicestypes.Service:
		return k.ServiceOutstandingRewards.Set(ctx, target.GetID(), rewards)
	default:
		return errors.Wrapf(restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target)
	}
}

// GetOutstandingRewards returns the outstanding rewards for a target
func (k *Keeper) GetOutstandingRewards(ctx context.Context, target restakingtypes.DelegationTarget) (types.OutstandingRewards, error) {
	switch target.(type) {
	case poolstypes.Pool:
		return k.PoolOutstandingRewards.Get(ctx, target.GetID())
	case operatorstypes.Operator:
		return k.OperatorOutstandingRewards.Get(ctx, target.GetID())
	case servicestypes.Service:
		return k.ServiceOutstandingRewards.Get(ctx, target.GetID())
	default:
		return types.OutstandingRewards{}, errors.Wrapf(
			restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target,
		)
	}
}

// --------------------------------------------------------------------------------------------------------------------

// SetCurrentRewards sets the current rewards for a target
func (k *Keeper) SetCurrentRewards(
	ctx context.Context, target restakingtypes.DelegationTarget, rewards types.CurrentRewards,
) error {
	switch target.(type) {
	case poolstypes.Pool:
		return k.PoolCurrentRewards.Set(ctx, target.GetID(), rewards)
	case operatorstypes.Operator:
		return k.OperatorCurrentRewards.Set(ctx, target.GetID(), rewards)
	case servicestypes.Service:
		return k.ServiceCurrentRewards.Set(ctx, target.GetID(), rewards)
	default:
		return errors.Wrapf(restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target)
	}
}

// HasCurrentRewards returns true if the current rewards exist for a target
func (k *Keeper) HasCurrentRewards(ctx context.Context, target restakingtypes.DelegationTarget) (bool, error) {
	switch target.(type) {
	case poolstypes.Pool:
		return k.PoolCurrentRewards.Has(ctx, target.GetID())
	case operatorstypes.Operator:
		return k.OperatorCurrentRewards.Has(ctx, target.GetID())
	case servicestypes.Service:
		return k.ServiceCurrentRewards.Has(ctx, target.GetID())
	default:
		return false, errors.Wrapf(restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target)
	}
}

// GetCurrentRewards returns the current rewards for a target
func (k *Keeper) GetCurrentRewards(ctx context.Context, target restakingtypes.DelegationTarget) (types.CurrentRewards, error) {
	switch target.(type) {
	case poolstypes.Pool:
		return k.PoolCurrentRewards.Get(ctx, target.GetID())
	case operatorstypes.Operator:
		return k.OperatorCurrentRewards.Get(ctx, target.GetID())
	case servicestypes.Service:
		return k.ServiceCurrentRewards.Get(ctx, target.GetID())
	default:
		return types.CurrentRewards{}, errors.Wrapf(
			restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target,
		)
	}
}

// DeleteCurrentRewards deletes the current rewards for a target
func (k *Keeper) DeleteCurrentRewards(ctx context.Context, target restakingtypes.DelegationTarget) error {
	switch target.(type) {
	case poolstypes.Pool:
		return k.PoolCurrentRewards.Remove(ctx, target.GetID())
	case operatorstypes.Operator:
		return k.OperatorCurrentRewards.Remove(ctx, target.GetID())
	case servicestypes.Service:
		return k.ServiceCurrentRewards.Remove(ctx, target.GetID())
	default:
		return errors.Wrapf(restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target)
	}
}

// --------------------------------------------------------------------------------------------------------------------

// SetHistoricalRewards sets the historical rewards for a target and period
func (k *Keeper) SetHistoricalRewards(
	ctx context.Context, target restakingtypes.DelegationTarget, period uint64,
	rewards types.HistoricalRewards,
) error {
	switch target.(type) {
	case poolstypes.Pool:
		return k.PoolHistoricalRewards.Set(ctx, collections.Join(target.GetID(), period), rewards)
	case operatorstypes.Operator:
		return k.OperatorHistoricalRewards.Set(ctx, collections.Join(target.GetID(), period), rewards)
	case servicestypes.Service:
		return k.ServiceHistoricalRewards.Set(ctx, collections.Join(target.GetID(), period), rewards)
	default:
		return errors.Wrapf(restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target)
	}
}

// GetHistoricalRewards returns the historical rewards for a target and period
func (k *Keeper) GetHistoricalRewards(
	ctx context.Context, target restakingtypes.DelegationTarget, period uint64,
) (types.HistoricalRewards, error) {
	switch target.(type) {
	case poolstypes.Pool:
		return k.PoolHistoricalRewards.Get(ctx, collections.Join(target.GetID(), period))
	case operatorstypes.Operator:
		return k.OperatorHistoricalRewards.Get(ctx, collections.Join(target.GetID(), period))
	case servicestypes.Service:
		return k.ServiceHistoricalRewards.Get(ctx, collections.Join(target.GetID(), period))
	default:
		return types.HistoricalRewards{}, errors.Wrapf(
			restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target,
		)
	}
}

// RemoveHistoricalRewards removes the historical rewards for a target and period
func (k *Keeper) RemoveHistoricalRewards(
	ctx context.Context, target restakingtypes.DelegationTarget, period uint64,
) error {
	switch target.(type) {
	case poolstypes.Pool:
		return k.PoolHistoricalRewards.Remove(ctx, collections.Join(target.GetID(), period))
	case operatorstypes.Operator:
		return k.OperatorHistoricalRewards.Remove(ctx, collections.Join(target.GetID(), period))
	case servicestypes.Service:
		return k.ServiceHistoricalRewards.Remove(ctx, collections.Join(target.GetID(), period))
	default:
		return errors.Wrapf(restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target)
	}
}

// DeleteHistoricalRewards deletes all historical rewards for a target
func (k *Keeper) DeleteHistoricalRewards(ctx context.Context, target restakingtypes.DelegationTarget) error {
	var collection collections.Map[collections.Pair[uint32, uint64], types.HistoricalRewards]
	switch target.(type) {
	case poolstypes.Pool:
		collection = k.PoolHistoricalRewards
	case operatorstypes.Operator:
		collection = k.OperatorHistoricalRewards
	case servicestypes.Service:
		collection = k.ServiceHistoricalRewards
	default:
		return errors.Wrapf(restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target)
	}

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

// GetOutstandingRewardsCoins returns the outstanding rewards coins for a target
func (k *Keeper) GetOutstandingRewardsCoins(ctx context.Context, target restakingtypes.DelegationTarget) (types.DecPools, error) {
	var (
		rewards types.OutstandingRewards
		err     error
	)
	switch target.(type) {
	case poolstypes.Pool:
		rewards, err = k.PoolOutstandingRewards.Get(ctx, target.GetID())
	case operatorstypes.Operator:
		rewards, err = k.OperatorOutstandingRewards.Get(ctx, target.GetID())
	case servicestypes.Service:
		rewards, err = k.ServiceOutstandingRewards.Get(ctx, target.GetID())
	default:
		return nil, errors.Wrapf(restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target)
	}
	if err != nil && !errors.IsOf(err, collections.ErrNotFound) {
		return nil, err
	}
	return rewards.Rewards, nil
}

// DeleteOutstandingRewards deletes the outstanding rewards for a target
func (k *Keeper) DeleteOutstandingRewards(ctx context.Context, target restakingtypes.DelegationTarget) error {
	switch target.(type) {
	case poolstypes.Pool:
		return k.PoolOutstandingRewards.Remove(ctx, target.GetID())
	case operatorstypes.Operator:
		return k.OperatorOutstandingRewards.Remove(ctx, target.GetID())
	case servicestypes.Service:
		return k.ServiceOutstandingRewards.Remove(ctx, target.GetID())
	default:
		return errors.Wrapf(restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target)
	}
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

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	delegator, err := k.accountKeeper.AddressCodec().BytesToString(delAddr)
	if err != nil {
		return nil, err
	}

	delegation, found, err := k.restakingKeeper.GetDelegationForTarget(sdkCtx, target, delegator)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, errors.Wrap(sdkerrors.ErrNotFound, "delegation not found")
	}

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
