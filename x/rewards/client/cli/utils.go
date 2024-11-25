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
	FeeAmount             string           `json:"fee_amount"`
}

type distributionJSON struct {
	Weight uint32          `json:"weight"`
	Type   json.RawMessage `json:"type"`
}

type ParsedRewardsPlan struct {
	types.RewardsPlan
	FeeAmount sdk.Coins
}

// ParseRewardsPlan parse a RewardsPlan from a json file.
func ParseRewardsPlan(cdc codec.Codec, path string) (ParsedRewardsPlan, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return ParsedRewardsPlan{}, err
	}

	var rewardsPlanJSON rewardsPlanJSON
	err = json.Unmarshal(contents, &rewardsPlanJSON)
	if err != nil {
		return ParsedRewardsPlan{}, err
	}

	// Parse the pools distribution types
	poolsDistributionType, err := parseDistributionType(cdc, rewardsPlanJSON.PoolsDistribution.Type)
	if err != nil {
		return ParsedRewardsPlan{}, fmt.Errorf("invalid pool distribution type %w", err)
	}
	poolsDistribution := types.Distribution{
		DelegationType: restakingtypes.DELEGATION_TYPE_POOL,
		Type:           poolsDistributionType,
		Weight:         rewardsPlanJSON.PoolsDistribution.Weight,
	}

	// Parse the operators distribution type
	operatorsDistributionType, err := parseDistributionType(cdc, rewardsPlanJSON.OperatorsDistribution.Type)
	if err != nil {
		return ParsedRewardsPlan{}, fmt.Errorf("invalid operator distribution type: %w", err)
	}
	operatorsDistribution := types.Distribution{
		DelegationType: restakingtypes.DELEGATION_TYPE_OPERATOR,
		Type:           operatorsDistributionType,
		Weight:         rewardsPlanJSON.OperatorsDistribution.Weight,
	}

	// Parse the users distribution type
	usersDistributionType, err := parseUserDistributionType(cdc, rewardsPlanJSON.UsersDistribution.Type)
	if err != nil {
		return ParsedRewardsPlan{}, fmt.Errorf("invalid user distribution type: %w", err)
	}
	usersDistribution := types.UsersDistribution{
		Type:   usersDistributionType,
		Weight: rewardsPlanJSON.UsersDistribution.Weight,
	}

	amountPerDay, err := sdk.ParseCoinsNormalized(rewardsPlanJSON.AmountPerDay)
	if err != nil {
		return ParsedRewardsPlan{}, fmt.Errorf("invalid amount per day: %w", err)
	}

	feeAmount, err := sdk.ParseCoinsNormalized(rewardsPlanJSON.FeeAmount)
	if err != nil {
		return ParsedRewardsPlan{}, fmt.Errorf("invalid fee amount: %w", err)
	}

	return ParsedRewardsPlan{
		RewardsPlan: types.NewRewardsPlan(1,
			rewardsPlanJSON.Description,
			rewardsPlanJSON.ServiceID,
			amountPerDay,
			rewardsPlanJSON.StartTime,
			rewardsPlanJSON.EndTime,
			poolsDistribution,
			operatorsDistribution,
			usersDistribution,
		),
		FeeAmount: feeAmount,
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
