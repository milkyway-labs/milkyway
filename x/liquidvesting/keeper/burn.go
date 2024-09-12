package keeper

import (
	"slices"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// IsBurner tells if a user have the permissions to burn tokens
// from a user's balance.
func (k *Keeper) IsBurner(ctx sdk.Context, user sdk.AccAddress) (bool, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return false, err
	}

	return slices.Contains(params.Burners, user.String()), nil
}

// BurnVestedRepresentation burns the vested staking representation
// from the user's balance.
// NOTE: If the coins are restaked they will be unstaked first.
func (k *Keeper) BurnVestedRepresentation(
	_ sdk.Context,
	_ sdk.AccAddress,
	_ sdk.Coins,
) error {
	panic("unimplemented")
}
