package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	sdkmath "cosmossdk.io/math"

	"github.com/milkyway-labs/milkyway/v9/x/liquidvesting/types"
	restakingtypes "github.com/milkyway-labs/milkyway/v9/x/restaking/types"
)

// beforeInsurancePercentageChanged is called before the insurance percentage
// parameter is changed. It iterates over all locked representation delegators
// and their delegations to withdraw their rewards from delegation targets that
// they have delegated locked tokens to.
func (k *Keeper) beforeInsurancePercentageChanged(ctx context.Context, oldPercentage, newPercentage sdkmath.LegacyDec) error {
	type delegationTargetCacheKey struct {
		delType  restakingtypes.DelegationType
		targetID uint32
	}
	delegationTargetCache := map[delegationTargetCacheKey]restakingtypes.DelegationTarget{}

	delegators, err := k.GetAllLockedRepresentationDelegators(ctx)
	if err != nil {
		return err
	}

	for _, delegator := range delegators {
		delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(delegator)
		if err != nil {
			return err
		}

		err = k.restakingKeeper.IterateUserDelegations(ctx, delegator, func(del restakingtypes.Delegation) (stop bool, err error) {
			// If the delegation has no locked shares, skip.
			if !types.HasLockedShares(del.Shares) {
				return false, nil
			}

			target, ok := delegationTargetCache[delegationTargetCacheKey{del.Type, del.TargetID}]
			if !ok {
				target, err = k.restakingKeeper.GetDelegationTarget(ctx, del.Type, del.TargetID)
				if err != nil {
					return true, err
				}
				delegationTargetCache[delegationTargetCacheKey{del.Type, del.TargetID}] = target
			}

			userInsuranceFund, err := k.GetUserInsuranceFundBalance(ctx, delegator)
			if err != nil {
				return true, err
			}
			oldCoveredLockedShares, err := types.GetCoveredLockedShares(target, del, userInsuranceFund, oldPercentage)
			if err != nil {
				return true, err
			}
			newCoveredLockedShares, err := types.GetCoveredLockedShares(target, del, userInsuranceFund, newPercentage)
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
				target, _, err := types.RemoveDelShares(target, uncoveredLockedShares)
				return target, err
			}
			// Initialize the delegation with the adjusted shares using the new insurance
			// percentage.
			k.restakingOverrider.GetDelegation = func(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32, delegator string) (restakingtypes.Delegation, bool, error) {
				del.Shares = types.DeductUncoveredLockedShares(del.Shares, newCoveredLockedShares)
				return del, true, nil
			}
			err = k.withRestakingOverrider(func() error {
				_, err := k.rewardsKeeper.WithdrawDelegationRewards(ctx, delAddr, del.Type, del.TargetID)
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
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *Keeper) SetParams(ctx context.Context, params types.Params) error {
	err := params.Validate()
	if err != nil {
		return err
	}

	isFirst := false // Whether the params are being set for the first time
	oldParams, err := k.params.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			oldParams = types.DefaultParams()
			isFirst = true
		} else {
			return err
		}
	}

	// If the insurance percentage has changed, we need to withdraw all delegators
	// restaking rewards who have delegated locked tokens.
	if !isFirst && !params.InsurancePercentage.Equal(oldParams.InsurancePercentage) {
		err = k.beforeInsurancePercentageChanged(ctx, oldParams.InsurancePercentage, params.InsurancePercentage)
		if err != nil {
			return err
		}
	}

	err = k.params.Set(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func (k *Keeper) GetParams(ctx context.Context) (types.Params, error) {
	params, err := k.params.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return types.DefaultParams(), nil
		} else {
			return types.Params{}, err
		}
	}
	return params, nil
}
