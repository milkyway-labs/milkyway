package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v6/x/services/keeper"
	"github.com/milkyway-labs/milkyway/v6/x/services/types"
)

func (suite *KeeperTestSuite) TestValidServicesInvariant() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		expBroken bool
	}{
		{
			name: "not found next service id breaks invariant",
			store: func(ctx sdk.Context) {
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
			},
			expBroken: true,
		},
		{
			name: "service with id equals to next service id breaks invariant",
			store: func(ctx sdk.Context) {
				err := suite.k.SetNextServiceID(ctx, 1)
				suite.Require().NoError(err)

				err = suite.k.SaveService(ctx, types.NewService(
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
			expBroken: true,
		},
		{
			name: "service with id higher than next service id breaks invariant",
			store: func(ctx sdk.Context) {
				err := suite.k.SetNextServiceID(ctx, 1)
				suite.Require().NoError(err)

				err = suite.k.SaveService(ctx, types.NewService(
					2,
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
			expBroken: true,
		},
		{
			name: "invalid service breaks invariant",
			store: func(ctx sdk.Context) {
				err := suite.k.SetNextServiceID(ctx, 1)
				suite.Require().NoError(err)

				err = suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_UNSPECIFIED,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
			},
			expBroken: true,
		},
		{
			name: "valid data does not break invariant",
			store: func(ctx sdk.Context) {
				err := suite.k.SetNextServiceID(ctx, 2)
				suite.Require().NoError(err)

				err = suite.k.SaveService(ctx, types.NewService(
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
			expBroken: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
			if tc.setup != nil {
				tc.setup()
			}
			if tc.store != nil {
				tc.store(ctx)
			}

			_, broken := keeper.ValidServicesInvariant(suite.k)(ctx)
			suite.Require().Equal(tc.expBroken, broken)
		})
	}
}
