package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/milkyway-labs/milkyway/v7/app/testutil"
	"github.com/milkyway-labs/milkyway/v7/utils"
	"github.com/milkyway-labs/milkyway/v7/x/distribution/types"
)

func (suite *KeeperTestSuite) TestVestingAccountRewards() {
	testCases := []struct {
		name                     string
		rewardsRatio             sdkmath.LegacyDec
		expNormalAccountRewards  string
		expVestingAccountRewards string
	}{
		{
			name:                     "0% vesting account rewards ratio",
			rewardsRatio:             sdkmath.LegacyZeroDec(),
			expNormalAccountRewards:  "500000.000000000000000000stake",
			expVestingAccountRewards: "",
		},
		{
			name:                     "10% vesting account rewards ratio",
			rewardsRatio:             utils.MustParseDec("0.1"),
			expNormalAccountRewards:  "476190.476190476190000000stake", // 1000000 * 1 / 2.1
			expVestingAccountRewards: "47619.047619047619000000stake",  // 1000000 * 0.1 / 2.1
		},
		{
			name:                     "50% vesting account rewards ratio",
			rewardsRatio:             utils.MustParseDec("0.5"),
			expNormalAccountRewards:  "400000.000000000000000000stake",
			expVestingAccountRewards: "200000.000000000000000000stake",
		},
		{
			name:                     "100% vesting account rewards ratio",
			rewardsRatio:             sdkmath.LegacyOneDec(),
			expNormalAccountRewards:  "333333.333333333333000000stake",
			expVestingAccountRewards: "333333.333333333333000000stake",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			// Override the rewards ratio
			types.VestingAccountRewardsRatio = tc.rewardsRatio

			ctx := suite.ctx.WithBlockTime(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))

			normalAddr := testutil.TestAddress(1)
			vestingAddr := testutil.TestAddress(2)
			vestingEndTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
			suite.createVestingAccount(
				ctx,
				testutil.TestAddress(10001).String(),
				vestingAddr.String(),
				utils.MustParseCoins("1000000stake"),
				vestingEndTime.Unix(),
				false,
				true,
			)

			valAddr := testutil.TestAddress(10000)
			validator := suite.createValidator(
				ctx,
				valAddr,
				stakingtypes.NewCommissionRates(utils.MustParseDec("0"), utils.MustParseDec("0.2"), utils.MustParseDec("0.01")),
				utils.MustParseCoin("1000000stake"),
				true,
			)
			suite.delegate(ctx, normalAddr.String(), validator.GetOperator(), utils.MustParseCoin("1000000stake"), true)
			suite.delegate(ctx, vestingAddr.String(), validator.GetOperator(), utils.MustParseCoin("1000000stake"), false)

			ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(5 * time.Second))
			err := suite.k.AllocateTokensToValidator(ctx, validator, utils.MustParseDecCoins("1000000stake"))
			suite.Require().NoError(err)

			querier := distrkeeper.NewQuerier(suite.k.Keeper)
			// Query the normal account's rewards
			cacheCtx, _ := ctx.CacheContext()
			rewards, err := querier.DelegationRewards(cacheCtx, &distrtypes.QueryDelegationRewardsRequest{
				DelegatorAddress: normalAddr.String(),
				ValidatorAddress: validator.GetOperator(),
			})
			suite.Require().NoError(err)
			suite.Require().Equal(tc.expNormalAccountRewards, rewards.Rewards.String())

			// Query the vesting account's rewards
			cacheCtx, _ = ctx.CacheContext()
			rewards, err = querier.DelegationRewards(cacheCtx, &distrtypes.QueryDelegationRewardsRequest{
				DelegatorAddress: vestingAddr.String(),
				ValidatorAddress: validator.GetOperator(),
			})
			suite.Require().NoError(err)
			suite.Require().Equal(tc.expVestingAccountRewards, rewards.Rewards.String())
		})
	}
}
