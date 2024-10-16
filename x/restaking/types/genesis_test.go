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
			name: "operator secured services record with invalid SecuredService returns error",
			genesis: types.NewGenesis(
				[]types.OperatorSecuredServicesRecord{
					{
						OperatorID:      1,
						SecuredServices: types.NewOperatorSecuredServices([]uint32{0}),
					},
				},
				nil,
				nil,
				nil,
				types.DefaultParams(),
			),
			shouldErr: true,
		},
		{
			name: "operator secured services record with invalid OperatorID returns error",
			genesis: types.NewGenesis(
				[]types.OperatorSecuredServicesRecord{
					{
						OperatorID:      0,
						SecuredServices: types.NewOperatorSecuredServices([]uint32{1}),
					},
				},
				nil,
				nil,
				nil,
				types.DefaultParams(),
			),
			shouldErr: true,
		},
		{
			name: "invalid service params record returns error",
			genesis: types.NewGenesis(
				nil,
				[]types.ServiceParamsRecord{
					{
						ServiceID: 1,
						Params:    types.NewServiceParams(sdkmath.LegacyNewDec(2), nil, nil),
					},
				},
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
				[]types.Delegation{
					types.NewPoolDelegation(
						0,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
					),
				},
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
				[]types.Delegation{
					types.NewServiceDelegation(
						0,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
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
				[]types.Delegation{
					types.NewOperatorDelegation(
						0,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
					),
				},
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
				[]types.OperatorSecuredServicesRecord{
					{
						OperatorID:      1,
						SecuredServices: types.NewOperatorSecuredServices([]uint32{2, 3, 5}),
					},
				},
				[]types.ServiceParamsRecord{
					{
						ServiceID: 2,
						Params: types.NewServiceParams(
							sdkmath.LegacyNewDecWithPrec(1, 2), []uint32{1, 2, 3}, []uint32{1, 5}),
					},
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
