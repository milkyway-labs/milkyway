package types

import (
	"fmt"
	"strconv"
	"time"

	coreaddress "cosmossdk.io/core/address"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/gogoproto/proto"

	"github.com/milkyway-labs/milkyway/utils"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
)

// MaxRewardsPlanDescriptionLength is the maximum length of a rewards plan description.
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

// NewRewardsPlan creates a new rewards plan.
func NewRewardsPlan(
	id uint64,
	description string,
	serviceID uint32,
	amtPerDay sdk.Coins,
	startTime,
	endTime time.Time,
	poolsDistribution Distribution,
	operatorsDistribution Distribution,
	usersDistribution UsersDistribution,
) RewardsPlan {
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

// IsActiveAt returns true if the plan is active at the given time.
func (plan RewardsPlan) IsActiveAt(t time.Time) bool {
	return !plan.StartTime.After(t) && plan.EndTime.After(t)
}

// MustGetRewardsPoolAddress returns the rewards pool address.
func (plan RewardsPlan) MustGetRewardsPoolAddress(addressCodec coreaddress.Codec) sdk.AccAddress {
	addr, err := addressCodec.StringToBytes(plan.RewardsPool)
	if err != nil {
		panic(err)
	}
	return addr
}

// Validate checks the plan for validity.
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

	err := plan.AmountPerDay.Validate()
	if err != nil {
		return fmt.Errorf("invalid amount per day: %w", err)
	}

	if !plan.EndTime.After(plan.StartTime) {
		return fmt.Errorf(
			"end time must be after start time: %s <= %s",
			plan.EndTime.Format(time.RFC3339),
			plan.StartTime.Format(time.RFC3339),
		)
	}

	_, err = sdk.AccAddressFromBech32(plan.RewardsPool)
	if err != nil {
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

// ParseRewardsPlanID tries parsing the given value as an rewards plan id
func ParseRewardsPlanID(value string) (uint64, error) {
	planID, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid rewards plan ID: %s", value)
	}
	return planID, nil
}

// --------------------------------------------------------------------------------------------------------------------

// DistributionType represents a generic distribution type
type DistributionType interface {
	proto.Message
	isDistributionType()
	Validate() error
}

// GetDistributionType returns the distribution type from the distribution
func GetDistributionType(unpacker codectypes.AnyUnpacker, distr Distribution) (DistributionType, error) {
	var distrType DistributionType
	err := unpacker.UnpackAny(distr.Type, &distrType)
	if err != nil {
		return nil, err
	}
	return distrType, nil
}

// NewDistributionWeight creates a new distribution weight
func NewDistributionWeight(targetID, weight uint32) DistributionWeight {
	return DistributionWeight{
		DelegationTargetID: targetID,
		Weight:             weight,
	}
}

// --------------------------------------------------------------------------------------------------------------------

// newBasicDistribution creates a new basic distribution
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

// NewBasicPoolsDistribution creates a new basic pools distribution
func NewBasicPoolsDistribution(weight uint32) Distribution {
	return newBasicDistribution(restakingtypes.DELEGATION_TYPE_POOL, weight)
}

// NewBasicOperatorsDistribution creates a new basic operators distribution
func NewBasicOperatorsDistribution(weight uint32) Distribution {
	return newBasicDistribution(restakingtypes.DELEGATION_TYPE_OPERATOR, weight)
}

// Validate checks the distribution for validity
func (t DistributionTypeBasic) Validate() error {
	return nil
}

// isDistributionType is a marker function
func (t DistributionTypeBasic) isDistributionType() {}

// --------------------------------------------------------------------------------------------------------------------

// newWeightedDistribution creates a new weighted distribution
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

// NewWeightedPoolsDistribution creates a new weighted pools distribution
func NewWeightedPoolsDistribution(weight uint32, weights []DistributionWeight) Distribution {
	return newWeightedDistribution(restakingtypes.DELEGATION_TYPE_POOL, weight, weights)
}

// NewWeightedOperatorsDistribution creates a new weighted operators distribution
func NewWeightedOperatorsDistribution(weight uint32, weights []DistributionWeight) Distribution {
	return newWeightedDistribution(restakingtypes.DELEGATION_TYPE_OPERATOR, weight, weights)
}

// Validate checks the distribution for validity
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

// isDistributionType is a marker function
func (t DistributionTypeWeighted) isDistributionType() {}

// --------------------------------------------------------------------------------------------------------------------

// newEgalitarianDistribution creates a new egalitarian distribution
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

// NewEgalitarianPoolsDistribution creates a new egalitarian pools distribution
func NewEgalitarianPoolsDistribution(weight uint32) Distribution {
	return newEgalitarianDistribution(restakingtypes.DELEGATION_TYPE_POOL, weight)
}

// NewEgalitarianOperatorsDistribution creates a new egalitarian operators distribution
func NewEgalitarianOperatorsDistribution(weight uint32) Distribution {
	return newEgalitarianDistribution(restakingtypes.DELEGATION_TYPE_OPERATOR, weight)
}

// Validate checks the distribution for validity
func (t DistributionTypeEgalitarian) Validate() error {
	return nil
}

// isDistributionType is a marker function
func (t DistributionTypeEgalitarian) isDistributionType() {}

// --------------------------------------------------------------------------------------------------------------------

// UsersDistributionType represents a generic users distribution type
type UsersDistributionType interface {
	proto.Message
	isUsersDistributionType()
	Validate() error
}

// GetUsersDistributionType returns the users distribution type from the distribution
func GetUsersDistributionType(unpacker codectypes.AnyUnpacker, distr UsersDistribution) (UsersDistributionType, error) {
	var distrType UsersDistributionType
	err := unpacker.UnpackAny(distr.Type, &distrType)
	if err != nil {
		return nil, err
	}
	return distrType, nil
}

// NewBasicUsersDistribution creates a new basic users distribution
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

// Validate checks the users distribution for validity
func (t UsersDistributionTypeBasic) Validate() error {
	return nil
}

// isUsersDistributionType is a marker function
func (t UsersDistributionTypeBasic) isUsersDistributionType() {}

// --------------------------------------------------------------------------------------------------------------------

// InitialAccumulatedCommission returns the initial accumulated commission (zero)
func InitialAccumulatedCommission() AccumulatedCommission {
	return AccumulatedCommission{}
}

// NewHistoricalRewards creates a new historical rewards
func NewHistoricalRewards(cumulativeRewardRatios ServicePools, referenceCount uint32) HistoricalRewards {
	return HistoricalRewards{
		CumulativeRewardRatios: cumulativeRewardRatios,
		ReferenceCount:         referenceCount,
	}
}

// NewCurrentRewards creates a new current rewards
func NewCurrentRewards(rewards ServicePools, period uint64) CurrentRewards {
	return CurrentRewards{
		Rewards: rewards,
		Period:  period,
	}
}

// NewDelegatorStartingInfo creates a new delegator starting info
func NewDelegatorStartingInfo(previousPeriod uint64, stakes sdk.DecCoins, height uint64) DelegatorStartingInfo {
	return DelegatorStartingInfo{
		PreviousPeriod: previousPeriod,
		Stakes:         stakes,
		Height:         height,
	}
}

// NewDelegationDelegatorReward creates a new delegation delegator reward
func NewDelegationDelegatorReward(
	delType restakingtypes.DelegationType, targetID uint32, rewards DecPools,
) DelegationDelegatorReward {
	return DelegationDelegatorReward{
		DelegationType:     delType,
		DelegationTargetID: targetID,
		Reward:             rewards,
	}
}

// NewPoolServiceTotalDelegatorShares creates a new pool service total delegator shares
func NewPoolServiceTotalDelegatorShares(poolID, serviceID uint32, shares sdk.DecCoins) PoolServiceTotalDelegatorShares {
	return PoolServiceTotalDelegatorShares{
		PoolID:    poolID,
		ServiceID: serviceID,
		Shares:    shares,
	}
}

// Validate validates the pool service total delegator shares
func (shares PoolServiceTotalDelegatorShares) Validate() error {
	if shares.PoolID == 0 {
		return fmt.Errorf("pool ID must not be 0")
	}
	if shares.ServiceID == 0 {
		return fmt.Errorf("service ID must not be 0")
	}
	err := shares.Shares.Validate()
	if err != nil {
		return fmt.Errorf("invalid pool service total delegator shares: %w", err)
	}
	return nil
}
