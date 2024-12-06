package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/milkyway-labs/milkyway/v3/x/liquidvesting/keeper"
	"github.com/milkyway-labs/milkyway/v3/x/liquidvesting/types"
)

func (suite *KeeperTestSuite) TestQuerier_InsuranceFund() {
	testCases := []struct {
		name       string
		store      func(ctx sdk.Context)
		expBalance sdk.Coins
	}{
		{
			name:       "empty insurance fund",
			expBalance: sdk.NewCoins(),
		},
		{
			name: "single deposit",
			store: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
				)
			},
			expBalance: sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
		},
		{
			name: "multiple deposits",
			store: func(ctx sdk.Context) {
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

			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			resp, err := querier.InsuranceFund(ctx, types.NewQueryInsuranceFundRequest())
			suite.Assert().NoError(err)
			suite.Assert().Equal(tc.expBalance, resp.Amount)
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_UserInsuranceFund() {
	testCases := []struct {
		name       string
		store      func(ctx sdk.Context)
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
			name: "single deposit",
			store: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
				)
			},
			request:    types.NewQueryUserInsuranceFundRequest("cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"),
			shouldErr:  false,
			expBalance: sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
			expUsed:    sdk.NewCoins(),
		},
		{
			name: "multiple deposits",
			store: func(ctx sdk.Context) {
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
			expUsed: sdk.NewCoins(),
		},
		{
			name: "with used amount",
			store: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
				)
				suite.mintLockedRepresentation(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
				)

				// Add other tokens
				suite.fundAccountInsuranceFund(ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(
						sdk.NewInt64Coin(IBCDenom, 1000),
						sdk.NewInt64Coin("stake", 1000),
					),
				)

				// Delegate to the pool
				suite.createPool(ctx, 1, LockedIBCDenom)
				_, err := suite.rk.DelegateToPool(ctx,
					sdk.NewInt64Coin(LockedIBCDenom, 1000),
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				)
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

			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			resp, err := querier.UserInsuranceFund(ctx, tc.request)
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
		store             func(ctx sdk.Context)
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
			store: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000), sdk.NewInt64Coin("stake", 1000)),
				)
				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
				)
			},
			request:   types.NewQueryUserInsuranceFundsRequest(nil),
			shouldErr: false,
			expInsuranceFunds: []types.UserInsuranceFundData{
				types.NewUserInsuranceFundData(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(
						sdk.NewInt64Coin(IBCDenom, 1000),
						sdk.NewInt64Coin("stake", 1000),
					),
					sdk.NewCoins(),
				),
				types.NewUserInsuranceFundData(
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
					sdk.NewCoins(),
				),
			},
		},
		{
			name: "respects handle pagination",
			store: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000), sdk.NewInt64Coin("stake", 1000)),
				)
				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
				)
			},
			request: types.NewQueryUserInsuranceFundsRequest(&query.PageRequest{
				Offset: 0,
				Limit:  1,
			}),
			shouldErr: false,
			expInsuranceFunds: []types.UserInsuranceFundData{
				types.NewUserInsuranceFundData(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(
						sdk.NewInt64Coin(IBCDenom, 1000),
						sdk.NewInt64Coin("stake", 1000),
					),
					sdk.NewCoins(),
				),
			},
		},
		{
			name: "with utilization",
			store: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
				)
				suite.mintLockedRepresentation(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
				)

				// Delegate to the pool
				suite.createPool(ctx, 1, LockedIBCDenom)
				_, err := suite.rk.DelegateToPool(ctx,
					sdk.NewInt64Coin(LockedIBCDenom, 1000),
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				)
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
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 20)),
				),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			resp, err := querier.UserInsuranceFunds(ctx, tc.request)
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
		store      func(ctx sdk.Context)
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
			name: "1% insurance fund",
			store: func(ctx sdk.Context) {
				suite.Assert().NoError(suite.k.SetParams(ctx, types.NewParams(
					math.LegacyMustNewDecFromStr("1"), nil, nil, nil, nil,
				)))
				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1)),
				)
			},
			request:    types.NewQueryUserRestakableAssetsRequest("cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"),
			shouldErr:  false,
			expBalance: sdk.NewCoins(sdk.NewInt64Coin(LockedIBCDenom, 100)),
		},
		{
			name: "5% insurance fund",
			store: func(ctx sdk.Context) {
				suite.Assert().NoError(suite.k.SetParams(ctx, types.NewParams(
					math.LegacyMustNewDecFromStr("5"), nil, nil, nil, nil,
				)))
				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1)),
				)
			},
			request:    types.NewQueryUserRestakableAssetsRequest("cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"),
			shouldErr:  false,
			expBalance: sdk.NewCoins(sdk.NewInt64Coin(LockedIBCDenom, 20)),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			resp, err := querier.UserRestakableAssets(ctx, tc.request)
			if tc.shouldErr {
				suite.Assert().Error(err)
			} else {
				suite.Assert().NoError(err)
				suite.Assert().Equal(tc.expBalance, resp.Amount)
			}
		})
	}
}
