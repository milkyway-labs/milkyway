package keeper

import (
	"context"
	"fmt"
	"slices"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gogotypes "github.com/cosmos/gogoproto/types"

	"github.com/milkyway-labs/milkyway/utils"
	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

// DistributionInfo stores information about a delegation target and its
// delegation value.
type DistributionInfo struct {
	DelegationTarget restakingtypes.DelegationTarget
	DelegationsValue math.LegacyDec
}

// GetLastRewardsAllocationTime returns the last time rewards were allocated.
// If there's no last rewards allocation time set yet, nil is returned.
func (k *Keeper) GetLastRewardsAllocationTime(ctx context.Context) (*time.Time, error) {
	ts, err := k.LastRewardsAllocationTime.Get(ctx)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	t, err := gogotypes.TimestampFromProto(&ts)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// SetLastRewardsAllocationTime sets the last time rewards were allocated.
func (k *Keeper) SetLastRewardsAllocationTime(ctx context.Context, t time.Time) error {
	ts, err := gogotypes.TimestampProto(t)
	if err != nil {
		return err
	}
	return k.LastRewardsAllocationTime.Set(ctx, *ts)
}

// AllocateRewards allocates restaking rewards to different entities based on
// active rewards plans. AllocateRewards skips rewards distribution when
// there's no last rewards allocation time set. In that case, AllocateRewards
// simply set the current block time as new last rewards allocation time.
func (k *Keeper) AllocateRewards(ctx context.Context) error {
	// Get last rewards allocation time and set the current block time as new
	// last rewards allocation time.
	lastAllocationTime, err := k.GetLastRewardsAllocationTime(ctx)
	if err != nil {
		return err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if err := k.SetLastRewardsAllocationTime(ctx, sdkCtx.BlockTime()); err != nil {
		return err
	}

	// If there's no last rewards allocation time set yet, it means this is
	// the first time AllocateRewards is called. In this case we just skip this
	// block for rewards allocation.
	if lastAllocationTime == nil {
		return nil
	}

	// Calculate time elapsed since the last rewards allocation to calculate
	// rewards amount to allocate in this block.
	timeSinceLastAllocation := sdkCtx.BlockTime().Sub(*lastAllocationTime)

	// TODO: clip elapsed time to prevent too much rewards allocation after possible chain halt?
	if timeSinceLastAllocation == 0 {
		return nil
	}

	pools := k.poolsKeeper.GetPools(sdkCtx)
	operators := k.operatorsKeeper.GetOperators(sdkCtx)

	// Iterate all rewards plan stored and allocate rewards by plan if it's
	// active(plan's start time <= current block time < plans' end time).
	err = k.RewardsPlans.Walk(ctx, nil, func(planID uint64, plan types.RewardsPlan) (stop bool, err error) {
		// Skip if the plan is not active at the current block time.
		if !plan.IsActiveAt(sdkCtx.BlockTime()) {
			return false, nil
		}

		err = k.AllocateRewardsByPlan(ctx, plan, timeSinceLastAllocation, pools, operators)
		if err != nil {
			return false, err
		}
		return false, nil
	})
	if err != nil {
		return err
	}
	return nil
}

// AllocateRewardsByPlan allocates rewards by a specific rewards plan.
func (k *Keeper) AllocateRewardsByPlan(
	ctx context.Context, plan types.RewardsPlan, timeSinceLastAllocation time.Duration,
	pools []poolstypes.Pool, operators []operatorstypes.Operator,
) error {
	// Calculate rewards amount for this block by following formula:
	// amountPerDay * timeSinceLastAllocation(ms) / 1 day(ms)
	rewards := sdk.NewDecCoinsFromCoins(plan.AmountPerDay...).
		MulDecTruncate(math.LegacyNewDec(timeSinceLastAllocation.Milliseconds())).
		QuoDecTruncate(math.LegacyNewDec((24 * time.Hour).Milliseconds()))

	// Truncate decimal and move the truncated rewards to the global rewards
	// pool.
	rewardsTruncated, _ := rewards.TruncateDecimal()

	// Use this truncated rewards so that we don't allocate more rewards than
	// what have been moved to the global rewards pool.
	rewards = sdk.NewDecCoinsFromCoins(rewardsTruncated...)

	// Check if the rewards pool has enough coins to allocate rewards.
	planRewardsPoolAddr := plan.MustGetRewardsPoolAddress(k.accountKeeper.AddressCodec())
	balances := k.bankKeeper.GetAllBalances(ctx, planRewardsPoolAddr)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if !balances.IsAllGTE(rewardsTruncated) {
		sdkCtx.Logger().Info(
			"Skipping rewards plan because its rewards pool has insufficient balances",
			"plan_id", plan.ID,
			"balances", balances.String(),
			"rewards", rewards.String(),
		)
		return nil
	}
	// Send the current block's rewards to the global rewards pool.
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, planRewardsPoolAddr, types.RewardsPoolName, rewardsTruncated)
	if err != nil {
		return err
	}

	service, found := k.servicesKeeper.GetService(sdkCtx, plan.ServiceID)
	if !found {
		return servicestypes.ErrServiceNotFound
	}

	// Get pools and operators that are eligible for this rewards plan's
	// rewards allocation, along with each entity's total delegation value.
	// If an entity's total delegation value is zero, then it won't be included
	// in this block's rewards allocation.
	serviceParams := k.restakingKeeper.GetServiceParams(sdkCtx, service.ID)

	eligiblePools := k.getEligiblePools(ctx, service, serviceParams, pools)
	poolDistrInfos, totalPoolsDelValues, err := k.getDistrInfos(ctx, eligiblePools)
	if err != nil {
		return err
	}

	eligibleOperators := k.getEligibleOperators(ctx, service, serviceParams, operators)
	operatorDistrInfos, totalOperatorsDelValues, err := k.getDistrInfos(ctx, eligibleOperators)
	if err != nil {
		return err
	}

	totalUsersDelValues, err := k.GetCoinsValue(ctx, service.Tokens)
	if err != nil {
		return err
	}

	totalDelValues := totalPoolsDelValues.Add(totalOperatorsDelValues).Add(totalUsersDelValues)

	// There's no delegations at all, so just skip.
	if totalDelValues.IsZero() {
		return nil
	}

	// Only sum weights with non-zero delegation values.
	totalWeightsNonZeroDelValues := uint32(0)
	if totalPoolsDelValues.IsPositive() {
		totalWeightsNonZeroDelValues += plan.PoolsDistribution.Weight
	}
	if totalOperatorsDelValues.IsPositive() {
		totalWeightsNonZeroDelValues += plan.OperatorsDistribution.Weight
	}
	if totalUsersDelValues.IsPositive() {
		totalWeightsNonZeroDelValues += plan.UsersDistribution.Weight
	}

	var poolsRewards, operatorsRewards, usersRewards sdk.DecCoins
	if totalWeightsNonZeroDelValues > 0 {
		// If weights are specified, then split rewards by
		// rewards * weight / totalWeights
		totalWeightsDec := math.LegacyNewDec(int64(totalWeightsNonZeroDelValues))
		poolsRewards = rewards.MulDecTruncate(math.LegacyNewDec(int64(plan.PoolsDistribution.Weight))).
			QuoDecTruncate(totalWeightsDec)
		operatorsRewards = rewards.MulDecTruncate(math.LegacyNewDec(int64(plan.OperatorsDistribution.Weight))).
			QuoDecTruncate(totalWeightsDec)
		usersRewards = rewards.MulDecTruncate(math.LegacyNewDec(int64(plan.UsersDistribution.Weight))).
			QuoDecTruncate(totalWeightsDec)
	} else {
		// If there's no weights specified, then distribute rewards based on their
		// total delegation values.
		poolsRewards = rewards.MulDecTruncate(totalPoolsDelValues).QuoDecTruncate(totalDelValues)
		operatorsRewards = rewards.MulDecTruncate(totalOperatorsDelValues).QuoDecTruncate(totalDelValues)
		usersRewards = rewards.MulDecTruncate(totalUsersDelValues).QuoDecTruncate(totalDelValues)
	}

	if poolsRewards.IsAllPositive() {
		err = k.allocateRewards(ctx, plan.PoolsDistribution, poolDistrInfos, poolsRewards)
		if err != nil {
			return err
		}
	}
	if operatorsRewards.IsAllPositive() {
		err = k.allocateRewards(ctx, plan.OperatorsDistribution, operatorDistrInfos, operatorsRewards)
		if err != nil {
			return err
		}
	}
	if usersRewards.IsAllPositive() {
		err = k.allocateRewardsToUsers(ctx, plan.UsersDistribution, service, totalUsersDelValues, usersRewards)
		if err != nil {
			return err
		}
	}
	return nil
}

// getEligiblePools returns a list of pools that are eligible for rewards
// allocation based on the given service's service params. If the service's
// service params don't have any whitelisted pools, then all pools are eligible
// for rewards allocation. Also, if the service is not whitelisted in the pools
// params, then no pools are eligible for rewards allocation.
func (k *Keeper) getEligiblePools(
	ctx context.Context, service servicestypes.Service,
	serviceParams restakingtypes.ServiceParams, pools []poolstypes.Pool,
) (eligiblePools []restakingtypes.DelegationTarget) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	poolsParams := k.poolsKeeper.GetParams(sdkCtx)
	if slices.Contains(poolsParams.AllowedServicesIDs, service.ID) {
		// If there's no whitelisted pools, that means all pools.
		if len(serviceParams.WhitelistedPoolsIDs) == 0 {
			return utils.Map(pools, func(p poolstypes.Pool) restakingtypes.DelegationTarget {
				return &p
			})
		}
		whitelistedPoolIDs := map[uint32]struct{}{}
		for _, poolID := range serviceParams.WhitelistedPoolsIDs {
			whitelistedPoolIDs[poolID] = struct{}{}
		}
		// Only include whitelisted pools.
		for _, pool := range pools {
			if _, ok := whitelistedPoolIDs[pool.ID]; ok {
				eligiblePools = append(eligiblePools, &pool)
			}
		}
	}
	return eligiblePools
}

// getEligibleOperators returns a list of operators that are eligible for rewards
// allocation based on the given service's service params. If the service's
// service params don't have any whitelisted operators, then all operators that
// have joined the service are eligible for rewards allocation.
func (k *Keeper) getEligibleOperators(
	ctx context.Context, service servicestypes.Service,
	serviceParams restakingtypes.ServiceParams, operators []operatorstypes.Operator,
) (eligibleOperators []restakingtypes.DelegationTarget) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	// TODO: can we optimize this? maybe by having a new index key
	for _, operator := range operators {
		operatorParams := k.restakingKeeper.GetOperatorParams(sdkCtx, operator.ID)
		if slices.Contains(operatorParams.JoinedServicesIDs, service.ID) &&
			(len(serviceParams.WhitelistedOperatorsIDs) == 0 ||
				slices.Contains(serviceParams.WhitelistedOperatorsIDs, operator.ID)) {
			eligibleOperators = append(eligibleOperators, &operator)
		}
	}
	return eligibleOperators
}

// getDistrInfos returns a list of DistributionInfo calculated based on each
// delegation target's delegation values. getDistrInfos also returns the total
// delegation values of all targets.
func (k *Keeper) getDistrInfos(
	ctx context.Context, targets []restakingtypes.DelegationTarget,
) (distrInfos []DistributionInfo, totalDelValues math.LegacyDec, err error) {
	totalDelValues = math.LegacyZeroDec()
	for _, target := range targets {
		delValue, err := k.GetCoinsValue(ctx, target.GetTokens())
		if err != nil {
			return nil, math.LegacyDec{}, err
		}
		// Skip if there's no delegations value. This can happen when there's
		// no tokens delegated at all or there's no price associated with the
		// delegated tokens.
		if delValue.IsZero() {
			continue
		}
		distrInfos = append(distrInfos, DistributionInfo{
			DelegationTarget: target,
			DelegationsValue: delValue,
		})
		totalDelValues = totalDelValues.Add(delValue)
	}
	return distrInfos, totalDelValues, nil
}

// allocateRewards allocates rewards to each delegation target based on the
// given distribution type. allocateRewards skips rewards allocation if there's
// no delegation targets specified in the distribution. If the distribution type
// is unknown, then an error is returned.
func (k *Keeper) allocateRewards(
	ctx context.Context, distr types.Distribution, distrInfos []DistributionInfo, rewards sdk.DecCoins,
) error {
	distrType, err := types.GetDistributionType(k.cdc, distr)
	if err != nil {
		return err
	}

	switch typ := distrType.(type) {
	case *types.DistributionTypeBasic:
		return k.allocateRewardsBasic(ctx, distrInfos, rewards)
	case *types.DistributionTypeWeighted:
		return k.allocateRewardsWeighted(ctx, distrInfos, rewards, typ.Weights)
	case *types.DistributionTypeEgalitarian:
		return k.allocateRewardsEgalitarian(ctx, distrInfos, rewards)
	default:
		return errors.Wrapf(sdkerrors.ErrInvalidType, "unknown distribution type: %T", typ)
	}
}

// allocateRewardsBasic allocates rewards to each delegation target based on
// their delegation values. Each delegation target receives rewards based on
// the following formula:
// targetRewards = rewards * targetDelegationsValue / totalDelegationsValue
func (k *Keeper) allocateRewardsBasic(
	ctx context.Context, distrInfos []DistributionInfo, rewards sdk.DecCoins,
) error {
	totalDelValues := math.LegacyZeroDec()
	for _, distrInfo := range distrInfos {
		totalDelValues = totalDelValues.Add(distrInfo.DelegationsValue)
	}
	for _, distrInfo := range distrInfos {
		targetRewards := rewards.MulDecTruncate(distrInfo.DelegationsValue).QuoDecTruncate(totalDelValues)
		err := k.allocateDelegationTargetRewards(ctx, distrInfo, targetRewards)
		if err != nil {
			return err
		}
	}
	return nil
}

// allocateRewardsWeighted allocates rewards to each delegation target based on
// their delegation values and weights. Each delegation target receives rewards
// based on the following formula:
// targetRewards = rewards * weight / totalWeights
func (k *Keeper) allocateRewardsWeighted(
	ctx context.Context, distrInfos []DistributionInfo, rewards sdk.DecCoins, weights []types.DistributionWeight,
) error {
	distrInfoByTargetID := map[uint32]DistributionInfo{}
	for _, distrInfo := range distrInfos {
		distrInfoByTargetID[distrInfo.DelegationTarget.GetID()] = distrInfo
	}

	totalWeights := math.LegacyZeroDec()
	for _, weight := range weights {
		if _, ok := distrInfoByTargetID[weight.DelegationTargetID]; !ok {
			// If there's no distrInfo for specified pool, skip it.
			continue
		}
		totalWeights = totalWeights.Add(math.LegacyNewDec(int64(weight.Weight)))
	}

	for _, weight := range weights {
		distrInfo, ok := distrInfoByTargetID[weight.DelegationTargetID]
		if !ok {
			continue
		}

		targetRewards := rewards.MulDecTruncate(math.LegacyNewDec(int64(weight.Weight))).QuoDecTruncate(totalWeights)
		err := k.allocateDelegationTargetRewards(ctx, distrInfo, targetRewards)
		if err != nil {
			return err
		}
	}
	return nil
}

// allocateRewardsEgalitarian allocates rewards to each delegation target equally.
// Each delegation target receives rewards based on the following formula:
// targetRewards = rewards / numTargets
func (k *Keeper) allocateRewardsEgalitarian(
	ctx context.Context, distrInfos []DistributionInfo, rewards sdk.DecCoins,
) error {
	numTargets := math.LegacyNewDec(int64(len(distrInfos)))
	for _, distrInfo := range distrInfos {
		targetRewards := rewards.QuoDecTruncate(numTargets)
		err := k.allocateDelegationTargetRewards(ctx, distrInfo, targetRewards)
		if err != nil {
			return err
		}
	}
	return nil
}

// allocateRewardsToUsers allocates rewards to users based on the given users
// distribution. allocateRewardsToUsers skips rewards allocation if there's no
// users specified in the distribution. If the users distribution type is unknown,
// then an error is returned.
func (k *Keeper) allocateRewardsToUsers(
	ctx context.Context,
	distr types.UsersDistribution,
	service servicestypes.Service,
	totalDelValues math.LegacyDec,
	rewards sdk.DecCoins,
) error {
	distrType, err := types.GetUsersDistributionType(k.cdc, distr)
	if err != nil {
		return err
	}

	switch distrType.(type) {
	case *types.UsersDistributionTypeBasic:
		return k.allocateDelegationTargetRewards(ctx, DistributionInfo{
			DelegationTarget: &service,
			DelegationsValue: totalDelValues,
		}, rewards)
	default:
		panic("unknown operators distribution type")
	}
}

// allocateDelegationTargetRewards allocates rewards to a specific delegation target.
func (k *Keeper) allocateDelegationTargetRewards(
	ctx context.Context, distrInfo DistributionInfo, rewards sdk.DecCoins,
) error {
	for _, token := range distrInfo.DelegationTarget.GetTokens() {
		tokenValue, err := k.GetCoinValue(ctx, token)
		if err != nil {
			return err
		}
		if tokenValue.IsZero() {
			continue
		}

		tokenRewards := rewards.MulDecTruncate(tokenValue).QuoDecTruncate(distrInfo.DelegationsValue)
		err = k.allocateRewardsPool(ctx, distrInfo.DelegationTarget, token.Denom, tokenRewards)
		if err != nil {
			return err
		}
	}
	return nil
}

// allocateRewardsPool allocates rewards to a specific delegation target's rewards pool.
func (k *Keeper) allocateRewardsPool(
	ctx context.Context, target restakingtypes.DelegationTarget, denom string, rewards sdk.DecCoins,
) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	shared := rewards
	if _, ok := target.(*operatorstypes.Operator); ok {
		// Split tokens between operator and delegators according to commission
		// TODO: optimize this read operation? we already read operator params in getEligibleOperators
		operatorParams, found, err := k.operatorsKeeper.GetOperatorParams(sdkCtx, target.GetID())
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("operator params for operator %d not found", target.GetID())
		}
		commission := rewards.MulDec(operatorParams.CommissionRate)
		shared = rewards.Sub(commission)

		// update current commission
		sdkCtx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCommission,
				sdk.NewAttribute(operatorstypes.AttributeKeyOperatorID, fmt.Sprint(target.GetID())),
				sdk.NewAttribute(types.AttributeKeyPool, denom),
				sdk.NewAttribute(sdk.AttributeKeyAmount, commission.String()),
			),
		)

		currentCommission, err := k.GetOperatorAccumulatedCommission(ctx, target.GetID())
		if err != nil {
			return err
		}
		currentCommission.Commissions = currentCommission.Commissions.Add(types.NewDecPool(denom, commission))
		err = k.OperatorAccumulatedCommissions.Set(ctx, target.GetID(), currentCommission)
		if err != nil {
			return err
		}
	}

	// Update current rewards
	currentRewards, err := k.GetCurrentRewards(ctx, target)
	if err != nil {
		return err
	}

	currentRewards.Rewards = currentRewards.Rewards.Add(types.NewDecPool(denom, shared))
	err = k.SetCurrentRewards(ctx, target, currentRewards)
	if err != nil {
		return err
	}

	// Update the outstanding rewards
	outstanding, err := k.GetOutstandingRewards(ctx, target)
	if err != nil {
		return err
	}

	outstanding.Rewards = outstanding.Rewards.Add(types.NewDecPool(denom, rewards))
	err = k.SetOutstandingRewards(ctx, target, outstanding)
	if err != nil {
		return err
	}

	// Emit the event
	delType, err := types.GetDelegationTargetType(target)
	if err != nil {
		return err
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRewards,
			sdk.NewAttribute(types.AttributeKeyDelegationType, delType.String()),
			sdk.NewAttribute(types.AttributeKeyDelegationTargetID, fmt.Sprint(target.GetID())),
			sdk.NewAttribute(types.AttributeKeyPool, denom),
			sdk.NewAttribute(sdk.AttributeKeyAmount, rewards.String()),
		),
	)

	return nil
}
