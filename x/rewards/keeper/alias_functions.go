package keeper

import (
	"context"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"

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

// get outstanding rewards
func (k *Keeper) GetPoolOutstandingRewardsCoins(ctx context.Context, poolID uint32) (sdk.DecCoins, error) {
	rewards, err := k.PoolOutstandingRewards.Get(ctx, poolID)
	if err != nil {
		return nil, err
	}

	return rewards.Rewards, nil
}

// get outstanding rewards
func (k *Keeper) GetOperatorOutstandingRewardsCoins(ctx context.Context, operatorID uint32) (types.DecPools, error) {
	rewards, err := k.OperatorOutstandingRewards.Get(ctx, operatorID)
	if err != nil {
		return nil, err
	}

	return rewards.Rewards, nil
}

// get outstanding rewards
func (k *Keeper) GetServiceOutstandingRewardsCoins(ctx context.Context, serviceID uint32) (types.DecPools, error) {
	rewards, err := k.ServiceOutstandingRewards.Get(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	return rewards.Rewards, nil
}

// get the delegator withdraw address, defaulting to the delegator address
// Note that it always returns the delegator address itself for now
// TODO: make it possible to set different address then the delegator address?
func (k *Keeper) GetDelegatorWithdrawAddr(_ context.Context, delAddr string) (sdk.AccAddress, error) {
	return sdk.AccAddressFromBech32(delAddr)
}

func (k *Keeper) PoolDelegationRewards(ctx context.Context, delegator string, poolID uint32) (sdk.DecCoins, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cacheCtx, _ := sdkCtx.CacheContext()
	pool, found := k.poolsKeeper.GetPool(cacheCtx, poolID)
	if !found {
		return nil, nil
	}
	del, found := k.restakingKeeper.GetPoolDelegation(cacheCtx, poolID, delegator)
	if !found {
		return nil, nil
	}
	endingPeriod, err := k.IncrementPoolPeriod(cacheCtx, pool)
	if err != nil {
		return nil, err
	}
	return k.CalculatePoolDelegationRewards(cacheCtx, pool, del, endingPeriod)
}

func (k *Keeper) OperatorDelegationRewards(ctx context.Context, delegator string, operatorID uint32) (sdk.DecCoins, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cacheCtx, _ := sdkCtx.CacheContext()
	operator, found := k.operatorsKeeper.GetOperator(cacheCtx, operatorID)
	if !found {
		return nil, nil
	}
	del, found := k.restakingKeeper.GetOperatorDelegation(cacheCtx, operatorID, delegator)
	if !found {
		return nil, nil
	}
	endingPeriod, err := k.IncrementOperatorPeriod(cacheCtx, operator)
	if err != nil {
		return nil, err
	}
	rewards, err := k.CalculateOperatorDelegationRewards(cacheCtx, operator, del, endingPeriod)
	if err != nil {
		return nil, err
	}
	rewardsSum := sdk.DecCoins{}
	for _, decPool := range rewards {
		rewardsSum = rewardsSum.Add(decPool.DecCoins...)
	}
	return rewardsSum, nil
}

func (k *Keeper) ServiceDelegationRewards(ctx context.Context, delegator string, serviceID uint32) (sdk.DecCoins, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cacheCtx, _ := sdkCtx.CacheContext()
	service, found := k.servicesKeeper.GetService(cacheCtx, serviceID)
	if !found {
		return nil, nil
	}
	del, found := k.restakingKeeper.GetServiceDelegation(cacheCtx, serviceID, delegator)
	if !found {
		return nil, nil
	}
	endingPeriod, err := k.IncrementServicePeriod(cacheCtx, service)
	if err != nil {
		return nil, err
	}
	rewards, err := k.CalculateServiceDelegationRewards(cacheCtx, service, del, endingPeriod)
	if err != nil {
		return nil, err
	}
	rewardsSum := sdk.DecCoins{}
	for _, decPool := range rewards {
		rewardsSum = rewardsSum.Add(decPool.DecCoins...)
	}
	return rewardsSum, nil
}
