package keeper

import (
	"context"
	"time"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

func (k *Keeper) CreateRewardsPlan(
	ctx context.Context, description string, serviceID uint32, amt sdk.Coins, startTime, endTime time.Time,
	poolsDistribution types.PoolsDistribution, operatorsDistribution types.OperatorsDistribution,
	usersDistribution types.UsersDistribution,
) (types.RewardsPlan, error) {
	// TODO: validate arguments.

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	_, found := k.servicesKeeper.GetService(sdkCtx, serviceID)
	if !found {
		return types.RewardsPlan{}, servicestypes.ErrServiceNotFound
	}

	err := k.validatePoolsDistribution(ctx, poolsDistribution)
	if err != nil {
		return types.RewardsPlan{}, err
	}
	err = k.validateOperatorsDistribution(ctx, operatorsDistribution)
	if err != nil {
		return types.RewardsPlan{}, err
	}

	// Get the next plan ID and increment it by 1
	planID, err := k.NextRewardsPlanID.Get(ctx)
	if err != nil {
		return types.RewardsPlan{}, err
	}
	if err := k.NextRewardsPlanID.Set(ctx, planID+1); err != nil {
		return types.RewardsPlan{}, err
	}

	plan := types.NewRewardsPlan(
		planID, description, serviceID, amt, startTime, endTime,
		poolsDistribution, operatorsDistribution,
		usersDistribution)
	if err := plan.Validate(k.cdc); err != nil {
		return types.RewardsPlan{}, err
	}

	if err := k.RewardsPlans.Set(ctx, planID, plan); err != nil {
		return types.RewardsPlan{}, err
	}

	k.createAccountIfNotExists(ctx, plan.MustGetRewardsPoolAddress())

	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateRewardsPlan,
			// TODO: add attributes
		),
	})

	return plan, nil
}

func (k *Keeper) validatePoolsDistribution(ctx context.Context, distribution types.PoolsDistribution) error {
	var distrType types.PoolsDistributionType
	err := k.cdc.UnpackAny(distribution.Type, &distrType)
	if err != nil {
		return err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	typ, ok := distrType.(*types.PoolsDistributionTypeWeighted)
	if !ok {
		// Only weighted distribution needs a validation.
		return nil
	}
	for _, weight := range typ.Weights {
		_, found := k.poolsKeeper.GetPool(sdkCtx, weight.PoolID)
		if !found {
			return errors.Wrapf(poolstypes.ErrPoolNotFound, "pool %d not found", weight.PoolID)
		}
	}
	return nil
}

func (k *Keeper) validateOperatorsDistribution(ctx context.Context, distribution types.OperatorsDistribution) error {
	var distrType types.OperatorsDistributionType
	err := k.cdc.UnpackAny(distribution.Type, &distrType)
	if err != nil {
		return err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	typ, ok := distrType.(*types.OperatorsDistributionTypeWeighted)
	if !ok {
		// Only weighted distribution needs a validation.
		return nil
	}
	for _, weight := range typ.Weights {
		_, found := k.operatorsKeeper.GetOperator(sdkCtx, weight.OperatorID)
		if !found {
			return errors.Wrapf(operatorstypes.ErrOperatorNotFound, "operator %d not found", weight.OperatorID)
		}
	}
	return nil
}
