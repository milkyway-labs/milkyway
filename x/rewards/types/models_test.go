package types_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/x/rewards/types"
)

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
