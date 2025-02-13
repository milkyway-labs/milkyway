package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v9/app/testutil"
	"github.com/milkyway-labs/milkyway/v9/utils"
	operatorstypes "github.com/milkyway-labs/milkyway/v9/x/operators/types"
	restakingtypes "github.com/milkyway-labs/milkyway/v9/x/restaking/types"
	rewardstypes "github.com/milkyway-labs/milkyway/v9/x/rewards/types"
	servicestypes "github.com/milkyway-labs/milkyway/v9/x/services/types"
)

func (suite *KeeperTestSuite) TestCoveredLockedSharesRewards() {
	ctx, _ := suite.ctx.CacheContext()
	ctx = ctx.WithBlockTime(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))

	// Call AllocateRewards to set the last rewards allocation time
	err := suite.rewardsKeeper.AllocateRewards(ctx)
	suite.Require().NoError(err)

	params, err := suite.k.GetParams(ctx)
	suite.Require().NoError(err)
	params.InsurancePercentage = utils.MustParseDec("2") // 2%
	err = suite.k.SetParams(ctx, params)
	suite.Require().NoError(err)

	// Create a service
	serviceAdmin := testutil.TestAddress(1)
	err = suite.sk.CreateService(ctx, servicestypes.NewService(
		1,
		servicestypes.SERVICE_STATUS_ACTIVE,
		"Service",
		"",
		"",
		"",
		serviceAdmin.String(),
		false,
	))
	suite.Require().NoError(err)

	// Create an operator
	operatorAdmin := testutil.TestAddress(2)
	err = suite.ok.CreateOperator(ctx, operatorstypes.NewOperator(
		1,
		operatorstypes.OPERATOR_STATUS_ACTIVE,
		"Operator",
		"",
		"",
		operatorAdmin.String(),
	))
	suite.Require().NoError(err)
	// The operator joins the service
	err = suite.restakingKeeper.AddServiceToOperatorJoinedServices(ctx, 1, 1)
	suite.Require().NoError(err)

	// A random delegator delegates 10 MILK to the operator
	delAddr1 := testutil.TestAddress(3)
	suite.fundAccount(ctx, delAddr1.String(), utils.MustParseCoins("10000000umilk"))
	_, err = suite.restakingKeeper.DelegateToOperator(ctx, 1, utils.MustParseCoins("10000000umilk"), delAddr1.String())
	suite.Require().NoError(err)

	// This is the test account
	delAddr2 := testutil.TestAddress(4)
	suite.mintLockedRepresentation(ctx, delAddr2.String(), utils.MustParseCoins("10000000umilk"))

	// Register both currencies so that rewards can be allocated
	suite.registerCurrency(ctx, "umilk", "MILK", 6, utils.MustParseDec("2"))
	suite.registerCurrency(ctx, "locked/umilk", "MILK", 6, utils.MustParseDec("2"))

	// The service admin creates a rewards plan
	plan, err := suite.rewardsKeeper.CreateRewardsPlan(
		ctx,
		"Plan",
		1,
		utils.MustParseCoin("100000000umilk"),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		rewardstypes.NewBasicPoolsDistribution(0),
		rewardstypes.NewBasicOperatorsDistribution(0),
		rewardstypes.NewBasicUsersDistribution(0),
	)
	suite.Require().NoError(err)
	suite.fundAccount(ctx, plan.RewardsPool, utils.MustParseCoins("100000000umilk"))

	testCases := []struct {
		name string
		run  func(ctx sdk.Context)
	}{
		{
			name: "uncovered locked shares rewards should be zero",
			run: func(ctx sdk.Context) {
				_, err := suite.restakingKeeper.DelegateToOperator(ctx, 1, utils.MustParseCoins("10000000locked/umilk"), delAddr2.String())
				suite.Require().NoError(err)

				// Increment the block height and allocate rewards
				ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(10 * time.Second))
				err = suite.rewardsKeeper.AllocateRewards(ctx)
				suite.Require().NoError(err)

				// The random delegator received all rewards
				rewards, err := suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr1, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Require().Equal("11574.000000000000000000umilk", rewards.Sum().String())

				// The test account received no rewards, because it doesn't have insurance fund
				rewards, err = suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr2, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Require().Empty(rewards)
			},
		},
		{
			name: "only covered locked shares receives rewards - fund and delegate",
			run: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx, delAddr2.String(), utils.MustParseCoins("100000umilk")) // 1%

				_, err := suite.restakingKeeper.DelegateToOperator(ctx, 1, utils.MustParseCoins("10000000locked/umilk"), delAddr2.String())
				suite.Require().NoError(err)

				// Increment the block height and allocate rewards
				ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(10 * time.Second))
				err = suite.rewardsKeeper.AllocateRewards(ctx)
				suite.Require().NoError(err)

				// The random delegator received 7716 = 11574 * 10 / 15
				rewards, err := suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr1, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("7716.000000000000000000umilk", rewards.Sum().String())

				// The test account received 3858 = 11574 * 5 / 15, since it has 1% of insurance
				// fund which can cover only the half of its total delegation amount, 10 locked
				// MILK
				rewards, err = suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr2, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("3858.000000000000000000umilk", rewards.Sum().String())
			},
		},
		{
			name: "only covered locked shares receives rewards - delegate and fund",
			run: func(ctx sdk.Context) {
				_, err := suite.restakingKeeper.DelegateToOperator(ctx, 1, utils.MustParseCoins("10000000locked/umilk"), delAddr2.String())
				suite.Require().NoError(err)

				suite.fundAccountInsuranceFund(ctx, delAddr2.String(), utils.MustParseCoins("100000umilk")) // 1%

				// Increment the block height and allocate rewards
				ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(10 * time.Second))
				err = suite.rewardsKeeper.AllocateRewards(ctx)
				suite.Require().NoError(err)

				// The random delegator received 7716 = 11574 * 10 / 15
				rewards, err := suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr1, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("7716.000000000000000000umilk", rewards.Sum().String())

				// The test account received 3858 = 11574 * 5 / 15, since it has 1% of insurance
				// fund which can cover only the half of its total delegation amount, 10 locked
				// MILK
				rewards, err = suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr2, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("3858.000000000000000000umilk", rewards.Sum().String())
			},
		},
		{
			name: "fully covered locked shares receives full rewards - fund and delegate",
			run: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx, delAddr2.String(), utils.MustParseCoins("300000umilk")) // 3%

				_, err := suite.restakingKeeper.DelegateToOperator(ctx, 1, utils.MustParseCoins("10000000locked/umilk"), delAddr2.String())
				suite.Require().NoError(err)

				// Increment the block height and allocate rewards
				ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(10 * time.Second))
				err = suite.rewardsKeeper.AllocateRewards(ctx)
				suite.Require().NoError(err)

				// Both users received the same rewards
				rewards, err := suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr1, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("5787.000000000000000000umilk", rewards.Sum().String())

				rewards, err = suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr2, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("5787.000000000000000000umilk", rewards.Sum().String())
			},
		},
		{
			name: "fully covered locked shares receives full rewards - delegate and fund",
			run: func(ctx sdk.Context) {
				_, err := suite.restakingKeeper.DelegateToOperator(ctx, 1, utils.MustParseCoins("10000000locked/umilk"), delAddr2.String())
				suite.Require().NoError(err)

				suite.fundAccountInsuranceFund(ctx, delAddr2.String(), utils.MustParseCoins("300000umilk")) // 3%

				// Increment the block height and allocate rewards
				ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(10 * time.Second))
				err = suite.rewardsKeeper.AllocateRewards(ctx)
				suite.Require().NoError(err)

				// Both users received the same rewards
				rewards, err := suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr1, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("5787.000000000000000000umilk", rewards.Sum().String())

				rewards, err = suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr2, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("5787.000000000000000000umilk", rewards.Sum().String())
			},
		},
		{
			name: "add insurance fund after receiving rewards",
			run: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx, delAddr2.String(), utils.MustParseCoins("100000umilk")) // 1%

				_, err := suite.restakingKeeper.DelegateToOperator(ctx, 1, utils.MustParseCoins("10000000locked/umilk"), delAddr2.String())
				suite.Require().NoError(err)

				// Increment the block height and allocate rewards
				ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(10 * time.Second))
				err = suite.rewardsKeeper.AllocateRewards(ctx)
				suite.Require().NoError(err)

				rewards, err := suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr2, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("3858.000000000000000000umilk", rewards.Sum().String())

				suite.fundAccountInsuranceFund(ctx, delAddr2.String(), utils.MustParseCoins("100000umilk")) // 1% more

				// Increment the block height and allocate rewards
				ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(10 * time.Second))
				err = suite.rewardsKeeper.AllocateRewards(ctx)
				suite.Require().NoError(err)

				rewards, err = suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr1, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("13503.000000000000000000umilk", rewards.Sum().String())

				rewards, err = suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr2, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("5787.000000000000000000umilk", rewards.Sum().String())
			},
		},
		{
			name: "withdraw insurance fund after receiving rewards",
			run: func(ctx sdk.Context) {
				suite.fundAccountInsuranceFund(ctx, delAddr2.String(), utils.MustParseCoins("200000umilk")) // 2%

				_, err := suite.restakingKeeper.DelegateToOperator(ctx, 1, utils.MustParseCoins("10000000locked/umilk"), delAddr2.String())
				suite.Require().NoError(err)

				// Increment the block height and allocate rewards
				ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(10 * time.Second))
				err = suite.rewardsKeeper.AllocateRewards(ctx)
				suite.Require().NoError(err)

				rewards, err := suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr2, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("5787.000000000000000000umilk", rewards.Sum().String())

				err = suite.k.WithdrawFromUserInsuranceFund(ctx, delAddr2.String(), utils.MustParseCoins("100000umilk")) // withdraw 1%
				suite.Require().NoError(err)

				// Increment the block height and allocate rewards
				ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(10 * time.Second))
				err = suite.rewardsKeeper.AllocateRewards(ctx)
				suite.Require().NoError(err)

				rewards, err = suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr1, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("13503.000000000000000000umilk", rewards.Sum().String())

				rewards, err = suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr2, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("3858.000000000000000000umilk", rewards.Sum().String())
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

	// Call AllocateRewards to set the last rewards allocation time
	err := suite.rewardsKeeper.AllocateRewards(ctx)
	suite.Require().NoError(err)

	params, err := suite.k.GetParams(ctx)
	suite.Require().NoError(err)
	params.InsurancePercentage = utils.MustParseDec("2") // 2%
	err = suite.k.SetParams(ctx, params)
	suite.Require().NoError(err)

	// Create a service
	serviceAdmin := testutil.TestAddress(1)
	err = suite.sk.CreateService(ctx, servicestypes.NewService(
		1,
		servicestypes.SERVICE_STATUS_ACTIVE,
		"Service",
		"",
		"",
		"",
		serviceAdmin.String(),
		false,
	))
	suite.Require().NoError(err)

	// Create an operator
	operatorAdmin := testutil.TestAddress(2)
	err = suite.ok.CreateOperator(ctx, operatorstypes.NewOperator(
		1,
		operatorstypes.OPERATOR_STATUS_ACTIVE,
		"Operator",
		"",
		"",
		operatorAdmin.String(),
	))
	suite.Require().NoError(err)
	// The operator joins the service
	err = suite.restakingKeeper.AddServiceToOperatorJoinedServices(ctx, 1, 1)
	suite.Require().NoError(err)

	// A random delegator delegates 10 MILK to the operator
	delAddr1 := testutil.TestAddress(3)
	suite.fundAccount(ctx, delAddr1.String(), utils.MustParseCoins("10000000umilk"))
	_, err = suite.restakingKeeper.DelegateToOperator(ctx, 1, utils.MustParseCoins("10000000umilk"), delAddr1.String())
	suite.Require().NoError(err)

	// This is the test account
	delAddr2 := testutil.TestAddress(4)
	suite.mintLockedRepresentation(ctx, delAddr2.String(), utils.MustParseCoins("10000000umilk"))
	suite.fundAccountInsuranceFund(ctx, delAddr2.String(), utils.MustParseCoins("200000umilk")) // 2%

	// Register both currencies so that rewards can be allocated
	suite.registerCurrency(ctx, "umilk", "MILK", 6, utils.MustParseDec("2"))
	suite.registerCurrency(ctx, "locked/umilk", "MILK", 6, utils.MustParseDec("2"))

	// The service admin creates a rewards plan
	plan, err := suite.rewardsKeeper.CreateRewardsPlan(
		ctx,
		"Plan",
		1,
		utils.MustParseCoin("100000000umilk"),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		rewardstypes.NewBasicPoolsDistribution(0),
		rewardstypes.NewBasicOperatorsDistribution(0),
		rewardstypes.NewBasicUsersDistribution(0),
	)
	suite.Require().NoError(err)
	suite.fundAccount(ctx, plan.RewardsPool, utils.MustParseCoins("100000000umilk"))

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

				_, err := suite.restakingKeeper.DelegateToOperator(ctx, 1, utils.MustParseCoins("10000000locked/umilk"), delAddr2.String())
				suite.Require().NoError(err)

				// Increment the block height and allocate rewards
				ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(10 * time.Second))
				err = suite.rewardsKeeper.AllocateRewards(ctx)
				suite.Require().NoError(err)

				// The random delegator received 7716 = 11574 * 10 / 15
				rewards, err := suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr1, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("7716.000000000000000000umilk", rewards.Sum().String())

				// The test account received 3858 = 11574 * 5 / 15, since it has 1% of insurance
				// fund which can cover only the half of its total delegation amount, 10 locked
				// MILK
				rewards, err = suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr2, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("3858.000000000000000000umilk", rewards.Sum().String())
			},
		},
		{
			name: "delegate, allocate, increase percentage and allocate",
			run: func(ctx sdk.Context) {
				_, err := suite.restakingKeeper.DelegateToOperator(ctx, 1, utils.MustParseCoins("10000000locked/umilk"), delAddr2.String())
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
				ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(10 * time.Second))
				err = suite.rewardsKeeper.AllocateRewards(ctx)
				suite.Require().NoError(err)

				// The random delegator received 7716(=11574 * 10 / 15) more, so the total will
				// be 13503
				rewards, err := suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr1, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("13503.000000000000000000umilk", rewards.Sum().String())

				// The test account received 3858 = 11574 * 5 / 15, since it has 0.2 MILK inside
				// insurance fund which can cover only the half of its total delegation amount(10
				// locked MILK) with the updated insurance percentage(4%)
				rewards, err = suite.rewardsKeeper.GetDelegationRewards(ctx, delAddr2, restakingtypes.DELEGATION_TYPE_OPERATOR, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal("3858.000000000000000000umilk", rewards.Sum().String())
			},
		},
		// TODO: delegate, change percentage and allocate
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			ctx, _ := ctx.CacheContext()
			tc.run(ctx)
		})
	}
}
