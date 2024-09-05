package keeper

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/rewards/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

func (k *Keeper) CreateRewardsPlan(
	ctx context.Context,
	description string,
	serviceID uint32,
	amt sdk.Coins,
	startTime,
	endTime time.Time,
	poolsDistribution types.Distribution,
	operatorsDistribution types.Distribution,
	usersDistribution types.UsersDistribution,
) (types.RewardsPlan, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	_, found := k.servicesKeeper.GetService(sdkCtx, serviceID)
	if !found {
		return types.RewardsPlan{}, servicestypes.ErrServiceNotFound
	}

	// Get the plan id to be used
	planID, err := k.NextRewardsPlanID.Get(ctx)
	if err != nil {
		return types.RewardsPlan{}, err
	}

	// Increment the plan id
	err = k.NextRewardsPlanID.Set(ctx, planID+1)
	if err != nil {
		return types.RewardsPlan{}, err
	}

	// Create the rewards plan
	plan := types.NewRewardsPlan(
		planID,
		description,
		serviceID,
		amt,
		startTime,
		endTime,
		poolsDistribution,
		operatorsDistribution,
		usersDistribution,
	)

	// Validate the plan
	err = plan.Validate(k.cdc)
	if err != nil {
		return types.RewardsPlan{}, err
	}

	// Validate the pools distribution
	err = k.validateDistributionDelegationTargets(ctx, poolsDistribution)
	if err != nil {
		return types.RewardsPlan{}, err
	}

	// Validate the operators distribution
	err = k.validateDistributionDelegationTargets(ctx, operatorsDistribution)
	if err != nil {
		return types.RewardsPlan{}, err
	}

	// We don't need to validate users distribution since there's
	// types.UsersDistributionTypeBasic only which doesn't need a validation.

	// Create a rewards pool account if it doesn't exist
	k.createAccountIfNotExists(ctx, plan.MustGetRewardsPoolAddress(k.accountKeeper.AddressCodec()))

	// Store the rewards plan
	err = k.RewardsPlans.Set(ctx, planID, plan)
	if err != nil {
		return types.RewardsPlan{}, err
	}

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

// GetRewardsPlan returns a rewards plan by ID.
func (k *Keeper) GetRewardsPlan(ctx context.Context, planID uint64) (types.RewardsPlan, error) {
	return k.RewardsPlans.Get(ctx, planID)
}

// terminateRewardsPlan removes a rewards plan and transfers the remaining
// rewards in the plan's rewards pool to the service's address.
func (k *Keeper) terminateRewardsPlan(ctx context.Context, plan types.RewardsPlan) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Transfer remaining rewards in the plan's rewards pool to the service's
	// address.
	rewardsPoolAddr := plan.MustGetRewardsPoolAddress(k.accountKeeper.AddressCodec())
	remaining := k.bankKeeper.GetAllBalances(ctx, rewardsPoolAddr)
	if remaining.IsAllPositive() {
		// Get the service's address.
		service, found := k.servicesKeeper.GetService(sdkCtx, plan.ServiceID)
		if !found {
			return servicestypes.ErrServiceNotFound
		}
		serviceAddr, err := k.accountKeeper.AddressCodec().StringToBytes(service.Address)
		if err != nil {
			return err
		}

		// Transfer all the remaining rewards to the service's address.
		err = k.bankKeeper.SendCoins(ctx, rewardsPoolAddr, serviceAddr, remaining)
		if err != nil {
			return err
		}
	}

	// Remove the plan.
	err := k.RewardsPlans.Remove(ctx, plan.ID)
	if err != nil {
		return err
	}

	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTerminateRewardsPlan,
			sdk.NewAttribute(types.AttributeKeyRewardsPlanID, fmt.Sprint(plan.ID)),
			sdk.NewAttribute(types.AttributeKeyRemainingRewards, remaining.String()),
		),
	})

	return nil
}

// TerminateEndedRewardsPlans terminates all rewards plans that have ended.
func (k *Keeper) TerminateEndedRewardsPlans(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	// Get the current block time
	blockTime := sdkCtx.BlockTime()

	// Iterate over all rewards plans
	err := k.RewardsPlans.Walk(ctx, nil, func(planID uint64, plan types.RewardsPlan) (stop bool, err error) {
		// If the plan has already ended, terminate it
		if !blockTime.Before(plan.EndTime) {
			err = k.terminateRewardsPlan(ctx, plan)
			if err != nil {
				return false, err
			}
		}
		return false, nil
	})
	if err != nil {
		return err
	}
	return nil
}
