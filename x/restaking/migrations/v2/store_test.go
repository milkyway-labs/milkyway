package v2_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	operatorstypes "github.com/milkyway-labs/milkyway/v9/x/operators/types"
	v2 "github.com/milkyway-labs/milkyway/v9/x/restaking/migrations/v2"
	"github.com/milkyway-labs/milkyway/v9/x/restaking/testutils"
	"github.com/milkyway-labs/milkyway/v9/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v9/x/services/types"
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

func TestMigrateStore_migratePreferences(t *testing.T) {
	testData := testutils.NewKeeperTestData(t)

	testCases := []struct {
		name           string
		oldPreferences []v2.UserPreferencesEntry
		shouldErr      bool
		expPreferences []types.UserPreferencesEntry
	}{
		{
			name: "preferences are migrated properly",
			oldPreferences: []v2.UserPreferencesEntry{
				{
					UserAddress: "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
					Preferences: v2.UserPreferences{
						TrustAccreditedServices:    false,
						TrustNonAccreditedServices: true,
						TrustedServicesIDs:         []uint32{1, 2, 3},
					},
				},
				{
					UserAddress: "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					Preferences: v2.UserPreferences{
						TrustAccreditedServices:    false,
						TrustNonAccreditedServices: true,
						TrustedServicesIDs:         []uint32{4, 5},
					},
				},
			},
			shouldErr: false,
			expPreferences: []types.UserPreferencesEntry{
				types.NewUserPreferencesEntry("cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn", types.NewUserPreferences([]types.TrustedServiceEntry{
					types.NewTrustedServiceEntry(1, nil),
					types.NewTrustedServiceEntry(2, nil),
					types.NewTrustedServiceEntry(3, nil),
				})),
				types.NewUserPreferencesEntry("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", types.NewUserPreferences([]types.TrustedServiceEntry{
					types.NewTrustedServiceEntry(4, nil),
					types.NewTrustedServiceEntry(5, nil),
				})),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx, _ := testData.Context.CacheContext()

			// Store the delegations
			store := testData.StoreService.OpenKVStore(ctx)
			for _, entry := range tc.oldPreferences {
				key := append(types.UserPreferencesPrefix, []byte(entry.UserAddress)...)
				err := store.Set(key, testData.Cdc.MustMarshal(&entry.Preferences))
				require.NoError(t, err)
			}

			err := v2.MigrateStore(ctx, testData.Keeper, testData.StoreService, testData.Cdc)
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			// Check the new preferences
			for _, entry := range tc.expPreferences {
				stored, err := testData.Keeper.GetUserPreferences(ctx, entry.UserAddress)
				require.NoError(t, err)
				require.Equal(t, entry.Preferences, stored)
			}
		})
	}
}
