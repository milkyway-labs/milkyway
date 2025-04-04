package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/v10/app/testutil"
	"github.com/milkyway-labs/milkyway/v10/utils"
	"github.com/milkyway-labs/milkyway/v10/x/investors/types"
)

func TestGenesisState_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		state     *types.GenesisState
		shouldErr bool
	}{
		{
			name: "valid genesis",
			state: types.NewGenesisState(
				utils.MustParseDec("0.5"),
				[]string{
					testutil.TestAddress(1).String(),
					testutil.TestAddress(2).String(),
				},
			),
			shouldErr: false,
		},
		{
			name: "invalid investors reward ratio",
			state: types.NewGenesisState(
				utils.MustParseDec("1.1"),
				[]string{
					testutil.TestAddress(1).String(),
					testutil.TestAddress(2).String(),
				},
			),
			shouldErr: true,
		},
		{
			name: "invalid investor address",
			state: types.NewGenesisState(
				utils.MustParseDec("0.5"),
				[]string{
					"invalid",
					testutil.TestAddress(2).String(),
				},
			),
			shouldErr: true,
		},
		{
			name: "duplicated investor address",
			state: types.NewGenesisState(
				utils.MustParseDec("0.5"),
				[]string{
					testutil.TestAddress(1).String(),
					testutil.TestAddress(1).String(),
				},
			),
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.state.Validate()
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
