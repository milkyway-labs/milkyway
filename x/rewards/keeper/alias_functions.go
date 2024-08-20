package keeper

import (
	"context"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gogotypes "github.com/cosmos/gogoproto/types"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

// GetLastRewardsAllocationTime returns the last time rewards were allocated.
// If there's no last rewards allocation time set yet, nil is returned.
func (k *Keeper) GetLastRewardsAllocationTime(ctx context.Context) (*time.Time, error) {
	ts, err := k.LastRewardsAllocationTime.Get(ctx)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	t, err := gogotypes.TimestampFromProto(&ts)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// SetLastRewardsAllocationTime sets the last time rewards were allocated.
func (k *Keeper) SetLastRewardsAllocationTime(ctx context.Context, t time.Time) error {
	ts, err := gogotypes.TimestampProto(t)
	if err != nil {
		return err
	}
	return k.LastRewardsAllocationTime.Set(ctx, *ts)
}

// GetRewardsPlan returns a rewards plan by ID.
func (k *Keeper) GetRewardsPlan(ctx context.Context, planID uint64) (types.RewardsPlan, error) {
	return k.RewardsPlans.Get(ctx, planID)
}

// SetDelegatorStartingInfo sets the delegator starting info for a delegator.
func (k *Keeper) SetDelegatorStartingInfo(
	ctx context.Context, target restakingtypes.DelegationTarget, del sdk.AccAddress, info types.DelegatorStartingInfo,
) error {
	switch target.(type) {
	case *poolstypes.Pool:
		return k.PoolDelegatorStartingInfos.Set(ctx, collections.Join(target.GetID(), del), info)
	case *operatorstypes.Operator:
		return k.OperatorDelegatorStartingInfos.Set(ctx, collections.Join(target.GetID(), del), info)
	case *servicestypes.Service:
		return k.ServiceDelegatorStartingInfos.Set(ctx, collections.Join(target.GetID(), del), info)
	default:
		return errors.Wrapf(restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target)
	}
}

func (k *Keeper) GetDelegatorStartingInfo(ctx context.Context, target restakingtypes.DelegationTarget, del sdk.AccAddress) (types.DelegatorStartingInfo, error) {
	switch target.(type) {
	case *poolstypes.Pool:
		return k.PoolDelegatorStartingInfos.Get(ctx, collections.Join(target.GetID(), del))
	case *operatorstypes.Operator:
		return k.OperatorDelegatorStartingInfos.Get(ctx, collections.Join(target.GetID(), del))
	case *servicestypes.Service:
		return k.ServiceDelegatorStartingInfos.Get(ctx, collections.Join(target.GetID(), del))
	default:
		return types.DelegatorStartingInfo{}, errors.Wrapf(
			restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target,
		)
	}
}

func (k *Keeper) HasDelegatorStartingInfo(ctx context.Context, target restakingtypes.DelegationTarget, del sdk.AccAddress) (bool, error) {
	switch target.(type) {
	case *poolstypes.Pool:
		return k.PoolDelegatorStartingInfos.Has(ctx, collections.Join(target.GetID(), del))
	case *operatorstypes.Operator:
		return k.OperatorDelegatorStartingInfos.Has(ctx, collections.Join(target.GetID(), del))
	case *servicestypes.Service:
		return k.ServiceDelegatorStartingInfos.Has(ctx, collections.Join(target.GetID(), del))
	default:
		return false, errors.Wrapf(restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target)
	}
}

func (k *Keeper) RemoveDelegatorStartingInfo(ctx context.Context, target restakingtypes.DelegationTarget, del sdk.AccAddress) error {
	switch target.(type) {
	case *poolstypes.Pool:
		return k.PoolDelegatorStartingInfos.Remove(ctx, collections.Join(target.GetID(), del))
	case *operatorstypes.Operator:
		return k.OperatorDelegatorStartingInfos.Remove(ctx, collections.Join(target.GetID(), del))
	case *servicestypes.Service:
		return k.ServiceDelegatorStartingInfos.Remove(ctx, collections.Join(target.GetID(), del))
	default:
		return errors.Wrapf(restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target)
	}
}

func (k *Keeper) SetOutstandingRewards(
	ctx context.Context, target restakingtypes.DelegationTarget, rewards types.OutstandingRewards,
) error {
	switch target.(type) {
	case *poolstypes.Pool:
		return k.PoolOutstandingRewards.Set(ctx, target.GetID(), rewards)
	case *operatorstypes.Operator:
		return k.OperatorOutstandingRewards.Set(ctx, target.GetID(), rewards)
	case *servicestypes.Service:
		return k.ServiceOutstandingRewards.Set(ctx, target.GetID(), rewards)
	default:
		return errors.Wrapf(restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target)
	}
}

func (k *Keeper) GetOutstandingRewards(ctx context.Context, target restakingtypes.DelegationTarget) (types.OutstandingRewards, error) {
	switch target.(type) {
	case *poolstypes.Pool:
		return k.PoolOutstandingRewards.Get(ctx, target.GetID())
	case *operatorstypes.Operator:
		return k.OperatorOutstandingRewards.Get(ctx, target.GetID())
	case *servicestypes.Service:
		return k.ServiceOutstandingRewards.Get(ctx, target.GetID())
	default:
		return types.OutstandingRewards{}, errors.Wrapf(
			restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target,
		)
	}
}

func (k *Keeper) GetCurrentRewards(ctx context.Context, target restakingtypes.DelegationTarget) (types.CurrentRewards, error) {
	switch target.(type) {
	case *poolstypes.Pool:
		return k.PoolCurrentRewards.Get(ctx, target.GetID())
	case *operatorstypes.Operator:
		return k.OperatorCurrentRewards.Get(ctx, target.GetID())
	case *servicestypes.Service:
		return k.ServiceCurrentRewards.Get(ctx, target.GetID())
	default:
		return types.CurrentRewards{}, errors.Wrapf(
			restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target,
		)
	}
}

func (k *Keeper) SetCurrentRewards(
	ctx context.Context, target restakingtypes.DelegationTarget, rewards types.CurrentRewards,
) error {
	switch target.(type) {
	case *poolstypes.Pool:
		return k.PoolCurrentRewards.Set(ctx, target.GetID(), rewards)
	case *operatorstypes.Operator:
		return k.OperatorCurrentRewards.Set(ctx, target.GetID(), rewards)
	case *servicestypes.Service:
		return k.ServiceCurrentRewards.Set(ctx, target.GetID(), rewards)
	default:
		return errors.Wrapf(restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target)
	}
}

func (k *Keeper) HasCurrentRewards(ctx context.Context, target restakingtypes.DelegationTarget) (bool, error) {
	switch target.(type) {
	case *poolstypes.Pool:
		return k.PoolCurrentRewards.Has(ctx, target.GetID())
	case *operatorstypes.Operator:
		return k.OperatorCurrentRewards.Has(ctx, target.GetID())
	case *servicestypes.Service:
		return k.ServiceCurrentRewards.Has(ctx, target.GetID())
	default:
		return false, errors.Wrapf(restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target)
	}
}

func (k *Keeper) GetHistoricalRewards(
	ctx context.Context, target restakingtypes.DelegationTarget, period uint64,
) (types.HistoricalRewards, error) {
	switch target.(type) {
	case *poolstypes.Pool:
		return k.PoolHistoricalRewards.Get(ctx, collections.Join(target.GetID(), period))
	case *operatorstypes.Operator:
		return k.OperatorHistoricalRewards.Get(ctx, collections.Join(target.GetID(), period))
	case *servicestypes.Service:
		return k.ServiceHistoricalRewards.Get(ctx, collections.Join(target.GetID(), period))
	default:
		return types.HistoricalRewards{}, errors.Wrapf(
			restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target,
		)
	}
}

func (k *Keeper) SetHistoricalRewards(
	ctx context.Context, target restakingtypes.DelegationTarget, period uint64,
	rewards types.HistoricalRewards,
) error {
	switch target.(type) {
	case *poolstypes.Pool:
		return k.PoolHistoricalRewards.Set(ctx, collections.Join(target.GetID(), period), rewards)
	case *operatorstypes.Operator:
		return k.OperatorHistoricalRewards.Set(ctx, collections.Join(target.GetID(), period), rewards)
	case *servicestypes.Service:
		return k.ServiceHistoricalRewards.Set(ctx, collections.Join(target.GetID(), period), rewards)
	default:
		return errors.Wrapf(restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target)
	}
}

func (k *Keeper) RemoveHistoricalRewards(
	ctx context.Context, target restakingtypes.DelegationTarget, period uint64,
) error {
	switch target.(type) {
	case *poolstypes.Pool:
		return k.PoolHistoricalRewards.Remove(ctx, collections.Join(target.GetID(), period))
	case *operatorstypes.Operator:
		return k.OperatorHistoricalRewards.Remove(ctx, collections.Join(target.GetID(), period))
	case *servicestypes.Service:
		return k.ServiceHistoricalRewards.Remove(ctx, collections.Join(target.GetID(), period))
	default:
		return errors.Wrapf(restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target)
	}
}

func (k *Keeper) GetOutstandingRewardsCoins(ctx context.Context, target restakingtypes.DelegationTarget) (types.DecPools, error) {
	var (
		rewards types.OutstandingRewards
		err     error
	)
	switch target.(type) {
	case *poolstypes.Pool:
		rewards, err = k.PoolOutstandingRewards.Get(ctx, target.GetID())
	case *operatorstypes.Operator:
		rewards, err = k.OperatorOutstandingRewards.Get(ctx, target.GetID())
	case *servicestypes.Service:
		rewards, err = k.ServiceOutstandingRewards.Get(ctx, target.GetID())
	default:
		return nil, errors.Wrapf(restakingtypes.ErrInvalidDelegationType, "invalid delegation target type %T", target)
	}
	if err != nil && !errors.IsOf(err, collections.ErrNotFound) {
		return nil, err
	}
	return rewards.Rewards, nil
}

// GetOperatorAccumulatedCommission returns the accumulated commission for an
// operator.
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

// GetDelegatorWithdrawAddr returns the delegator's withdraw address if set,
// otherwise the delegator address is returned.
func (k *Keeper) GetDelegatorWithdrawAddr(ctx context.Context, delAddr sdk.AccAddress) (sdk.AccAddress, error) {
	addr, err := k.DelegatorWithdrawAddrs.Get(ctx, delAddr)
	if err != nil && errors.IsOf(err, collections.ErrNotFound) {
		return delAddr, nil
	}
	return addr, err
}

// GetDelegationRewards returns the rewards for a delegation within a cached
// context.
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

	del, found := k.restakingKeeper.GetDelegationForTarget(sdkCtx, target, delegator)
	if !found {
		return nil, errors.Wrap(sdkerrors.ErrNotFound, "delegation not found")
	}

	cacheCtx, _ := sdkCtx.CacheContext()
	endingPeriod, err := k.IncrementDelegationTargetPeriod(cacheCtx, target)
	if err != nil {
		return nil, err
	}

	rewards, err := k.CalculateDelegationRewards(cacheCtx, target, del, endingPeriod)
	if err != nil {
		return nil, err
	}

	return rewards, nil
}

// GetPoolDelegationRewards returns the rewards for a pool delegation within a
// cached context.
func (k *Keeper) GetPoolDelegationRewards(
	ctx context.Context, delAddr sdk.AccAddress, poolID uint32,
) (types.DecPools, error) {
	return k.GetDelegationRewards(ctx, delAddr, restakingtypes.DELEGATION_TYPE_POOL, poolID)
}

// GetOperatorDelegationRewards returns the rewards for an operator delegation
// within a cached context.
func (k *Keeper) GetOperatorDelegationRewards(
	ctx context.Context, delAddr sdk.AccAddress, operatorID uint32,
) (types.DecPools, error) {
	return k.GetDelegationRewards(ctx, delAddr, restakingtypes.DELEGATION_TYPE_OPERATOR, operatorID)
}

// GetServiceDelegationRewards returns the rewards for a service delegation
// within a cached context.
func (k *Keeper) GetServiceDelegationRewards(
	ctx context.Context, delAddr sdk.AccAddress, serviceID uint32,
) (types.DecPools, error) {
	return k.GetDelegationRewards(ctx, delAddr, restakingtypes.DELEGATION_TYPE_SERVICE, serviceID)
}

// createAccountIfNotExists creates an account if it does not exist.
func (k *Keeper) createAccountIfNotExists(ctx context.Context, address sdk.AccAddress) {
	if !k.accountKeeper.HasAccount(ctx, address) {
		defer telemetry.IncrCounter(1, "new", "account")
		k.accountKeeper.SetAccount(ctx, k.accountKeeper.NewAccountWithAddress(ctx, address))
	}
}
