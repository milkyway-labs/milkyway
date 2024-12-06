package simulation

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/milkyway-labs/milkyway/v3/utils"
	restakingtypes "github.com/milkyway-labs/milkyway/v3/x/restaking/types"
	"github.com/milkyway-labs/milkyway/v3/x/rewards/types"
)

func RandomParams(r *rand.Rand, rewardsPlanCreationFeeDenoms []string) types.Params {
	rewardsPlanCreationFees := sdk.NewCoins()
	for _, denom := range rewardsPlanCreationFeeDenoms {
		rewardsPlanCreationFees = rewardsPlanCreationFees.Add(
			sdk.NewInt64Coin(denom, r.Int63()))
	}

	return types.NewParams(rewardsPlanCreationFees)
}

func RandomDistributionType(r *rand.Rand) types.DistributionType {
	switch r.Intn(3) {
	case 1:
		weightsCount := r.Intn(5) + 1

		var distributionWeights []types.DistributionWeight
		for i := 0; i < weightsCount; i++ {
			distributionWeights = append(distributionWeights, types.DistributionWeight{
				Weight:             r.Uint32(),
				DelegationTargetID: r.Uint32(),
			})
		}
		return &types.DistributionTypeWeighted{
			Weights: distributionWeights,
		}
	case 2:
		return &types.DistributionTypeEgalitarian{}
	default:
		return &types.DistributionTypeBasic{}
	}
}

func RandomDistribution(r *rand.Rand) types.Distribution {
	return types.NewDistribution(
		restakingtypes.DELEGATION_TYPE_POOL,
		r.Uint32(),
		RandomDistributionType(r),
	)
}

func RandomUsersDistributionType(_ *rand.Rand) types.UsersDistributionType {
	return &types.UsersDistributionTypeBasic{}
}

func RandomUsersDistribution(r *rand.Rand) types.UsersDistribution {
	return types.NewUsersDistribution(r.Uint32(), RandomUsersDistributionType(r))
}

func RandomRewardsPlan(r *rand.Rand, serviceID uint32, amtPerDeyDenoms []string) types.RewardsPlan {
	randomAmountPerDays := sdk.NewCoins()
	for _, denom := range amtPerDeyDenoms {
		randomAmountPerDays = randomAmountPerDays.Add(
			sdk.NewInt64Coin(denom, r.Int63()))
	}

	return types.NewRewardsPlan(
		r.Uint64(),
		simtypes.RandStringOfLength(r, 32),
		serviceID,
		randomAmountPerDays,
		simtypes.RandTimestamp(r),
		simtypes.RandTimestamp(r),
		RandomDistribution(r),
		RandomDistribution(r),
		RandomUsersDistribution(r),
	)
}

func RandomRewardsPlans(r *rand.Rand, allowedDenoms []string) []types.RewardsPlan {
	// Get a random numer of rewards plans to create
	rewardsPlanCount := r.Intn(30)

	// Generate the rewards plans
	var rewardsPlans []types.RewardsPlan
	for id := 0; id < rewardsPlanCount; id++ {
		rewardsPlans = append(rewardsPlans, RandomRewardsPlan(r, uint32(r.Intn(10)+1), allowedDenoms))
	}

	return rewardsPlans
}

func RandomDelegatorWithdrawInfos(r *rand.Rand, accs []simtypes.Account) []types.DelegatorWithdrawInfo {
	count := r.Intn(len(accs))

	var infos []types.DelegatorWithdrawInfo
	for i := 0; i < count; i++ {
		randomAccount, _ := simtypes.RandomAcc(r, accs)
		infos = append(infos, types.DelegatorWithdrawInfo{
			DelegatorAddress: randomAccount.Address.String(),
			WithdrawAddress:  randomAccount.Address.String(),
		})
	}

	return utils.RemoveDuplicatesFunc(infos, func(i types.DelegatorWithdrawInfo) string {
		return i.DelegatorAddress
	})
}
