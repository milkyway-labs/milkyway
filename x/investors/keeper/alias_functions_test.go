package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/milkyway-labs/milkyway/v10/app/testutil"
	"github.com/milkyway-labs/milkyway/v10/utils"
)

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
	balancesAfter := suite.bk.GetAllBalances(ctx, investorAddr)
	rewards := balancesAfter.Sub(balancesBefore...)
	suite.Assert().Equal("166666stake", rewards.String()) // 333333 * 0.5(used previous ratio)
}

func (suite *KeeperTestSuite) TestSetVestingInvestor() {
	ctx, _ := suite.ctx.CacheContext()
	ctx = ctx.WithBlockTime(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))

	normalAddr := testutil.TestAddress(1)
	suite.fundAccount(ctx, normalAddr.String(), utils.MustParseCoins("1000000stake"))
	activeVestingAddr := testutil.TestAddress(2)
	suite.createVestingAccount(
		ctx,
		testutil.TestAddress(10000).String(),
		activeVestingAddr.String(),
		utils.MustParseCoins("1000000stake"),
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		false,
		true,
	)
	endedVestingAddr := testutil.TestAddress(3)
	suite.createVestingAccount(
		ctx,
		testutil.TestAddress(10000).String(),
		endedVestingAddr.String(),
		utils.MustParseCoins("1000000stake"),
		time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC).Unix(),
		false,
		true,
	)

	// endedVestingAddr's vesting schedule ends
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC))

	testCases := []struct {
		name      string
		addr      string
		shouldErr bool
	}{
		{
			name:      "vesting account returns no error",
			addr:      activeVestingAddr.String(),
			shouldErr: false,
		},
		{
			name:      "ended vesting account returns no error",
			addr:      endedVestingAddr.String(),
			shouldErr: false,
		},
		{
			name:      "invalid address returns error",
			addr:      "invalid",
			shouldErr: true,
		},
		{
			name:      "non existent address returns error",
			addr:      testutil.TestAddress(4).String(),
			shouldErr: true,
		},
		{
			name:      "non-vesting account returns error",
			addr:      normalAddr.String(),
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			err := suite.k.SetVestingInvestor(ctx, tc.addr)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
