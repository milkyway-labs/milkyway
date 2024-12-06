package types

import (
	"fmt"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/milkyway-labs/milkyway/v3/utils"
)

func NewGenesisState(
	params Params,
	nextRewardsPlanID uint64,
	rewardsPlans []RewardsPlan,
	lastRewardsAllocationTime *time.Time,
	delegatorWithdrawInfos []DelegatorWithdrawInfo,
	poolsRecords,
	operatorsRecords,
	servicesRecords DelegationTypeRecords,
	operatorAccumulatedCommissionRecords []OperatorAccumulatedCommissionRecord,
	poolServiceTotalDelShares []PoolServiceTotalDelegatorShares,
) *GenesisState {
	return &GenesisState{
		Params:                          params,
		NextRewardsPlanID:               nextRewardsPlanID,
		RewardsPlans:                    rewardsPlans,
		LastRewardsAllocationTime:       lastRewardsAllocationTime,
		DelegatorWithdrawInfos:          delegatorWithdrawInfos,
		PoolsRecords:                    poolsRecords,
		OperatorsRecords:                operatorsRecords,
		ServicesRecords:                 servicesRecords,
		OperatorAccumulatedCommissions:  operatorAccumulatedCommissionRecords,
		PoolServiceTotalDelegatorShares: poolServiceTotalDelShares,
	}
}

// DefaultGenesis returns the default genesis state.
func DefaultGenesis() *GenesisState {
	return NewGenesisState(
		DefaultParams(), 1, []RewardsPlan{}, nil, []DelegatorWithdrawInfo{},
		NewDelegationTypeRecords(
			[]OutstandingRewardsRecord{}, []HistoricalRewardsRecord{}, []CurrentRewardsRecord{},
			[]DelegatorStartingInfoRecord{}),
		NewDelegationTypeRecords(
			[]OutstandingRewardsRecord{}, []HistoricalRewardsRecord{}, []CurrentRewardsRecord{},
			[]DelegatorStartingInfoRecord{}),
		NewDelegationTypeRecords(
			[]OutstandingRewardsRecord{}, []HistoricalRewardsRecord{}, []CurrentRewardsRecord{},
			[]DelegatorStartingInfoRecord{}),
		[]OperatorAccumulatedCommissionRecord{},
		[]PoolServiceTotalDelegatorShares{},
	)
}

// Validate checks that the genesis state is valid.
func (genState *GenesisState) Validate(unpacker codectypes.AnyUnpacker) error {
	// Validate params
	err := genState.Params.Validate()
	if err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	if genState.NextRewardsPlanID == 0 {
		return fmt.Errorf("invalid next rewards plan ID: %d", genState.NextRewardsPlanID)
	}

	// Check for duplicate distribution plans
	if duplicate := findDuplicateRewardsPlans(genState.RewardsPlans); duplicate != nil {
		return fmt.Errorf("duplicated rewards plan: %d", duplicate.ID)
	}

	for i, plan := range genState.RewardsPlans {
		err = plan.Validate(unpacker)
		if err != nil {
			return fmt.Errorf("invalid rewards plan at index %d: %w", i, err)
		}
	}

	for _, shares := range genState.PoolServiceTotalDelegatorShares {
		err = shares.Validate()
		if err != nil {
			return fmt.Errorf("invalid pool service total delegator shares: %w", err)
		}
	}
	return nil
}

func NewDelegationTypeRecords(
	outstandingRewards []OutstandingRewardsRecord, historicalRewards []HistoricalRewardsRecord,
	currentRewards []CurrentRewardsRecord, delegatorStartingInfos []DelegatorStartingInfoRecord,
) DelegationTypeRecords {
	return DelegationTypeRecords{
		OutstandingRewards:     outstandingRewards,
		HistoricalRewards:      historicalRewards,
		CurrentRewards:         currentRewards,
		DelegatorStartingInfos: delegatorStartingInfos,
	}
}

// findDuplicateRewardsPlans returns the first duplicated rewards plan in the slice.
// If no duplicates are found, it returns nil instead.
func findDuplicateRewardsPlans(plans []RewardsPlan) *RewardsPlan {
	return utils.FindDuplicateFunc(plans, func(a, b RewardsPlan) bool {
		return a.ID == b.ID
	})
}
