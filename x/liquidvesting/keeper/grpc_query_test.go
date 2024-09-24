package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/keeper"
	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
)

func (suite *KeeperTestSuite) TestQuerier_InsuranceFund() {
	user1 := "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"
	user2 := "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"

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
				suite.fundAccountInsuranceFund(ctx, user1, sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 1000)))
			},
			expBalance: sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 1000)),
		},
		{
			name: "multiple deposits",
			setup: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx, user1, sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 1000)))
				suite.fundAccountInsuranceFund(ctx, user2, sdk.NewCoins(
					sdk.NewInt64Coin(iBCDenom, 1000), sdk.NewInt64Coin("stake", 1000)))
			},
			expBalance: sdk.NewCoins(
				sdk.NewInt64Coin(iBCDenom, 2000),
				sdk.NewInt64Coin("stake", 1000),
			),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, _ := suite.ctx.CacheContext()

			if tc.setup != nil {
				tc.setup(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			resp, err := querier.InsuranceFund(ctx, types.NewQueryInsuranceFundRequest())
			suite.Assert().NoError(err)
			suite.Assert().Equal(tc.expBalance, resp.Amount)
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_UserInsuranceFund() {
	user1 := "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"
	user2 := "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"

	testCases := []struct {
		name       string
		setup      func(ctx sdk.Context)
		shouldErr  bool
		request    *types.QueryUserInsuranceFundRequest
		expBalance sdk.Coins
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
				suite.fundAccountInsuranceFund(ctx, user1, sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 1000)))
			},
			request:    types.NewQueryUserInsuranceFundRequest(user1),
			shouldErr:  false,
			expBalance: sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 1000)),
		},
		{
			name: "multiple deposits",
			setup: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx, user1, sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 1000)))
				suite.fundAccountInsuranceFund(ctx, user2, sdk.NewCoins(
					sdk.NewInt64Coin(iBCDenom, 1000), sdk.NewInt64Coin("stake", 1000)))
			},
			request:   types.NewQueryUserInsuranceFundRequest(user2),
			shouldErr: false,
			expBalance: sdk.NewCoins(
				sdk.NewInt64Coin(iBCDenom, 1000),
				sdk.NewInt64Coin("stake", 1000),
			),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, _ := suite.ctx.CacheContext()

			if tc.setup != nil {
				tc.setup(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			resp, err := querier.UserInsuranceFund(ctx, tc.request)
			if tc.shouldErr {
				suite.Assert().Error(err)
			} else {
				suite.Assert().NoError(err)
				suite.Assert().Equal(tc.expBalance, resp.Amount)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_UserInsuranceFunds() {
	user1 := "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"
	user2 := "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"

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
				suite.fundAccountInsuranceFund(ctx, user1, sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 1000)))
				suite.fundAccountInsuranceFund(ctx, user2, sdk.NewCoins(
					sdk.NewInt64Coin(iBCDenom, 1000), sdk.NewInt64Coin("stake", 1000)))
			},
			request:   types.NewQueryUserInsuranceFundsRequest(nil),
			shouldErr: false,
			expInsuranceFunds: []types.UserInsuranceFundData{
				types.NewUserInsuranceFundData(user1, sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 1000))),
				types.NewUserInsuranceFundData(user2, sdk.NewCoins(
					sdk.NewInt64Coin(iBCDenom, 1000), sdk.NewInt64Coin("stake", 1000))),
			},
		},
		{
			name: "respects handle pagination",
			setup: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx, user1, sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 1000)))
				suite.fundAccountInsuranceFund(ctx, user2, sdk.NewCoins(
					sdk.NewInt64Coin(iBCDenom, 1000), sdk.NewInt64Coin("stake", 1000)))
			},
			request: types.NewQueryUserInsuranceFundsRequest(&query.PageRequest{
				Offset: 0,
				Limit:  1,
			}),
			shouldErr: false,
			expInsuranceFunds: []types.UserInsuranceFundData{
				types.NewUserInsuranceFundData(user1, sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 1000))),
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, _ := suite.ctx.CacheContext()

			if tc.setup != nil {
				tc.setup(ctx)
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
	user1 := "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"

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
				suite.fundAccountInsuranceFund(ctx, user1, sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 1)))
			},
			request:    types.NewQueryUserRestakableAssetsRequest(user1),
			shouldErr:  false,
			expBalance: sdk.NewCoins(sdk.NewInt64Coin(vestedDenom, 100)),
		},
		{
			name: "5% insurance fund",
			setup: func(ctx sdk.Context) {
				suite.Assert().NoError(suite.k.SetParams(ctx, types.NewParams(
					math.LegacyMustNewDecFromStr("5"), nil, nil,
				)))
				suite.fundAccountInsuranceFund(ctx, user1, sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 1)))
			},
			request:    types.NewQueryUserRestakableAssetsRequest(user1),
			shouldErr:  false,
			expBalance: sdk.NewCoins(sdk.NewInt64Coin(vestedDenom, 20)),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, _ := suite.ctx.CacheContext()

			if tc.setup != nil {
				tc.setup(ctx)
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
