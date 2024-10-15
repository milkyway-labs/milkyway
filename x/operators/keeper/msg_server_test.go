package keeper_test

import (
	"fmt"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/milkyway-labs/milkyway/x/operators/keeper"
	"github.com/milkyway-labs/milkyway/x/operators/types"
)

func (suite *KeeperTestSuite) TestMsgServer_RegisterOperator() {
	testCases := []struct {
		name        string
		setup       func()
		store       func(ctx sdk.Context)
		setupCtx    func(ctx sdk.Context) sdk.Context
		msg         *types.MsgRegisterOperator
		shouldErr   bool
		expResponse *types.MsgRegisterOperatorResponse
		expEvents   sdk.Events
		check       func(ctx sdk.Context)
	}{
		{
			name: "non existing next operator id returns error",
			msg: types.NewMsgRegisterOperator(
				"MilkyWay Operator",
				"https://milkyway.com",
				"https://milkyway.com/picture",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "invalid operator returns error",
			setup: func() {
				suite.k.SetNextOperatorID(suite.ctx, 1)
				suite.k.SetParams(suite.ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000))),
					6*time.Hour,
				))
			},
			msg: types.NewMsgRegisterOperator(
				types.DoNotModify,
				"https://milkyway.com",
				"https://milkyway.com/picture",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "operator registered successfully",
			store: func(ctx sdk.Context) {
				suite.k.SetNextOperatorID(suite.ctx, 2)
				suite.k.SetParams(suite.ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000))),
					6*time.Hour,
				))

				// Send funds to the user
				suite.fundAccount(ctx, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4", sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(200_000_000))))
			},
			msg: types.NewMsgRegisterOperator(
				"MilkyWay Operator",
				"https://milkyway.com",
				"https://milkyway.com/picture",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: false,
			expResponse: &types.MsgRegisterOperatorResponse{
				NewOperatorID: 2,
			},
			expEvents: []sdk.Event{
				sdk.NewEvent(
					types.EventTypeRegisterOperator,
					sdk.NewAttribute(types.AttributeKeyOperatorID, "2"),
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the operator was stored
				stored, found := suite.k.GetOperator(ctx, 2)
				suite.Require().True(found)
				suite.Require().Equal(types.NewOperator(
					2,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				), stored)

				// Make sure the operator account has been created
				hasAccount := suite.ak.HasAccount(ctx, types.GetOperatorAddress(2))
				suite.Require().True(hasAccount)

				// Make sure the newly registered operator has the default params
				params, found, err := suite.k.GetOperatorParams(ctx, 2)
				suite.Require().True(found)
				suite.Require().NoError(err)
				suite.Require().Equal(types.DefaultOperatorParams(), params)

				// Make sure the next operator id has incremented
				nextID, err := suite.k.GetNextOperatorID(ctx)
				suite.Require().NoError(err)
				suite.Require().Equal(uint32(3), nextID)

				// Make sure the user's funds were deducted
				userAddress, err := sdk.AccAddressFromBech32("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().NoError(err)
				balance := suite.bk.GetBalance(ctx, userAddress, "uatom")
				suite.Require().Equal(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000)), balance)

				// Make sure the community pool was funded
				poolBalance := suite.bk.GetBalance(ctx, authtypes.NewModuleAddress(authtypes.FeeCollectorName), "uatom")
				suite.Require().Equal(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000)), poolBalance)
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

			msgServer := keeper.NewMsgServer(suite.k)
			res, err := msgServer.RegisterOperator(sdk.WrapSDKContext(ctx), tc.msg)
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

func (suite *KeeperTestSuite) TestMsgServer_UpdateOperator() {
	testCases := []struct {
		name        string
		setup       func()
		store       func(ctx sdk.Context)
		setupCtx    func(ctx sdk.Context) sdk.Context
		msg         *types.MsgUpdateOperator
		shouldErr   bool
		expResponse *types.MsgUpdateOperatorResponse
		expEvents   sdk.Events
		check       func(ctx sdk.Context)
	}{
		{
			name: "not found operator returns error",
			msg: types.NewMsgUpdateOperator(
				1,
				"MilkyWay Operator",
				"https://milkyway.com",
				"https://milkyway.com/picture",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "non admin sender returns error",
			store: func(ctx sdk.Context) {
				suite.k.SaveOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
			},
			msg: types.NewMsgUpdateOperator(
				1,
				"MilkyWay Operator",
				"https://milkyway.com",
				"https://milkyway.com/picture",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "invalid operator returns error",
			store: func(ctx sdk.Context) {
				suite.k.SaveOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
			},
			msg: types.NewMsgUpdateOperator(
				1,
				"MilkyWay Operator",
				"https://milkyway.com",
				"https://milkyway.com/picture",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "operator updated successfully",
			store: func(ctx sdk.Context) {
				suite.k.SaveOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
			},
			msg: types.NewMsgUpdateOperator(
				1,
				"MilkyWay Updated Operator",
				"https://milkyway.zone",
				"https://milkyway.zone/picture",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr:   false,
			expResponse: &types.MsgUpdateOperatorResponse{},
			expEvents: []sdk.Event{
				sdk.NewEvent(
					types.EventTypeUpdateOperator,
					sdk.NewAttribute(types.AttributeKeyOperatorID, "1"),
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the operator was updated
				stored, found := suite.k.GetOperator(ctx, 1)
				suite.Require().True(found)
				suite.Require().Equal(types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Updated Operator",
					"https://milkyway.zone",
					"https://milkyway.zone/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				), stored)
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

			msgServer := keeper.NewMsgServer(suite.k)
			res, err := msgServer.UpdateOperator(sdk.WrapSDKContext(ctx), tc.msg)
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

func (suite *KeeperTestSuite) TestMsgServer_DeactivateOperator() {
	testCases := []struct {
		name        string
		setup       func()
		store       func(ctx sdk.Context)
		setupCtx    func(ctx sdk.Context) sdk.Context
		msg         *types.MsgDeactivateOperator
		shouldErr   bool
		expResponse *types.MsgDeactivateOperatorResponse
		expEvents   sdk.Events
		check       func(ctx sdk.Context)
	}{
		{
			name: "not found operator returns error",
			msg: types.NewMsgDeactivateOperator(
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "non admin sender returns error",
			store: func(ctx sdk.Context) {
				suite.k.SaveOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				))
			},
			msg: types.NewMsgDeactivateOperator(
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "already inactivating operator returns error",
			store: func(ctx sdk.Context) {
				suite.k.SaveOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_INACTIVATING,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
			},
			msg: types.NewMsgDeactivateOperator(
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "already inactive operator returns error",
			store: func(ctx sdk.Context) {
				suite.k.SaveOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_INACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
			},
			msg: types.NewMsgDeactivateOperator(
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "operator inactivation started successfully",
			setupCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
			},
			store: func(ctx sdk.Context) {
				suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000))),
					6*time.Hour,
				))
				suite.k.SaveOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
			},
			msg: types.NewMsgDeactivateOperator(
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr:   false,
			expResponse: &types.MsgDeactivateOperatorResponse{},
			expEvents: []sdk.Event{
				sdk.NewEvent(
					types.EventTypeStartOperatorInactivation,
					sdk.NewAttribute(types.AttributeKeyOperatorID, "1"),
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the operator was updated
				stored, found := suite.k.GetOperator(ctx, 1)
				suite.Require().True(found)
				suite.Require().Equal(types.NewOperator(
					1,
					types.OPERATOR_STATUS_INACTIVATING,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				), stored)

				// Make sure the operator was added to the inactivating queue
				inactivatingOperators, _ := suite.k.GetInactivatingOperators(ctx)
				suite.Require().Len(inactivatingOperators, 1)
				suite.Require().Equal(types.NewUnbondingOperator(
					1,
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Add(6*time.Hour),
				), inactivatingOperators[0])
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

			msgServer := keeper.NewMsgServer(suite.k)
			res, err := msgServer.DeactivateOperator(sdk.WrapSDKContext(ctx), tc.msg)
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

func (suite *KeeperTestSuite) TestMsgServer_TransferOperatorOwnership() {
	testCases := []struct {
		name        string
		setup       func()
		store       func(ctx sdk.Context)
		setupCtx    func(ctx sdk.Context) sdk.Context
		msg         *types.MsgTransferOperatorOwnership
		shouldErr   bool
		expResponse *types.MsgTransferOperatorOwnershipResponse
		expEvents   sdk.Events
		check       func(ctx sdk.Context)
	}{
		{
			name: "not found operator returns error",
			msg: types.NewMsgTransferOperatorOwnership(
				1,
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "non admin sender returns error",
			store: func(ctx sdk.Context) {
				suite.k.SaveOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				))
			},
			msg: types.NewMsgTransferOperatorOwnership(
				1,
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "operator ownership transferred successfully",
			store: func(ctx sdk.Context) {
				suite.k.SaveOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
			},
			msg: types.NewMsgTransferOperatorOwnership(
				1,
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr:   false,
			expResponse: &types.MsgTransferOperatorOwnershipResponse{},
			expEvents: []sdk.Event{
				sdk.NewEvent(
					types.EventTypeTransferOperatorOwnership,
					sdk.NewAttribute(types.AttributeKeyOperatorID, "1"),
					sdk.NewAttribute(types.AttributeKeyNewAdmin, "cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn"),
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the operator was updated
				stored, found := suite.k.GetOperator(ctx, 1)
				suite.Require().True(found)
				suite.Require().Equal(types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
				), stored)
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

			msgServer := keeper.NewMsgServer(suite.k)
			res, err := msgServer.TransferOperatorOwnership(sdk.WrapSDKContext(ctx), tc.msg)
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

func (suite *KeeperTestSuite) TestMsgServer_SetOperatorParams() {
	testOperatorId := uint32(2)
	operatorAdmin := "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"

	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		msg       *types.MsgSetOperatorParams
		shouldErr bool
		expEvents sdk.Events
		check     func(ctx sdk.Context)
	}{
		{
			name: "set invalid params fails",
			msg: types.NewMsgSetOperatorParams(
				operatorAdmin,
				testOperatorId,
				sdkmath.LegacyNewDec(-1),
			),
			shouldErr: true,
		},
		{
			name: "not admin can't set params",
			msg: types.NewMsgSetOperatorParams(
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
				testOperatorId,
				sdkmath.LegacyMustNewDecFromStr("0.2"),
			),
			shouldErr: true,
		},
		{
			name: "set params for not existing operator fails",
			msg: types.NewMsgSetOperatorParams(
				operatorAdmin,
				3,
				sdkmath.LegacyMustNewDecFromStr("0.2"),
			),
			shouldErr: true,
		},
		{
			name: "set params works properly",
			msg: types.NewMsgSetOperatorParams(
				operatorAdmin,
				testOperatorId,
				sdkmath.LegacyMustNewDecFromStr("0.2"),
			),
			expEvents: []sdk.Event{
				sdk.NewEvent(
					types.EventTypeSetOperatorParams,
					sdk.NewAttribute(types.AttributeKeyOperatorID, fmt.Sprintf("%d", testOperatorId)),
					sdk.NewAttribute(types.AttributeKeyNewCommissionRate, "0.200000000000000000"),
				),
			},
			shouldErr: false,
			check: func(ctx sdk.Context) {
				params, found, err := suite.k.GetOperatorParams(ctx, testOperatorId)
				suite.Require().Nil(err)
				suite.Require().True(found)
				suite.Require().Equal(types.NewOperatorParams(
					sdkmath.LegacyMustNewDecFromStr("0.2"),
				), params)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx := suite.ctx

			if tc.store != nil {
				tc.store(ctx)
			}

			// Register a test operator
			err := suite.k.RegisterOperator(ctx, types.NewOperator(
				testOperatorId,
				types.OPERATOR_STATUS_ACTIVE,
				"MilkyWay Operator",
				"https://milkyway.com",
				"https://milkyway.com/picture",
				operatorAdmin,
			))
			suite.Require().NoError(err)

			msgServer := keeper.NewMsgServer(suite.k)
			res, err := msgServer.SetOperatorParams(ctx, tc.msg)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				for _, event := range tc.expEvents {
					suite.Assert().Contains(ctx.EventManager().Events(), event)
				}

				if tc.check != nil {
					tc.check(ctx)
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgServer_UpdateParams() {
	testCases := []struct {
		name        string
		setup       func()
		store       func(ctx sdk.Context)
		setupCtx    func(ctx sdk.Context) sdk.Context
		msg         *types.MsgUpdateParams
		shouldErr   bool
		expResponse *types.MsgUpdateParamsResponse
		expEvents   sdk.Events
		check       func(ctx sdk.Context)
	}{
		{
			name: "invalid authority return error",
			msg: types.NewMsgUpdateParams(
				types.DefaultParams(),
				"invalid",
			),
			shouldErr: true,
		},
		{
			name: "valid data returns no error",
			msg: types.NewMsgUpdateParams(
				types.DefaultParams(),
				authtypes.NewModuleAddress("gov").String(),
			),
			shouldErr:   false,
			expResponse: &types.MsgUpdateParamsResponse{},
			expEvents:   sdk.Events{},
			check: func(ctx sdk.Context) {
				params := suite.k.GetParams(ctx)
				suite.Require().Equal(types.DefaultParams(), params)
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

			msgServer := keeper.NewMsgServer(suite.k)
			res, err := msgServer.UpdateParams(sdk.WrapSDKContext(ctx), tc.msg)
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
