package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v10/app/testutil"
	"github.com/milkyway-labs/milkyway/v10/utils"
	restakingtypes "github.com/milkyway-labs/milkyway/v10/x/restaking/types"
	rewardstypes "github.com/milkyway-labs/milkyway/v10/x/rewards/types"
)

func (suite *KeeperTestSuite) TestCoveredLockedSharesRewards() {
	ctx, _ := suite.ctx.CacheContext()
	ctx = ctx.WithBlockTime(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))

	// Register both currencies so that rewards can be allocated
	suite.registerCurrency(ctx, "stake", "STAKE", 6, utils.MustParseDec("2"))
	suite.registerCurrency(ctx, "locked/stake", "STAKE", 6, utils.MustParseDec("2"))

	// Call AllocateRewards to set the last rewards allocation time
	err := suite.rewardsKeeper.AllocateRewards(ctx)
	suite.Require().NoError(err)

	params, err := suite.k.GetParams(ctx)
	suite.Require().NoError(err)
	params.InsurancePercentage = utils.MustParseDec("2") // 2%
	err = suite.k.SetParams(ctx, params)
	suite.Require().NoError(err)

	// Create a service and an operator
	suite.createService(ctx, 1)
	suite.createOperator(ctx, 1)

	// The operator joins the service
	err = suite.restakingKeeper.AddServiceToOperatorJoinedServices(ctx, 1, 1)
	suite.Require().NoError(err)

	// A random delegator delegates 10 MILK to the operator
	delAddr1 := testutil.TestAddress(3)
	suite.fundAccount(ctx, delAddr1.String(), utils.MustParseCoins("10000000stake"))
	_, err = suite.restakingKeeper.DelegateToOperator(ctx, 1, utils.MustParseCoins("10000000stake"), delAddr1.String())
	suite.Require().NoError(err)

	// This is the test account
	delAddr2 := testutil.TestAddress(4)
	suite.mintLockedRepresentation(ctx, delAddr2.String(), utils.MustParseCoins("10000000stake"))

	// The service admin creates a rewards plan
	plan, err := suite.rewardsKeeper.CreateRewardsPlan(
		ctx,
		"Plan",
		1,
		utils.MustParseCoin("100000000stake"),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		rewardstypes.NewBasicPoolsDistribution(0),
		rewardstypes.NewBasicOperatorsDistribution(0),
		rewardstypes.NewBasicUsersDistribution(0),
	)
	suite.Require().NoError(err)
	suite.fundAccount(ctx, plan.RewardsPool, utils.MustParseCoins("100000000stake"))

	testCases := []struct {
		name string
		run  func(ctx sdk.Context)
	}{
		{
			name: "uncovered locked shares rewards should be zero",
			run: func(ctx sdk.Context) {
				_, err := suite.restakingKeeper.DelegateToOperator(ctx, 1, utils.MustParseCoins("10000000locked/stake"), delAddr2.String())
				suite.Require().NoError(err)

				// Increment the block height and allocate rewards
				ctx = suite.allocateRewards(ctx, 10*time.Second)

				// The random delegator received all rewards
				rewards, err := suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr1, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Require().Equal("11574.000000000000000000stake", rewards.Sum().String())

				// The test account received no rewards, because it doesn't have insurance fund
				rewards, err = suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr2, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Require().Empty(rewards)
			},
		},
		{
			name: "only covered locked shares receives rewards - fund and delegate",
			run: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx, delAddr2.String(), utils.MustParseCoins("100000stake")) // 1%

				_, err := suite.restakingKeeper.DelegateToOperator(ctx, 1, utils.MustParseCoins("10000000locked/stake"), delAddr2.String())
				suite.Require().NoError(err)

				// Increment the block height and allocate rewards
				ctx = suite.allocateRewards(ctx, 10*time.Second)

				// The random delegator received 7716 = 11574 * 10 / 15
				rewards, err := suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr1, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("7716.000000000000000000stake", rewards.Sum().String())

				// The test account received 3858 = 11574 * 5 / 15, since it has 1% of insurance
				// fund which can cover only the half of its total delegation amount, 10 locked
				// MILK
				rewards, err = suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr2, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("3858.000000000000000000stake", rewards.Sum().String())
			},
		},
		{
			name: "only covered locked shares receives rewards - delegate and fund",
			run: func(ctx sdk.Context) {
				_, err := suite.restakingKeeper.DelegateToOperator(ctx, 1, utils.MustParseCoins("10000000locked/stake"), delAddr2.String())
				suite.Require().NoError(err)

				suite.fundAccountInsuranceFund(ctx, delAddr2.String(), utils.MustParseCoins("100000stake")) // 1%

				// Increment the block height and allocate rewards
				ctx = suite.allocateRewards(ctx, 10*time.Second)

				// The random delegator received 7716 = 11574 * 10 / 15
				rewards, err := suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr1, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("7716.000000000000000000stake", rewards.Sum().String())

				// The test account received 3858 = 11574 * 5 / 15, since it has 1% of insurance
				// fund which can cover only the half of its total delegation amount, 10 locked
				// MILK
				rewards, err = suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr2, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("3858.000000000000000000stake", rewards.Sum().String())
			},
		},
		{
			name: "fully covered locked shares receives full rewards - fund and delegate",
			run: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx, delAddr2.String(), utils.MustParseCoins("300000stake")) // 3%

				_, err := suite.restakingKeeper.DelegateToOperator(ctx, 1, utils.MustParseCoins("10000000locked/stake"), delAddr2.String())
				suite.Require().NoError(err)

				// Increment the block height and allocate rewards
				ctx = suite.allocateRewards(ctx, 10*time.Second)

				// Both users received the same rewards
				rewards, err := suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr1, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("5787.000000000000000000stake", rewards.Sum().String())

				rewards, err = suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr2, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("5787.000000000000000000stake", rewards.Sum().String())
			},
		},
		{
			name: "fully covered locked shares receives full rewards - delegate and fund",
			run: func(ctx sdk.Context) {
				_, err := suite.restakingKeeper.DelegateToOperator(ctx, 1, utils.MustParseCoins("10000000locked/stake"), delAddr2.String())
				suite.Require().NoError(err)

				suite.fundAccountInsuranceFund(ctx, delAddr2.String(), utils.MustParseCoins("300000stake")) // 3%

				// Increment the block height and allocate rewards
				ctx = suite.allocateRewards(ctx, 10*time.Second)

				// Both users received the same rewards
				rewards, err := suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr1, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("5787.000000000000000000stake", rewards.Sum().String())

				rewards, err = suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr2, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("5787.000000000000000000stake", rewards.Sum().String())
			},
		},
		{
			name: "add insurance fund after receiving rewards",
			run: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx, delAddr2.String(), utils.MustParseCoins("100000stake")) // 1%

				_, err := suite.restakingKeeper.DelegateToOperator(ctx, 1, utils.MustParseCoins("10000000locked/stake"), delAddr2.String())
				suite.Require().NoError(err)

				// Increment the block height and allocate rewards
				ctx = suite.allocateRewards(ctx, 10*time.Second)

				rewards, err := suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr2, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("3858.000000000000000000stake", rewards.Sum().String())

				suite.fundAccountInsuranceFund(ctx, delAddr2.String(), utils.MustParseCoins("100000stake")) // 1% more

				// Increment the block height and allocate rewards
				ctx = suite.allocateRewards(ctx, 10*time.Second)

				rewards, err = suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr1, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("13503.000000000000000000stake", rewards.Sum().String())

				rewards, err = suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr2, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("5787.000000000000000000stake", rewards.Sum().String())
			},
		},
		{
			name: "withdraw insurance fund after receiving rewards",
			run: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx, delAddr2.String(), utils.MustParseCoins("200000stake")) // 2%

				_, err := suite.restakingKeeper.DelegateToOperator(ctx, 1, utils.MustParseCoins("10000000locked/stake"), delAddr2.String())
				suite.Require().NoError(err)

				// Increment the block height and allocate rewards
				ctx = suite.allocateRewards(ctx, 10*time.Second)

				rewards, err := suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr2, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("5787.000000000000000000stake", rewards.Sum().String())

				err = suite.k.WithdrawFromUserInsuranceFund(ctx, delAddr2.String(), utils.MustParseCoins("100000stake")) // withdraw 1%
				suite.Require().NoError(err)

				// Increment the block height and allocate rewards
				ctx = suite.allocateRewards(ctx, 10*time.Second)

				rewards, err = suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr1, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("13503.000000000000000000stake", rewards.Sum().String())

				rewards, err = suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr2, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("3858.000000000000000000stake", rewards.Sum().String())
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			ctx, _ := ctx.CacheContext()
			tc.run(ctx)
		})
	}
}

func (suite *KeeperTestSuite) TestCoveredLockedSharesRewards_UpdateInsurancePercentage() {
	ctx, _ := suite.ctx.CacheContext()
	ctx = ctx.WithBlockTime(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))

	// Register both currencies so that rewards can be allocated
	suite.registerCurrency(ctx, "stake", "STAKE", 6, utils.MustParseDec("2"))
	suite.registerCurrency(ctx, "locked/stake", "STAKE", 6, utils.MustParseDec("2"))

	// Call AllocateRewards to set the last rewards allocation time
	err := suite.rewardsKeeper.AllocateRewards(ctx)
	suite.Require().NoError(err)

	params, err := suite.k.GetParams(ctx)
	suite.Require().NoError(err)
	params.InsurancePercentage = utils.MustParseDec("2") // 2%
	err = suite.k.SetParams(ctx, params)
	suite.Require().NoError(err)

	// Create a service and an operator
	suite.createService(ctx, 1)
	suite.createOperator(ctx, 1)

	// The operator joins the service
	err = suite.restakingKeeper.AddServiceToOperatorJoinedServices(ctx, 1, 1)
	suite.Require().NoError(err)

	// A random delegator delegates 10 MILK to the operator
	delAddr1 := testutil.TestAddress(3)
	suite.fundAccount(ctx, delAddr1.String(), utils.MustParseCoins("10000000stake"))
	_, err = suite.restakingKeeper.DelegateToOperator(ctx, 1, utils.MustParseCoins("10000000stake"), delAddr1.String())
	suite.Require().NoError(err)

	// This is the test account
	delAddr2 := testutil.TestAddress(4)
	suite.mintLockedRepresentation(ctx, delAddr2.String(), utils.MustParseCoins("10000000stake"))
	suite.fundAccountInsuranceFund(ctx, delAddr2.String(), utils.MustParseCoins("200000stake")) // 2%

	// The service admin creates a rewards plan
	plan, err := suite.rewardsKeeper.CreateRewardsPlan(
		ctx,
		"Plan",
		1,
		utils.MustParseCoin("100000000stake"),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		rewardstypes.NewBasicPoolsDistribution(0),
		rewardstypes.NewBasicOperatorsDistribution(0),
		rewardstypes.NewBasicUsersDistribution(0),
	)
	suite.Require().NoError(err)
	suite.fundAccount(ctx, plan.RewardsPool, utils.MustParseCoins("100000000stake"))

	testCases := []struct {
		name string
		run  func(ctx sdk.Context)
	}{
		{
			name: "increase percentage, delegate and allocate",
			run: func(ctx sdk.Context) {
				// Double the insurance percentage
				params.InsurancePercentage = utils.MustParseDec("4") // 4%
				err = suite.k.SetParams(ctx, params)
				suite.Require().NoError(err)

				_, err := suite.restakingKeeper.DelegateToOperator(ctx, 1, utils.MustParseCoins("10000000locked/stake"), delAddr2.String())
				suite.Require().NoError(err)

				// Increment the block height and allocate rewards
				ctx = suite.allocateRewards(ctx, 10*time.Second)

				// The random delegator received 7716 = 11574 * 10 / 15
				rewards, err := suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr1, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("7716.000000000000000000stake", rewards.Sum().String())

				// The test account received 3858 = 11574 * 5 / 15, since it has 1% of insurance
				// fund which can cover only the half of its total delegation amount, 10 locked
				// MILK
				rewards, err = suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr2, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("3858.000000000000000000stake", rewards.Sum().String())
			},
		},
		{
			name: "delegate, allocate, increase percentage and allocate",
			run: func(ctx sdk.Context) {
				_, err := suite.restakingKeeper.DelegateToOperator(ctx, 1, utils.MustParseCoins("10000000locked/stake"), delAddr2.String())
				suite.Require().NoError(err)

				// Increment the block height and allocate rewards
				ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(10 * time.Second))
				err = suite.rewardsKeeper.AllocateRewards(ctx)
				suite.Require().NoError(err)

				// Double the insurance percentage
				params.InsurancePercentage = utils.MustParseDec("4") // 4%
				err = suite.k.SetParams(ctx, params)
				suite.Require().NoError(err)

				// Increment the block height and allocate rewards again
				ctx = suite.allocateRewards(ctx, 10*time.Second)

				// The random delegator received 7716(=11574 * 10 / 15) more, so the total will
				// be 13503
				rewards, err := suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr1, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("13503.000000000000000000stake", rewards.Sum().String())

				// The test account received 3858 = 11574 * 5 / 15, since it has 0.2 MILK inside
				// insurance fund which can cover only the half of its total delegation amount(10
				// locked MILK) with the updated insurance percentage(4%)
				rewards, err = suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr2, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("3858.000000000000000000stake", rewards.Sum().String())
			},
		},
		{
			name: "delegate, increase percentage and allocate",
			run: func(ctx sdk.Context) {
				_, err := suite.restakingKeeper.DelegateToOperator(ctx, 1, utils.MustParseCoins("10000000locked/stake"), delAddr2.String())
				suite.Require().NoError(err)

				// Double the insurance percentage
				params.InsurancePercentage = utils.MustParseDec("4") // 4%
				err = suite.k.SetParams(ctx, params)
				suite.Require().NoError(err)

				// Increment the block height and allocate rewards
				ctx = suite.allocateRewards(ctx, 10*time.Second)

				// The random delegator received 7716 = 11574 * 10 / 15
				rewards, err := suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr1, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("7716.000000000000000000stake", rewards.Sum().String())

				// The test account received 3858 = 11574 * 5 / 15, since it has 1% of insurance
				// fund which can cover only the half of its total delegation amount, 10 locked
				// MILK
				rewards, err = suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr2, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("3858.000000000000000000stake", rewards.Sum().String())
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			ctx, _ := ctx.CacheContext()
			tc.run(ctx)
		})
	}
}
