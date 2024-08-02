package keeper

import (
	"context"
	"errors"
	"time"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"

	"github.com/milkyway-labs/milkyway/x/rewards/types"
)

func (k *Keeper) GetLastRewardsAllocationTime(ctx context.Context) (*time.Time, error) {
	ts, err := k.LastRewardsAllocationTime.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
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
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return nil, err
	}

	return rewards.Rewards, nil
}

// get outstanding rewards
func (k *Keeper) GetOperatorOutstandingRewardsCoins(ctx context.Context, operatorID uint32) (types.DecPools, error) {
	rewards, err := k.OperatorOutstandingRewards.Get(ctx, operatorID)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return nil, err
	}

	return rewards.Rewards, nil
}

// get outstanding rewards
func (k *Keeper) GetServiceOutstandingRewardsCoins(ctx context.Context, serviceID uint32) (types.DecPools, error) {
	rewards, err := k.ServiceOutstandingRewards.Get(ctx, serviceID)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return nil, err
	}

	return rewards.Rewards, nil
}

// get accumulated commission for an operator
func (k Keeper) GetOperatorAccumulatedCommission(ctx context.Context, operatorID uint32) (commission types.MultiAccumulatedCommission, err error) {
	commission, err = k.OperatorAccumulatedCommissions.Get(ctx, operatorID)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return types.MultiAccumulatedCommission{}, nil
		}
		return types.MultiAccumulatedCommission{}, err
	}
	return
}

// get the delegator withdraw address, defaulting to the delegator address
func (k *Keeper) GetDelegatorWithdrawAddr(ctx context.Context, delAddr sdk.AccAddress) (sdk.AccAddress, error) {
	addr, err := k.DelegatorWithdrawAddrs.Get(ctx, delAddr)
	if err != nil && errors.Is(err, collections.ErrNotFound) {
		return delAddr, nil
	}
	return addr, err
}
