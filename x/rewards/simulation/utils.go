package simulation

import (
	"math/rand"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/milkyway-labs/milkyway/v3/testutils/simtesting"
	"github.com/milkyway-labs/milkyway/v3/utils"
	poolstypes "github.com/milkyway-labs/milkyway/v3/x/pools/types"
	restakingtypes "github.com/milkyway-labs/milkyway/v3/x/restaking/types"
	"github.com/milkyway-labs/milkyway/v3/x/rewards/keeper"
	"github.com/milkyway-labs/milkyway/v3/x/rewards/types"
	servicestypes "github.com/milkyway-labs/milkyway/v3/x/services/types"
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

func RandomDistribution(r *rand.Rand, delegationType restakingtypes.DelegationType) types.Distribution {
	return types.NewDistribution(
		delegationType,
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
		RandomDistribution(r, restakingtypes.DELEGATION_TYPE_POOL),
		RandomDistribution(r, restakingtypes.DELEGATION_TYPE_OPERATOR),
		RandomUsersDistribution(r),
	)
}

func RandomRewardsPlans(r *rand.Rand, services []servicestypes.Service, allowedDenoms []string) []types.RewardsPlan {
	// We can't create a rewards if there are no services
	if len(services) == 0 {
		return nil
	}

	// Get a random numer of rewards plans to create
	rewardsPlanCount := r.Intn(30)

	// Generate the rewards plans
	var rewardsPlans []types.RewardsPlan
	for id := 0; id < rewardsPlanCount; id++ {
		serviceIndex := r.Intn(len(services))
		rewardsPlans = append(rewardsPlans, RandomRewardsPlan(r, services[serviceIndex].ID, allowedDenoms))
	}

	return rewardsPlans
}

func GetRandomExistingRewardsPlan(r *rand.Rand, ctx sdk.Context, k *keeper.Keeper) (types.RewardsPlan, bool) {
	var plans []types.RewardsPlan
	k.RewardsPlans.Walk(ctx, nil, func(key uint64, p types.RewardsPlan) (bool, error) {
		plans = append(plans, p)
		return false, nil
	})

	if len(plans) == 0 {
		return types.RewardsPlan{}, false
	}

	return plans[r.Intn(len(plans))], true
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

func RandomDecPools(r *rand.Rand, availableDenoms []string) types.DecPools {
	pools := types.NewDecPools()

	// Pick a random subset of denoms
	denoms := simtesting.RandomSubSlice(r, availableDenoms)
	if len(denoms) == 0 {
		return pools
	}

	for _, denom := range denoms {
		// Generate a random amount
		amount := simtypes.RandomAmount(r, math.NewIntFromUint64(r.Uint64()))
		// Ignore if zero
		if amount.IsZero() {
			continue
		}

		// Create a DecPool with the random amount
		pool := types.NewDecPool(denom, sdk.NewDecCoins(
			sdk.NewDecCoin(denom, amount),
		))
		pools = pools.Add(pool)
	}

	return pools
}

func RandomServicePools(
	r *rand.Rand,
	servicesGenesis servicestypes.GenesisState,
	availableDenoms []string,
) types.ServicePools {
	var servicePools types.ServicePools

	services := simtesting.RandomSubSlice(r, servicesGenesis.Services)
	for _, service := range services {
		servicePools = append(servicePools, types.ServicePool{
			ServiceID: service.ID,
			DecPools:  RandomDecPools(r, availableDenoms),
		})
	}

	return servicePools
}

func RandomOutstandingRewardsRecords(r *rand.Rand, availableDenoms []string) []types.OutstandingRewardsRecord {
	var outstandingRewardsRecords []types.OutstandingRewardsRecord

	count := r.Intn(10)
	for i := 0; i < count; i++ {
		// Pick a random subset of the available denoms
		denoms := simtesting.RandomSubSlice(r, availableDenoms)
		if len(denoms) == 0 {
			continue
		}

		outstandingRewards := RandomDecPools(r, availableDenoms)
		// Ignore empty outstanding rewards
		if outstandingRewards.IsEmpty() {
			continue
		}

		outstandingRewardsRecords = append(outstandingRewardsRecords, types.OutstandingRewardsRecord{
			DelegationTargetID: simtesting.RandomPositiveUint32(r),
			OutstandingRewards: outstandingRewards,
		})
	}

	return outstandingRewardsRecords
}

func RandomHistoricalRewardsRecords(
	r *rand.Rand,
	servicesGenesis servicestypes.GenesisState,
	availableDenoms []string,
) []types.HistoricalRewardsRecord {
	var historicalRewardsRecords []types.HistoricalRewardsRecord

	count := r.Intn(10)
	for i := 0; i < count; i++ {
		servicePools := RandomServicePools(r, servicesGenesis, availableDenoms)
		// Ignore empty service pools
		if len(servicePools) == 0 {
			continue
		}

		historicalRewardsRecords = append(historicalRewardsRecords, types.HistoricalRewardsRecord{
			DelegationTargetID: simtesting.RandomPositiveUint32(r),
			Period:             r.Uint64(),
			Rewards: types.HistoricalRewards{
				CumulativeRewardRatios: servicePools,
				ReferenceCount:         r.Uint32(),
			},
		})
	}

	return historicalRewardsRecords
}

func RandomCurrentRewardsRecords(
	r *rand.Rand,
	servicesGenesis servicestypes.GenesisState,
	availableDenoms []string,
) []types.CurrentRewardsRecord {
	var currentRewardsRecords []types.CurrentRewardsRecord

	count := r.Intn(10)
	for i := 0; i < count; i++ {
		currentRewards := types.CurrentRewards{
			Rewards: RandomServicePools(r, servicesGenesis, availableDenoms),
			Period:  r.Uint64(),
		}
		// Ignore CurrentRewards if empty
		if len(currentRewards.Rewards) == 0 {
			continue
		}

		currentRewardsRecords = append(currentRewardsRecords, types.CurrentRewardsRecord{
			DelegationTargetID: simtesting.RandomPositiveUint32(r),
			Rewards:            currentRewards,
		})
	}

	return currentRewardsRecords
}

func RandomDelegatorStartingInfoRecords(
	r *rand.Rand,
	availableDenoms []string,
) []types.DelegatorStartingInfoRecord {
	var delegatorStartingInfoRecords []types.DelegatorStartingInfoRecord

	accounts := simtypes.RandomAccounts(r, r.Intn(10))
	for _, account := range accounts {
		record := types.DelegatorStartingInfoRecord{
			DelegatorAddress:   account.Address.String(),
			DelegationTargetID: simtesting.RandomPositiveUint32(r),
			StartingInfo: types.DelegatorStartingInfo{
				PreviousPeriod: simtesting.RandomPositiveUint64(r),
				Stakes:         simtesting.RandomDecCoins(r, availableDenoms, math.LegacyNewDec(r.Int63())),
				Height:         simtesting.RandomPositiveUint64(r),
			},
		}

		delegatorStartingInfoRecords = append(delegatorStartingInfoRecords, record)
	}

	return delegatorStartingInfoRecords
}

func RandomOperatorAccumulatedCommissionRecords(r *rand.Rand, availableDenoms []string) []types.OperatorAccumulatedCommissionRecord {
	var records []types.OperatorAccumulatedCommissionRecord

	count := r.Intn(10)
	for i := 0; i < count; i++ {
		randomDenoms := simtesting.RandomSubSlice(r, availableDenoms)
		if len(randomDenoms) == 0 {
			continue
		}

		records = append(records, types.OperatorAccumulatedCommissionRecord{
			OperatorID: simtesting.RandomPositiveUint32(r),
			Accumulated: types.AccumulatedCommission{
				Commissions: RandomDecPools(r, availableDenoms),
			},
		})
	}

	return records
}

func RandomPoolServiceTotalDelegatorShares(
	r *rand.Rand,
	poolsGenesis poolstypes.GenesisState,
	servicesGenesis servicestypes.GenesisState,
	availableDenoms []string,
) []types.PoolServiceTotalDelegatorShares {
	var records []types.PoolServiceTotalDelegatorShares

	services := simtesting.RandomSubSlice(r, servicesGenesis.Services)
	pools := simtesting.RandomSubSlice(r, poolsGenesis.Pools)
	for _, service := range services {
		for _, pool := range pools {
			randomDenom := availableDenoms[r.Intn(len(availableDenoms))]
			decCoin := sdk.NewDecCoinFromDec(randomDenom, simtypes.RandomDecAmount(
				r, math.LegacyNewDecFromInt(math.NewIntFromUint64(simtesting.RandomPositiveUint64(r))),
			))

			records = append(records, types.PoolServiceTotalDelegatorShares{
				PoolID:    pool.ID,
				ServiceID: service.ID,
				Shares:    sdk.NewDecCoins(decCoin),
			})
		}
	}

	return records
}
