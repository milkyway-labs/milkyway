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
				[]types.Delegation{
					types.NewPoolDelegation(
						0,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
					),
				},
				types.DefaultParams(),
			),
			shouldErr: true,
		},
		{
			name: "invalid service delegation entry returns error",
			genesis: types.NewGenesis(
				[]types.Delegation{
					types.NewServiceDelegation(
						0,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
					),
				},
				types.DefaultParams(),
			),
			shouldErr: true,
		},
		{
			name: "invalid operator delegation entry returns error",
			genesis: types.NewGenesis(
				[]types.Delegation{
					types.NewOperatorDelegation(
						0,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
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
