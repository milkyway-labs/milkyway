package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	bankkeeper "github.com/milkyway-labs/milkyway/x/bank/keeper"
	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
	liquidvestingtypes "github.com/milkyway-labs/milkyway/x/liquidvesting/types"
	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

func (suite *KeeperTestSuite) TestSendRestriction() {
	testCases := []struct {
		name         string
		buildMessage func(ctx sdk.Context) *banktypes.MsgSend
		shouldErr    bool
		check        func(ctx sdk.Context)
	}{
		{
			name: "send of vested tokens is not allowed",
			buildMessage: func(ctx sdk.Context) *banktypes.MsgSend {
				coins := sdk.NewCoins(sdk.NewInt64Coin("stake", 1000))
				suite.mintVestedRepresentation("cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre", coins)
				vestedDenom, err := types.GetVestedRepresentationDenom("stake")
				suite.Require().NoError(err)

				// Ensure we can transfer the vested representation
				suite.bk.SetSendEnabled(ctx, vestedDenom, true)

				senderAccAddr := sdk.MustAccAddressFromBech32("cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre")
				receiverAccAddr := sdk.MustAccAddressFromBech32("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				return banktypes.NewMsgSend(senderAccAddr, receiverAccAddr, sdk.NewCoins(sdk.NewInt64Coin(vestedDenom, 1000)))
			},
			shouldErr: true,
		},
		{
			name: "send vested tokens to pool is allowed",
			buildMessage: func(ctx sdk.Context) *banktypes.MsgSend {
				coins := sdk.NewCoins(sdk.NewInt64Coin("stake", 1000))
				suite.mintVestedRepresentation("cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre", coins)
				vestedDenom, err := types.GetVestedRepresentationDenom("stake")
				suite.Require().NoError(err)

				// Ensure we can transfer the vested representation
				suite.bk.SetSendEnabled(ctx, vestedDenom, true)

				// Create a pool for the vested denom
				pool := poolstypes.NewPool(1, vestedDenom)
				err = suite.pk.SavePool(ctx, pool)
				suite.Require().NoError(err)

				senderAccAddr := sdk.MustAccAddressFromBech32("cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre")
				receiverAccAddr := sdk.MustAccAddressFromBech32(pool.GetAddress())
				return banktypes.NewMsgSend(senderAccAddr, receiverAccAddr, sdk.NewCoins(sdk.NewInt64Coin(vestedDenom, 1000)))
			},
			shouldErr: false,
			check: func(ctx sdk.Context) {
				vestedDenom, err := types.GetVestedRepresentationDenom("stake")
				suite.Require().NoError(err)

				pool, found := suite.pk.GetPool(ctx, 1)
				suite.Require().True(found)

				poolCoins := suite.bk.GetAllBalances(ctx, sdk.MustAccAddressFromBech32(pool.GetAddress()))
				suite.Require().Equal(sdk.NewCoins(sdk.NewInt64Coin(vestedDenom, 1000)), poolCoins)
			},
		},
		{
			name: "send vested tokens to operator is allowed",
			buildMessage: func(ctx sdk.Context) *banktypes.MsgSend {
				coins := sdk.NewCoins(sdk.NewInt64Coin("stake", 1000))
				suite.mintVestedRepresentation("cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre", coins)
				vestedDenom, err := types.GetVestedRepresentationDenom("stake")
				suite.Require().NoError(err)

				// Ensure we can transfer the vested representation
				suite.bk.SetSendEnabled(ctx, vestedDenom, true)

				// Create a operator
				operator := operatorstypes.NewOperator(
					1,
					operatorstypes.OPERATOR_STATUS_ACTIVE,
					"test operator",
					"",
					"",
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				)
				err = suite.ok.SaveOperator(ctx, operator)
				suite.Require().NoError(err)

				senderAccAddr := sdk.MustAccAddressFromBech32("cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre")
				receiverAccAddr := sdk.MustAccAddressFromBech32(operator.GetAddress())
				return banktypes.NewMsgSend(senderAccAddr, receiverAccAddr, sdk.NewCoins(sdk.NewInt64Coin(vestedDenom, 1000)))
			},
			shouldErr: false,
			check: func(ctx sdk.Context) {
				vestedDenom, err := types.GetVestedRepresentationDenom("stake")
				suite.Require().NoError(err)

				operator, found := suite.ok.GetOperator(ctx, 1)
				suite.Require().True(found)

				poolCoins := suite.bk.GetAllBalances(ctx, sdk.MustAccAddressFromBech32(operator.GetAddress()))
				suite.Require().Equal(sdk.NewCoins(sdk.NewInt64Coin(vestedDenom, 1000)), poolCoins)
			},
		},
		{
			name: "send vested tokens to service is allowed",
			buildMessage: func(ctx sdk.Context) *banktypes.MsgSend {
				coins := sdk.NewCoins(sdk.NewInt64Coin("stake", 1000))
				suite.mintVestedRepresentation("cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre", coins)
				vestedDenom, err := types.GetVestedRepresentationDenom("stake")
				suite.Require().NoError(err)

				// Ensure we can transfer the vested representation
				suite.bk.SetSendEnabled(ctx, vestedDenom, true)

				// Create a service
				service := servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_ACTIVE,
					"test operator",
					"",
					"",
					"",
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				)
				err = suite.sk.SaveService(ctx, service)
				suite.Require().NoError(err)

				senderAccAddr := sdk.MustAccAddressFromBech32("cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre")
				receiverAccAddr := sdk.MustAccAddressFromBech32(service.GetAddress())
				return banktypes.NewMsgSend(senderAccAddr, receiverAccAddr, sdk.NewCoins(sdk.NewInt64Coin(vestedDenom, 1000)))
			},
			shouldErr: false,
			check: func(ctx sdk.Context) {
				vestedDenom, err := types.GetVestedRepresentationDenom("stake")
				suite.Require().NoError(err)

				service, found := suite.sk.GetService(ctx, 1)
				suite.Require().True(found)

				poolCoins := suite.bk.GetAllBalances(ctx, sdk.MustAccAddressFromBech32(service.GetAddress()))
				suite.Require().Equal(sdk.NewCoins(sdk.NewInt64Coin(vestedDenom, 1000)), poolCoins)
			},
		},
	}

	for _, tc := range testCases {
		suite.SetupTest()
		ctx, _ := suite.ctx.CacheContext()
		suite.Run(tc.name, func() {
			msg := tc.buildMessage(ctx)
			msgServer := bankkeeper.NewMsgServerImpl(suite.bk)

			_, err := msgServer.Send(ctx, msg)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				if tc.check != nil {
					tc.check(ctx)
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_SendRegistrionFn() {
	testCase := []struct {
		name      string
		store     func(ctx sdk.Context)
		from      string
		to        string
		amount    sdk.Coins
		shouldErr bool
		expTo     string
	}{
		{
			name:      "sending normal coins from user to user works properly",
			from:      "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			to:        "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			amount:    sdk.NewCoins(sdk.NewInt64Coin("ibc/1", 100)),
			shouldErr: false,
			expTo:     "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
		},
		{
			name:      "sending vested representation from user to user returns error",
			from:      "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			to:        "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			amount:    sdk.NewCoins(sdk.NewInt64Coin("vested/stake", 100)),
			shouldErr: true,
		},
		{
			name: "sending normal coins between restaking targets is not allowed",
			store: func(ctx sdk.Context) {
				// Create a test service and operator
				suite.createService(testServiceId)
				suite.createOperator(testOperatorId)
			},
			from:      servicestypes.GetServiceAddress(testServiceId).String(),
			to:        operatorstypes.GetOperatorAddress(testOperatorId).String(),
			amount:    sdk.NewCoins(sdk.NewInt64Coin("stake", 100)),
			shouldErr: true,
		},
		{
			name: "sending normal coins between restaking targets is not allowed",
			store: func(ctx sdk.Context) {
				// Create a test service and operator
				suite.createPool(testPoolId, vestedIBCDenom)
				suite.createOperator(testOperatorId)
			},
			from:      poolstypes.GetPoolAddress(testPoolId).String(),
			to:        operatorstypes.GetOperatorAddress(testOperatorId).String(),
			amount:    sdk.NewCoins(sdk.NewInt64Coin("vested/stake", 100)),
			shouldErr: true,
		},
		{
			name: "sending coins from the module account is allowed",
			from: authtypes.NewModuleAddress(liquidvestingtypes.ModuleName).String(),
			to:   "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			amount: sdk.NewCoins(
				sdk.NewInt64Coin("ibc/1", 100),
				sdk.NewInt64Coin("vested/stake", 200),
			),
			shouldErr: false,
			expTo:     "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
		},
		{
			name: "sending coins to the module account is allowed",
			from: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			to:   authtypes.NewModuleAddress(liquidvestingtypes.ModuleName).String(),
			amount: sdk.NewCoins(
				sdk.NewInt64Coin("ibc/1", 100),
				sdk.NewInt64Coin("vested/stake", 200),
			),
			shouldErr: false,
			expTo:     authtypes.NewModuleAddress(liquidvestingtypes.ModuleName).String(),
		},
	}

	for _, tc := range testCase {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			from, err := sdk.AccAddressFromBech32(tc.from)
			suite.Require().NoError(err)

			to, err := sdk.AccAddressFromBech32(tc.to)
			suite.Require().NoError(err)

			receivedTo, err := suite.k.SendRestrictionFn(ctx, from, to, tc.amount)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expTo, receivedTo.String())
			}
		})
	}
}
