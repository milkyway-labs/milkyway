package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/milkyway-labs/milkyway/v7/app/testutil"
	"github.com/milkyway-labs/milkyway/v7/utils"
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

			valOwnerAddr := testutil.TestAddress(10000)
			validator := suite.createValidator(
				ctx,
				valOwnerAddr,
				stakingtypes.NewCommissionRates(utils.MustParseDec("0"), utils.MustParseDec("0.2"), utils.MustParseDec("0.01")),
				utils.MustParseCoin("1000000stake"),
				true,
			)
			suite.delegate(ctx, normalAddr.String(), validator.GetOperator(), utils.MustParseCoin("1000000stake"), true)
			suite.delegate(ctx, investorAddr.String(), validator.GetOperator(), utils.MustParseCoin("1000000stake"), false)

			ctx = suite.allocateTokensToValidator(ctx, validator, utils.MustParseDecCoins("1000000stake"), true)

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

	suite.delegate(ctx, normalAddr.String(), validator.GetOperator(), utils.MustParseCoin("1000000stake"), true)
	suite.delegate(ctx, investorAddr.String(), validator.GetOperator(), utils.MustParseCoin("1000000stake"), false)

	// Allocate 1000000stake as rewards
	ctx = suite.allocateTokensToValidator(ctx, validator, utils.MustParseDecCoins("1000000stake"), true)

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
	withdrawn, err := suite.dk.WithdrawDelegationRewards(ctx, investorAddr, sdk.ValAddress(valOwnerAddr))
	suite.Require().NoError(err)
	suite.Require().Equal("200000stake", withdrawn.String())

	// Allocate 1000000stake as rewards again
	ctx = suite.allocateTokensToValidator(ctx, validator, utils.MustParseDecCoins("1000000stake"), true)

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
