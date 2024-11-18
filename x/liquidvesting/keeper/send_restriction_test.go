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

const (
	restaker       = "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"
	testServiceId  = 1
	testPoolId     = 1
	testOperatorId = 1
)

func (suite *KeeperTestSuite) TestBankHooks_TestPoolRestaking() {
	testCases := []struct {
		name          string
		setup         func()
		msg           *restakingtypes.MsgDelegatePool
		expectedUsage sdk.Coins
		shouldErr     bool
	}{
		{
			name: "no insurance fund",
			setup: func() {
				// Simulate the minting of the staking representation
				suite.mintVestedRepresentation(
					restaker,
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegatePool(
				sdk.NewInt64Coin(vestedIBCDenom, 300),
				restaker,
			),
			shouldErr: true,
		},
		{
			name: "insufficient insurance fund",
			setup: func() {
				// Fund the user's insurance fund
				insuranceFundCoins := sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1))
				suite.fundAccountInsuranceFund(suite.ctx, restaker, insuranceFundCoins)

				// Simulate the minting of the staking representation
				suite.mintVestedRepresentation(
					restaker,
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegatePool(
				sdk.NewInt64Coin(vestedIBCDenom, 300),
				restaker,
			),
			shouldErr: true,
		},
		{
			name: "covered funds already restaked",
			setup: func() {
				// Create a test service and operator
				suite.createService(testServiceId)
				suite.createOperator(testOperatorId)

				insuranceFundCoins := sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6))
				suite.fundAccountInsuranceFund(suite.ctx, restaker, insuranceFundCoins)
				// Mint the staked representation
				suite.mintVestedRepresentation(
					restaker,
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)

				// Delegates the funds covered by the insurance fund to the service and the operator
				msgSrv := restakingkeeper.NewMsgServer(suite.rk)
				_, err := msgSrv.DelegateService(suite.ctx, restakingtypes.NewMsgDelegateService(
					testServiceId, sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 150)), restaker),
				)
				suite.Assert().NoError(err)
				_, err = msgSrv.DelegateOperator(suite.ctx, restakingtypes.NewMsgDelegateOperator(
					testOperatorId, sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 150)), restaker),
				)
				suite.Assert().NoError(err)
			},
			msg: restakingtypes.NewMsgDelegatePool(
				sdk.NewInt64Coin(vestedIBCDenom, 150),
				restaker,
			),
			shouldErr: true,
		},
		{
			name: "restake correctly",
			setup: func() {
				// Add the 2% to the insurance fund
				insuranceFundCoins := sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6))
				suite.fundAccountInsuranceFund(suite.ctx, restaker, insuranceFundCoins)

				// Simulate the minting of the staking representation
				suite.mintVestedRepresentation(
					restaker,
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegatePool(
				sdk.NewInt64Coin(vestedIBCDenom, 300),
				restaker,
			),
			expectedUsage: sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6)),
			shouldErr:     false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			suite.createPool(testPoolId, vestedIBCDenom)

			tc.setup()
			msgServer := restakingkeeper.NewMsgServer(suite.rk)

			_, err := msgServer.DelegatePool(suite.ctx, tc.msg)
			if tc.shouldErr {
				suite.Assert().Error(err)
			} else {
				suite.Assert().NoError(err)
				insuranceFund, err := suite.k.GetUserInsuranceFund(suite.ctx, sdk.MustAccAddressFromBech32(restaker))
				suite.Assert().NoError(err)
				suite.Assert().Equal(tc.expectedUsage, insuranceFund.Used)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestBankHooks_TestServiceRestaking() {
	testCases := []struct {
		name          string
		setup         func()
		msg           *restakingtypes.MsgDelegateService
		expectedUsage sdk.Coins
		shouldErr     bool
	}{
		{
			name: "no insurance fund",
			setup: func() {
				// Simulate the minting of the staking representation
				suite.mintVestedRepresentation(
					restaker,
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegateService(
				testServiceId,
				sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 300)),
				restaker,
			),
			shouldErr: true,
		},
		{
			name: "insufficient insurance fund",
			setup: func() {
				// Fund the user's insurance fund
				insuranceFundCoins := sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1))
				suite.fundAccountInsuranceFund(suite.ctx, restaker, insuranceFundCoins)

				// Simulate the minting of the staked representation
				suite.mintVestedRepresentation(
					restaker,
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegateService(
				testServiceId,
				sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 300)),
				restaker,
			),
			shouldErr: true,
		},
		{
			name: "covered funds already restaked",
			setup: func() {
				// Create a test pool and operator
				suite.createPool(testPoolId, vestedIBCDenom)
				suite.createOperator(testOperatorId)

				insuranceFundCoins := sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6))
				suite.fundAccountInsuranceFund(suite.ctx, restaker, insuranceFundCoins)
				// Mint the staked representation
				suite.mintVestedRepresentation(
					restaker,
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)

				// Delegates the funds covered by the insurance fund to a pool and an operator
				msgSrv := restakingkeeper.NewMsgServer(suite.rk)
				_, err := msgSrv.DelegatePool(suite.ctx, restakingtypes.NewMsgDelegatePool(
					sdk.NewInt64Coin(vestedIBCDenom, 150), restaker))
				suite.Assert().NoError(err)
				_, err = msgSrv.DelegateOperator(suite.ctx, restakingtypes.NewMsgDelegateOperator(
					testOperatorId, sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 150)), restaker),
				)
				suite.Assert().NoError(err)
			},
			msg: restakingtypes.NewMsgDelegateService(
				testServiceId,
				sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 150)),
				restaker,
			),
			shouldErr: true,
		},
		{
			name: "restake correctly",
			setup: func() {
				insuranceFundCoins := sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6))
				suite.fundAccountInsuranceFund(suite.ctx, restaker, insuranceFundCoins)
				// Mint the staked representation
				suite.mintVestedRepresentation(
					restaker,
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegateService(
				testServiceId,
				sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 300)),
				restaker,
			),
			expectedUsage: sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6)),
			shouldErr:     false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			suite.createService(testServiceId)

			tc.setup()
			msgServer := restakingkeeper.NewMsgServer(suite.rk)

			_, err := msgServer.DelegateService(suite.ctx, tc.msg)
			if tc.shouldErr {
				suite.Assert().Error(err)
			} else {
				suite.Assert().NoError(err)
				insuranceFund, err := suite.k.GetUserInsuranceFund(suite.ctx, sdk.MustAccAddressFromBech32(restaker))
				suite.Assert().NoError(err)
				suite.Assert().Equal(tc.expectedUsage, insuranceFund.Used)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestBankHooks_TestOperatorRestaking() {
	testCases := []struct {
		name          string
		setup         func()
		msg           *restakingtypes.MsgDelegateOperator
		expectedUsage sdk.Coins
		shouldErr     bool
	}{
		{
			name: "no insurance fund",
			setup: func() {
				// Simulate the minting of the staking representation
				suite.mintVestedRepresentation(
					restaker,
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegateOperator(
				testOperatorId,
				sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 300)),
				restaker,
			),
			shouldErr: true,
		},
		{
			name: "insufficient insurance fund",
			setup: func() {
				// Fund the user's insurance fund
				insuranceFundCoins := sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1))
				suite.fundAccountInsuranceFund(suite.ctx, restaker, insuranceFundCoins)

				// Simulate the minting of the staked representation
				suite.mintVestedRepresentation(
					restaker,
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegateOperator(
				testOperatorId,
				sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 300)),
				restaker,
			),
			shouldErr: true,
		},
		{
			name: "covered funds already restaked",
			setup: func() {
				suite.createPool(testPoolId, vestedIBCDenom)
				suite.createService(testServiceId)

				insuranceFundCoins := sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6))
				suite.fundAccountInsuranceFund(suite.ctx, restaker, insuranceFundCoins)
				// Mint the staked representation
				suite.mintVestedRepresentation(
					restaker,
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)

				// Delegates the funds covered by the insurance fund to a pool and an operator
				msgSrv := restakingkeeper.NewMsgServer(suite.rk)
				_, err := msgSrv.DelegatePool(suite.ctx, restakingtypes.NewMsgDelegatePool(
					sdk.NewInt64Coin(vestedIBCDenom, 150), restaker))
				suite.Assert().NoError(err)
				_, err = msgSrv.DelegateService(suite.ctx, restakingtypes.NewMsgDelegateService(
					testServiceId, sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 150)), restaker),
				)
				suite.Assert().NoError(err)
			},
			msg: restakingtypes.NewMsgDelegateOperator(
				testOperatorId,
				sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 150)),
				restaker,
			),
			shouldErr: true,
		},
		{
			name: "restake correctly",
			setup: func() {
				insuranceFundCoins := sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6))
				suite.fundAccountInsuranceFund(suite.ctx, restaker, insuranceFundCoins)
				// Mint the staked representation
				suite.mintVestedRepresentation(
					restaker,
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegateOperator(
				testOperatorId,
				sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 300)),
				restaker,
			),
			expectedUsage: sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6)),
			shouldErr:     false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			suite.createOperator(testOperatorId)

			tc.setup()
			msgServer := restakingkeeper.NewMsgServer(suite.rk)

			_, err := msgServer.DelegateOperator(suite.ctx, tc.msg)
			if tc.shouldErr {
				suite.Assert().Error(err)
			} else {
				suite.Assert().NoError(err)
				insuranceFund, err := suite.k.GetUserInsuranceFund(suite.ctx, sdk.MustAccAddressFromBech32(restaker))
				suite.Assert().NoError(err)
				suite.Assert().Equal(tc.expectedUsage, insuranceFund.Used)
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
