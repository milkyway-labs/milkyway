package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/v2/x/assets/types"
)

func TestGenesis_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		genesis *types.GenesisState
		expErr  bool
	}{
		{
			name:    "default genesis returns no error",
			genesis: types.DefaultGenesis(),
			expErr:  false,
		},
		{
			name: "invalid asset returns error",
			genesis: types.NewGenesisState(
				[]types.Asset{
					types.NewAsset("@#$%", "bitcoin", 1),
				},
			),
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.genesis.Validate()
			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
