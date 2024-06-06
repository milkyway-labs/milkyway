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

				err := suite.k.CreateService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
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
					),
					types.NewService(
						2,
						types.SERVICE_STATUS_INACTIVE,
						"Inertia",
						"AVS-based Liquid Restaking Platform",
						"https://inertia.zone",
						"https://inertia.zone/logo.png",
						"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
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
