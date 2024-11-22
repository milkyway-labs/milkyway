package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

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
			store: func(ctx sdk.Context) {
				err := suite.k.SetNextOperatorID(ctx, 1)
				suite.Require().NoError(err)

				err = suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000))),
					6*time.Hour,
				))
				suite.Require().NoError(err)
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
				err := suite.k.SetNextOperatorID(ctx, 2)
				suite.Require().NoError(err)

				err = suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000))),
					6*time.Hour,
				))
				suite.Require().NoError(err)

				// Send funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(200_000_000))),
				)
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
				stored, found, err := suite.k.GetOperator(ctx, 2)
				suite.Require().NoError(err)
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
				params, err := suite.k.GetOperatorParams(ctx, 2)
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
				poolBalance := suite.bk.GetBalance(ctx, authtypes.NewModuleAddress(distrtypes.ModuleName), "uatom")
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
			res, err := msgServer.RegisterOperator(ctx, tc.msg)
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
				err := suite.k.SaveOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)
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
				err := suite.k.SaveOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)
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
				err := suite.k.SaveOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)
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
				stored, found, err := suite.k.GetOperator(ctx, 1)
				suite.Require().NoError(err)
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
			res, err := msgServer.UpdateOperator(ctx, tc.msg)
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
				err := suite.k.SaveOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				))
				suite.Require().NoError(err)
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
				err := suite.k.SaveOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_INACTIVATING,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)
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
				err := suite.k.SaveOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_INACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)
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
				err := suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000))),
					6*time.Hour,
				))
				suite.Require().NoError(err)

				err = suite.k.SaveOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)
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
				stored, found, err := suite.k.GetOperator(ctx, 1)
				suite.Require().NoError(err)
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
			res, err := msgServer.DeactivateOperator(ctx, tc.msg)
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

func (suite *KeeperTestSuite) TestMsgServer_ReactivateOperator() {
	testCases := []struct {
		name        string
		setup       func()
		store       func(ctx sdk.Context)
		msg         *types.MsgReactivateOperator
		shouldErr   bool
		expResponse *types.MsgReactivateOperatorResponse
		expEvents   sdk.Events
		check       func(ctx sdk.Context)
	}{
		{
			name:      "not found operator returns error",
			msg:       types.NewMsgReactivateOperator(1, "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd"),
			shouldErr: true,
		},
		{
			name: "non admin sender returns error",
			store: func(ctx sdk.Context) {
				err := suite.k.SaveOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_INACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgReactivateOperator(
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "active operator can't be reactivated",
			store: func(ctx sdk.Context) {
				err := suite.k.SaveOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgReactivateOperator(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "inactivating operator can't be reactivated",
			store: func(ctx sdk.Context) {
				err := suite.k.SaveOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_INACTIVATING,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgReactivateOperator(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "operator reactivated successfully",
			store: func(ctx sdk.Context) {
				err := suite.k.SaveOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_INACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgReactivateOperator(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr:   false,
			expResponse: &types.MsgReactivateOperatorResponse{},
			expEvents: []sdk.Event{
				sdk.NewEvent(
					types.EventTypeReactivateOperator,
					sdk.NewAttribute(types.AttributeKeyOperatorID, "1"),
				),
			},
			check: func(ctx sdk.Context) {
				operator, found, err := suite.k.GetOperator(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().Equal(types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				), operator)

				// Check hook called
				called := suite.hooks.CalledMap["AfterOperatorReactivated"]
				suite.Require().True(called)
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
			if tc.store != nil {
				tc.store(ctx)
			}

			msgServer := keeper.NewMsgServer(suite.k)
			res, err := msgServer.ReactivateOperator(ctx, tc.msg)
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
				err := suite.k.SaveOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				))
				suite.Require().NoError(err)
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
				err := suite.k.SaveOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)
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
				stored, found, err := suite.k.GetOperator(ctx, 1)
				suite.Require().NoError(err)
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
			res, err := msgServer.TransferOperatorOwnership(ctx, tc.msg)
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

func (suite *KeeperTestSuite) TestMsgServer_DeleteOperator() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		msg       *types.MsgDeleteOperator
		shouldErr bool
		expEvents sdk.Events
		check     func(ctx sdk.Context)
	}{
		{
			name: "not found operator returns error",
			msg: types.NewMsgDeleteOperator(
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "non admin sender returns error",
			store: func(ctx sdk.Context) {
				err := suite.k.SaveOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_INACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgDeleteOperator(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "active operator returns error",
			store: func(ctx sdk.Context) {
				err := suite.k.SaveOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgDeleteOperator(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "inactivating operator returns error",
			store: func(ctx sdk.Context) {
				err := suite.k.SaveOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_INACTIVATING,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgDeleteOperator(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "operator deleted successfully",
			store: func(ctx sdk.Context) {
				err := suite.k.SaveOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_INACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgDeleteOperator(
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: false,
			expEvents: []sdk.Event{
				sdk.NewEvent(
					types.EventTypeDeleteOperator,
					sdk.NewAttribute(types.AttributeKeyOperatorID, "1"),
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the operator was updated
				_, found, err := suite.k.GetOperator(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().False(found)

				// Ensure the hook has been called
				suite.Require().True(suite.hooks.CalledMap["BeforeOperatorDeleted"])
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			msgServer := keeper.NewMsgServer(suite.k)
			res, err := msgServer.DeleteOperator(ctx, tc.msg)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
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
			store: func(ctx sdk.Context) {
				// Register a test operator
				err := suite.k.CreateOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					operatorAdmin,
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgSetOperatorParams(
				1,
				types.NewOperatorParams(sdkmath.LegacyNewDec(-1)),
				operatorAdmin,
			),
			shouldErr: true,
		},
		{
			name: "not admin can't set params",
			store: func(ctx sdk.Context) {
				// Register a test operator
				err := suite.k.CreateOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					operatorAdmin,
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgSetOperatorParams(
				1,
				types.NewOperatorParams(sdkmath.LegacyMustNewDecFromStr("0.2")),
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
			),
			shouldErr: true,
		},
		{
			name: "set params for not existing operator fails",
			store: func(ctx sdk.Context) {
				// Register a test operator
				err := suite.k.CreateOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					operatorAdmin,
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgSetOperatorParams(
				3,
				types.NewOperatorParams(sdkmath.LegacyMustNewDecFromStr("0.2")),
				operatorAdmin,
			),
			shouldErr: true,
		},
		{
			name: "set params works properly",
			store: func(ctx sdk.Context) {
				// Register a test operator
				err := suite.k.CreateOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					operatorAdmin,
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgSetOperatorParams(
				1,
				types.NewOperatorParams(sdkmath.LegacyMustNewDecFromStr("0.2")),
				operatorAdmin,
			),
			expEvents: []sdk.Event{
				sdk.NewEvent(
					types.EventTypeSetOperatorParams,
				),
			},
			shouldErr: false,
			check: func(ctx sdk.Context) {
				params, err := suite.k.GetOperatorParams(ctx, 1)
				suite.Require().Nil(err)
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

			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

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
				params, err := suite.k.GetParams(ctx)
				suite.Require().NoError(err)
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
			res, err := msgServer.UpdateParams(ctx, tc.msg)
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
