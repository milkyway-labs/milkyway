package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v3/x/liquidvesting/types"
	restakingtypes "github.com/milkyway-labs/milkyway/v3/x/restaking/types"
)

func (suite *KeeperTestSuite) TestKeeper_EndBlocker() {
	lockedStakeDenom, err := types.GetLockedRepresentationDenom("stake")
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
				err = suite.rk.SetParams(ctx, restakingtypes.NewParams(
					7*24*time.Hour,
					nil,
					restakingtypes.DefaultRestakingCap,
				))
				suite.Require().NoError(err)

				// Add some tokens to the user's insurance fund so they can restake the locked representation
				suite.fundAccountInsuranceFund(ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewInt64Coin("stake", 100)),
				)

				// Fund the account
				suite.mintLockedRepresentation(ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)),
				)
				suite.fundAccount(ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewInt64Coin("stake2", 200)),
				)

				// Delegate some locked representation to pool, service and operator
				suite.createPool(ctx, 1, lockedStakeDenom)
				_, err = suite.rk.DelegateToPool(ctx,
					sdk.NewInt64Coin(lockedStakeDenom, 200),
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				)
				suite.Assert().NoError(err)

				suite.createService(ctx, 1)
				_, err = suite.rk.DelegateToService(ctx,
					1,
					sdk.NewCoins(sdk.NewInt64Coin(lockedStakeDenom, 300)),
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				)
				suite.Assert().NoError(err)

				suite.createOperator(ctx, 1)
				_, err = suite.rk.DelegateToOperator(ctx,
					1,
					sdk.NewCoins(sdk.NewInt64Coin(lockedStakeDenom, 300)),
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				)
				suite.Assert().NoError(err)

				// Burn all the coins
				err = suite.k.BurnLockedRepresentation(ctx,
					sdk.MustAccAddressFromBech32("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
					sdk.NewCoins(sdk.NewInt64Coin(lockedStakeDenom, 1000)),
				)
				suite.Assert().NoError(err)
			},
			updateCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(ctx.BlockTime().Add(3 * 24 * time.Hour)) // 3 days later
			},
			shouldErr: false,
			check: func(ctx sdk.Context) {
				// The user shouldn't have the locked representation
				userBalance := suite.bk.GetAllBalances(ctx, sdk.MustAccAddressFromBech32("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"))
				suite.Assert().Equal(sdk.NewCoins(sdk.NewInt64Coin("stake2", 200)), userBalance)

				// The burn queue should contain our record
				toBurnCoins, err := suite.k.GetUnbondedCoinsFromQueue(ctx, ctx.BlockTime().Add(4*24*time.Hour))
				suite.Require().NoError(err)
				suite.Assert().Len(toBurnCoins, 1)

				suite.Assert().Equal("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4", toBurnCoins[0].DelegatorAddress)
				suite.Assert().Equal(sdk.NewCoins(sdk.NewInt64Coin(lockedStakeDenom, 800)), toBurnCoins[0].Amount)

				// The user insurance fund signal that there are still 20 coins used
				// to cover the restaking position
				usedUserInsuranceFund, err := suite.k.GetUserUsedInsuranceFund(ctx, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Assert().NoError(err)
				suite.Assert().Equal(sdk.NewCoins(sdk.NewInt64Coin("stake", 16)), usedUserInsuranceFund)
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
				err = suite.rk.SetParams(ctx, restakingtypes.NewParams(
					7*24*time.Hour,
					nil,
					restakingtypes.DefaultRestakingCap,
				))
				suite.Require().NoError(err)

				// Add some tokens to the user's insurance fund so they can restake
				// the locked representation
				suite.fundAccountInsuranceFund(ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewInt64Coin("stake", 100)),
				)

				// Fund the account
				suite.mintLockedRepresentation(ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)),
				)
				suite.fundAccount(ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewInt64Coin("stake2", 200)),
				)

				// Delegate some locked representation to pool, service and operator
				suite.createPool(ctx, 1, lockedStakeDenom)
				_, err = suite.rk.DelegateToPool(ctx, sdk.NewInt64Coin(lockedStakeDenom, 200), "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Assert().NoError(err)

				suite.createService(ctx, 1)
				_, err = suite.rk.DelegateToService(ctx,
					1,
					sdk.NewCoins(sdk.NewInt64Coin(lockedStakeDenom, 300)),
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				)
				suite.Assert().NoError(err)

				suite.createOperator(ctx, 1)
				_, err = suite.rk.DelegateToOperator(ctx,
					1,
					sdk.NewCoins(sdk.NewInt64Coin(lockedStakeDenom, 300)),
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				)
				suite.Assert().NoError(err)

				// Burn all the coins
				err = suite.k.BurnLockedRepresentation(ctx,
					sdk.MustAccAddressFromBech32("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
					sdk.NewCoins(sdk.NewInt64Coin(lockedStakeDenom, 700)),
				)
				suite.Assert().NoError(err)
			},
			updateCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(ctx.BlockTime().Add(7 * 24 * time.Hour)) // 7 days later
			},
			shouldErr: false,
			check: func(ctx sdk.Context) {
				// The user shouldn't have the locked representation
				userBalance := suite.bk.GetAllBalances(ctx, sdk.MustAccAddressFromBech32("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"))
				suite.Assert().Equal(sdk.NewCoins(sdk.NewInt64Coin("stake2", 200)), userBalance)

				// The burn queue should be empty
				unbondingQueue, err := suite.k.GetUnbondedCoinsFromQueue(ctx, ctx.BlockTime())
				suite.Require().NoError(err)
				suite.Assert().Len(unbondingQueue, 0)

				// The user insurance fund should update properly
				userUsedInsuranceFund, err := suite.k.GetUserUsedInsuranceFund(ctx, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Assert().NoError(err)
				suite.Assert().Equal(sdk.NewCoins(sdk.NewInt64Coin("stake", 6)), userUsedInsuranceFund)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
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
			err = suite.rk.CompleteMatureUnbondingDelegations(ctx)
			suite.Assert().NoError(err)

			// run our end block logic
			err = suite.k.CompleteBurnCoins(ctx)
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
