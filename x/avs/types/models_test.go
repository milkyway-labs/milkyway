package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/x/avs/types"
)

func TestAVS_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		avs       types.AVS
		shouldErr bool
	}{
		{
			name:      "invalid ID returns error",
			avs:       types.NewAVS(0, "MilkyWay", "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			shouldErr: true,
		},
		{
			name:      "invalid name returns error",
			avs:       types.NewAVS(1, "", "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			shouldErr: true,
		},
		{
			name:      "invalid address returns error",
			avs:       types.NewAVS(1, "MilkyWay", "invalid_address"),
			shouldErr: true,
		},
		{
			name:      "valid AVS returns no error",
			avs:       types.NewAVS(1, "MilkyWay", "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.avs.Validate()
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
