package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	operatorstypes "github.com/milkyway-labs/milkyway/v6/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/v6/x/pools/types"
	"github.com/milkyway-labs/milkyway/v6/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v6/x/services/types"
)

func (suite *KeeperTestSuite) TestKeeper_CompleteMatureUnbondingDelegations() {
	testCases := []struct {
		name      string
		setupCtx  func(ctx sdk.Context) sdk.Context
		store     func(ctx sdk.Context)
		updateCtx func(ctx sdk.Context) sdk.Context
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name: "no mature unbonding delegations",
			setupCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.
					WithBlockHeight(10).
					WithBlockTime(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
			},
			store: func(ctx sdk.Context) {
				delegationAmount := sdk.NewCoin("umilk", sdkmath.NewInt(100))
				delegatorAddress := "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"

				// Set the unbonding delegation time to 7 days
				err := suite.k.SetParams(ctx, types.Params{
					UnbondingTime: 7 * 24 * time.Hour,
				})
				suite.Require().NoError(err)

				// Create a pool
				err = suite.pk.SavePool(ctx, poolstypes.NewPool(1, delegationAmount.Denom))
				suite.Require().NoError(err)

				// Send some tokens to the user
				suite.fundAccount(ctx, delegatorAddress, sdk.NewCoins(delegationAmount))

				// Delegate the tokens to the pool
				_, err = suite.k.DelegateToPool(ctx, delegationAmount, delegatorAddress)
				suite.Require().NoError(err)

				// Make sure the user has no tokens
				delegatorAddr, err := sdk.AccAddressFromBech32(delegatorAddress)
				suite.Require().NoError(err)

				balances := suite.bk.GetAllBalances(ctx, delegatorAddr)
				suite.Require().True(balances.IsZero())

				// Unbond the delegation
				completionTime, err := suite.k.UndelegateFromPool(ctx, delegationAmount, delegatorAddress)
				suite.Require().NoError(err)
				suite.Require().Equal(ctx.BlockTime().Add(7*24*time.Hour), completionTime)
			},
			updateCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(ctx.BlockTime().Add(3 * 24 * time.Hour)) // 3 days later
			},
			shouldErr: false,
			check: func(ctx sdk.Context) {
				// Make sure the delegation is still there and unbonding
				ubd, found, err := suite.k.GetUnbondingDelegation(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					types.DELEGATION_TYPE_POOL,
					1,
				)
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().Equal(types.NewPoolUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					1,
					10,
					time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					1,
				), ubd)
			},
		},
		{
			name: "mature unbonding delegation is matured properly",
			setupCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.
					WithBlockHeight(10).
					WithBlockTime(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
			},
			store: func(ctx sdk.Context) {
				delegationAmount := sdk.NewCoin("umilk", sdkmath.NewInt(100))
				delegatorAddress := "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"

				// Set the unbonding delegation time to 7 days
				err := suite.k.SetParams(ctx, types.Params{
					UnbondingTime: 7 * 24 * time.Hour,
				})
				suite.Require().NoError(err)

				// Create a pool
				err = suite.pk.SavePool(ctx, poolstypes.NewPool(1, delegationAmount.Denom))
				suite.Require().NoError(err)

				// Send some tokens to the user
				suite.fundAccount(ctx, delegatorAddress, sdk.NewCoins(delegationAmount))

				// Delegate the tokens to the pool
				_, err = suite.k.DelegateToPool(ctx, delegationAmount, delegatorAddress)
				suite.Require().NoError(err)

				// Make sure the user has no tokens
				delegatorAddr, err := sdk.AccAddressFromBech32(delegatorAddress)
				suite.Require().NoError(err)

				balances := suite.bk.GetAllBalances(ctx, delegatorAddr)
				suite.Require().True(balances.IsZero())

				// Unbond the delegation
				completionTime, err := suite.k.UndelegateFromPool(ctx, delegationAmount, delegatorAddress)
				suite.Require().NoError(err)
				suite.Require().Equal(ctx.BlockTime().Add(7*24*time.Hour), completionTime)
			},
			updateCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.
					WithBlockHeight(100).
					WithBlockTime(ctx.BlockTime().Add(7*24*time.Hour + time.Minute)) // 7 days later
			},
			shouldErr: false,
			check: func(ctx sdk.Context) {
				// Make sure the delegation is no longer unbonding
				_, found, err := suite.k.GetUnbondingDelegation(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					types.DELEGATION_TYPE_POOL,
					1,
				)
				suite.Require().NoError(err)
				suite.Require().False(found)

				// Make sure the user has the tokens
				delegatorAddr, err := sdk.AccAddressFromBech32("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().NoError(err)

				balances := suite.bk.GetAllBalances(ctx, delegatorAddr)
				suite.Require().Equal(sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))), balances)
			},
		},
		{
			name: "multiple partial unbonding delegations are matured properly",
			setupCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.
					WithBlockHeight(10).
					WithBlockTime(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
			},
			store: func(ctx sdk.Context) {
				milkBalance := sdk.NewCoin("umilk", sdkmath.NewInt(1_000))
				delegatorAddress := "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"

				// Set the unbonding delegation time to 7 days
				err := suite.k.SetParams(ctx, types.Params{
					UnbondingTime: 7 * 24 * time.Hour,
				})
				suite.Require().NoError(err)

				// Create a pool
				err = suite.pk.SavePool(ctx, poolstypes.NewPool(1, milkBalance.Denom))
				suite.Require().NoError(err)

				// Create an operator
				err = suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
					1,
					operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)

				// Create a service
				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is an AVS of a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)

				// Send some tokens to the user
				suite.fundAccount(ctx, delegatorAddress, sdk.NewCoins(milkBalance))

				// Delegate the tokens
				poolDelegationAmount := sdk.NewCoin("umilk", sdkmath.NewInt(200))
				_, err = suite.k.DelegateToPool(ctx, poolDelegationAmount, delegatorAddress)
				suite.Require().NoError(err)

				operatorDelegationAmount := sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(300)))
				_, err = suite.k.DelegateToOperator(ctx, 1, operatorDelegationAmount, delegatorAddress)
				suite.Require().NoError(err)

				serviceDelegationAmount := sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(500)))
				_, err = suite.k.DelegateToService(ctx, 1, serviceDelegationAmount, delegatorAddress)
				suite.Require().NoError(err)

				// Make sure the user has no tokens
				delegatorAddr, err := sdk.AccAddressFromBech32(delegatorAddress)
				suite.Require().NoError(err)

				balances := suite.bk.GetAllBalances(ctx, delegatorAddr)
				suite.Require().True(balances.IsZero())

				// Unbond the delegations
				poolUnbondingAmount := sdk.NewCoin("umilk", sdkmath.NewInt(100))
				completionTime, err := suite.k.UndelegateFromPool(ctx, poolUnbondingAmount, delegatorAddress)
				suite.Require().NoError(err)
				suite.Require().Equal(ctx.BlockTime().Add(7*24*time.Hour), completionTime)

				operatorUnbondingAmount := sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(150)))
				completionTime, err = suite.k.UndelegateFromOperator(ctx, 1, operatorUnbondingAmount, delegatorAddress)
				suite.Require().NoError(err)
				suite.Require().Equal(ctx.BlockTime().Add(7*24*time.Hour), completionTime)

				serviceUnbondingAmount := sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(250)))
				completionTime, err = suite.k.UndelegateFromService(ctx, 1, serviceUnbondingAmount, delegatorAddress)
				suite.Require().NoError(err)
				suite.Require().Equal(ctx.BlockTime().Add(7*24*time.Hour), completionTime)
			},
			updateCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.
					WithBlockHeight(100).
					WithBlockTime(ctx.BlockTime().Add(7 * 24 * time.Hour)) // 7 days later
			},
			shouldErr: false,
			check: func(ctx sdk.Context) {
				// Make sure the delegations are no longer unbonding
				_, found, err := suite.k.GetUnbondingDelegation(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					types.DELEGATION_TYPE_POOL,
					1,
				)
				suite.Require().NoError(err)
				suite.Require().False(found)

				_, found, err = suite.k.GetUnbondingDelegation(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					types.DELEGATION_TYPE_OPERATOR,
					1,
				)
				suite.Require().NoError(err)
				suite.Require().False(found)

				_, found, err = suite.k.GetUnbondingDelegation(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					types.DELEGATION_TYPE_SERVICE,
					1,
				)
				suite.Require().NoError(err)
				suite.Require().False(found)

				// Make sure the user has the tokens
				delegatorAddr, err := sdk.AccAddressFromBech32("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().NoError(err)

				balances := suite.bk.GetAllBalances(ctx, delegatorAddr)
				suite.Require().Equal(sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(500))), balances)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			// Setup the context
			ctx := suite.ctx
			if tc.setupCtx != nil {
				ctx = tc.setupCtx(ctx)
			}

			// Store the data
			if tc.store != nil {
				tc.store(ctx)
			}

			// Update the context (simulates the passing of time)
			if tc.updateCtx != nil {
				ctx = tc.updateCtx(ctx)
			}

			// Run the function
			err := suite.k.CompleteMatureUnbondingDelegations(ctx)
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
