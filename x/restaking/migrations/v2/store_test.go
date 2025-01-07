package v2_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	operatorstypes "github.com/milkyway-labs/milkyway/v7/x/operators/types"
	v2 "github.com/milkyway-labs/milkyway/v7/x/restaking/migrations/v2"
	"github.com/milkyway-labs/milkyway/v7/x/restaking/testutils"
	servicestypes "github.com/milkyway-labs/milkyway/v7/x/services/types"
)

func TestMigrateStore_removeNotAllowedJoinedServices(t *testing.T) {
	testData := testutils.NewKeeperTestData(t)

	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name: "operator joined service without allow list",
			store: func(ctx sdk.Context) {
				err := testData.OperatorsKeeper.SaveOperator(ctx, operatorstypes.NewOperator(
					2,
					operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				require.NoError(t, err)

				err = testData.ServicesKeeper.SaveService(ctx, servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				require.NoError(t, err)

				err = testData.Keeper.AddServiceToOperatorJoinedServices(ctx, 2, 1)
				require.NoError(t, err)
			},
			shouldErr: false,
			check: func(ctx sdk.Context) {
				hasJoined, err := testData.Keeper.HasOperatorJoinedService(ctx, 2, 1)
				require.NoError(t, err)
				require.True(t, hasJoined)
			},
		},
		{
			name: "operator joined service with allow list and is allowed",
			store: func(ctx sdk.Context) {
				err := testData.OperatorsKeeper.SaveOperator(ctx, operatorstypes.NewOperator(
					2,
					operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				require.NoError(t, err)

				err = testData.ServicesKeeper.SaveService(ctx, servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				require.NoError(t, err)

				err = testData.Keeper.AddOperatorToServiceAllowList(ctx, 1, 2)
				require.NoError(t, err)

				err = testData.Keeper.AddServiceToOperatorJoinedServices(ctx, 2, 1)
				require.NoError(t, err)
			},
			shouldErr: false,
			check: func(ctx sdk.Context) {
				hasJoined, err := testData.Keeper.HasOperatorJoinedService(ctx, 2, 1)
				require.NoError(t, err)
				require.True(t, hasJoined)
			},
		},
		{
			name: "operator joined service with allow list and should not be allowed",
			store: func(ctx sdk.Context) {
				err := testData.OperatorsKeeper.SaveOperator(ctx, operatorstypes.NewOperator(
					2,
					operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				require.NoError(t, err)

				err = testData.ServicesKeeper.SaveService(ctx, servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				require.NoError(t, err)

				err = testData.Keeper.AddOperatorToServiceAllowList(ctx, 1, 1)
				require.NoError(t, err)

				err = testData.Keeper.AddServiceToOperatorJoinedServices(ctx, 2, 1)
				require.NoError(t, err)
			},
			shouldErr: false,
			check: func(ctx sdk.Context) {
				hasJoined, err := testData.Keeper.HasOperatorJoinedService(ctx, 2, 1)
				require.NoError(t, err)
				require.False(t, hasJoined)
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

			err := v2.MigrateStore(ctx, testData.Keeper)
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}
