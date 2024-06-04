package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/x/services/types"
)

func TestValidateGenesis(t *testing.T) {
	testCases := []struct {
		name      string
		genesis   *types.GenesisState
		shouldErr bool
	}{
		{
			name: "invalid next Service ID returns error",
			genesis: &types.GenesisState{
				NextAVSID: 0,
				Services:  nil,
				Params:    types.DefaultParams(),
			},
			shouldErr: true,
		},
		{
			name: "duplicated service returns error",
			genesis: &types.GenesisState{
				NextAVSID: 1,
				Services: []types.Service{
					{
						ID:     1,
						Status: types.AVS_STATUS_CREATED,
						Admin:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						Name:   "MilkyWay",
					},
					{
						ID:     1,
						Status: types.AVS_STATUS_CREATED,
						Admin:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						Name:   "MilkyWay",
					}},
				Params: types.DefaultParams(),
			},
			shouldErr: true,
		},
		{
			name: "invalid service returns error",
			genesis: &types.GenesisState{
				NextAVSID: 1,
				Services: []types.Service{
					{
						ID:     1,
						Status: types.AVS_STATUS_CREATED,
						Admin:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						Name:   "MilkyWay",
					},
					{
						ID:     2,
						Status: types.AVS_STATUS_UNSPECIFIED,
						Admin:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						Name:   "IBC Relaying",
					},
				},
				Params: types.DefaultParams(),
			},
			shouldErr: true,
		},
		{
			name: "invalid params returns error",
			genesis: &types.GenesisState{
				NextAVSID: 1,
				Services:  nil,
				Params: types.Params{
					AvsRegistrationFee: sdk.Coins{sdk.Coin{Denom: "", Amount: sdkmath.NewInt(10)}},
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
				NextAVSID: 1,
				Services: []types.Service{
					{
						ID:     1,
						Status: types.AVS_STATUS_CREATED,
						Admin:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						Name:   "MilkyWay",
					},
					{
						ID:     2,
						Status: types.AVS_STATUS_REGISTERED,
						Admin:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						Name:   "IBC Relaying",
					},
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
