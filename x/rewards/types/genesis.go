package types

import (
	"fmt"
	"time"

	"github.com/milkyway-labs/milkyway/utils"
)

func NewGenesisState(
	params Params,
	nextRewardsPlanID uint64,
	rewardsPlans []RewardsPlan,
	lastRewardsAllocationTime *time.Time,
	poolOutstandingRewardsRecords []PoolOutstandingRewardsRecord,
	poolHistoricalRewardsRecords []PoolHistoricalRewardsRecord,
	poolCurrentRewardsRecords []PoolCurrentRewardsRecord,
	poolDelegatorStartingInfoRecords []PoolDelegatorStartingInfoRecord,
	operatorOutstandingRewardsRecords []OperatorOutstandingRewardsRecord,
	operatorAccumulatedCommissionRecords []OperatorAccumulatedCommissionRecord,
	operatorHistoricalRewardsRecords []OperatorHistoricalRewardsRecord,
	operatorCurrentRewardsRecords []OperatorCurrentRewardsRecord,
	operatorDelegatorStartingInfoRecords []OperatorDelegatorStartingInfoRecord,
	serviceOutstandingRewardsRecords []ServiceOutstandingRewardsRecord,
	serviceHistoricalRewardsRecords []ServiceHistoricalRewardsRecord,
	serviceCurrentRewardsRecords []ServiceCurrentRewardsRecord,
	serviceDelegatorStartingInfoRecords []ServiceDelegatorStartingInfoRecord,
) *GenesisState {
	return &GenesisState{
		Params:                         params,
		NextRewardsPlanID:              nextRewardsPlanID,
		RewardsPlans:                   rewardsPlans,
		LastRewardsAllocationTime:      lastRewardsAllocationTime,
		PoolOutstandingRewards:         poolOutstandingRewardsRecords,
		PoolHistoricalRewards:          poolHistoricalRewardsRecords,
		PoolCurrentRewards:             poolCurrentRewardsRecords,
		PoolDelegatorStartingInfos:     poolDelegatorStartingInfoRecords,
		OperatorOutstandingRewards:     operatorOutstandingRewardsRecords,
		OperatorAccumulatedCommissions: operatorAccumulatedCommissionRecords,
		OperatorHistoricalRewards:      operatorHistoricalRewardsRecords,
		OperatorCurrentRewards:         operatorCurrentRewardsRecords,
		OperatorDelegatorStartingInfos: operatorDelegatorStartingInfoRecords,
		ServiceOutstandingRewards:      serviceOutstandingRewardsRecords,
		ServiceHistoricalRewards:       serviceHistoricalRewardsRecords,
		ServiceCurrentRewards:          serviceCurrentRewardsRecords,
		ServiceDelegatorStartingInfos:  serviceDelegatorStartingInfoRecords,
	}
}

// DefaultGenesis returns the default genesis state.
func DefaultGenesis() *GenesisState {
	return NewGenesisState(
		DefaultParams(), 1, []RewardsPlan{}, nil, []PoolOutstandingRewardsRecord{},
		[]PoolHistoricalRewardsRecord{}, []PoolCurrentRewardsRecord{}, []PoolDelegatorStartingInfoRecord{},
		[]OperatorOutstandingRewardsRecord{}, []OperatorAccumulatedCommissionRecord{}, []OperatorHistoricalRewardsRecord{},
		[]OperatorCurrentRewardsRecord{}, []OperatorDelegatorStartingInfoRecord{}, []ServiceOutstandingRewardsRecord{},
		[]ServiceHistoricalRewardsRecord{}, []ServiceCurrentRewardsRecord{}, []ServiceDelegatorStartingInfoRecord{})
}

// Validate checks that the genesis state is valid.
func (genState *GenesisState) Validate() error {
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

	// TODO: add more validations

	return nil
}

// findDuplicateRewardsPlans returns the first duplicated rewards plan in the slice.
// If no duplicates are found, it returns nil instead.
func findDuplicateRewardsPlans(plans []RewardsPlan) *RewardsPlan {
	return utils.FindDuplicateFunc(plans, func(a, b RewardsPlan) bool {
		return a.ID == b.ID
	})
}
