package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/milkyway-labs/milkyway/v10/app/testutil"
	"github.com/milkyway-labs/milkyway/v10/utils"
	distrkeeper "github.com/milkyway-labs/milkyway/v10/x/distribution/keeper"
	"github.com/milkyway-labs/milkyway/v10/x/investors"
	"github.com/milkyway-labs/milkyway/v10/x/investors/types"
)

func (suite *KeeperTestSuite) TestInvestorsRewardRatio() {
	testCases := []struct {
		name                 string
		investorsRewardRatio sdkmath.LegacyDec
		expectedRewards      string
	}{
		{
			name:                 "0% investors reward ratio",
			investorsRewardRatio: sdkmath.LegacyZeroDec(),
			expectedRewards:      "",
		},
		{
			name:                 "10% investors reward ratio",
			investorsRewardRatio: utils.MustParseDec("0.1"),
			expectedRewards:      "33333stake", // 333333 * 0.1
		},
		{
			name:                 "50% investors reward ratio",
			investorsRewardRatio: utils.MustParseDec("0.5"),
			expectedRewards:      "166666stake", // 333333 * 0.5
		},
		{
			name:                 "100% investors reward ratio",
			investorsRewardRatio: sdkmath.LegacyOneDec(),
			expectedRewards:      "333333stake",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx := suite.ctx.WithBlockTime(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))

			err := suite.k.SetInvestorsRewardRatio(ctx, tc.investorsRewardRatio)
			suite.Require().NoError(err)

			normalAddr := testutil.TestAddress(1)
			investorAddr := testutil.TestAddress(2)
			vestingEndTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
			suite.createVestingAccount(
				ctx,
				testutil.TestAddress(10001).String(),
				investorAddr.String(),
				utils.MustParseCoins("1000000stake"),
				vestingEndTime.Unix(),
				false,
				true,
			)
			err = suite.k.SetVestingInvestor(ctx, investorAddr.String())
			suite.Require().NoError(err)

			valOwnerAddr := testutil.TestAddress(10000)
			validator := suite.createValidator(
				ctx,
				valOwnerAddr,
				stakingtypes.NewCommissionRates(utils.MustParseDec("0"), utils.MustParseDec("0.2"), utils.MustParseDec("0.01")),
				utils.MustParseCoin("1000000stake"),
				true,
			)
			valAddr := sdk.ValAddress(valOwnerAddr)
			suite.delegate(ctx, normalAddr.String(), validator.GetOperator(), utils.MustParseCoin("1000000stake"), true)
			suite.delegate(ctx, investorAddr.String(), validator.GetOperator(), utils.MustParseCoin("1000000stake"), false)

			ctx = suite.allocateTokensToValidator(ctx, valAddr, utils.MustParseDecCoins("1000000stake"), true)

			// Check the investor's rewards
			var rewards sdk.Coins
			rewards, ctx = suite.withdrawRewardsAndIncrementBlockHeight(ctx, investorAddr, valAddr)
			suite.Require().Equal(tc.expectedRewards, rewards.String())
		})
	}
}

func (suite *KeeperTestSuite) TestCommunityPool() {
	ctx, _ := suite.ctx.CacheContext()
	ctx = suite.ctx.WithBlockTime(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))

	err := suite.k.SetInvestorsRewardRatio(ctx, utils.MustParseDec("0.5")) // 50%
	suite.Require().NoError(err)

	normalAddr := testutil.TestAddress(1)
	investorAddr := testutil.TestAddress(2)
	vestingEndTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	suite.createVestingAccount(
		ctx,
		testutil.TestAddress(10001).String(),
		investorAddr.String(),
		utils.MustParseCoins("1000000stake"),
		vestingEndTime.Unix(),
		false,
		true,
	)
	err = suite.k.SetVestingInvestor(ctx, investorAddr.String())
	suite.Require().NoError(err)

	valOwnerAddr := testutil.TestAddress(10000)
	validator := suite.createValidator(
		ctx,
		valOwnerAddr,
		stakingtypes.NewCommissionRates(utils.MustParseDec("0"), utils.MustParseDec("0.2"), utils.MustParseDec("0.01")),
		utils.MustParseCoin("1000000stake"),
		true,
	)
	valAddr := sdk.ValAddress(valOwnerAddr)

	suite.delegate(ctx, normalAddr.String(), validator.GetOperator(), utils.MustParseCoin("1000000stake"), true)
	suite.delegate(ctx, investorAddr.String(), validator.GetOperator(), utils.MustParseCoin("1000000stake"), false)

	ctx = suite.allocateTokensToValidator(ctx, valAddr, utils.MustParseDecCoins("1000000stake"), true)

	// Check the investor's rewards
	var rewards sdk.Coins
	rewards, ctx = suite.withdrawRewardsAndIncrementBlockHeight(ctx, investorAddr, valAddr)
	suite.Assert().Equal("166666stake", rewards.String()) // 333333 * 0.5

	// Check the community pool
	res, err := distrkeeper.NewQuerier(suite.dk).CommunityPool(ctx, &distrtypes.QueryCommunityPoolRequest{})
	suite.Require().NoError(err)
	suite.Assert().Equal("166667.333333333333000000stake", res.Pool.String()) // 333333 - (333333 * 0.5) + dust

	// The module should have no balances after the end blocker
	moduleBalances := suite.bk.GetAllBalances(ctx, suite.ak.GetModuleAddress(types.ModuleName))
	suite.Assert().Empty(moduleBalances)
}

func (suite *KeeperTestSuite) TestUpdateInvestorsRewardRatio() {
	ctx, _ := suite.ctx.CacheContext()
	ctx = ctx.WithBlockTime(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))

	err := suite.k.SetInvestorsRewardRatio(ctx, utils.MustParseDec("0.5")) // 50%
	suite.Require().NoError(err)

	valOwnerAddr := testutil.TestAddress(10000)
	validator := suite.createValidator(
		ctx,
		valOwnerAddr,
		stakingtypes.NewCommissionRates(utils.MustParseDec("0"), utils.MustParseDec("0.2"), utils.MustParseDec("0.01")),
		utils.MustParseCoin("1000000stake"),
		true,
	)
	valAddr := sdk.ValAddress(valOwnerAddr)

	normalAddr := testutil.TestAddress(1)
	investorAddr := testutil.TestAddress(2)
	vestingEndTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	suite.createVestingAccount(
		ctx,
		testutil.TestAddress(10001).String(),
		investorAddr.String(),
		utils.MustParseCoins("1000000stake"),
		vestingEndTime.Unix(),
		false,
		true,
	)
	err = suite.k.SetVestingInvestor(ctx, investorAddr.String())
	suite.Require().NoError(err)

	suite.delegate(ctx, normalAddr.String(), validator.GetOperator(), utils.MustParseCoin("1000000stake"), true)
	suite.delegate(ctx, investorAddr.String(), validator.GetOperator(), utils.MustParseCoin("1000000stake"), false)

	ctx = suite.allocateTokensToValidator(ctx, valAddr, utils.MustParseDecCoins("1000000stake"), true)

	// Investor's rewards are automatically withdrawn upon the ratio update
	balancesBefore := suite.bk.GetAllBalances(ctx, investorAddr)
	err = suite.k.UpdateInvestorsRewardRatio(ctx, utils.MustParseDec("1")) // 100%
	suite.Require().NoError(err)
	err = investors.EndBlocker(ctx, suite.k)
	suite.Require().NoError(err)
	balancesAfter := suite.bk.GetAllBalances(ctx, investorAddr)
	rewards := balancesAfter.Sub(balancesBefore...)
	suite.Assert().Equal("166666stake", rewards.String()) // 333333 * 0.5(used previous ratio)
}

func (suite *KeeperTestSuite) TestVestingEndedInvestorsReward() {
	ctx, _ := suite.ctx.CacheContext()
	ctx = ctx.WithBlockTime(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))

	err := suite.k.SetInvestorsRewardRatio(ctx, utils.MustParseDec("0.5")) // 50%
	suite.Require().NoError(err)

	valOwnerAddr := testutil.TestAddress(10000)
	validator := suite.createValidator(
		ctx,
		valOwnerAddr,
		stakingtypes.NewCommissionRates(utils.MustParseDec("0"), utils.MustParseDec("0.2"), utils.MustParseDec("0.01")),
		utils.MustParseCoin("1000000stake"),
		true,
	)
	valAddr := sdk.ValAddress(valOwnerAddr)

	normalAddr := testutil.TestAddress(1)
	investorAddr := testutil.TestAddress(2)
	vestingEndTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	suite.createVestingAccount(
		ctx,
		testutil.TestAddress(10001).String(),
		investorAddr.String(),
		utils.MustParseCoins("1000000stake"),
		vestingEndTime.Unix(),
		false,
		true,
	)
	err = suite.k.SetVestingInvestor(ctx, investorAddr.String())
	suite.Require().NoError(err)

	suite.delegate(ctx, normalAddr.String(), validator.GetOperator(), utils.MustParseCoin("1000000stake"), true)
	suite.delegate(ctx, investorAddr.String(), validator.GetOperator(), utils.MustParseCoin("1000000stake"), false)

	// Allocate 1000000stake as rewards
	ctx = suite.allocateTokensToValidator(ctx, valAddr, utils.MustParseDecCoins("1000000stake"), true)

	// Now the investor's vesting period is over, the investor should receive normal
	// rewards
	balancesBefore := suite.bk.GetAllBalances(ctx, investorAddr)
	ctx = ctx.WithBlockTime(vestingEndTime)
	err = suite.k.RemoveVestingEndedInvestors(ctx)
	suite.Require().NoError(err)
	err = investors.EndBlocker(ctx, suite.k)
	suite.Require().NoError(err)
	balancesAfter := suite.bk.GetAllBalances(ctx, investorAddr)
	rewards := balancesAfter.Sub(balancesBefore...)
	suite.Assert().Equal("166666stake", rewards.String())
	isVestingInvestor, err := suite.k.IsVestingInvestor(ctx, investorAddr.String())
	suite.Require().NoError(err)
	suite.Assert().False(isVestingInvestor)

	// Allocate 1000000stake as rewards
	ctx = suite.allocateTokensToValidator(ctx, valAddr, utils.MustParseDecCoins("1000000stake"), true)

	// Check the investor's rewards
	// The investor should receive ~333333stake = 1000000 / 3
	rewards, ctx = suite.withdrawRewardsAndIncrementBlockHeight(ctx, investorAddr, valAddr)
	suite.Assert().Equal("333333stake", rewards.String())
}
