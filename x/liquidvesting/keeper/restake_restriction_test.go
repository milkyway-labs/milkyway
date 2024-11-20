package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	restakingkeeper "github.com/milkyway-labs/milkyway/x/restaking/keeper"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
)

func (suite *KeeperTestSuite) TestRestakeRestriction_TestPoolRestaking() {
	testCases := []struct {
		name          string
		store         func(ctx sdk.Context)
		msg           *restakingtypes.MsgDelegatePool
		expectedUsage sdk.Coins
		shouldErr     bool
	}{
		{
			name: "no insurance fund",
			store: func(ctx sdk.Context) {
				// Create the pool
				suite.createPool(ctx, 1, vestedIBCDenom)

				// Simulate the minting of the staking representation
				suite.mintVestedRepresentation(
					ctx,
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
				// Create the pool
				suite.createPool(ctx, 1, vestedIBCDenom)

				// Fund the user's insurance fund
				suite.fundAccountInsuranceFund(
					ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1)),
				)

				// Simulate the minting of the staking representation
				suite.mintVestedRepresentation(
					ctx,
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
				// Create the pool
				suite.createPool(ctx, 1, vestedIBCDenom)

				// Create a test service and operator
				suite.createService(ctx, 1)
				suite.createOperator(ctx, 1)

				// Fund the user account
				suite.fundAccountInsuranceFund(
					ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6)),
				)

				// Mint the staked representation
				suite.mintVestedRepresentation(
					ctx,
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
				// Add the 2% to the insurance fund
				insuranceFundCoins := sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6))
				suite.fundAccountInsuranceFund(ctx, "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre", insuranceFundCoins)

				// Simulate the minting of the staking representation
				suite.mintVestedRepresentation(
					ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegatePool(
				sdk.NewInt64Coin(vestedIBCDenom, 300),
				"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			),
			expectedUsage: sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6)),
			shouldErr:     false,
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

				usedInsuranceFund, err := suite.k.GetUserUsedInsuranceFund(ctx, "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre")
				suite.Assert().NoError(err)
				suite.Assert().Equal(tc.expectedUsage, usedInsuranceFund)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestRestakeRestriction_TestServiceRestaking() {
	testCases := []struct {
		name          string
		store         func(ctx sdk.Context)
		msg           *restakingtypes.MsgDelegateService
		expectedUsage sdk.Coins
		shouldErr     bool
	}{
		{
			name: "no insurance fund",
			store: func(ctx sdk.Context) {
				// Create a service
				suite.createService(ctx, 1)

				// Simulate the minting of the staking representation
				suite.mintVestedRepresentation(
					ctx,
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
				// Create a service
				suite.createService(ctx, 1)

				// Fund the user's insurance fund
				insuranceFundCoins := sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1))
				suite.fundAccountInsuranceFund(ctx, "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre", insuranceFundCoins)

				// Simulate the minting of the staked representation
				suite.mintVestedRepresentation(
					ctx,
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
				// Create a service
				suite.createService(ctx, 1)

				// Create a test pool and operator
				suite.createPool(ctx, 1, vestedIBCDenom)
				suite.createOperator(ctx, 1)

				// Fund the user account
				suite.fundAccountInsuranceFund(
					ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6)),
				)

				// Mint the staked representation
				suite.mintVestedRepresentation(
					ctx,
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
				// Create a service
				suite.createService(ctx, 1)

				// Fund the user account
				suite.fundAccountInsuranceFund(
					ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6)),
				)

				// Mint the staked representation
				suite.mintVestedRepresentation(
					ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegateService(
				1,
				sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 300)),
				"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			),
			expectedUsage: sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6)),
			shouldErr:     false,
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

				usedInsuranceFund, err := suite.k.GetUserUsedInsuranceFund(ctx, "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre")
				suite.Assert().NoError(err)
				suite.Assert().Equal(tc.expectedUsage, usedInsuranceFund)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestRestakeRestriction_TestOperatorRestaking() {
	testCases := []struct {
		name          string
		store         func(ctx sdk.Context)
		msg           *restakingtypes.MsgDelegateOperator
		expectedUsage sdk.Coins
		shouldErr     bool
	}{
		{
			name: "no insurance fund",
			store: func(ctx sdk.Context) {
				// Create an operator
				suite.createOperator(ctx, 1)

				// Simulate the minting of the staking representation
				suite.mintVestedRepresentation(
					ctx,
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
				// Create an operator
				suite.createOperator(ctx, 1)

				// Fund the user's insurance fund
				suite.fundAccountInsuranceFund(
					ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1)),
				)

				// Simulate the minting of the staked representation
				suite.mintVestedRepresentation(
					ctx,
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
				// Create an operator
				suite.createOperator(ctx, 1)

				suite.createPool(ctx, 1, vestedIBCDenom)
				suite.createService(ctx, 1)

				suite.fundAccountInsuranceFund(
					ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6)),
				)

				// Mint the staked representation
				suite.mintVestedRepresentation(
					ctx,
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
				// Create an operator
				suite.createOperator(ctx, 1)

				suite.fundAccountInsuranceFund(
					ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6)),
				)

				// Mint the staked representation
				suite.mintVestedRepresentation(
					ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegateOperator(
				1,
				sdk.NewCoins(sdk.NewInt64Coin(vestedIBCDenom, 300)),
				"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			),
			expectedUsage: sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 6)),
			shouldErr:     false,
		},
	}

	for _, tc := range testCases {
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

				usedInsuranceFund, err := suite.k.GetUserUsedInsuranceFund(ctx, "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre")
				suite.Assert().NoError(err)
				suite.Assert().Equal(tc.expectedUsage, usedInsuranceFund)
			}
		})
	}
}
