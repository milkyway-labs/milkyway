package keeper_test

import (
	"time"

	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/milkyway-labs/milkyway/v7/app/testutil"
	"github.com/milkyway-labs/milkyway/v7/utils"
)

func (suite *KeeperTestSuite) TestVestingAccountRewards() {
	ctx, _ := suite.ctx.CacheContext()
	ctx = ctx.WithBlockTime(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))

	valAddr := testutil.TestAddress(10000)
	validator := suite.createValidator(
		ctx,
		valAddr,
		stakingtypes.NewCommissionRates(utils.MustParseDec("0"), utils.MustParseDec("0.2"), utils.MustParseDec("0.01")),
		utils.MustParseCoin("1000000stake"),
		true,
	)

	delAddr1 := testutil.TestAddress(1)
	suite.delegate(ctx, delAddr1.String(), validator.GetOperator(), utils.MustParseCoin("1000000stake"), true)
	delAddr2 := testutil.TestAddress(2)
	vestingEndTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	suite.createVestingAccount(
		ctx,
		testutil.TestAddress(10001).String(),
		delAddr2.String(),
		utils.MustParseCoins("1000000stake"),
		vestingEndTime.Unix(),
		false,
		true,
	)
	suite.delegate(ctx, delAddr2.String(), validator.GetOperator(), utils.MustParseCoin("1000000stake"), false)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(5 * time.Second))
	err := suite.k.AllocateTokensToValidator(ctx, validator, utils.MustParseDecCoins("1000000stake"))
	suite.Require().NoError(err)

	querier := distrkeeper.NewQuerier(suite.k.Keeper)
	// Delegator 1 receive 40% of rewards
	cacheCtx, _ := ctx.CacheContext()
	rewards, err := querier.DelegationRewards(cacheCtx, &distrtypes.QueryDelegationRewardsRequest{
		DelegatorAddress: delAddr1.String(),
		ValidatorAddress: validator.GetOperator(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal("400000.000000000000000000stake", rewards.Rewards.String())

	// Delegator 2 has a vesting account so it receives 20% of rewards
	cacheCtx, _ = ctx.CacheContext()
	rewards, err = querier.DelegationRewards(cacheCtx, &distrtypes.QueryDelegationRewardsRequest{
		DelegatorAddress: delAddr2.String(),
		ValidatorAddress: validator.GetOperator(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal("200000.000000000000000000stake", rewards.Rewards.String())

	// The rest of the rewards was allocated to the validator
	cacheCtx, _ = ctx.CacheContext()
	rewards, err = querier.DelegationRewards(cacheCtx, &distrtypes.QueryDelegationRewardsRequest{
		DelegatorAddress: valAddr.String(),
		ValidatorAddress: validator.GetOperator(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal("400000.000000000000000000stake", rewards.Rewards.String())
}
