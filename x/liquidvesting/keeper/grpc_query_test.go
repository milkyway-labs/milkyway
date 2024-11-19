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
				suite.fundAccountInsuranceFund(ctx, user1, sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)))
			},
			expBalance: sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
		},
		{
			name: "multiple deposits",
			setup: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx, user1, sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)))
				suite.fundAccountInsuranceFund(ctx, user2, sdk.NewCoins(
					sdk.NewInt64Coin(IBCDenom, 1000), sdk.NewInt64Coin("stake", 1000)))
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
	user1 := "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"
	user2 := "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"

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
				suite.fundAccountInsuranceFund(ctx, user1, sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)))
			},
			request:    types.NewQueryUserInsuranceFundRequest(user1),
			shouldErr:  false,
			expBalance: sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
			expUsed:    sdk.NewCoins(),
		},
		{
			name: "multiple deposits",
			setup: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx, user1, sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)))
				suite.fundAccountInsuranceFund(ctx, user2, sdk.NewCoins(
					sdk.NewInt64Coin(IBCDenom, 1000), sdk.NewInt64Coin("stake", 1000)))
			},
			request:   types.NewQueryUserInsuranceFundRequest(user2),
			shouldErr: false,
			expBalance: sdk.NewCoins(
				sdk.NewInt64Coin(IBCDenom, 1000),
				sdk.NewInt64Coin("stake", 1000),
			),
			expUsed: sdk.NewCoins(),
		},
		{
			name: "with used amount",
			setup: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx, user1, sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)))
				suite.mintVestedRepresentation(user1, sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)))

				// Add other tokens
				suite.fundAccountInsuranceFund(ctx, user2, sdk.NewCoins(
					sdk.NewInt64Coin(IBCDenom, 1000), sdk.NewInt64Coin("stake", 1000)))

				// Delegate to the pool
				suite.createPool(1, vestedIBCDenom)
				_, err := suite.rk.DelegateToPool(ctx, sdk.NewInt64Coin(vestedIBCDenom, 1000), user1)
				suite.Require().NoError(err)
			},
			request:   types.NewQueryUserInsuranceFundRequest(user1),
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
				suite.fundAccountInsuranceFund(ctx, user1, sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)))
				suite.fundAccountInsuranceFund(ctx, user2, sdk.NewCoins(
					sdk.NewInt64Coin(IBCDenom, 1000), sdk.NewInt64Coin("stake", 1000)))
			},
			request:   types.NewQueryUserInsuranceFundsRequest(nil),
			shouldErr: false,
			expInsuranceFunds: []types.UserInsuranceFundData{
				types.NewUserInsuranceFundData(user1, sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)), sdk.NewCoins()),
				types.NewUserInsuranceFundData(user2,
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000), sdk.NewInt64Coin("stake", 1000)), sdk.NewCoins()),
			},
		},
		{
			name: "respects handle pagination",
			setup: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx, user1, sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)))
				suite.fundAccountInsuranceFund(ctx, user2, sdk.NewCoins(
					sdk.NewInt64Coin(IBCDenom, 1000), sdk.NewInt64Coin("stake", 1000)))
			},
			request: types.NewQueryUserInsuranceFundsRequest(&query.PageRequest{
				Offset: 0,
				Limit:  1,
			}),
			shouldErr: false,
			expInsuranceFunds: []types.UserInsuranceFundData{
				types.NewUserInsuranceFundData(user1, sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)), sdk.NewCoins()),
			},
		},
		{
			name: "with utilization",
			setup: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx, user1, sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)))
				suite.mintVestedRepresentation(user1, sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)))

				// Delegate to the pool
				suite.createPool(1, vestedIBCDenom)
				_, err := suite.rk.DelegateToPool(ctx, sdk.NewInt64Coin(vestedIBCDenom, 1000), user1)
				suite.Require().NoError(err)
			},
			request: types.NewQueryUserInsuranceFundsRequest(&query.PageRequest{
				Offset: 0,
				Limit:  1,
			}),
			shouldErr: false,
			expInsuranceFunds: []types.UserInsuranceFundData{
				types.NewUserInsuranceFundData(user1,
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)), sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 20))),
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
					math.LegacyMustNewDecFromStr("1"), nil, nil, nil,
				)))
				suite.fundAccountInsuranceFund(ctx, user1, sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1)))
			},
			request:    types.NewQueryUserRestakableAssetsRequest(user1),
			shouldErr:  false,
			expBalance: sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 100)),
		},
		{
			name: "5% insurance fund",
			setup: func(ctx sdk.Context) {
				suite.Assert().NoError(suite.k.SetParams(ctx, types.NewParams(
					math.LegacyMustNewDecFromStr("5"), nil, nil, nil,
				)))
				suite.fundAccountInsuranceFund(ctx, user1, sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1)))
			},
			request:    types.NewQueryUserRestakableAssetsRequest(user1),
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
