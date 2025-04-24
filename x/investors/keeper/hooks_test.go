package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/milkyway-labs/milkyway/v11/app/testutil"
	"github.com/milkyway-labs/milkyway/v11/utils"
)

func (suite *KeeperTestSuite) TestBeforeDelegationRewardsWithdrawnHook_WithdrawAddressHasNoEffect() {
	ctx := suite.ctx.WithBlockTime(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))

	valOwnerAddr := testutil.TestAddress(10000)
	validator := suite.createValidator(
		ctx,
		valOwnerAddr,
		stakingtypes.NewCommissionRates(utils.MustParseDec("0"), utils.MustParseDec("0.2"), utils.MustParseDec("0.01")),
		utils.MustParseCoin("1000000stake"),
		true,
	)
	valAddr := sdk.ValAddress(valOwnerAddr)

	investorAddr := testutil.TestAddress(1)
	normalAddr := testutil.TestAddress(2)
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
	err := suite.k.SetVestingInvestor(ctx, investorAddr.String())
	suite.Require().NoError(err)
	suite.delegate(ctx, investorAddr.String(), validator.GetOperator(), utils.MustParseCoin("1000000stake"), false)
	// The investor sets its withdraw address to a normal account, trying to bypass
	// the rewards reduction logic
	err = suite.dk.SetWithdrawAddr(ctx, investorAddr, normalAddr)
	suite.Require().NoError(err)

	ctx = suite.allocateTokensToValidator(ctx, valAddr, utils.MustParseDecCoins("1000000stake"), true)

	// Even if the withdraw address is a normal account, the hook checks if the
	// original delegator is a vesting investor so it will still lower the rewards of
	// the investor
	balancesBefore := suite.bk.GetAllBalances(ctx, normalAddr)
	rewards, err := suite.dk.WithdrawDelegationRewards(ctx, investorAddr, valAddr)
	suite.Require().NoError(err)
	suite.Require().Equal("250000stake", rewards.String())
	balancesAfter := suite.bk.GetAllBalances(ctx, normalAddr)
	withdrawnRewards := balancesAfter.Sub(balancesBefore...)
	suite.Require().Equal("250000stake", withdrawnRewards.String())
}
