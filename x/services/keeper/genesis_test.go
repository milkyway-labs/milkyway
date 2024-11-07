package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/services/types"
)

func (suite *KeeperTestSuite) TestKeeper_ExportGenesis() {
	testCases := []struct {
		name       string
		setup      func()
		store      func(ctx sdk.Context)
		expGenesis *types.GenesisState
	}{
		{
			name: "next service id is exported properly",
			store: func(ctx sdk.Context) {
				suite.k.SetNextServiceID(ctx, 10)
				suite.k.SetParams(ctx, types.DefaultParams())
			},
			expGenesis: &types.GenesisState{
				NextServiceID: 10,
				Services:      nil,
				Params:        types.DefaultParams(),
			},
		},
		{
			name: "services data are exported properly",
			store: func(ctx sdk.Context) {
				suite.k.SetNextServiceID(ctx, 1)
				suite.k.SetParams(ctx, types.DefaultParams())

				err := suite.k.SaveService(ctx, types.NewService(
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

				err = suite.k.SaveService(ctx, types.NewService(
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
			expGenesis: &types.GenesisState{
				NextServiceID: 1,
				Services: []types.Service{
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
				Params: types.DefaultParams(),
			},
		},
		{
			name: "params are exported properly",
			store: func(ctx sdk.Context) {
				suite.k.SetNextServiceID(ctx, 1)
				suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(1_000_000_000))),
				))
			},
			expGenesis: &types.GenesisState{
				NextServiceID: 1,
				Services:      nil,
				Params: types.NewParams(
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(1_000_000_000))),
				),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.setup != nil {
				tc.setup()
			}
			if tc.store != nil {
				tc.store(ctx)
			}

			genesis := suite.k.ExportGenesis(ctx)
			suite.Require().Equal(tc.expGenesis, genesis)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_InitGenesis() {
	testCases := []struct {
		name      string
		genesis   *types.GenesisState
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name: "next service id is initialized properly",
			genesis: types.NewGenesisState(
				10,
				nil,
				types.DefaultParams(),
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				nextServiceID, err := suite.k.GetNextServiceID(ctx)
				suite.Require().NoError(err)
				suite.Require().Equal(uint32(10), nextServiceID)
			},
		},
		{
			name: "services data are initialized properly",
			genesis: types.NewGenesisState(
				1,
				[]types.Service{
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
				},
				types.DefaultParams(),
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				var services []types.Service
				suite.k.IterateServices(ctx, func(service types.Service) (stop bool) {
					services = append(services, service)
					return false
				})

				suite.Require().Len(services, 1)
				suite.Require().Equal(types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				), services[0])
			},
		},
		{
			name: "params are initialized properly",
			genesis: types.NewGenesisState(
				1,
				nil,
				types.NewParams(
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(1_000_000_000))),
				),
			),
			check: func(ctx sdk.Context) {
				params := suite.k.GetParams(ctx)
				suite.Require().Equal(types.NewParams(
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(1_000_000_000))),
				), params)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()

			err := suite.k.InitGenesis(ctx, tc.genesis)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}
