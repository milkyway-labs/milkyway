package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	restakingkeeper "github.com/milkyway-labs/milkyway/x/restaking/keeper"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
)

func (suite *KeeperTestSuite) TestKeeper_AddToInsuranceFund() {
	testCases := []struct {
		name                string
		deposits            map[string]sdk.Coins
		expectedTotalAmount sdk.Coins
	}{
		{
			name: "add multiple amounts",
			deposits: map[string]sdk.Coins{
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn": sdk.NewCoins(sdk.NewInt64Coin("stake", 100)),
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd": sdk.NewCoins(sdk.NewInt64Coin("stake", 200)),
			},
			expectedTotalAmount: sdk.NewCoins(sdk.NewInt64Coin("stake", 300)),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			for address, amount := range tc.deposits {
				// Mint the coins that should be in the module
				suite.Assert().NoError(
					suite.bk.MintCoins(suite.ctx, types.ModuleName, amount))
				accAddress := sdk.MustAccAddressFromBech32(address)
				suite.Assert().NoError(
					suite.k.AddToUserInsuranceFund(suite.ctx, accAddress, amount))
			}

			for address, expectedAmount := range tc.deposits {
				amount, err := suite.k.GetUserInsuranceFundBalance(suite.ctx, address)
				suite.Assert().NoError(err)
				suite.Assert().Equal(expectedAmount, amount)
			}

			balance, err := suite.k.GetInsuranceFundBalance(suite.ctx)
			suite.Assert().NoError(err)
			suite.Assert().Equal(tc.expectedTotalAmount, balance)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_WithdrawFromInsuranceFund() {
	testCases := []struct {
		name      string
		from      string
		amount    sdk.Coins
		setup     func()
		shouldErr bool
	}{
		{
			name: "withdraw more then deposited",
			from: "cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
			amount: sdk.NewCoins(
				sdk.NewInt64Coin("stake", 100),
				sdk.NewInt64Coin("stake2", 50),
			),
			setup: func() {
				suite.fundAccountInsuranceFund(suite.ctx,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewCoins(
						sdk.NewInt64Coin("stake", 50),
						sdk.NewInt64Coin("stake2", 50),
					))
			},
			shouldErr: true,
		},
		{
			name: "withdraw correctly",
			from: "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			amount: sdk.NewCoins(
				sdk.NewInt64Coin("stake", 200),
				sdk.NewInt64Coin("stake2", 100),
			),
			setup: func() {
				suite.fundAccountInsuranceFund(suite.ctx,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewCoins(
						sdk.NewInt64Coin("stake", 200),
						sdk.NewInt64Coin("stake2", 100),
					))
			},
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.setup()
			accAddr := sdk.MustAccAddressFromBech32(tc.from)
			err := suite.k.WithdrawFromUserInsuranceFund(suite.ctx, accAddr, tc.amount)
			if tc.shouldErr {
				suite.Assert().Error(err)
			} else {
				suite.Assert().NoError(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_UsedInsuranceFundIsUpdatedCorrectly() {
	suite.SetupTest()
	ctx := suite.ctx
	restaker := "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd"

	// Set insurance percentage to 2%
	suite.k.SetParams(ctx, types.NewParams(math.LegacyNewDec(2), nil, nil, nil))

	// Fund the restaker insurance fund
	suite.fundAccountInsuranceFund(ctx, restaker, sdk.NewCoins(sdk.NewInt64Coin("stake", 20)))

	// Mint some vested representation to the restaker
	suite.mintVestedRepresentation(restaker, sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)))

	// Check that the used insurance fund is 0
	used, err := suite.k.GetUserUsedInsuranceFund(ctx, restaker)
	suite.Require().NoError(err)
	suite.Require().True(used.IsZero())

	// Create a pool
	pool := poolstypes.NewPool(1, "vested/stake")
	err = suite.pk.SavePool(ctx, pool)
	suite.Require().NoError(err)

	// Restake some coins to the pool
	restakingMsgService := restakingkeeper.NewMsgServer(suite.rk)
	_, err = restakingMsgService.DelegatePool(ctx, restakingtypes.NewMsgDelegatePool(sdk.NewInt64Coin("vested/stake", 100), restaker))
	suite.Require().NoError(err)

	// Check that the used insurance fund is 2stake
	used, err = suite.k.GetUserUsedInsuranceFund(ctx, restaker)
	suite.Require().NoError(err)
	suite.Require().Equal(sdk.NewCoins(sdk.NewInt64Coin("stake", 2)), used)

	// Undelegate 1vested/stake
	_, err = restakingMsgService.UndelegatePool(ctx, restakingtypes.NewMsgUndelegatePool(sdk.NewInt64Coin("vested/stake", 1), restaker))
	suite.Require().NoError(err)

	// Check that the used insurance fund is still 2stake
	used, err = suite.k.GetUserUsedInsuranceFund(ctx, restaker)
	suite.Require().NoError(err)
	suite.Require().Equal(sdk.NewCoins(sdk.NewInt64Coin("stake", 2)), used)

	// Wait for the unbonding period to expire
	restakingParams := suite.rk.GetParams(ctx)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 100).
		WithBlockTime(ctx.BlockTime().Add(restakingParams.UnbondingTime))
	// Trigger the unbonding
	suite.rk.CompleteMatureUnbondingDelegations(ctx)

	// Ensure the user has no pending undelegations
	unbondingDelegations := suite.rk.GetAllUserUnbondingDelegations(ctx, restaker)
	suite.Require().Empty(unbondingDelegations)

	// Check that the user used insurance fund is still 2stake
	used, err = suite.k.GetUserUsedInsuranceFund(ctx, restaker)
	suite.Require().NoError(err)
	suite.Require().Equal(sdk.NewCoins(sdk.NewInt64Coin("stake", 2)), used)
}

func (suite *KeeperTestSuite) TestKeeper_InsuranceFundUpdatesCorreclyWithCompleteUnbond() {
	suite.SetupTest()
	ctx := suite.ctx
	restaker := "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd"

	// Set insurance percentage to 2%
	suite.k.SetParams(ctx, types.NewParams(math.LegacyNewDec(2), nil, nil, nil))

	// Fund the restaker insurance fund
	suite.fundAccountInsuranceFund(ctx, restaker, sdk.NewCoins(sdk.NewInt64Coin("stake", 20)))

	// Mint some vested representation to the restaker
	suite.mintVestedRepresentation(restaker, sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)))

	// Check that the used insurance fund is 0
	used, err := suite.k.GetUserUsedInsuranceFund(ctx, restaker)
	suite.Require().NoError(err)
	suite.Require().True(used.IsZero())

	// Create a pool
	pool := poolstypes.NewPool(1, "vested/stake")
	err = suite.pk.SavePool(ctx, pool)
	suite.Require().NoError(err)

	// Restake some coins to the pool
	restakingMsgService := restakingkeeper.NewMsgServer(suite.rk)
	_, err = restakingMsgService.DelegatePool(ctx, restakingtypes.NewMsgDelegatePool(sdk.NewInt64Coin("vested/stake", 1000), restaker))
	suite.Require().NoError(err)

	// Check that the used insurance fund is 20stake
	used, err = suite.k.GetUserUsedInsuranceFund(ctx, restaker)
	suite.Require().NoError(err)
	suite.Require().Equal(sdk.NewCoins(sdk.NewInt64Coin("stake", 20)), used)

	// Undelegate all the vested representations
	_, err = restakingMsgService.UndelegatePool(ctx, restakingtypes.NewMsgUndelegatePool(sdk.NewInt64Coin("vested/stake", 1000), restaker))
	suite.Require().NoError(err)

	// Check that the used insurance fund is still 20stake
	used, err = suite.k.GetUserUsedInsuranceFund(ctx, restaker)
	suite.Require().NoError(err)
	suite.Require().Equal(sdk.NewCoins(sdk.NewInt64Coin("stake", 20)), used)

	// Wait for the unbonding period to expire
	restakingParams := suite.rk.GetParams(ctx)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 100).
		WithBlockTime(ctx.BlockTime().Add(restakingParams.UnbondingTime))
	// Trigger the unbonding
	suite.rk.CompleteMatureUnbondingDelegations(ctx)

	// Ensure the user has no pending undelegations
	unbondingDelegations := suite.rk.GetAllUserUnbondingDelegations(ctx, restaker)
	suite.Require().Empty(unbondingDelegations)

	// Check that the user used insurance fund is zero
	used, err = suite.k.GetUserUsedInsuranceFund(ctx, restaker)
	suite.Require().NoError(err)
	suite.Require().True(used.IsZero())
}
