package types_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/v12/x/restaking/types"
)

func TestParams_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		params    types.Params
		shouldErr bool
	}{
		{
			name:      "invalid unbonding time returns error",
			params:    types.NewParams(0, nil, types.DefaultRestakingCap, types.DefaultMaxEntries),
			shouldErr: true,
		},
		{
			name:      "invalid denom returns error",
			params:    types.NewParams(5, []string{"1denom"}, types.DefaultRestakingCap, types.DefaultMaxEntries),
			shouldErr: true,
		},
		{
			name:      "empty denom returns error",
			params:    types.NewParams(5, []string{""}, types.DefaultRestakingCap, types.DefaultMaxEntries),
			shouldErr: true,
		},
		{
			name:      "negative restaking cap returns error",
			params:    types.NewParams(5, nil, math.LegacyNewDec(-1), types.DefaultMaxEntries),
			shouldErr: true,
		},
		{
			name:      "zero max entries returns error",
			params:    types.NewParams(5, nil, types.DefaultRestakingCap, 0),
			shouldErr: true,
		},
		{
			name:      "default params return no error",
			params:    types.DefaultParams(),
			shouldErr: false,
		},
		{
			name:      "valid params return no error",
			params:    types.NewParams(5*time.Hour, nil, math.LegacyNewDec(100000), types.DefaultMaxEntries),
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
