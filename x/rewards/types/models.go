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
	_ PoolsDistributionType     = (*PoolsDistributionTypeBasic)(nil)
	_ PoolsDistributionType     = (*PoolsDistributionTypeWeighted)(nil)
	_ PoolsDistributionType     = (*PoolsDistributionTypeEgalitarian)(nil)
	_ OperatorsDistributionType = (*OperatorsDistributionTypeBasic)(nil)
	_ OperatorsDistributionType = (*OperatorsDistributionTypeWeighted)(nil)
	_ OperatorsDistributionType = (*OperatorsDistributionTypeEgalitarian)(nil)
	_ UsersDistributionType     = (*UsersDistributionTypeBasic)(nil)
)

// GetRewardsPoolAddress generates a rewards pool address from plan id.
func GetRewardsPoolAddress(planID uint64) sdk.AccAddress {
	return address.Module(ModuleName, []byte(fmt.Sprintf("rewards-plan-%d", planID)))
}

func NewRewardsPlan(
	id uint64, description string, serviceID uint32, amtPerDay sdk.Coins, startTime, endTime time.Time,
	poolsDistribution PoolsDistribution, operatorsDistribution OperatorsDistribution,
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

func (plan RewardsPlan) TotalWeights() uint32 {
	return plan.PoolsDistribution.Weight + plan.OperatorsDistribution.Weight + plan.UsersDistribution.Weight
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

	poolsDistrType, err := GetPoolsDistributionType(unpacker, plan.PoolsDistribution)
	if err != nil {
		return fmt.Errorf("get pools distribution type: %w", err)
	}
	err = poolsDistrType.Validate()
	if err != nil {
		return fmt.Errorf("invalid pools distribution type: %w", err)
	}

	operatorsDistrType, err := GetOperatorsDistributionType(unpacker, plan.OperatorsDistribution)
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

type PoolsDistributionType interface {
	proto.Message
	isPoolsDistributionType()
	Validate() error
}

func GetPoolsDistributionType(unpacker codectypes.AnyUnpacker, distr PoolsDistribution) (PoolsDistributionType, error) {
	var distrType PoolsDistributionType
	err := unpacker.UnpackAny(distr.Type, &distrType)
	if err != nil {
		return nil, err
	}
	return distrType, nil
}

func NewBasicPoolsDistribution(weight uint32) PoolsDistribution {
	a, err := codectypes.NewAnyWithValue(&PoolsDistributionTypeBasic{})
	if err != nil {
		panic(err)
	}
	return PoolsDistribution{
		Weight: weight,
		Type:   a,
	}
}

func (t PoolsDistributionTypeBasic) Validate() error {
	return nil
}

func NewWeightedPoolsDistribution(weight uint32, poolWeights []PoolDistributionWeight) PoolsDistribution {
	a, err := codectypes.NewAnyWithValue(&PoolsDistributionTypeWeighted{Weights: poolWeights})
	if err != nil {
		panic(err)
	}
	return PoolsDistribution{
		Weight: weight,
		Type:   a,
	}
}

func (t PoolsDistributionTypeWeighted) Validate() error {
	duplicate := utils.FindDuplicateFunc(t.Weights, func(a, b PoolDistributionWeight) bool {
		return a.PoolID == b.PoolID
	})
	if duplicate != nil {
		return fmt.Errorf("duplicated weight for the same pool ID: %d", duplicate.PoolID)
	}

	for _, weight := range t.Weights {
		if weight.Weight == 0 {
			return fmt.Errorf("weight must be positive: %d", weight.Weight)
		}
		if weight.PoolID == 0 {
			return fmt.Errorf("invalid pool ID: %d", weight.PoolID)
		}
	}
	return nil
}

func NewEgalitarianPoolsDistribution(weight uint32) PoolsDistribution {
	a, err := codectypes.NewAnyWithValue(&PoolsDistributionTypeEgalitarian{})
	if err != nil {
		panic(err)
	}
	return PoolsDistribution{
		Weight: weight,
		Type:   a,
	}
}

func (t PoolsDistributionTypeEgalitarian) Validate() error {
	return nil
}

func NewPoolDistributionWeight(poolID, weight uint32) PoolDistributionWeight {
	return PoolDistributionWeight{
		PoolID: poolID,
		Weight: weight,
	}
}

func (t PoolsDistributionTypeBasic) isPoolsDistributionType() {}

func (t PoolsDistributionTypeWeighted) isPoolsDistributionType() {}

func (t PoolsDistributionTypeEgalitarian) isPoolsDistributionType() {}

type OperatorsDistributionType interface {
	proto.Message
	isOperatorsDistributionType()
	Validate() error
}

func GetOperatorsDistributionType(unpacker codectypes.AnyUnpacker, distr OperatorsDistribution) (OperatorsDistributionType, error) {
	var distrType OperatorsDistributionType
	err := unpacker.UnpackAny(distr.Type, &distrType)
	if err != nil {
		return nil, err
	}
	return distrType, nil
}

func NewBasicOperatorsDistribution(weight uint32) OperatorsDistribution {
	a, err := codectypes.NewAnyWithValue(&OperatorsDistributionTypeBasic{})
	if err != nil {
		panic(err)
	}
	return OperatorsDistribution{
		Weight: weight,
		Type:   a,
	}
}

func (t OperatorsDistributionTypeBasic) Validate() error {
	return nil
}

func NewWeightedOperatorsDistribution(
	weight uint32, operatorWeights []OperatorDistributionWeight) OperatorsDistribution {
	a, err := codectypes.NewAnyWithValue(&OperatorsDistributionTypeWeighted{Weights: operatorWeights})
	if err != nil {
		panic(err)
	}
	return OperatorsDistribution{
		Weight: weight,
		Type:   a,
	}
}

func (t OperatorsDistributionTypeWeighted) Validate() error {
	duplicate := utils.FindDuplicateFunc(t.Weights, func(a, b OperatorDistributionWeight) bool {
		return a.OperatorID == b.OperatorID
	})
	if duplicate != nil {
		return fmt.Errorf("duplicated weight for the same operator ID: %d", duplicate.OperatorID)
	}

	for _, weight := range t.Weights {
		if weight.Weight == 0 {
			return fmt.Errorf("weight must be positive: %d", weight.Weight)
		}
		if weight.OperatorID == 0 {
			return fmt.Errorf("invalid operator ID: %d", weight.OperatorID)
		}
	}
	return nil
}

func NewEgalitarianOperatorsDistribution(weight uint32) OperatorsDistribution {
	a, err := codectypes.NewAnyWithValue(&OperatorsDistributionTypeEgalitarian{})
	if err != nil {
		panic(err)
	}
	return OperatorsDistribution{
		Weight: weight,
		Type:   a,
	}
}

func (t OperatorsDistributionTypeEgalitarian) Validate() error {
	return nil
}

func NewOperatorDistributionWeight(operatorID, weight uint32) OperatorDistributionWeight {
	return OperatorDistributionWeight{
		OperatorID: operatorID,
		Weight:     weight,
	}
}

func (t OperatorsDistributionTypeBasic) isOperatorsDistributionType() {}

func (t OperatorsDistributionTypeWeighted) isOperatorsDistributionType() {}

func (t OperatorsDistributionTypeEgalitarian) isOperatorsDistributionType() {}

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
