package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/x/services/types"
)

func TestAVS_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		avs       types.Service
		shouldErr bool
	}{
		{
			name: "invalid status returns error",
			avs: types.Service{
				ID:     1,
				Status: types.SERVICE_STATUS_UNSPECIFIED,
				Admin:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Name:   "MilkyWay",
			},
			shouldErr: true,
		},
		{
			name: "invalid ID returns error",
			avs: types.Service{
				ID:     0,
				Status: types.SERVICE_STATUS_CREATED,
				Admin:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Name:   "MilkyWay",
			},
			shouldErr: true,
		},
		{
			name: "invalid name returns error",
			avs: types.Service{
				ID:     1,
				Status: types.SERVICE_STATUS_CREATED,
				Admin:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Name:   "",
			},
			shouldErr: true,
		},
		{
			name: "invalid address returns error",
			avs: types.Service{
				ID:     1,
				Status: types.SERVICE_STATUS_CREATED,
				Admin:  "",
				Name:   "MilkyWay",
			},
			shouldErr: true,
		},
		{
			name: "valid Service returns no error",
			avs: types.Service{
				ID:     1,
				Status: types.SERVICE_STATUS_CREATED,
				Admin:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Name:   "MilkyWay",
			},
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
