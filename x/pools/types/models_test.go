package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/x/pools/types"
)

func TestPool_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		pool      types.Pool
		shouldErr bool
	}{
		{
			name:      "invalid pool id returns error",
			pool:      types.NewPool(0, "uatom"),
			shouldErr: true,
		},
		{
			name:      "invalid pool denom returns error",
			pool:      types.NewPool(1, "uatom!"),
			shouldErr: true,
		},
		{
			name: "valid pool does not return errors",
			pool: types.NewPool(1, "uatom"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.pool.Validate()
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
