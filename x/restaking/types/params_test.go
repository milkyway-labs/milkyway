package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/v2/x/restaking/types"
)

func TestParams_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		params    types.Params
		shouldErr bool
	}{
		{
			name:      "invalid unbonding time returns error",
			params:    types.NewParams(0, nil),
			shouldErr: true,
		},
		{
			name:      "invalid denom returns error",
			params:    types.NewParams(5, []string{"1denom"}),
			shouldErr: true,
		},
		{
			name:      "empty denom returns error",
			params:    types.NewParams(5, []string{""}),
			shouldErr: true,
		},
		{
			name:      "default params return no error",
			params:    types.DefaultParams(),
			shouldErr: false,
		},
		{
			name:      "valid params return no error",
			params:    types.NewParams(5*time.Hour, nil),
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
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
