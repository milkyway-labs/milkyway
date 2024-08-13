package types_test

import (
	"strings"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

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
			"valid rewards plan",
			nil,
			"",
		},
		{
			"invalid plan ID returns error",
			func(plan *types.RewardsPlan) {
				plan.ID = 0
			},
			"invalid plan ID: 0",
		},
		{
			"too long description returns error",
			func(plan *types.RewardsPlan) {
				plan.Description = strings.Repeat("A", types.MaxRewardsPlanDescriptionLength+1)
			},
			"too long description",
		},
		{
			"invalid service ID returns error",
			func(plan *types.RewardsPlan) {
				plan.ServiceID = 0
			},
			"invalid service ID: 0",
		},
		{
			"invalid amount per day returns error",
			func(plan *types.RewardsPlan) {
				plan.AmountPerDay = sdk.Coins{sdk.Coin{Denom: "umilk", Amount: math.ZeroInt()}}
			},
			"invalid amount per day: coin 0umilk amount is not positive",
		},
		{
			"end time must be after start time",
			func(plan *types.RewardsPlan) {
				plan.StartTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
				plan.EndTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
			},
			"end time must be after start time: 2024-01-01T00:00:00Z <= 2024-01-01T00:00:00Z",
		},
		{
			"invalid rewards pool returns error",
			func(plan *types.RewardsPlan) {
				plan.RewardsPool = "invalid"
			},
			"invalid rewards pool: decoding bech32 failed: invalid bech32 string length 7",
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
			"plan is active at start time",
			startTime,
			true,
		},
		{
			"plan is inactive at end time",
			endTime,
			false,
		},
		{
			"plan is inactive before start time",
			startTime.AddDate(0, 0, -1),
			false,
		},
		{
			"plan is active between start time and end time",
			startTime.AddDate(0, 0, 1),
			true,
		},
		{
			"plan is inactive after end time",
			endTime.AddDate(0, 0, 1),
			false,
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
			"basic distribution type returns no error",
			&types.DistributionTypeBasic{},
			"",
		},
		{
			"invalid delegation target ID returns error",
			&types.DistributionTypeWeighted{
				Weights: []types.DistributionWeight{
					{
						DelegationTargetID: 0,
						Weight:             1,
					},
				},
			},
			"invalid delegation target ID: 0",
		},
		{
			"invalid weight returns error",
			&types.DistributionTypeWeighted{
				Weights: []types.DistributionWeight{
					{
						DelegationTargetID: 1,
						Weight:             0,
					},
				},
			},
			"weight must be positive: 0",
		},
		{
			"duplicated delegation target ID returns error",
			&types.DistributionTypeWeighted{
				Weights: []types.DistributionWeight{
					{
						DelegationTargetID: 1,
						Weight:             1,
					},
					{
						DelegationTargetID: 1,
						Weight:             2,
					},
				},
			},
			"duplicated weight for the same delegation target ID: 1",
		},
		{
			"egalitarian distribution type returns no error",
			&types.DistributionTypeEgalitarian{},
			"",
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
			"basic distribution type returns no error",
			&types.UsersDistributionTypeBasic{},
			"",
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
