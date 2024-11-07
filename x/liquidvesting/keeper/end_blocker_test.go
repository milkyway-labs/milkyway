package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
)

func (suite *KeeperTestSuite) TestKeeper_EndBlocker() {
	testAccount := "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"
	vestedStake, err := types.GetVestedRepresentationDenom("stake")
	suite.Assert().NoError(err)

	testCases := []struct {
		name      string
		setupCtx  func(sdk.Context) sdk.Context
		store     func(sdk.Context)
		updateCtx func(sdk.Context) sdk.Context
		shouldErr bool
		check     func(sdk.Context)
	}{
		{
			name:      "run without coins to burn",
			shouldErr: false,
		},
		{
			name: "coins are not burn until unbonded",
			setupCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.
					WithBlockHeight(10).
					WithBlockTime(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
			},
			store: func(ctx sdk.Context) {
				// Set the unbonding delegation time to 7 days
				suite.rk.SetParams(ctx, restakingtypes.NewParams(7*24*time.Hour, nil))

				// Add some tokens to the user's insurance fund so they can restake
				// the vested representation
				suite.fundAccountInsuranceFund(ctx, testAccount, sdk.NewCoins(sdk.NewInt64Coin("stake", 100)))
				// Fund the account
				suite.mintVestedRepresentation(testAccount, sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)))
				suite.fundAccount(ctx, testAccount, sdk.NewCoins(sdk.NewInt64Coin("stake2", 200)))

				// Delegate some vested representation to pool, service and operator
				suite.createPool(1, vestedStake)
				_, err := suite.rk.DelegateToPool(ctx, sdk.NewInt64Coin(vestedStake, 200), testAccount)
				suite.Assert().NoError(err)

				suite.createService(1)
				_, err = suite.rk.DelegateToService(ctx, 1,
					sdk.NewCoins(sdk.NewInt64Coin(vestedStake, 300)), testAccount)
				suite.Assert().NoError(err)

				suite.createOperator(1)
				_, err = suite.rk.DelegateToOperator(ctx, 1,
					sdk.NewCoins(sdk.NewInt64Coin(vestedStake, 300)), testAccount)
				suite.Assert().NoError(err)

				// Burn all the coins
				err = suite.k.BurnVestedRepresentation(ctx,
					sdk.MustAccAddressFromBech32(testAccount),
					sdk.NewCoins(sdk.NewInt64Coin(vestedStake, 1000)))
				suite.Assert().NoError(err)
			},
			updateCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(ctx.BlockTime().Add(3 * 24 * time.Hour)) // 3 days later
			},
			shouldErr: false,
			check: func(ctx sdk.Context) {
				// The user shouldn't have the vested representation
				userBalance := suite.bk.GetAllBalances(ctx, sdk.MustAccAddressFromBech32(testAccount))
				suite.Assert().Equal(sdk.NewCoins(sdk.NewInt64Coin("stake2", 200)), userBalance)

				// The burn queue should contain our record
				toBurnCoins := suite.k.GetUnbondedCoinsFromQueue(ctx, ctx.BlockTime().Add(4*24*time.Hour))
				suite.Assert().Len(toBurnCoins, 1)

				suite.Assert().Equal(testAccount, toBurnCoins[0].DelegatorAddress)
				suite.Assert().Equal(
					sdk.NewCoins(sdk.NewInt64Coin(vestedStake, 800)),
					toBurnCoins[0].Amount)

				// The user insurance fund signal that there are still 20 coins used
				// to cover the restaking position
				userInsuranceFund, err := suite.k.GetUserInsuranceFund(ctx, sdk.MustAccAddressFromBech32(testAccount))
				suite.Assert().NoError(err)
				suite.Assert().Equal(sdk.NewCoins(sdk.NewInt64Coin("stake", 16)), userInsuranceFund.Used)
			},
		},
		{
			name: "burns enqueued coins",
			setupCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.
					WithBlockHeight(10).
					WithBlockTime(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
			},
			store: func(ctx sdk.Context) {
				// Set the unbonding delegation time to 7 days
				suite.rk.SetParams(ctx, restakingtypes.NewParams(7*24*time.Hour, nil))

				// Add some tokens to the user's insurance fund so they can restake
				// the vested representation
				suite.fundAccountInsuranceFund(ctx, testAccount, sdk.NewCoins(sdk.NewInt64Coin("stake", 100)))
				// Fund the account
				suite.mintVestedRepresentation(testAccount, sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)))
				suite.fundAccount(ctx, testAccount, sdk.NewCoins(sdk.NewInt64Coin("stake2", 200)))

				// Delegate some vested representation to pool, service and operator
				suite.createPool(1, vestedStake)
				_, err := suite.rk.DelegateToPool(ctx, sdk.NewInt64Coin(vestedStake, 200), testAccount)
				suite.Assert().NoError(err)

				suite.createService(1)
				_, err = suite.rk.DelegateToService(ctx, 1,
					sdk.NewCoins(sdk.NewInt64Coin(vestedStake, 300)), testAccount)
				suite.Assert().NoError(err)

				suite.createOperator(1)
				_, err = suite.rk.DelegateToOperator(ctx, 1,
					sdk.NewCoins(sdk.NewInt64Coin(vestedStake, 300)), testAccount)
				suite.Assert().NoError(err)

				// Burn all the coins
				err = suite.k.BurnVestedRepresentation(ctx,
					sdk.MustAccAddressFromBech32(testAccount),
					sdk.NewCoins(sdk.NewInt64Coin(vestedStake, 700)))
				suite.Assert().NoError(err)
			},
			updateCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(ctx.BlockTime().Add(7 * 24 * time.Hour)) // 7 days later
			},
			shouldErr: false,
			check: func(ctx sdk.Context) {
				// The user shouldn't have the vested representation
				userBalance := suite.bk.GetAllBalances(ctx, sdk.MustAccAddressFromBech32(testAccount))
				suite.Assert().Equal(sdk.NewCoins(sdk.NewInt64Coin("stake2", 200)), userBalance)

				// The burn queue should be empty
				suite.Assert().Len(suite.k.GetUnbondedCoinsFromQueue(ctx, ctx.BlockTime()), 0)

				// The user insurance fund should update properly
				userInsuranceFund, err := suite.k.GetUserInsuranceFund(ctx, sdk.MustAccAddressFromBech32(testAccount))
				suite.Assert().NoError(err)
				suite.Assert().Equal(sdk.NewCoins(sdk.NewInt64Coin("stake", 6)), userInsuranceFund.Used)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx := suite.ctx

			if tc.setupCtx != nil {
				ctx = tc.setupCtx(ctx)
			}

			if tc.store != nil {
				tc.store(ctx)
			}

			if tc.updateCtx != nil {
				ctx = tc.updateCtx(ctx)
			}

			// run the restaking keep end block logic
			suite.Assert().NoError(suite.rk.CompleteMatureUnbondingDelegations(ctx))

			// run our end block logic
			err := suite.k.CompleteBurnCoins(ctx)

			if tc.shouldErr {
				suite.Assert().Error(err)
			} else {
				suite.Assert().NoError(err)
				if tc.check != nil {
					tc.check(ctx)
				}
			}
		})
	}
}
