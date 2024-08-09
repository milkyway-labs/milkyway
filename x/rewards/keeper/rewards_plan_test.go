package keeper_test

import (
	"time"

	"github.com/milkyway-labs/milkyway/utils"
	rewardskeeper "github.com/milkyway-labs/milkyway/x/rewards/keeper"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
)

func (s *KeeperTestSuite) TestCreateRewardsPlan_PoolOrOperatorNotFound() {
	service, _ := s.setupSampleServiceAndOperator()

	// Create an active rewards plan.
	amtPerDay := utils.MustParseCoins("100_000000service")
	planStartTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	planEndTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	rewardsMsgServer := rewardskeeper.NewMsgServer(s.App.RewardsKeeper)

	// There's no pool 1 yet.
	_, err := rewardsMsgServer.CreateRewardsPlan(s.Ctx, types.NewMsgCreateRewardsPlan(
		service.Admin, "Rewards Plan", service.ID, amtPerDay, planStartTime, planEndTime,
		types.NewWeightedPoolsDistribution(1, []types.PoolDistributionWeight{
			types.NewPoolDistributionWeight(1, 1),
			types.NewPoolDistributionWeight(2, 3),
		}),
		types.NewWeightedOperatorsDistribution(1, []types.OperatorDistributionWeight{
			types.NewOperatorDistributionWeight(1, 3),
			types.NewOperatorDistributionWeight(2, 2),
		}),
		types.NewBasicUsersDistribution(1)))
	s.Require().EqualError(err, "pool 1 not found: pool not found: not found")

	s.DelegatePool(utils.MustParseCoin("100_000000umilk"), utils.TestAddress(1).String(), true)
	s.DelegatePool(utils.MustParseCoin("100_000000uinit"), utils.TestAddress(2).String(), true)

	// After users delegates to pools, the pools are created, but there's no
	// operator 2 this time.
	_, err = rewardsMsgServer.CreateRewardsPlan(s.Ctx, types.NewMsgCreateRewardsPlan(
		service.Admin, "Rewards Plan", service.ID, amtPerDay, planStartTime, planEndTime,
		types.NewWeightedPoolsDistribution(1, []types.PoolDistributionWeight{
			types.NewPoolDistributionWeight(1, 1),
			types.NewPoolDistributionWeight(2, 3),
		}),
		types.NewWeightedOperatorsDistribution(1, []types.OperatorDistributionWeight{
			types.NewOperatorDistributionWeight(1, 3),
			types.NewOperatorDistributionWeight(2, 2),
		}),
		types.NewBasicUsersDistribution(1)))
	s.Require().EqualError(err, "operator 2 not found: operator not found: not found")
}
