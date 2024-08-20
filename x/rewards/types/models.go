package types

import (
	"fmt"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/gogoproto/proto"

	"github.com/milkyway-labs/milkyway/utils"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
)

const MaxRewardsPlanDescriptionLength = 1000

var (
	_ DistributionType      = (*DistributionTypeBasic)(nil)
	_ DistributionType      = (*DistributionTypeWeighted)(nil)
	_ DistributionType      = (*DistributionTypeEgalitarian)(nil)
	_ DistributionType      = (*DistributionTypeBasic)(nil)
	_ DistributionType      = (*DistributionTypeWeighted)(nil)
	_ DistributionType      = (*DistributionTypeEgalitarian)(nil)
	_ UsersDistributionType = (*UsersDistributionTypeBasic)(nil)
)

// GetRewardsPoolAddress generates a rewards pool address from plan id.
func GetRewardsPoolAddress(planID uint64) sdk.AccAddress {
	return address.Module(ModuleName, []byte(fmt.Sprintf("rewards-plan-%d", planID)))
}

func NewRewardsPlan(
	id uint64, description string, serviceID uint32, amtPerDay sdk.Coins, startTime, endTime time.Time,
	poolsDistribution Distribution, operatorsDistribution Distribution,
	usersDistribution UsersDistribution) RewardsPlan {
	return RewardsPlan{
		ID:                    id,
		Description:           description,
		ServiceID:             serviceID,
		AmountPerDay:          amtPerDay,
		StartTime:             startTime,
		EndTime:               endTime,
		RewardsPool:           GetRewardsPoolAddress(id).String(),
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

func (plan RewardsPlan) Validate(unpacker codectypes.AnyUnpacker) error {
	if plan.ID == 0 {
		return fmt.Errorf("invalid plan ID: %d", plan.ID)
	}
	if len(plan.Description) > MaxRewardsPlanDescriptionLength {
		return fmt.Errorf("too long description")
	}
	if plan.ServiceID == 0 {
		return fmt.Errorf("invalid service ID: %d", plan.ServiceID)
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

	if plan.PoolsDistribution.DelegationType != restakingtypes.DELEGATION_TYPE_POOL {
		return fmt.Errorf("pools distribution has invalid delegation type: %v", plan.PoolsDistribution.DelegationType)
	}
	poolsDistrType, err := GetDistributionType(unpacker, plan.PoolsDistribution)
	if err != nil {
		return fmt.Errorf("get pools distribution type: %w", err)
	}
	err = poolsDistrType.Validate()
	if err != nil {
		return fmt.Errorf("invalid pools distribution type: %w", err)
	}

	if plan.OperatorsDistribution.DelegationType != restakingtypes.DELEGATION_TYPE_OPERATOR {
		return fmt.Errorf("operators distribution has invalid delegation type: %v", plan.OperatorsDistribution.DelegationType)
	}
	operatorsDistrType, err := GetDistributionType(unpacker, plan.OperatorsDistribution)
	if err != nil {
		return fmt.Errorf("get operators distribution type: %w", err)
	}
	err = operatorsDistrType.Validate()
	if err != nil {
		return fmt.Errorf("invalid operators distribution type: %w", err)
	}

	usersDistrType, err := GetUsersDistributionType(unpacker, plan.UsersDistribution)
	if err != nil {
		return fmt.Errorf("get users distribution type: %w", err)
	}
	err = usersDistrType.Validate()
	if err != nil {
		return fmt.Errorf("invalid users distribution type: %w", err)
	}

	return nil
}

type DistributionType interface {
	proto.Message
	isDistributionType()
	Validate() error
}

func GetDistributionType(unpacker codectypes.AnyUnpacker, distr Distribution) (DistributionType, error) {
	var distrType DistributionType
	err := unpacker.UnpackAny(distr.Type, &distrType)
	if err != nil {
		return nil, err
	}
	return distrType, nil
}

func newBasicDistribution(delType restakingtypes.DelegationType, weight uint32) Distribution {
	a, err := codectypes.NewAnyWithValue(&DistributionTypeBasic{})
	if err != nil {
		panic(err)
	}
	return Distribution{
		DelegationType: delType,
		Weight:         weight,
		Type:           a,
	}
}

func NewBasicPoolsDistribution(weight uint32) Distribution {
	return newBasicDistribution(restakingtypes.DELEGATION_TYPE_POOL, weight)
}

func NewBasicOperatorsDistribution(weight uint32) Distribution {
	return newBasicDistribution(restakingtypes.DELEGATION_TYPE_OPERATOR, weight)
}

func (t DistributionTypeBasic) Validate() error {
	return nil
}

func newWeightedDistribution(delType restakingtypes.DelegationType, weight uint32, weights []DistributionWeight) Distribution {
	a, err := codectypes.NewAnyWithValue(&DistributionTypeWeighted{Weights: weights})
	if err != nil {
		panic(err)
	}
	return Distribution{
		DelegationType: delType,
		Weight:         weight,
		Type:           a,
	}
}

func NewWeightedPoolsDistribution(weight uint32, weights []DistributionWeight) Distribution {
	return newWeightedDistribution(restakingtypes.DELEGATION_TYPE_POOL, weight, weights)
}

func NewWeightedOperatorsDistribution(weight uint32, weights []DistributionWeight) Distribution {
	return newWeightedDistribution(restakingtypes.DELEGATION_TYPE_OPERATOR, weight, weights)
}

func (t DistributionTypeWeighted) Validate() error {
	duplicate := utils.FindDuplicateFunc(t.Weights, func(a, b DistributionWeight) bool {
		return a.DelegationTargetID == b.DelegationTargetID
	})
	if duplicate != nil {
		return fmt.Errorf("duplicated weight for the same delegation target ID: %d", duplicate.DelegationTargetID)
	}

	for _, weight := range t.Weights {
		if weight.Weight == 0 {
			return fmt.Errorf("weight must be positive: %d", weight.Weight)
		}
		if weight.DelegationTargetID == 0 {
			return fmt.Errorf("invalid delegation target ID: %d", weight.DelegationTargetID)
		}
	}
	return nil
}

func NewDistributionWeight(targetID, weight uint32) DistributionWeight {
	return DistributionWeight{
		DelegationTargetID: targetID,
		Weight:             weight,
	}
}

func newEgalitarianDistribution(delType restakingtypes.DelegationType, weight uint32) Distribution {
	a, err := codectypes.NewAnyWithValue(&DistributionTypeEgalitarian{})
	if err != nil {
		panic(err)
	}
	return Distribution{
		DelegationType: delType,
		Weight:         weight,
		Type:           a,
	}
}

func NewEgalitarianPoolsDistribution(weight uint32) Distribution {
	return newEgalitarianDistribution(restakingtypes.DELEGATION_TYPE_POOL, weight)
}

func NewEgalitarianOperatorsDistribution(weight uint32) Distribution {
	return newEgalitarianDistribution(restakingtypes.DELEGATION_TYPE_OPERATOR, weight)
}

func (t DistributionTypeEgalitarian) Validate() error {
	return nil
}

func (t DistributionTypeBasic) isDistributionType() {}

func (t DistributionTypeWeighted) isDistributionType() {}

func (t DistributionTypeEgalitarian) isDistributionType() {}

type UsersDistributionType interface {
	proto.Message
	isUsersDistributionType()
	Validate() error
}

func GetUsersDistributionType(unpacker codectypes.AnyUnpacker, distr UsersDistribution) (UsersDistributionType, error) {
	var distrType UsersDistributionType
	err := unpacker.UnpackAny(distr.Type, &distrType)
	if err != nil {
		return nil, err
	}
	return distrType, nil
}

func NewBasicUsersDistribution(weight uint32) UsersDistribution {
	a, err := codectypes.NewAnyWithValue(&UsersDistributionTypeBasic{})
	if err != nil {
		panic(err)
	}
	return UsersDistribution{
		Weight: weight,
		Type:   a,
	}
}

func (t UsersDistributionTypeBasic) Validate() error {
	return nil
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

// create a new DelegationDelegatorReward
func NewDelegationDelegatorReward(
	delType restakingtypes.DelegationType, targetID uint32, rewards DecPools) DelegationDelegatorReward {
	return DelegationDelegatorReward{
		DelegationType:     delType,
		DelegationTargetID: targetID,
		Reward:             rewards,
	}
}
