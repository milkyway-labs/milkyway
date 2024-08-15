package types_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/utils"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
)

func TestRewardsPlan_Validate(t *testing.T) {
	ir := codectypes.NewInterfaceRegistry()
	types.RegisterInterfaces(ir)

	testCases := []struct {
		name        string
		malleate    func(plan *types.RewardsPlan)
		expectedErr string
	}{
		{
			name:        "valid rewards plan",
			malleate:    nil,
			expectedErr: "",
		},
		{
			name: "invalid plan ID returns error",
			malleate: func(plan *types.RewardsPlan) {
				plan.ID = 0
			},
			expectedErr: "invalid plan ID: 0",
		},
		{
			name: "too long description returns error",
			malleate: func(plan *types.RewardsPlan) {
				plan.Description = strings.Repeat("A", types.MaxRewardsPlanDescriptionLength+1)
			},
			expectedErr: "too long description",
		},
		{
			name: "invalid service ID returns error",
			malleate: func(plan *types.RewardsPlan) {
				plan.ServiceID = 0
			},
			expectedErr: "invalid service ID: 0",
		},
		{
			name: "invalid amount per day returns error",
			malleate: func(plan *types.RewardsPlan) {
				plan.AmountPerDay = sdk.Coins{sdk.Coin{Denom: "umilk", Amount: math.ZeroInt()}}
			},
			expectedErr: "invalid amount per day: coin 0umilk amount is not positive",
		},
		{
			name: "end time must be after start time",
			malleate: func(plan *types.RewardsPlan) {
				plan.StartTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
				plan.EndTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
			},
			expectedErr: "end time must be after start time: 2024-01-01T00:00:00Z <= 2024-01-01T00:00:00Z",
		},
		{
			name: "invalid rewards pool returns error",
			malleate: func(plan *types.RewardsPlan) {
				plan.RewardsPool = "invalid"
			},
			expectedErr: "invalid rewards pool: decoding bech32 failed: invalid bech32 string length 7",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			plan := types.NewRewardsPlan(
				1, "Plan", 1, utils.MustParseCoins("100_000000umilk"),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				types.NewBasicPoolsDistribution(0), types.NewBasicOperatorsDistribution(0),
				types.NewBasicUsersDistribution(0))
			if tc.malleate != nil {
				tc.malleate(&plan)
			}
			err := plan.Validate(ir)
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}

func TestRewardsPlan_IsActiveAt(t *testing.T) {
	startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	plan := types.NewRewardsPlan(
		1, "Plan", 1, sdk.NewCoins(sdk.NewInt64Coin("umilk", 100_000000)), startTime, endTime,
		types.NewBasicPoolsDistribution(1), types.NewBasicOperatorsDistribution(1), types.NewBasicUsersDistribution(1))

	testCases := []struct {
		name     string
		date     time.Time
		isActive bool
	}{
		{
			name:     "plan is active at start time",
			date:     startTime,
			isActive: true,
		},
		{
			name:     "plan is inactive at end time",
			date:     endTime,
			isActive: false,
		},
		{
			name:     "plan is inactive before start time",
			date:     startTime.AddDate(0, 0, -1),
			isActive: false,
		},
		{
			name:     "plan is active between start time and end time",
			date:     startTime.AddDate(0, 0, 1),
			isActive: true,
		},
		{
			name:     "plan is inactive after end time",
			date:     endTime.AddDate(0, 0, 1),
			isActive: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isActive := plan.IsActiveAt(tc.date)
			require.Equal(t, tc.isActive, isActive)
		})
	}
}

func TestDistribution_Validate(t *testing.T) {
	testCases := []struct {
		name        string
		distrType   types.DistributionType
		expectedErr string
	}{
		{
			name:        "basic distribution type returns no error",
			distrType:   &types.DistributionTypeBasic{},
			expectedErr: "",
		},
		{
			name: "invalid delegation target ID returns error",
			distrType: &types.DistributionTypeWeighted{
				Weights: []types.DistributionWeight{
					types.NewDistributionWeight(0, 1),
				},
			},
			expectedErr: "invalid delegation target ID: 0",
		},
		{
			name: "invalid weight returns error",
			distrType: &types.DistributionTypeWeighted{
				Weights: []types.DistributionWeight{
					types.NewDistributionWeight(1, 0),
				},
			},
			expectedErr: "weight must be positive: 0",
		},
		{
			name: "duplicated delegation target ID returns error",
			distrType: &types.DistributionTypeWeighted{
				Weights: []types.DistributionWeight{
					types.NewDistributionWeight(1, 1),
					types.NewDistributionWeight(1, 2),
				},
			},
			expectedErr: "duplicated weight for the same delegation target ID: 1",
		},
		{
			name:        "egalitarian distribution type returns no error",
			distrType:   &types.DistributionTypeEgalitarian{},
			expectedErr: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.distrType.Validate()
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}

func TestUsersDistribution_Validate(t *testing.T) {
	testCases := []struct {
		name        string
		distrType   types.UsersDistributionType
		expectedErr string
	}{
		{
			name:        "basic distribution type returns no error",
			distrType:   &types.UsersDistributionTypeBasic{},
			expectedErr: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.distrType.Validate()
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}
