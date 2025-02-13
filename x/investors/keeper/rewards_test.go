package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
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

			querier := distrkeeper.NewQuerier(suite.dk)
			// Query the normal account's rewards
			cacheCtx, _ := ctx.CacheContext()
			rewards, err := querier.DelegationRewards(cacheCtx, &distrtypes.QueryDelegationRewardsRequest{
				DelegatorAddress: normalAddr.String(),
				ValidatorAddress: validator.GetOperator(),
			})
			suite.Require().NoError(err)
			suite.Require().Equal(tc.expNormalAccountRewards, rewards.Rewards.String())

			// Query the investor's rewards
			cacheCtx, _ = ctx.CacheContext()
			rewards, err = querier.DelegationRewards(cacheCtx, &distrtypes.QueryDelegationRewardsRequest{
				DelegatorAddress: investorAddr.String(),
				ValidatorAddress: validator.GetOperator(),
			})
			suite.Require().NoError(err)
			suite.Require().Equal(tc.expInvestorRewards, rewards.Rewards.String())
		})
	}
}

func (suite *KeeperTestSuite) TestUpdateInvestorsRewardRatio() {
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

	querier := distrkeeper.NewQuerier(suite.dk)
	// Query the normal account's rewards
	cacheCtx, _ := ctx.CacheContext()
	rewards, err := querier.DelegationRewards(cacheCtx, &distrtypes.QueryDelegationRewardsRequest{
		DelegatorAddress: investorAddr.String(),
		ValidatorAddress: validator.GetOperator(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal("200000.000000000000000000stake", rewards.Rewards.String())

	// Update the reward ratio to 100%
	err = suite.k.UpdateInvestorsRewardRatio(ctx, sdkmath.LegacyOneDec()) // 100%
	suite.Require().NoError(err)

	// Withdraw the rewards, it shouldn't be changed
	withdrawn, err := suite.dk.WithdrawDelegationRewards(ctx, investorAddr, valAddr)
	suite.Require().NoError(err)
	suite.Require().Equal("200000stake", withdrawn.String())

	// Allocate 1000000stake as rewards again
	ctx = suite.allocateTokensToValidator(ctx, valAddr, utils.MustParseDecCoins("1000000stake"), true)

	// The investor receives more rewards than before
	cacheCtx, _ = ctx.CacheContext()
	rewards, err = querier.DelegationRewards(cacheCtx, &distrtypes.QueryDelegationRewardsRequest{
		DelegatorAddress: investorAddr.String(),
		ValidatorAddress: validator.GetOperator(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal("333333.333333333333000000stake", rewards.Rewards.String())

	// The normal account's rewards are now 400000stake(already given) +
	// 333333.3stake(newly allocated) = 733333.3stake
	cacheCtx, _ = ctx.CacheContext()
	rewards, err = querier.DelegationRewards(cacheCtx, &distrtypes.QueryDelegationRewardsRequest{
		DelegatorAddress: normalAddr.String(),
		ValidatorAddress: validator.GetOperator(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal("733333.333333333333000000stake", rewards.Rewards.String())
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

	querier := distrkeeper.NewQuerier(suite.dk)
	// Query the investor's rewards
	cacheCtx, _ := ctx.CacheContext()
	rewards, err := querier.DelegationRewards(cacheCtx, &distrtypes.QueryDelegationRewardsRequest{
		DelegatorAddress: investorAddr.String(),
		ValidatorAddress: validator.GetOperator(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal("200000.000000000000000000stake", rewards.Rewards.String())

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
	cacheCtx, _ = ctx.CacheContext()
	rewards, err = querier.DelegationRewards(cacheCtx, &distrtypes.QueryDelegationRewardsRequest{
		DelegatorAddress: investorAddr.String(),
		ValidatorAddress: validator.GetOperator(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal("333333.333333333333000000stake", rewards.Rewards.String())

	// Query the normal account's rewards
	// The normal account received ~333333stake = 1000000 / 3 for this block,
	// so the accumulated rewards should be ~733333stake = 400000 + 333333
	cacheCtx, _ = ctx.CacheContext()
	rewards, err = querier.DelegationRewards(cacheCtx, &distrtypes.QueryDelegationRewardsRequest{
		DelegatorAddress: normalAddr.String(),
		ValidatorAddress: validator.GetOperator(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal("733333.333333333333000000stake", rewards.Rewards.String())
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
