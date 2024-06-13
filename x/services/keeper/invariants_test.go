package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/services/keeper"
	"github.com/milkyway-labs/milkyway/x/services/types"
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
				suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					types.GetServiceAddress(1).String(),
				))
			},
			expBroken: true,
		},
		{
			name: "service with id equals to next service id breaks invariant",
			store: func(ctx sdk.Context) {
				suite.k.SetNextServiceID(ctx, 1)
				suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					types.GetServiceAddress(1).String(),
				))
			},
			expBroken: true,
		},
		{
			name: "service with id higher than next service id breaks invariant",
			store: func(ctx sdk.Context) {
				suite.k.SetNextServiceID(ctx, 1)
				suite.k.SaveService(ctx, types.NewService(
					2,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					types.GetServiceAddress(2).String(),
				))
			},
			expBroken: true,
		},
		{
			name: "invalid service breaks invariant",
			store: func(ctx sdk.Context) {
				suite.k.SetNextServiceID(ctx, 1)
				suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_UNSPECIFIED,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					types.GetServiceAddress(1).String(),
				))
			},
			expBroken: true,
		},
		{
			name: "valid data does not break invariant",
			store: func(ctx sdk.Context) {
				suite.k.SetNextServiceID(ctx, 2)
				suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					types.GetServiceAddress(1).String(),
				))
			},
			expBroken: false,
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

			_, broken := keeper.ValidServicesInvariant(suite.k)(ctx)
			suite.Require().Equal(tc.expBroken, broken)
		})
	}
}
