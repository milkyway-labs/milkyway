package keeper

import (
	"context"
	"time"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/rewards/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

func (k *Keeper) CreateRewardsPlan(
	ctx context.Context, description string, serviceID uint32, amt sdk.Coins, startTime, endTime time.Time,
	poolsDistribution types.Distribution, operatorsDistribution types.Distribution,
	usersDistribution types.UsersDistribution,
) (types.RewardsPlan, error) {
	// TODO: validate arguments.

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	_, found := k.servicesKeeper.GetService(sdkCtx, serviceID)
	if !found {
		return types.RewardsPlan{}, servicestypes.ErrServiceNotFound
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
		poolsDistribution, operatorsDistribution, usersDistribution)
	if err := plan.Validate(k.cdc); err != nil {
		return types.RewardsPlan{}, err
	}

	err = k.validateDistributionDelegationTargets(ctx, poolsDistribution)
	if err != nil {
		return types.RewardsPlan{}, err
	}
	err = k.validateDistributionDelegationTargets(ctx, operatorsDistribution)
	if err != nil {
		return types.RewardsPlan{}, err
	}
	// We don't need to validate users distribution since there's
	// types.UsersDistributionTypeBasic only which doesn't need a validation.

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

// validateDistributionDelegationTargets validates types.Distribution and
// returns an error if any of delegation targets specified is not found.
func (k *Keeper) validateDistributionDelegationTargets(ctx context.Context, distribution types.Distribution) error {
	var distrType types.DistributionType
	err := k.cdc.UnpackAny(distribution.Type, &distrType)
	if err != nil {
		return err
	}
	typ, ok := distrType.(*types.DistributionTypeWeighted)
	if !ok {
		// Only weighted distribution needs a validation.
		return nil
	}
	for _, weight := range typ.Weights {
		_, err = k.GetDelegationTarget(ctx, distribution.DelegationType, weight.DelegationTargetID)
		if err != nil {
			return errors.Wrapf(err, "cannot get delegation target %d", weight.DelegationTargetID)
		}
	}
	return nil
}
