package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/milkyway-labs/milkyway/v9/app/testutil"
	"github.com/milkyway-labs/milkyway/v9/utils"
)

func (suite *KeeperTestSuite) TestInvestorsRewardRatio() {
	testCases := []struct {
		name                    string
		investorsRewardRatio    sdkmath.LegacyDec
		expNormalAccountRewards string
		expInvestorRewards      string
	}{
		{
			name:                    "0% investors reward ratio",
			investorsRewardRatio:    sdkmath.LegacyZeroDec(),
			expNormalAccountRewards: "500000.000000000000000000stake",
			expInvestorRewards:      "",
		},
		{
			name:                    "10% investors reward ratio",
			investorsRewardRatio:    utils.MustParseDec("0.1"),
			expNormalAccountRewards: "476190.476190476190000000stake", // 1000000 * 1 / 2.1
			expInvestorRewards:      "47619.047619047619000000stake",  // 1000000 * 0.1 / 2.1
		},
		{
			name:                    "50% investors reward ratio",
			investorsRewardRatio:    utils.MustParseDec("0.5"),
			expNormalAccountRewards: "400000.000000000000000000stake",
			expInvestorRewards:      "200000.000000000000000000stake",
		},
		{
			name:                    "100% investors reward ratio",
			investorsRewardRatio:    sdkmath.LegacyOneDec(),
			expNormalAccountRewards: "333333.333333333333000000stake",
			expInvestorRewards:      "333333.333333333333000000stake",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx := suite.ctx.WithBlockTime(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))

			err := suite.k.UpdateInvestorsRewardRatio(ctx, tc.investorsRewardRatio)
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
			err = suite.k.SetVestingInvestor(ctx, investorAddr)
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

			rewards := suite.delegationRewards(ctx, normalAddr.String(), validator.GetOperator())
			suite.Require().Equal(tc.expNormalAccountRewards, rewards.String())

			// Query the investor's rewards
			rewards = suite.delegationRewards(ctx, investorAddr.String(), validator.GetOperator())
			suite.Require().Equal(tc.expInvestorRewards, rewards.String())
		})
	}
}

func (suite *KeeperTestSuite) TestUpdateInvestorsRewardRatio() {
	ctx, _ := suite.ctx.CacheContext()
	ctx = ctx.WithBlockTime(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))

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
	err := suite.k.SetVestingInvestor(ctx, investorAddr)
	suite.Require().NoError(err)

	testCases := []struct {
		name         string
		initialRatio sdkmath.LegacyDec
		newRatio     sdkmath.LegacyDec
		check        func(ctx sdk.Context)
	}{
		{
			name:         "increase reward ratio",
			initialRatio: utils.MustParseDec("0.5"),
			newRatio:     utils.MustParseDec("1"),
			check: func(ctx sdk.Context) {
				// The investor receives more rewards than before
				rewards := suite.delegationRewards(ctx, investorAddr.String(), validator.GetOperator())
				suite.Assert().Equal("333333.333333333333000000stake", rewards.String())

				// The normal account's rewards are now 400000stake(already given) +
				// 333333.3stake(newly allocated) = 733333.3stake
				rewards = suite.delegationRewards(ctx, normalAddr.String(), validator.GetOperator())
				suite.Assert().Equal("733333.333333333333000000stake", rewards.String())
			},
		},
		{
			name:         "decrease reward ratio",
			initialRatio: utils.MustParseDec("1"),
			newRatio:     utils.MustParseDec("0.5"),
			check: func(ctx sdk.Context) {
				// The investor receives fewer rewards than before
				rewards := suite.delegationRewards(ctx, investorAddr.String(), validator.GetOperator())
				suite.Assert().Equal("200000.000000000000000000stake", rewards.String())

				// The normal account's rewards are now 333333.3stake(already given) +
				// 400000stake(newly allocated) = 733333.3stake
				rewards = suite.delegationRewards(ctx, normalAddr.String(), validator.GetOperator())
				suite.Assert().Equal("733333.333333333333000000stake", rewards.String())
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			ctx, _ := ctx.CacheContext()

			err := suite.k.UpdateInvestorsRewardRatio(ctx, tc.initialRatio)
			suite.Require().NoError(err)

			suite.delegate(ctx, normalAddr.String(), validator.GetOperator(), utils.MustParseCoin("1000000stake"), true)
			suite.delegate(ctx, investorAddr.String(), validator.GetOperator(), utils.MustParseCoin("1000000stake"), false)

			// Allocate 1000000stake as rewards
			ctx = suite.allocateTokensToValidator(ctx, valAddr, utils.MustParseDecCoins("1000000stake"), true)

			// Update the reward ratio
			err = suite.k.UpdateInvestorsRewardRatio(ctx, tc.newRatio)
			suite.Require().NoError(err)

			// Allocate 1000000stake as rewards again
			ctx = suite.allocateTokensToValidator(ctx, valAddr, utils.MustParseDecCoins("1000000stake"), true)

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestVestingEndedInvestorsReward() {
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

	// Allocate 1000000stake as rewards
	ctx = suite.allocateTokensToValidator(ctx, valAddr, utils.MustParseDecCoins("1000000stake"), true)

	// Query the investor's rewards
	rewards := suite.delegationRewards(ctx, investorAddr.String(), validator.GetOperator())
	suite.Require().Equal("200000.000000000000000000stake", rewards.String())

	// Now the investor's vesting period is over, the investor should receive normal
	// rewards
	ctx = ctx.WithBlockTime(vestingEndTime)
	err = suite.k.RemoveVestingEndedInvestors(ctx)
	suite.Require().NoError(err)
	suite.Require().False(suite.isVestingInvestor(ctx, investorAddr))

	// Allocate 1000000stake as rewards
	ctx = suite.allocateTokensToValidator(ctx, valAddr, utils.MustParseDecCoins("1000000stake"), true)

	// Query the investor's rewards
	// The investor should receive ~333333stake = 1000000 / 3
	rewards = suite.delegationRewards(ctx, investorAddr.String(), validator.GetOperator())
	suite.Require().Equal("333333.333333333333000000stake", rewards.String())

	// Query the normal account's rewards
	// The normal account received ~333333stake = 1000000 / 3 for this block,
	// so the accumulated rewards should be ~733333stake = 400000 + 333333
	rewards = suite.delegationRewards(ctx, normalAddr.String(), validator.GetOperator())
	suite.Require().Equal("733333.333333333333000000stake", rewards.String())
}

func (suite *KeeperTestSuite) TestUnbond() {
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

	investorAddr := testutil.TestAddress(1)
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

	suite.delegate(ctx, investorAddr.String(), validator.GetOperator(), utils.MustParseCoin("1000000stake"), false)

	// Unbonding should work with vesting investors
	_, err = stakingkeeper.NewMsgServerImpl(suite.sk).Undelegate(
		ctx,
		stakingtypes.NewMsgUndelegate(investorAddr.String(), valAddr.String(), utils.MustParseCoin("1000000stake")),
	)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) TestRedelegate() {
	ctx, _ := suite.ctx.CacheContext()
	ctx = ctx.WithBlockTime(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))

	err := suite.k.UpdateInvestorsRewardRatio(ctx, sdkmath.LegacyNewDecWithPrec(5, 1)) // 50%
	suite.Require().NoError(err)

	valOwnerAddr1 := testutil.TestAddress(10000)
	validator1 := suite.createValidator(
		ctx,
		valOwnerAddr1,
		stakingtypes.NewCommissionRates(utils.MustParseDec("0"), utils.MustParseDec("0.2"), utils.MustParseDec("0.01")),
		utils.MustParseCoin("1000000stake"),
		true,
	)
	valAddr1 := sdk.ValAddress(valOwnerAddr1)

	valOwnerAddr2 := testutil.TestAddress(10001)
	validator2 := suite.createValidator(
		ctx,
		valOwnerAddr2,
		stakingtypes.NewCommissionRates(utils.MustParseDec("0"), utils.MustParseDec("0.2"), utils.MustParseDec("0.01")),
		utils.MustParseCoin("1000000stake"),
		true,
	)
	valAddr2 := sdk.ValAddress(valOwnerAddr2)

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

	// The normal user delegates 1000000stake to both validator 1 and 2
	suite.delegate(ctx, normalAddr.String(), validator1.GetOperator(), utils.MustParseCoin("1000000stake"), true)
	suite.delegate(ctx, normalAddr.String(), validator2.GetOperator(), utils.MustParseCoin("1000000stake"), true)
	// The investor delegates 1000000stake to validator 1
	suite.delegate(ctx, investorAddr.String(), validator1.GetOperator(), utils.MustParseCoin("1000000stake"), false)

	// Allocate 1000000stake as rewards to both validators
	ctx = suite.allocateTokensToValidator(ctx, valAddr1, utils.MustParseDecCoins("1000000stake"), true)
	ctx = suite.allocateTokensToValidator(ctx, valAddr2, utils.MustParseDecCoins("1000000stake"), true)

	// Redelegate the investor's delegation from validator 1 to validator 2
	_, err = stakingkeeper.NewMsgServerImpl(suite.sk).BeginRedelegate(
		ctx,
		stakingtypes.NewMsgBeginRedelegate(
			investorAddr.String(),
			validator1.GetOperator(),
			validator2.GetOperator(),
			utils.MustParseCoin("1000000stake"),
		),
	)
	suite.Require().NoError(err)
	// Previous rewards are withdrawn due to the redelegation
	suite.Assert().Equal("200000stake", suite.bk.GetAllBalances(ctx, investorAddr).String())

	// Allocate 1000000stake as rewards to both validators again
	ctx = suite.allocateTokensToValidator(ctx, valAddr1, utils.MustParseDecCoins("1000000stake"), true)
	ctx = suite.allocateTokensToValidator(ctx, valAddr2, utils.MustParseDecCoins("1000000stake"), true)

	// No rewards from the validator 1 since the investor has redelegated all
	// tokens from it
	rewards := suite.delegationRewards(ctx, investorAddr.String(), validator1.GetOperator())
	suite.Assert().Equal("", rewards.String())
	// The investor receives 2000000stake from the validator 2 after the
	// redelegation
	rewards = suite.delegationRewards(ctx, investorAddr.String(), validator2.GetOperator())
	suite.Assert().Equal("200000.000000000000000000stake", rewards.String())

	// 4000000stake(before redelegation) + 5000000stake(after redelegation)
	rewards = suite.delegationRewards(ctx, normalAddr.String(), validator1.GetOperator())
	suite.Assert().Equal("900000.000000000000000000stake", rewards.String())
	// 5000000stake(before redelegation) + 4000000stake(after redelegation)
	rewards = suite.delegationRewards(ctx, normalAddr.String(), validator2.GetOperator())
	suite.Assert().Equal("900000.000000000000000000stake", rewards.String())
}
