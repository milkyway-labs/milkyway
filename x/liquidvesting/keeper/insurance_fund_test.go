package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
)

func (suite *KeeperTestSuite) TestKeeper_AddToInsuranceFund() {
	testCases := []struct {
		name                string
		deposits            map[string]sdk.Coins
		expectedTotalAmount sdk.Coins
	}{
		{
			name: "add multiple amounts",
			deposits: map[string]sdk.Coins{
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn": sdk.NewCoins(sdk.NewInt64Coin("stake", 100)),
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd": sdk.NewCoins(sdk.NewInt64Coin("stake", 200)),
			},
			expectedTotalAmount: sdk.NewCoins(sdk.NewInt64Coin("stake", 300)),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			// Cache the context
			ctx, _ := suite.ctx.CacheContext()

			for address, amount := range tc.deposits {
				// Mint the coins that should be in the module
				err := suite.bk.MintCoins(ctx, types.ModuleName, amount)
				suite.Assert().NoError(err)

				accAddress, err := sdk.AccAddressFromBech32(address)
				suite.Require().NoError(err)
				err = suite.k.AddToUserInsuranceFund(ctx, accAddress, amount)
				suite.Assert().NoError(err)
			}

			for address, expectedAmount := range tc.deposits {
				accAddress, err := sdk.AccAddressFromBech32(address)
				suite.Require().NoError(err)

				amount, err := suite.k.GetUserInsuranceFundBalance(ctx, accAddress)
				suite.Assert().NoError(err)
				suite.Assert().Equal(expectedAmount, amount)
			}

			balance, err := suite.k.GetInsuranceFundBalance(ctx)
			suite.Assert().NoError(err)
			suite.Assert().Equal(tc.expectedTotalAmount, balance)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_WithdrawFromInsuranceFund() {
	testCases := []struct {
		name      string
		from      string
		amount    sdk.Coins
		store     func(ctx sdk.Context)
		shouldErr bool
	}{
		{
			name: "withdraw more then deposited",
			from: "cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
			amount: sdk.NewCoins(
				sdk.NewInt64Coin("stake", 100),
				sdk.NewInt64Coin("stake2", 50),
			),
			store: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewCoins(
						sdk.NewInt64Coin("stake", 50),
						sdk.NewInt64Coin("stake2", 50),
					))
			},
			shouldErr: true,
		},
		{
			name: "withdraw correctly",
			from: "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			amount: sdk.NewCoins(
				sdk.NewInt64Coin("stake", 200),
				sdk.NewInt64Coin("stake2", 100),
			),
			store: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewCoins(
						sdk.NewInt64Coin("stake", 200),
						sdk.NewInt64Coin("stake2", 100),
					),
				)
			},
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			userAddr, err := sdk.AccAddressFromBech32(tc.from)
			suite.Require().NoError(err)

			err = suite.k.WithdrawFromUserInsuranceFund(ctx, userAddr, tc.amount)
			if tc.shouldErr {
				suite.Assert().Error(err)
			} else {
				suite.Assert().NoError(err)
			}
		})
	}
}
