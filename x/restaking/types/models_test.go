package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
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
			preferences: types.UserPreferences{
				TrustedServicesIDs: []uint32{0},
			},
			shouldErr: true,
		},
		{
			name: "valid preferences returns no error",
			preferences: types.NewUserPreferences(
				false,
				true,
				[]uint32{1, 2, 6, 7},
			),
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

func TestUserPreferences_IsServiceTrusted(t *testing.T) {
	testCases := []struct {
		name        string
		preferences types.UserPreferences
		service     servicestypes.Service
		trusted     bool
	}{
		{
			name:        "user does not trust any services - accredited service",
			preferences: types.NewUserPreferences(false, false, nil),
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
			trusted: false,
		},
		{
			name:        "user does not trust any services - non-accredited service",
			preferences: types.NewUserPreferences(false, false, nil),
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
			trusted: false,
		},
		{
			name:        "user only trusts accredited services - accredited service",
			preferences: types.NewUserPreferences(false, true, nil),
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
			trusted: true,
		},
		{
			name:        "user only trusts accredited services - non-accredited service",
			preferences: types.NewUserPreferences(false, true, nil),
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
			trusted: false,
		},
		{
			name:        "user only trusts non-accredited services - accredited service",
			preferences: types.NewUserPreferences(true, false, nil),
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
			trusted: false,
		},
		{
			name:        "user only trusts non-accredited services - non-accredited service",
			preferences: types.NewUserPreferences(true, false, nil),
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
			trusted: true,
		},
		{
			name:        "user trusts both accredited and non-accredited services - accredited service",
			preferences: types.NewUserPreferences(true, true, nil),
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
			trusted: true,
		},
		{
			name:        "user trusts both accredited and non-accredited services - non-accredited service",
			preferences: types.NewUserPreferences(true, true, nil),
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
			trusted: true,
		},
		{
			name:        "user trusts only specified services - specified and accredited service",
			preferences: types.NewUserPreferences(false, false, []uint32{1}),
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
			trusted: true,
		},
		{
			name:        "user trusts only specified services - specified and non-accredited service",
			preferences: types.NewUserPreferences(false, false, []uint32{1}),
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
			trusted: true,
		},
		{
			name:        "user trusts only specified services - not specified and accredited service",
			preferences: types.NewUserPreferences(false, false, []uint32{1}),
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
			trusted: false,
		},
		{
			name:        "user trusts only specified services - not specified and non-accredited service",
			preferences: types.NewUserPreferences(false, false, []uint32{1}),
			service: servicestypes.NewService(
				2,
				servicestypes.SERVICE_STATUS_ACTIVE,
				"MilkyWay",
				"MilkyWay is an AVS of a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				false,
			),
			trusted: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			trusted := tc.preferences.IsServiceTrusted(tc.service)
			require.Equal(t, tc.trusted, trusted)
		})
	}
}
