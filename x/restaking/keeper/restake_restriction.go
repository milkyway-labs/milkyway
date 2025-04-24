package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v11/x/restaking/types"
)

// SetRestakeRestriction sets a function that checks if a restake operation is allowed.
func (k *Keeper) SetRestakeRestriction(restriction types.RestakeRestrictionFn) {
	if k.restakeRestriction != nil {
		panic("restake restriction already set")
	}

	k.restakeRestriction = restriction
}

// ValidateRestake returns nil if the restake operation is allowed, otherwise returns an error.
func (k *Keeper) ValidateRestake(ctx context.Context, restakerAddress string, restakedAmount sdk.Coins, target types.DelegationTarget) error {
	// Check against the restaking cap only if it is non-zero
	restakingCap, err := k.RestakingCap(ctx)
	if err != nil {
		return err
	}
	if !restakingCap.IsZero() {
		totalRestakedAssets, err := k.GetTotalRestakedAssets(ctx)
		if err != nil {
			return err
		}

		// Add newly restaked amount to the total restaked assets
		totalRestakedAssets = totalRestakedAssets.Add(restakedAmount...)

		totalRestakedValue, err := k.GetCoinsValue(ctx, totalRestakedAssets)
		if err != nil {
			return err
		}
		if totalRestakedValue.GT(restakingCap) {
			return types.ErrRestakingCapExceeded.Wrapf(
				"total restaked value %s is greater than the cap %s",
				totalRestakedValue,
				restakingCap,
			)
		}
	}

	if k.restakeRestriction == nil {
		return nil
	}

	return k.restakeRestriction(ctx, restakerAddress, restakedAmount, target)
}
