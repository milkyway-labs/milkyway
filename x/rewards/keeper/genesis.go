package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v12/x/rewards/types"
)

// ExportGenesis returns the GenesisState associated with the given context
func (k *Keeper) ExportGenesis(ctx sdk.Context) (*types.GenesisState, error) {
	// Get the params
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	// Get the next rewards plan ID
	nextRewardsPlanID, err := k.NextRewardsPlanID.Get(ctx)
	if err != nil {
		return nil, err
	}

	// Get the rewards plans
	rewardsPlans, err := k.GetRewardsPlans(ctx)
	if err != nil {
		return nil, err
	}

	// Get the last rewards allocation time
	lastRewardsAllocationTime, err := k.GetLastRewardsAllocationTime(ctx)
	if err != nil {
		return nil, err
	}

	var delegatorWithdrawInfos []types.DelegatorWithdrawInfo
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
		return nil, err
	}

	// Get the pool records info
	poolsRecords, err := k.exportDelegationsRecords(
		ctx,
		k.PoolOutstandingRewards,
		k.PoolHistoricalRewards,
		k.PoolCurrentRewards,
		k.PoolDelegatorStartingInfos,
	)
	if err != nil {
		return nil, err
	}

	// Get the operator records info
	operatorRecords, err := k.exportDelegationsRecords(
		ctx,
		k.OperatorOutstandingRewards,
		k.OperatorHistoricalRewards,
		k.OperatorCurrentRewards,
		k.OperatorDelegatorStartingInfos,
	)
	if err != nil {
		return nil, err
	}

	// Get the service records info
	serviceRecords, err := k.exportDelegationsRecords(
		ctx,
		k.ServiceOutstandingRewards,
		k.ServiceHistoricalRewards,
		k.ServiceCurrentRewards,
		k.ServiceDelegatorStartingInfos,
	)
	if err != nil {
		return nil, err
	}

	var operatorAccumulatedCommissions []types.OperatorAccumulatedCommissionRecord
	err = k.OperatorAccumulatedCommissions.Walk(ctx, nil, func(operatorID uint32, commission types.AccumulatedCommission) (stop bool, err error) {
		operatorAccumulatedCommissions = append(operatorAccumulatedCommissions, types.OperatorAccumulatedCommissionRecord{
			OperatorID:  operatorID,
			Accumulated: commission,
		})
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	var totalShares []types.PoolServiceTotalDelegatorShares
	err = k.PoolServiceTotalDelegatorShares.Walk(ctx, nil, func(_ collections.Pair[uint32, uint32], shares types.PoolServiceTotalDelegatorShares) (stop bool, err error) {
		totalShares = append(totalShares, shares)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return types.NewGenesisState(
		params,
		nextRewardsPlanID,
		rewardsPlans,
		lastRewardsAllocationTime,
		delegatorWithdrawInfos,
		poolsRecords,
		operatorRecords,
		serviceRecords,
		operatorAccumulatedCommissions,
		totalShares,
	), nil
}

// exportDelegationsRecords exports the delegation records for a specific type of delegation
func (k *Keeper) exportDelegationsRecords(
	ctx sdk.Context,
	outstandingRewardsCollection collections.Map[uint32, types.OutstandingRewards],
	historicalRewardsCollection collections.Map[collections.Pair[uint32, uint64], types.HistoricalRewards],
	currentRewardsCollection collections.Map[uint32, types.CurrentRewards],
	startingInfoCollection collections.Map[collections.Pair[uint32, sdk.AccAddress], types.DelegatorStartingInfo],
) (types.DelegationTypeRecords, error) {
	// Get the outstanding rewards
	var outstandingRewards []types.OutstandingRewardsRecord
	err := outstandingRewardsCollection.Walk(ctx, nil, func(targetID uint32, rewards types.OutstandingRewards) (stop bool, err error) {
		outstandingRewards = append(outstandingRewards, types.OutstandingRewardsRecord{
			DelegationTargetID: targetID,
			OutstandingRewards: rewards.Rewards,
		})
		return false, nil
	})
	if err != nil {
		return types.DelegationTypeRecords{}, err
	}

	// Get the historical rewards
	var historicalRewards []types.HistoricalRewardsRecord
	err = historicalRewardsCollection.Walk(ctx, nil, func(key collections.Pair[uint32, uint64], rewards types.HistoricalRewards) (stop bool, err error) {
		targetID := key.K1()
		period := key.K2()
		historicalRewards = append(historicalRewards, types.HistoricalRewardsRecord{
			DelegationTargetID: targetID,
			Period:             period,
			Rewards:            rewards,
		})
		return false, nil
	})
	if err != nil {
		return types.DelegationTypeRecords{}, err
	}

	// Get the current rewards
	var currentRewards []types.CurrentRewardsRecord
	err = currentRewardsCollection.Walk(ctx, nil, func(targetID uint32, rewards types.CurrentRewards) (stop bool, err error) {
		currentRewards = append(currentRewards, types.CurrentRewardsRecord{
			DelegationTargetID: targetID,
			Rewards:            rewards,
		})
		return false, nil
	})
	if err != nil {
		return types.DelegationTypeRecords{}, err
	}

	// Get the delegator starting infos
	var delegatorStartingInfos []types.DelegatorStartingInfoRecord
	err = startingInfoCollection.Walk(ctx, nil, func(key collections.Pair[uint32, sdk.AccAddress], info types.DelegatorStartingInfo) (stop bool, err error) {
		targetID := key.K1()
		delAddr := key.K2()

		delegator, err := k.accountKeeper.AddressCodec().BytesToString(delAddr)
		if err != nil {
			return false, err
		}

		delegatorStartingInfos = append(delegatorStartingInfos, types.DelegatorStartingInfoRecord{
			DelegatorAddress:   delegator,
			DelegationTargetID: targetID,
			StartingInfo:       info,
		})
		return false, nil
	})
	if err != nil {
		return types.DelegationTypeRecords{}, err
	}

	// Return the delegation records
	return types.NewDelegationTypeRecords(
		outstandingRewards,
		historicalRewards,
		currentRewards,
		delegatorStartingInfos,
	), nil
}

// --------------------------------------------------------------------------------------------------------------------

// InitGenesis initializes the state from a GenesisState
func (k *Keeper) InitGenesis(ctx sdk.Context, state *types.GenesisState) error {
	var totalOutstandingRewards sdk.DecCoins

	// Store params
	if err := k.Params.Set(ctx, state.Params); err != nil {
		return err
	}

	// Set the next distribution plan ID
	if err := k.NextRewardsPlanID.Set(ctx, state.NextRewardsPlanID); err != nil {
		return err
	}

	// Store the rewards plans
	for _, plan := range state.RewardsPlans {
		if err := k.RewardsPlans.Set(ctx, plan.ID, plan); err != nil {
			return err
		}
	}

	// Store the last rewards allocation time
	if state.LastRewardsAllocationTime != nil {
		if err := k.SetLastRewardsAllocationTime(ctx, *state.LastRewardsAllocationTime); err != nil {
			return err
		}
	}

	// Store the delegator withdraw addresses
	for _, info := range state.DelegatorWithdrawInfos {
		delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(info.DelegatorAddress)
		if err != nil {
			return err
		}
		withdrawAddr, err := k.accountKeeper.AddressCodec().StringToBytes(info.WithdrawAddress)
		if err != nil {
			return err
		}
		err = k.DelegatorWithdrawAddrs.Set(ctx, delAddr, withdrawAddr)
		if err != nil {
			return err
		}
	}

	// Initialize the pool records
	poolsTotalOutstandingRewards, err := k.initializeDelegationRecords(
		ctx,
		state.PoolsRecords,
		k.PoolOutstandingRewards,
		k.PoolHistoricalRewards,
		k.PoolCurrentRewards,
		k.PoolDelegatorStartingInfos,
	)
	if err != nil {
		return err
	}
	totalOutstandingRewards = totalOutstandingRewards.Add(poolsTotalOutstandingRewards...)

	// Initialize the operator records
	operatorsTotalOutstandingRewards, err := k.initializeDelegationRecords(
		ctx,
		state.OperatorsRecords,
		k.OperatorOutstandingRewards,
		k.OperatorHistoricalRewards,
		k.OperatorCurrentRewards,
		k.OperatorDelegatorStartingInfos,
	)
	if err != nil {
		return err
	}
	totalOutstandingRewards = totalOutstandingRewards.Add(operatorsTotalOutstandingRewards...)

	// Store the operator accumulated commissions
	for _, record := range state.OperatorAccumulatedCommissions {
		err := k.OperatorAccumulatedCommissions.Set(ctx, record.OperatorID, record.Accumulated)
		if err != nil {
			return err
		}
	}

	// Initialize the service records
	servicesTotalOutstandingRewards, err := k.initializeDelegationRecords(
		ctx,
		state.ServicesRecords,
		k.ServiceOutstandingRewards,
		k.ServiceHistoricalRewards,
		k.ServiceCurrentRewards,
		k.ServiceDelegatorStartingInfos,
	)
	if err != nil {
		return err
	}
	totalOutstandingRewards = totalOutstandingRewards.Add(servicesTotalOutstandingRewards...)

	// Check module holdings. Sum of all outstanding rewards must not be
	// greater than the holdings of the rewards pool module account.
	rewardsPoolAcc := k.accountKeeper.GetModuleAccount(ctx, types.RewardsPoolName)
	if rewardsPoolAcc == nil {
		return fmt.Errorf("rewards pool module account has not been set")
	}

	totalOutstandingRewardsTruncated, _ := totalOutstandingRewards.TruncateDecimal()

	// Get the rewards pool balances based on the denoms that have outstanding rewards
	// This is to avoid the call to GetAllBalances which can be exploited by a malicious user
	// since it iterates unboundedly over the full address balance
	rewardsPoolBalances := sdk.NewCoins()
	for _, outstandingReward := range totalOutstandingRewards {
		rewardsPoolBalance := k.bankKeeper.GetBalance(ctx, rewardsPoolAcc.GetAddress(), outstandingReward.Denom)
		rewardsPoolBalances = rewardsPoolBalances.Add(rewardsPoolBalance)
	}

	// Save the rewards pool module account if balances are zero.
	// This code is taken from Cosmos SDK.
	if rewardsPoolBalances.IsZero() {
		k.accountKeeper.SetModuleAccount(ctx, rewardsPoolAcc)
	}

	if totalOutstandingRewardsTruncated.IsAnyGT(rewardsPoolBalances) {
		return fmt.Errorf("rewards pool module balance does not match the module holdings: %s < %s",
			rewardsPoolBalances,
			totalOutstandingRewardsTruncated)
	}

	for _, delShares := range state.PoolServiceTotalDelegatorShares {
		err = k.SetPoolServiceTotalDelegatorShares(ctx, delShares.PoolID, delShares.ServiceID, delShares.Shares)
		if err != nil {
			return err
		}
	}

	return nil
}

// initializeDelegationRecords initializes the delegation records for a specific type of delegation
func (k *Keeper) initializeDelegationRecords(
	ctx sdk.Context,
	records types.DelegationTypeRecords,
	outstandingRewardsCollection collections.Map[uint32, types.OutstandingRewards],
	historicalRewardsCollection collections.Map[collections.Pair[uint32, uint64], types.HistoricalRewards],
	currentRewardsCollection collections.Map[uint32, types.CurrentRewards],
	startingInfoCollection collections.Map[collections.Pair[uint32, sdk.AccAddress], types.DelegatorStartingInfo],
) (totalOutstandingRewards sdk.DecCoins, err error) {

	// Store the outstanding rewards
	for _, outstanding := range records.OutstandingRewards {
		err = outstandingRewardsCollection.Set(ctx, outstanding.DelegationTargetID, types.OutstandingRewards{Rewards: outstanding.OutstandingRewards})
		if err != nil {
			return nil, err
		}
		totalOutstandingRewards = totalOutstandingRewards.Add(outstanding.OutstandingRewards.Sum()...)
	}

	// Store the historical rewards
	for _, historical := range records.HistoricalRewards {
		err = historicalRewardsCollection.Set(ctx, collections.Join(historical.DelegationTargetID, historical.Period), historical.Rewards)
		if err != nil {
			return nil, err
		}
	}

	// Store the current rewards
	for _, current := range records.CurrentRewards {
		err = currentRewardsCollection.Set(ctx, current.DelegationTargetID, current.Rewards)
		if err != nil {
			return nil, err
		}
	}

	// Store the delegators starting info
	for _, startingInfo := range records.DelegatorStartingInfos {
		delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(startingInfo.DelegatorAddress)
		if err != nil {
			return nil, err
		}

		err = startingInfoCollection.Set(ctx, collections.Join(startingInfo.DelegationTargetID, sdk.AccAddress(delAddr)), startingInfo.StartingInfo)
		if err != nil {
			return nil, err
		}
	}

	return totalOutstandingRewards, nil
}
