package keeper_test

import (
	"time"

	"cosmossdk.io/collections"

	"github.com/milkyway-labs/milkyway/v7/app/testutil"
	"github.com/milkyway-labs/milkyway/v7/utils"
	rewardskeeper "github.com/milkyway-labs/milkyway/v7/x/rewards/keeper"
	"github.com/milkyway-labs/milkyway/v7/x/rewards/types"
)

func (suite *KeeperTestSuite) TestCreateRewardsPlan_PoolOrOperatorNotFound() {
	// Cache the context to avoid errors
	ctx, _ := suite.ctx.CacheContext()

	service, _ := suite.setupSampleServiceAndOperator(ctx)

	// Create an active rewards plan.
	amtPerDay := utils.MustParseCoins("100_000000service")
	planStartTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	planEndTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	rewardsMsgServer := rewardskeeper.NewMsgServer(suite.keeper)

	// There's no pool 1 yet.
	_, err := rewardsMsgServer.CreateRewardsPlan(ctx, types.NewMsgCreateRewardsPlan(
		service.ID,
		"Rewards Plan",
		amtPerDay,
		planStartTime,
		planEndTime,
		types.NewWeightedPoolsDistribution(1, []types.DistributionWeight{
			types.NewDistributionWeight(1, 1),
			types.NewDistributionWeight(2, 3),
		}),
		types.NewWeightedOperatorsDistribution(1, []types.DistributionWeight{
			types.NewDistributionWeight(1, 3),
			types.NewDistributionWeight(2, 2),
		}),
		types.NewBasicUsersDistribution(1),
		nil,
		service.Admin,
	))
	suite.Require().EqualError(err, "cannot get delegation target 1: pool not found: not found")

	suite.DelegatePool(ctx, utils.MustParseCoin("100_000000umilk"), testutil.TestAddress(1).String(), true)
	suite.DelegatePool(ctx, utils.MustParseCoin("100_000000uinit"), testutil.TestAddress(2).String(), true)

	// After users delegates to pools, the pools are created, but there's no
	// operator 2 this time.
	_, err = rewardsMsgServer.CreateRewardsPlan(ctx, types.NewMsgCreateRewardsPlan(
		service.ID,
		"Rewards Plan",
		amtPerDay,
		planStartTime,
		planEndTime,
		types.NewWeightedPoolsDistribution(1, []types.DistributionWeight{
			types.NewDistributionWeight(1, 1),
			types.NewDistributionWeight(2, 3),
		}),
		types.NewWeightedOperatorsDistribution(1, []types.DistributionWeight{
			types.NewDistributionWeight(1, 3),
			types.NewDistributionWeight(2, 2),
		}),
		types.NewBasicUsersDistribution(1),
		nil,
		service.Admin,
	))
	suite.Require().EqualError(err, "cannot get delegation target 2: operator not found: not found")
}

func (suite *KeeperTestSuite) TestTerminateEndedRewardsPlans() {
	// Cache the context to avoid errors
	ctx, _ := suite.ctx.CacheContext()

	service, _ := suite.setupSampleServiceAndOperator(ctx)

	// Create an active rewards plan.
	plan := suite.CreateBasicRewardsPlan(
		ctx,
		service.ID,
		utils.MustParseCoins("100_000000service"),
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		utils.MustParseCoins("10000_000000service"),
	)

	rewardsPoolAddr := plan.MustGetRewardsPoolAddress(suite.accountKeeper.AddressCodec())
	remaining := suite.bankKeeper.GetAllBalances(ctx, rewardsPoolAddr)
	suite.Require().Equal("10000000000service", remaining.String())

	// Change the block time so that the plan becomes no more active.
	ctx = ctx.WithBlockTime(time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC))

	// Terminate the ended rewards plans
	err := suite.keeper.TerminateEndedRewardsPlans(ctx)
	suite.Require().NoError(err)

	// The plan is removed.
	_, err = suite.keeper.GetRewardsPlan(ctx, plan.ID)
	suite.Require().ErrorIs(err, collections.ErrNotFound)

	// All remaining rewards are transferred to the service's address.
	remaining = suite.bankKeeper.GetAllBalances(ctx, rewardsPoolAddr)
	suite.Require().True(remaining.IsZero())

	// Check the service's address balances.
	serviceAddr, err := suite.accountKeeper.AddressCodec().StringToBytes(service.Address)
	suite.Require().NoError(err)
	serviceBalances := suite.bankKeeper.GetAllBalances(ctx, serviceAddr)
	suite.Require().Equal("10000000000service", serviceBalances.String())
}
