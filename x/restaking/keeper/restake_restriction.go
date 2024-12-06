package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v3/x/restaking/types"
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
	if k.restakeRestriction == nil {
		return nil
	}

	return k.restakeRestriction(ctx, restakerAddress, restakedAmount, target)
}
