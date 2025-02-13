package types_test

import (
	"strings"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectestutil "github.com/cosmos/cosmos-sdk/codec/testutil"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v9/utils"
	"github.com/milkyway-labs/milkyway/v9/x/rewards/types"
)

func TestRewardsPlan_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		plan      types.RewardsPlan
		shouldErr bool
	}{
		{
			name: "valid plan returns no error",
			plan: types.NewRewardsPlan(
				1,
				"Plan",
				1,
				utils.MustParseCoin("100_000000umilk"),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				types.NewBasicPoolsDistribution(0),
				types.NewBasicOperatorsDistribution(0),
				types.NewBasicUsersDistribution(0),
			),
			shouldErr: false,
		},
		{
			name: "invalid plan ID returns error",
			plan: types.NewRewardsPlan(
				0,
				"Plan",
				1,
				utils.MustParseCoin("100_000000umilk"),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				types.NewBasicPoolsDistribution(0),
				types.NewBasicOperatorsDistribution(0),
				types.NewBasicUsersDistribution(0),
			),
			shouldErr: true,
		},
		{
			name: "too long description returns error",
			plan: types.NewRewardsPlan(
				1,
				strings.Repeat("A", types.MaxRewardsPlanDescriptionLength+1),
				1,
				utils.MustParseCoin("100_000000umilk"),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				types.NewBasicPoolsDistribution(0),
				types.NewBasicOperatorsDistribution(0),
				types.NewBasicUsersDistribution(0),
			),
			shouldErr: true,
		},
		{
			name: "invalid service ID returns error",
			plan: types.NewRewardsPlan(
				1,
				"Plan",
				0,
				utils.MustParseCoin("100_000000umilk"),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				types.NewBasicPoolsDistribution(0),
				types.NewBasicOperatorsDistribution(0),
				types.NewBasicUsersDistribution(0),
			),
			shouldErr: true,
		},
		{
			name: "invalid amount per day returns error",
			plan: types.NewRewardsPlan(
				1,
				"Plan",
				1,
				sdk.Coin{Denom: "umilk", Amount: math.ZeroInt()},
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				types.NewBasicPoolsDistribution(0),
				types.NewBasicOperatorsDistribution(0),
				types.NewBasicUsersDistribution(0),
			),
			shouldErr: true,
		},
		{
			name: "end time must be after start time",
			plan: types.NewRewardsPlan(
				1,
				"Plan",
				1,
				utils.MustParseCoin("100_000000umilk"),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				types.NewBasicPoolsDistribution(0),
				types.NewBasicOperatorsDistribution(0),
				types.NewBasicUsersDistribution(0),
			),
			shouldErr: true,
		},
		{
			name: "invalid rewards pool returns error",
			plan: types.RewardsPlan{
				ID:                    1,
				Description:           "Plan",
				ServiceID:             1,
				AmountPerDay:          utils.MustParseCoin("100_000000umilk"),
				StartTime:             time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				EndTime:               time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				RewardsPool:           "invalid",
				PoolsDistribution:     types.NewBasicPoolsDistribution(0),
				OperatorsDistribution: types.NewBasicOperatorsDistribution(0),
				UsersDistribution:     types.NewBasicUsersDistribution(0),
			},
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			interfaceRegistry := codectestutil.CodecOptions{AccAddressPrefix: "cosmo", ValAddressPrefix: "cosmovaloper"}.NewInterfaceRegistry()
			cdc := codec.NewProtoCodec(interfaceRegistry)

			err := tc.plan.Validate(cdc)
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRewardsPlan_IsActiveAt(t *testing.T) {
	testCases := []struct {
		name      string
		plan      types.RewardsPlan
		date      time.Time
		expActive bool
	}{
		{
			name: "plan is inactive before start time",
			plan: types.NewRewardsPlan(
				1,
				"Plan",
				1,
				sdk.NewInt64Coin("umilk", 100_000000),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				types.NewBasicPoolsDistribution(1),
				types.NewBasicOperatorsDistribution(1),
				types.NewBasicUsersDistribution(1),
			),
			date:      time.Date(2023, 12, 31, 23, 59, 59, 999, time.UTC),
			expActive: false,
		},
		{
			name: "plan is active at start time",
			plan: types.NewRewardsPlan(
				1,
				"Plan",
				1,
				sdk.NewInt64Coin("umilk", 100_000000),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				types.NewBasicPoolsDistribution(1),
				types.NewBasicOperatorsDistribution(1),
				types.NewBasicUsersDistribution(1),
			),
			date:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expActive: true,
		},
		{
			name: "plan is active between start time and end time",
			plan: types.NewRewardsPlan(
				1,
				"Plan",
				1,
				sdk.NewInt64Coin("umilk", 100_000000),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				types.NewBasicPoolsDistribution(1),
				types.NewBasicOperatorsDistribution(1),
				types.NewBasicUsersDistribution(1),
			),
			date:      time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			expActive: true,
		},
		{
			name: "plan is inactive at end time",
			plan: types.NewRewardsPlan(
				1,
				"Plan",
				1,
				sdk.NewInt64Coin("umilk", 100_000000),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				types.NewBasicPoolsDistribution(1),
				types.NewBasicOperatorsDistribution(1),
				types.NewBasicUsersDistribution(1),
			),
			date:      time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expActive: false,
		},
		{
			name: "plan is inactive after end time",
			plan: types.NewRewardsPlan(
				1,
				"Plan",
				1,
				sdk.NewInt64Coin("umilk", 100_000000),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				types.NewBasicPoolsDistribution(1),
				types.NewBasicOperatorsDistribution(1),
				types.NewBasicUsersDistribution(1),
			),
			date:      time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
			expActive: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isActive := tc.plan.IsActiveAt(tc.date)
			require.Equal(t, tc.expActive, isActive)
		})
	}
}

func TestDistribution_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		distrType types.DistributionType
		shouldErr bool
	}{
		{
			name:      "basic distribution type returns no error",
			distrType: &types.DistributionTypeBasic{},
			shouldErr: false,
		},
		{
			name: "invalid delegation target ID returns error",
			distrType: &types.DistributionTypeWeighted{
				Weights: []types.DistributionWeight{
					types.NewDistributionWeight(0, 1),
				},
			},
			shouldErr: true,
		},
		{
			name: "invalid weight returns error",
			distrType: &types.DistributionTypeWeighted{
				Weights: []types.DistributionWeight{
					types.NewDistributionWeight(1, 0),
				},
			},
			shouldErr: true,
		},
		{
			name: "duplicated delegation target ID returns error",
			distrType: &types.DistributionTypeWeighted{
				Weights: []types.DistributionWeight{
					types.NewDistributionWeight(1, 1),
					types.NewDistributionWeight(1, 2),
				},
			},
			shouldErr: true,
		},
		{
			name:      "egalitarian distribution type returns no error",
			distrType: &types.DistributionTypeEgalitarian{},
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.distrType.Validate()
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUsersDistribution_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		distrType types.UsersDistributionType
		shouldErr bool
	}{
		{
			name:      "basic distribution type returns no error",
			distrType: &types.UsersDistributionTypeBasic{},
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.distrType.Validate()
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
