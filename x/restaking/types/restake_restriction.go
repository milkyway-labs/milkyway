package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RestakeRestrictionFn is a function that checks if a restake operation is allowed.
type RestakeRestrictionFn func(ctx sdk.Context, restakerAddrees string, restakedAmount sdk.Coins, target DelegationTarget) error
