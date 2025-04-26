package cli_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/v12/x/restaking/client/cli"
	"github.com/milkyway-labs/milkyway/v12/x/restaking/types"
)

func TestParseTrustedServiceEntry(t *testing.T) {
	testCases := []struct {
		name      string
		value     string
		shouldErr bool
		expected  types.TrustedServiceEntry
	}{
		{
			name:      "basic entry",
			value:     "1-1,2,3",
			shouldErr: false,
			expected:  types.NewTrustedServiceEntry(1, []uint32{1, 2, 3}),
		},
		{
			name:      "specifying one pool ID is valid",
			value:     "1-1",
			shouldErr: false,
			expected:  types.NewTrustedServiceEntry(1, []uint32{1}),
		},
		{
			name:      "specifying only the service ID is valid",
			value:     "1",
			shouldErr: false,
			expected:  types.NewTrustedServiceEntry(1, nil),
		},
		{
			name:      "malformed input returns an error #1",
			value:     "1-",
			shouldErr: true,
		},
		{
			name:      "malformed input returns an error #2",
			value:     "-",
			shouldErr: true,
		},
		{
			name:      "malformed input returns an error #3",
			value:     "",
			shouldErr: true,
		},
		{
			name:      "malformed input returns an error #4",
			value:     "1-2-3",
			shouldErr: true,
		},
		{
			name:      "specifying more than one service ID returns an error",
			value:     "1,2-1,2,3",
			shouldErr: true,
		},
		{
			name:      "specifying no service ID returns an error",
			value:     "-1,2,3",
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entry, err := cli.ParseTrustedServiceEntry(tc.value)
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, entry)
			}
		})
	}
}
