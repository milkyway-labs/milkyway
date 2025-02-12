package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/v9/x/services/types"
)

func TestValidateGenesis(t *testing.T) {
	testCases := []struct {
		name      string
		genesis   *types.GenesisState
		shouldErr bool
	}{
		{
			name: "invalid next service ID returns error",
			genesis: &types.GenesisState{
				NextServiceID: 0,
				Services:      nil,
				Params:        types.DefaultParams(),
			},
			shouldErr: true,
		},
		{
			name: "duplicated service returns error",
			genesis: &types.GenesisState{
				NextServiceID: 1,
				Services: []types.Service{
					types.NewService(
						1,
						types.SERVICE_STATUS_ACTIVE,
						"MilkyWay",
						"MilkyWay is an AVS of a restaking platform",
						"https://milkyway.com",
						"https://milkyway.com/logo.png",
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						false,
					),
					types.NewService(
						1,
						types.SERVICE_STATUS_ACTIVE,
						"MilkyWay",
						"MilkyWay is an AVS of a restaking platform",
						"https://milkyway.com",
						"https://milkyway.com/logo.png",
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						false,
					),
				},
				Params: types.DefaultParams(),
			},
			shouldErr: true,
		},
		{
			name: "invalid service returns error",
			genesis: &types.GenesisState{
				NextServiceID: 1,
				Services: []types.Service{
					types.NewService(
						1,
						types.SERVICE_STATUS_ACTIVE,
						"MilkyWay",
						"MilkyWay is an AVS of a restaking platform",
						"https://milkyway.com",
						"https://milkyway.com/logo.png",
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						false,
					),
					types.NewService(
						2,
						types.SERVICE_STATUS_UNSPECIFIED,
						"Tucana",
						"",
						"",
						"",
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						false,
					),
				},
				Params: types.DefaultParams(),
			},
			shouldErr: true,
		},
		{
			name: "invalid service params returns error",
			genesis: &types.GenesisState{
				NextServiceID: 1,
				Services:      nil,
				ServicesParams: []types.ServiceParamsRecord{
					types.NewServiceParamsRecord(0, types.NewServiceParams([]string{"umilk"})),
				},
				Params: types.DefaultParams(),
			},
			shouldErr: true,
		},
		{
			name: "invalid params returns error",
			genesis: &types.GenesisState{
				NextServiceID: 1,
				Services:      nil,
				Params: types.Params{
					ServiceRegistrationFee: sdk.Coins{sdk.Coin{Denom: "", Amount: sdkmath.NewInt(10)}},
				},
			},
			shouldErr: true,
		},
		{
			name:      "default genesis is valid",
			genesis:   types.DefaultGenesis(),
			shouldErr: false,
		},
		{
			name: "valid genesis returns no error",
			genesis: &types.GenesisState{
				NextServiceID: 1,
				Services: []types.Service{
					types.NewService(
						1,
						types.SERVICE_STATUS_ACTIVE,
						"MilkyWay",
						"MilkyWay is an AVS of a restaking platform",
						"https://milkyway.com",
						"https://milkyway.com/logo.png",
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						false,
					),
					types.NewService(
						2,
						types.SERVICE_STATUS_ACTIVE,
						"Tucana",
						"",
						"",
						"",
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						false,
					),
				},
				Params: types.NewParams(
					sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(10))),
				),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
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
