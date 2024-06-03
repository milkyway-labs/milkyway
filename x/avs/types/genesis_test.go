package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/x/avs/types"
)

func TestValidateGenesis(t *testing.T) {
	testCases := []struct {
		name      string
		genesis   *types.GenesisState
		shouldErr bool
	}{
		{
			name: "duplicated service returns error",
			genesis: &types.GenesisState{
				Services: []types.AVS{
					types.NewAVS(1, "MilkyWay", "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
					types.NewAVS(1, "MilkyWay", "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
				},
				Params: types.DefaultParams(),
			},
			shouldErr: true,
		},
		{
			name: "invalid service returns error",
			genesis: &types.GenesisState{
				Services: []types.AVS{
					types.NewAVS(1, "MilkyWay", "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
					types.NewAVS(2, "IBC Relaying", ""),
				},
				Params: types.DefaultParams(),
			},
			shouldErr: true,
		},
		{
			name: "invalid params returns error",
			genesis: &types.GenesisState{
				Services: nil,
				Params: types.Params{
					AvsRegistrationFee: sdk.Coins{sdk.Coin{Denom: "", Amount: sdkmath.NewInt(10)}},
				},
			},
			shouldErr: true,
		},
		{
			name:      "default genesis is valid",
			genesis:   types.DefaultGenesisState(),
			shouldErr: false,
		},
		{
			name: "valid genesis returns no error",
			genesis: &types.GenesisState{
				Services: []types.AVS{
					types.NewAVS(1, "MilkyWay", "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
					types.NewAVS(2, "IBC Relaying", "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
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
			err := types.ValidateGenesis(tc.genesis)
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
