package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/keeper"
	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
)

func (suite *KeeperTestSuite) TestQuerier_InsuranceFund() {
	testCases := []struct {
		name       string
		setup      func(ctx sdk.Context)
		expBalance sdk.Coins
	}{
		{
			name:       "empty insurance fund",
			expBalance: sdk.NewCoins(),
		},
		{
			name: "single deposit",
			setup: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
				)
			},
			expBalance: sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
		},
		{
			name: "multiple deposits",
			setup: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
				)
				suite.fundAccountInsuranceFund(ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000), sdk.NewInt64Coin("stake", 1000)),
				)
			},
			expBalance: sdk.NewCoins(
				sdk.NewInt64Coin(IBCDenom, 2000),
				sdk.NewInt64Coin("stake", 1000),
			),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			if tc.setup != nil {
				tc.setup(suite.ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			resp, err := querier.InsuranceFund(suite.ctx, types.NewQueryInsuranceFundRequest())
			suite.Assert().NoError(err)
			suite.Assert().Equal(tc.expBalance, resp.Amount)
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_UserInsuranceFund() {
	testCases := []struct {
		name       string
		setup      func(ctx sdk.Context)
		shouldErr  bool
		request    *types.QueryUserInsuranceFundRequest
		expBalance sdk.Coins
		expUsed    sdk.Coins
	}{
		{
			name:      "empty request",
			request:   nil,
			shouldErr: true,
		},
		{
			name:      "invalid address",
			request:   types.NewQueryUserInsuranceFundRequest("invalid"),
			shouldErr: true,
		},
		{
			name: "single deposit",
			setup: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
				)
			},
			request:    types.NewQueryUserInsuranceFundRequest("cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"),
			shouldErr:  false,
			expBalance: sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
		},
		{
			name: "multiple deposits",
			setup: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
				)
				suite.fundAccountInsuranceFund(ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000), sdk.NewInt64Coin("stake", 1000)),
				)
			},
			request:   types.NewQueryUserInsuranceFundRequest("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			shouldErr: false,
			expBalance: sdk.NewCoins(
				sdk.NewInt64Coin(IBCDenom, 1000),
				sdk.NewInt64Coin("stake", 1000),
			),
		},
		{
			name: "with used amount",
			setup: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
				)
				suite.mintVestedRepresentation(
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
				)

				// Add other tokens
				suite.fundAccountInsuranceFund(ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000), sdk.NewInt64Coin("stake", 1000)),
				)

				// Delegate to the pool
				suite.createPool(1, vestedIBCDenom)
				_, err := suite.rk.DelegateToPool(ctx, sdk.NewInt64Coin(vestedIBCDenom, 1000), "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre")
				suite.Require().NoError(err)
			},
			request:   types.NewQueryUserInsuranceFundRequest("cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"),
			shouldErr: false,
			expBalance: sdk.NewCoins(
				sdk.NewInt64Coin(IBCDenom, 1000),
			),
			expUsed: sdk.NewCoins(
				sdk.NewInt64Coin(IBCDenom, 20),
			),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			if tc.setup != nil {
				tc.setup(suite.ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			resp, err := querier.UserInsuranceFund(suite.ctx, tc.request)
			if tc.shouldErr {
				suite.Assert().Error(err)
			} else {
				suite.Assert().NoError(err)
				suite.Assert().Equal(tc.expBalance, resp.Balance)
				suite.Assert().Equal(tc.expUsed, resp.Used)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_UserInsuranceFunds() {
	testCases := []struct {
		name              string
		setup             func(ctx sdk.Context)
		shouldErr         bool
		request           *types.QueryUserInsuranceFundsRequest
		expInsuranceFunds []types.UserInsuranceFundData
	}{
		{
			name:      "empty request",
			request:   nil,
			shouldErr: true,
		},
		{
			name: "no pagination",
			setup: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
				)
				suite.fundAccountInsuranceFund(ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000), sdk.NewInt64Coin("stake", 1000)),
				)
			},
			request:   types.NewQueryUserInsuranceFundsRequest(nil),
			shouldErr: false,
			expInsuranceFunds: []types.UserInsuranceFundData{
				types.NewUserInsuranceFundData(
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					types.NewInsuranceFund(
						sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
						nil,
					),
				),
				types.NewUserInsuranceFundData(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					types.NewInsuranceFund(
						sdk.NewCoins(
							sdk.NewInt64Coin(IBCDenom, 1000),
							sdk.NewInt64Coin("stake", 1000),
						),
						nil,
					),
				),
			},
		},
		{
			name: "respects handle pagination",
			setup: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
				)
				suite.fundAccountInsuranceFund(ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000), sdk.NewInt64Coin("stake", 1000)),
				)
			},
			request: types.NewQueryUserInsuranceFundsRequest(&query.PageRequest{
				Offset: 0,
				Limit:  1,
			}),
			shouldErr: false,
			expInsuranceFunds: []types.UserInsuranceFundData{
				types.NewUserInsuranceFundData(
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					types.NewInsuranceFund(
						sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
						nil,
					),
				),
			},
		},
		{
			name: "with utilization",
			setup: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
				)
				suite.mintVestedRepresentation(
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
				)

				// Delegate to the pool
				suite.createPool(1, vestedIBCDenom)
				_, err := suite.rk.DelegateToPool(ctx, sdk.NewInt64Coin(vestedIBCDenom, 1000), "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre")
				suite.Require().NoError(err)
			},
			request: types.NewQueryUserInsuranceFundsRequest(&query.PageRequest{
				Offset: 0,
				Limit:  1,
			}),
			shouldErr: false,
			expInsuranceFunds: []types.UserInsuranceFundData{
				types.NewUserInsuranceFundData(
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					types.NewInsuranceFund(
						sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
						sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 20)),
					),
				),
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			if tc.setup != nil {
				tc.setup(suite.ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			resp, err := querier.UserInsuranceFunds(suite.ctx, tc.request)
			if tc.shouldErr {
				suite.Assert().Error(err)
			} else {
				suite.Assert().NoError(err)
				suite.Assert().Equal(tc.expInsuranceFunds, resp.InsuranceFunds)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_UserRestakableAssets() {
	testCases := []struct {
		name       string
		setup      func(ctx sdk.Context)
		shouldErr  bool
		request    *types.QueryUserRestakableAssetsRequest
		expBalance sdk.Coins
	}{
		{
			name:      "empty request",
			request:   nil,
			shouldErr: true,
		},
		{
			name:      "invalid address",
			request:   types.NewQueryUserRestakableAssetsRequest("invalid"),
			shouldErr: true,
		},
		{
			name: "1% insurance fund",
			setup: func(ctx sdk.Context) {
				suite.Assert().NoError(suite.k.SetParams(ctx, types.NewParams(
					math.LegacyMustNewDecFromStr("1"), nil, nil,
				)))
				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1)),
				)
			},
			request:    types.NewQueryUserRestakableAssetsRequest("cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"),
			shouldErr:  false,
			expBalance: sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 100)),
		},
		{
			name: "5% insurance fund",
			setup: func(ctx sdk.Context) {
				suite.Assert().NoError(suite.k.SetParams(ctx, types.NewParams(
					math.LegacyMustNewDecFromStr("5"), nil, nil,
				)))
				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1)),
				)
			},
			request:    types.NewQueryUserRestakableAssetsRequest("cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"),
			shouldErr:  false,
			expBalance: sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 20)),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			if tc.setup != nil {
				tc.setup(suite.ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			resp, err := querier.UserRestakableAssets(suite.ctx, tc.request)
			if tc.shouldErr {
				suite.Assert().Error(err)
			} else {
				suite.Assert().NoError(err)
				suite.Assert().Equal(tc.expBalance, resp.Amount)
			}
		})
	}
}
