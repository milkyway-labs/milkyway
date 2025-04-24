package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/milkyway-labs/milkyway/v11/x/services/types"
)

func (suite *KeeperTestSuite) TestKeeper_Services() {
	testCases := []struct {
		name        string
		store       func(ctx sdk.Context)
		request     *types.QueryServicesRequest
		shouldErr   bool
		expServices []types.Service
	}{
		{
			name: "query without pagination returns data properly",
			store: func(ctx sdk.Context) {
				err := suite.k.CreateService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)

				err = suite.k.CreateService(ctx, types.NewService(
					2,
					types.SERVICE_STATUS_INACTIVE,
					"Inertia",
					"AVS-based Liquid Restaking Platform",
					"https://inertia.zone",
					"https://inertia.zone/logo.png",
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					false,
				))
				suite.Require().NoError(err)
			},
			request:   types.NewQueryServicesRequest(nil),
			shouldErr: false,
			expServices: []types.Service{
				types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				),
				types.NewService(
					2,
					types.SERVICE_STATUS_INACTIVE,
					"Inertia",
					"AVS-based Liquid Restaking Platform",
					"https://inertia.zone",
					"https://inertia.zone/logo.png",
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					false,
				),
			},
		},
		{
			name: "query with pagination returns data properly",
			store: func(ctx sdk.Context) {
				err := suite.k.CreateService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)

				err = suite.k.CreateService(ctx, types.NewService(
					2,
					types.SERVICE_STATUS_INACTIVE,
					"Inertia",
					"AVS-based Liquid Restaking Platform",
					"https://inertia.zone",
					"https://inertia.zone/logo.png",
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					false,
				))
				suite.Require().NoError(err)
			},
			request: types.NewQueryServicesRequest(&query.PageRequest{
				Offset: 1,
				Limit:  1,
			}),
			shouldErr: false,
			expServices: []types.Service{
				types.NewService(
					2,
					types.SERVICE_STATUS_INACTIVE,
					"Inertia",
					"AVS-based Liquid Restaking Platform",
					"https://inertia.zone",
					"https://inertia.zone/logo.png",
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					false,
				),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			res, err := suite.k.Services(ctx, tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expServices, res.Services)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_Service() {
	testCases := []struct {
		name       string
		store      func(ctx sdk.Context)
		request    *types.QueryServiceRequest
		shouldErr  bool
		expService types.Service
	}{
		{
			name:      "invalid service id returns error",
			request:   types.NewQueryServiceRequest(0),
			shouldErr: true,
		},
		{
			name:      "not found service returns error",
			request:   types.NewQueryServiceRequest(1),
			shouldErr: true,
		},
		{
			name: "found service is returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.CreateService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
			},
			request:   types.NewQueryServiceRequest(1),
			shouldErr: false,
			expService: types.NewService(
				1,
				types.SERVICE_STATUS_ACTIVE,
				"MilkyWay",
				"MilkyWay is an AVS of a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				false,
			),
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			res, err := suite.k.Service(ctx, tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expService, res.Service)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_ServiceParams() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		request   *types.QueryServiceParamsRequest
		shouldErr bool
		expParams types.ServiceParams
	}{
		{
			name:      "invalid service id returns error",
			request:   types.NewQueryServiceParamsRequest(0),
			shouldErr: true,
		},
		{
			name:      "not found service returns error",
			request:   types.NewQueryServiceParamsRequest(1),
			shouldErr: true,
		},
		{
			name: "default params are returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.CreateService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
			},
			request:   types.NewQueryServiceParamsRequest(1),
			shouldErr: false,
			expParams: types.DefaultServiceParams(),
		},
		{
			name: "found params are returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.CreateService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
				// Set some custom parameters
				err = suite.k.SetServiceParams(ctx, 1, types.NewServiceParams([]string{"umilk"}))
				suite.Require().NoError(err)
			},
			request:   types.NewQueryServiceParamsRequest(1),
			shouldErr: false,
			expParams: types.NewServiceParams([]string{"umilk"}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			res, err := suite.k.ServiceParams(ctx, tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expParams, res.ServiceParams)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQueryServer_Params() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		request   *types.QueryParamsRequest
		expParams types.Params
	}{
		{
			name: "params are returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(1000)))),
				)
				suite.Require().NoError(err)
			},
			request:   types.NewQueryParamsRequest(),
			expParams: types.NewParams(sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(1000)))),
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			res, err := suite.k.Params(ctx, tc.request)
			suite.Require().NoError(err)
			suite.Require().Equal(tc.expParams, res.Params)
		})
	}
}
