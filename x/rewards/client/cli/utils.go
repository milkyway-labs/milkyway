package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
)

type rewardsPlanJSON struct {
	Description           string           `json:"description"`
	ServiceID             uint32           `json:"service_id"`
	AmountPerDay          string           `json:"amount_per_day"`
	StartTime             time.Time        `json:"start_time"`
	EndTime               time.Time        `json:"end_time"`
	PoolsDistribution     distributionJSON `json:"pools_distribution"`
	OperatorsDistribution distributionJSON `json:"operators_distribution"`
	UsersDistribution     distributionJSON `json:"users_distribution"`
}

type distributionJSON struct {
	Weight uint32          `json:"weight"`
	Type   json.RawMessage `json:"type"`
}

type RewardsPlan struct {
	Description           string
	ServiceID             uint32
	AmountPerDay          sdk.Coins
	StartTime             time.Time
	EndTime               time.Time
	PoolsDistribution     types.Distribution
	OperatorsDistribution types.Distribution
	UsersDistribution     types.UsersDistribution
}

func (plan *RewardsPlan) Validate(unpacker codectypes.AnyUnpacker) error {
	if len(plan.Description) > types.MaxRewardsPlanDescriptionLength {
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

	if plan.PoolsDistribution.DelegationType != restakingtypes.DELEGATION_TYPE_POOL {
		return fmt.Errorf("pools distribution has invalid delegation type: %v", plan.PoolsDistribution.DelegationType)
	}

	poolsDistrType, err := types.GetDistributionType(unpacker, plan.PoolsDistribution)
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

	operatorsDistrType, err := types.GetDistributionType(unpacker, plan.OperatorsDistribution)
	if err != nil {
		return fmt.Errorf("get operators distribution type: %w", err)
	}

	err = operatorsDistrType.Validate()
	if err != nil {
		return fmt.Errorf("invalid operators distribution type: %w", err)
	}

	usersDistrType, err := types.GetUsersDistributionType(unpacker, plan.UsersDistribution)
	if err != nil {
		return fmt.Errorf("get users distribution type: %w", err)
	}

	err = usersDistrType.Validate()
	if err != nil {
		return fmt.Errorf("invalid users distribution type: %w", err)
	}

	return nil
}

// ParseRewardsPlan parse a RewardsPlan from a json file.
func ParseRewardsPlan(cdc codec.Codec, path string) (RewardsPlan, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return RewardsPlan{}, err
	}

	var rewardsPlanJSON rewardsPlanJSON
	err = json.Unmarshal(contents, &rewardsPlanJSON)
	if err != nil {
		return RewardsPlan{}, err
	}

	// Parse the pools distribution types
	poolsDistributionType, err := parseDistributionType(cdc, rewardsPlanJSON.PoolsDistribution.Type)
	if err != nil {
		return RewardsPlan{}, fmt.Errorf("invalid pool distribution type %w", err)
	}
	poolsDistribution := types.Distribution{
		DelegationType: restakingtypes.DELEGATION_TYPE_POOL,
		Type:           poolsDistributionType,
		Weight:         rewardsPlanJSON.PoolsDistribution.Weight,
	}

	// Parse the operators distribution type
	operatorsDistributionType, err := parseDistributionType(cdc, rewardsPlanJSON.OperatorsDistribution.Type)
	if err != nil {
		return RewardsPlan{}, fmt.Errorf("invalid operator distribution type: %w", err)
	}
	operatorsDistribution := types.Distribution{
		DelegationType: restakingtypes.DELEGATION_TYPE_OPERATOR,
		Type:           operatorsDistributionType,
		Weight:         rewardsPlanJSON.OperatorsDistribution.Weight,
	}

	// Parse the users distribution type
	usersDistributionType, err := parseUserDistributionType(cdc, rewardsPlanJSON.UsersDistribution.Type)
	if err != nil {
		return RewardsPlan{}, fmt.Errorf("invalid user distribution type: %w", err)
	}
	usersDistribution := types.UsersDistribution{
		Type:   usersDistributionType,
		Weight: rewardsPlanJSON.UsersDistribution.Weight,
	}

	amountPerDay, err := sdk.ParseCoinsNormalized(rewardsPlanJSON.AmountPerDay)
	if err != nil {
		return RewardsPlan{}, fmt.Errorf("invalid amount per day: %w", err)
	}

	return RewardsPlan{
		Description:           rewardsPlanJSON.Description,
		ServiceID:             rewardsPlanJSON.ServiceID,
		AmountPerDay:          amountPerDay,
		StartTime:             rewardsPlanJSON.StartTime,
		EndTime:               rewardsPlanJSON.EndTime,
		PoolsDistribution:     poolsDistribution,
		OperatorsDistribution: operatorsDistribution,
		UsersDistribution:     usersDistribution,
	}, nil
}

// parseDistributionType parses a types.DistributionType from a json.RawMessage.
func parseDistributionType(cdc codec.Codec, rawJSON json.RawMessage) (*codectypes.Any, error) {
	var poolsDistributionType types.DistributionType
	err := cdc.UnmarshalInterfaceJSON(rawJSON, &poolsDistributionType)
	if err != nil {
		return nil, fmt.Errorf("invalid distribution type: %w", err)
	}
	return codectypes.NewAnyWithValue(poolsDistributionType)
}

// parseUserDistributionType parses a types.UsersDistributionType from a json.RawMessage.
func parseUserDistributionType(cdc codec.Codec, rawJSON json.RawMessage) (*codectypes.Any, error) {
	var userDistributionType types.UsersDistributionType
	err := cdc.UnmarshalInterfaceJSON(rawJSON, &userDistributionType)
	if err != nil {
		return nil, fmt.Errorf("invalid user distribution type: %w", err)
	}
	return codectypes.NewAnyWithValue(userDistributionType)
}
