package types_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/v7/x/restaking/types"
)

func TestGenesis_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		genesis   *types.GenesisState
		shouldErr bool
	}{
		{
			name: "operator joined services record with invalid JoinedService returns error",
			genesis: types.NewGenesis(
				[]types.OperatorJoinedServices{
					types.NewOperatorJoinedServices(1, []uint32{0}),
				},
				nil,
				nil,
				nil,
				nil,
				nil,
				types.DefaultParams(),
			),
			shouldErr: true,
		},
		{
			name: "operator joined services record with invalid OperatorID returns error",
			genesis: types.NewGenesis(
				[]types.OperatorJoinedServices{
					types.NewOperatorJoinedServices(0, []uint32{1}),
				},
				nil,
				nil,
				nil,
				nil,
				nil,
				types.DefaultParams(),
			),
			shouldErr: true,
		},
		{
			name: "duplicated service allowed operator returns error",
			genesis: types.NewGenesis(
				nil,
				[]types.ServiceAllowedOperators{
					types.NewServiceAllowedOperators(1, []uint32{1, 2, 3}),
					types.NewServiceAllowedOperators(1, []uint32{1, 2, 3}),
				},
				nil,
				nil,
				nil,
				nil,
				types.DefaultParams(),
			),
			shouldErr: true,
		},
		{
			name: "service allowed operators with invalid service ID returns error",
			genesis: types.NewGenesis(
				nil,
				[]types.ServiceAllowedOperators{
					types.NewServiceAllowedOperators(0, []uint32{1, 2, 3}),
				},
				nil,
				nil,
				nil,
				nil,
				types.DefaultParams(),
			),
			shouldErr: true,
		},
		{
			name: "service allowed operators with invalid list returns error",
			genesis: types.NewGenesis(
				nil,
				[]types.ServiceAllowedOperators{
					types.NewServiceAllowedOperators(1, []uint32{0, 1}),
				},
				nil,
				nil,
				nil,
				nil,
				types.DefaultParams(),
			),
			shouldErr: true,
		},
		{
			name: "duplicated service securing pool returns error",
			genesis: types.NewGenesis(
				nil,
				nil,
				[]types.ServiceSecuringPools{
					types.NewServiceSecuringPools(1, []uint32{1, 2, 3}),
					types.NewServiceSecuringPools(1, []uint32{1, 2, 3}),
				},
				nil,
				nil,
				nil,
				types.DefaultParams(),
			),
			shouldErr: true,
		},
		{
			name: "service securing pools with invalid service ID returns error",
			genesis: types.NewGenesis(
				nil,
				nil,
				[]types.ServiceSecuringPools{
					types.NewServiceSecuringPools(0, []uint32{1, 2, 3}),
				},
				nil,
				nil,
				nil,
				types.DefaultParams(),
			),
			shouldErr: true,
		},
		{
			name: "service securing pools with invalid list returns error",
			genesis: types.NewGenesis(
				nil,
				nil,
				[]types.ServiceSecuringPools{
					types.NewServiceSecuringPools(1, []uint32{0, 1}),
				},
				nil,
				nil,
				nil,
				types.DefaultParams(),
			),
			shouldErr: true,
		},
		{
			name: "invalid pool delegation entry returns error",
			genesis: types.NewGenesis(
				nil,
				nil,
				nil,
				[]types.Delegation{
					types.NewPoolDelegation(
						0,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
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
				nil,
				nil,
				[]types.Delegation{
					types.NewServiceDelegation(
						0,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
					),
				},
				nil,
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
				nil,
				[]types.Delegation{
					types.NewOperatorDelegation(
						0,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
					),
				},
				nil,
				nil,
				types.DefaultParams(),
			),
			shouldErr: true,
		},
		{
			name: "invalid unbonding delegation returns error",
			genesis: types.NewGenesis(
				nil,
				nil,
				nil,
				nil,
				[]types.UnbondingDelegation{
					types.NewPoolUnbondingDelegation(
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						0,
						1,
						time.Now(),
						sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
						1,
					),
				},
				nil,
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
				nil,
				nil,
				nil,
				types.NewParams(0, nil, types.DefaultRestakingCap, types.DefaultMaxEntries),
			),
			shouldErr: true,
		},
		{
			name: "invalid user preferences entry returns error",
			genesis: types.NewGenesis(
				nil,
				nil,
				nil,
				nil,
				nil,
				[]types.UserPreferencesEntry{
					types.NewUserPreferencesEntry(
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						types.NewUserPreferences([]types.TrustedServiceEntry{
							types.NewTrustedServiceEntry(0, []uint32{1}),
						}),
					),
				},
				types.DefaultParams(),
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
				[]types.OperatorJoinedServices{
					types.NewOperatorJoinedServices(1, []uint32{2, 3, 5}),
				},
				[]types.ServiceAllowedOperators{
					types.NewServiceAllowedOperators(1, []uint32{1, 2, 3}),
					types.NewServiceAllowedOperators(2, []uint32{5, 6, 7}),
				},
				[]types.ServiceSecuringPools{
					types.NewServiceSecuringPools(3, []uint32{1, 2, 3}),
					types.NewServiceSecuringPools(4, []uint32{5, 6, 7}),
				},
				[]types.Delegation{
					types.NewPoolDelegation(
						1,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
					),
					types.NewServiceDelegation(
						2,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
					),
					types.NewOperatorDelegation(
						3,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
					),
				},
				[]types.UnbondingDelegation{
					types.NewPoolUnbondingDelegation(
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						1,
						1,
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
						sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
						1,
					),
					types.NewOperatorUnbondingDelegation(
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						1,
						1,
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
						sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
						1,
					),
					types.NewServiceUnbondingDelegation(
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						1,
						1,
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
						sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
						1,
					),
				},
				[]types.UserPreferencesEntry{
					types.NewUserPreferencesEntry(
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						types.NewUserPreferences([]types.TrustedServiceEntry{
							types.NewTrustedServiceEntry(1, nil),
							types.NewTrustedServiceEntry(2, nil),
							types.NewTrustedServiceEntry(3, nil),
						}),
					),
				},
				types.NewParams(5*24*time.Hour, nil, sdkmath.LegacyNewDec(100000), types.DefaultMaxEntries),
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
