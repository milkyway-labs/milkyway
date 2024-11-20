package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	liquidvestingtypes "github.com/milkyway-labs/milkyway/x/liquidvesting/types"
	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	restakingkeeper "github.com/milkyway-labs/milkyway/x/restaking/keeper"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

func (suite *KeeperTestSuite) TestBankHooks_TestPoolRestaking() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		msg       *restakingtypes.MsgDelegatePool
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name: "no insurance fund",
			store: func(ctx sdk.Context) {
				suite.createPool(ctx, 1, vestedIBCDenom)

				// Simulate the minting of the staking representation
				suite.mintVestedRepresentation(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegatePool(
				sdk.NewInt64Coin(vestedIBCDenom, 300),
				"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			),
			shouldErr: true,
		},
		{
			name: "insufficient insurance fund",
			store: func(ctx sdk.Context) {
				suite.createPool(ctx, 1, vestedIBCDenom)

				// Fund the user's insurance fund
				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1)),
				)

				// Simulate the minting of the staking representation
				suite.mintVestedRepresentation(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegatePool(
				sdk.NewInt64Coin(vestedIBCDenom, 300),
				"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			),
			shouldErr: true,
		},
		{
			name: "covered funds already restaked",
			store: func(ctx sdk.Context) {
				suite.createPool(ctx, 1, vestedIBCDenom)

				// Create a test service and operator
				suite.createService(ctx, 1)
				suite.createOperator(ctx, 1)

				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6)),
				)

				// Mint the staked representation
				suite.mintVestedRepresentation(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)

				// Delegates the funds covered by the insurance fund to the service and the operator
				msgSrv := restakingkeeper.NewMsgServer(suite.rk)
				_, err := msgSrv.DelegateService(ctx, restakingtypes.NewMsgDelegateService(
					1,
					sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 150)),
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				))
				suite.Assert().NoError(err)
				_, err = msgSrv.DelegateOperator(ctx, restakingtypes.NewMsgDelegateOperator(
					1,
					sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 150)),
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				))
				suite.Assert().NoError(err)
			},
			msg: restakingtypes.NewMsgDelegatePool(
				sdk.NewInt64Coin(vestedIBCDenom, 150),
				"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			),
			shouldErr: true,
		},
		{
			name: "restake correctly",
			store: func(ctx sdk.Context) {
				suite.createPool(ctx, 1, vestedIBCDenom)

				// Add the 2% to the insurance fund
				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6)),
				)

				// Simulate the minting of the staking representation
				suite.mintVestedRepresentation(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegatePool(
				sdk.NewInt64Coin(vestedIBCDenom, 300),
				"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				userAddr, err := sdk.AccAddressFromBech32("cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre")
				suite.Assert().NoError(err)

				insuranceFund, err := suite.k.GetUserInsuranceFund(ctx, userAddr)
				suite.Assert().NoError(err)
				suite.Assert().Equal(sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6)), insuranceFund.Used)
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

			msgServer := restakingkeeper.NewMsgServer(suite.rk)
			_, err := msgServer.DelegatePool(ctx, tc.msg)
			if tc.shouldErr {
				suite.Assert().Error(err)
			} else {
				suite.Assert().NoError(err)
			}

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestBankHooks_TestServiceRestaking() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		msg       *restakingtypes.MsgDelegateService
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name: "no insurance fund",
			store: func(ctx sdk.Context) {
				suite.createService(ctx, 1)

				// Simulate the minting of the staking representation
				suite.mintVestedRepresentation(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegateService(
				1,
				sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 300)),
				"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			),
			shouldErr: true,
		},
		{
			name: "insufficient insurance fund",
			store: func(ctx sdk.Context) {
				suite.createService(ctx, 1)

				// Fund the user's insurance fund
				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1)),
				)

				// Simulate the minting of the staked representation
				suite.mintVestedRepresentation(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegateService(
				1,
				sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 300)),
				"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			),
			shouldErr: true,
		},
		{
			name: "covered funds already restaked",
			store: func(ctx sdk.Context) {
				suite.createService(ctx, 1)

				// Create a test pool and operator
				suite.createPool(ctx, 1, vestedIBCDenom)
				suite.createOperator(ctx, 1)

				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6)),
				)

				// Mint the staked representation
				suite.mintVestedRepresentation(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)

				// Delegates the funds covered by the insurance fund to a pool and an operator
				msgSrv := restakingkeeper.NewMsgServer(suite.rk)
				_, err := msgSrv.DelegatePool(ctx, restakingtypes.NewMsgDelegatePool(
					sdk.NewInt64Coin(vestedIBCDenom, 150),
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				))
				suite.Assert().NoError(err)

				_, err = msgSrv.DelegateOperator(ctx, restakingtypes.NewMsgDelegateOperator(
					1,
					sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 150)),
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				))
				suite.Assert().NoError(err)
			},
			msg: restakingtypes.NewMsgDelegateService(
				1,
				sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 150)),
				"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			),
			shouldErr: true,
		},
		{
			name: "restake correctly",
			store: func(ctx sdk.Context) {
				suite.createService(ctx, 1)

				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6)),
				)

				// Mint the staked representation
				suite.mintVestedRepresentation(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegateService(
				1,
				sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 300)),
				"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				userAddr, err := sdk.AccAddressFromBech32("cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre")
				suite.Assert().NoError(err)

				insuranceFund, err := suite.k.GetUserInsuranceFund(ctx, userAddr)
				suite.Assert().NoError(err)
				suite.Assert().Equal(sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6)), insuranceFund.Used)
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

			msgServer := restakingkeeper.NewMsgServer(suite.rk)
			_, err := msgServer.DelegateService(ctx, tc.msg)
			if tc.shouldErr {
				suite.Assert().Error(err)
			} else {
				suite.Assert().NoError(err)
			}

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestBankHooks_TestOperatorRestaking() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		msg       *restakingtypes.MsgDelegateOperator
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name: "no insurance fund",
			store: func(ctx sdk.Context) {
				suite.createOperator(ctx, 1)

				// Simulate the minting of the staking representation
				suite.mintVestedRepresentation(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegateOperator(
				1,
				sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 300)),
				"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			),
			shouldErr: true,
		},
		{
			name: "insufficient insurance fund",
			store: func(ctx sdk.Context) {
				suite.createOperator(ctx, 1)

				// Fund the user's insurance fund
				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1)),
				)

				// Simulate the minting of the staked representation
				suite.mintVestedRepresentation(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegateOperator(
				1,
				sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 300)),
				"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			),
			shouldErr: true,
		},
		{
			name: "covered funds already restaked",
			store: func(ctx sdk.Context) {
				suite.createOperator(ctx, 1)

				suite.createPool(ctx, 1, vestedIBCDenom)
				suite.createService(ctx, 1)

				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6)),
				)

				// Mint the staked representation
				suite.mintVestedRepresentation(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)

				// Delegates the funds covered by the insurance fund to a pool and an operator
				msgSrv := restakingkeeper.NewMsgServer(suite.rk)
				_, err := msgSrv.DelegatePool(ctx, restakingtypes.NewMsgDelegatePool(
					sdk.NewInt64Coin(vestedIBCDenom, 150),
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				))
				suite.Assert().NoError(err)

				_, err = msgSrv.DelegateService(ctx, restakingtypes.NewMsgDelegateService(
					1,
					sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 150)),
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				))
				suite.Assert().NoError(err)
			},
			msg: restakingtypes.NewMsgDelegateOperator(
				1,
				sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 150)),
				"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			),
			shouldErr: true,
		},
		{
			name: "restake correctly",
			store: func(ctx sdk.Context) {
				suite.createOperator(ctx, 1)

				suite.fundAccountInsuranceFund(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6)),
				)

				// Mint the staked representation
				suite.mintVestedRepresentation(ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegateOperator(
				1,
				sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 300)),
				"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				userAddr, err := sdk.AccAddressFromBech32("cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre")
				suite.Assert().NoError(err)

				insuranceFund, err := suite.k.GetUserInsuranceFund(ctx, userAddr)
				suite.Assert().NoError(err)
				suite.Assert().Equal(sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6)), insuranceFund.Used)
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

			msgServer := restakingkeeper.NewMsgServer(suite.rk)
			_, err := msgServer.DelegateOperator(ctx, tc.msg)
			if tc.shouldErr {
				suite.Assert().Error(err)
			} else {
				suite.Assert().NoError(err)
			}

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_SendRegistrionFn() {
	testCase := []struct {
		name      string
		store     func(ctx sdk.Context)
		from      string
		to        string
		amount    sdk.Coins
		shouldErr bool
		expTo     string
	}{
		{
			name:      "sending normal coins from user to user works properly",
			from:      "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			to:        "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			amount:    sdk.NewCoins(sdk.NewInt64Coin("ibc/1", 100)),
			shouldErr: false,
			expTo:     "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
		},
		{
			name:      "sending vested representation from user to user returns error",
			from:      "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			to:        "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			amount:    sdk.NewCoins(sdk.NewInt64Coin("vested/stake", 100)),
			shouldErr: true,
		},
		{
			name: "sending normal coins between restaking targets is not allowed",
			store: func(ctx sdk.Context) {
				// Create a test service and operator
				suite.createService(testServiceId)
				suite.createOperator(testOperatorId)
			},
			from:      servicestypes.GetServiceAddress(testServiceId).String(),
			to:        operatorstypes.GetOperatorAddress(testOperatorId).String(),
			amount:    sdk.NewCoins(sdk.NewInt64Coin("stake", 100)),
			shouldErr: true,
		},
		{
			name: "sending normal coins between restaking targets is not allowed",
			store: func(ctx sdk.Context) {
				// Create a test service and operator
				suite.createPool(testPoolId, vestedIBCDenom)
				suite.createOperator(testOperatorId)
			},
			from:      poolstypes.GetPoolAddress(testPoolId).String(),
			to:        operatorstypes.GetOperatorAddress(testOperatorId).String(),
			amount:    sdk.NewCoins(sdk.NewInt64Coin("vested/stake", 100)),
			shouldErr: true,
		},
		{
			name: "sending coins from the module account is allowed",
			from: authtypes.NewModuleAddress(liquidvestingtypes.ModuleName).String(),
			to:   "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			amount: sdk.NewCoins(
				sdk.NewInt64Coin("ibc/1", 100),
				sdk.NewInt64Coin("vested/stake", 200),
			),
			shouldErr: false,
			expTo:     "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
		},
		{
			name: "sending coins to the module account is allowed",
			from: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			to:   authtypes.NewModuleAddress(liquidvestingtypes.ModuleName).String(),
			amount: sdk.NewCoins(
				sdk.NewInt64Coin("ibc/1", 100),
				sdk.NewInt64Coin("vested/stake", 200),
			),
			shouldErr: false,
			expTo:     authtypes.NewModuleAddress(liquidvestingtypes.ModuleName).String(),
		},
	}

	for _, tc := range testCase {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			from, err := sdk.AccAddressFromBech32(tc.from)
			suite.Require().NoError(err)

			to, err := sdk.AccAddressFromBech32(tc.to)
			suite.Require().NoError(err)

			receivedTo, err := suite.k.SendRestrictionFn(ctx, from, to, tc.amount)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expTo, receivedTo.String())
			}
		})
	}
}
