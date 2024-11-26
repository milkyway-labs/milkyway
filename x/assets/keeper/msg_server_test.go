package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v2/x/assets/keeper"
	"github.com/milkyway-labs/milkyway/v2/x/assets/types"
)

func (suite *KeeperTestSuite) TestMsgServer_RegisterAsset() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		setupCtx  func(ctx sdk.Context) sdk.Context
		msg       *types.MsgRegisterAsset
		shouldErr bool
		expEvents sdk.Events
		check     func(ctx sdk.Context)
	}{
		{
			name:      "invalid asset returns error",
			msg:       types.NewMsgRegisterAsset(suite.authority, types.NewAsset("umilk", "@#$%", 0)),
			shouldErr: true,
		},
		{
			name:      "invalid authority returns error",
			msg:       types.NewMsgRegisterAsset("invalid", types.NewAsset("umilk", "MILK", 6)),
			shouldErr: true,
		},
		{
			name: "valid asset is registered properly",
			msg:  types.NewMsgRegisterAsset(suite.authority, types.NewAsset("umilk", "MILK", 6)),
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeRegisterAsset,
					sdk.NewAttribute(types.AttributeKeyDenom, "umilk"),
					sdk.NewAttribute(types.AttributeKeyTicker, "MILK"),
					sdk.NewAttribute(types.AttributeKeyExponent, "6"),
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the asset is stored
				asset, err := suite.keeper.GetAsset(ctx, "umilk")
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewAsset("umilk", "MILK", 6), asset)
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
			_, err := msgServer.RegisterAsset(ctx, tc.msg)
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

func (suite *KeeperTestSuite) TestMsgServer_DeregisterAsset() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		setupCtx  func(ctx sdk.Context) sdk.Context
		msg       *types.MsgDeregisterAsset
		shouldErr bool
		expEvents sdk.Events
		check     func(ctx sdk.Context)
	}{
		{
			name:      "invalid authority returns error",
			msg:       types.NewMsgDeregisterAsset("invalid", "umilk"),
			shouldErr: true,
		},
		{
			name:      "asset not found returns error",
			msg:       types.NewMsgDeregisterAsset(suite.authority, "umilk"),
			shouldErr: true,
		},
		{
			name: "valid asset is deregistered properly",
			store: func(ctx sdk.Context) {
				err := suite.keeper.SetAsset(ctx, types.NewAsset("umilk", "MILK", 6))
				suite.Require().NoError(err)
			},
			msg: types.NewMsgDeregisterAsset(suite.authority, "umilk"),
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeDeregisterAsset,
					sdk.NewAttribute(types.AttributeKeyDenom, "umilk"),
				),
			},
			check: func(ctx sdk.Context) {
				// Make sure the asset is removed
				_, err := suite.keeper.GetAsset(ctx, "umilk")
				suite.Require().Error(err)
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
			_, err := msgServer.DeregisterAsset(ctx, tc.msg)
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
