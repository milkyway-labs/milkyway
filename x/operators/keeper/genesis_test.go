package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

func (suite *KeeperTestSuite) TestKeeper_ExportGenesis() {
	testCases := []struct {
		name       string
		setup      func()
		setupCtx   func(ctx sdk.Context) sdk.Context
		store      func(ctx sdk.Context)
		expGenesis *types.GenesisState
	}{
		{
			name: "next operator id is exported properly",
			store: func(ctx sdk.Context) {
				suite.k.SetNextOperatorID(ctx, 10)
				suite.k.SetParams(ctx, types.DefaultParams())
			},
			expGenesis: &types.GenesisState{
				NextOperatorID: 10,
				Operators:      nil,
				Params:         types.DefaultParams(),
			},
		},
		{
			name: "operators data are exported properly",
			store: func(ctx sdk.Context) {
				suite.k.SetNextOperatorID(ctx, 10)
				suite.k.SetParams(ctx, types.DefaultParams())

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
			expGenesis: &types.GenesisState{
				NextOperatorID: 10,
				Operators: []types.Operator{
					types.NewOperator(
						1,
						types.OPERATOR_STATUS_ACTIVE,
						"MilkyWay Operator",
						"https://milkyway.com",
						"https://milkyway.com/picture",
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						types.DefaultOperatorParams(),
					),
					types.NewOperator(
						2,
						types.OPERATOR_STATUS_INACTIVATING,
						"Inertia",
						"https://inertia.zone",
						"",
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						types.DefaultOperatorParams(),
					),
				},
				Params: types.DefaultParams(),
			},
		},
		{
			name: "inactivating operators data are exported properly",
			setupCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
			},
			store: func(ctx sdk.Context) {
				suite.k.SetNextOperatorID(ctx, 10)
				suite.k.SetParams(ctx, types.DefaultParams())

				activeValidator := types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					types.DefaultOperatorParams(),
				)
				err := suite.k.RegisterOperator(ctx, activeValidator)
				suite.Require().NoError(err)

				suite.k.StartOperatorInactivation(ctx, activeValidator)

				err = suite.k.RegisterOperator(ctx, types.NewOperator(
					2,
					types.OPERATOR_STATUS_ACTIVE,
					"Inertia",
					"https://inertia.zone",
					"",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					types.DefaultOperatorParams(),
				))
				suite.Require().NoError(err)
			},
			expGenesis: &types.GenesisState{
				NextOperatorID: 10,
				Operators: []types.Operator{
					types.NewOperator(
						1,
						types.OPERATOR_STATUS_INACTIVATING,
						"MilkyWay Operator",
						"https://milkyway.com",
						"https://milkyway.com/picture",
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						types.DefaultOperatorParams(),
					),
					types.NewOperator(
						2,
						types.OPERATOR_STATUS_ACTIVE,
						"Inertia",
						"https://inertia.zone",
						"",
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						types.DefaultOperatorParams(),
					),
				},
				UnbondingOperators: []types.UnbondingOperator{
					types.NewUnbondingOperator(
						1,
						time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC).Add(types.DefaultParams().DeactivationTime),
					),
				},
				Params: types.DefaultParams(),
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
			if tc.setupCtx != nil {
				ctx = tc.setupCtx(ctx)
			}
			if tc.store != nil {
				tc.store(ctx)
			}

			genesis, err := suite.k.ExportGenesis(ctx)
			suite.Require().NoError(err)
			suite.Require().Equal(tc.expGenesis, genesis)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_InitGenesis() {
	testCases := []struct {
		name     string
		setup    func()
		setupCtx func(ctx sdk.Context)
		store    func(ctx sdk.Context)
		genesis  *types.GenesisState
		check    func(ctx sdk.Context)
	}{
		{
			name: "next operator id is initialized properly",
			genesis: &types.GenesisState{
				NextOperatorID: 10,
				Params:         types.DefaultParams(),
			},
			check: func(ctx sdk.Context) {
				nextOperatorID, err := suite.k.GetNextOperatorID(ctx)
				suite.Require().NoError(err)
				suite.Require().Equal(uint32(10), nextOperatorID)
			},
		},
		{
			name: "operators are initialized properly",
			genesis: &types.GenesisState{
				Operators: []types.Operator{
					types.NewOperator(
						1,
						types.OPERATOR_STATUS_ACTIVE,
						"MilkyWay Operator",
						"https://milkyway.com",
						"https://milkyway.com/picture",
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						types.DefaultOperatorParams(),
					),
					types.NewOperator(
						2,
						types.OPERATOR_STATUS_INACTIVATING,
						"Inertia",
						"https://inertia.zone",
						"",
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						types.DefaultOperatorParams(),
					),
				},
				Params: types.DefaultParams(),
			},
			check: func(ctx sdk.Context) {
				operator1, found := suite.k.GetOperator(ctx, 1)
				suite.Require().True(found)
				suite.Require().Equal(types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					types.DefaultOperatorParams(),
				), operator1)

				operator2, found := suite.k.GetOperator(ctx, 2)
				suite.Require().True(found)
				suite.Require().Equal(types.NewOperator(
					2,
					types.OPERATOR_STATUS_INACTIVATING,
					"Inertia",
					"https://inertia.zone",
					"",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					types.DefaultOperatorParams(),
				), operator2)
			},
		},
		{
			name: "unbonding operators are initialized properly",
			genesis: &types.GenesisState{
				UnbondingOperators: []types.UnbondingOperator{
					types.NewUnbondingOperator(
						1,
						time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
					),
					types.NewUnbondingOperator(
						2,
						time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
					),
				},
			},
			check: func(ctx sdk.Context) {
				inactivatingOperators, _ := suite.k.GetInactivatingOperators(ctx)
				suite.Require().Equal([]types.UnbondingOperator{
					types.NewUnbondingOperator(
						1,
						time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
					),
					types.NewUnbondingOperator(
						2,
						time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
					),
				}, inactivatingOperators)
			},
		},
		{
			name: "params are initialized properly",
			genesis: &types.GenesisState{
				Params: types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000))),
					10*time.Hour,
				),
			},
			check: func(ctx sdk.Context) {
				params := suite.k.GetParams(ctx)
				suite.Require().Equal(types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000))),
					10*time.Hour,
				), params)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			if tc.setup != nil {
				tc.setup()
			}
			ctx, _ := suite.ctx.CacheContext()
			if tc.setupCtx != nil {
				tc.setupCtx(ctx)
			}
			if tc.store != nil {
				tc.store(ctx)
			}

			suite.k.InitGenesis(ctx, *tc.genesis)

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}
