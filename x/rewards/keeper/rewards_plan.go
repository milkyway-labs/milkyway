package keeper

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

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

	// TODO: check if pools, operators exist

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
	if err := plan.Validate(); err != nil {
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
