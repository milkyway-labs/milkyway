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
		service   types.Service
		shouldErr bool
	}{
		{
			name: "invalid status returns error",
			service: types.NewService(
				1,
				types.SERVICE_STATUS_UNSPECIFIED,
				"MilkyWay",
				"MilkyWay is an AVS of a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				types.GetServiceAddress(1).String(),
			),
			shouldErr: true,
		},
		{
			name: "invalid ID returns error",
			service: types.NewService(
				0,
				types.SERVICE_STATUS_ACTIVE,
				"MilkyWay",
				"MilkyWay is an AVS of a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				types.GetServiceAddress(1).String(),
			),
			shouldErr: true,
		},
		{
			name: "invalid name returns error",
			service: types.NewService(
				1,
				types.SERVICE_STATUS_ACTIVE,
				"",
				"MilkyWay is an AVS of a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				types.GetServiceAddress(1).String(),
			),
			shouldErr: true,
		},
		{
			name: "invalid admin address returns error",
			service: types.NewService(
				1,
				types.SERVICE_STATUS_ACTIVE,
				"MilkyWay",
				"MilkyWay is an AVS of a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				"",
				types.GetServiceAddress(1).String(),
			),
			shouldErr: true,
		},
		{
			name: "invalid address returns error",
			service: types.NewService(
				1,
				types.SERVICE_STATUS_ACTIVE,
				"MilkyWay",
				"MilkyWay is an AVS of a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				"",
			),
			shouldErr: true,
		},
		{
			name: "valid service returns no error",
			service: types.NewService(
				1,
				types.SERVICE_STATUS_ACTIVE,
				"MilkyWay",
				"MilkyWay is an AVS of a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				types.GetServiceAddress(1).String(),
			),
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.service.Validate()
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
		service   types.Service
		update    types.ServiceUpdate
		expResult types.Service
	}{
		{
			name: "update name",
			service: types.NewService(
				1,
				types.SERVICE_STATUS_ACTIVE,
				"MilkyWay",
				"MilkyWay is an AVS of a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				types.GetServiceAddress(1).String(),
			),
			update: types.NewServiceUpdate(
				"MilkyWay2",
				types.DoNotModify,
				types.DoNotModify,
				types.DoNotModify,
			),
			expResult: types.NewService(
				1,
				types.SERVICE_STATUS_ACTIVE,
				"MilkyWay2",
				"MilkyWay is an AVS of a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				types.GetServiceAddress(1).String(),
			),
		},
		{
			name: "update description",
			service: types.NewService(
				1,
				types.SERVICE_STATUS_ACTIVE,
				"MilkyWay",
				"MilkyWay is an AVS of a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				types.GetServiceAddress(1).String(),
			),
			update: types.NewServiceUpdate(
				types.DoNotModify,
				"New description",
				types.DoNotModify,
				types.DoNotModify,
			),
			expResult: types.NewService(
				1,
				types.SERVICE_STATUS_ACTIVE,
				"MilkyWay",
				"New description",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				types.GetServiceAddress(1).String(),
			),
		},
		{
			name: "update website",
			service: types.NewService(
				1,
				types.SERVICE_STATUS_ACTIVE,
				"MilkyWay",
				"MilkyWay is an AVS of a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				types.GetServiceAddress(1).String(),
			),
			update: types.NewServiceUpdate(
				types.DoNotModify,
				types.DoNotModify,
				"https://example.com",
				types.DoNotModify,
			),
			expResult: types.NewService(
				1,
				types.SERVICE_STATUS_ACTIVE,
				"MilkyWay",
				"MilkyWay is an AVS of a restaking platform",
				"https://example.com",
				"https://milkyway.com/logo.png",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				types.GetServiceAddress(1).String(),
			),
		},
		{
			name: "update picture URL",
			service: types.NewService(
				1,
				types.SERVICE_STATUS_ACTIVE,
				"MilkyWay",
				"MilkyWay is an AVS of a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				types.GetServiceAddress(1).String(),
			),
			update: types.NewServiceUpdate(
				types.DoNotModify,
				types.DoNotModify,
				types.DoNotModify,
				"https://example.com/picture.jpg",
			),
			expResult: types.NewService(
				1,
				types.SERVICE_STATUS_ACTIVE,
				"MilkyWay",
				"MilkyWay is an AVS of a restaking platform",
				"https://milkyway.com",
				"https://example.com/picture.jpg",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				types.GetServiceAddress(1).String(),
			),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := tc.service.Update(tc.update)
			require.Equal(t, tc.expResult, result)
		})
	}
}
