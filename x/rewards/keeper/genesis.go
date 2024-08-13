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

	poolOutstandingRewards := []types.OutstandingRewardsRecord{}
	err = k.PoolOutstandingRewards.Walk(ctx, nil, func(poolID uint32, rewards types.OutstandingRewards) (stop bool, err error) {
		poolOutstandingRewards = append(poolOutstandingRewards, types.OutstandingRewardsRecord{
			DelegationTargetID: poolID,
			OutstandingRewards: rewards.Rewards,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	poolHistoricalRewards := []types.HistoricalRewardsRecord{}
	err = k.PoolHistoricalRewards.Walk(ctx, nil, func(key collections.Pair[uint32, uint64], rewards types.HistoricalRewards) (stop bool, err error) {
		poolID := key.K1()
		period := key.K2()
		poolHistoricalRewards = append(poolHistoricalRewards, types.HistoricalRewardsRecord{
			DelegationTargetID: poolID,
			Period:             period,
			Rewards:            rewards,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	poolCurrentRewards := []types.CurrentRewardsRecord{}
	err = k.PoolCurrentRewards.Walk(ctx, nil, func(poolID uint32, rewards types.CurrentRewards) (stop bool, err error) {
		poolCurrentRewards = append(poolCurrentRewards, types.CurrentRewardsRecord{
			DelegationTargetID: poolID,
			Rewards:            rewards,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	poolDelegatorStartingInfos := []types.DelegatorStartingInfoRecord{}
	err = k.PoolDelegatorStartingInfos.Walk(ctx, nil, func(key collections.Pair[uint32, sdk.AccAddress], info types.DelegatorStartingInfo) (stop bool, err error) {
		poolID := key.K1()
		delAddr := key.K2()
		delegator, err := k.accountKeeper.AddressCodec().BytesToString(delAddr)
		if err != nil {
			return false, err
		}
		poolDelegatorStartingInfos = append(poolDelegatorStartingInfos, types.DelegatorStartingInfoRecord{
			DelegatorAddress:   delegator,
			DelegationTargetID: poolID,
			StartingInfo:       info,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	operatorOutstandingRewards := []types.OutstandingRewardsRecord{}
	err = k.OperatorOutstandingRewards.Walk(ctx, nil, func(operatorID uint32, rewards types.OutstandingRewards) (stop bool, err error) {
		operatorOutstandingRewards = append(operatorOutstandingRewards, types.OutstandingRewardsRecord{
			DelegationTargetID: operatorID,
			OutstandingRewards: rewards.Rewards,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	operatorAccumulatedCommissions := []types.OperatorAccumulatedCommissionRecord{}
	err = k.OperatorAccumulatedCommissions.Walk(ctx, nil, func(operatorID uint32, commission types.AccumulatedCommission) (stop bool, err error) {
		operatorAccumulatedCommissions = append(operatorAccumulatedCommissions, types.OperatorAccumulatedCommissionRecord{
			OperatorID:  operatorID,
			Accumulated: commission,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	operatorHistoricalRewards := []types.HistoricalRewardsRecord{}
	err = k.OperatorHistoricalRewards.Walk(ctx, nil, func(key collections.Pair[uint32, uint64], rewards types.HistoricalRewards) (stop bool, err error) {
		operatorID := key.K1()
		period := key.K2()
		operatorHistoricalRewards = append(operatorHistoricalRewards, types.HistoricalRewardsRecord{
			DelegationTargetID: operatorID,
			Period:             period,
			Rewards:            rewards,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	operatorCurrentRewards := []types.CurrentRewardsRecord{}
	err = k.OperatorCurrentRewards.Walk(ctx, nil, func(operatorID uint32, rewards types.CurrentRewards) (stop bool, err error) {
		operatorCurrentRewards = append(operatorCurrentRewards, types.CurrentRewardsRecord{
			DelegationTargetID: operatorID,
			Rewards:            rewards,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	operatorDelegatorStartingInfos := []types.DelegatorStartingInfoRecord{}
	err = k.OperatorDelegatorStartingInfos.Walk(ctx, nil, func(key collections.Pair[uint32, sdk.AccAddress], info types.DelegatorStartingInfo) (stop bool, err error) {
		operatorID := key.K1()
		delAddr := key.K2()
		delegator, err := k.accountKeeper.AddressCodec().BytesToString(delAddr)
		if err != nil {
			return false, err
		}
		operatorDelegatorStartingInfos = append(operatorDelegatorStartingInfos, types.DelegatorStartingInfoRecord{
			DelegatorAddress:   delegator,
			DelegationTargetID: operatorID,
			StartingInfo:       info,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	serviceOutstandingRewards := []types.OutstandingRewardsRecord{}
	err = k.ServiceOutstandingRewards.Walk(ctx, nil, func(serviceID uint32, rewards types.OutstandingRewards) (stop bool, err error) {
		serviceOutstandingRewards = append(serviceOutstandingRewards, types.OutstandingRewardsRecord{
			DelegationTargetID: serviceID,
			OutstandingRewards: rewards.Rewards,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	serviceHistoricalRewards := []types.HistoricalRewardsRecord{}
	err = k.ServiceHistoricalRewards.Walk(ctx, nil, func(key collections.Pair[uint32, uint64], rewards types.HistoricalRewards) (stop bool, err error) {
		serviceID := key.K1()
		period := key.K2()
		serviceHistoricalRewards = append(serviceHistoricalRewards, types.HistoricalRewardsRecord{
			DelegationTargetID: serviceID,
			Period:             period,
			Rewards:            rewards,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	serviceCurrentRewards := []types.CurrentRewardsRecord{}
	err = k.ServiceCurrentRewards.Walk(ctx, nil, func(serviceID uint32, rewards types.CurrentRewards) (stop bool, err error) {
		serviceCurrentRewards = append(serviceCurrentRewards, types.CurrentRewardsRecord{
			DelegationTargetID: serviceID,
			Rewards:            rewards,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	serviceDelegatorStartingInfos := []types.DelegatorStartingInfoRecord{}
	err = k.ServiceDelegatorStartingInfos.Walk(ctx, nil, func(key collections.Pair[uint32, sdk.AccAddress], info types.DelegatorStartingInfo) (stop bool, err error) {
		serviceID := key.K1()
		delAddr := key.K2()
		delegator, err := k.accountKeeper.AddressCodec().BytesToString(delAddr)
		if err != nil {
			return false, err
		}
		serviceDelegatorStartingInfos = append(serviceDelegatorStartingInfos, types.DelegatorStartingInfoRecord{
			DelegatorAddress:   delegator,
			DelegationTargetID: serviceID,
			StartingInfo:       info,
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
			ctx, record.DelegationTargetID, types.OutstandingRewards{Rewards: record.OutstandingRewards}); err != nil {
			panic(err)
		}
	}

	for _, record := range state.PoolHistoricalRewards {
		if err := k.PoolHistoricalRewards.Set(
			ctx, collections.Join(record.DelegationTargetID, record.Period), record.Rewards); err != nil {
			panic(err)
		}
	}

	for _, record := range state.PoolCurrentRewards {
		if err := k.PoolCurrentRewards.Set(ctx, record.DelegationTargetID, record.Rewards); err != nil {
			panic(err)
		}
	}

	for _, record := range state.PoolDelegatorStartingInfos {
		delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(record.DelegatorAddress)
		if err != nil {
			panic(err)
		}
		if err := k.PoolDelegatorStartingInfos.Set(
			ctx, collections.Join(record.DelegationTargetID, sdk.AccAddress(delAddr)), record.StartingInfo); err != nil {
			panic(err)
		}
	}

	for _, record := range state.OperatorOutstandingRewards {
		if err := k.OperatorOutstandingRewards.Set(
			ctx, record.DelegationTargetID, types.OutstandingRewards{Rewards: record.OutstandingRewards}); err != nil {
			panic(err)
		}
	}

	for _, record := range state.OperatorAccumulatedCommissions {
		err := k.OperatorAccumulatedCommissions.Set(ctx, record.OperatorID, record.Accumulated)
		if err != nil {
			panic(err)
		}
	}

	for _, record := range state.OperatorHistoricalRewards {
		if err := k.OperatorHistoricalRewards.Set(
			ctx, collections.Join(record.DelegationTargetID, record.Period), record.Rewards); err != nil {
			panic(err)
		}
	}

	for _, record := range state.OperatorCurrentRewards {
		if err := k.OperatorCurrentRewards.Set(ctx, record.DelegationTargetID, record.Rewards); err != nil {
			panic(err)
		}
	}

	for _, record := range state.OperatorDelegatorStartingInfos {
		delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(record.DelegatorAddress)
		if err != nil {
			panic(err)
		}
		if err := k.OperatorDelegatorStartingInfos.Set(
			ctx, collections.Join(record.DelegationTargetID, sdk.AccAddress(delAddr)), record.StartingInfo); err != nil {
			panic(err)
		}
	}

	for _, record := range state.ServiceOutstandingRewards {
		if err := k.ServiceOutstandingRewards.Set(
			ctx, record.DelegationTargetID, types.OutstandingRewards{Rewards: record.OutstandingRewards}); err != nil {
			panic(err)
		}
	}

	for _, record := range state.ServiceHistoricalRewards {
		if err := k.ServiceHistoricalRewards.Set(
			ctx, collections.Join(record.DelegationTargetID, record.Period), record.Rewards); err != nil {
			panic(err)
		}
	}

	for _, record := range state.ServiceCurrentRewards {
		if err := k.ServiceCurrentRewards.Set(ctx, record.DelegationTargetID, record.Rewards); err != nil {
			panic(err)
		}
	}

	for _, record := range state.ServiceDelegatorStartingInfos {
		delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(record.DelegatorAddress)
		if err != nil {
			panic(err)
		}
		if err := k.ServiceDelegatorStartingInfos.Set(
			ctx, collections.Join(record.DelegationTargetID, sdk.AccAddress(delAddr)), record.StartingInfo); err != nil {
			panic(err)
		}
	}

	// TODO: check module holdings
}
