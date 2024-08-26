package keeper_test

import (
	"context"
	"time"

	"github.com/milkyway-labs/milkyway/app/testutil"
	"github.com/milkyway-labs/milkyway/utils"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
)

func (suite *KeeperTestSuite) TestMsgCreateRewardsPlan() {
	service, _ := suite.setupSampleServiceAndOperator()
	msgCreateRewardsPlan := types.NewMsgCreateRewardsPlan(
		testutil.TestAddress(1).String(),
		"Rewards Plan",
		service.ID,
		utils.MustParseCoins("100_000000service"),
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		types.NewBasicPoolsDistribution(0),
		types.NewBasicOperatorsDistribution(0),
		types.NewBasicUsersDistribution(0),
	)

	testCases := []struct {
		name        string
		preRun      func(ctx context.Context)
		msg         *types.MsgCreateRewardsPlan
		check       func(ctx context.Context, resp *types.MsgCreateRewardsPlanResponse)
		expectedErr string
	}{
		{
			name: "service not found returns error",
			msg: types.NewMsgCreateRewardsPlan(
				msgCreateRewardsPlan.Sender,
				msgCreateRewardsPlan.Description,
				2,
				msgCreateRewardsPlan.Amount,
				msgCreateRewardsPlan.StartTime,
				msgCreateRewardsPlan.EndTime,
				msgCreateRewardsPlan.PoolsDistribution,
				msgCreateRewardsPlan.OperatorsDistribution,
				msgCreateRewardsPlan.UsersDistribution,
			),
			expectedErr: "service not found: not found",
		},
		{
			name: "success",
			msg:  msgCreateRewardsPlan,
			check: func(ctx context.Context, resp *types.MsgCreateRewardsPlanResponse) {
				suite.Require().Equal(uint64(1), resp.NewRewardsPlanID)
				_, err := suite.keeper.GetRewardsPlan(ctx, resp.NewRewardsPlanID)
				suite.Require().NoError(err)
			},
		},
		{
			name: "rewards plan creation fee is charged",
			preRun: func(ctx context.Context) {
				// Change rewards plan creation fee to 100 $MILK.
				params, err := suite.keeper.Params.Get(ctx)
				suite.Require().NoError(err)
				params.RewardsPlanCreationFee = utils.MustParseCoins("100_000000umilk")
				err = suite.keeper.Params.Set(ctx, params)
				suite.Require().NoError(err)

				// Fund the sender account enough coins to pay the fee.
				suite.FundAccount(msgCreateRewardsPlan.Sender, utils.MustParseCoins("500_000000umilk"))
			},
			msg: msgCreateRewardsPlan,
			check: func(ctx context.Context, resp *types.MsgCreateRewardsPlanResponse) {
				// Check that the balance is decreased by amount of the fee.
				senderAddr, err := suite.App.AccountKeeper.AddressCodec().StringToBytes(msgCreateRewardsPlan.Sender)
				suite.Require().NoError(err)
				balances := suite.App.BankKeeper.GetAllBalances(ctx, senderAddr)
				suite.Require().Equal("400000000umilk", balances.String())
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			ctx, _ := suite.Ctx.CacheContext()
			if tc.preRun != nil {
				tc.preRun(ctx)
			}
			resp, err := suite.msgServer.CreateRewardsPlan(ctx, tc.msg)
			if tc.expectedErr == "" {
				suite.Require().NoError(err)
				if tc.check != nil {
					tc.check(ctx, resp)
				}
			} else {
				suite.Require().EqualError(err, tc.expectedErr)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgSetWithdrawAddress() {
	testCases := []struct {
		name        string
		msg         *types.MsgSetWithdrawAddress
		check       func(ctx context.Context)
		expectedErr string
	}{
		{
			name: "success",
			msg: types.NewMsgSetWithdrawAddress(
				testutil.TestAddress(1).String(),
				testutil.TestAddress(2).String(),
			),
			check: func(ctx context.Context) {
				withdrawAddr, err := suite.keeper.GetDelegatorWithdrawAddr(ctx, testutil.TestAddress(1))
				suite.Require().NoError(err)
				suite.Require().Equal(testutil.TestAddress(2), withdrawAddr)
			},
			expectedErr: "",
		},
		{
			name: "invalid sender address returns error",
			msg: types.NewMsgSetWithdrawAddress(
				"invalid",
				testutil.TestAddress(2).String(),
			),
			expectedErr: "invalid sender address: decoding bech32 failed: invalid bech32 string length 7: invalid address",
		},
		{
			name: "invalid withdraw address returns error",
			msg: types.NewMsgSetWithdrawAddress(
				testutil.TestAddress(1).String(),
				"invalid",
			),
			expectedErr: "invalid withdraw address: decoding bech32 failed: invalid bech32 string length 7: invalid address",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			ctx, _ := suite.Ctx.CacheContext()
			_, err := suite.msgServer.SetWithdrawAddress(ctx, tc.msg)
			if tc.expectedErr == "" {
				suite.Require().NoError(err)
				if tc.check != nil {
					tc.check(ctx)
				}
			} else {
				suite.Require().EqualError(err, tc.expectedErr)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgWithdrawDelegatorReward() {
	service, operator := suite.setupSampleServiceAndOperator()

	suite.CreateBasicRewardsPlan(
		service.ID, utils.MustParseCoins("100_000000service"),
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		utils.MustParseCoins("10000_000000service"))

	delAddr := testutil.TestAddress(1)

	suite.DelegateService(service.ID, utils.MustParseCoins("100_000000umilk"), delAddr.String(), true)

	suite.allocateRewards(10 * time.Second)

	testCases := []struct {
		name        string
		msg         *types.MsgWithdrawDelegatorReward
		check       func(ctx context.Context)
		expectedErr string
	}{
		{
			name: "success",
			msg: types.NewMsgWithdrawDelegatorReward(
				delAddr.String(), restakingtypes.DELEGATION_TYPE_SERVICE, service.ID,
			),
			check: func(ctx context.Context) {
				balances := suite.App.BankKeeper.GetAllBalances(ctx, delAddr)
				suite.Require().Equal("11574service", balances.String())
			},
			expectedErr: "",
		},
		{
			name: "invalid delegator address returns error",
			msg: types.NewMsgWithdrawDelegatorReward(
				"invalid", restakingtypes.DELEGATION_TYPE_SERVICE, service.ID,
			),
			expectedErr: "invalid delegator address: decoding bech32 failed: invalid bech32 string length 7: invalid address",
		},
		{
			name: "invalid delegation type returns error",
			msg: types.NewMsgWithdrawDelegatorReward(
				delAddr.String(), restakingtypes.DelegationType(5), service.ID,
			),
			expectedErr: "invalid delegation type: 5: invalid delegation type",
		},
		{
			name: "invalid target ID returns error",
			msg: types.NewMsgWithdrawDelegatorReward(
				delAddr.String(), restakingtypes.DELEGATION_TYPE_SERVICE, 0,
			),
			expectedErr: "invalid delegation target ID: 0: invalid request",
		},
		{
			name: "delegation not found",
			msg: types.NewMsgWithdrawDelegatorReward(
				delAddr.String(), restakingtypes.DELEGATION_TYPE_OPERATOR, operator.ID,
			),
			expectedErr: "delegation not found: 1, cosmos103vfz2vlvjyl3v2qalnlpnvtecdrdaxhs725g07fcw9acfkwsaps2jwxt9: not found",
		},
		{
			name: "delegation not found #2",
			msg: types.NewMsgWithdrawDelegatorReward(
				delAddr.String(), restakingtypes.DELEGATION_TYPE_SERVICE, 3,
			),
			expectedErr: "service not found: not found",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			ctx, _ := suite.Ctx.CacheContext()
			_, err := suite.msgServer.WithdrawDelegatorReward(ctx, tc.msg)
			if tc.expectedErr == "" {
				suite.Require().NoError(err)
				if tc.check != nil {
					tc.check(ctx)
				}
			} else {
				suite.Require().EqualError(err, tc.expectedErr)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgWithdrawOperatorCommission() {
	service, operator := suite.setupSampleServiceAndOperator()

	suite.CreateBasicRewardsPlan(
		service.ID, utils.MustParseCoins("100_000000service"),
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		utils.MustParseCoins("10000_000000service"))

	delAddr := testutil.TestAddress(1)

	suite.DelegateOperator(operator.ID, utils.MustParseCoins("100_000000umilk"), delAddr.String(), true)

	suite.allocateRewards(10 * time.Second)

	testCases := []struct {
		name        string
		msg         *types.MsgWithdrawOperatorCommission
		check       func(ctx context.Context)
		expectedErr string
	}{
		{
			name: "success",
			msg:  types.NewMsgWithdrawOperatorCommission(operator.Admin, operator.ID),
			check: func(ctx context.Context) {
				adminAddr, err := suite.App.AccountKeeper.AddressCodec().StringToBytes(operator.Admin)
				suite.Require().NoError(err)
				balances := suite.App.BankKeeper.GetAllBalances(ctx, adminAddr)
				suite.Require().Equal("1157service", balances.String())
			},
			expectedErr: "",
		},
		{
			name:        "invalid sender address returns error",
			msg:         types.NewMsgWithdrawOperatorCommission("invalid", operator.ID),
			expectedErr: "invalid sender address: decoding bech32 failed: invalid bech32 string length 7: invalid address",
		},
		{
			name:        "only admin can withdraw commission",
			msg:         types.NewMsgWithdrawOperatorCommission(testutil.TestAddress(1).String(), operator.ID),
			expectedErr: "only operator admin can withdraw operator commission: unauthorized",
		},
		{
			name:        "operator not found",
			msg:         types.NewMsgWithdrawOperatorCommission(operator.Admin, 3),
			expectedErr: "operator not found: not found",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			ctx, _ := suite.Ctx.CacheContext()
			_, err := suite.msgServer.WithdrawOperatorCommission(ctx, tc.msg)
			if tc.expectedErr == "" {
				suite.Require().NoError(err)
				if tc.check != nil {
					tc.check(ctx)
				}
			} else {
				suite.Require().EqualError(err, tc.expectedErr)
			}
		})
	}
}
