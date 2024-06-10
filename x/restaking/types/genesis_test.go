package types_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

func TestGenesis_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		genesis   *types.GenesisState
		shouldErr bool
	}{
		{
			name: "invalid pool delegation entry returns error",
			genesis: types.NewGenesis(
				[]types.PoolDelegationEntry{
					types.NewPoolDelegationEntry(
						0,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewCoin("umilk", sdkmath.NewInt(100)),
					),
				},
				nil,
				nil,
				types.DefaultParams(),
			),
			shouldErr: true,
		},
		{
			name: "invalid service delegation entry returns error",
			genesis: types.NewGenesis(
				nil,
				[]types.ServiceDelegationEntry{
					types.NewServiceDelegationEntry(
						0,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewCoin("umilk", sdkmath.NewInt(100)),
					),
				},
				nil,
				types.DefaultParams(),
			),
			shouldErr: true,
		},
		{
			name: "invalid operator delegation entry returns error",
			genesis: types.NewGenesis(
				nil,
				nil,
				[]types.OperatorDelegationEntry{
					types.NewOperatorDelegationEntry(
						0,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewCoin("umilk", sdkmath.NewInt(100)),
					),
				},
				types.DefaultParams(),
			),
			shouldErr: true,
		},
		{
			name: "invalid params return error",
			genesis: types.NewGenesis(
				nil,
				nil,
				nil,
				types.NewParams(0),
			),
			shouldErr: true,
		},
		{
			name:      "default genesis returns no error",
			genesis:   types.DefaultGenesis(),
			shouldErr: false,
		},
		{
			name: "valid genesis returns no error",
			genesis: types.NewGenesis(
				[]types.PoolDelegationEntry{
					types.NewPoolDelegationEntry(
						1,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewCoin("umilk", sdkmath.NewInt(100)),
					),
				},
				[]types.ServiceDelegationEntry{
					types.NewServiceDelegationEntry(
						2,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewCoin("umilk", sdkmath.NewInt(100)),
					),
				},
				[]types.OperatorDelegationEntry{
					types.NewOperatorDelegationEntry(
						3,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewCoin("umilk", sdkmath.NewInt(100)),
					),
				},
				types.NewParams(5*24*time.Hour),
			),
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.genesis.Validate()
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// --------------------------------------------------------------------------------------------------------------------

func TestPoolDelegationEntry_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		entry     types.PoolDelegationEntry
		shouldErr bool
	}{
		{
			name: "invalid pool id returns error",
			entry: types.NewPoolDelegationEntry(
				0,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				sdk.NewCoin("umilk", sdkmath.NewInt(100)),
			),
			shouldErr: true,
		},
		{
			name: "invalid user address returns error",
			entry: types.NewPoolDelegationEntry(
				1,
				"",
				sdk.NewCoin("umilk", sdkmath.NewInt(100)),
			),
			shouldErr: true,
		},
		{
			name: "invalid amount returns error",
			entry: types.NewPoolDelegationEntry(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				sdk.Coin{},
			),
			shouldErr: true,
		},
		{
			name: "valid entry returns no error",
			entry: types.NewPoolDelegationEntry(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				sdk.NewCoin("umilk", sdkmath.NewInt(100)),
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

func TestServiceDelegationEntry_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		entry     types.ServiceDelegationEntry
		shouldErr bool
	}{
		{
			name: "invalid service id returns error",
			entry: types.NewServiceDelegationEntry(
				0,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				sdk.NewCoin("umilk", sdkmath.NewInt(100)),
			),
			shouldErr: true,
		},
		{
			name: "invalid user address returns error",
			entry: types.NewServiceDelegationEntry(
				1,
				"",
				sdk.NewCoin("umilk", sdkmath.NewInt(100)),
			),
			shouldErr: true,
		},
		{
			name: "invalid amount returns error",
			entry: types.NewServiceDelegationEntry(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				sdk.Coin{},
			),
			shouldErr: true,
		},
		{
			name: "valid entry returns no error",
			entry: types.NewServiceDelegationEntry(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				sdk.NewCoin("umilk", sdkmath.NewInt(100)),
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

func TestOperatorDelegationEntry_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		entry     types.OperatorDelegationEntry
		shouldErr bool
	}{
		{
			name: "invalid operator id returns error",
			entry: types.NewOperatorDelegationEntry(
				0,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				sdk.NewCoin("umilk", sdkmath.NewInt(100)),
			),
			shouldErr: true,
		},
		{
			name: "invalid user address returns error",
			entry: types.NewOperatorDelegationEntry(
				1,
				"",
				sdk.NewCoin("umilk", sdkmath.NewInt(100)),
			),
			shouldErr: true,
		},
		{
			name: "invalid amount returns error",
			entry: types.NewOperatorDelegationEntry(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				sdk.Coin{},
			),
			shouldErr: true,
		},
		{
			name: "valid entry returns no error",
			entry: types.NewOperatorDelegationEntry(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				sdk.NewCoin("umilk", sdkmath.NewInt(100)),
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
