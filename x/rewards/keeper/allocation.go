package keeper

import (
	"context"
	"fmt"
	"slices"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

func (k *Keeper) AllocateRewards(ctx context.Context) error {
	lastAllocationTime, err := k.GetLastRewardsAllocationTime(ctx)
	if err != nil {
		return err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if err := k.SetLastRewardsAllocationTime(ctx, sdkCtx.BlockTime()); err != nil {
		return err
	}

	if lastAllocationTime == nil {
		// If there's no last rewards allocation time set yet, it means this is
		// the first time AllocateRewards is called. In this case we just set
		// the current block time as the last rewards allocation time and skip
		// this block.
		return nil
	}

	timeSinceLastAllocation := sdkCtx.BlockTime().Sub(*lastAllocationTime)
	// TODO: clip elapsed time to prevent too much rewards allocation after
	//       possible chain halt?
	if timeSinceLastAllocation == 0 {
		return nil
	}

	err = k.RewardsPlans.Walk(ctx, nil, func(planID uint64, plan types.RewardsPlan) (stop bool, err error) {
		// Skip if the plan is not active at the current block time.
		if !plan.IsActiveAt(sdkCtx.BlockTime()) {
			return false, nil
		}

		err = k.AllocateRewardsByPlan(ctx, plan, timeSinceLastAllocation)
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

func (k *Keeper) AllocateRewardsByPlan(
	ctx context.Context, plan types.RewardsPlan, timeSinceLastAllocation time.Duration) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	rewards := sdk.NewDecCoinsFromCoins(plan.AmountPerDay...).
		MulDecTruncate(math.LegacyNewDec(timeSinceLastAllocation.Milliseconds())).
		QuoDecTruncate(math.LegacyNewDec((24 * time.Hour).Milliseconds()))

	// Check if the rewards pool has enough coins to allocate rewards.
	rewardsPoolAddr := plan.MustGetRewardsPoolAddress()
	balances := k.bankKeeper.GetAllBalances(ctx, rewardsPoolAddr)
	_, hasNeg := sdk.NewDecCoinsFromCoins(balances...).SafeSub(rewards)
	if hasNeg {
		sdkCtx.Logger().Info(
			"Skipping rewards plan because its rewards pool has insufficient balances",
			"plan_id", plan.ID,
			"balances", balances.String(),
			"rewards", rewards.String(),
		)
		return nil
	}

	service, found := k.servicesKeeper.GetService(sdkCtx, plan.ServiceID)
	if !found {
		return servicestypes.ErrServiceNotFound
	}

	serviceParams := k.restakingKeeper.GetServiceParams(sdkCtx, service.ID)
	pools := k.getPoolsForRewardsAllocation(ctx, service, serviceParams)
	operators := k.getOperatorsForRewardsAllocation(ctx, service, serviceParams)

	poolDistrInfos, totalPoolsDelValues, err := k.getPoolDistrInfos(ctx, pools)
	if err != nil {
		return err
	}

	operatorDistrInfos, totalOperatorsDelValues, err := k.getOperatorDistrInfos(ctx, operators)
	if err != nil {
		return err
	}

	totalUsersDelValues, err := k.GetCoinsValue(ctx, service.Tokens)
	if err != nil {
		return err
	}

	totalDelValues := totalPoolsDelValues.Add(totalOperatorsDelValues).Add(totalUsersDelValues)
	if totalDelValues.IsZero() {
		// There's no delegations at all, so just skip.
		return nil
	}

	var poolsRewards, operatorsRewards, usersRewards sdk.DecCoins

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
		err = k.allocateRewardsToPools(ctx, plan.PoolsDistribution, poolDistrInfos, poolsRewards)
		if err != nil {
			return err
		}
	}
	if operatorsRewards.IsAllPositive() {
		err = k.allocateRewardsToOperators(ctx, plan.OperatorsDistribution, operatorDistrInfos, operatorsRewards)
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

func (k *Keeper) getPoolsForRewardsAllocation(
	ctx context.Context, service servicestypes.Service,
	serviceParams restakingtypes.ServiceParams) []poolstypes.Pool {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	var pools []poolstypes.Pool
	poolsParams := k.poolsKeeper.GetParams(sdkCtx)
	if slices.Contains(poolsParams.AllowedServiceIDs, service.ID) {
		// If there's no whitelisted pools, that means all pools.
		if len(serviceParams.WhitelistedPoolIDs) == 0 {
			return k.poolsKeeper.GetPools(sdkCtx)
		}
		for _, poolID := range serviceParams.WhitelistedPoolIDs {
			pool, found := k.poolsKeeper.GetPool(sdkCtx, poolID)
			if !found {
				// TODO: panic here if we're sure that this never happens
				k.Logger(ctx).Warn("whitelisted pool not found", "pool_id", poolID)
				continue
			}
			pools = append(pools, pool)
		}
	}
	return pools
}

func (k *Keeper) getOperatorsForRewardsAllocation(
	ctx context.Context, service servicestypes.Service,
	serviceParams restakingtypes.ServiceParams) []operatorstypes.Operator {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	// TODO: can we optimize this? maybe by having a new index key
	var operators []operatorstypes.Operator
	k.operatorsKeeper.IterateOperators(sdkCtx, func(operator operatorstypes.Operator) (stop bool) {
		operatorParams := k.restakingKeeper.GetOperatorParams(sdkCtx, operator.ID)
		if slices.Contains(operatorParams.JoinedServiceIDs, service.ID) &&
			(len(serviceParams.WhitelistedOperatorIDs) == 0 ||
				slices.Contains(serviceParams.WhitelistedOperatorIDs, operator.ID)) {
			operators = append(operators, operator)
		}
		return false
	})
	return operators
}

func (k *Keeper) getPoolDistrInfos(
	ctx context.Context, pools []poolstypes.Pool,
) (distrInfos []PoolDistributionInfo, totalDelValues math.LegacyDec, err error) {
	totalDelValues = math.LegacyZeroDec()
	for _, pool := range pools {
		delValue, err := k.GetCoinValue(ctx, sdk.NewCoin(pool.Denom, pool.Tokens))
		if err != nil {
			return nil, math.LegacyDec{}, err
		}
		if delValue.IsZero() {
			continue
		}
		distrInfos = append(distrInfos, PoolDistributionInfo{
			Pool:             &pool,
			DelegationsValue: delValue,
		})
		totalDelValues = totalDelValues.Add(delValue)
	}
	return distrInfos, totalDelValues, nil
}

func (k *Keeper) getOperatorDistrInfos(
	ctx context.Context, operators []operatorstypes.Operator,
) (distrInfos []OperatorDistributionInfo, totalDelValues math.LegacyDec, err error) {
	totalDelValues = math.LegacyZeroDec()
	for _, operator := range operators {
		delValue, err := k.GetCoinsValue(ctx, operator.Tokens)
		if err != nil {
			return nil, math.LegacyDec{}, err
		}
		if delValue.IsZero() {
			continue
		}
		distrInfos = append(distrInfos, OperatorDistributionInfo{
			Operator:         &operator,
			DelegationsValue: delValue,
		})
		totalDelValues = totalDelValues.Add(delValue)
	}
	return distrInfos, totalDelValues, nil
}

func (k *Keeper) allocateRewardsToPools(
	ctx context.Context, distr types.PoolsDistribution, poolDistrInfos []PoolDistributionInfo,
	rewards sdk.DecCoins) error {
	var poolsDistrType types.PoolsDistributionType
	err := k.cdc.UnpackAny(distr.Type, &poolsDistrType)
	if err != nil {
		return err
	}
	switch typ := poolsDistrType.(type) {
	case *types.PoolsDistributionTypeBasic:
		return k.allocateRewardsToPoolsBasic(ctx, poolDistrInfos, rewards)
	case *types.PoolsDistributionTypeWeighted:
		return k.allocateRewardsToPoolsWeighted(ctx, poolDistrInfos, rewards, typ.Weights)
	case *types.PoolsDistributionTypeEgalitarian:
		return k.allocateRewardsToPoolsEgalitarian(ctx, poolDistrInfos, rewards)
	default:
		panic("unknown pools distribution type")
	}
}

func (k *Keeper) allocateRewardsToPoolsBasic(
	ctx context.Context, distrInfos []PoolDistributionInfo, rewards sdk.DecCoins) error {
	totalDelValues := math.LegacyZeroDec()
	for _, distrInfo := range distrInfos {
		totalDelValues = totalDelValues.Add(distrInfo.DelegationsValue)
	}
	for _, distrInfo := range distrInfos {
		poolRewards := rewards.MulDecTruncate(distrInfo.DelegationsValue).QuoDecTruncate(totalDelValues)
		err := k.allocateRewardsPool(ctx, types.NewDelegationTarget(distrInfo.Pool), distrInfo.Pool.Denom, poolRewards)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *Keeper) allocateRewardsToPoolsWeighted(
	ctx context.Context, distrInfos []PoolDistributionInfo, rewards sdk.DecCoins,
	weights []types.PoolDistributionWeight) error {
	distrInfoByPoolID := map[uint32]PoolDistributionInfo{}
	for _, distrInfo := range distrInfos {
		distrInfoByPoolID[distrInfo.Pool.ID] = distrInfo
	}

	totalWeights := math.LegacyZeroDec()
	for _, weight := range weights {
		if _, ok := distrInfoByPoolID[weight.PoolID]; !ok {
			// If there's no distrInfo for specified pool, skip it.
			continue
		}
		totalWeights = totalWeights.Add(math.LegacyNewDec(int64(weight.Weight)))
	}

	for _, weight := range weights {
		distrInfo, ok := distrInfoByPoolID[weight.PoolID]
		if !ok {
			continue
		}

		poolRewards := rewards.MulDecTruncate(math.LegacyNewDec(int64(weight.Weight))).QuoDecTruncate(totalWeights)
		err := k.allocateRewardsPool(ctx, types.NewDelegationTarget(distrInfo.Pool), distrInfo.Pool.Denom, poolRewards)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *Keeper) allocateRewardsToPoolsEgalitarian(
	ctx context.Context, distrInfos []PoolDistributionInfo, rewards sdk.DecCoins) error {
	numPools := math.LegacyNewDec(int64(len(distrInfos)))
	for _, distrInfo := range distrInfos {
		poolRewards := rewards.QuoDecTruncate(numPools)
		err := k.allocateRewardsPool(ctx, types.NewDelegationTarget(distrInfo.Pool), distrInfo.Pool.Denom, poolRewards)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *Keeper) allocateRewardsToOperators(
	ctx context.Context, distr types.OperatorsDistribution, distrInfos []OperatorDistributionInfo,
	rewards sdk.DecCoins) error {
	var operatorsDistrType types.OperatorsDistributionType
	err := k.cdc.UnpackAny(distr.Type, &operatorsDistrType)
	if err != nil {
		return err
	}
	switch typ := operatorsDistrType.(type) {
	case *types.OperatorsDistributionTypeBasic:
		return k.allocateRewardsToOperatorsBasic(ctx, distrInfos, rewards)
	case *types.OperatorsDistributionTypeWeighted:
		return k.allocateRewardsToOperatorsWeighted(ctx, distrInfos, rewards, typ.Weights)
	case *types.OperatorsDistributionTypeEgalitarian:
		return k.allocateRewardsToOperatorsEgalitarian(ctx, distrInfos, rewards)
	default:
		panic("unknown operators distribution type")
	}
}

func (k *Keeper) allocateRewardsToOperatorsBasic(
	ctx context.Context, distrInfos []OperatorDistributionInfo, rewards sdk.DecCoins) error {
	totalDelValues := math.LegacyZeroDec()
	for _, distrInfo := range distrInfos {
		totalDelValues = totalDelValues.Add(distrInfo.DelegationsValue)
	}
	for _, distrInfo := range distrInfos {
		operatorRewards := rewards.MulDecTruncate(distrInfo.DelegationsValue).QuoDecTruncate(totalDelValues)
		err := k.allocateRewardsToOperator(ctx, distrInfo, operatorRewards)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *Keeper) allocateRewardsToOperatorsWeighted(
	ctx context.Context, distrInfos []OperatorDistributionInfo, rewards sdk.DecCoins,
	weights []types.OperatorDistributionWeight) error {
	distrInfoByOperatorID := map[uint32]OperatorDistributionInfo{}
	for _, distrInfo := range distrInfos {
		distrInfoByOperatorID[distrInfo.Operator.ID] = distrInfo
	}

	totalWeights := math.LegacyZeroDec()
	for _, weight := range weights {
		if _, ok := distrInfoByOperatorID[weight.OperatorID]; !ok {
			// If there's no distrInfo for specified operator, skip it.
			continue
		}
		totalWeights = totalWeights.Add(math.LegacyNewDec(int64(weight.Weight)))
	}

	for _, weight := range weights {
		distrInfo, ok := distrInfoByOperatorID[weight.OperatorID]
		if !ok {
			continue
		}

		operatorRewards := rewards.MulDecTruncate(math.LegacyNewDec(int64(weight.Weight))).QuoDecTruncate(totalWeights)
		err := k.allocateRewardsToOperator(ctx, distrInfo, operatorRewards)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *Keeper) allocateRewardsToOperatorsEgalitarian(
	ctx context.Context, distrInfos []OperatorDistributionInfo, rewards sdk.DecCoins) error {
	numOperators := math.LegacyNewDec(int64(len(distrInfos)))
	for _, distrInfo := range distrInfos {
		operatorRewards := rewards.QuoDecTruncate(numOperators)
		err := k.allocateRewardsToOperator(ctx, distrInfo, operatorRewards)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *Keeper) allocateRewardsToOperator(
	ctx context.Context, distrInfo OperatorDistributionInfo, rewards sdk.DecCoins) error {
	for _, token := range distrInfo.Operator.Tokens {
		tokenValue, err := k.GetCoinValue(ctx, token)
		if err != nil {
			return err
		}
		if tokenValue.IsZero() {
			continue
		}
		tokenRewards := rewards.MulDecTruncate(tokenValue).QuoDecTruncate(distrInfo.DelegationsValue)
		err = k.allocateRewardsPool(ctx, types.NewDelegationTarget(distrInfo.Operator), token.Denom, tokenRewards)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *Keeper) allocateRewardsPool(
	ctx context.Context, target *types.DelegationTarget, denom string, rewards sdk.DecCoins) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	shared := rewards
	if target.Type() == restakingtypes.DELEGATION_TYPE_OPERATOR {
		// split tokens between operator and delegators according to commission
		// TODO: optimize this read operation? we already read operator params in
		//       getOperatorsForRewardsAllocation
		operatorParams := k.restakingKeeper.GetOperatorParams(sdkCtx, target.GetID())
		commission := rewards.MulDec(operatorParams.CommissionRate)
		shared = rewards.Sub(commission)

		// update current commission
		sdkCtx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCommission,
				sdk.NewAttribute(sdk.AttributeKeyAmount, commission.String()),
				sdk.NewAttribute(types.AttributeKeyDelegationType, restakingtypes.DELEGATION_TYPE_OPERATOR.String()),
				sdk.NewAttribute(types.AttributeKeyDelegationTargetID, fmt.Sprint(target.GetID())),
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

	currentRewards, err := k.GetCurrentRewards(ctx, target)
	if err != nil {
		return err
	}
	currentRewards.Rewards = currentRewards.Rewards.Add(types.NewDecPool(denom, shared))
	err = k.SetCurrentRewards(ctx, target, currentRewards)
	if err != nil {
		return err
	}

	outstanding, err := k.GetOutstandingRewards(ctx, target)
	if err != nil {
		return err
	}
	outstanding.Rewards = outstanding.Rewards.Add(types.NewDecPool(denom, rewards))
	err = k.SetOutstandingRewards(ctx, target, outstanding)
	if err != nil {
		return err
	}

	// update outstanding rewards
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRewards,
			sdk.NewAttribute(types.AttributeKeyDelegationType, target.Type().String()),
			sdk.NewAttribute(types.AttributeKeyDelegationTargetID, fmt.Sprint(target.GetID())),
			sdk.NewAttribute(types.AttributeKeyPool, denom),
			sdk.NewAttribute(sdk.AttributeKeyAmount, rewards.String()),
		),
	)

	return nil
}

func (k *Keeper) allocateRewardsToUsers(
	ctx context.Context, distr types.UsersDistribution, service servicestypes.Service, totalDelValues math.LegacyDec,
	rewards sdk.DecCoins) error {
	var usersDistrType types.UsersDistributionType
	err := k.cdc.UnpackAny(distr.Type, &usersDistrType)
	if err != nil {
		return err
	}
	switch usersDistrType.(type) {
	case *types.UsersDistributionTypeBasic:
		return k.allocateRewardsToService(ctx, service, totalDelValues, rewards)
	default:
		panic("unknown operators distribution type")
	}
}

func (k *Keeper) allocateRewardsToService(
	ctx context.Context, service servicestypes.Service, totalDelValues math.LegacyDec,
	rewards sdk.DecCoins) error {
	for _, token := range service.Tokens {
		tokenValue, err := k.GetCoinValue(ctx, token)
		if err != nil {
			return err
		}
		if tokenValue.IsZero() {
			continue
		}
		tokenRewards := rewards.MulDecTruncate(tokenValue).QuoDecTruncate(totalDelValues)
		err = k.allocateRewardsPool(ctx, types.NewDelegationTarget(&service), token.Denom, tokenRewards)
		if err != nil {
			return err
		}
	}
	return nil
}

type PoolDistributionInfo struct {
	Pool             *poolstypes.Pool
	DelegationsValue math.LegacyDec
}

type OperatorDistributionInfo struct {
	Operator         *operatorstypes.Operator
	DelegationsValue math.LegacyDec
}