package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/milkyway-labs/milkyway/x/services/keeper"
	"github.com/milkyway-labs/milkyway/x/services/types"
)

func (suite *KeeperTestSuite) TestMsgServer_CreateService() {
	testCases := []struct {
		name        string
		setup       func()
		store       func(ctx sdk.Context)
		setupCtx    func(ctx sdk.Context) sdk.Context
		msg         *types.MsgCreateService
		shouldErr   bool
		expResponse *types.MsgCreateServiceResponse
		expEvents   sdk.Events
		check       func(ctx sdk.Context)
	}{
		{
			name: "non existing next service id returns error",
			store: func(ctx sdk.Context) {
				suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000))),
				))
			},
			msg: types.NewMsgCreateService(
				"MilkyWay",
				"MilkyWay is a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "invalid service returns error",
			store: func(ctx sdk.Context) {
				suite.k.SetNextServiceID(ctx, 1)
				suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000))),
				))
			},
			msg: types.NewMsgCreateService(
				"MilkyWay",
				types.DoNotModify,
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "user without enough funds return error",
			store: func(ctx sdk.Context) {
				suite.k.SetNextServiceID(ctx, 1)
				suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000))),
				))
			},
			msg: types.NewMsgCreateService(
				"MilkyWay",
				"MilkyWay is a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "valid service is created properly",
			store: func(ctx sdk.Context) {
				suite.k.SetNextServiceID(ctx, 1)
				suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000))),
				))
				suite.fundAccount(ctx, "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(200_000))))
			},
			msg: types.NewMsgCreateService(
				"MilkyWay",
				"MilkyWay is a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: false,
			expResponse: &types.MsgCreateServiceResponse{
				NewServiceID: 1,
			},
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeCreateService,
					sdk.NewAttribute(types.AttributeKeyServiceID, "1"),
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the service has been stored
				stored, found := suite.k.GetService(ctx, 1)
				suite.Require().True(found)
				suite.Require().Equal(types.NewService(
					1,
					types.SERVICE_STATUS_CREATED,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				), stored)

				// Make sure the service account has been created properly
				hasAccount := suite.ak.HasAccount(ctx, types.GetServiceAddress(1))
				suite.Require().True(hasAccount)

				// Make sure the next service id has been incremented
				nextServiceID, err := suite.k.GetNextServiceID(ctx)
				suite.Require().NoError(err)
				suite.Require().Equal(uint32(2), nextServiceID)

				// Make sure the user was charged for the fee
				userAddress, err := sdk.AccAddressFromBech32("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd")
				suite.Require().NoError(err)
				balance := suite.bk.GetBalance(ctx, userAddress, "uatom")
				suite.Require().Equal(sdk.NewCoin("uatom", sdkmath.NewInt(100_000)), balance)

				// Make sure the fee was transferred to the module account
				poolBalance := suite.bk.GetBalance(ctx, authtypes.NewModuleAddress(authtypes.FeeCollectorName), "uatom")
				suite.Require().Equal(sdk.NewCoin("uatom", sdkmath.NewInt(100_000)), poolBalance)
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
			res, err := msgServer.CreateService(ctx, tc.msg)
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

func (suite *KeeperTestSuite) TestMsgServer_UpdateService() {
	testCases := []struct {
		name        string
		setup       func()
		store       func(ctx sdk.Context)
		setupCtx    func(ctx sdk.Context) sdk.Context
		msg         *types.MsgUpdateService
		shouldErr   bool
		expResponse *types.MsgUpdateServiceResponse
		expEvents   sdk.Events
		check       func(ctx sdk.Context)
	}{
		{
			name: "not found service returns error",
			msg: types.NewMsgUpdateService(
				1,
				"MilkyWay",
				types.DoNotModify,
				types.DoNotModify,
				types.DoNotModify,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "non admin user returns error",
			store: func(ctx sdk.Context) {
				suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_CREATED,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				))
			},
			msg: types.NewMsgUpdateService(
				1,
				"MilkyWay Modular Restaking",
				types.DoNotModify,
				types.DoNotModify,
				types.DoNotModify,
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
			),
			shouldErr: true,
		},
		{
			name: "invalid service returns error",
			store: func(ctx sdk.Context) {
				suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_CREATED,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				))
			},
			msg: types.NewMsgUpdateService(
				1,
				"",
				types.DoNotModify,
				types.DoNotModify,
				types.DoNotModify,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "valid service is updated properly",
			store: func(ctx sdk.Context) {
				suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_CREATED,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				))
			},
			msg: types.NewMsgUpdateService(
				1,
				"MilkyWay Modular Restaking",
				types.DoNotModify,
				types.DoNotModify,
				types.DoNotModify,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr:   false,
			expResponse: &types.MsgUpdateServiceResponse{},
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeUpdateService,
					sdk.NewAttribute(types.AttributeKeyServiceID, "1"),
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the service was updated
				stored, found := suite.k.GetService(ctx, 1)
				suite.Require().True(found)
				suite.Require().Equal(types.NewService(
					1,
					types.SERVICE_STATUS_CREATED,
					"MilkyWay Modular Restaking",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
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
			res, err := msgServer.UpdateService(ctx, tc.msg)
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

func (suite *KeeperTestSuite) TestMsgServer_ActivateService() {
	testCases := []struct {
		name        string
		setup       func()
		store       func(ctx sdk.Context)
		setupCtx    func(ctx sdk.Context) sdk.Context
		msg         *types.MsgActivateService
		shouldErr   bool
		expResponse *types.MsgActivateServiceResponse
		expEvents   sdk.Events
		check       func(ctx sdk.Context)
	}{
		{
			name: "service not found returns error",
			msg: types.NewMsgActivateService(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "non admin user returns error",
			store: func(ctx sdk.Context) {
				suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_CREATED,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				))
			},
			msg: types.NewMsgActivateService(
				1,
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
			),
			shouldErr: true,
		},
		{
			name: "already active service returns error",
			store: func(ctx sdk.Context) {
				suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				))
			},
			msg: types.NewMsgActivateService(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "service with status CREATED is activated properly",
			store: func(ctx sdk.Context) {
				suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_CREATED,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				))
			},
			msg: types.NewMsgActivateService(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr:   false,
			expResponse: &types.MsgActivateServiceResponse{},
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeActivateService,
					sdk.NewAttribute(types.AttributeKeyServiceID, "1"),
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
			res, err := msgServer.ActivateService(ctx, tc.msg)
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

func (suite *KeeperTestSuite) TestMsgServer_DeactivateService() {
	testCases := []struct {
		name        string
		setup       func()
		store       func(ctx sdk.Context)
		setupCtx    func(ctx sdk.Context) sdk.Context
		msg         *types.MsgDeactivateService
		shouldErr   bool
		expResponse *types.MsgDeactivateServiceResponse
		expEvents   sdk.Events
		check       func(ctx sdk.Context)
	}{
		{
			name: "service not found returns error",
			msg: types.NewMsgDeactivateService(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "non admin user returns error",
			store: func(ctx sdk.Context) {
				suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				))
			},
			msg: types.NewMsgDeactivateService(
				1,
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
			),
			shouldErr: true,
		},
		{
			name: "valid service is deactivated properly",
			store: func(ctx sdk.Context) {
				suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				))
			},
			msg: types.NewMsgDeactivateService(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr:   false,
			expResponse: &types.MsgDeactivateServiceResponse{},
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeDeactivateService,
					sdk.NewAttribute(types.AttributeKeyServiceID, "1"),
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the service was deactivated
				stored, found := suite.k.GetService(ctx, 1)
				suite.Require().True(found)
				suite.Require().Equal(types.NewService(
					1,
					types.SERVICE_STATUS_INACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
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
			res, err := msgServer.DeactivateService(ctx, tc.msg)
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

func (suite *KeeperTestSuite) TestMsgServer_TransferServiceOwnership() {
	testCases := []struct {
		name        string
		setup       func()
		store       func(ctx sdk.Context)
		setupCtx    func(ctx sdk.Context) sdk.Context
		msg         *types.MsgTransferServiceOwnership
		shouldErr   bool
		expResponse *types.MsgTransferServiceOwnershipResponse
		expEvents   sdk.Events
		check       func(ctx sdk.Context)
	}{
		{
			name: "not found service returns error",
			msg: types.NewMsgTransferServiceOwnership(
				1,
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "non admin sender returns error",
			store: func(ctx sdk.Context) {
				suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				))
			},
			msg: types.NewMsgTransferServiceOwnership(
				1,
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "service ownership transferred successfully",
			store: func(ctx sdk.Context) {
				suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				))
			},
			msg: types.NewMsgTransferServiceOwnership(
				1,
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr:   false,
			expResponse: &types.MsgTransferServiceOwnershipResponse{},
			expEvents: []sdk.Event{
				sdk.NewEvent(
					types.EventTypeTransferServiceOwnership,
					sdk.NewAttribute(types.AttributeKeyServiceID, "1"),
					sdk.NewAttribute(types.AttributeKeyNewAdmin, "cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn"),
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the service was updated
				stored, found := suite.k.GetService(ctx, 1)
				suite.Require().True(found)
				suite.Require().Equal(types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
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
			res, err := msgServer.TransferServiceOwnership(sdk.WrapSDKContext(ctx), tc.msg)
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

func (suite *KeeperTestSuite) TestMsgServer_SetServiceParams() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		msg       *types.MsgSetServiceParams
		shouldErr bool
		expEvents sdk.Events
		check     func(ctx sdk.Context)
	}{
		{
			name: "not found service returns error",
			msg: types.NewMsgSetServiceParams(
				1,
				types.NewServiceParams(sdkmath.LegacyMustNewDecFromStr("0.2")),
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "non admin user returns error",
			store: func(ctx sdk.Context) {
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_CREATED,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgSetServiceParams(
				1,
				types.NewServiceParams(sdkmath.LegacyMustNewDecFromStr("0.2")),
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
			),
			shouldErr: true,
		},
		{
			name: "invalid service params returns error",
			store: func(ctx sdk.Context) {
				suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_CREATED,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				))
			},
			msg: types.NewMsgSetServiceParams(
				1,
				types.NewServiceParams(sdkmath.LegacyMustNewDecFromStr("2")),
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "valid params are set properly",
			store: func(ctx sdk.Context) {
				suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_CREATED,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				))
			},
			msg: types.NewMsgSetServiceParams(
				1,
				types.NewServiceParams(sdkmath.LegacyMustNewDecFromStr("0.7")),
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: false,
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeSetServiceParams,
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the service was updated
				stored, err := suite.k.GetServiceParams(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().Equal(
					types.NewServiceParams(sdkmath.LegacyMustNewDecFromStr("0.7")), stored)
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
			res, err := msgServer.SetServiceParams(ctx, tc.msg)
			if tc.shouldErr {
				suite.Require().Error(err)
				suite.Require().Nil(res)
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
