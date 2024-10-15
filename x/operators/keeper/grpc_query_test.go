package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

func (suite *KeeperTestSuite) TestQueryServer_Operator() {
	testCases := []struct {
		name        string
		store       func(ctx sdk.Context)
		request     *types.QueryOperatorRequest
		shouldErr   bool
		expResponse *types.QueryOperatorResponse
	}{
		{
			name:      "not found operator returns error",
			request:   types.NewQueryOperatorRequest(1),
			shouldErr: true,
		},
		{
			name: "existing operator is returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.RegisterOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)
			},
			request:   types.NewQueryOperatorRequest(1),
			shouldErr: false,
			expResponse: &types.QueryOperatorResponse{
				Operator: types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			if tc.store != nil {
				tc.store(suite.ctx)
			}

			res, err := suite.k.Operator(sdk.WrapSDKContext(suite.ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expResponse, res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQueryServer_OperatorParams() {
	testCases := []struct {
		name        string
		store       func(ctx sdk.Context)
		request     *types.QueryOperatorParamsRequest
		shouldErr   bool
		expResponse *types.QueryOperatorParamsResponse
	}{
		{
			name:      "not found operator returns error",
			request:   types.NewQueryOperatorParamsRequest(1),
			shouldErr: true,
		},
		{
			name: "default operator params are returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.RegisterOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)
			},
			request:   types.NewQueryOperatorParamsRequest(1),
			shouldErr: false,
			expResponse: &types.QueryOperatorParamsResponse{
				OperatorParams: types.DefaultOperatorParams(),
			},
		},
		{
			name: "updated operator params are returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.RegisterOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)

				err = suite.k.SaveOperatorParams(ctx, 1, types.NewOperatorParams(
					sdkmath.LegacyMustNewDecFromStr("0.2"),
				))
			},
			request:   types.NewQueryOperatorParamsRequest(1),
			shouldErr: false,
			expResponse: &types.QueryOperatorParamsResponse{
				OperatorParams: types.NewOperatorParams(
					sdkmath.LegacyMustNewDecFromStr("0.2"),
				),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			if tc.store != nil {
				tc.store(suite.ctx)
			}

			res, err := suite.k.OperatorParams(suite.ctx, tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expResponse, res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQueryServer_Operators() {
	testCases := []struct {
		name         string
		store        func(ctx sdk.Context)
		request      *types.QueryOperatorsRequest
		shouldErr    bool
		expOperators []types.Operator
	}{
		{
			name: "query without pagination returns data properly",
			store: func(ctx sdk.Context) {
				err := suite.k.RegisterOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)

				err = suite.k.RegisterOperator(ctx, types.NewOperator(
					2,
					types.OPERATOR_STATUS_INACTIVATING,
					"Inertia",
					"https://inertia.zone",
					"",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				))
				suite.Require().NoError(err)
			},
			request:   types.NewQueryOperatorsRequest(nil),
			shouldErr: false,
			expOperators: []types.Operator{
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
		},
		{
			name: "query with pagination returns data properly",
			store: func(ctx sdk.Context) {
				err := suite.k.RegisterOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)

				err = suite.k.RegisterOperator(ctx, types.NewOperator(
					2,
					types.OPERATOR_STATUS_INACTIVATING,
					"Inertia",
					"https://inertia.zone",
					"",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				))
				suite.Require().NoError(err)
			},
			request: types.NewQueryOperatorsRequest(&query.PageRequest{
				Offset: 1,
				Limit:  1,
			}),
			shouldErr: false,
			expOperators: []types.Operator{
				types.NewOperator(
					2,
					types.OPERATOR_STATUS_INACTIVATING,
					"Inertia",
					"https://inertia.zone",
					"",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
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

			res, err := suite.k.Operators(ctx, tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expOperators, res.Operators)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQueryServer_Params() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		request   *types.QueryParamsRequest
		shouldErr bool
		expParams types.Params
	}{
		{
			name: "existing params are returned properly",
			store: func(ctx sdk.Context) {
				suite.k.SetParams(ctx, types.DefaultParams())
			},
			request:   types.NewQueryParamsRequest(),
			shouldErr: false,
			expParams: types.DefaultParams(),
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
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expParams, res.Params)
			}
		})
	}
}
