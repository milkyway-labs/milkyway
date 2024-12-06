package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/v3/x/pools/types"
)

func TestGenesis_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		genesis   *types.GenesisState
		shouldErr bool
	}{
		{
			name:      "invalid next pool id returns error",
			genesis:   types.NewGenesis(0, nil),
			shouldErr: true,
		},
		{
			name: "invalid pool returns error",
			genesis: types.NewGenesis(
				1,
				[]types.Pool{
					types.NewPool(0, "uatom"),
				},
			),
			shouldErr: true,
		},
		{
			name:      "default genesis does not return errors",
			genesis:   types.DefaultGenesis(),
			shouldErr: false,
		},
		{
			name: "valid genesis does not return errors",
			genesis: types.NewGenesis(
				1,
				[]types.Pool{
					types.NewPool(1, "uatom"),
				},
			),
			shouldErr: false,
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
