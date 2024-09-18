package keeper

import (
	"slices"

	"github.com/milkyway-labs/milkyway/utils"
	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"

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
	ctx sdk.Context,
	accAddress sdk.AccAddress,
	amount sdk.Coins,
) error {
	isBurner, err := k.IsBurner(ctx, accAddress)
	if err != nil {
		return nil
	}
	if !isBurner {
		return types.ErrNotBurner
	}

	// Get the user balance
	userBalance := k.bankKeeper.GetAllBalances(ctx, accAddress)
	userBalance = utils.IntersectCoinsByDenom(userBalance, amount)
	// The amount to burn is not in the user balance, check if we can remove that
	// amount from the user's delegations.
	if !userBalance.IsAllGTE(amount) {
		// Compute the amount that should be unstaked
		toUnstake := amount.Sub(userBalance...)
		err := k.restakingKeeper.UndelegateRestakedAssets(ctx, accAddress, toUnstake)
		if err != nil {
			return err
		}

		// TODO: Enqueue the burn after the undelegation completes
	}

	// Burn the coins
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, accAddress, types.ModuleName, amount); err != nil {
		return err
	}
	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, amount); err != nil {
		return err
	}

	return nil
}
