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
			for address, amount := range tc.deposits {
				// Mint the coins that should be in the module
				suite.Assert().NoError(
					suite.bk.MintCoins(suite.ctx, types.ModuleName, amount))
				accAddress := sdk.MustAccAddressFromBech32(address)
				suite.Assert().NoError(
					suite.k.AddToUserInsuranceFund(suite.ctx, accAddress, amount))
			}

			for address, expectedAmount := range tc.deposits {
				accAddress := sdk.MustAccAddressFromBech32(address)
				amount, err := suite.k.GetUserInsuranceFundBalance(suite.ctx, accAddress)
				suite.Assert().NoError(err)
				suite.Assert().Equal(expectedAmount, amount)
			}

			balance, err := suite.k.GetInsuranceFundBalance(suite.ctx)
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
		setup     func()
		shouldErr bool
	}{
		{
			name: "withdraw more then deposited",
			from: "cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
			amount: sdk.NewCoins(
				sdk.NewInt64Coin("stake", 100),
				sdk.NewInt64Coin("stake2", 50),
			),
			setup: func() {
				suite.fundAccountInsuranceFund(suite.ctx,
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
			setup: func() {
				suite.fundAccountInsuranceFund(suite.ctx,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewCoins(
						sdk.NewInt64Coin("stake", 200),
						sdk.NewInt64Coin("stake2", 100),
					))
			},
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			accAddr := sdk.MustAccAddressFromBech32(tc.from)
			err := suite.k.WithdrawFromUserInsuranceFund(suite.ctx, accAddr, tc.amount)
			if tc.shouldErr {
				suite.Assert().Error(err)
			} else {
				suite.Assert().NoError(err)
			}
		})
	}
}
