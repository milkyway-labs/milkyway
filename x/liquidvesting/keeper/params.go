package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	sdkmath "cosmossdk.io/math"

	"github.com/milkyway-labs/milkyway/v9/x/liquidvesting/types"
	operatorstypes "github.com/milkyway-labs/milkyway/v9/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/v9/x/pools/types"
	restakingtypes "github.com/milkyway-labs/milkyway/v9/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v9/x/services/types"
)

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

			k.restakingOverrider.GetDelegationTarget = func(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32) (restakingtypes.DelegationTarget, error) {
				uncoveredLockedShares := types.UncoveredLockedShares(target.GetDelegatorShares(), newTargetCoveredLockedShares)

				switch target := target.(type) {
				case poolstypes.Pool:
					target, _, err = target.RemoveDelShares(uncoveredLockedShares)
					if err != nil {
						return nil, err
					}
					return target, nil
				case operatorstypes.Operator:
					target, _ = target.RemoveDelShares(uncoveredLockedShares)
					return target, nil
				case servicestypes.Service:
					target, _ = target.RemoveDelShares(uncoveredLockedShares)
					return target, nil
				default:
					return nil, fmt.Errorf("invalid target type %T", target)
				}
			}
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

			err = k.TargetsCoveredLockedShares.Set(
				ctx,
				collections.Join(int32(del.Type), del.TargetID),
				types.TargetCoveredLockedShares{Shares: newTargetCoveredLockedShares},
			)
			if err != nil {
				return true, err
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

	oldParams, err := k.GetParams(ctx)
	if err != nil {
		return err
	}

	// If the insurance percentage has changed, we need to withdraw all delegators
	// restaking rewards who have delegated locked tokens.
	if !params.InsurancePercentage.Equal(oldParams.InsurancePercentage) {
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
