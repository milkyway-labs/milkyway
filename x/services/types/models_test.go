package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/x/services/types"
)

func TestParseServiceID(t *testing.T) {
	testCases := []struct {
		name     string
		id       string
		expID    uint32
		expError bool
	}{
		{
			name:  "valid ID returns no error",
			id:    "1",
			expID: 1,
		},
		{
			name:     "invalid ID returns error",
			id:       "invalid",
			expError: true,
		},
		{
			name:     "empty ID returns error",
			id:       "",
			expError: true,
		},
		{
			name:     "negative ID returns error",
			id:       "-1",
			expError: true,
		},
		{
			name:     "zero ID returns no error",
			id:       "0",
			expError: false,
			expID:    0,
		},
		{
			name:     "max uint32 returns no error",
			id:       "4294967295",
			expError: false,
			expID:    4294967295,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			id, err := types.ParseServiceID(tc.id)
			if tc.expError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expID, id)
			}
		})
	}
}

func TestService_Validate(t *testing.T) {
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

func TestService_Update(t *testing.T) {
	testCases := []struct {
		name      string
		avs       types.Service
		update    types.ServiceUpdate
		expResult types.Service
	}{
		{
			name: "update name",
			avs: types.Service{
				ID:     1,
				Status: types.SERVICE_STATUS_CREATED,
				Admin:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Name:   "MilkyWay",
			},
			update: types.NewServiceUpdate(
				"MilkyWay2",
				types.DoNotModify,
				types.DoNotModify,
				types.DoNotModify,
			),
			expResult: types.Service{
				ID:     1,
				Status: types.SERVICE_STATUS_CREATED,
				Admin:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Name:   "MilkyWay2",
			},
		},
		{
			name: "update description",
			avs: types.Service{
				ID:     1,
				Status: types.SERVICE_STATUS_CREATED,
				Admin:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Name:   "MilkyWay",
			},
			update: types.NewServiceUpdate(
				types.DoNotModify,
				"New description",
				types.DoNotModify,
				types.DoNotModify,
			),
			expResult: types.Service{
				ID:          1,
				Status:      types.SERVICE_STATUS_CREATED,
				Admin:       "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Name:        "MilkyWay",
				Description: "New description",
			},
		},
		{
			name: "update website",
			avs: types.Service{
				ID:     1,
				Status: types.SERVICE_STATUS_CREATED,
				Admin:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Name:   "MilkyWay",
			},
			update: types.NewServiceUpdate(
				types.DoNotModify,
				types.DoNotModify,
				"https://example.com",
				types.DoNotModify,
			),
			expResult: types.Service{
				ID:      1,
				Status:  types.SERVICE_STATUS_CREATED,
				Admin:   "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Name:    "MilkyWay",
				Website: "https://example.com",
			},
		},
		{
			name: "update picture URL",
			avs: types.Service{
				ID:     1,
				Status: types.SERVICE_STATUS_CREATED,
				Admin:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Name:   "MilkyWay",
			},
			update: types.NewServiceUpdate(
				types.DoNotModify,
				types.DoNotModify,
				types.DoNotModify,
				"https://example.com/picture.jpg",
			),
			expResult: types.Service{
				ID:         1,
				Status:     types.SERVICE_STATUS_CREATED,
				Admin:      "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Name:       "MilkyWay",
				PictureURL: "https://example.com/picture.jpg",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := tc.avs.Update(tc.update)
			require.Equal(t, tc.expResult, result)
		})
	}
}
