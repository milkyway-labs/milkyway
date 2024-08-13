package types

import (
	"fmt"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/milkyway-labs/milkyway/utils"
)

func NewGenesisState(
	params Params,
	nextRewardsPlanID uint64,
	rewardsPlans []RewardsPlan,
	lastRewardsAllocationTime *time.Time,
	delegatorWithdrawInfos []DelegatorWithdrawInfo,
	poolOutstandingRewardsRecords []OutstandingRewardsRecord,
	poolHistoricalRewardsRecords []HistoricalRewardsRecord,
	poolCurrentRewardsRecords []CurrentRewardsRecord,
	poolDelegatorStartingInfoRecords []DelegatorStartingInfoRecord,
	operatorOutstandingRewardsRecords []OutstandingRewardsRecord,
	operatorAccumulatedCommissionRecords []OperatorAccumulatedCommissionRecord,
	operatorHistoricalRewardsRecords []HistoricalRewardsRecord,
	operatorCurrentRewardsRecords []CurrentRewardsRecord,
	operatorDelegatorStartingInfoRecords []DelegatorStartingInfoRecord,
	serviceOutstandingRewardsRecords []OutstandingRewardsRecord,
	serviceHistoricalRewardsRecords []HistoricalRewardsRecord,
	serviceCurrentRewardsRecords []CurrentRewardsRecord,
	serviceDelegatorStartingInfoRecords []DelegatorStartingInfoRecord,
) *GenesisState {
	return &GenesisState{
		Params:                         params,
		NextRewardsPlanID:              nextRewardsPlanID,
		RewardsPlans:                   rewardsPlans,
		LastRewardsAllocationTime:      lastRewardsAllocationTime,
		DelegatorWithdrawInfos:         delegatorWithdrawInfos,
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
		DefaultParams(), 1, []RewardsPlan{}, nil, []DelegatorWithdrawInfo{}, []OutstandingRewardsRecord{},
		[]HistoricalRewardsRecord{}, []CurrentRewardsRecord{}, []DelegatorStartingInfoRecord{},
		[]OutstandingRewardsRecord{}, []OperatorAccumulatedCommissionRecord{}, []HistoricalRewardsRecord{},
		[]CurrentRewardsRecord{}, []DelegatorStartingInfoRecord{}, []OutstandingRewardsRecord{},
		[]HistoricalRewardsRecord{}, []CurrentRewardsRecord{}, []DelegatorStartingInfoRecord{})
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
	return nil
}

// findDuplicateRewardsPlans returns the first duplicated rewards plan in the slice.
// If no duplicates are found, it returns nil instead.
func findDuplicateRewardsPlans(plans []RewardsPlan) *RewardsPlan {
	return utils.FindDuplicateFunc(plans, func(a, b RewardsPlan) bool {
		return a.ID == b.ID
	})
}
