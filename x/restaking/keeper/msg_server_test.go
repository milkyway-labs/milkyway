package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	"github.com/milkyway-labs/milkyway/x/restaking/keeper"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

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
					sdk.NewAttribute(types.AttributeKeyOperatorID, "1"),
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
					sdk.NewAttribute(types.AttributeKeyServiceID, "1"),
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