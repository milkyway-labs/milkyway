package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/milkyway-labs/milkyway/v6/x/liquidvesting/keeper"
	"github.com/milkyway-labs/milkyway/v6/x/liquidvesting/types"
)

func (suite *KeeperTestSuite) TestMsgServer_MintLockedRepresentation() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		msg       *types.MsgMintLockedRepresentation
		shouldErr bool
		expEvents sdk.Events
		check     func(ctx sdk.Context)
	}{
		{
			name: "burner can't mint",
			store: func(ctx sdk.Context) {
				// Store the params
				err := suite.k.SetParams(ctx, types.NewParams(
					math.LegacyMustNewDecFromStr("2.0"),
					[]string{"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"},
					[]string{"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"},
					nil,
				))
				suite.Assert().NoError(err)
			},
			msg: types.NewMsgMintLockedRepresentation(
				"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
				sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
			),
			shouldErr: true,
		},
		{
			name: "can't mint locked representation of locked representation coin",
			store: func(ctx sdk.Context) {
				// Store the params
				err := suite.k.SetParams(ctx, types.NewParams(
					math.LegacyMustNewDecFromStr("2.0"),
					[]string{"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"},
					[]string{"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"},
					nil,
				))
				suite.Assert().NoError(err)
			},
			msg: types.NewMsgMintLockedRepresentation(
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
				sdk.NewCoins(sdk.NewInt64Coin(LockedIBCDenom, 1000)),
			),
			shouldErr: true,
		},
		{
			name: "mint properly",
			store: func(ctx sdk.Context) {
				// Store the params
				err := suite.k.SetParams(ctx, types.NewParams(
					math.LegacyMustNewDecFromStr("2.0"),
					[]string{"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"},
					[]string{"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"},
					nil,
				))
				suite.Assert().NoError(err)
			},
			msg: types.NewMsgMintLockedRepresentation(
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
				sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
			),
			shouldErr: false,
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeMintLockedRepresentation,
					sdk.NewAttribute(sdk.AttributeKeySender, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
					sdk.NewAttribute(sdk.AttributeKeyAmount, "1000"+LockedIBCDenom),
					sdk.NewAttribute(types.AttributeKeyReceiver, "cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn"),
				),
			},
			check: func(ctx sdk.Context) {
				balances := suite.bk.GetAllBalances(ctx, sdk.MustAccAddressFromBech32("cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn"))
				suite.Assert().Equal(
					sdk.NewCoins(sdk.NewInt64Coin(LockedIBCDenom, 1000)),
					balances,
				)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			// Cache the context
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			msgServer := keeper.NewMsgServer(suite.k)
			_, err := msgServer.MintLockedRepresentation(ctx, tc.msg)

			if tc.shouldErr {
				suite.Assert().Error(err)
			} else {
				suite.Assert().NoError(err)
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

func (suite *KeeperTestSuite) TestMsgServer_BurnLockedRepresentation() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		msg       *types.MsgBurnLockedRepresentation
		shouldErr bool
		expEvents sdk.Events
		check     func(ctx sdk.Context)
	}{
		{
			name: "minter can't burn",
			store: func(ctx sdk.Context) {
				err := suite.k.SetParams(ctx, types.NewParams(
					math.LegacyMustNewDecFromStr("2.0"),
					[]string{"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"},
					[]string{"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"},
					nil,
				))
				suite.Assert().NoError(err)

				suite.mintLockedRepresentation(ctx,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
				)
			},
			msg: types.NewMsgBurnLockedRepresentation(
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
				sdk.NewCoins(sdk.NewInt64Coin(LockedIBCDenom, 1000)),
			),

			shouldErr: true,
		},
		{
			name: "can't burn normal coins",
			store: func(ctx sdk.Context) {
				err := suite.k.SetParams(ctx, types.NewParams(
					math.LegacyMustNewDecFromStr("2.0"),
					[]string{"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"},
					[]string{"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"},
					nil,
				))
				suite.Assert().NoError(err)

				suite.fundAccount(ctx,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
				)
				suite.mintLockedRepresentation(ctx,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
				)
			},
			msg: types.NewMsgBurnLockedRepresentation(
				"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
				sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
			),
			shouldErr: true,
		},
		{
			name: "burn properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetParams(ctx, types.NewParams(
					math.LegacyMustNewDecFromStr("2.0"),
					[]string{"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"},
					[]string{"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"},
					nil,
				))
				suite.Assert().NoError(err)

				suite.mintLockedRepresentation(ctx,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
				)
			},
			msg: types.NewMsgBurnLockedRepresentation(
				"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
				sdk.NewCoins(sdk.NewInt64Coin(LockedIBCDenom, 1000)),
			),
			shouldErr: false,
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeBurnLockedRepresentation,
					sdk.NewAttribute(sdk.AttributeKeySender, "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"),
					sdk.NewAttribute(sdk.AttributeKeyAmount, "1000"+LockedIBCDenom),
					sdk.NewAttribute(types.AttributeKeyUser, "cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn"),
				),
			},
			check: func(ctx sdk.Context) {
				userAddr, err := sdk.AccAddressFromBech32("cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn")
				suite.Require().NoError(err)

				userBalance := suite.bk.GetAllBalances(ctx, userAddr)
				suite.Assert().Empty(userBalance)
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
			_, err := msgServer.BurnLockedRepresentation(ctx, tc.msg)

			if tc.shouldErr {
				suite.Assert().Error(err)
			} else {
				suite.Assert().NoError(err)
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

func (suite *KeeperTestSuite) TestMsgServer_WithdrawInsuranceFund() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		msg       *types.MsgWithdrawInsuranceFund
		shouldErr bool
		check     func(ctx sdk.Context)
		expEvents sdk.Events
	}{
		{
			name: "can't withdraw without deposits",
			msg: types.NewMsgWithdrawInsuranceFund(
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
				sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 100)),
			),
			shouldErr: true,
		},
		{
			name: "can't withdraw more then deposited",
			store: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx, "cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn", sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10)))
			},
			msg:       types.NewMsgWithdrawInsuranceFund("cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn", sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 100))),
			shouldErr: true,
		},
		{
			name: "can't withdraw more then available",
			store: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10)),
				)

				// Delegate to pool to simulate insurance fund utilization
				suite.mintLockedRepresentation(ctx,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 1000)),
				)

				suite.createPool(ctx, 1, LockedIBCDenom)

				_, err := suite.rk.DelegateToPool(ctx,
					sdk.NewInt64Coin(LockedIBCDenom, 500),
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
				)
				suite.Assert().NoError(err)
			},
			msg: types.NewMsgWithdrawInsuranceFund(
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
				sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10)),
			),
			shouldErr: true,
		},
		{
			name: "withdraw correctly",
			store: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10)),
				)
			},
			msg: types.NewMsgWithdrawInsuranceFund(
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
				sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10)),
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				userAddr, err := sdk.AccAddressFromBech32("cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn")
				suite.Require().NoError(err)

				balances := suite.bk.GetAllBalances(ctx, userAddr)
				suite.Assert().Equal(sdk.NewCoins(sdk.NewInt64Coin(IBCDenom, 10)), balances)
			},
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeWithdrawInsuranceFund,
					sdk.NewAttribute(sdk.AttributeKeySender, "cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn"),
					sdk.NewAttribute(sdk.AttributeKeyAmount, "10"+IBCDenom),
				),
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			msgServer := keeper.NewMsgServer(suite.k)
			_, err := msgServer.WithdrawInsuranceFund(ctx, tc.msg)

			if tc.shouldErr {
				suite.Assert().Error(err)
			} else {
				suite.Assert().NoError(err)
				for _, event := range tc.expEvents {
					suite.Require().Contains(ctx.EventManager().Events(), event)
				}
			}

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgServer_UpdateParams() {
	testCases := []struct {
		name      string
		setup     func(ctx sdk.Context)
		msg       *types.MsgUpdateParams
		shouldErr bool
		expEvents sdk.Events
		check     func(ctx sdk.Context)
	}{
		{
			name: "invalid authority return error",
			msg: types.NewMsgUpdateParams(
				"invalid",
				types.DefaultParams(),
			),
			shouldErr: true,
		},
		{
			name: "valid data returns no error",
			msg: types.NewMsgUpdateParams(
				authtypes.NewModuleAddress("gov").String(),
				types.DefaultParams(),
			),
			shouldErr: false,
			expEvents: sdk.Events{},
			check: func(ctx sdk.Context) {
				params, err := suite.k.GetParams(ctx)
				suite.Assert().NoError(err)
				suite.Assert().Equal(types.DefaultParams(), params)
			},
		},
		{
			name: "invalid allowed channels returns error",
			msg: types.NewMsgUpdateParams(
				authtypes.NewModuleAddress("gov").String(),
				types.NewParams(math.LegacyNewDec(2), nil, nil, []string{"invalid-channel"}),
			),
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
			if tc.setup != nil {
				tc.setup(ctx)
			}

			msgServer := keeper.NewMsgServer(suite.k)
			_, err := msgServer.UpdateParams(ctx, tc.msg)

			if tc.shouldErr {
				suite.Assert().Error(err)
			} else {
				suite.Assert().NoError(err)
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
