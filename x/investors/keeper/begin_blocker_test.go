package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/milkyway-labs/milkyway/v9/app/testutil"
	"github.com/milkyway-labs/milkyway/v9/utils"
)

func (suite *KeeperTestSuite) TestRemoveVestingEndedInvestors() {
	ctx, _ := suite.ctx.CacheContext()
	ctx = ctx.WithBlockTime(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))

	faucet := testutil.TestAddress(10000)

	vestingEndTime1 := ctx.BlockTime().Add(365 * 24 * time.Hour) // = T
	vestingEndTime2 := vestingEndTime1.Add(time.Second)          // = T'

	// Vesting ends at T
	investorAddr1 := testutil.TestAddress(1)
	suite.createVestingAccount(
		ctx,
		faucet.String(),
		investorAddr1.String(),
		utils.MustParseCoins("1000000stake"),
		vestingEndTime1.Unix(),
		false,
		true,
	)
	err := suite.k.SetVestingInvestor(ctx, investorAddr1)
	suite.Require().NoError(err)

	// Vesting ends at T
	investorAddr2 := testutil.TestAddress(2)
	suite.createVestingAccount(
		ctx,
		faucet.String(),
		investorAddr2.String(),
		utils.MustParseCoins("1000000stake"),
		vestingEndTime1.Unix(),
		false,
		true,
	)
	err = suite.k.SetVestingInvestor(ctx, investorAddr2)
	suite.Require().NoError(err)

	// Vesting ends at T'
	investorAddr3 := testutil.TestAddress(3)
	suite.createVestingAccount(
		ctx,
		faucet.String(),
		investorAddr3.String(),
		utils.MustParseCoins("1000000stake"),
		vestingEndTime2.Unix(),
		false,
		true,
	)
	err = suite.k.SetVestingInvestor(ctx, investorAddr3)
	suite.Require().NoError(err)

	// All investors are still in vesting period so nothing happens
	err = suite.k.RemoveVestingEndedInvestors(ctx)
	suite.Require().NoError(err)
	suite.Require().True(suite.isVestingInvestor(ctx, investorAddr1))
	suite.Require().True(suite.isVestingInvestor(ctx, investorAddr2))
	suite.Require().True(suite.isVestingInvestor(ctx, investorAddr3))

	// Move to T, both investor 1 and 2 should be removed from the list
	ctx = ctx.WithBlockTime(vestingEndTime1)
	err = suite.k.RemoveVestingEndedInvestors(ctx)
	suite.Require().NoError(err)
	suite.Require().False(suite.isVestingInvestor(ctx, investorAddr1))
	suite.Require().False(suite.isVestingInvestor(ctx, investorAddr2))
	suite.Require().True(suite.isVestingInvestor(ctx, investorAddr3))

	// Move to T', now the investor 3 should be removed from the list, too
	ctx = ctx.WithBlockTime(vestingEndTime2)
	err = suite.k.RemoveVestingEndedInvestors(ctx)
	suite.Require().NoError(err)
	suite.Require().False(suite.isVestingInvestor(ctx, investorAddr1))
	suite.Require().False(suite.isVestingInvestor(ctx, investorAddr2))
	suite.Require().False(suite.isVestingInvestor(ctx, investorAddr3))
}

func (suite *KeeperTestSuite) TestRemoveVestingInvestor_RewardsAfterVestingEnded() {
	ctx, _ := suite.ctx.CacheContext()
	ctx = ctx.WithBlockTime(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))

	err := suite.k.UpdateInvestorsRewardRatio(ctx, sdkmath.LegacyNewDecWithPrec(5, 1)) // 50%
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
	err = suite.k.SetVestingInvestor(ctx, investorAddr)
	suite.Require().NoError(err)

	suite.delegate(ctx, normalAddr.String(), validator.GetOperator(), utils.MustParseCoin("1000000stake"), true)
	suite.delegate(ctx, investorAddr.String(), validator.GetOperator(), utils.MustParseCoin("1000000stake"), false)

	ctx = suite.allocateTokensToValidator(ctx, valAddr, utils.MustParseDecCoins("1000000stake"), true)

	// Query the investor's rewards
	rewards := suite.delegationRewards(ctx, investorAddr.String(), validator.GetOperator())
	suite.Require().Equal("200000.000000000000000000stake", rewards.String())

	// Remove the investor from the vesting investors list
	err = suite.k.RemoveVestingInvestor(ctx, investorAddr)
	suite.Require().NoError(err)

	// Rewards are withdrawn
	balances := suite.bk.GetAllBalances(ctx, investorAddr)
	suite.Require().Equal("200000stake", balances.String())

	// Query the normal account's rewards, it shouldn't be changed
	rewards = suite.delegationRewards(ctx, normalAddr.String(), validator.GetOperator())
	suite.Require().Equal("400000.000000000000000000stake", rewards.String())

	ctx = suite.allocateTokensToValidator(ctx, valAddr, utils.MustParseDecCoins("1000000stake"), true)

	// Query the investor's rewards
	rewards = suite.delegationRewards(ctx, investorAddr.String(), validator.GetOperator())
	suite.Require().Equal("333333.333333333333000000stake", rewards.String())
}
