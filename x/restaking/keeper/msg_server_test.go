package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	"github.com/milkyway-labs/milkyway/x/restaking/keeper"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
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
					1, operatorstypes.OPERATOR_STATUS_ACTIVE,
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
					1, operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)

				err = suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_INACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
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
				))
				suite.Require().NoError(err)

				err = suite.k.AddServiceToOperator(ctx, 1, 1)
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
				))
				suite.Require().NoError(err)

				err = suite.k.AddServiceToOperator(ctx, 1, 1)
				suite.Require().NoError(err)
			},
			msg: types.NewMsgLeaveService(
				1,
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			check: func(ctx sdk.Context) {
				joinedServices, err := suite.k.GetOperatorJoinedServices(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().Empty(joinedServices.ServiceIDs)
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

func (suite *KeeperTestSuite) TestMsgServer_AllowOperator() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		msg       *types.MsgAllowOperator
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
			msg: types.NewMsgAllowOperator(
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
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgAllowOperator(
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
				))
				suite.Require().NoError(err)
			},
			msg:       types.NewMsgAllowOperator(1, 1, "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd"),
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
				))
				suite.Require().NoError(err)

				err = suite.k.AddOperatorToServiceAllowList(ctx, 1, 1)
				suite.Require().NoError(err)
			},
			msg: types.NewMsgAllowOperator(
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
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgAllowOperator(
				1,
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				configured, err := suite.k.IsServiceOpertorsAllowListConfigured(ctx, 1)
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
			_, err := msgServer.AllowOperator(ctx, tc.msg)
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
		msg       *types.MsgRemoveAllowedOperator
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
			msg: types.NewMsgRemoveAllowedOperator(
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
				))
				suite.Require().NoError(err)

				err = suite.k.AddOperatorToServiceAllowList(ctx, 1, 1)
				suite.Require().NoError(err)
			},
			msg:       types.NewMsgRemoveAllowedOperator(1, 1, "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd"),
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
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgRemoveAllowedOperator(
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
				))
				suite.Require().NoError(err)

				err = suite.k.AddOperatorToServiceAllowList(ctx, 1, 1)
				suite.Require().NoError(err)
				err = suite.k.AddOperatorToServiceAllowList(ctx, 1, 2)
				suite.Require().NoError(err)
			},
			msg:       types.NewMsgRemoveAllowedOperator(1, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
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
			_, err := msgServer.RemoveAllowedOperator(ctx, tc.msg)
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
				suite.pk.SavePool(ctx, poolstypes.NewPool(1, "utia"))
			},
			msg:       types.NewMsgBorrowPoolSecurity(1, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			shouldErr: true,
		},
		{
			name: "non existing pool returns error",
			store: func(ctx sdk.Context) {
				suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
			},
			msg:       types.NewMsgBorrowPoolSecurity(1, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			shouldErr: true,
		},
		{
			name: "only service admin can allow borrow a security from a new pool",
			store: func(ctx sdk.Context) {
				suite.pk.SavePool(ctx, poolstypes.NewPool(1, "utia"))
				suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
			},
			msg:       types.NewMsgBorrowPoolSecurity(1, 1, "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd"),
			shouldErr: true,
		},
		{
			name: "security is borrowed properly",
			store: func(ctx sdk.Context) {
				suite.pk.SavePool(ctx, poolstypes.NewPool(1, "utia"))
				suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
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
					sdk.NewAttribute(types.AttributeKeyPoolID, "1"),
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
				suite.ok.SaveOperator(ctx, operatorstypes.NewOperator(
					1, operatorstypes.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
			},
			msg:       types.NewMsgCeasePoolSecurityBorrow(1, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			shouldErr: true,
		},
		{
			name: "only service admin can cease pool security borrow",
			store: func(ctx sdk.Context) {
				suite.pk.SavePool(ctx, poolstypes.NewPool(1, "utia"))
				suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.k.AddPoolToServiceSecuringPools(ctx, 1, 1)
			},
			msg:       types.NewMsgCeasePoolSecurityBorrow(1, 1, "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd"),
			shouldErr: true,
		},
		{
			name: "security is ceased properly",
			store: func(ctx sdk.Context) {
				suite.pk.SavePool(ctx, poolstypes.NewPool(1, "utia"))
				suite.sk.SaveService(ctx, servicestypes.NewService(
					1, servicestypes.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.k.AddPoolToServiceSecuringPools(ctx, 1, 1)
				suite.k.AddPoolToServiceSecuringPools(ctx, 1, 2)
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
					sdk.NewAttribute(types.AttributeKeyPoolID, "1"),
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
				suite.k.SetParams(ctx, types.NewParams(7*24*time.Hour))

				// Create the pool
				err := suite.pk.SavePool(ctx, poolstypes.Pool{
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
				delegation, found := suite.k.GetPoolDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
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
				// Check the delegation
				delegation, found := suite.k.GetPoolDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
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
				suite.ok.SaveOperator(ctx, operatorstypes.Operator{
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
			name: "valid amount is delegated properly",
			store: func(ctx sdk.Context) {
				// Create the operator
				suite.ok.SaveOperator(ctx, operatorstypes.Operator{
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
				suite.k.SetParams(ctx, types.NewParams(7*24*time.Hour))

				// Create the operator
				suite.ok.SaveOperator(ctx, operatorstypes.Operator{
					ID:      1,
					Status:  operatorstypes.OPERATOR_STATUS_ACTIVE,
					Address: operatorstypes.GetOperatorAddress(1).String(),
				})

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)

				// Delegate some funds
				msgServer := keeper.NewMsgServer(suite.k)
				_, err := msgServer.DelegateOperator(ctx, &types.MsgDelegateOperator{
					OperatorID: 1,
					Delegator:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					Amount:     sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				})
				suite.Require().NoError(err)

				// Check the delegation
				delegation, found := suite.k.GetOperatorDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
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
				// Check the delegation
				delegation, found := suite.k.GetOperatorDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
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
				suite.sk.SaveService(ctx, servicestypes.Service{
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
			name: "valid amount is delegated properly",
			store: func(ctx sdk.Context) {
				// Create the service
				suite.sk.SaveService(ctx, servicestypes.Service{
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
				suite.k.SetParams(ctx, types.NewParams(7*24*time.Hour))

				// Create the service
				suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      1,
					Status:  servicestypes.SERVICE_STATUS_ACTIVE,
					Address: servicestypes.GetServiceAddress(1).String(),
				})

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)

				// Delegate some funds
				msgServer := keeper.NewMsgServer(suite.k)
				_, err := msgServer.DelegateService(ctx, &types.MsgDelegateService{
					ServiceID: 1,
					Delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					Amount:    sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				})
				suite.Require().NoError(err)

				// Check the delegation
				delegation, found := suite.k.GetServiceDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
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
				// Check the delegation
				delegation, found := suite.k.GetServiceDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
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
