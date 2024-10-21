package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
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

func TestOperatorJoinedServices_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		entry     types.OperatorJoinedServices
		shouldErr bool
	}{
		{
			name:      "empty operator joined services is valid",
			entry:     types.NewEmptyOperatorJoinedServices(),
			shouldErr: false,
		},
		{
			name:      "joined services with a service id equal to 0 is invalid",
			entry:     types.NewOperatorJoinedServices([]uint32{0}),
			shouldErr: true,
		},
		{
			name:      "joined services with a duplicated service id is invalid",
			entry:     types.NewOperatorJoinedServices([]uint32{1, 1}),
			shouldErr: true,
		},
		{
			name:      "joined services validates properly",
			entry:     types.NewOperatorJoinedServices([]uint32{1, 2}),
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

func TestOperatorJoinedServices_Add(t *testing.T) {
	testCases := []struct {
		name         string
		entry        types.OperatorJoinedServices
		newServiceID uint32
		shouldErr    bool
	}{
		{
			name:         "add 0 should fail",
			entry:        types.NewEmptyOperatorJoinedServices(),
			newServiceID: 0,
			shouldErr:    true,
		},
		{
			name:         "add already present service id should fail",
			entry:        types.NewOperatorJoinedServices([]uint32{1}),
			newServiceID: 1,
			shouldErr:    true,
		},
		{
			name:         "add correctly to empty",
			entry:        types.NewEmptyOperatorJoinedServices(),
			newServiceID: 1,
			shouldErr:    false,
		},
		{
			name:         "add correctly",
			entry:        types.NewOperatorJoinedServices([]uint32{1, 2}),
			newServiceID: 3,
			shouldErr:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.entry.Validate()
			require.NoError(t, err)

			err = tc.entry.Add(tc.newServiceID)
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
