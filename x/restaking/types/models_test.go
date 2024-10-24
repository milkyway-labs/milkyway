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
