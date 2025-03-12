package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v9/x/liquidvesting/types"
	restakingtypes "github.com/milkyway-labs/milkyway/v9/x/restaking/types"
)

// beforeInsurancePercentageChanged is called before the insurance percentage
// parameter is changed. It iterates over all locked representation delegators
// and their delegations to withdraw their rewards from delegation targets that
// they have delegated locked tokens to.
func (k *Keeper) beforeInsurancePercentageChanged(ctx context.Context, oldPercentage, newPercentage sdkmath.LegacyDec) error {
	delegators, err := k.GetAllLockedRepresentationDelegators(ctx)
	if err != nil {
		return err
	}

	delTargetCache := delegationTargetCache{}

	for _, delegator := range delegators {
		insuranceFund, err := k.GetUserInsuranceFundBalance(ctx, delegator)
		if err != nil {
			return err
		}

		activeLockedTokens, err := k.GetAllUserActiveLockedRepresentations(ctx, delegator)
		if err != nil {
			return err
		}

		err = k.WithdrawAllUserRestakingRewardsWithCache(
			ctx,
			delegator,
			func(del restakingtypes.Delegation) bool { return true },
			func() (sdk.Coins, sdkmath.LegacyDec, sdk.DecCoins) {
				return insuranceFund, oldPercentage, activeLockedTokens
			},
			func() (sdk.Coins, sdkmath.LegacyDec, sdk.DecCoins) {
				return insuranceFund, newPercentage, activeLockedTokens
			},
			delTargetCache,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetParams sets the module parameters. It also checks if the insurance
// percentage has changed and withdraws all restakers restaking rewards who
// have delegated locked tokens if it has.
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

// GetParams returns the module parameters. If the parameters are not found, it
// returns the default parameters.
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
