package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/operators/keeper"
	"github.com/milkyway-labs/milkyway/x/operators/types"
)

func (suite *KeeperTestSuite) TestValidOperatorsInvariant() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		expBroken bool
	}{
		{
			name: "not found next operator id breaks invariant",
			store: func(ctx sdk.Context) {
				err := suite.k.RegisterOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					types.DefaultOperatorParams(),
				))
				suite.Require().NoError(err)

				err = suite.k.RegisterOperator(ctx, types.NewOperator(
					2,
					types.OPERATOR_STATUS_INACTIVATING,
					"Inertia",
					"https://inertia.zone",
					"",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					types.DefaultOperatorParams(),
				))
				suite.Require().NoError(err)
			},
			expBroken: true,
		},
		{
			name: "operator with id equals to next operator id breaks invariant",
			store: func(ctx sdk.Context) {
				suite.k.SetNextOperatorID(ctx, 1)
				err := suite.k.RegisterOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					types.DefaultOperatorParams(),
				))
				suite.Require().NoError(err)
			},
			expBroken: true,
		},
		{
			name: "operator with id higher than next operator id breaks invariant",
			store: func(ctx sdk.Context) {
				suite.k.SetNextOperatorID(ctx, 1)
				err := suite.k.RegisterOperator(ctx, types.NewOperator(
					2,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					types.DefaultOperatorParams(),
				))
				suite.Require().NoError(err)
			},
			expBroken: true,
		},
		{
			name: "invalid operator breaks invariant",
			store: func(ctx sdk.Context) {
				suite.k.SetNextOperatorID(ctx, 1)
				err := suite.k.RegisterOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_UNSPECIFIED,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					types.DefaultOperatorParams(),
				))
				suite.Require().NoError(err)
			},
			expBroken: true,
		},
		{
			name: "valid data does not break invariant",
			store: func(ctx sdk.Context) {
				suite.k.SetNextOperatorID(ctx, 3)
				err := suite.k.RegisterOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					types.DefaultOperatorParams(),
				))
				suite.Require().NoError(err)

				err = suite.k.RegisterOperator(ctx, types.NewOperator(
					2,
					types.OPERATOR_STATUS_INACTIVATING,
					"Inertia",
					"https://inertia.zone",
					"",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					types.DefaultOperatorParams(),
				))
				suite.Require().NoError(err)
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

			_, broken := keeper.ValidOperatorsInvariant(suite.k)(ctx)
			suite.Require().Equal(tc.expBroken, broken)
		})
	}
}
