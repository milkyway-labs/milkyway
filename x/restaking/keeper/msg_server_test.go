package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	operatorstypes "github.com/milkyway-labs/milkyway/v9/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/v9/x/pools/types"
	"github.com/milkyway-labs/milkyway/v9/x/restaking/keeper"
	"github.com/milkyway-labs/milkyway/v9/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v9/x/services/types"
)

func (suite *KeeperTestSuite) TestMsgServer_JoinService() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		setupCtx  func(ctx sdk.Context) sdk.Context
		msg       *types.MsgJoinService
		shouldErr bool
		expEvents sdk.Events
		check     func(ctx sdk.Context)
	}{
		{
			name: "non-existent operator id returns an error",
			msg: types.NewMsgJoinService(
				1,
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "non-existent service id returns an error",
			store: func(ctx sdk.Context) {
				err := suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
					1,
					operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgJoinService(
				1,
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "only operator admin can perform JoinService",
			store: func(ctx sdk.Context) {
				err := suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
					1,
					operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)

				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgJoinService(
				1,
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "service not active returns error",
			store: func(ctx sdk.Context) {
				err := suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
					1,
					operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)

				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_INACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgJoinService(
				1,
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "can't join inactive service",
			store: func(ctx sdk.Context) {
				err := suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
					1,
					operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)
				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_INACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
			},
			msg: &types.MsgJoinService{
				Sender:     "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				OperatorID: 1,
				ServiceID:  1,
			},
			shouldErr: true,
		},
		{
			name: "can't join created service",
			store: func(ctx sdk.Context) {
				err := suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
					1,
					operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)
				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_INACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
			},
			msg: &types.MsgJoinService{
				Sender:     "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				OperatorID: 1,
				ServiceID:  1,
			},
			shouldErr: true,
		},
		{
			name: "can't join service due to allow list",
			store: func(ctx sdk.Context) {
				err := suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
					1,
					operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)

				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)

				err = suite.k.AddOperatorToServiceAllowList(ctx, 1, 2)
				suite.Require().NoError(err)
			},
			msg: types.NewMsgJoinService(
				1,
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "join is performed correctly - no allow list",
			store: func(ctx sdk.Context) {
				err := suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
					1,
					operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)

				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgJoinService(
				1,
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: false,
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeJoinService,
					sdk.NewAttribute(operatorstypes.AttributeKeyOperatorID, "1"),
					sdk.NewAttribute(servicestypes.AttributeKeyServiceID, "1"),
				),
			},
		},
		{
			name: "join is performed correctly - allow list",
			store: func(ctx sdk.Context) {
				err := suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
					1,
					operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)

				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)

				err = suite.k.AddOperatorToServiceAllowList(ctx, 1, 1)
				suite.Require().NoError(err)
			},
			msg: types.NewMsgJoinService(
				1,
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: false,
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeJoinService,
					sdk.NewAttribute(operatorstypes.AttributeKeyOperatorID, "1"),
					sdk.NewAttribute(servicestypes.AttributeKeyServiceID, "1"),
				),
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
			_, err := msgServer.JoinService(ctx, tc.msg)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
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

func (suite *KeeperTestSuite) TestMsgServer_LeaveService() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		setupCtx  func(ctx sdk.Context) sdk.Context
		msg       *types.MsgLeaveService
		shouldErr bool
		expEvents sdk.Events
		check     func(ctx sdk.Context)
	}{
		{
			name: "non-existent operator id returns an error",
			msg: types.NewMsgLeaveService(
				1,
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "non-existent service id returns an error",
			store: func(ctx sdk.Context) {
				err := suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
					1, operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgLeaveService(
				1,
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "only operator admin can perform LeaveService",
			store: func(ctx sdk.Context) {
				err := suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
					1, operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)

				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				err = suite.k.AddServiceToOperatorJoinedServices(ctx, 1, 1)
				suite.Require().NoError(err)
			},
			msg: types.NewMsgLeaveService(
				1,
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "valid update is done properly",
			store: func(ctx sdk.Context) {
				err := suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
					1, operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)

				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)

				err = suite.k.AddServiceToOperatorJoinedServices(ctx, 1, 1)
				suite.Require().NoError(err)
			},
			msg: types.NewMsgLeaveService(
				1,
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			check: func(ctx sdk.Context) {
				joinedServices, err := suite.k.HasOperatorJoinedService(ctx, 1, 1)
				suite.Require().NoError(err)
				suite.Require().False(joinedServices)
			},
			shouldErr: false,
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeLeaveService,
					sdk.NewAttribute(operatorstypes.AttributeKeyOperatorID, "1"),
					sdk.NewAttribute(servicestypes.AttributeKeyServiceID, "1"),
				),
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
			_, err := msgServer.LeaveService(ctx, tc.msg)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
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

func (suite *KeeperTestSuite) TestMsgServer_AddOperatorToAllowList() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		msg       *types.MsgAddOperatorToAllowList
		shouldErr bool
		expEvents sdk.Events
		check     func(ctx sdk.Context)
	}{
		{
			name: "non existing service returns error",
			store: func(ctx sdk.Context) {
				err := suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
					1, operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgAddOperatorToAllowList(
				1,
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "non existing operator returns error",
			store: func(ctx sdk.Context) {
				err := suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					false,
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgAddOperatorToAllowList(
				1,
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "only service admin can allow an operator",
			store: func(ctx sdk.Context) {
				err := suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
					1, operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)

				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					false,
				))
				suite.Require().NoError(err)
			},
			msg:       types.NewMsgAddOperatorToAllowList(1, 1, "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd"),
			shouldErr: true,
		},
		{
			name: "allow already allowed operator fails",
			store: func(ctx sdk.Context) {
				err := suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
					1, operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)

				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					false,
				))
				suite.Require().NoError(err)

				err = suite.k.AddOperatorToServiceAllowList(ctx, 1, 1)
				suite.Require().NoError(err)
			},
			msg: types.NewMsgAddOperatorToAllowList(
				1,
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "operator is allowed properly",
			store: func(ctx sdk.Context) {
				err := suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
					1, operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)

				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					false,
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgAddOperatorToAllowList(
				1,
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				configured, err := suite.k.IsServiceOperatorsAllowListConfigured(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().True(configured)
				whitelisted, err := suite.k.CanOperatorValidateService(ctx, 1, 1)
				suite.Require().NoError(err)
				suite.Require().True(whitelisted)
			},
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeAllowOperator,
					sdk.NewAttribute(operatorstypes.AttributeKeyOperatorID, "1"),
					sdk.NewAttribute(servicestypes.AttributeKeyServiceID, "1"),
				),
			},
		},
		{
			name: "adding the first operator to allow list makes other operators to leave",
			store: func(ctx sdk.Context) {
				err := suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
					1, operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)
				err = suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
					2, operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator 2",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)

				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					false,
				))
				suite.Require().NoError(err)

				err = suite.k.AddServiceToOperatorJoinedServices(ctx, 1, 1)
				err = suite.k.AddServiceToOperatorJoinedServices(ctx, 2, 1)
				suite.Require().NoError(err)
			},
			msg: types.NewMsgAddOperatorToAllowList(
				1,
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				configured, err := suite.k.IsServiceOperatorsAllowListConfigured(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().True(configured)
				whitelisted, err := suite.k.CanOperatorValidateService(ctx, 1, 2)
				suite.Require().NoError(err)
				suite.Require().False(whitelisted)
				whitelisted, err = suite.k.CanOperatorValidateService(ctx, 1, 1)
				suite.Require().NoError(err)
				suite.Require().True(whitelisted)
				joined, err := suite.k.HasOperatorJoinedService(ctx, 1, 1)
				suite.Require().NoError(err)
				suite.Require().True(joined)
				joined, err = suite.k.HasOperatorJoinedService(ctx, 1, 2)
				suite.Require().NoError(err)
				suite.Require().False(joined)
			},
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeAllowOperator,
					sdk.NewAttribute(operatorstypes.AttributeKeyOperatorID, "1"),
					sdk.NewAttribute(servicestypes.AttributeKeyServiceID, "1"),
				),
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

			msgServer := keeper.NewMsgServer(suite.k)
			_, err := msgServer.AddOperatorToAllowList(ctx, tc.msg)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				for _, event := range tc.expEvents {
					events := ctx.EventManager().Events()
					suite.Require().Contains(events, event)
				}

				if tc.check != nil {
					tc.check(ctx)
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgServer_RemoveAllowedOperator() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		msg       *types.MsgRemoveOperatorFromAllowlist
		shouldErr bool
		expEvents sdk.Events
		check     func(ctx sdk.Context)
	}{
		{
			name: "non existing service returns error",
			store: func(ctx sdk.Context) {
				err := suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
					1, operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgRemoveOperatorFromAllowList(
				1,
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "only service admin can remove an operator",
			store: func(ctx sdk.Context) {
				err := suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
					1, operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)

				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					false,
				))
				suite.Require().NoError(err)

				err = suite.k.AddOperatorToServiceAllowList(ctx, 1, 1)
				suite.Require().NoError(err)
			},
			msg:       types.NewMsgRemoveOperatorFromAllowList(1, 1, "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd"),
			shouldErr: true,
		},
		{
			name: "remove not allowed operator fails",
			store: func(ctx sdk.Context) {
				err := suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
					1, operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)

				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					false,
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgRemoveOperatorFromAllowList(
				1,
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "operator is removed properly",
			store: func(ctx sdk.Context) {
				err := suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
					1, operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)

				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					false,
				))
				suite.Require().NoError(err)

				err = suite.k.AddOperatorToServiceAllowList(ctx, 1, 1)
				suite.Require().NoError(err)
				err = suite.k.AddOperatorToServiceAllowList(ctx, 1, 2)
				suite.Require().NoError(err)
			},
			msg:       types.NewMsgRemoveOperatorFromAllowList(1, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				canValidate, err := suite.k.CanOperatorValidateService(ctx, 1, 1)
				suite.Require().NoError(err)
				suite.Require().False(canValidate)
			},
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeRemoveAllowedOperator,
					sdk.NewAttribute(operatorstypes.AttributeKeyOperatorID, "1"),
					sdk.NewAttribute(servicestypes.AttributeKeyServiceID, "1"),
				),
			},
		},
		{
			name: "operator is removed properly, leaves the service automatically",
			store: func(ctx sdk.Context) {
				err := suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
					1, operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)

				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					false,
				))
				suite.Require().NoError(err)

				err = suite.k.AddOperatorToServiceAllowList(ctx, 1, 1)
				suite.Require().NoError(err)
				err = suite.k.AddOperatorToServiceAllowList(ctx, 1, 2)
				suite.Require().NoError(err)

				err = suite.k.AddServiceToOperatorJoinedServices(ctx, 1, 1)
				suite.Require().NoError(err)
			},
			msg:       types.NewMsgRemoveOperatorFromAllowList(1, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				canValidate, err := suite.k.CanOperatorValidateService(ctx, 1, 1)
				suite.Require().NoError(err)
				suite.Require().False(canValidate)

				joined, err := suite.k.HasOperatorJoinedService(ctx, 1, 1)
				suite.Require().NoError(err)
				suite.Require().False(joined)
			},
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeRemoveAllowedOperator,
					sdk.NewAttribute(operatorstypes.AttributeKeyOperatorID, "1"),
					sdk.NewAttribute(servicestypes.AttributeKeyServiceID, "1"),
				),
			},
		},
		{
			name: "operator is removed and service allows all operators so it doesn't leave the service",
			store: func(ctx sdk.Context) {
				err := suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
					1, operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)

				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					false,
				))
				suite.Require().NoError(err)

				err = suite.k.AddOperatorToServiceAllowList(ctx, 1, 1)
				suite.Require().NoError(err)

				err = suite.k.AddServiceToOperatorJoinedServices(ctx, 1, 1)
				suite.Require().NoError(err)
			},
			msg:       types.NewMsgRemoveOperatorFromAllowList(1, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				canValidate, err := suite.k.CanOperatorValidateService(ctx, 1, 1)
				suite.Require().NoError(err)
				suite.Require().True(canValidate)

				joined, err := suite.k.HasOperatorJoinedService(ctx, 1, 1)
				suite.Require().NoError(err)
				suite.Require().True(joined)
			},
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeRemoveAllowedOperator,
					sdk.NewAttribute(operatorstypes.AttributeKeyOperatorID, "1"),
					sdk.NewAttribute(servicestypes.AttributeKeyServiceID, "1"),
				),
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

			msgServer := keeper.NewMsgServer(suite.k)
			_, err := msgServer.RemoveOperatorFromAllowlist(ctx, tc.msg)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				for _, event := range tc.expEvents {
					events := ctx.EventManager().Events()
					suite.Require().Contains(events, event)
				}

				if tc.check != nil {
					tc.check(ctx)
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgServer_BorrowPoolSecurity() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		msg       *types.MsgBorrowPoolSecurity
		shouldErr bool
		expEvents sdk.Events
		check     func(ctx sdk.Context)
	}{
		{
			name: "non existing service returns error",
			store: func(ctx sdk.Context) {
				err := suite.pk.SavePool(ctx, poolstypes.NewPool(1, "utia"))
				suite.Require().NoError(err)
			},
			msg:       types.NewMsgBorrowPoolSecurity(1, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			shouldErr: true,
		},
		{
			name: "non existing pool returns error",
			store: func(ctx sdk.Context) {
				err := suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					false,
				))
				suite.Require().NoError(err)
			},
			msg:       types.NewMsgBorrowPoolSecurity(1, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			shouldErr: true,
		},
		{
			name: "only service admin can allow borrow a security from a new pool",
			store: func(ctx sdk.Context) {
				err := suite.pk.SavePool(ctx, poolstypes.NewPool(1, "utia"))
				suite.Require().NoError(err)
				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					false,
				))
				suite.Require().NoError(err)
			},
			msg:       types.NewMsgBorrowPoolSecurity(1, 1, "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd"),
			shouldErr: true,
		},
		{
			name: "borrow from already present pool fails",
			store: func(ctx sdk.Context) {
				err := suite.pk.SavePool(ctx, poolstypes.NewPool(1, "utia"))
				suite.Require().NoError(err)
				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					false,
				))
				suite.Require().NoError(err)
				err = suite.k.AddPoolToServiceSecuringPools(ctx, 1, 1)
				suite.Require().NoError(err)
			},
			msg:       types.NewMsgBorrowPoolSecurity(1, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			shouldErr: true,
		},
		{
			name: "created service can't borrow",
			store: func(ctx sdk.Context) {
				err := suite.pk.SavePool(ctx, poolstypes.NewPool(1, "utia"))
				suite.Require().NoError(err)
				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_CREATED,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					false,
				))
				suite.Require().NoError(err)
			},
			msg:       types.NewMsgBorrowPoolSecurity(1, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			shouldErr: true,
		},
		{
			name: "inactive service can't borrow",
			store: func(ctx sdk.Context) {
				err := suite.pk.SavePool(ctx, poolstypes.NewPool(1, "utia"))
				suite.Require().NoError(err)
				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_INACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					false,
				))
				suite.Require().NoError(err)
			},
			msg:       types.NewMsgBorrowPoolSecurity(1, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			shouldErr: true,
		},
		{
			name: "security is borrowed properly",
			store: func(ctx sdk.Context) {
				err := suite.pk.SavePool(ctx, poolstypes.NewPool(1, "utia"))
				suite.Require().NoError(err)
				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					false,
				))
				suite.Require().NoError(err)
			},
			msg:       types.NewMsgBorrowPoolSecurity(1, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				secured, err := suite.k.IsServiceSecuredByPool(ctx, 1, 1)
				suite.Require().NoError(err)
				suite.Require().True(secured)
			},
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeBorrowPoolSecurity,
					sdk.NewAttribute(
						servicestypes.AttributeKeyServiceID, "1"),
					sdk.NewAttribute(poolstypes.AttributeKeyPoolID, "1"),
				),
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

			msgServer := keeper.NewMsgServer(suite.k)
			_, err := msgServer.BorrowPoolSecurity(ctx, tc.msg)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				for _, event := range tc.expEvents {
					events := ctx.EventManager().Events()
					suite.Require().Contains(events, event)
				}

				if tc.check != nil {
					tc.check(ctx)
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgServer_CeasePoolSecurityBorrow() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		msg       *types.MsgCeasePoolSecurityBorrow
		shouldErr bool
		expEvents sdk.Events
		check     func(ctx sdk.Context)
	}{
		{
			name: "non existing service returns error",
			store: func(ctx sdk.Context) {
				err := suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
					1, operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)
			},
			msg:       types.NewMsgCeasePoolSecurityBorrow(1, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			shouldErr: true,
		},
		{
			name: "only service admin can cease pool security borrow",
			store: func(ctx sdk.Context) {
				err := suite.pk.SavePool(ctx, poolstypes.NewPool(1, "utia"))
				suite.Require().NoError(err)
				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					false,
				))
				suite.Require().NoError(err)
				err = suite.k.AddPoolToServiceSecuringPools(ctx, 1, 1)
				suite.Require().NoError(err)
			},
			msg:       types.NewMsgCeasePoolSecurityBorrow(1, 1, "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd"),
			shouldErr: true,
		},
		{
			name: "cease from not borrowing pool fails",
			store: func(ctx sdk.Context) {
				err := suite.pk.SavePool(ctx, poolstypes.NewPool(1, "utia"))
				suite.Require().NoError(err)
				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					false,
				))
				suite.Require().NoError(err)
			},
			msg:       types.NewMsgCeasePoolSecurityBorrow(1, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			shouldErr: true,
		},
		{
			name: "security is ceased properly",
			store: func(ctx sdk.Context) {
				err := suite.pk.SavePool(ctx, poolstypes.NewPool(1, "utia"))
				suite.Require().NoError(err)
				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					false,
				))
				suite.Require().NoError(err)
				err = suite.k.AddPoolToServiceSecuringPools(ctx, 1, 1)
				suite.Require().NoError(err)
				err = suite.k.AddPoolToServiceSecuringPools(ctx, 1, 2)
				suite.Require().NoError(err)
			},
			msg:       types.NewMsgCeasePoolSecurityBorrow(1, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				secured, err := suite.k.IsServiceSecuredByPool(ctx, 1, 1)
				suite.Require().NoError(err)
				suite.Require().False(secured)
			},
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeCeasePoolSecurityBorrow,
					sdk.NewAttribute(
						servicestypes.AttributeKeyServiceID, "1"),
					sdk.NewAttribute(poolstypes.AttributeKeyPoolID, "1"),
				),
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

			msgServer := keeper.NewMsgServer(suite.k)
			_, err := msgServer.CeasePoolSecurityBorrow(ctx, tc.msg)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				for _, event := range tc.expEvents {
					events := ctx.EventManager().Events()
					suite.Require().Contains(events, event)
				}

				if tc.check != nil {
					tc.check(ctx)
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgServer_DelegatePool() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		setupCtx  func(ctx sdk.Context) sdk.Context
		msg       *types.MsgDelegatePool
		shouldErr bool
		expEvents sdk.Events
		check     func(ctx sdk.Context)
	}{
		{
			name: "invalid amount returns error",
			store: func(ctx sdk.Context) {
				// Create the pool
				err := suite.pk.SavePool(ctx, poolstypes.Pool{
					ID:              1,
					Denom:           "umilk",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(20),
					DelegatorShares: sdkmath.LegacyNewDec(100),
				})
				suite.Require().NoError(err)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
			},
			msg: &types.MsgDelegatePool{
				Delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Amount:    sdk.NewCoin("umilk", sdkmath.NewInt(0)),
			},
			shouldErr: true,
		},
		{
			name: "not allowed denom returns error",
			store: func(ctx sdk.Context) {
				// Configure the allowed restakable denoms
				err := suite.k.SetRestakableDenoms(ctx, []string{"uinit"})
				suite.Require().NoError(err)

				// Create the pool
				err = suite.pk.SavePool(ctx, poolstypes.Pool{
					ID:              1,
					Denom:           "umilk",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(20),
					DelegatorShares: sdkmath.LegacyNewDec(100),
				})
				suite.Require().NoError(err)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
			},
			msg: &types.MsgDelegatePool{
				Delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Amount:    sdk.NewCoin("umilk", sdkmath.NewInt(10)),
			},
			shouldErr: true,
		},
		{
			name: "valid amount is delegated properly",
			store: func(ctx sdk.Context) {
				// Create the pool
				err := suite.pk.SavePool(ctx, poolstypes.Pool{
					ID:              1,
					Denom:           "umilk",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(20),
					DelegatorShares: sdkmath.LegacyNewDec(100),
				})
				suite.Require().NoError(err)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
			},
			msg: &types.MsgDelegatePool{
				Delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Amount:    sdk.NewCoin("umilk", sdkmath.NewInt(100)),
			},
			shouldErr: false,
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeDelegatePool,
					sdk.NewAttribute(types.AttributeKeyDelegator, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
					sdk.NewAttribute(sdk.AttributeKeyAmount, "100umilk"),
					sdk.NewAttribute(types.AttributeKeyNewShares, "500.000000000000000000pool/1/umilk"),
				),
			},
		},
		{
			name: "allowed denom is delegated properly",
			store: func(ctx sdk.Context) {
				// Configure the allowed restakable denoms
				err := suite.k.SetRestakableDenoms(ctx, []string{"umilk"})
				suite.Require().NoError(err)

				// Create the pool
				err = suite.pk.SavePool(ctx, poolstypes.Pool{
					ID:              1,
					Denom:           "umilk",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(20),
					DelegatorShares: sdkmath.LegacyNewDec(100),
				})
				suite.Require().NoError(err)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
			},
			msg: &types.MsgDelegatePool{
				Delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Amount:    sdk.NewCoin("umilk", sdkmath.NewInt(100)),
			},
			shouldErr: false,
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeDelegatePool,
					sdk.NewAttribute(types.AttributeKeyDelegator, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
					sdk.NewAttribute(sdk.AttributeKeyAmount, "100umilk"),
					sdk.NewAttribute(types.AttributeKeyNewShares, "500.000000000000000000pool/1/umilk"),
				),
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

			msgServer := keeper.NewMsgServer(suite.k)
			_, err := msgServer.DelegatePool(ctx, tc.msg)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
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

func (suite *KeeperTestSuite) TestMsgServer_UndelegatePool() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		setupCtx  func(ctx sdk.Context) sdk.Context
		msg       *types.MsgUndelegatePool
		shouldErr bool
		expEvents sdk.Events
		check     func(ctx sdk.Context)
	}{
		{
			name: "non existing delegation returns error",
			msg: &types.MsgUndelegatePool{
				Delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Amount:    sdk.NewCoin("umilk", sdkmath.NewInt(100)),
			},
			shouldErr: true,
		},
		{
			name: "existing delegation is unbonded properly",
			setupCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
			},
			store: func(ctx sdk.Context) {
				// Set the unbonding time to 1 week
				err := suite.k.SetParams(ctx, types.NewParams(7*24*time.Hour, nil, types.DefaultRestakingCap, types.DefaultMaxEntries))
				suite.Require().NoError(err)

				// Create the pool
				err = suite.pk.SavePool(ctx, poolstypes.Pool{
					ID:      1,
					Denom:   "umilk",
					Address: poolstypes.GetPoolAddress(1).String(),
				})
				suite.Require().NoError(err)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)

				// Delegate some funds
				msgServer := keeper.NewMsgServer(suite.k)
				_, err = msgServer.DelegatePool(ctx, &types.MsgDelegatePool{
					Delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					Amount:    sdk.NewCoin("umilk", sdkmath.NewInt(100)),
				})
				suite.Require().NoError(err)

				// Check the delegation
				delegation, found, err := suite.k.GetPoolDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().Equal(sdk.NewDecCoins(sdk.NewDecCoin("pool/1/umilk", sdkmath.NewInt(100))), delegation.Shares)
			},
			msg: &types.MsgUndelegatePool{
				Delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Amount:    sdk.NewCoin("umilk", sdkmath.NewInt(100)),
			},
			shouldErr: false,
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeUnbondPool,
					sdk.NewAttribute(sdk.AttributeKeyAmount, "100umilk"),
					sdk.NewAttribute(types.AttributeKeyDelegator, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
					sdk.NewAttribute(types.AttributeKeyCompletionTime, "2024-01-08T00:00:00Z"),
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the gas charged is at least BaseDelegationDenomCost
				// The 36950 is obtained by running this test with BaseDelegationDenomCost set to 0
				suite.Require().GreaterOrEqual(ctx.GasMeter().GasConsumed(), 36950+types.BaseDelegationDenomCost)

				// Check the delegation
				delegation, found, err := suite.k.GetPoolDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().NoError(err)
				suite.Require().False(found)
				suite.Require().Empty(delegation.Shares)
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

			// Reset the gas meter
			ctx = ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())

			msgServer := keeper.NewMsgServer(suite.k)
			_, err := msgServer.UndelegatePool(ctx, tc.msg)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
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

func (suite *KeeperTestSuite) TestMsgServer_DelegateOperator() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		setupCtx  func(ctx sdk.Context) sdk.Context
		msg       *types.MsgDelegateOperator
		shouldErr bool
		expEvents sdk.Events
		check     func(ctx sdk.Context)
	}{
		{
			name: "invalid amount returns error",
			store: func(ctx sdk.Context) {
				// Create the operator
				err := suite.ok.SaveOperator(ctx, operatorstypes.Operator{
					ID:      1,
					Status:  operatorstypes.OPERATOR_STATUS_ACTIVE,
					Address: operatorstypes.GetOperatorAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(20)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				})
				suite.Require().NoError(err)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
			},
			msg: &types.MsgDelegateOperator{
				OperatorID: 1,
				Delegator:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Amount:     sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(0))),
			},
			shouldErr: true,
		},
		{
			name: "not allowed denom returns error",
			store: func(ctx sdk.Context) {
				// Configure the allowed restakable denoms
				err := suite.k.SetRestakableDenoms(ctx, []string{"uinit"})
				suite.Require().NoError(err)

				// Create the operator
				err = suite.ok.SaveOperator(ctx, operatorstypes.Operator{
					ID:      1,
					Status:  operatorstypes.OPERATOR_STATUS_ACTIVE,
					Address: operatorstypes.GetOperatorAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(20)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				})
				suite.Require().NoError(err)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
			},
			msg: &types.MsgDelegateOperator{
				OperatorID: 1,
				Delegator:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Amount:     sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			},
			shouldErr: true,
		},
		{
			name: "valid amount is delegated properly",
			store: func(ctx sdk.Context) {
				// Create the operator
				err := suite.ok.SaveOperator(ctx, operatorstypes.Operator{
					ID:      1,
					Status:  operatorstypes.OPERATOR_STATUS_ACTIVE,
					Address: operatorstypes.GetOperatorAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(20)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				})
				suite.Require().NoError(err)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
			},
			msg: &types.MsgDelegateOperator{
				OperatorID: 1,
				Delegator:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Amount:     sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			},
			shouldErr: false,
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeDelegateOperator,
					sdk.NewAttribute(types.AttributeKeyDelegator, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
					sdk.NewAttribute(operatorstypes.AttributeKeyOperatorID, "1"),
					sdk.NewAttribute(sdk.AttributeKeyAmount, "100umilk"),
					sdk.NewAttribute(types.AttributeKeyNewShares, "500.000000000000000000operator/1/umilk"),
				),
			},
		},
		{
			name: "allowed denom is delegated properly",
			store: func(ctx sdk.Context) {
				// Configure the allowed restakable denoms
				err := suite.k.SetRestakableDenoms(ctx, []string{"umilk"})
				suite.Require().NoError(err)

				// Create the operator
				err = suite.ok.SaveOperator(ctx, operatorstypes.Operator{
					ID:      1,
					Status:  operatorstypes.OPERATOR_STATUS_ACTIVE,
					Address: operatorstypes.GetOperatorAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(20)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				})
				suite.Require().NoError(err)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
			},
			msg: &types.MsgDelegateOperator{
				OperatorID: 1,
				Delegator:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Amount:     sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			},
			shouldErr: false,
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeDelegateOperator,
					sdk.NewAttribute(types.AttributeKeyDelegator, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
					sdk.NewAttribute(operatorstypes.AttributeKeyOperatorID, "1"),
					sdk.NewAttribute(sdk.AttributeKeyAmount, "100umilk"),
					sdk.NewAttribute(types.AttributeKeyNewShares, "500.000000000000000000operator/1/umilk"),
				),
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

			msgServer := keeper.NewMsgServer(suite.k)
			_, err := msgServer.DelegateOperator(ctx, tc.msg)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
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

func (suite *KeeperTestSuite) TestMsgServer_UndelegateOperator() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		setupCtx  func(ctx sdk.Context) sdk.Context
		msg       *types.MsgUndelegateOperator
		shouldErr bool
		expEvents sdk.Events
		check     func(ctx sdk.Context)
	}{
		{
			name: "non existing delegation returns error",
			msg: &types.MsgUndelegateOperator{
				Delegator:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				OperatorID: 1,
				Amount:     sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			},
			shouldErr: true,
		},
		{
			name: "existing delegation is unbonded properly",
			setupCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
			},
			store: func(ctx sdk.Context) {
				// Set the unbonding time to 1 week
				err := suite.k.SetParams(ctx, types.NewParams(7*24*time.Hour, nil, types.DefaultRestakingCap, types.DefaultMaxEntries))
				suite.Require().NoError(err)

				// Create the operator
				err = suite.ok.SaveOperator(ctx, operatorstypes.Operator{
					ID:      1,
					Status:  operatorstypes.OPERATOR_STATUS_ACTIVE,
					Address: operatorstypes.GetOperatorAddress(1).String(),
				})
				suite.Require().NoError(err)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)

				// Delegate some funds
				msgServer := keeper.NewMsgServer(suite.k)
				_, err = msgServer.DelegateOperator(ctx, &types.MsgDelegateOperator{
					OperatorID: 1,
					Delegator:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					Amount:     sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				})
				suite.Require().NoError(err)

				// Check the delegation
				delegation, found, err := suite.k.GetOperatorDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().Equal(sdk.NewDecCoins(sdk.NewDecCoin("operator/1/umilk", sdkmath.NewInt(100))), delegation.Shares)
			},
			msg: &types.MsgUndelegateOperator{
				Delegator:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				OperatorID: 1,
				Amount:     sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			},
			shouldErr: false,
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeUnbondOperator,
					sdk.NewAttribute(sdk.AttributeKeyAmount, "100umilk"),
					sdk.NewAttribute(types.AttributeKeyDelegator, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
					sdk.NewAttribute(operatorstypes.AttributeKeyOperatorID, "1"),
					sdk.NewAttribute(types.AttributeKeyCompletionTime, "2024-01-08T00:00:00Z"),
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the gas charged is at least BaseDelegationDenomCost
				// The 36690 is obtained by running this test with BaseDelegationDenomCost set to 0
				suite.Require().GreaterOrEqual(ctx.GasMeter().GasConsumed(), 36690+types.BaseDelegationDenomCost)

				// Check the delegation
				delegation, found, err := suite.k.GetOperatorDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().NoError(err)
				suite.Require().False(found)
				suite.Require().Empty(delegation.Shares)
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

			// Reset the gas meter
			ctx = ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())

			msgServer := keeper.NewMsgServer(suite.k)
			_, err := msgServer.UndelegateOperator(ctx, tc.msg)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
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

func (suite *KeeperTestSuite) TestMsgServer_DelegateService() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		setupCtx  func(ctx sdk.Context) sdk.Context
		msg       *types.MsgDelegateService
		shouldErr bool
		expEvents sdk.Events
		check     func(ctx sdk.Context)
	}{
		{
			name: "invalid amount returns error",
			store: func(ctx sdk.Context) {
				// Create the service
				err := suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      1,
					Status:  servicestypes.SERVICE_STATUS_ACTIVE,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(20)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				})
				suite.Require().NoError(err)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
			},
			msg: &types.MsgDelegateService{
				ServiceID: 1,
				Delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Amount:    sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(0))),
			},
			shouldErr: true,
		},
		{
			name: "not allowed denom returns error",
			store: func(ctx sdk.Context) {
				// Configure the allowed restakable denoms
				err := suite.k.SetRestakableDenoms(ctx, []string{"uinit"})
				suite.Require().NoError(err)

				// Create the service
				err = suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      1,
					Status:  servicestypes.SERVICE_STATUS_ACTIVE,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(20)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				})
				suite.Require().NoError(err)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
			},
			msg: &types.MsgDelegateService{
				ServiceID: 1,
				Delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Amount:    sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			},
			shouldErr: true,
		},
		{
			name: "not allowed service denom returns error",
			store: func(ctx sdk.Context) {
				// Create the service
				err := suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      1,
					Status:  servicestypes.SERVICE_STATUS_ACTIVE,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(20)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				})
				suite.Require().NoError(err)

				// Configure the service parameters
				err = suite.sk.SetServiceParams(ctx, 1, servicestypes.NewServiceParams([]string{"uinit"}))
				suite.Require().NoError(err)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
			},
			msg: &types.MsgDelegateService{
				ServiceID: 1,
				Delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Amount:    sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			},
			shouldErr: true,
		},
		{
			name: "not allowed service denom when intersecting allowed denoms returns error",
			store: func(ctx sdk.Context) {
				// Configure the allowed restakable denoms
				err := suite.k.SetRestakableDenoms(ctx, []string{"uinit", "umilk"})
				suite.Require().NoError(err)

				// Create the service
				err = suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      1,
					Status:  servicestypes.SERVICE_STATUS_ACTIVE,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(20)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				})
				suite.Require().NoError(err)

				// Configure the service parameters
				err = suite.sk.SetServiceParams(ctx, 1, servicestypes.NewServiceParams([]string{"uinit"}))
				suite.Require().NoError(err)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
			},
			msg: &types.MsgDelegateService{
				ServiceID: 1,
				Delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Amount:    sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			},
			shouldErr: true,
		},
		{
			name: "valid amount is delegated properly",
			store: func(ctx sdk.Context) {
				// Create the service
				err := suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      1,
					Status:  servicestypes.SERVICE_STATUS_ACTIVE,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(20)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				})
				suite.Require().NoError(err)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
			},
			msg: &types.MsgDelegateService{
				ServiceID: 1,
				Delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Amount:    sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			},
			shouldErr: false,
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeDelegateService,
					sdk.NewAttribute(types.AttributeKeyDelegator, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
					sdk.NewAttribute(servicestypes.AttributeKeyServiceID, "1"),
					sdk.NewAttribute(sdk.AttributeKeyAmount, "100umilk"),
					sdk.NewAttribute(types.AttributeKeyNewShares, "500.000000000000000000service/1/umilk"),
				),
			},
		},
		{
			name: "allowed denom is delegated properly",
			store: func(ctx sdk.Context) {
				// Configure the allowed restakable denoms
				err := suite.k.SetRestakableDenoms(ctx, []string{"umilk"})
				suite.Require().NoError(err)

				// Create the service
				err = suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      1,
					Status:  servicestypes.SERVICE_STATUS_ACTIVE,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(20)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				})
				suite.Require().NoError(err)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
			},
			msg: &types.MsgDelegateService{
				ServiceID: 1,
				Delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Amount:    sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			},
			shouldErr: false,
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeDelegateService,
					sdk.NewAttribute(types.AttributeKeyDelegator, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
					sdk.NewAttribute(servicestypes.AttributeKeyServiceID, "1"),
					sdk.NewAttribute(sdk.AttributeKeyAmount, "100umilk"),
					sdk.NewAttribute(types.AttributeKeyNewShares, "500.000000000000000000service/1/umilk"),
				),
			},
		},
		{
			name: "allowed service denom is delegated properly",
			store: func(ctx sdk.Context) {
				// Create the service
				err := suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      1,
					Status:  servicestypes.SERVICE_STATUS_ACTIVE,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(20)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				})
				suite.Require().NoError(err)

				// Configure the service's allowed restakable denoms
				err = suite.sk.SetServiceParams(ctx, 1, servicestypes.NewServiceParams([]string{"umilk"}))
				suite.Require().NoError(err)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
			},
			msg: &types.MsgDelegateService{
				ServiceID: 1,
				Delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Amount:    sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			},
			shouldErr: false,
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeDelegateService,
					sdk.NewAttribute(types.AttributeKeyDelegator, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
					sdk.NewAttribute(servicestypes.AttributeKeyServiceID, "1"),
					sdk.NewAttribute(sdk.AttributeKeyAmount, "100umilk"),
					sdk.NewAttribute(types.AttributeKeyNewShares, "500.000000000000000000service/1/umilk"),
				),
			},
		},
		{
			name: "allowed denom after intersecting allowed denoms is delegated properly",
			store: func(ctx sdk.Context) {
				// Configure the allowed restakable denoms
				err := suite.k.SetRestakableDenoms(ctx, []string{"umilk", "uinit"})
				suite.Require().NoError(err)

				// Create the service
				err = suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      1,
					Status:  servicestypes.SERVICE_STATUS_ACTIVE,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(20)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				})
				suite.Require().NoError(err)

				// Configure the service's allowed restakable denoms
				err = suite.sk.SetServiceParams(ctx, 1, servicestypes.NewServiceParams([]string{"umilk"}))
				suite.Require().NoError(err)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
			},
			msg: &types.MsgDelegateService{
				ServiceID: 1,
				Delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				Amount:    sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			},
			shouldErr: false,
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeDelegateService,
					sdk.NewAttribute(types.AttributeKeyDelegator, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
					sdk.NewAttribute(servicestypes.AttributeKeyServiceID, "1"),
					sdk.NewAttribute(sdk.AttributeKeyAmount, "100umilk"),
					sdk.NewAttribute(types.AttributeKeyNewShares, "500.000000000000000000service/1/umilk"),
				),
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

			msgServer := keeper.NewMsgServer(suite.k)
			_, err := msgServer.DelegateService(ctx, tc.msg)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
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

func (suite *KeeperTestSuite) TestMsgServer_UndelegateService() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		setupCtx  func(ctx sdk.Context) sdk.Context
		msg       *types.MsgUndelegateService
		shouldErr bool
		expEvents sdk.Events
		check     func(ctx sdk.Context)
	}{
		{
			name: "non existing delegation returns error",
			msg: &types.MsgUndelegateService{
				Delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				ServiceID: 1,
				Amount:    sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			},
			shouldErr: true,
		},
		{
			name: "existing delegation is unbonded properly",
			setupCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
			},
			store: func(ctx sdk.Context) {
				// Set the unbonding time to 1 week
				err := suite.k.SetParams(ctx, types.NewParams(7*24*time.Hour, nil, types.DefaultRestakingCap, types.DefaultMaxEntries))
				suite.Require().NoError(err)

				// Create the service
				err = suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      1,
					Status:  servicestypes.SERVICE_STATUS_ACTIVE,
					Address: servicestypes.GetServiceAddress(1).String(),
				})
				suite.Require().NoError(err)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)

				// Delegate some funds
				msgServer := keeper.NewMsgServer(suite.k)
				_, err = msgServer.DelegateService(ctx, &types.MsgDelegateService{
					ServiceID: 1,
					Delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					Amount:    sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				})
				suite.Require().NoError(err)

				// Check the delegation
				delegation, found, err := suite.k.GetServiceDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().Equal(sdk.NewDecCoins(sdk.NewDecCoin("service/1/umilk", sdkmath.NewInt(100))), delegation.Shares)
			},
			msg: &types.MsgUndelegateService{
				Delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				ServiceID: 1,
				Amount:    sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			},
			shouldErr: false,
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeUnbondService,
					sdk.NewAttribute(sdk.AttributeKeyAmount, "100umilk"),
					sdk.NewAttribute(types.AttributeKeyDelegator, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
					sdk.NewAttribute(servicestypes.AttributeKeyServiceID, "1"),
					sdk.NewAttribute(types.AttributeKeyCompletionTime, "2024-01-08T00:00:00Z"),
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the gas charged is at least BaseDelegationDenomCost
				// The 36680 is obtained by running this test with BaseDelegationDenomCost set to 0
				suite.Require().GreaterOrEqual(ctx.GasMeter().GasConsumed(), 36680+types.BaseDelegationDenomCost)

				// Check the delegation
				delegation, found, err := suite.k.GetServiceDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().NoError(err)
				suite.Require().False(found)
				suite.Require().Empty(delegation.Shares)
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

			// Reset the gas meter
			ctx = ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())

			msgServer := keeper.NewMsgServer(suite.k)
			_, err := msgServer.UndelegateService(ctx, tc.msg)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
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

func (suite *KeeperTestSuite) TestMsgServer_SetUserPreferences() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		setupCtx  func(ctx sdk.Context) sdk.Context
		msg       *types.MsgSetUserPreferences
		shouldErr bool
		expEvents sdk.Events
		check     func(ctx sdk.Context)
	}{
		{
			name: "not found services return error",
			msg: types.NewMsgSetUserPreferences(
				types.NewUserPreferences([]types.TrustedServiceEntry{
					types.NewTrustedServiceEntry(1, nil),
					types.NewTrustedServiceEntry(2, nil),
				}),
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "invalid user preferences return error",
			msg: types.NewMsgSetUserPreferences(
				types.NewUserPreferences([]types.TrustedServiceEntry{
					types.NewTrustedServiceEntry(0, nil),
				}),
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "new user preferences are set properly",
			store: func(ctx sdk.Context) {
				// Store the services
				err := suite.sk.SaveService(ctx, servicestypes.NewService(
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

				err = suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      2,
					Address: servicestypes.GetServiceAddress(2).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("utia", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/2/utia", sdkmath.LegacyNewDec(150)),
					),
				})
				suite.Require().NoError(err)
			},
			msg: types.NewMsgSetUserPreferences(
				types.NewUserPreferences([]types.TrustedServiceEntry{
					types.NewTrustedServiceEntry(1, nil),
					types.NewTrustedServiceEntry(2, nil),
				}),
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: false,
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeSetUserPreferences,
					sdk.NewAttribute(types.AttributeKeyUser, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the preferences are stored properly
				stored, err := suite.k.GetUserPreferences(ctx, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewUserPreferences([]types.TrustedServiceEntry{
					types.NewTrustedServiceEntry(1, nil),
					types.NewTrustedServiceEntry(2, nil),
				}), stored)
			},
		},
		{
			name: "existing user preferences are updated properly",
			store: func(ctx sdk.Context) {
				// Store the services
				err := suite.sk.SaveService(ctx, servicestypes.NewService(
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

				err = suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      2,
					Address: servicestypes.GetServiceAddress(2).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("utia", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/2/utia", sdkmath.LegacyNewDec(150)),
					),
				})
				suite.Require().NoError(err)

				// Set the user preferences
				preferences := types.NewUserPreferences([]types.TrustedServiceEntry{
					types.NewTrustedServiceEntry(1, nil),
					types.NewTrustedServiceEntry(2, nil),
				})
				err = suite.k.SetUserPreferences(ctx, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4", preferences)
				suite.Require().NoError(err)
			},
			msg: types.NewMsgSetUserPreferences(
				types.NewUserPreferences(nil),
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: false,
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeSetUserPreferences,
					sdk.NewAttribute(types.AttributeKeyUser, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the preferences are stored properly
				stored, err := suite.k.GetUserPreferences(ctx, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewUserPreferences(nil), stored)
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
			_, err := msgServer.SetUserPreferences(ctx, tc.msg)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
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

func (suite *KeeperTestSuite) TestMsgServer_UpdateParams() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		setupCtx  func(ctx sdk.Context) sdk.Context
		msg       *types.MsgUpdateParams
		shouldErr bool
		expEvents sdk.Events
		check     func(ctx sdk.Context)
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
			shouldErr: false,
			expEvents: sdk.Events{},
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
			_, err := msgServer.UpdateParams(ctx, tc.msg)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
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
