package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/keeper"
	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
)

func (suite *KeeperTestSuite) TestMsgServer_MintVestedRepresentation() {
	burnerAccount := "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"
	minterAccount := "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"
	testAccount := "cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn"

	testCases := []struct {
		name      string
		msg       *types.MsgMintVestedRepresentation
		shouldErr bool
		expEvents sdk.Events
		check     func(ctx sdk.Context)
	}{
		{
			name: "burner can't mint",
			msg: types.NewMsgMintVestedRepresentation(burnerAccount, testAccount,
				sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 1000))),
			shouldErr: true,
		},
		{
			name: "can't mint vested representation of vested representation coin",
			msg: types.NewMsgMintVestedRepresentation(minterAccount, testAccount,
				sdk.NewCoins(sdk.NewInt64Coin(vestedDenom, 1000))),
			shouldErr: true,
		},
		{
			name: "mint properly",
			msg: types.NewMsgMintVestedRepresentation(minterAccount, testAccount,
				sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 1000))),
			shouldErr: false,
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeMintVestedRepresentation,
					sdk.NewAttribute(sdk.AttributeKeySender, minterAccount),
					sdk.NewAttribute(sdk.AttributeKeyAmount, "1000"+vestedDenom),
					sdk.NewAttribute(types.AttributeKeyReceiver, testAccount),
				),
			},
			check: func(ctx sdk.Context) {
				balances := suite.bk.GetAllBalances(ctx, sdk.MustAccAddressFromBech32(testAccount))
				suite.Assert().Equal(
					sdk.NewCoins(sdk.NewInt64Coin(vestedDenom, 1000)),
					balances,
				)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()

			suite.Assert().NoError(suite.k.SetParams(ctx, types.NewParams(
				math.LegacyMustNewDecFromStr("2.0"),
				[]string{burnerAccount},
				[]string{minterAccount},
			)))

			msgServer := keeper.NewMsgServer(suite.k)
			_, err := msgServer.MintVestedRepresentation(ctx, tc.msg)

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

func (suite *KeeperTestSuite) TestMsgServer_BurnVestedRepresentation() {
	burnerAccount := "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"
	minterAccount := "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"
	testAccount := "cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn"

	testCases := []struct {
		name      string
		setup     func(ctx sdk.Context)
		msg       *types.MsgBurnVestedRepresentation
		shouldErr bool
		expEvents sdk.Events
		check     func(ctx sdk.Context)
	}{
		{
			name: "minter can't burn",
			msg: types.NewMsgBurnVestedRepresentation(minterAccount, testAccount,
				sdk.NewCoins(sdk.NewInt64Coin(vestedDenom, 1000))),
			setup: func(ctx sdk.Context) {
				suite.mintVestedRepresentation(testAccount,
					sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 1000)))
			},
			shouldErr: true,
		},
		{
			name: "can't burn normal coins",
			msg: types.NewMsgBurnVestedRepresentation(burnerAccount, testAccount,
				sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 1000))),
			setup: func(ctx sdk.Context) {
				suite.fundAccount(ctx, testAccount,
					sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 1000)))
				suite.mintVestedRepresentation(testAccount,
					sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 1000)))
			},
			shouldErr: true,
		},
		{
			name: "burn properly",
			msg: types.NewMsgBurnVestedRepresentation(burnerAccount, testAccount,
				sdk.NewCoins(sdk.NewInt64Coin(vestedDenom, 1000))),
			shouldErr: false,
			setup: func(ctx sdk.Context) {
				suite.mintVestedRepresentation(testAccount,
					sdk.NewCoins(sdk.NewInt64Coin(iBCDenom, 1000)))
			},
			expEvents: sdk.Events{
				sdk.NewEvent(
					types.EventTypeBurnVestedRepresentation,
					sdk.NewAttribute(sdk.AttributeKeySender, burnerAccount),
					sdk.NewAttribute(sdk.AttributeKeyAmount, "1000"+vestedDenom),
					sdk.NewAttribute(types.AttributeKeyUser, testAccount),
				),
			},
			check: func(ctx sdk.Context) {
				balances := suite.bk.GetAllBalances(ctx, sdk.MustAccAddressFromBech32(testAccount))
				suite.Assert().Equal(
					sdk.NewCoins(),
					balances,
				)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
			suite.Assert().NoError(suite.k.SetParams(ctx, types.NewParams(
				math.LegacyMustNewDecFromStr("2.0"),
				[]string{burnerAccount},
				[]string{minterAccount},
			)))

			if tc.setup != nil {
				tc.setup(ctx)
			}

			msgServer := keeper.NewMsgServer(suite.k)
			_, err := msgServer.BurnVestedRepresentation(ctx, tc.msg)

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
