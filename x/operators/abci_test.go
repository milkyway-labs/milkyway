package operators_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/v2/x/operators"
	"github.com/milkyway-labs/milkyway/v2/x/operators/testutils"
	"github.com/milkyway-labs/milkyway/v2/x/operators/types"
)

func TestBeginBlocker(t *testing.T) {
	data := testutils.NewKeeperTestData(t)
	operatorsKeeper := data.Keeper

	testCases := []struct {
		name      string
		setupCtx  func(ctx sdk.Context) sdk.Context
		store     func(ctx sdk.Context)
		updateCtx func(ctx sdk.Context) sdk.Context
		check     func(ctx sdk.Context)
	}{
		{
			name: "operator inactivation is not completed before time",
			setupCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
			},
			store: func(ctx sdk.Context) {
				err := operatorsKeeper.SetParams(ctx, types.NewParams(nil, 6*time.Hour))
				require.NoError(t, err)

				err = operatorsKeeper.StartOperatorInactivation(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				require.NoError(t, err)
			},
			updateCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(time.Date(2024, 1, 1, 17, 59, 59, 999, time.UTC))
			},
			check: func(ctx sdk.Context) {
				kvStore := data.StoreService.OpenKVStore(ctx)
				endTime := time.Date(2024, 1, 1, 18, 0, 0, 0, time.UTC)

				hasKey, err := kvStore.Has(types.InactivatingOperatorQueueKey(1, endTime))
				require.NoError(t, err)
				require.True(t, hasKey)
			},
		},
		{
			name: "operator inactivation is completed at exact time",
			setupCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(time.Date(2020, 1, 1, 12, 0, 0, 1, time.UTC))
			},
			store: func(ctx sdk.Context) {
				err := operatorsKeeper.SetParams(ctx, types.NewParams(nil, 6*time.Hour))
				require.NoError(t, err)

				err = operatorsKeeper.StartOperatorInactivation(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				require.NoError(t, err)
			},
			updateCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(time.Date(2020, 1, 1, 18, 0, 0, 0, time.UTC))
			},
			check: func(ctx sdk.Context) {
				// Make sure the operator is still inactivating
				operator, found, err := operatorsKeeper.GetOperator(ctx, 1)
				require.NoError(t, err)
				require.True(t, found)
				require.Equal(t, types.OPERATOR_STATUS_INACTIVATING, operator.Status)

				// Make sure the operator is still in the inactivating queue
				kvStore := data.StoreService.OpenKVStore(ctx)
				endTime := time.Date(2020, 1, 1, 18, 0, 0, 0, time.UTC)

				hasKey, err := kvStore.Has(types.InactivatingOperatorQueueKey(1, endTime))
				require.NoError(t, err)
				require.False(t, hasKey)
			},
		},
		{
			name: "operator inactivation is completed after time",
			setupCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(time.Date(2020, 1, 1, 12, 0, 0, 1, time.UTC))
			},
			store: func(ctx sdk.Context) {
				err := operatorsKeeper.SetParams(ctx, types.NewParams(nil, 6*time.Hour))
				require.NoError(t, err)

				err = operatorsKeeper.StartOperatorInactivation(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				require.NoError(t, err)
			},
			updateCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(time.Date(2020, 1, 1, 20, 0, 0, 0, time.UTC))
			},
			check: func(ctx sdk.Context) {
				// Make sure the operator is inactive
				operator, found, err := operatorsKeeper.GetOperator(ctx, 1)
				require.NoError(t, err)
				require.True(t, found)
				require.Equal(t, types.OPERATOR_STATUS_INACTIVE, operator.Status)

				// Make sure the operator is not in the inactivating queue
				kvStore := data.StoreService.OpenKVStore(ctx)
				endTime := time.Date(2020, 1, 1, 18, 0, 0, 0, time.UTC)

				hasKey, err := kvStore.Has(types.InactivatingOperatorQueueKey(1, endTime))
				require.NoError(t, err)
				require.False(t, hasKey)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx, _ := data.Context.CacheContext()
			if tc.setupCtx != nil {
				ctx = tc.setupCtx(ctx)
			}
			if tc.store != nil {
				tc.store(ctx)
			}

			if tc.updateCtx != nil {
				ctx = tc.updateCtx(ctx)
			}

			err := operators.BeginBlocker(ctx, operatorsKeeper)
			require.NoError(t, err)

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}
