package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/milkyway-labs/milkyway/v2/x/services/keeper"
	"github.com/milkyway-labs/milkyway/v2/x/services/types"
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
				err := suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000))),
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgCreateService(
				"MilkyWay",
				"MilkyWay is a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				nil,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "invalid service returns error",
			store: func(ctx sdk.Context) {
				err := suite.k.SetNextServiceID(ctx, 1)
				suite.Require().NoError(err)

				err = suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000))),
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgCreateService(
				"MilkyWay",
				types.DoNotModify,
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				nil,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "user without enough funds return error",
			store: func(ctx sdk.Context) {
				err := suite.k.SetNextServiceID(ctx, 1)
				suite.Require().NoError(err)

				err = suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000))),
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgCreateService(
				"MilkyWay",
				"MilkyWay is a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				nil,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "valid service is created properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetNextServiceID(ctx, 1)
				suite.Require().NoError(err)

				err = suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000))),
				))
				suite.Require().NoError(err)

				suite.fundAccount(ctx,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(200_000))),
				)
			},
			msg: types.NewMsgCreateService(
				"MilkyWay",
				"MilkyWay is a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000))),
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
				stored, found, err := suite.k.GetService(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().Equal(types.NewService(
					1,
					types.SERVICE_STATUS_CREATED,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
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
				poolBalance := suite.bk.GetBalance(ctx, authtypes.NewModuleAddress(distrtypes.ModuleName), "uatom")
				suite.Require().Equal(sdk.NewCoin("uatom", sdkmath.NewInt(100_000)), poolBalance)
			},
		},
		{
			name: "service is created and fee is charged - one of many fees denoms",
			store: func(ctx sdk.Context) {
				err := suite.k.SetNextServiceID(ctx, 1)
				suite.Require().NoError(err)

				err = suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(
						sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000)),
						sdk.NewCoin("utia", sdkmath.NewInt(30_000_000)),
						sdk.NewCoin("milktia", sdkmath.NewInt(80_000_000)),
					),
				))
				suite.Require().NoError(err)

				suite.fundAccount(ctx,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewCoins(
						sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000)),
						sdk.NewCoin("utia", sdkmath.NewInt(100_000_000)),
						sdk.NewCoin("milktia", sdkmath.NewInt(100_000_000)),
					),
				)
			},
			msg: types.NewMsgCreateService(
				"MilkyWay",
				"MilkyWay is a restaking platform",
				"https://milkyway.com",
				"https://milkyway.com/logo.png",
				sdk.NewCoins(
					sdk.NewCoin("uatom", sdkmath.NewInt(20_000_000)),
					sdk.NewCoin("utia", sdkmath.NewInt(15_000_000)),
					sdk.NewCoin("milktia", sdkmath.NewInt(80_000_000)),
				),
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
				stored, found, err := suite.k.GetService(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().Equal(types.NewService(
					1,
					types.SERVICE_STATUS_CREATED,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				), stored)

				// Make sure the service account has been created properly
				hasAccount := suite.ak.HasAccount(ctx, types.GetServiceAddress(1))
				suite.Require().True(hasAccount)

				// Make sure the next service id has been incremented
				nextServiceID, err := suite.k.GetNextServiceID(ctx)
				suite.Require().NoError(err)
				suite.Require().Equal(uint32(2), nextServiceID)

				// Make sure the user's funds were deducted
				userAddress, err := sdk.AccAddressFromBech32("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd")
				suite.Require().NoError(err)
				balance := suite.bk.GetAllBalances(ctx, userAddress)
				suite.Require().Equal(sdk.NewCoins(
					sdk.NewCoin("uatom", sdkmath.NewInt(80_000_000)),
					sdk.NewCoin("utia", sdkmath.NewInt(85_000_000)),
					sdk.NewCoin("milktia", sdkmath.NewInt(20_000_000)),
				), balance)

				// Make sure the community pool was funded
				poolBalance := suite.bk.GetAllBalances(ctx, authtypes.NewModuleAddress(distrtypes.ModuleName))
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
			suite.SetupTest()

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
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_CREATED,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
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
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_CREATED,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
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
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_CREATED,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
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
				stored, found, err := suite.k.GetService(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().Equal(types.NewService(
					1,
					types.SERVICE_STATUS_CREATED,
					"MilkyWay Modular Restaking",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
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
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_CREATED,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
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
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
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
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_CREATED,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
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
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
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
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
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
				stored, found, err := suite.k.GetService(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().Equal(types.NewService(
					1,
					types.SERVICE_STATUS_INACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
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

func (suite *KeeperTestSuite) TestMsgServer_DeleteService() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		msg       *types.MsgDeleteService
		shouldErr bool
		expEvents sdk.Events
		check     func(ctx sdk.Context)
	}{
		{
			name: "service not found returns error",
			msg: types.NewMsgDeleteService(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "non admin user returns error",
			store: func(ctx sdk.Context) {
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_INACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgDeleteService(
				1,
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
			),
			shouldErr: true,
		},
		{
			name: "active service can't be deleted",
			store: func(ctx sdk.Context) {
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgDeleteService(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "created service can be deleted",
			store: func(ctx sdk.Context) {
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_CREATED,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgDeleteService(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: false,
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeDeleteService,
					sdk.NewAttribute(types.AttributeKeyServiceID, "1"),
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the service was removed
				_, found, err := suite.k.GetService(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().False(found)
			},
		},
		{
			name: "inactive service is deleted properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_INACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgDeleteService(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: false,
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeDeleteService,
					sdk.NewAttribute(types.AttributeKeyServiceID, "1"),
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the service was removed
				_, found, err := suite.k.GetService(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().False(found)
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
			if tc.store != nil {
				tc.store(ctx)
			}

			msgServer := keeper.NewMsgServer(suite.k)
			res, err := msgServer.DeleteService(ctx, tc.msg)
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
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
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
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
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
				stored, found, err := suite.k.GetService(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().Equal(types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					false,
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
			res, err := msgServer.TransferServiceOwnership(ctx, tc.msg)
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
		name        string
		setup       func()
		store       func(ctx sdk.Context)
		setupCtx    func(ctx sdk.Context) sdk.Context
		msg         *types.MsgSetServiceParams
		shouldErr   bool
		expResponse *types.MsgSetServiceParamsResponse
		expEvents   sdk.Events
		check       func(ctx sdk.Context)
	}{
		{
			name: "not found service returns error",
			msg: types.NewMsgSetServiceParams(
				1,
				types.DefaultServiceParams(),
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "non admin sender returns error",
			store: func(ctx sdk.Context) {
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgSetServiceParams(
				1,
				types.DefaultServiceParams(),
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "set invalid params returns error",
			store: func(ctx sdk.Context) {
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgSetServiceParams(
				1,
				types.NewServiceParams([]string{"1stake"}),
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "service params updated successfully",
			store: func(ctx sdk.Context) {
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_ACTIVE,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					false,
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgSetServiceParams(
				1,
				types.NewServiceParams([]string{"umilk"}),
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr:   false,
			expResponse: &types.MsgSetServiceParamsResponse{},
			expEvents: []sdk.Event{
				sdk.NewEvent(
					types.EventTypeSetServiceParams,
					sdk.NewAttribute(types.AttributeKeyServiceID, "1"),
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the service was updated
				stored, err := suite.k.GetServiceParams(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().Equal(
					types.NewServiceParams([]string{"umilk"}),
					stored,
				)
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
			res, err := msgServer.SetServiceParams(ctx, tc.msg)
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

func (suite *KeeperTestSuite) TestMsgServer_AccreditService() {
	testCases := []struct {
		name        string
		setup       func()
		store       func(ctx sdk.Context)
		setupCtx    func(ctx sdk.Context) sdk.Context
		msg         *types.MsgAccreditService
		shouldErr   bool
		expResponse *types.MsgAccreditServiceResponse
		expEvents   sdk.Events
		check       func(ctx sdk.Context)
	}{
		{
			name: "invalid authority return error",
			msg: types.NewMsgAccreditService(
				1,
				"invalid",
			),
			shouldErr: true,
		},
		{
			name: "not found service returns error",
			msg: types.NewMsgAccreditService(
				1,
				authtypes.NewModuleAddress("gov").String(),
			),
			shouldErr: true,
		},
		{
			name: "valid data returns no error",
			store: func(ctx sdk.Context) {
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_CREATED,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"",
					false,
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgAccreditService(
				1,
				authtypes.NewModuleAddress("gov").String(),
			),
			shouldErr:   false,
			expResponse: &types.MsgAccreditServiceResponse{},
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeAccreditService,
					sdk.NewAttribute(types.AttributeKeyServiceID, "1"),
				),
			},
			check: func(ctx sdk.Context) {
				service, found, err := suite.k.GetService(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().True(service.Accredited)
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
			res, err := msgServer.AccreditService(ctx, tc.msg)
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

func (suite *KeeperTestSuite) TestMsgService_RevokeServiceAccreditation() {
	testCases := []struct {
		name        string
		setup       func()
		store       func(ctx sdk.Context)
		setupCtx    func(ctx sdk.Context) sdk.Context
		msg         *types.MsgRevokeServiceAccreditation
		shouldErr   bool
		expResponse *types.MsgRevokeServiceAccreditationResponse
		expEvents   sdk.Events
		check       func(ctx sdk.Context)
	}{
		{
			name: "invalid authority return error",
			msg: types.NewMsgRevokeServiceAccreditation(
				1,
				"invalid",
			),
			shouldErr: true,
		},
		{
			name: "not found service returns error",
			msg: types.NewMsgRevokeServiceAccreditation(
				1,
				authtypes.NewModuleAddress("gov").String(),
			),
			shouldErr: true,
		},
		{
			name: "valid data returns no error",
			store: func(ctx sdk.Context) {
				err := suite.k.SaveService(ctx, types.NewService(
					1,
					types.SERVICE_STATUS_CREATED,
					"MilkyWay",
					"MilkyWay is a restaking platform",
					"https://milkyway.com",
					"https://milkyway.com/logo.png",
					"",
					true,
				))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgRevokeServiceAccreditation(
				1,
				authtypes.NewModuleAddress("gov").String(),
			),
			shouldErr:   false,
			expResponse: &types.MsgRevokeServiceAccreditationResponse{},
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeRevokeServiceAccreditation,
					sdk.NewAttribute(types.AttributeKeyServiceID, "1"),
				),
			},
			check: func(ctx sdk.Context) {
				service, found, err := suite.k.GetService(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().False(service.Accredited)
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
			res, err := msgServer.RevokeServiceAccreditation(ctx, tc.msg)
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
