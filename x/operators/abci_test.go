package operators_test

import (
	"testing"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	db "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/app"
	"github.com/milkyway-labs/milkyway/x/operators"
	"github.com/milkyway-labs/milkyway/x/operators/keeper"
	"github.com/milkyway-labs/milkyway/x/operators/types"
)

func TestBeginBlocker(t *testing.T) {
	// Define store keys
	keys := storetypes.NewKVStoreKeys(types.StoreKey)

	// Create an in-memory db
	memDB := db.NewMemDB()
	ms := store.NewCommitMultiStore(memDB, log.NewNopLogger(), metrics.NewNoOpMetrics())
	for _, key := range keys {
		ms.MountStoreWithDB(key, storetypes.StoreTypeIAVL, memDB)
	}

	err := ms.LoadLatestVersion()
	require.NoError(t, err)

	ctx := sdk.NewContext(ms, tmproto.Header{ChainID: "test-chain"}, false, log.NewNopLogger())
	cdc, _ := app.MakeCodecs()

	operatorsKeeper := keeper.NewKeeper(cdc, keys[types.StoreKey],
		runtime.NewKVStoreService(keys[types.StoreKey]), nil, nil, "")

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
				operatorsKeeper.SetParams(ctx, types.NewParams(nil, 6*time.Hour))
				operatorsKeeper.StartOperatorInactivation(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					types.DefaultOperatorParams(),
				))
			},
			updateCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(time.Date(2024, 1, 1, 17, 59, 59, 999, time.UTC))
			},
			check: func(ctx sdk.Context) {
				kvStore := ctx.KVStore(keys[types.StoreKey])
				endTime := time.Date(2024, 1, 1, 18, 0, 0, 0, time.UTC)
				require.True(t, kvStore.Has(types.InactivatingOperatorQueueKey(1, endTime)))
			},
		},
		{
			name: "operator inactivation is completed at exact time",
			setupCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(time.Date(2020, 1, 1, 12, 0, 0, 1, time.UTC))
			},
			store: func(ctx sdk.Context) {
				operatorsKeeper.SetParams(ctx, types.NewParams(nil, 6*time.Hour))
				operatorsKeeper.StartOperatorInactivation(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					types.DefaultOperatorParams(),
				))
			},
			updateCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(time.Date(2020, 1, 1, 18, 0, 0, 0, time.UTC))
			},
			check: func(ctx sdk.Context) {
				// Make sure the operator is still inactivating
				operator, found := operatorsKeeper.GetOperator(ctx, 1)
				require.True(t, found)
				require.Equal(t, types.OPERATOR_STATUS_INACTIVATING, operator.Status)

				// Make sure the operator is still in the inactivating queue
				kvStore := ctx.KVStore(keys[types.StoreKey])
				endTime := time.Date(2020, 1, 1, 18, 0, 0, 0, time.UTC)
				require.False(t, kvStore.Has(types.InactivatingOperatorQueueKey(1, endTime)))
			},
		},
		{
			name: "operator inactivation is completed after time",
			setupCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(time.Date(2020, 1, 1, 12, 0, 0, 1, time.UTC))
			},
			store: func(ctx sdk.Context) {
				operatorsKeeper.SetParams(ctx, types.NewParams(nil, 6*time.Hour))
				operatorsKeeper.StartOperatorInactivation(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					types.DefaultOperatorParams(),
				))
			},
			updateCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(time.Date(2020, 1, 1, 20, 0, 0, 0, time.UTC))
			},
			check: func(ctx sdk.Context) {
				// Make sure the operator is inactive
				operator, found := operatorsKeeper.GetOperator(ctx, 1)
				require.True(t, found)
				require.Equal(t, types.OPERATOR_STATUS_INACTIVE, operator.Status)

				// Make sure the operator is not in the inactivating queue
				kvStore := ctx.KVStore(keys[types.StoreKey])
				endTime := time.Date(2020, 1, 1, 18, 0, 0, 0, time.UTC)
				require.False(t, kvStore.Has(types.InactivatingOperatorQueueKey(1, endTime)))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx, _ := ctx.CacheContext()
			if tc.setupCtx != nil {
				ctx = tc.setupCtx(ctx)
			}
			if tc.store != nil {
				tc.store(ctx)
			}

			if tc.updateCtx != nil {
				ctx = tc.updateCtx(ctx)
			}

			operators.BeginBlocker(ctx, operatorsKeeper)

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}
