package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"

	"github.com/milkyway-labs/milkyway/v9/utils"
	"github.com/milkyway-labs/milkyway/v9/x/liquidvesting/types"
	operatorstypes "github.com/milkyway-labs/milkyway/v9/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/v9/x/pools/types"
	restakingtypes "github.com/milkyway-labs/milkyway/v9/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v9/x/services/types"
)

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
		// TODO: store locked shares delegators and iterate through them only
		// TODO: DON'T USE GetAllDelegations!!!
		delegations, err := k.restakingKeeper.GetAllDelegations(ctx)
		if err != nil {
			return err
		}
		for _, delegation := range delegations {
			// Only withdraw rewards for delegations that have locked shares
			hasLockedShares := false
			for _, share := range delegation.Shares {
				tokenDenom := utils.GetTokenDenomFromSharesDenom(share.Denom)
				if types.IsLockedRepresentationDenom(tokenDenom) {
					hasLockedShares = true
					break
				}
			}
			if !hasLockedShares {
				continue
			}

			delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(delegation.UserAddress)
			if err != nil {
				return err
			}

			// TODO: optimize GetDelegationTarget
			target, err := k.restakingKeeper.GetDelegationTarget(ctx, delegation.Type, delegation.TargetID)
			if err != nil {
				return err
			}

			userInsuranceFund, err := k.GetUserInsuranceFundBalance(ctx, delegation.UserAddress)
			if err != nil {
				return err
			}
			oldCoveredLockedShares, err := types.GetCoveredLockedShares(target, delegation, userInsuranceFund, oldParams.InsurancePercentage)
			if err != nil {
				return err
			}
			newCoveredLockedShares, err := types.GetCoveredLockedShares(target, delegation, userInsuranceFund, params.InsurancePercentage)
			if err != nil {
				return err
			}
			oldTargetCoveredLockedShares, err := k.GetTargetCoveredLockedShares(ctx, delegation.Type, delegation.TargetID)
			if err != nil {
				return err
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
				delegation.Shares = types.DeductUncoveredLockedShares(delegation.Shares, newCoveredLockedShares)
				return delegation, true, nil
			}
			err = k.withRestakingOverrider(func() error {
				_, err := k.rewardsKeeper.WithdrawDelegationRewards(ctx, delAddr, delegation.Type, delegation.TargetID)
				return err
			})
			if err != nil {
				return err
			}

			err = k.TargetsCoveredLockedShares.Set(
				ctx,
				collections.Join(int32(delegation.Type), delegation.TargetID),
				types.CoveredLockedShares{Shares: newTargetCoveredLockedShares},
			)
			if err != nil {
				return err
			}
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
