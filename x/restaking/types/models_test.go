package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/v7/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v7/x/services/types"
)

func TestPoolDelegation_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		entry     types.Delegation
		shouldErr bool
	}{
		{
			name: "invalid pool id returns error",
			entry: types.NewPoolDelegation(
				0,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
			),
			shouldErr: true,
		},
		{
			name: "invalid user address returns error",
			entry: types.NewPoolDelegation(
				1,
				"",
				sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
			),
			shouldErr: true,
		},
		{
			name: "invalid shares returns error",
			entry: types.NewPoolDelegation(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				sdk.DecCoins{sdk.DecCoin{Denom: "umilk", Amount: sdkmath.LegacyNewDec(-100)}},
			),
			shouldErr: true,
		},
		{
			name: "valid entry returns no error",
			entry: types.NewPoolDelegation(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
			),
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.entry.Validate()
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// --------------------------------------------------------------------------------------------------------------------

func TestServiceDelegation_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		entry     types.Delegation
		shouldErr bool
	}{
		{
			name: "invalid service id returns error",
			entry: types.NewServiceDelegation(
				0,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
			),
			shouldErr: true,
		},
		{
			name: "invalid user address returns error",
			entry: types.NewServiceDelegation(
				1,
				"",
				sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
			),
			shouldErr: true,
		},
		{
			name: "invalid shares returns error",
			entry: types.NewServiceDelegation(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				sdk.DecCoins{sdk.DecCoin{Denom: "umilk", Amount: sdkmath.LegacyNewDec(-100)}},
			),
			shouldErr: true,
		},
		{
			name: "valid entry returns no error",
			entry: types.NewServiceDelegation(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
			),
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.entry.Validate()
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// --------------------------------------------------------------------------------------------------------------------

func TestOperatorDelegation_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		entry     types.Delegation
		shouldErr bool
	}{
		{
			name: "invalid operator id returns error",
			entry: types.NewOperatorDelegation(
				0,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
			),
			shouldErr: true,
		},
		{
			name: "invalid user address returns error",
			entry: types.NewOperatorDelegation(
				1,
				"",
				sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
			),
			shouldErr: true,
		},
		{
			name: "invalid shares returns error",
			entry: types.NewOperatorDelegation(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				sdk.DecCoins{sdk.DecCoin{Denom: "umilk", Amount: sdkmath.LegacyNewDec(-100)}},
			),
			shouldErr: true,
		},
		{
			name: "valid entry returns no error",
			entry: types.NewOperatorDelegation(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
			),
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.entry.Validate()
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// --------------------------------------------------------------------------------------------------------------------

func TestUserPreferences_Validate(t *testing.T) {
	testCases := []struct {
		name        string
		preferences types.UserPreferences
		shouldErr   bool
	}{
		{
			name: "invalid service id returns error",
			preferences: types.NewUserPreferences([]types.TrustedServiceEntry{
				types.NewTrustedServiceEntry(0, nil),
			}),
			shouldErr: true,
		},
		{
			name: "invalid pool id returns error",
			preferences: types.NewUserPreferences([]types.TrustedServiceEntry{
				types.NewTrustedServiceEntry(1, []uint32{0}),
			}),
			shouldErr: true,
		},
		{
			name: "duplicated entry for the same service id returns error",
			preferences: types.NewUserPreferences([]types.TrustedServiceEntry{
				types.NewTrustedServiceEntry(1, nil),
				types.NewTrustedServiceEntry(1, nil),
			}),
			shouldErr: true,
		},
		{
			name: "valid preferences returns no error",
			preferences: types.NewUserPreferences([]types.TrustedServiceEntry{
				types.NewTrustedServiceEntry(1, []uint32{1, 2}),
			}),
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.preferences.Validate()
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUserPreferences_TrustedServicesIDs(t *testing.T) {
	testCases := []struct {
		name        string
		preferences types.UserPreferences
		service     servicestypes.Service
		expTrusted  bool
	}{
		{
			name: "user trusts only specified services - specified and accredited service",
			preferences: types.NewUserPreferences([]types.TrustedServiceEntry{
				types.NewTrustedServiceEntry(1, nil),
			}),
			service: servicestypes.NewService(
				1,
				servicestypes.SERVICE_STATUS_ACTIVE,
				"MilkyWay",
				"MilkyWay is an AVS of a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				true,
			),
			expTrusted: true,
		},
		{
			name: "user trusts only specified services - specified and non-accredited service",
			preferences: types.NewUserPreferences([]types.TrustedServiceEntry{
				types.NewTrustedServiceEntry(1, nil),
			}),
			service: servicestypes.NewService(
				1,
				servicestypes.SERVICE_STATUS_ACTIVE,
				"MilkyWay",
				"MilkyWay is an AVS of a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				false,
			),
			expTrusted: true,
		},
		{
			name: "user trusts only specified services - not specified and accredited service",
			preferences: types.NewUserPreferences([]types.TrustedServiceEntry{
				types.NewTrustedServiceEntry(1, nil),
			}),
			service: servicestypes.NewService(
				2,
				servicestypes.SERVICE_STATUS_ACTIVE,
				"MilkyWay",
				"MilkyWay is an AVS of a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				true,
			),
			expTrusted: false,
		},
		{
			name: "user trusts only specified services - not specified and non-accredited service",
			preferences: types.NewUserPreferences([]types.TrustedServiceEntry{
				types.NewTrustedServiceEntry(1, nil),
			}), service: servicestypes.NewService(
				2,
				servicestypes.SERVICE_STATUS_ACTIVE,
				"MilkyWay",
				"MilkyWay is an AVS of a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				false,
			),
			expTrusted: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			trusted := tc.preferences.TrustedServicesIDs()
			if tc.expTrusted {
				require.Contains(t, trusted, tc.service.ID)
			} else {
				require.NotContains(t, trusted, tc.service.ID)
			}
		})
	}
}

// --------------------------------------------------------------------------------------------------------------------

func TestComputeChangedServices(t *testing.T) {
	testCases := []struct {
		name                string
		before              types.UserPreferences
		after               types.UserPreferences
		expectedServicesIDs []uint32
	}{
		{
			name: "no changes",
			before: types.NewUserPreferences([]types.TrustedServiceEntry{
				types.NewTrustedServiceEntry(1, []uint32{1, 2}),
			}),
			after: types.NewUserPreferences([]types.TrustedServiceEntry{
				types.NewTrustedServiceEntry(1, []uint32{1, 2}),
			}),
			expectedServicesIDs: nil,
		},
		{
			name: "service removed",
			before: types.NewUserPreferences([]types.TrustedServiceEntry{
				types.NewTrustedServiceEntry(1, []uint32{1, 2}),
			}),
			after:               types.NewUserPreferences(nil),
			expectedServicesIDs: []uint32{1},
		},
		{
			name:   "service added",
			before: types.NewUserPreferences(nil),
			after: types.NewUserPreferences([]types.TrustedServiceEntry{
				types.NewTrustedServiceEntry(1, []uint32{1, 2}),
			}),
			expectedServicesIDs: []uint32{1},
		},
		{
			name: "pool within service removed",
			before: types.NewUserPreferences([]types.TrustedServiceEntry{
				types.NewTrustedServiceEntry(1, []uint32{1, 2}),
			}),
			after: types.NewUserPreferences([]types.TrustedServiceEntry{
				types.NewTrustedServiceEntry(1, []uint32{1}),
			}),
			expectedServicesIDs: []uint32{1},
		},
		{
			name: "pool within service added",
			before: types.NewUserPreferences([]types.TrustedServiceEntry{
				types.NewTrustedServiceEntry(1, []uint32{1}),
			}),
			after: types.NewUserPreferences([]types.TrustedServiceEntry{
				types.NewTrustedServiceEntry(1, []uint32{1, 2}),
			}),
			expectedServicesIDs: []uint32{1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entries := types.ComputeChangedServicesIDs(tc.before, tc.after)
			require.Equal(t, tc.expectedServicesIDs, entries)
		})
	}
}
