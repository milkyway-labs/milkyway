package v2_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/v9/utils"
	v2 "github.com/milkyway-labs/milkyway/v9/x/rewards/migrations/v2"
	"github.com/milkyway-labs/milkyway/v9/x/rewards/testutils"
	"github.com/milkyway-labs/milkyway/v9/x/rewards/types"
)

func Test_MigratePlan(t *testing.T) {
	testData := testutils.NewKeeperTestData(t)

	testCases := []struct {
		name        string
		store       func(ctx sdk.Context)
		legacyPlans []v2.RewardsPlan
		shouldErr   bool
		expPlans    []types.RewardsPlan
		check       func(ctx sdk.Context)
	}{
		{
			name: "single plan with multiple denoms",
			store: func(ctx sdk.Context) {
				// Store the next ID
				store := testData.StoreService.OpenKVStore(ctx)
				err := store.Set(types.NextRewardsPlanIDKey, sdk.Uint64ToBigEndian(3))
				require.NoError(t, err)
			},
			legacyPlans: []v2.RewardsPlan{
				{
					ID:                    2,
					Description:           "Plan 2",
					ServiceID:             2,
					AmountPerDay:          utils.MustParseCoins("100_000000service1,100_000000service2"),
					StartTime:             time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
					EndTime:               time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
					PoolsDistribution:     types.NewEgalitarianPoolsDistribution(1),
					OperatorsDistribution: types.NewEgalitarianOperatorsDistribution(1),
					UsersDistribution:     types.NewBasicUsersDistribution(1),
				},
			},
			shouldErr: false,
			expPlans: []types.RewardsPlan{
				types.NewRewardsPlan(
					2,
					"Plan 2",
					2,
					utils.MustParseCoin("100_000000service1"),
					time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
					types.NewEgalitarianPoolsDistribution(1),
					types.NewEgalitarianOperatorsDistribution(1),
					types.NewBasicUsersDistribution(1),
				),
				types.NewRewardsPlan(
					3,
					"Plan 2",
					2,
					utils.MustParseCoin("100_000000service2"),
					time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
					types.NewEgalitarianPoolsDistribution(1),
					types.NewEgalitarianOperatorsDistribution(1),
					types.NewBasicUsersDistribution(1),
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the next plan id is now 4
				store := testData.StoreService.OpenKVStore(ctx)
				bz, err := store.Get(types.NextRewardsPlanIDKey)
				require.NoError(t, err)
				require.EqualValues(t, 4, sdk.BigEndianToUint64(bz))
			},
		},
		{
			name: "plan with multiple denoms as the first plan",
			store: func(ctx sdk.Context) {
				// Store the next ID
				store := testData.StoreService.OpenKVStore(ctx)
				err := store.Set(types.NextRewardsPlanIDKey, sdk.Uint64ToBigEndian(11))
				require.NoError(t, err)
			},
			legacyPlans: []v2.RewardsPlan{
				{
					ID:                    2,
					Description:           "Plan 2",
					ServiceID:             2,
					AmountPerDay:          utils.MustParseCoins("100_000000service1,100_000000service2"),
					StartTime:             time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
					EndTime:               time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
					PoolsDistribution:     types.NewEgalitarianPoolsDistribution(1),
					OperatorsDistribution: types.NewEgalitarianOperatorsDistribution(1),
					UsersDistribution:     types.NewBasicUsersDistribution(1),
				},
				{
					ID:                    10,
					Description:           "Plan 10",
					ServiceID:             10,
					AmountPerDay:          utils.MustParseCoins("100_000000service10"),
					StartTime:             time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
					EndTime:               time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
					PoolsDistribution:     types.NewEgalitarianPoolsDistribution(1),
					OperatorsDistribution: types.NewEgalitarianOperatorsDistribution(1),
					UsersDistribution:     types.NewBasicUsersDistribution(1),
				},
			},
			shouldErr: false,
			expPlans: []types.RewardsPlan{
				types.NewRewardsPlan(
					2,
					"Plan 2",
					2,
					utils.MustParseCoin("100_000000service1"),
					time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
					types.NewEgalitarianPoolsDistribution(1),
					types.NewEgalitarianOperatorsDistribution(1),
					types.NewBasicUsersDistribution(1),
				),
				types.NewRewardsPlan(
					10,
					"Plan 10",
					10,
					utils.MustParseCoin("100_000000service10"),
					time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
					types.NewEgalitarianPoolsDistribution(1),
					types.NewEgalitarianOperatorsDistribution(1),
					types.NewBasicUsersDistribution(1),
				),
				types.NewRewardsPlan(
					11,
					"Plan 2",
					2,
					utils.MustParseCoin("100_000000service2"),
					time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
					types.NewEgalitarianPoolsDistribution(1),
					types.NewEgalitarianOperatorsDistribution(1),
					types.NewBasicUsersDistribution(1),
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the next plan id is now 4
				store := testData.StoreService.OpenKVStore(ctx)
				bz, err := store.Get(types.NextRewardsPlanIDKey)
				require.NoError(t, err)
				require.EqualValues(t, 12, sdk.BigEndianToUint64(bz))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx, _ := testData.Context.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			// Store the old plans
			for _, plan := range tc.legacyPlans {
				store := testData.StoreService.OpenKVStore(ctx)
				err := store.Set(v2.PlanStoreKey(plan.ID), testData.Cdc.MustMarshal(&plan))
				require.NoError(t, err)
			}

			// Migrate the plans
			err := v2.MigrateStore(ctx, testData.StoreService, testData.Cdc)
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			// Check the stored plans
			var storedPlans []types.RewardsPlan
			err = testData.Keeper.RewardsPlans.Walk(ctx, nil, func(key uint64, value types.RewardsPlan) (stop bool, err error) {
				storedPlans = append(storedPlans, value)
				return false, nil
			})
			require.NoError(t, err)
		})
	}

}
