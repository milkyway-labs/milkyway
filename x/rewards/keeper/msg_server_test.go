package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/milkyway-labs/milkyway/app/testutil"
	"github.com/milkyway-labs/milkyway/utils"
	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	"github.com/milkyway-labs/milkyway/x/rewards/keeper"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

func (suite *KeeperTestSuite) TestMsgCreateRewardsPlan() {
	testCases := []struct {
		name        string
		setup       func()
		store       func(ctx sdk.Context)
		setupCtx    func(ctx sdk.Context) sdk.Context
		msg         *types.MsgCreateRewardsPlan
		shouldErr   bool
		expResponse *types.MsgCreateRewardsPlanResponse
		expEvents   sdk.Events
		check       func(ctx sdk.Context)
	}{
		{
			name: "service not found returns error",
			msg: types.NewMsgCreateRewardsPlan(
				1,
				"Rewards Plan",
				utils.MustParseCoins("100_000000service"),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				types.NewBasicPoolsDistribution(0),
				types.NewBasicOperatorsDistribution(0),
				types.NewBasicUsersDistribution(0),
				sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100_000_000))),
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "sender different from admin returns error",
			store: func(ctx sdk.Context) {
				// Create a service
				_, _ = suite.setupSampleServiceAndOperator(ctx)

				// Create a service
				_, _ = suite.setupSampleServiceAndOperator(ctx)

				// Change rewards plan creation fee to 100 $MILK.
				err := suite.keeper.Params.Set(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewInt64Coin("umilk", 100_000_000)),
				))
				suite.Require().NoError(err)

				// Set the next plan id
				err = suite.keeper.NextRewardsPlanID.Set(ctx, 1)
				suite.Require().NoError(err)

				// Fund the sender account enough coins to pay the fee.
				suite.FundAccount(ctx, "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", utils.MustParseCoins("500_000000umilk"))
			},
			msg: types.NewMsgCreateRewardsPlan(
				1,
				"Rewards Plan",
				utils.MustParseCoins("100_000000service"),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				types.NewBasicPoolsDistribution(0),
				types.NewBasicOperatorsDistribution(0),
				types.NewBasicUsersDistribution(0),
				sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100_000_000))),
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "service is created and fee is charged",
			store: func(ctx sdk.Context) {
				// Create a service
				_, _ = suite.setupSampleServiceAndOperator(ctx)

				// Change rewards plan creation fee to 100 $MILK.
				err := suite.keeper.Params.Set(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewInt64Coin("umilk", 100_000_000)),
				))
				suite.Require().NoError(err)

				// Set the next plan id
				err = suite.keeper.NextRewardsPlanID.Set(ctx, 1)
				suite.Require().NoError(err)

				// Fund the sender account enough coins to pay the fee.
				suite.FundAccount(
					ctx,
					testutil.TestAddress(10000).String(),
					utils.MustParseCoins("500_000000umilk"),
				)
			},
			msg: types.NewMsgCreateRewardsPlan(
				1,
				"Rewards Plan",
				utils.MustParseCoins("100_000000service"),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				types.NewBasicPoolsDistribution(0),
				types.NewBasicOperatorsDistribution(0),
				types.NewBasicUsersDistribution(0),
				sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100_000_000))),
				testutil.TestAddress(10000).String(),
			),
			shouldErr: false,
			expResponse: &types.MsgCreateRewardsPlanResponse{
				NewRewardsPlanID: 1,
			},
			expEvents: []sdk.Event{
				sdk.NewEvent(
					types.EventTypeCreateRewardsPlan,
					sdk.NewAttribute(types.AttributeKeyRewardsPlanID, "1"),
					sdk.NewAttribute(servicestypes.AttributeKeyServiceID, "1"),
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the next plan id has been increased
				nextPlanID, err := suite.keeper.NextRewardsPlanID.Get(ctx)
				suite.Require().NoError(err)

				suite.Require().Equal(uint64(2), nextPlanID)

				// Make sure the rewards plan has been created
				_, err = suite.keeper.GetRewardsPlan(ctx, 1)
				suite.Require().NoError(err)

				// Make sure the balance is decreased by amount of the fee
				senderAddr, err := sdk.AccAddressFromBech32(testutil.TestAddress(10000).String())
				suite.Require().NoError(err)

				balances := suite.bankKeeper.GetAllBalances(ctx, senderAddr)
				suite.Require().Equal("400000000umilk", balances.String())
			},
		},
		{
			name: "service is created and fee is charged - one of many fees denoms",
			store: func(ctx sdk.Context) {
				// Create a service
				_, _ = suite.setupSampleServiceAndOperator(ctx)

				// Change rewards plan creation fee to 100 $MILK.
				err := suite.keeper.Params.Set(ctx, types.NewParams(
					sdk.NewCoins(
						sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000)),
						sdk.NewCoin("utia", sdkmath.NewInt(30_000_000)),
						sdk.NewCoin("milktia", sdkmath.NewInt(80_000_000)),
					),
				))
				suite.Require().NoError(err)

				// Set the next plan id
				err = suite.keeper.NextRewardsPlanID.Set(ctx, 1)
				suite.Require().NoError(err)

				// Fund the sender account enough coins to pay the fee.
				suite.FundAccount(
					ctx,
					testutil.TestAddress(10000).String(),
					sdk.NewCoins(
						sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000)),
						sdk.NewCoin("utia", sdkmath.NewInt(100_000_000)),
						sdk.NewCoin("milktia", sdkmath.NewInt(100_000_000)),
					),
				)
			},
			msg: types.NewMsgCreateRewardsPlan(
				1,
				"Rewards Plan",
				utils.MustParseCoins("100_000000service"),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				types.NewBasicPoolsDistribution(0),
				types.NewBasicOperatorsDistribution(0),
				types.NewBasicUsersDistribution(0),
				sdk.NewCoins(
					sdk.NewCoin("uatom", sdkmath.NewInt(20_000_000)),
					sdk.NewCoin("utia", sdkmath.NewInt(15_000_000)),
					sdk.NewCoin("milktia", sdkmath.NewInt(80_000_000)),
				),
				testutil.TestAddress(10000).String(),
			),
			shouldErr: false,
			expResponse: &types.MsgCreateRewardsPlanResponse{
				NewRewardsPlanID: 1,
			},
			expEvents: []sdk.Event{
				sdk.NewEvent(
					types.EventTypeCreateRewardsPlan,
					sdk.NewAttribute(types.AttributeKeyRewardsPlanID, "1"),
					sdk.NewAttribute(servicestypes.AttributeKeyServiceID, "1"),
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the next plan id has been increased
				nextPlanID, err := suite.keeper.NextRewardsPlanID.Get(ctx)
				suite.Require().NoError(err)

				suite.Require().Equal(uint64(2), nextPlanID)

				// Make sure the rewards plan has been created
				_, err = suite.keeper.GetRewardsPlan(ctx, 1)
				suite.Require().NoError(err)

				// Make sure the user's funds were deducted
				userAddress, err := sdk.AccAddressFromBech32(testutil.TestAddress(10000).String())
				suite.Require().NoError(err)
				balance := suite.bankKeeper.GetAllBalances(ctx, userAddress)
				suite.Require().Equal(sdk.NewCoins(
					sdk.NewCoin("uatom", sdkmath.NewInt(80_000_000)),
					sdk.NewCoin("utia", sdkmath.NewInt(85_000_000)),
					sdk.NewCoin("milktia", sdkmath.NewInt(20_000_000)),
				), balance)

				// Make sure the community pool was funded
				poolBalance := suite.bankKeeper.GetAllBalances(ctx, authtypes.NewModuleAddress(distrtypes.ModuleName))
				suite.Require().Equal(sdk.NewCoins(
					sdk.NewCoin("uatom", sdkmath.NewInt(20_000_000)),
					sdk.NewCoin("utia", sdkmath.NewInt(15_000_000)),
					sdk.NewCoin("milktia", sdkmath.NewInt(80_000_000)),
				), poolBalance)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.setup != nil {
				tc.setup()
			}
			if tc.setupCtx != nil {
				ctx = tc.setupCtx(ctx)
			}
			if tc.store != nil {
				tc.store(ctx)
			}

			msgServer := keeper.NewMsgServer(suite.keeper)
			res, err := msgServer.CreateRewardsPlan(ctx, tc.msg)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expResponse, res)
				for _, event := range tc.expEvents {
					suite.Require().Contains(ctx.EventManager().Events(), event)
				}

				if tc.check != nil {
					tc.check(ctx)
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgEditRewardsPlan() {
	testCases := []struct {
		name        string
		setup       func()
		store       func(ctx sdk.Context)
		setupCtx    func(ctx sdk.Context) sdk.Context
		msg         *types.MsgEditRewardsPlan
		shouldErr   bool
		expResponse *types.MsgEditRewardsPlanResponse
		expEvents   sdk.Events
		check       func(ctx sdk.Context)
	}{
		{
			name: "service not found returns error",
			msg: types.NewMsgEditRewardsPlan(
				1,
				"Rewards Plan",
				utils.MustParseCoins("100_000000service"),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				types.NewBasicPoolsDistribution(0),
				types.NewBasicOperatorsDistribution(0),
				types.NewBasicUsersDistribution(0),
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "sender different from service admin returns error",
			store: func(ctx sdk.Context) {
				// Create a service
				_, _ = suite.setupSampleServiceAndOperator(ctx)

				// Set the next plan id
				err := suite.keeper.NextRewardsPlanID.Set(ctx, 1)
				suite.Require().NoError(err)

				// Create a rewards plan
				_, err = suite.keeper.CreateRewardsPlan(
					ctx,
					"Rewards Plan",
					1,
					utils.MustParseCoins("100_000000service"),
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					types.NewBasicPoolsDistribution(0),
					types.NewBasicOperatorsDistribution(0),
					types.NewBasicUsersDistribution(0),
				)
				suite.Require().NoError(err)
			},
			msg: types.NewMsgEditRewardsPlan(
				1,
				"Rewards Plan",
				utils.MustParseCoins("100_000000service"),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				types.NewBasicPoolsDistribution(0),
				types.NewBasicOperatorsDistribution(0),
				types.NewBasicUsersDistribution(0),
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "end time before start time returns error",
			setupCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
			},
			store: func(ctx sdk.Context) {
				// Create a service
				_, _ = suite.setupSampleServiceAndOperator(ctx)

				// Set the next plan id
				err := suite.keeper.NextRewardsPlanID.Set(ctx, 1)
				suite.Require().NoError(err)

				// Create a rewards plan
				_, err = suite.keeper.CreateRewardsPlan(
					ctx,
					"Rewards Plan",
					1,
					utils.MustParseCoins("100_000000service"),
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					types.NewBasicPoolsDistribution(0),
					types.NewBasicOperatorsDistribution(0),
					types.NewBasicUsersDistribution(0),
				)
				suite.Require().NoError(err)
			},
			msg: types.NewMsgEditRewardsPlan(
				1,
				"Rewards Plan",
				utils.MustParseCoins("100_000000service"),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
				types.NewBasicPoolsDistribution(0),
				types.NewBasicOperatorsDistribution(0),
				types.NewBasicUsersDistribution(0),
				testutil.TestAddress(10000).String(),
			),
			shouldErr: true,
		},
		{
			name: "edit inactive rewards plan returns error",
			setupCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC))
			},
			store: func(ctx sdk.Context) {
				// Create a service
				_, _ = suite.setupSampleServiceAndOperator(ctx)

				// Set the next plan id
				err := suite.keeper.NextRewardsPlanID.Set(ctx, 1)
				suite.Require().NoError(err)

				// Create a rewards plan
				_, err = suite.keeper.CreateRewardsPlan(
					ctx,
					"Rewards Plan",
					1,
					utils.MustParseCoins("100_000000service"),
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					types.NewBasicPoolsDistribution(0),
					types.NewBasicOperatorsDistribution(0),
					types.NewBasicUsersDistribution(0),
				)
				suite.Require().NoError(err)
			},
			msg: types.NewMsgEditRewardsPlan(
				1,
				"Rewards Plan - Edited",
				utils.MustParseCoins("100_000000service"),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				types.NewBasicPoolsDistribution(0),
				types.NewBasicOperatorsDistribution(0),
				types.NewBasicUsersDistribution(0),
				testutil.TestAddress(10000).String(),
			),
			shouldErr: true,
		},
		{
			name: "edit rewards plan successfully",
			setupCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC))
			},
			store: func(ctx sdk.Context) {
				// Create a service
				_, _ = suite.setupSampleServiceAndOperator(ctx)

				// Set the next plan id
				err := suite.keeper.NextRewardsPlanID.Set(ctx, 1)
				suite.Require().NoError(err)

				// Create a rewards plan
				_, err = suite.keeper.CreateRewardsPlan(
					ctx,
					"Rewards Plan",
					1,
					utils.MustParseCoins("100_000000service"),
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					types.NewBasicPoolsDistribution(0),
					types.NewBasicOperatorsDistribution(0),
					types.NewBasicUsersDistribution(0),
				)
				suite.Require().NoError(err)
			},
			msg: types.NewMsgEditRewardsPlan(
				1,
				"Rewards Plan - Edited",
				utils.MustParseCoins("200_000000service"),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
				types.NewBasicPoolsDistribution(1),
				types.NewBasicOperatorsDistribution(2),
				types.NewBasicUsersDistribution(3),
				testutil.TestAddress(10000).String(),
			),
			shouldErr:   false,
			expResponse: &types.MsgEditRewardsPlanResponse{},
			expEvents: []sdk.Event{
				sdk.NewEvent(
					types.EventTypeEditRewardsPlan,
					sdk.NewAttribute(types.AttributeKeyRewardsPlanID, "1"),
					sdk.NewAttribute(servicestypes.AttributeKeyServiceID, "1"),
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the rewards plan has been edited
				plan, err := suite.keeper.GetRewardsPlan(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().Equal(utils.MustParseCoins("200_000000service"), plan.AmountPerDay)
				suite.Require().Equal(
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					plan.StartTime,
				)
				suite.Require().Equal(
					time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
					plan.EndTime,
				)
				// Check pools distribution
				poolsDistributionType, err := types.GetDistributionType(suite.cdc, plan.PoolsDistribution)
				suite.Require().IsType(&types.DistributionTypeBasic{}, poolsDistributionType)
				suite.Require().Equal(uint32(1), plan.PoolsDistribution.Weight)

				// Check operators distribution
				operatorsDistributionType, err := types.GetDistributionType(suite.cdc, plan.OperatorsDistribution)
				suite.Require().IsType(&types.DistributionTypeBasic{}, operatorsDistributionType)
				suite.Require().Equal(uint32(2), plan.OperatorsDistribution.Weight)

				// Check users distribution
				suite.Require().Equal(uint32(3), plan.UsersDistribution.Weight)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx := suite.ctx
			if tc.setup != nil {
				tc.setup()
			}
			if tc.setupCtx != nil {
				ctx = tc.setupCtx(ctx)
			}
			if tc.store != nil {
				tc.store(ctx)
			}

			msgServer := keeper.NewMsgServer(suite.keeper)
			res, err := msgServer.EditRewardsPlan(ctx, tc.msg)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expResponse, res)
				for _, event := range tc.expEvents {
					suite.Require().Contains(ctx.EventManager().Events(), event)
				}

				if tc.check != nil {
					tc.check(ctx)
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgSetWithdrawAddress() {
	testCases := []struct {
		name        string
		setup       func()
		store       func(ctx sdk.Context)
		setupCtx    func(ctx sdk.Context) sdk.Context
		msg         *types.MsgSetWithdrawAddress
		shouldErr   bool
		expResponse *types.MsgSetWithdrawAddressResponse
		expEvents   sdk.Events
		check       func(ctx sdk.Context)
	}{
		{
			name: "invalid sender address returns error",
			msg: types.NewMsgSetWithdrawAddress(
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				"invalid",
			),
			shouldErr: true,
		},
		{
			name: "invalid withdraw address returns error",
			msg: types.NewMsgSetWithdrawAddress(
				"invalid",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "success",
			msg: types.NewMsgSetWithdrawAddress(
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr:   false,
			expResponse: &types.MsgSetWithdrawAddressResponse{},
			expEvents: []sdk.Event{
				sdk.NewEvent(
					types.EventTypeSetWithdrawAddress,
					sdk.NewAttribute(sdk.AttributeKeySender, "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd"),
					sdk.NewAttribute(types.AttributeKeyWithdrawAddress, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the withdrawal address has been set
				delegatorAddr, err := sdk.AccAddressFromBech32("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd")
				suite.Require().NoError(err)

				withdrawAddr, err := suite.keeper.GetDelegatorWithdrawAddr(ctx, delegatorAddr)
				suite.Require().NoError(err)
				suite.Require().Equal("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4", withdrawAddr.String())
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.setup != nil {
				tc.setup()
			}
			if tc.setupCtx != nil {
				ctx = tc.setupCtx(ctx)
			}
			if tc.store != nil {
				tc.store(ctx)
			}

			msgServer := keeper.NewMsgServer(suite.keeper)
			res, err := msgServer.SetWithdrawAddress(ctx, tc.msg)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expResponse, res)
				for _, event := range tc.expEvents {
					suite.Require().Contains(ctx.EventManager().Events(), event)
				}

				if tc.check != nil {
					tc.check(ctx)
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgWithdrawDelegatorReward() {
	testCases := []struct {
		name        string
		setup       func()
		setupCtx    func(ctx sdk.Context) sdk.Context
		store       func(ctx sdk.Context)
		updateCtx   func(ctx sdk.Context) sdk.Context
		msg         *types.MsgWithdrawDelegatorReward
		shouldErr   bool
		expResponse *types.MsgWithdrawDelegatorRewardResponse
		expEvents   sdk.Events
		check       func(ctx sdk.Context)
	}{
		{
			name: "invalid delegator address returns error",
			msg: types.NewMsgWithdrawDelegatorReward(
				restakingtypes.DELEGATION_TYPE_SERVICE,
				1,
				"invalid",
			),
			shouldErr: true,
		},
		{
			name: "invalid delegation type returns error",
			msg: types.NewMsgWithdrawDelegatorReward(
				restakingtypes.DelegationType(5),
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "invalid target ID returns error",
			msg: types.NewMsgWithdrawDelegatorReward(
				restakingtypes.DELEGATION_TYPE_SERVICE,
				0,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "operator delegation not found returns error",
			msg: types.NewMsgWithdrawDelegatorReward(
				restakingtypes.DELEGATION_TYPE_OPERATOR,
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "service delegation not found returns error",
			msg: types.NewMsgWithdrawDelegatorReward(
				restakingtypes.DELEGATION_TYPE_SERVICE,
				3,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "found delegation rewards is withdrawn properly",
			store: func(ctx sdk.Context) {
				// Create a service and its reward plan
				service, _ := suite.setupSampleServiceAndOperator(ctx)
				suite.CreateBasicRewardsPlan(
					ctx,
					service.ID,
					utils.MustParseCoins("100_000000service"),
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					utils.MustParseCoins("10000_000000service"),
				)

				// Delegate some tokens to the service
				suite.DelegateService(
					ctx,
					1,
					utils.MustParseCoins("100_000000umilk"),
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					true,
				)
			},
			updateCtx: func(ctx sdk.Context) sdk.Context {
				// Allocate rewards
				return suite.allocateRewards(ctx, 10*time.Second)
			},
			msg: types.NewMsgWithdrawDelegatorReward(
				restakingtypes.DELEGATION_TYPE_SERVICE,
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: false,
			expResponse: &types.MsgWithdrawDelegatorRewardResponse{
				Amount: sdk.NewCoins(sdk.NewCoin("service", sdkmath.NewInt(11574))),
			},
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeWithdrawRewards,
					sdk.NewAttribute(types.AttributeKeyDelegationType, restakingtypes.DELEGATION_TYPE_SERVICE.String()),
					sdk.NewAttribute(types.AttributeKeyDelegationTargetID, "1"),
					sdk.NewAttribute(restakingtypes.AttributeKeyDelegator, "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd"),
					sdk.NewAttribute(sdk.AttributeKeyAmount, "11574service"),
					sdk.NewAttribute(types.AttributeKeyAmountPerPool, "denom:\"umilk\" coins:<denom:\"service\" amount:\"11574\" > "),
				),
			},
			check: func(ctx sdk.Context) {
				delegatorAddress, err := sdk.AccAddressFromBech32("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd")
				suite.Require().NoError(err)

				balances := suite.bankKeeper.GetAllBalances(ctx, delegatorAddress)
				suite.Require().Equal("11574service", balances.String())
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.setup != nil {
				tc.setup()
			}
			if tc.setupCtx != nil {
				ctx = tc.setupCtx(ctx)
			}
			if tc.store != nil {
				tc.store(ctx)
			}
			if tc.updateCtx != nil {
				ctx = tc.updateCtx(ctx)
			}

			msgServer := keeper.NewMsgServer(suite.keeper)
			res, err := msgServer.WithdrawDelegatorReward(ctx, tc.msg)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expResponse, res)
				for _, event := range tc.expEvents {
					suite.Require().Contains(ctx.EventManager().Events(), event)
				}

				if tc.check != nil {
					tc.check(ctx)
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgWithdrawOperatorCommission() {
	testCases := []struct {
		name        string
		setup       func()
		setupCtx    func(ctx sdk.Context) sdk.Context
		store       func(ctx sdk.Context)
		updateCtx   func(ctx sdk.Context) sdk.Context
		msg         *types.MsgWithdrawOperatorCommission
		shouldErr   bool
		expResponse *types.MsgWithdrawOperatorCommissionResponse
		expEvents   sdk.Events
		check       func(ctx sdk.Context)
	}{
		{
			name:      "invalid sender address returns error",
			msg:       types.NewMsgWithdrawOperatorCommission(1, "invalid"),
			shouldErr: true,
		},
		{
			name:      "operator not found returns error",
			msg:       types.NewMsgWithdrawOperatorCommission(1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			shouldErr: true,
		},
		{
			name: "non-admin address returns error",
			store: func(ctx sdk.Context) {
				// Create a service and a reward plan
				service, operator := suite.setupSampleServiceAndOperator(ctx)
				suite.CreateBasicRewardsPlan(
					ctx,
					service.ID,
					utils.MustParseCoins("100_000000service"),
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					utils.MustParseCoins("10000_000000service"),
				)

				// Delegate to an operator
				suite.DelegateOperator(
					ctx,
					operator.ID,
					utils.MustParseCoins("100_000000umilk"),
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					true,
				)
			},
			updateCtx: func(ctx sdk.Context) sdk.Context {
				// Allocate rewards
				return suite.allocateRewards(ctx, 10*time.Second)
			},
			msg: types.NewMsgWithdrawOperatorCommission(
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},

		{
			name: "existing operator commission is withdrawn properly",
			store: func(ctx sdk.Context) {
				// Create a service and its rewards plan
				service, operator := suite.setupSampleServiceAndOperator(ctx)
				suite.CreateBasicRewardsPlan(
					ctx,
					service.ID,
					utils.MustParseCoins("100_000000service"),
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					utils.MustParseCoins("10000_000000service"),
				)

				// Delegate to the operator
				suite.DelegateOperator(
					ctx,
					operator.ID,
					utils.MustParseCoins("100_000000umilk"),
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					true,
				)
			},
			updateCtx: func(ctx sdk.Context) sdk.Context {
				// Allocate rewards
				return suite.allocateRewards(ctx, 10*time.Second)
			},
			msg: types.NewMsgWithdrawOperatorCommission(1, testutil.TestAddress(10001).String()),
			expResponse: &types.MsgWithdrawOperatorCommissionResponse{
				Amount: sdk.NewCoins(sdk.NewCoin("service", sdkmath.NewInt(1157))),
			},
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeWithdrawCommission,
					sdk.NewAttribute(operatorstypes.AttributeKeyOperatorID, "1"),
					sdk.NewAttribute(sdk.AttributeKeyAmount, "1157service"),
					sdk.NewAttribute(types.AttributeKeyAmountPerPool, "denom:\"umilk\" coins:<denom:\"service\" amount:\"1157\" > "),
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the funds have been sent to the admin
				adminAddress := testutil.TestAddress(10001)
				balances := suite.bankKeeper.GetAllBalances(ctx, adminAddress)
				suite.Require().Equal("1157service", balances.String())
			},
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.setup != nil {
				tc.setup()
			}
			if tc.setupCtx != nil {
				ctx = tc.setupCtx(ctx)
			}
			if tc.store != nil {
				tc.store(ctx)
			}
			if tc.updateCtx != nil {
				ctx = tc.updateCtx(ctx)
			}

			msgServer := keeper.NewMsgServer(suite.keeper)
			res, err := msgServer.WithdrawOperatorCommission(ctx, tc.msg)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expResponse, res)
				for _, event := range tc.expEvents {
					suite.Require().Contains(ctx.EventManager().Events(), event)
				}

				if tc.check != nil {
					tc.check(ctx)
				}
			}
		})
	}
}
