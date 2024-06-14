package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

func TestSplitInactivatingOperatorQueueKey(t *testing.T) {
	testCases := []struct {
		name          string
		key           []byte
		shouldErr     bool
		expOperatorID uint32
		expTime       time.Time
	}{
		{
			name:      "invalid key panics",
			key:       []byte("invalid"),
			shouldErr: true,
		},
		{
			name:          "valid key is parsed properly",
			key:           types.InactivatingOperatorQueueKey(1, time.Date(2025, 1, 1, 12, 30, 0, 0, time.UTC)),
			shouldErr:     false,
			expOperatorID: 1,
			expTime:       time.Date(2025, 1, 1, 12, 30, 0, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if tc.shouldErr {
				require.Panics(t, func() {
					types.SplitInactivatingOperatorQueueKey(tc.key)
				})
			} else {
				operatorID, parsedTime := types.SplitInactivatingOperatorQueueKey(tc.key)
				require.Equal(t, tc.expOperatorID, operatorID)
				require.True(t, parsedTime.Equal(tc.expTime))
			}
		})
	}
}
