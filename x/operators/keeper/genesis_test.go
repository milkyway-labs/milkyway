package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v5/x/operators/types"
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
				err := suite.k.SetNextOperatorID(ctx, 10)
				suite.Require().NoError(err)

				err = suite.k.SetParams(ctx, types.DefaultParams())
				suite.Require().NoError(err)
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
				err := suite.k.SetNextOperatorID(ctx, 10)
				suite.Require().NoError(err)

				err = suite.k.SetParams(ctx, types.DefaultParams())
				suite.Require().NoError(err)

				err = suite.k.CreateOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)

				err = suite.k.CreateOperator(ctx, types.NewOperator(
					2,
					types.OPERATOR_STATUS_INACTIVATING,
					"Inertia",
					"https://inertia.zone",
					"",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
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
					),
					types.NewOperator(
						2,
						types.OPERATOR_STATUS_INACTIVATING,
						"Inertia",
						"https://inertia.zone",
						"",
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
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
				err := suite.k.SetNextOperatorID(ctx, 10)
				suite.Require().NoError(err)

				err = suite.k.SetParams(ctx, types.DefaultParams())
				suite.Require().NoError(err)

				activeValidator := types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				)
				err = suite.k.CreateOperator(ctx, activeValidator)
				suite.Require().NoError(err)

				err = suite.k.StartOperatorInactivation(ctx, activeValidator)
				suite.Require().NoError(err)

				err = suite.k.CreateOperator(ctx, types.NewOperator(
					2,
					types.OPERATOR_STATUS_ACTIVE,
					"Inertia",
					"https://inertia.zone",
					"",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
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
					),
					types.NewOperator(
						2,
						types.OPERATOR_STATUS_ACTIVE,
						"Inertia",
						"https://inertia.zone",
						"",
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
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
		{
			name: "operators params are exported properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetNextOperatorID(ctx, 10)
				suite.Require().NoError(err)

				err = suite.k.SetParams(ctx, types.DefaultParams())
				suite.Require().NoError(err)

				err = suite.k.CreateOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)

				// Set some custom params
				err = suite.k.SaveOperatorParams(ctx, 1, types.NewOperatorParams(
					sdkmath.LegacyMustNewDecFromStr("0.1"),
				))
				suite.Require().NoError(err)

				err = suite.k.CreateOperator(ctx, types.NewOperator(
					2,
					types.OPERATOR_STATUS_INACTIVATING,
					"Inertia",
					"https://inertia.zone",
					"",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				))
				suite.Require().NoError(err)

				// Set some custom params
				err = suite.k.SaveOperatorParams(ctx, 2, types.NewOperatorParams(
					sdkmath.LegacyMustNewDecFromStr("0.2"),
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
					),
					types.NewOperator(
						2,
						types.OPERATOR_STATUS_INACTIVATING,
						"Inertia",
						"https://inertia.zone",
						"",
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					),
				},
				OperatorsParams: []types.OperatorParamsRecord{
					types.NewOperatorParamsRecord(1, types.NewOperatorParams(
						sdkmath.LegacyMustNewDecFromStr("0.1"),
					)),
					types.NewOperatorParamsRecord(2, types.NewOperatorParams(
						sdkmath.LegacyMustNewDecFromStr("0.2"),
					)),
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

			genesis := suite.k.ExportGenesis(ctx)
			suite.Require().Equal(tc.expGenesis, genesis)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_InitGenesis() {
	testCases := []struct {
		name      string
		setup     func()
		setupCtx  func(ctx sdk.Context)
		store     func(ctx sdk.Context)
		genesis   *types.GenesisState
		check     func(ctx sdk.Context)
		shouldErr bool
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
					),
					types.NewOperator(
						2,
						types.OPERATOR_STATUS_INACTIVATING,
						"Inertia",
						"https://inertia.zone",
						"",
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					),
				},
				Params: types.DefaultParams(),
			},
			check: func(ctx sdk.Context) {
				operator, err := suite.k.GetOperator(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				), operator)

				// Ensure that the operator has the default params
				params, err := suite.k.GetOperatorParams(ctx, 1)
				suite.Require().Equal(types.DefaultOperatorParams(), params)

				operator, err = suite.k.GetOperator(ctx, 2)
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewOperator(
					2,
					types.OPERATOR_STATUS_INACTIVATING,
					"Inertia",
					"https://inertia.zone",
					"",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				), operator)

				// Ensure that the operator has the default params
				params, err = suite.k.GetOperatorParams(ctx, 2)
				suite.Require().NoError(err)
				suite.Require().Equal(types.DefaultOperatorParams(), params)
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
				params, err := suite.k.GetParams(ctx)
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000))),
					10*time.Hour,
				), params)
			},
		},
		{
			name: "operator params without associated operator should err",
			genesis: &types.GenesisState{
				Operators: []types.Operator{
					types.NewOperator(
						1,
						types.OPERATOR_STATUS_ACTIVE,
						"MilkyWay Operator",
						"https://milkyway.com",
						"https://milkyway.com/picture",
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					),
				},
				OperatorsParams: []types.OperatorParamsRecord{
					types.NewOperatorParamsRecord(2, types.DefaultOperatorParams()),
				},
			},
			shouldErr: true,
		},
		{
			name: "operator params are initialized properly",
			genesis: &types.GenesisState{
				Operators: []types.Operator{
					types.NewOperator(
						1,
						types.OPERATOR_STATUS_ACTIVE,
						"MilkyWay Operator",
						"https://milkyway.com",
						"https://milkyway.com/picture",
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					),
				},
				OperatorsParams: []types.OperatorParamsRecord{
					types.NewOperatorParamsRecord(1, types.NewOperatorParams(
						sdkmath.LegacyMustNewDecFromStr("0.2"),
					)),
				},
			},
			check: func(ctx sdk.Context) {
				params, err := suite.k.GetOperatorParams(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().Equal(params, types.NewOperatorParams(
					sdkmath.LegacyMustNewDecFromStr("0.2"),
				))
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

			err := suite.k.InitGenesis(ctx, tc.genesis)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				if tc.check != nil {
					tc.check(ctx)
				}
			}
		})
	}
}
