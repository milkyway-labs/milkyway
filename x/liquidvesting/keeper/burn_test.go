package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v3/x/liquidvesting/types"
)

func (suite *KeeperTestSuite) TestKeeper_TestBurn() {
	lockedStakeDenom, err := types.GetLockedRepresentationDenom("stake")
	suite.Assert().NoError(err)

	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		account   string
		amount    sdk.Coins
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name: "burn non locked representation fails",
			store: func(ctx sdk.Context) {
				suite.fundAccount(ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(
						sdk.NewInt64Coin("stake", 1000),
						sdk.NewInt64Coin("stake2", 200),
					),
				)
			},
			account:   "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			amount:    sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)),
			shouldErr: true,
		},
		{
			name: "burn with no funds fails",
			store: func(ctx sdk.Context) {
				suite.mintLockedRepresentation(ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewInt64Coin("test", 1000)),
				)
			},
			account:   "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			amount:    sdk.NewCoins(sdk.NewInt64Coin(lockedStakeDenom, 1000)),
			shouldErr: true,
		},
		{
			name: "burn from user balance",
			store: func(ctx sdk.Context) {
				suite.mintLockedRepresentation(ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)),
				)
				suite.fundAccount(ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewInt64Coin("stake2", 200)),
				)
			},
			account:   "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			amount:    sdk.NewCoins(sdk.NewInt64Coin(lockedStakeDenom, 1000)),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				userAddr, err := sdk.AccAddressFromBech32("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().NoError(err)

				// Make sure the user only has non locked coins
				coins := suite.bk.GetAllBalances(ctx, userAddr)
				suite.Assert().Equal(sdk.NewCoins(sdk.NewInt64Coin("stake2", 200)), coins)
			},
		},
		{
			name: "burn from delegations",
			store: func(ctx sdk.Context) {
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
			},
			account:   "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			amount:    sdk.NewCoins(sdk.NewInt64Coin(lockedStakeDenom, 1000)),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				userAddr, err := sdk.AccAddressFromBech32("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().NoError(err)

				coins := suite.bk.GetAllBalances(ctx, userAddr)
				suite.Assert().Equal(sdk.NewCoins(sdk.NewInt64Coin("stake2", 200)), coins)

				// Compute when the unbond will end
				unbondingTime, err := suite.rk.UnbondingTime(ctx)
				suite.Require().NoError(err)

				unbondEnd := ctx.BlockHeader().Time.Add(unbondingTime)

				// Deque the values
				values, err := suite.k.GetUnbondedCoinsFromQueue(ctx, unbondEnd)
				suite.Require().NoError(err)
				suite.Assert().Len(values, 1)

				// Check that we are burning the coins that have been delegated
				toBurnCoins := values[0].Amount
				suite.Assert().Equal(sdk.NewCoins(sdk.NewInt64Coin(lockedStakeDenom, 800)), toBurnCoins)
			},
		},
		{
			name: "burn more coins then owned delegations",
			store: func(ctx sdk.Context) {
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
			},
			account:   "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			amount:    sdk.NewCoins(sdk.NewInt64Coin(lockedStakeDenom, 2000)),
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			err = suite.k.BurnLockedRepresentation(ctx, sdk.MustAccAddressFromBech32(tc.account), tc.amount)

			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				if tc.check != nil {
					tc.check(ctx)
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_TestIsBurner() {
	testCases := []struct {
		name        string
		store       func(ctx sdk.Context)
		account     string
		shouldErr   bool
		expIsBurner bool
	}{
		{
			name:        "not burner should fail",
			account:     "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			shouldErr:   false,
			expIsBurner: false,
		},
		{
			name: "valid burner",
			store: func(ctx sdk.Context) {
				err := suite.k.SetParams(ctx, types.NewParams(
					math.LegacyMustNewDecFromStr("2.0"),
					[]string{"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"},
					nil,
					nil,
				))
				suite.Assert().NoError(err)
			},
			account:     "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			expIsBurner: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			isBurner, err := suite.k.IsBurner(ctx, sdk.MustAccAddressFromBech32(tc.account))
			if tc.shouldErr {
				suite.Assert().Error(err)
			} else {
				suite.Assert().NoError(err)
				suite.Assert().Equal(tc.expIsBurner, isBurner)
			}
		})
	}
}
