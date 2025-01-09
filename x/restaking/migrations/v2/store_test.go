package v2_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	operatorstypes "github.com/milkyway-labs/milkyway/v7/x/operators/types"
	v2 "github.com/milkyway-labs/milkyway/v7/x/restaking/migrations/v2"
	"github.com/milkyway-labs/milkyway/v7/x/restaking/testutils"
	"github.com/milkyway-labs/milkyway/v7/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v7/x/services/types"
)

func TestMigrateStore_setDelegationByTargetIDValues(t *testing.T) {
	testData := testutils.NewKeeperTestData(t)

	testCases := []struct {
		name        string
		delegations []types.Delegation
		shouldErr   bool
	}{
		{
			name: "delegations with all prefixes",
			delegations: []types.Delegation{
				types.NewPoolDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(50))),
				),
				types.NewOperatorDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
				),
				types.NewServiceDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
				),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx, _ := testData.Context.CacheContext()

			// Store the delegations
			store := testData.StoreService.OpenKVStore(ctx)
			for _, delegation := range tc.delegations {
				getDelegationKey, _, err := types.GetDelegationKeyBuilders(delegation)
				require.NoError(t, err)

				err = store.Set(getDelegationKey(delegation.UserAddress, delegation.TargetID), testData.Cdc.MustMarshal(&delegation))
				require.NoError(t, err)
			}

			// Make sure that the delegation-by-target-id keys do not exist
			for _, delegation := range tc.delegations {
				_, getDelegationByTargetIDKey, err := types.GetDelegationKeyBuilders(delegation)
				require.NoError(t, err)

				stored, err := store.Get(getDelegationByTargetIDKey(delegation.TargetID, delegation.UserAddress))
				require.NoError(t, err)
				require.Empty(t, stored)
			}

			// Run the migrations and check the results
			err := v2.MigrateStore(ctx, testData.Keeper, testData.StoreService, testData.Cdc)
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Check the delegations-by-target-id keys
				for _, delegation := range tc.delegations {
					_, getDelegationByTargetIDKey, err := types.GetDelegationKeyBuilders(delegation)
					require.NoError(t, err)

					stored, err := store.Get(getDelegationByTargetIDKey(delegation.TargetID, delegation.UserAddress))
					require.NoError(t, err)
					require.NotEmpty(t, stored)
				}
			}
		})
	}
}

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

			err := v2.MigrateStore(ctx, testData.Keeper, testData.StoreService, testData.Cdc)
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
