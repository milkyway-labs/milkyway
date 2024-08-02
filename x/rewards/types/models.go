package types

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/gogoproto/proto"
)

// GetRewardsPoolAddress generates a rewards pool address from plan id.
func GetRewardsPoolAddress(planID uint64) sdk.AccAddress {
	return address.Module(ModuleName, []byte(fmt.Sprintf("rewards-plan-%d", planID)))
}

func NewRewardsPlan(
	id uint64, description string, serviceID uint32, amtPerDay sdk.Coins, startTime, endTime time.Time,
	rewardsPool string, poolsDistribution PoolsDistribution, operatorsDistribution OperatorsDistribution,
	usersDistribution UsersDistribution) RewardsPlan {
	return RewardsPlan{
		ID:                    id,
		Description:           description,
		ServiceID:             serviceID,
		AmountPerDay:          amtPerDay,
		StartTime:             startTime,
		EndTime:               endTime,
		RewardsPool:           rewardsPool,
		PoolsDistribution:     poolsDistribution,
		OperatorsDistribution: operatorsDistribution,
		UsersDistribution:     usersDistribution,
	}
}

func (plan RewardsPlan) IsActiveAt(t time.Time) bool {
	return !plan.StartTime.After(t) && plan.EndTime.After(t)
}

func (plan RewardsPlan) MustGetRewardsPoolAddress() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(plan.RewardsPool)
	if err != nil {
		panic(err)
	}
	return addr
}

func (plan RewardsPlan) TotalWeights() uint32 {
	return plan.PoolsDistribution.Weight + plan.OperatorsDistribution.Weight + plan.UsersDistribution.Weight
}

func (plan RewardsPlan) Validate() error {
	if plan.ID == 0 {
		return fmt.Errorf("invalid plan ID")
	}
	// TODO: validate description
	if plan.ServiceID == 0 {
		return fmt.Errorf("invalid service ID")
	}
	if err := plan.AmountPerDay.Validate(); err != nil {
		return fmt.Errorf("invalid amount per day: %w", err)
	}
	if !plan.EndTime.After(plan.StartTime) {
		return fmt.Errorf(
			"end time must be after start time: %s <= %s",
			plan.EndTime.Format(time.RFC3339), plan.StartTime.Format(time.RFC3339))
	}
	if _, err := sdk.AccAddressFromBech32(plan.RewardsPool); err != nil {
		return fmt.Errorf("invalid rewards pool: %w", err)
	}

	// TODO: need to validate distribution types?
	return nil
}

func NewBasicPoolsDistribution(weight uint32) PoolsDistribution {
	a, err := types.NewAnyWithValue(&PoolsDistributionTypeBasic{})
	if err != nil {
		panic(err)
	}
	return PoolsDistribution{
		Weight: weight,
		Type:   a,
	}
}

func NewWeightedPoolsDistribution(weight uint32, poolWeights []PoolDistributionWeight) PoolsDistribution {
	a, err := types.NewAnyWithValue(&PoolsDistributionTypeWeighted{Weights: poolWeights})
	if err != nil {
		panic(err)
	}
	return PoolsDistribution{
		Weight: weight,
		Type:   a,
	}
}

func NewEgalitarianPoolsDistribution(weight uint32) PoolsDistribution {
	a, err := types.NewAnyWithValue(&PoolsDistributionTypeEgalitarian{})
	if err != nil {
		panic(err)
	}
	return PoolsDistribution{
		Weight: weight,
		Type:   a,
	}
}

type PoolsDistributionType interface {
	proto.Message
	isPoolsDistributionType()
}

func (t PoolsDistributionTypeBasic) isPoolsDistributionType() {}

func (t PoolsDistributionTypeWeighted) isPoolsDistributionType() {}

func (t PoolsDistributionTypeEgalitarian) isPoolsDistributionType() {}

type OperatorsDistributionType interface {
	proto.Message
	isOperatorsDistributionType()
}

func NewBasicOperatorsDistribution(weight uint32) OperatorsDistribution {
	a, err := types.NewAnyWithValue(&OperatorsDistributionTypeBasic{})
	if err != nil {
		panic(err)
	}
	return OperatorsDistribution{
		Weight: weight,
		Type:   a,
	}
}

func NewWeightedOperatorsDistribution(
	weight uint32, operatorWeights []OperatorDistributionWeight) OperatorsDistribution {
	a, err := types.NewAnyWithValue(&OperatorsDistributionTypeWeighted{Weights: operatorWeights})
	if err != nil {
		panic(err)
	}
	return OperatorsDistribution{
		Weight: weight,
		Type:   a,
	}
}

func NewEgalitarianOperatorsDistribution(weight uint32) OperatorsDistribution {
	a, err := types.NewAnyWithValue(&OperatorsDistributionTypeEgalitarian{})
	if err != nil {
		panic(err)
	}
	return OperatorsDistribution{
		Weight: weight,
		Type:   a,
	}
}

func (t OperatorsDistributionTypeBasic) isOperatorsDistributionType() {}

func (t OperatorsDistributionTypeWeighted) isOperatorsDistributionType() {}

func (t OperatorsDistributionTypeEgalitarian) isOperatorsDistributionType() {}

type UsersDistributionType interface {
	proto.Message
	isUsersDistributionType()
}

func NewBasicUsersDistribution(weight uint32) UsersDistribution {
	a, err := types.NewAnyWithValue(&UsersDistributionTypeBasic{})
	if err != nil {
		panic(err)
	}
	return UsersDistribution{
		Weight: weight,
		Type:   a,
	}
}

func (t UsersDistributionTypeBasic) isUsersDistributionType() {}

// return the initial accumulated commission (zero)
func InitialAccumulatedCommission() AccumulatedCommission {
	return AccumulatedCommission{}
}

func NewHistoricalRewards(cumulativeRewardRatios DecPools, referenceCount uint32) HistoricalRewards {
	return HistoricalRewards{
		CumulativeRewardRatios: cumulativeRewardRatios,
		ReferenceCount:         referenceCount,
	}
}

func NewCurrentRewards(rewards DecPools, period uint64) CurrentRewards {
	return CurrentRewards{
		Rewards: rewards,
		Period:  period,
	}
}

// create a new DelegatorStartingInfo
func NewDelegatorStartingInfo(previousPeriod uint64, stakes sdk.DecCoins, height uint64) DelegatorStartingInfo {
	return DelegatorStartingInfo{
		PreviousPeriod: previousPeriod,
		Stakes:         stakes,
		Height:         height,
	}
}
