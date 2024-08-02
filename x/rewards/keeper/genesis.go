package keeper

import (
	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/rewards/types"
)

// ExportGenesis returns the GenesisState associated with the given context
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params, err := k.Params.Get(ctx)
	if err != nil {
		panic(err)
	}

	nextRewardsPlanID, err := k.NextRewardsPlanID.Get(ctx)
	if err != nil {
		panic(err)
	}

	rewardsPlans := []types.RewardsPlan{}
	err = k.RewardsPlans.Walk(ctx, nil, func(id uint64, plan types.RewardsPlan) (stop bool, err error) {
		rewardsPlans = append(rewardsPlans, plan)
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	lastRewardsAllocationTime, err := k.GetLastRewardsAllocationTime(ctx)
	if err != nil {
		panic(err)
	}

	delegatorWithdrawInfos := []types.DelegatorWithdrawInfo{}
	err = k.DelegatorWithdrawAddrs.Walk(ctx, nil, func(delAddr, withdrawAddr sdk.AccAddress) (stop bool, err error) {
		delegator, err := k.accountKeeper.AddressCodec().BytesToString(delAddr)
		if err != nil {
			return false, err
		}
		withdraw, err := k.accountKeeper.AddressCodec().BytesToString(withdrawAddr)
		if err != nil {
			return false, err
		}
		delegatorWithdrawInfos = append(delegatorWithdrawInfos, types.DelegatorWithdrawInfo{
			DelegatorAddress: delegator,
			WithdrawAddress:  withdraw,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	poolOutstandingRewards := []types.PoolOutstandingRewardsRecord{}
	err = k.PoolOutstandingRewards.Walk(ctx, nil, func(poolID uint32, rewards types.OutstandingRewards) (stop bool, err error) {
		poolOutstandingRewards = append(poolOutstandingRewards, types.PoolOutstandingRewardsRecord{
			PoolID:             poolID,
			OutstandingRewards: rewards.Rewards,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	poolHistoricalRewards := []types.PoolHistoricalRewardsRecord{}
	err = k.PoolHistoricalRewards.Walk(ctx, nil, func(key collections.Pair[uint32, uint64], rewards types.HistoricalRewards) (stop bool, err error) {
		poolID := key.K1()
		period := key.K2()
		poolHistoricalRewards = append(poolHistoricalRewards, types.PoolHistoricalRewardsRecord{
			PoolID:  poolID,
			Period:  period,
			Rewards: rewards,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	poolCurrentRewards := []types.PoolCurrentRewardsRecord{}
	err = k.PoolCurrentRewards.Walk(ctx, nil, func(poolID uint32, rewards types.CurrentRewards) (stop bool, err error) {
		poolCurrentRewards = append(poolCurrentRewards, types.PoolCurrentRewardsRecord{
			PoolID:  poolID,
			Rewards: rewards,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	poolDelegatorStartingInfos := []types.PoolDelegatorStartingInfoRecord{}
	err = k.PoolDelegatorStartingInfos.Walk(ctx, nil, func(key collections.Pair[uint32, sdk.AccAddress], info types.DelegatorStartingInfo) (stop bool, err error) {
		poolID := key.K1()
		delAddr := key.K2()
		delegator, err := k.accountKeeper.AddressCodec().BytesToString(delAddr)
		if err != nil {
			return false, err
		}
		poolDelegatorStartingInfos = append(poolDelegatorStartingInfos, types.PoolDelegatorStartingInfoRecord{
			DelegatorAddress: delegator,
			PoolID:           poolID,
			StartingInfo:     info,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	operatorOutstandingRewards := []types.OperatorOutstandingRewardsRecord{}
	err = k.OperatorOutstandingRewards.Walk(ctx, nil, func(operatorID uint32, rewards types.MultiOutstandingRewards) (stop bool, err error) {
		operatorOutstandingRewards = append(operatorOutstandingRewards, types.OperatorOutstandingRewardsRecord{
			OperatorID:         operatorID,
			OutstandingRewards: rewards.Rewards,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	operatorAccumulatedCommissions := []types.OperatorAccumulatedCommissionRecord{}
	err = k.OperatorAccumulatedCommissions.Walk(ctx, nil, func(operatorID uint32, commission types.MultiAccumulatedCommission) (stop bool, err error) {
		operatorAccumulatedCommissions = append(operatorAccumulatedCommissions, types.OperatorAccumulatedCommissionRecord{
			OperatorId:  operatorID,
			Accumulated: commission,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	operatorHistoricalRewards := []types.OperatorHistoricalRewardsRecord{}
	err = k.OperatorHistoricalRewards.Walk(ctx, nil, func(key collections.Pair[uint32, uint64], rewards types.MultiHistoricalRewards) (stop bool, err error) {
		operatorID := key.K1()
		period := key.K2()
		operatorHistoricalRewards = append(operatorHistoricalRewards, types.OperatorHistoricalRewardsRecord{
			OperatorID: operatorID,
			Period:     period,
			Rewards:    rewards,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	operatorCurrentRewards := []types.OperatorCurrentRewardsRecord{}
	err = k.OperatorCurrentRewards.Walk(ctx, nil, func(operatorID uint32, rewards types.MultiCurrentRewards) (stop bool, err error) {
		operatorCurrentRewards = append(operatorCurrentRewards, types.OperatorCurrentRewardsRecord{
			OperatorID: operatorID,
			Rewards:    rewards,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	operatorDelegatorStartingInfos := []types.OperatorDelegatorStartingInfoRecord{}
	err = k.OperatorDelegatorStartingInfos.Walk(ctx, nil, func(key collections.Pair[uint32, sdk.AccAddress], info types.MultiDelegatorStartingInfo) (stop bool, err error) {
		operatorID := key.K1()
		delAddr := key.K2()
		delegator, err := k.accountKeeper.AddressCodec().BytesToString(delAddr)
		if err != nil {
			return false, err
		}
		operatorDelegatorStartingInfos = append(operatorDelegatorStartingInfos, types.OperatorDelegatorStartingInfoRecord{
			DelegatorAddress: delegator,
			OperatorID:       operatorID,
			StartingInfo:     info,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	serviceOutstandingRewards := []types.ServiceOutstandingRewardsRecord{}
	err = k.ServiceOutstandingRewards.Walk(ctx, nil, func(serviceID uint32, rewards types.MultiOutstandingRewards) (stop bool, err error) {
		serviceOutstandingRewards = append(serviceOutstandingRewards, types.ServiceOutstandingRewardsRecord{
			ServiceID:          serviceID,
			OutstandingRewards: rewards.Rewards,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	serviceHistoricalRewards := []types.ServiceHistoricalRewardsRecord{}
	err = k.ServiceHistoricalRewards.Walk(ctx, nil, func(key collections.Pair[uint32, uint64], rewards types.MultiHistoricalRewards) (stop bool, err error) {
		serviceID := key.K1()
		period := key.K2()
		serviceHistoricalRewards = append(serviceHistoricalRewards, types.ServiceHistoricalRewardsRecord{
			ServiceID: serviceID,
			Period:    period,
			Rewards:   rewards,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	serviceCurrentRewards := []types.ServiceCurrentRewardsRecord{}
	err = k.ServiceCurrentRewards.Walk(ctx, nil, func(serviceID uint32, rewards types.MultiCurrentRewards) (stop bool, err error) {
		serviceCurrentRewards = append(serviceCurrentRewards, types.ServiceCurrentRewardsRecord{
			ServiceID: serviceID,
			Rewards:   rewards,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	serviceDelegatorStartingInfos := []types.ServiceDelegatorStartingInfoRecord{}
	err = k.ServiceDelegatorStartingInfos.Walk(ctx, nil, func(key collections.Pair[uint32, sdk.AccAddress], info types.MultiDelegatorStartingInfo) (stop bool, err error) {
		serviceID := key.K1()
		delAddr := key.K2()
		delegator, err := k.accountKeeper.AddressCodec().BytesToString(delAddr)
		if err != nil {
			return false, err
		}
		serviceDelegatorStartingInfos = append(serviceDelegatorStartingInfos, types.ServiceDelegatorStartingInfoRecord{
			DelegatorAddress: delegator,
			ServiceID:        serviceID,
			StartingInfo:     info,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	return types.NewGenesisState(
		params,
		nextRewardsPlanID,
		rewardsPlans,
		lastRewardsAllocationTime,
		delegatorWithdrawInfos,
		poolOutstandingRewards,
		poolHistoricalRewards,
		poolCurrentRewards,
		poolDelegatorStartingInfos,
		operatorOutstandingRewards,
		operatorAccumulatedCommissions,
		operatorHistoricalRewards,
		operatorCurrentRewards,
		operatorDelegatorStartingInfos,
		serviceOutstandingRewards,
		serviceHistoricalRewards,
		serviceCurrentRewards,
		serviceDelegatorStartingInfos,
	)
}

// --------------------------------------------------------------------------------------------------------------------

// InitGenesis initializes the state from a GenesisState
func (k *Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) {
	// Store params
	if err := k.Params.Set(ctx, state.Params); err != nil {
		panic(err)
	}

	// Set the next distribution plan ID
	if err := k.NextRewardsPlanID.Set(ctx, state.NextRewardsPlanID); err != nil {
		panic(err)
	}

	// Store the rewards plans
	for _, plan := range state.RewardsPlans {
		if err := k.RewardsPlans.Set(ctx, plan.ID, plan); err != nil {
			panic(err)
		}
	}

	if state.LastRewardsAllocationTime != nil {
		if err := k.SetLastRewardsAllocationTime(ctx, *state.LastRewardsAllocationTime); err != nil {
			panic(err)
		}
	}

	for _, info := range state.DelegatorWithdrawInfos {
		delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(info.DelegatorAddress)
		if err != nil {
			panic(err)
		}
		withdrawAddr, err := k.accountKeeper.AddressCodec().StringToBytes(info.WithdrawAddress)
		if err != nil {
			panic(err)
		}
		err = k.DelegatorWithdrawAddrs.Set(ctx, delAddr, withdrawAddr)
		if err != nil {
			panic(err)
		}
	}

	for _, record := range state.PoolOutstandingRewards {
		if err := k.PoolOutstandingRewards.Set(
			ctx, record.PoolID, types.OutstandingRewards{Rewards: record.OutstandingRewards}); err != nil {
			panic(err)
		}
	}

	for _, record := range state.PoolHistoricalRewards {
		if err := k.PoolHistoricalRewards.Set(
			ctx, collections.Join(record.PoolID, record.Period), record.Rewards); err != nil {
			panic(err)
		}
	}

	for _, record := range state.PoolCurrentRewards {
		if err := k.PoolCurrentRewards.Set(ctx, record.PoolID, record.Rewards); err != nil {
			panic(err)
		}
	}

	for _, record := range state.PoolDelegatorStartingInfos {
		delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(record.DelegatorAddress)
		if err != nil {
			panic(err)
		}
		if err := k.PoolDelegatorStartingInfos.Set(
			ctx, collections.Join(record.PoolID, sdk.AccAddress(delAddr)), record.StartingInfo); err != nil {
			panic(err)
		}
	}

	for _, record := range state.OperatorOutstandingRewards {
		if err := k.OperatorOutstandingRewards.Set(
			ctx, record.OperatorID, types.MultiOutstandingRewards{Rewards: record.OutstandingRewards}); err != nil {
			panic(err)
		}
	}

	for _, record := range state.OperatorHistoricalRewards {
		if err := k.OperatorHistoricalRewards.Set(
			ctx, collections.Join(record.OperatorID, record.Period), record.Rewards); err != nil {
			panic(err)
		}
	}

	for _, record := range state.OperatorCurrentRewards {
		if err := k.OperatorCurrentRewards.Set(ctx, record.OperatorID, record.Rewards); err != nil {
			panic(err)
		}
	}

	for _, record := range state.OperatorDelegatorStartingInfos {
		delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(record.DelegatorAddress)
		if err != nil {
			panic(err)
		}
		if err := k.OperatorDelegatorStartingInfos.Set(
			ctx, collections.Join(record.OperatorID, sdk.AccAddress(delAddr)), record.StartingInfo); err != nil {
			panic(err)
		}
	}

	for _, record := range state.ServiceOutstandingRewards {
		if err := k.ServiceOutstandingRewards.Set(
			ctx, record.ServiceID, types.MultiOutstandingRewards{Rewards: record.OutstandingRewards}); err != nil {
			panic(err)
		}
	}

	for _, record := range state.ServiceHistoricalRewards {
		if err := k.ServiceHistoricalRewards.Set(
			ctx, collections.Join(record.ServiceID, record.Period), record.Rewards); err != nil {
			panic(err)
		}
	}

	for _, record := range state.ServiceCurrentRewards {
		if err := k.ServiceCurrentRewards.Set(ctx, record.ServiceID, record.Rewards); err != nil {
			panic(err)
		}
	}

	for _, record := range state.ServiceDelegatorStartingInfos {
		delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(record.DelegatorAddress)
		if err != nil {
			panic(err)
		}
		if err := k.ServiceDelegatorStartingInfos.Set(
			ctx, collections.Join(record.ServiceID, sdk.AccAddress(delAddr)), record.StartingInfo); err != nil {
			panic(err)
		}
	}

	// TODO: check module holdings
}
