package keeper_test

import (
	"time"

	"github.com/milkyway-labs/milkyway/v7/app/testutil"
	"github.com/milkyway-labs/milkyway/v7/utils"
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
	err := suite.k.TrySetVestingInvestor(ctx, investorAddr1)
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
	err = suite.k.TrySetVestingInvestor(ctx, investorAddr2)
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
	err = suite.k.TrySetVestingInvestor(ctx, investorAddr3)
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
