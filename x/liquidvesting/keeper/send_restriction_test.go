package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/milkyway-labs/milkyway/v10/x/liquidvesting/types"
	liquidvestingtypes "github.com/milkyway-labs/milkyway/v10/x/liquidvesting/types"
	operatorstypes "github.com/milkyway-labs/milkyway/v10/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/v10/x/pools/types"
	servicestypes "github.com/milkyway-labs/milkyway/v10/x/services/types"
)

func MustGetLockedDenom(denom string) string {
	lockedDenom, err := types.GetLockedRepresentationDenom(denom)
	if err != nil {
		panic(err)
	}

	return lockedDenom
}

func (suite *KeeperTestSuite) TestKeeper_BankSend() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		msg       *banktypes.MsgSend
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name: "send of locked tokens is not allowed",
			store: func(ctx sdk.Context) {
				suite.mintLockedRepresentation(
					ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)),
				)

				lockedDenom, err := types.GetLockedRepresentationDenom("stake")
				suite.Require().NoError(err)

				// Ensure we can transfer the locked representation
				suite.bk.SetSendEnabled(ctx, lockedDenom, true)
			},
			msg: banktypes.NewMsgSend(
				sdk.MustAccAddressFromBech32("cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"),
				sdk.MustAccAddressFromBech32("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
				sdk.NewCoins(sdk.NewInt64Coin(MustGetLockedDenom("stake"), 1000)),
			),
			shouldErr: true,
		},
		{
			name: "send locked tokens to pool is allowed",
			store: func(ctx sdk.Context) {
				suite.mintLockedRepresentation(
					ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)),
				)

				lockedDenom, err := types.GetLockedRepresentationDenom("stake")
				suite.Require().NoError(err)

				// Ensure we can transfer the locked representation
				suite.bk.SetSendEnabled(ctx, lockedDenom, true)

				// Create a pool for the locked denom
				err = suite.pk.SavePool(ctx, poolstypes.NewPool(1, lockedDenom))
				suite.Require().NoError(err)
			},
			msg: banktypes.NewMsgSend(
				sdk.MustAccAddressFromBech32("cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"),
				poolstypes.GetPoolAddress(1),
				sdk.NewCoins(sdk.NewInt64Coin(MustGetLockedDenom("stake"), 1000)),
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				lockedDenom, err := types.GetLockedRepresentationDenom("stake")
				suite.Require().NoError(err)

				pool, err := suite.pk.GetPool(ctx, 1)
				suite.Require().NoError(err)

				poolCoins := suite.bk.GetAllBalances(ctx, sdk.MustAccAddressFromBech32(pool.GetAddress()))
				suite.Require().Equal(sdk.NewCoins(sdk.NewInt64Coin(lockedDenom, 1000)), poolCoins)
			},
		},
		{
			name: "send locked tokens to operator is allowed",
			store: func(ctx sdk.Context) {
				suite.mintLockedRepresentation(
					ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)),
				)

				lockedDenom, err := types.GetLockedRepresentationDenom("stake")
				suite.Require().NoError(err)

				// Ensure we can transfer the locked representation
				suite.bk.SetSendEnabled(ctx, lockedDenom, true)

				// Create an operator
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
			},
			msg: banktypes.NewMsgSend(
				sdk.MustAccAddressFromBech32("cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"),
				operatorstypes.GetOperatorAddress(1),
				sdk.NewCoins(sdk.NewInt64Coin(MustGetLockedDenom("stake"), 1000)),
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				lockedDenom, err := types.GetLockedRepresentationDenom("stake")
				suite.Require().NoError(err)

				operator, err := suite.ok.GetOperator(ctx, 1)
				suite.Require().NoError(err)

				poolCoins := suite.bk.GetAllBalances(ctx, sdk.MustAccAddressFromBech32(operator.GetAddress()))
				suite.Require().Equal(sdk.NewCoins(sdk.NewInt64Coin(lockedDenom, 1000)), poolCoins)
			},
		},
		{
			name: "send locked tokens to service is allowed",
			store: func(ctx sdk.Context) {
				suite.mintLockedRepresentation(
					ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)),
				)

				lockedDenom, err := types.GetLockedRepresentationDenom("stake")
				suite.Require().NoError(err)

				// Ensure we can transfer the locked representation
				suite.bk.SetSendEnabled(ctx, lockedDenom, true)

				// Create a service
				service := servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_ACTIVE,
					"test operator",
					"",
					"",
					"",
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					true,
				)
				err = suite.sk.SaveService(ctx, service)
				suite.Require().NoError(err)
			},
			msg: banktypes.NewMsgSend(
				sdk.MustAccAddressFromBech32("cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"),
				servicestypes.GetServiceAddress(1),
				sdk.NewCoins(sdk.NewInt64Coin(MustGetLockedDenom("stake"), 1000)),
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				lockedDenom, err := types.GetLockedRepresentationDenom("stake")
				suite.Require().NoError(err)

				service, err := suite.sk.GetService(ctx, 1)
				suite.Require().NoError(err)

				poolCoins := suite.bk.GetAllBalances(ctx, sdk.MustAccAddressFromBech32(service.GetAddress()))
				suite.Require().Equal(sdk.NewCoins(sdk.NewInt64Coin(lockedDenom, 1000)), poolCoins)
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

			msgServer := bankkeeper.NewMsgServerImpl(suite.bk)
			_, err := msgServer.Send(ctx, tc.msg)
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

func (suite *KeeperTestSuite) TestKeeper_SendRestrictionFn() {
	testCase := []struct {
		name      string
		store     func(ctx sdk.Context)
		from      sdk.AccAddress
		to        sdk.AccAddress
		amount    sdk.Coins
		shouldErr bool
		expTo     sdk.AccAddress
	}{
		{
			name:      "sending normal coins from user to user works properly",
			from:      sdk.MustAccAddressFromBech32("cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"),
			to:        sdk.MustAccAddressFromBech32("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			amount:    sdk.NewCoins(sdk.NewInt64Coin("ibc/1", 100)),
			shouldErr: false,
			expTo:     sdk.MustAccAddressFromBech32("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
		},
		{
			name:      "sending locked representation from user to user returns error",
			from:      sdk.MustAccAddressFromBech32("cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre"),
			to:        sdk.MustAccAddressFromBech32("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			amount:    sdk.NewCoins(sdk.NewInt64Coin("locked/stake", 100)),
			shouldErr: true,
		},
		{
			name: "sending normal coins between restaking targets is not allowed",
			store: func(ctx sdk.Context) {
				// Create a test service and operator
				suite.createService(ctx, 1)
				suite.createOperator(ctx, 1)
			},
			from:      servicestypes.GetServiceAddress(1),
			to:        operatorstypes.GetOperatorAddress(1),
			amount:    sdk.NewCoins(sdk.NewInt64Coin("stake", 100)),
			shouldErr: true,
		},
		{
			name: "sending normal coins between restaking targets is not allowed",
			store: func(ctx sdk.Context) {
				// Create a test service and operator
				suite.createPool(ctx, 1, LockedIBCDenom)
				suite.createOperator(ctx, 1)
			},
			from:      poolstypes.GetPoolAddress(1),
			to:        operatorstypes.GetOperatorAddress(1),
			amount:    sdk.NewCoins(sdk.NewInt64Coin("locked/stake", 100)),
			shouldErr: true,
		},
		{
			name: "sending coins from the module account is allowed",
			from: authtypes.NewModuleAddress(liquidvestingtypes.ModuleName),
			to:   sdk.MustAccAddressFromBech32("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			amount: sdk.NewCoins(
				sdk.NewInt64Coin("ibc/1", 100),
				sdk.NewInt64Coin("locked/stake", 200),
			),
			shouldErr: false,
			expTo:     sdk.MustAccAddressFromBech32("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
		},
		{
			name: "sending coins to the module account is allowed",
			from: sdk.MustAccAddressFromBech32("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			to:   authtypes.NewModuleAddress(liquidvestingtypes.ModuleName),
			amount: sdk.NewCoins(
				sdk.NewInt64Coin("ibc/1", 100),
				sdk.NewInt64Coin("locked/stake", 200),
			),
			shouldErr: false,
			expTo:     authtypes.NewModuleAddress(liquidvestingtypes.ModuleName),
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

			receivedTo, err := suite.k.SendRestrictionFn(ctx, tc.from, tc.to, tc.amount)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expTo, receivedTo)
			}
		})
	}
}
