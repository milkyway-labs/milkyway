package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

func TestParams_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		params    types.Params
		shouldErr bool
	}{
		{
			name: "invalid deactivation time returns error",
			params: types.Params{
				DeactivationTime: 0,
			},
			shouldErr: true,
		},
		{
			name:      "default params do not return errors",
			params:    types.DefaultParams(),
			shouldErr: false,
		},
		{
			name:      "valid params do not return errors",
			params:    types.NewParams(60),
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.params.Validate()
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
