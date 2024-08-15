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

	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
)

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

func (k *Keeper) SetLastRewardsAllocationTime(ctx context.Context, t time.Time) error {
	ts, err := gogotypes.TimestampProto(t)
	if err != nil {
		return err
	}
	return k.LastRewardsAllocationTime.Set(ctx, *ts)
}

func (k *Keeper) GetRewardsPlan(ctx context.Context, planID uint64) (types.RewardsPlan, error) {
	return k.RewardsPlans.Get(ctx, planID)
}

func (k *Keeper) SetDelegatorStartingInfo(
	ctx context.Context, target types.DelegationTarget, del sdk.AccAddress, info types.DelegatorStartingInfo,
) error {
	switch target.Type() {
	case restakingtypes.DELEGATION_TYPE_POOL:
		return k.PoolDelegatorStartingInfos.Set(ctx, collections.Join(target.GetID(), del), info)
	case restakingtypes.DELEGATION_TYPE_OPERATOR:
		return k.OperatorDelegatorStartingInfos.Set(ctx, collections.Join(target.GetID(), del), info)
	case restakingtypes.DELEGATION_TYPE_SERVICE:
		return k.ServiceDelegatorStartingInfos.Set(ctx, collections.Join(target.GetID(), del), info)
	default:
		panic("unknown delegation type")
	}
}

func (k *Keeper) GetDelegatorStartingInfo(ctx context.Context, target types.DelegationTarget, del sdk.AccAddress) (types.DelegatorStartingInfo, error) {
	switch target.Type() {
	case restakingtypes.DELEGATION_TYPE_POOL:
		return k.PoolDelegatorStartingInfos.Get(ctx, collections.Join(target.GetID(), del))
	case restakingtypes.DELEGATION_TYPE_OPERATOR:
		return k.OperatorDelegatorStartingInfos.Get(ctx, collections.Join(target.GetID(), del))
	case restakingtypes.DELEGATION_TYPE_SERVICE:
		return k.ServiceDelegatorStartingInfos.Get(ctx, collections.Join(target.GetID(), del))
	default:
		panic("unknown delegation type")
	}
}

func (k *Keeper) HasDelegatorStartingInfo(ctx context.Context, target types.DelegationTarget, del sdk.AccAddress) (bool, error) {
	switch target.Type() {
	case restakingtypes.DELEGATION_TYPE_POOL:
		return k.PoolDelegatorStartingInfos.Has(ctx, collections.Join(target.GetID(), del))
	case restakingtypes.DELEGATION_TYPE_OPERATOR:
		return k.OperatorDelegatorStartingInfos.Has(ctx, collections.Join(target.GetID(), del))
	case restakingtypes.DELEGATION_TYPE_SERVICE:
		return k.ServiceDelegatorStartingInfos.Has(ctx, collections.Join(target.GetID(), del))
	default:
		panic("unknown delegation type")
	}
}

func (k *Keeper) RemoveDelegatorStartingInfo(ctx context.Context, target types.DelegationTarget, del sdk.AccAddress) error {
	switch target.Type() {
	case restakingtypes.DELEGATION_TYPE_POOL:
		return k.PoolDelegatorStartingInfos.Remove(ctx, collections.Join(target.GetID(), del))
	case restakingtypes.DELEGATION_TYPE_OPERATOR:
		return k.OperatorDelegatorStartingInfos.Remove(ctx, collections.Join(target.GetID(), del))
	case restakingtypes.DELEGATION_TYPE_SERVICE:
		return k.ServiceDelegatorStartingInfos.Remove(ctx, collections.Join(target.GetID(), del))
	default:
		panic("unknown delegation type")
	}
}

func (k *Keeper) SetOutstandingRewards(
	ctx context.Context, target types.DelegationTarget, rewards types.OutstandingRewards,
) error {
	switch target.Type() {
	case restakingtypes.DELEGATION_TYPE_POOL:
		return k.PoolOutstandingRewards.Set(ctx, target.GetID(), rewards)
	case restakingtypes.DELEGATION_TYPE_OPERATOR:
		return k.OperatorOutstandingRewards.Set(ctx, target.GetID(), rewards)
	case restakingtypes.DELEGATION_TYPE_SERVICE:
		return k.ServiceOutstandingRewards.Set(ctx, target.GetID(), rewards)
	default:
		panic("unknown delegation type")
	}
}

func (k *Keeper) GetOutstandingRewards(ctx context.Context, target types.DelegationTarget) (types.OutstandingRewards, error) {
	switch target.Type() {
	case restakingtypes.DELEGATION_TYPE_POOL:
		return k.PoolOutstandingRewards.Get(ctx, target.GetID())
	case restakingtypes.DELEGATION_TYPE_OPERATOR:
		return k.OperatorOutstandingRewards.Get(ctx, target.GetID())
	case restakingtypes.DELEGATION_TYPE_SERVICE:
		return k.ServiceOutstandingRewards.Get(ctx, target.GetID())
	default:
		panic("unknown delegation type")
	}
}

func (k *Keeper) GetCurrentRewards(ctx context.Context, target types.DelegationTarget) (types.CurrentRewards, error) {
	switch target.Type() {
	case restakingtypes.DELEGATION_TYPE_POOL:
		return k.PoolCurrentRewards.Get(ctx, target.GetID())
	case restakingtypes.DELEGATION_TYPE_OPERATOR:
		return k.OperatorCurrentRewards.Get(ctx, target.GetID())
	case restakingtypes.DELEGATION_TYPE_SERVICE:
		return k.ServiceCurrentRewards.Get(ctx, target.GetID())
	default:
		panic("unknown delegation type")
	}
}

func (k *Keeper) SetCurrentRewards(
	ctx context.Context, target types.DelegationTarget, rewards types.CurrentRewards,
) error {
	switch target.Type() {
	case restakingtypes.DELEGATION_TYPE_POOL:
		return k.PoolCurrentRewards.Set(ctx, target.GetID(), rewards)
	case restakingtypes.DELEGATION_TYPE_OPERATOR:
		return k.OperatorCurrentRewards.Set(ctx, target.GetID(), rewards)
	case restakingtypes.DELEGATION_TYPE_SERVICE:
		return k.ServiceCurrentRewards.Set(ctx, target.GetID(), rewards)
	default:
		panic("unknown delegation type")
	}
}

func (k *Keeper) HasCurrentRewards(ctx context.Context, target types.DelegationTarget) (bool, error) {
	switch target.Type() {
	case restakingtypes.DELEGATION_TYPE_POOL:
		return k.PoolCurrentRewards.Has(ctx, target.GetID())
	case restakingtypes.DELEGATION_TYPE_OPERATOR:
		return k.OperatorCurrentRewards.Has(ctx, target.GetID())
	case restakingtypes.DELEGATION_TYPE_SERVICE:
		return k.ServiceCurrentRewards.Has(ctx, target.GetID())
	default:
		panic("unknown delegation type")
	}
}

func (k *Keeper) GetHistoricalRewards(
	ctx context.Context, target types.DelegationTarget, period uint64,
) (types.HistoricalRewards, error) {
	switch target.Type() {
	case restakingtypes.DELEGATION_TYPE_POOL:
		return k.PoolHistoricalRewards.Get(ctx, collections.Join(target.GetID(), period))
	case restakingtypes.DELEGATION_TYPE_OPERATOR:
		return k.OperatorHistoricalRewards.Get(ctx, collections.Join(target.GetID(), period))
	case restakingtypes.DELEGATION_TYPE_SERVICE:
		return k.ServiceHistoricalRewards.Get(ctx, collections.Join(target.GetID(), period))
	default:
		panic("unknown delegation type")
	}
}

func (k *Keeper) SetHistoricalRewards(
	ctx context.Context, target types.DelegationTarget, period uint64,
	rewards types.HistoricalRewards,
) error {
	switch target.Type() {
	case restakingtypes.DELEGATION_TYPE_POOL:
		return k.PoolHistoricalRewards.Set(ctx, collections.Join(target.GetID(), period), rewards)
	case restakingtypes.DELEGATION_TYPE_OPERATOR:
		return k.OperatorHistoricalRewards.Set(ctx, collections.Join(target.GetID(), period), rewards)
	case restakingtypes.DELEGATION_TYPE_SERVICE:
		return k.ServiceHistoricalRewards.Set(ctx, collections.Join(target.GetID(), period), rewards)
	default:
		panic("unknown delegation type")
	}
}

func (k *Keeper) RemoveHistoricalRewards(
	ctx context.Context, target types.DelegationTarget, period uint64,
) error {
	switch target.Type() {
	case restakingtypes.DELEGATION_TYPE_POOL:
		return k.PoolHistoricalRewards.Remove(ctx, collections.Join(target.GetID(), period))
	case restakingtypes.DELEGATION_TYPE_OPERATOR:
		return k.OperatorHistoricalRewards.Remove(ctx, collections.Join(target.GetID(), period))
	case restakingtypes.DELEGATION_TYPE_SERVICE:
		return k.ServiceHistoricalRewards.Remove(ctx, collections.Join(target.GetID(), period))
	default:
		panic("unknown delegation type")
	}
}

func (k *Keeper) GetOutstandingRewardsCoins(ctx context.Context, target types.DelegationTarget) (types.DecPools, error) {
	var (
		rewards types.OutstandingRewards
		err     error
	)
	switch target.Type() {
	case restakingtypes.DELEGATION_TYPE_POOL:
		rewards, err = k.PoolOutstandingRewards.Get(ctx, target.GetID())
	case restakingtypes.DELEGATION_TYPE_OPERATOR:
		rewards, err = k.OperatorOutstandingRewards.Get(ctx, target.GetID())
	case restakingtypes.DELEGATION_TYPE_SERVICE:
		rewards, err = k.ServiceOutstandingRewards.Get(ctx, target.GetID())
	default:
		panic("unknown delegation type")
	}
	if err != nil && !errors.IsOf(err, collections.ErrNotFound) {
		return nil, err
	}
	return rewards.Rewards, nil
}

// get accumulated commission for an operator
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

// get the delegator withdraw address, defaulting to the delegator address
func (k *Keeper) GetDelegatorWithdrawAddr(ctx context.Context, delAddr sdk.AccAddress) (sdk.AccAddress, error) {
	addr, err := k.DelegatorWithdrawAddrs.Get(ctx, delAddr)
	if err != nil && errors.IsOf(err, collections.ErrNotFound) {
		return delAddr, nil
	}
	return addr, err
}

func (k *Keeper) DelegationRewards(
	ctx context.Context, delAddr sdk.AccAddress, delType restakingtypes.DelegationType, targetID uint32,
) (types.DecPools, error) {
	target, err := k.GetDelegationTarget(ctx, delType, targetID)
	if err != nil {
		return nil, err
	}

	del, found := k.GetDelegation(ctx, target, delAddr)
	if !found {
		return nil, errors.Wrap(sdkerrors.ErrNotFound, "delegation not found")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
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

func (k *Keeper) PoolDelegationRewards(ctx context.Context, delAddr sdk.AccAddress, poolID uint32) (types.DecPools, error) {
	return k.DelegationRewards(ctx, delAddr, restakingtypes.DELEGATION_TYPE_POOL, poolID)
}

func (k *Keeper) OperatorDelegationRewards(ctx context.Context, delAddr sdk.AccAddress, operatorID uint32) (types.DecPools, error) {
	return k.DelegationRewards(ctx, delAddr, restakingtypes.DELEGATION_TYPE_OPERATOR, operatorID)
}

func (k *Keeper) ServiceDelegationRewards(ctx context.Context, delAddr sdk.AccAddress, serviceID uint32) (types.DecPools, error) {
	return k.DelegationRewards(ctx, delAddr, restakingtypes.DELEGATION_TYPE_SERVICE, serviceID)
}

// createAccountIfNotExists creates an account if it does not exist
func (k *Keeper) createAccountIfNotExists(ctx context.Context, address sdk.AccAddress) {
	if !k.accountKeeper.HasAccount(ctx, address) {
		defer telemetry.IncrCounter(1, "new", "account")
		k.accountKeeper.SetAccount(ctx, k.accountKeeper.NewAccountWithAddress(ctx, address))
	}
}
