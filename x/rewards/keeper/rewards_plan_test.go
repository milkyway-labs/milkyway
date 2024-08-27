package keeper_test

import (
	"time"

	"github.com/milkyway-labs/milkyway/app/testutil"
	"github.com/milkyway-labs/milkyway/utils"
	rewardskeeper "github.com/milkyway-labs/milkyway/x/rewards/keeper"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
)

func (suite *KeeperTestSuite) TestCreateRewardsPlan_PoolOrOperatorNotFound() {
	// Cache the context to avoid errors
	ctx, _ := suite.Ctx.CacheContext()

	service, _ := suite.setupSampleServiceAndOperator(ctx)

	// Create an active rewards plan.
	amtPerDay := utils.MustParseCoins("100_000000service")
	planStartTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	planEndTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	rewardsMsgServer := rewardskeeper.NewMsgServer(suite.App.RewardsKeeper)

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
		service.Admin,
	))
	suite.Require().EqualError(err, "cannot get delegation target 2: operator not found: not found")
}
