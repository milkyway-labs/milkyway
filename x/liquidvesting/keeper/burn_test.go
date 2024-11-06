package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
)

func (suite *KeeperTestSuite) TestKeeper_TestBurn() {
	testAccount := "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"
	vestedStake, err := types.GetVestedRepresentationDenom("stake")
	suite.Assert().NoError(err)

	testCases := []struct {
		name      string
		setup     func(ctx sdk.Context)
		account   string
		amount    sdk.Coins
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name: "burn non vested representation fails",
			setup: func(ctx sdk.Context) {
				suite.fundAccount(ctx, testAccount, sdk.NewCoins(
					sdk.NewInt64Coin("stake", 1000), sdk.NewInt64Coin("stake2", 200)))
			},
			account:   testAccount,
			amount:    sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)),
			shouldErr: true,
		},
		{
			name: "burn with no funds fails",
			setup: func(ctx sdk.Context) {
				suite.mintVestedRepresentation(testAccount, sdk.NewCoins(sdk.NewInt64Coin("test", 1000)))
			},
			account:   testAccount,
			amount:    sdk.NewCoins(sdk.NewInt64Coin(vestedStake, 1000)),
			shouldErr: true,
		},
		{
			name: "burn from user balance",
			setup: func(ctx sdk.Context) {
				suite.mintVestedRepresentation(testAccount, sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)))
				suite.fundAccount(ctx, testAccount, sdk.NewCoins(sdk.NewInt64Coin("stake2", 200)))
			},
			account:   testAccount,
			amount:    sdk.NewCoins(sdk.NewInt64Coin(vestedStake, 1000)),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				coins := suite.bk.GetAllBalances(ctx, sdk.MustAccAddressFromBech32(testAccount))
				// we should have only the non vested coins
				suite.Assert().Equal(sdk.NewCoins(sdk.NewInt64Coin("stake2", 200)), coins)
			},
		},
		{
			name: "burn from delegations",
			setup: func(ctx sdk.Context) {
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
			},
			account:   testAccount,
			amount:    sdk.NewCoins(sdk.NewInt64Coin(vestedStake, 1000)),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				coins := suite.bk.GetAllBalances(ctx, sdk.MustAccAddressFromBech32(testAccount))
				suite.Assert().Equal(sdk.NewCoins(sdk.NewInt64Coin("stake2", 200)), coins)

				// Compute when the unbond will end
				unbondingTime, err := suite.rk.UnbondingTime(ctx)
				suite.Require().NoError(err)
				unbodnEnd := ctx.BlockHeader().Time.Add(unbondingTime)
				// Deque the values
				values := suite.k.GetUnbondedCoinsFromQueue(ctx, unbodnEnd)
				suite.Assert().Len(values, 1)
				// Check that we are burning the coins that have been delegated
				toBurnCoins := values[0].Amount
				suite.Assert().Equal(sdk.NewCoins(sdk.NewInt64Coin(vestedStake, 800)), toBurnCoins)
			},
		},
		{
			name: "burn more coins then owned delegations",
			setup: func(ctx sdk.Context) {
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
			},
			account:   testAccount,
			amount:    sdk.NewCoins(sdk.NewInt64Coin(vestedStake, 2000)),
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			if tc.setup != nil {
				tc.setup(suite.ctx)
			}

			err := suite.k.BurnVestedRepresentation(suite.ctx,
				sdk.MustAccAddressFromBech32(tc.account), tc.amount)

			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				if tc.check != nil {
					tc.check(suite.ctx)
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_TestIsBurner() {
	burnerAccount := "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"
	testAccount := "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"

	testCases := []struct {
		name     string
		account  string
		isBurner bool
	}{
		{
			name:     "not burner should fail",
			account:  testAccount,
			isBurner: false,
		},
		{
			name:     "valid burner",
			account:  burnerAccount,
			isBurner: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			suite.Assert().NoError(
				suite.k.SetParams(suite.ctx, types.NewParams(
					math.LegacyMustNewDecFromStr("2.0"),
					[]string{burnerAccount},
					nil)))

			isBurner, err := suite.k.IsBurner(suite.ctx, sdk.MustAccAddressFromBech32(tc.account))
			suite.Assert().NoError(err)
			suite.Assert().Equal(tc.isBurner, isBurner)
		})
	}
}
