package v2

import (
	"encoding/binary"

	corestoretypes "cosmossdk.io/core/store"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v7/x/rewards/types"
)

// MigrateStore performs in-place store migrations from v1 to v2. The migrations include:
// - Migrate AmountPerDay from sdk.Coins to sdk.Coin
func MigrateStore(ctx sdk.Context, storeService corestoretypes.KVStoreService, cdc codec.BinaryCodec) error {
	return migratePlans(ctx, storeService, cdc)
}

// getNextPlanID returns the next plan ID and increments the next plan ID in the store.
func getNextPlanID(ctx sdk.Context, storeService corestoretypes.KVStoreService) (uint64, error) {
	store := storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.NextRewardsPlanIDKey)
	if err != nil {
		return 0, nil
	}

	// Get the next ID
	var nextID uint64 = 1
	if bz != nil {
		nextID = binary.BigEndian.Uint64(bz)
	}

	// Update the next ID
	err = store.Set(types.NextRewardsPlanIDKey, sdk.Uint64ToBigEndian(nextID+1))
	if err != nil {
		return 0, err
	}

	return nextID, nil
}

// migratePlans migrates the legacy RewardsPlan to the new types.RewardsPlan.
// It does so by iterating over all the legacy plans, and creating a new types.RewardsPlan for each
// existing RewardsPlan. If an existing RewardsPlan has multiple denoms, it will create a new types.RewardsPlan
// for each denom.
//
// After the iteration is done, the following will be true:
// - All the new plans will be saved in the store
// - The next plan ID will be updated to be equal to the highest plan ID + 1, if any new plans were created
func migratePlans(ctx sdk.Context, storeService corestoretypes.KVStoreService, cdc codec.BinaryCodec) error {
	store := storeService.OpenKVStore(ctx)
	iterator, err := store.Iterator(types.RewardsPlanKeyPrefix, storetypes.PrefixEndBytes(types.RewardsPlanKeyPrefix))
	if err != nil {
		return err
	}

	// Get all the legacy plans
	var legacyPlans []RewardsPlan
	for ; iterator.Valid(); iterator.Next() {
		var plan RewardsPlan
		if err := cdc.Unmarshal(iterator.Value(), &plan); err != nil {
			return err
		}
		legacyPlans = append(legacyPlans, plan)
	}

	// Close the iterator
	err = iterator.Close()
	if err != nil {
		return err
	}

	// Convert each legacy plans to n new plans, splitting the denoms
	var newPlans []types.RewardsPlan
	for _, plan := range legacyPlans {
		// Create the basic plan
		newPlans = append(newPlans, mapRewardsPlan(plan, plan.ID, plan.AmountPerDay[0]))

		if len(plan.AmountPerDay) == 1 {
			continue
		}

		// Create additional plans for each denoms after the first one
		for _, coin := range plan.AmountPerDay[1:] {
			planID, err := getNextPlanID(ctx, storeService)
			if err != nil {
				return err
			}

			newPlans = append(newPlans, mapRewardsPlan(plan, planID, coin))
		}
	}

	// Save the new plans
	for _, plan := range newPlans {
		bz, err := cdc.Marshal(&plan)
		if err != nil {
			return err
		}

		err = store.Set(PlanStoreKey(plan.ID), bz)
		if err != nil {
			return err
		}
	}

	return nil
}

// mapRewardsPlan converts a legacy RewardsPlan to a new types.RewardsPlan with a single denom.
// The new plan will have the given ID and the given amount per day.
func mapRewardsPlan(plan RewardsPlan, id uint64, amountPerDay sdk.Coin) types.RewardsPlan {
	return types.NewRewardsPlan(
		id,
		plan.Description,
		plan.ServiceID,
		amountPerDay,
		plan.StartTime,
		plan.EndTime,
		plan.PoolsDistribution,
		plan.OperatorsDistribution,
		plan.UsersDistribution,
	)
}
