package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	restakingkeeper "github.com/milkyway-labs/milkyway/x/restaking/keeper"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
)

const (
	iBCDenom       = "ibc/D79E7D83AB399BFFF93433E54FAA480C191248FC556924A2A8351AE2638B3877"
	vestedDenom    = "vested/" + iBCDenom
	restaker       = "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"
	testServiceId  = 1
	testPoolId     = 1
	testOperatorId = 1
)

func (suite *KeeperTestSuite) TestBankHooks_TestPoolRestaking() {
	testCases := []struct {
		name      string
		setup     func()
		msg       *restakingtypes.MsgDelegatePool
		shouldErr bool
		errorMsg  string
	}{
		{
			name: "no insurance fund",
			setup: func() {
				// Simulate the minting of the staking representation
				suite.mintVestedRepresentation(
					restaker,
					sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegatePool(
				sdk.NewInt64Coin(vestedDenom, 300),
				restaker,
			),
			shouldErr: true,
			errorMsg:  "insufficient insurance fund",
		},
		{
			name: "insufficient insurance fund",
			setup: func() {
				// Fund the user's insurance fund
				insuranceFundCoins := sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 1))
				suite.fundAccountInsuranceFund(suite.ctx, restaker, insuranceFundCoins)

				// Simulate the minting of the staking representation
				suite.mintVestedRepresentation(
					restaker,
					sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegatePool(
				sdk.NewInt64Coin(vestedDenom, 300),
				restaker,
			),
			shouldErr: true,
			errorMsg:  "required: 6",
		},
		{
			name: "covered funds already restaked",
			setup: func() {
				// Create a test service and operator
				suite.createService(testServiceId)
				suite.createOperator(testOperatorId)

				insuranceFundCoins := sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 6))
				suite.fundAccountInsuranceFund(suite.ctx, restaker, insuranceFundCoins)
				// Mint the staked representation
				suite.mintVestedRepresentation(
					restaker,
					sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 10000)),
				)

				// Delegates the funds covered by the insurance fund to the service and the operator
				msgSrv := restakingkeeper.NewMsgServer(suite.rk)
				_, err := msgSrv.DelegateService(suite.ctx, restakingtypes.NewMsgDelegateService(
					testServiceId, sdk.NewCoins(sdk.NewInt64Coin(vestedDenom, 150)), restaker),
				)
				suite.Assert().NoError(err)
				_, err = msgSrv.DelegateOperator(suite.ctx, restakingtypes.NewMsgDelegateOperator(
					testOperatorId, sdk.NewCoins(sdk.NewInt64Coin(vestedDenom, 150)), restaker),
				)
				suite.Assert().NoError(err)
			},
			msg: restakingtypes.NewMsgDelegatePool(
				sdk.NewInt64Coin(vestedDenom, 150),
				restaker,
			),
			shouldErr: true,
			errorMsg:  "required: 9",
		},
		{
			name: "restake correctly",
			setup: func() {
				// Add the 2% to the insurance fund
				insuranceFundCoins := sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 6))
				suite.fundAccountInsuranceFund(suite.ctx, restaker, insuranceFundCoins)

				// Simulate the minting of the staking representation
				suite.mintVestedRepresentation(
					restaker,
					sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegatePool(
				sdk.NewInt64Coin(vestedDenom, 300),
				restaker,
			),
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			suite.createPool(testPoolId, vestedDenom)

			ctx, _ := suite.ctx.CacheContext()
			tc.setup()
			msgServer := restakingkeeper.NewMsgServer(suite.rk)

			_, err := msgServer.DelegatePool(ctx, tc.msg)
			if tc.shouldErr {
				if tc.errorMsg == "" {
					suite.Require().Error(err)
				} else {
					suite.Assert().ErrorContains(err, tc.errorMsg)
				}
			} else {
				suite.Assert().NoError(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestBankHooks_TestServiceRestaking() {
	testCases := []struct {
		name      string
		setup     func()
		msg       *restakingtypes.MsgDelegateService
		shouldErr bool
		errorMsg  string
	}{
		{
			name: "no insurance fund",
			setup: func() {
				// Simulate the minting of the staking representation
				suite.mintVestedRepresentation(
					restaker,
					sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegateService(
				testServiceId,
				sdk.NewCoins(sdk.NewInt64Coin(vestedDenom, 300)),
				restaker,
			),
			shouldErr: true,
			errorMsg:  "insufficient insurance fund",
		},
		{
			name: "insufficient insurance fund",
			setup: func() {
				// Fund the user's insurance fund
				insuranceFundCoins := sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 1))
				suite.fundAccountInsuranceFund(suite.ctx, restaker, insuranceFundCoins)

				// Simulate the minting of the staked representation
				suite.mintVestedRepresentation(
					restaker,
					sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegateService(
				testServiceId,
				sdk.NewCoins(sdk.NewInt64Coin(vestedDenom, 300)),
				restaker,
			),
			shouldErr: true,
			errorMsg:  "required: 6",
		},
		{
			name: "covered funds already restaked",
			setup: func() {
				// Create a test pool and operator
				suite.createPool(testPoolId, vestedDenom)
				suite.createOperator(testOperatorId)

				insuranceFundCoins := sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 6))
				suite.fundAccountInsuranceFund(suite.ctx, restaker, insuranceFundCoins)
				// Mint the staked representation
				suite.mintVestedRepresentation(
					restaker,
					sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 10000)),
				)

				// Delegates the funds covered by the insurance fund to a pool and an operator
				msgSrv := restakingkeeper.NewMsgServer(suite.rk)
				_, err := msgSrv.DelegatePool(suite.ctx, restakingtypes.NewMsgDelegatePool(
					sdk.NewInt64Coin(vestedDenom, 150), restaker))
				suite.Assert().NoError(err)
				_, err = msgSrv.DelegateOperator(suite.ctx, restakingtypes.NewMsgDelegateOperator(
					testOperatorId, sdk.NewCoins(sdk.NewInt64Coin(vestedDenom, 150)), restaker),
				)
				suite.Assert().NoError(err)
			},
			msg: restakingtypes.NewMsgDelegateService(
				testServiceId,
				sdk.NewCoins(sdk.NewInt64Coin(vestedDenom, 150)),
				restaker,
			),
			shouldErr: true,
			errorMsg:  "required: 9",
		},
		{
			name: "restake correctly",
			setup: func() {
				insuranceFundCoins := sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 6))
				suite.fundAccountInsuranceFund(suite.ctx, restaker, insuranceFundCoins)
				// Mint the staked representation
				suite.mintVestedRepresentation(
					restaker,
					sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegateService(
				testServiceId,
				sdk.NewCoins(sdk.NewInt64Coin(vestedDenom, 300)),
				restaker,
			),
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			suite.createService(testServiceId)

			ctx, _ := suite.ctx.CacheContext()
			tc.setup()
			msgServer := restakingkeeper.NewMsgServer(suite.rk)

			_, err := msgServer.DelegateService(ctx, tc.msg)
			if tc.shouldErr {
				if tc.errorMsg == "" {
					suite.Require().Error(err)
				} else {
					suite.Assert().ErrorContains(err, tc.errorMsg)
				}
			} else {
				suite.Assert().NoError(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestBankHooks_TestOperatorRestaking() {
	testCases := []struct {
		name      string
		setup     func()
		msg       *restakingtypes.MsgDelegateOperator
		shouldErr bool
		errorMsg  string
	}{
		{
			name: "no insurance fund",
			setup: func() {
				// Simulate the minting of the staking representation
				suite.mintVestedRepresentation(
					restaker,
					sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegateOperator(
				testOperatorId,
				sdk.NewCoins(sdk.NewInt64Coin(vestedDenom, 300)),
				restaker,
			),
			shouldErr: true,
			errorMsg:  "insufficient insurance fund",
		},
		{
			name: "insufficient insurance fund",
			setup: func() {
				// Fund the user's insurance fund
				insuranceFundCoins := sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 1))
				suite.fundAccountInsuranceFund(suite.ctx, restaker, insuranceFundCoins)

				// Simulate the minting of the staked representation
				suite.mintVestedRepresentation(
					restaker,
					sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegateOperator(
				testOperatorId,
				sdk.NewCoins(sdk.NewInt64Coin(vestedDenom, 300)),
				restaker,
			),
			shouldErr: true,
			errorMsg:  "required: 6",
		},
		{
			name: "covered funds already restaked",
			setup: func() {
				suite.createPool(testPoolId, vestedDenom)
				suite.createService(testServiceId)

				insuranceFundCoins := sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 6))
				suite.fundAccountInsuranceFund(suite.ctx, restaker, insuranceFundCoins)
				// Mint the staked representation
				suite.mintVestedRepresentation(
					restaker,
					sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 10000)),
				)

				// Delegates the funds covered by the insurance fund to a pool and an operator
				msgSrv := restakingkeeper.NewMsgServer(suite.rk)
				_, err := msgSrv.DelegatePool(suite.ctx, restakingtypes.NewMsgDelegatePool(
					sdk.NewInt64Coin(vestedDenom, 150), restaker))
				suite.Assert().NoError(err)
				_, err = msgSrv.DelegateService(suite.ctx, restakingtypes.NewMsgDelegateService(
					testServiceId, sdk.NewCoins(sdk.NewInt64Coin(vestedDenom, 150)), restaker),
				)
				suite.Assert().NoError(err)
			},
			msg: restakingtypes.NewMsgDelegateOperator(
				testOperatorId,
				sdk.NewCoins(sdk.NewInt64Coin(vestedDenom, 150)),
				restaker,
			),
			shouldErr: true,
			errorMsg:  "required: 9",
		},
		{
			name: "restake correctly",
			setup: func() {
				insuranceFundCoins := sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 6))
				suite.fundAccountInsuranceFund(suite.ctx, restaker, insuranceFundCoins)
				// Mint the staked representation
				suite.mintVestedRepresentation(
					restaker,
					sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 10000)),
				)
			},
			msg: restakingtypes.NewMsgDelegateOperator(
				testOperatorId,
				sdk.NewCoins(sdk.NewInt64Coin(vestedDenom, 300)),
				restaker,
			),
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			suite.createOperator(testOperatorId)

			ctx, _ := suite.ctx.CacheContext()
			tc.setup()
			msgServer := restakingkeeper.NewMsgServer(suite.rk)

			_, err := msgServer.DelegateOperator(ctx, tc.msg)
			if tc.shouldErr {
				if tc.errorMsg == "" {
					suite.Require().Error(err)
				} else {
					suite.Assert().ErrorContains(err, tc.errorMsg)
				}
			} else {
				suite.Assert().NoError(err)
			}
		})
	}
}