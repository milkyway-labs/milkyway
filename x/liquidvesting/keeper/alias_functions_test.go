package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	operatorstypes "github.com/milkyway-labs/milkyway/v7/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/v7/x/pools/types"
	restakingkeeper "github.com/milkyway-labs/milkyway/v7/x/restaking/keeper"
	restakingtypes "github.com/milkyway-labs/milkyway/v7/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v7/x/services/types"
)

func (suite *KeeperTestSuite) TestKeeper_GetAllUserRestakedLockedRepresentations() {
	testCases := []struct {
		name     string
		user     string
		setup    func(sdk.Context)
		expected sdk.DecCoins
	}{
		{
			name:     "user has no restaked vesting representations",
			user:     "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			expected: sdk.NewDecCoins(),
		},
		{
			name: "computed amount is correct",
			user: "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			setup: func(ctx sdk.Context) {
				pool := poolstypes.NewPool(1, "locked/stake")
				err := suite.pk.SavePool(ctx, pool)
				suite.Require().NoError(err)

				service := servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_ACTIVE,
					"test service",
					"",
					"",
					"",
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					false,
				)
				err = suite.sk.SaveService(ctx, service)
				suite.Require().NoError(err)

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

				// Fund the account
				suite.fundAccountInsuranceFund(
					ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin("stake", 20)),
				)
				suite.mintLockedRepresentation(
					ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)),
				)

				restakingService := restakingkeeper.NewMsgServer(suite.rk)

				// Perform some delegations
				_, err = restakingService.DelegatePool(ctx, restakingtypes.NewMsgDelegatePool(
					sdk.NewInt64Coin("locked/stake", 200),
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				))
				suite.Require().NoError(err)

				_, err = restakingService.DelegateService(ctx, restakingtypes.NewMsgDelegateService(
					1,
					sdk.NewCoins(sdk.NewInt64Coin("locked/stake", 200)),
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				))
				suite.Require().NoError(err)

				_, err = restakingService.DelegateOperator(ctx, restakingtypes.NewMsgDelegateOperator(
					1,
					sdk.NewCoins(sdk.NewInt64Coin("locked/stake", 600)),
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				))
				suite.Require().NoError(err)
			},
			expected: sdk.NewDecCoins(sdk.NewInt64DecCoin("locked/stake", 1000)),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, _ := suite.ctx.CacheContext()

			if tc.setup != nil {
				tc.setup(ctx)
			}

			coins, err := suite.k.GetAllUserRestakedLockedRepresentations(ctx, tc.user)
			suite.Require().NoError(err)
			suite.Require().Equal(tc.expected, coins)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_GetAllUserUnbondingLockedRepresentations() {
	testCases := []struct {
		name     string
		user     string
		setup    func(sdk.Context)
		expected sdk.Coins
	}{
		{
			name:     "user has no vesting representations undelegations",
			user:     "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			expected: sdk.NewCoins(),
		},
		{
			name: "computed amount is correct",
			user: "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			setup: func(ctx sdk.Context) {
				pool := poolstypes.NewPool(1, "locked/stake")
				err := suite.pk.SavePool(ctx, pool)
				suite.Require().NoError(err)

				service := servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_ACTIVE,
					"test service",
					"",
					"",
					"",
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					true,
				)
				err = suite.sk.SaveService(ctx, service)
				suite.Require().NoError(err)

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

				// Fund the account
				suite.fundAccountInsuranceFund(
					ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin("stake", 20)),
				)
				suite.mintLockedRepresentation(
					ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)),
				)

				restakingService := restakingkeeper.NewMsgServer(suite.rk)

				// Perform some delegations
				_, err = restakingService.DelegatePool(ctx, restakingtypes.NewMsgDelegatePool(
					sdk.NewInt64Coin("locked/stake", 200),
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				))
				suite.Require().NoError(err)

				_, err = restakingService.DelegateService(ctx, restakingtypes.NewMsgDelegateService(
					1,
					sdk.NewCoins(sdk.NewInt64Coin("locked/stake", 200)),
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				))
				suite.Require().NoError(err)

				_, err = restakingService.DelegateOperator(ctx, restakingtypes.NewMsgDelegateOperator(
					1,
					sdk.NewCoins(sdk.NewInt64Coin("locked/stake", 600)),
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				))
				suite.Require().NoError(err)

				// Perform some undelegations
				_, err = restakingService.UndelegatePool(ctx, restakingtypes.NewMsgUndelegatePool(
					sdk.NewInt64Coin("locked/stake", 100),
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				))
				suite.Require().NoError(err)

				_, err = restakingService.UndelegateService(ctx, restakingtypes.NewMsgUndelegateService(
					1,
					sdk.NewCoins(sdk.NewInt64Coin("locked/stake", 100)),
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				))
				suite.Require().NoError(err)

				_, err = restakingService.UndelegateOperator(ctx, restakingtypes.NewMsgUndelegateOperator(
					1,
					sdk.NewCoins(sdk.NewInt64Coin("locked/stake", 400)),
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				))
				suite.Require().NoError(err)
			},
			expected: sdk.NewCoins(sdk.NewInt64Coin("locked/stake", 600)),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, _ := suite.ctx.CacheContext()

			if tc.setup != nil {
				tc.setup(ctx)
			}

			coins := suite.k.GetAllUserUnbondingLockedRepresentations(ctx, tc.user)
			suite.Require().Equal(tc.expected, coins)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_GetAllUserActiveLockedRepresentations() {
	testCases := []struct {
		name     string
		user     string
		setup    func(sdk.Context)
		expected sdk.DecCoins
	}{
		{
			name: "user has no vesting representations delegations and undelegations",
			user: "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
		},
		{
			name: "computed amount is correct",
			user: "cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
			setup: func(ctx sdk.Context) {
				pool := poolstypes.NewPool(1, "locked/stake")
				err := suite.pk.SavePool(ctx, pool)
				suite.Require().NoError(err)

				service := servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_ACTIVE,
					"test service",
					"",
					"",
					"",
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					true,
				)
				err = suite.sk.SaveService(ctx, service)
				suite.Require().NoError(err)

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

				// Fund the account
				suite.fundAccountInsuranceFund(
					ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin("stake", 20)),
				)
				suite.mintLockedRepresentation(
					ctx,
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
					sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)),
				)

				restakingService := restakingkeeper.NewMsgServer(suite.rk)

				// Perform some delegations
				_, err = restakingService.DelegatePool(ctx, restakingtypes.NewMsgDelegatePool(
					sdk.NewInt64Coin("locked/stake", 200),
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				))
				suite.Require().NoError(err)

				_, err = restakingService.DelegateService(ctx, restakingtypes.NewMsgDelegateService(
					1,
					sdk.NewCoins(sdk.NewInt64Coin("locked/stake", 200)),
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				))
				suite.Require().NoError(err)

				_, err = restakingService.DelegateOperator(ctx, restakingtypes.NewMsgDelegateOperator(
					1,
					sdk.NewCoins(sdk.NewInt64Coin("locked/stake", 600)),
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				))
				suite.Require().NoError(err)

				// Perform some undelegations
				_, err = restakingService.UndelegatePool(ctx, restakingtypes.NewMsgUndelegatePool(
					sdk.NewInt64Coin("locked/stake", 100),
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				))
				suite.Require().NoError(err)

				_, err = restakingService.UndelegateService(ctx, restakingtypes.NewMsgUndelegateService(
					1,
					sdk.NewCoins(sdk.NewInt64Coin("locked/stake", 100)),
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				))
				suite.Require().NoError(err)

				_, err = restakingService.UndelegateOperator(ctx, restakingtypes.NewMsgUndelegateOperator(
					1,
					sdk.NewCoins(sdk.NewInt64Coin("locked/stake", 400)),
					"cosmos1pgzph9rze2j2xxavx4n7pdhxlkgsq7raqh8hre",
				))
				suite.Require().NoError(err)
			},
			expected: sdk.NewDecCoins(sdk.NewInt64DecCoin("locked/stake", 1000)),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, _ := suite.ctx.CacheContext()

			if tc.setup != nil {
				tc.setup(ctx)
			}

			coins, err := suite.k.GetAllUserActiveLockedRepresentations(ctx, tc.user)
			suite.Require().NoError(err)
			suite.Require().Equal(tc.expected, coins)
		})
	}
}
