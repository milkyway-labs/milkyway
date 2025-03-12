package keeper_test

import (
	"fmt"
	"time"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	operatorstypes "github.com/milkyway-labs/milkyway/v10/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/v10/x/pools/types"
	"github.com/milkyway-labs/milkyway/v10/x/restaking/keeper"
	"github.com/milkyway-labs/milkyway/v10/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v10/x/services/types"
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

func (suite *KeeperTestSuite) TestGasConsumption_UndelegateVsEndBlockProcessing() {
	for _, numDenomsPerDelegation := range []int{1, 10, 30, 100, 500} {
		suite.Run(fmt.Sprintf("with%dDenoms", numDenomsPerDelegation), func() {
			ctx, _ := suite.ctx.CacheContext()

			ctx = ctx.
				WithBlockHeight(10).
				WithBlockTime(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))

			// Set the unbonding time to 7 days
			params, err := suite.k.GetParams(ctx)
			suite.Require().NoError(err)
			params.UnbondingTime = 7 * 24 * time.Hour
			err = suite.k.SetParams(ctx, params)
			suite.Require().NoError(err)

			// Create a delegator address
			delegator := "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"

			// Generate denoms to delegate
			denoms := make([]string, numDenomsPerDelegation)
			for i := range denoms {
				denoms[i] = fmt.Sprintf("denom%d", i)
			}

			// Create a service to delegate to
			const serviceID = 1
			serviceAddress := servicestypes.GetServiceAddress(serviceID).String()
			err = suite.sk.SaveService(ctx, servicestypes.Service{
				ID:              serviceID,
				Status:          servicestypes.SERVICE_STATUS_ACTIVE,
				Address:         serviceAddress,
				Tokens:          sdk.NewCoins(),
				DelegatorShares: sdk.NewDecCoins(),
			})
			suite.Require().NoError(err)

			// Fund the delegator account with sufficient balance
			initialBalance := sdk.NewCoins()
			for _, denom := range denoms {
				initialBalance = initialBalance.Add(sdk.NewInt64Coin(denom, 1000_000000))
			}
			suite.fundAccount(ctx, delegator, initialBalance)

			// Prepare the total amount to delegate
			delAmt := sdk.NewCoins()
			for _, denom := range denoms {
				delAmt = delAmt.Add(sdk.NewInt64Coin(denom, 100_000000))
			}

			// Delegate multiple denominations to the service using MsgServer
			//
			// NOTE: We don't include this as part of the gas comparison because delegations can
			// be cheaply accumulated over arbitrary amounts of time and undelegated in a batch in a single block.
			// The attack vector hinges on this burst of undelegations to pack >1 block worth of unmetered gas
			// into a single block of gas consumption.
			msgServer := keeper.NewMsgServer(suite.k)
			_, err = msgServer.DelegateService(ctx, types.NewMsgDelegateService(serviceID, delAmt, delegator))
			suite.Require().NoError(err)

			// --- Undelegation Gas Tracking ---

			// Measure gas consumption during initial undelegation
			ctx = ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())

			// Undelegate the denominations using MsgServer
			_, err = msgServer.UndelegateService(ctx, types.NewMsgUndelegateService(serviceID, delAmt, delegator))
			suite.Require().NoError(err)

			// Calculate gas used during undelegations
			gasUsedForUndelegation := ctx.GasMeter().GasConsumed()
			fmt.Println("Gas used for undelegation:", gasUsedForUndelegation)

			// --- EndBlock Unbond Completion Gas Tracking ---

			// Advance context time to when undelegations mature
			ctx = ctx.WithBlockTime(ctx.BlockTime().Add(7 * 24 * time.Hour))

			ctx = ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())

			// Measure gas consumption during end block processing
			// NOTE: we isolate the component of the EndBlock call we want to test
			err = suite.k.CompleteMatureUnbondingDelegations(ctx)
			suite.Require().NoError(err)

			// Calculate gas used during end block processing
			gasUsedForEndBlock := ctx.GasMeter().GasConsumed()
			fmt.Println("Gas used for end block:", gasUsedForEndBlock)

			// --- Gas Consumption Comparison ---

			suite.Require().GreaterOrEqual(gasUsedForUndelegation, gasUsedForEndBlock)
		})
	}
}
