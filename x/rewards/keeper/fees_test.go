package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/rewards/types"
)

func (suite *KeeperTestSuite) TestKeeper_PayRegistrationFee() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		user      string
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name: "user without enough balance of any denom returns error",
			store: func(ctx sdk.Context) {
				// Set the module params
				err := suite.keeper.Params.Set(ctx, types.NewParams(
					sdk.NewCoins(
						sdk.NewInt64Coin("milktia", 20_000_000),
						sdk.NewInt64Coin("tia", 15_000_000),
						sdk.NewInt64Coin("umilk", 10_000_000),
					),
				))
				suite.Require().NoError(err)

				// Fund the user account
				suite.FundAccount(
					ctx,
					"cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
					sdk.NewCoins(
						sdk.NewInt64Coin("milktia", 1_000_000),
						sdk.NewInt64Coin("tia", 1_000_000),
						sdk.NewInt64Coin("umilk", 1_000_000),
					),
				)
			},
			user:      "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
			shouldErr: true,
		},
		{
			name: "user with enough balance of a single denom pays the fee",
			store: func(ctx sdk.Context) {
				// Set the module params
				err := suite.keeper.Params.Set(ctx, types.NewParams(
					sdk.NewCoins(
						sdk.NewInt64Coin("milktia", 20_000_000),
						sdk.NewInt64Coin("tia", 15_000_000),
						sdk.NewInt64Coin("umilk", 10_000_000),
					),
				))
				suite.Require().NoError(err)

				// Fund the user account
				suite.FundAccount(
					ctx,
					"cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
					sdk.NewCoins(
						sdk.NewInt64Coin("milktia", 10_000_000),
						sdk.NewInt64Coin("tia", 8_000_000),
						sdk.NewInt64Coin("umilk", 10_000_000),
					),
				)
			},
			user:      "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
			shouldErr: false,
			check: func(ctx sdk.Context) {
				// Make sure the user balance has been updated correctly
				userAddr, err := sdk.AccAddressFromBech32("cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn")
				suite.Require().NoError(err)

				balance := suite.bankKeeper.GetAllBalances(ctx, userAddr)
				suite.Require().Equal(
					sdk.NewCoins(
						sdk.NewInt64Coin("milktia", 10_000_000),
						sdk.NewInt64Coin("tia", 8_000_000),
					),
					balance,
				)
			},
		},
		{
			name: "user with enough balance of multiple denoms pays for the fist one that is found",
			store: func(ctx sdk.Context) {
				// Set the module params
				err := suite.keeper.Params.Set(ctx, types.NewParams(
					sdk.NewCoins(
						sdk.NewInt64Coin("milktia", 20_000_000),
						sdk.NewInt64Coin("tia", 15_000_000),
						sdk.NewInt64Coin("umilk", 10_000_000),
					),
				))
				suite.Require().NoError(err)

				// Fund the user account
				suite.FundAccount(
					ctx,
					"cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
					sdk.NewCoins(
						sdk.NewInt64Coin("milktia", 30_000_000),
						sdk.NewInt64Coin("tia", 30_000_000),
						sdk.NewInt64Coin("umilk", 30_000_000),
					),
				)
			},
			user:      "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
			shouldErr: false,
			check: func(ctx sdk.Context) {
				// Make sure the user balance has been updated correctly
				userAddr, err := sdk.AccAddressFromBech32("cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn")
				suite.Require().NoError(err)

				balance := suite.bankKeeper.GetAllBalances(ctx, userAddr)
				suite.Require().Equal(
					sdk.NewCoins(
						sdk.NewInt64Coin("milktia", 10_000_000),
						sdk.NewInt64Coin("tia", 30_000_000),
						sdk.NewInt64Coin("umilk", 30_000_000),
					),
					balance,
				)
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

			err := suite.keeper.PayRegistrationFees(ctx, tc.user)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}
