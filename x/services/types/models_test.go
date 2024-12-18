package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/v5/x/services/types"
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
				false,
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
				false,
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
				false,
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
				false,
			),
			shouldErr: true,
		},
		{
			name: "invalid address returns error",
			service: types.Service{
				ID:      1,
				Status:  types.SERVICE_STATUS_ACTIVE,
				Name:    "MilkyWay",
				Admin:   "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				Address: "",
			},
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
				false,
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
				false,
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
				false,
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
				false,
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
				false,
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
				false,
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
				false,
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
				false,
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
				false,
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

func TestService_SharesFromTokens(t *testing.T) {
	testCases := []struct {
		name      string
		service   types.Service
		tokens    sdk.Coins
		shouldErr bool
		expShares sdk.DecCoins
	}{
		{
			name: "service with no delegation shares returns error",
			service: types.Service{
				ID:              1,
				Address:         types.GetServiceAddress(1).String(),
				DelegatorShares: sdk.NewDecCoins(),
				Tokens:          sdk.NewCoins(),
			},
			tokens: sdk.NewCoins(
				sdk.NewCoin("umilk", sdkmath.NewInt(100)),
			),
			shouldErr: true,
		},
		{
			name: "shares are computed properly for non empty operator",
			service: types.Service{
				ID:      1,
				Address: types.GetServiceAddress(1).String(),
				Tokens: sdk.NewCoins(
					sdk.NewCoin("umilk", sdkmath.NewInt(50)),
				),
				DelegatorShares: sdk.NewDecCoins(
					sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(100)),
				),
			},
			tokens: sdk.NewCoins(
				sdk.NewCoin("umilk", sdkmath.NewInt(20)),
			),
			shouldErr: false,
			expShares: sdk.NewDecCoins(
				sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(40)),
			),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			shares, err := tc.service.SharesFromTokens(tc.tokens)
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expShares, shares)
			}
		})
	}
}

func TestService_TokensFromShares(t *testing.T) {
	testCases := []struct {
		name      string
		service   types.Service
		shares    sdk.DecCoins
		expTokens sdk.DecCoins
	}{
		{
			name: "service with shares returns correct amount",
			service: types.Service{
				ID:      1,
				Address: types.GetServiceAddress(1).String(),
				Tokens: sdk.NewCoins(
					sdk.NewCoin("umilk", sdkmath.NewInt(70)),
				),
				DelegatorShares: sdk.NewDecCoins(
					sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(140)),
				),
			},
			shares: sdk.NewDecCoins(
				sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(40)),
			),
			expTokens: sdk.NewDecCoins(
				sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(20)),
			),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tokens := tc.service.TokensFromShares(tc.shares)
			require.Equal(t, tc.expTokens, tokens)
		})
	}
}
