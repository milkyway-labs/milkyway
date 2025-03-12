package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v9/x/liquidvesting/types"
	restakingtypes "github.com/milkyway-labs/milkyway/v9/x/restaking/types"
)

// delegationTargetCacheKey is the map key for the delegation target cache.
type delegationTargetCacheKey struct {
	Type     restakingtypes.DelegationType
	TargetID uint32
}

// delegationTargetCache represents a cache for delegation targets.
type delegationTargetCache map[delegationTargetCacheKey]restakingtypes.DelegationTarget

// coveredLockedSharesParamsFn represents a function to derive the insurance fund
// and insurance percentage to calculate the covered locked shares.
type coveredLockedSharesParamsFn func() (
	insuranceFund sdk.Coins,
	insurancePercentage sdkmath.LegacyDec,
	activeLockedTokens sdk.DecCoins,
)

// WithdrawUserLockedRestakingRewards withdraws all restaking rewards from the
// delegation targets that the user has delegated locked tokens to.
func (k *Keeper) WithdrawUserLockedRestakingRewards(
	ctx context.Context,
	user string,
	filterFn func(del restakingtypes.Delegation) bool,
	oldParamsFn coveredLockedSharesParamsFn,
	newParamsFn coveredLockedSharesParamsFn,
) error {
	return k.WithdrawAllUserRestakingRewardsWithCache(
		ctx,
		user,
		filterFn,
		oldParamsFn,
		newParamsFn,
		delegationTargetCache{},
	)
}

// WithdrawAllUserRestakingRewardsWithCache withdraws all restaking rewards from
// the delegation targets that the user has delegated locked tokens to. If filterFn
// returns false, the delegation will be skipped. oldParamsFn and newParamsFn are
// functions that return the insurance fund, insurance percentage, and active locked
// tokens to calculate the covered locked shares.
func (k *Keeper) WithdrawAllUserRestakingRewardsWithCache(
	ctx context.Context,
	user string,
	filterFn func(del restakingtypes.Delegation) bool,
	oldParamsFn coveredLockedSharesParamsFn,
	newParamsFn coveredLockedSharesParamsFn,
	delTargetCache delegationTargetCache,
) error {
	userAddr, err := k.accountKeeper.AddressCodec().StringToBytes(user)
	if err != nil {
		return err
	}

	err = k.restakingKeeper.IterateUserDelegations(ctx, user, func(del restakingtypes.Delegation) (stop bool, err error) {
		// If the delegation has no locked shares, skip.
		if !types.HasLockedShares(del.Shares) {
			return false, nil
		}

		// If the filter returns false, skip the delegation.
		if !filterFn(del) {
			return false, nil
		}

		target, ok := delTargetCache[delegationTargetCacheKey{del.Type, del.TargetID}]
		if !ok {
			target, err = k.restakingKeeper.GetDelegationTarget(ctx, del.Type, del.TargetID)
			if err != nil {
				return true, err
			}
			delTargetCache[delegationTargetCacheKey{del.Type, del.TargetID}] = target
		}

		oldInsuranceFund, oldInsurancePercentage, oldActiveLockedTokens := oldParamsFn()
		oldCoveredLockedShares, err := types.GetCoveredLockedShares(
			target,
			del,
			oldInsuranceFund,
			oldInsurancePercentage,
			oldActiveLockedTokens,
		)
		if err != nil {
			return true, err
		}
		newInsuranceFund, newInsurancePercentage, newActiveLockedTokens := newParamsFn()
		newCoveredLockedShares, err := types.GetCoveredLockedShares(
			target,
			del,
			newInsuranceFund,
			newInsurancePercentage,
			newActiveLockedTokens,
		)
		if err != nil {
			return true, err
		}
		oldTargetCoveredLockedShares, err := k.GetTargetCoveredLockedShares(ctx, del.Type, del.TargetID)
		if err != nil {
			return true, err
		}
		newTargetCoveredLockedShares := oldTargetCoveredLockedShares.Add(newCoveredLockedShares...).Sub(oldCoveredLockedShares)

		// When initializing a new delegation inside WithdrawDelegationRewards, adjust
		// the delegation target's total delegation shares by deducting the uncovered
		// locked shares calculated with the new insurance percentage.
		k.restakingOverrider.GetDelegationTarget = func(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32) (restakingtypes.DelegationTarget, error) {
			uncoveredLockedShares := types.UncoveredLockedShares(target.GetDelegatorShares(), newTargetCoveredLockedShares)
			return types.DelegationTargetWithDeductedShares(target, uncoveredLockedShares)
		}
		// Initialize the delegation with the adjusted shares using the new insurance
		// percentage.
		k.restakingOverrider.GetDelegation = func(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32, delegator string) (restakingtypes.Delegation, bool, error) {
			uncoveredLockedShares := types.UncoveredLockedShares(del.Shares, newCoveredLockedShares)
			del.Shares = del.Shares.Sub(uncoveredLockedShares)
			return del, true, nil
		}
		err = k.withRestakingOverrider(func() error {
			_, err := k.rewardsKeeper.WithdrawDelegationRewards(ctx, userAddr, del.Type, del.TargetID)
			return err
		})
		if err != nil {
			return true, err
		}

		if newTargetCoveredLockedShares.IsZero() {
			err = k.RemoveTargetCoveredLockedShares(ctx, del.Type, del.TargetID)
			if err != nil {
				return true, err
			}
		} else {
			err = k.SetTargetCoveredLockedShares(ctx, del.Type, del.TargetID, newTargetCoveredLockedShares)
			if err != nil {
				return true, err
			}
		}
		return false, nil
	})
	return err
}
