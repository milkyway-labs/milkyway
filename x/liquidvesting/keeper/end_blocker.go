package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v4/x/liquidvesting/types"
)

// CompleteBurnCoins runs the endblocker logic to burn the coins after the
// undelegation completes.
func (k *Keeper) CompleteBurnCoins(ctx sdk.Context) error {
	// Remove all the information about the coins to burn.
	coinsToBurn, err := k.DequeueAllBurnCoinsFromUnbondingQueue(ctx, ctx.BlockHeader().Time)
	if err != nil {
		return err
	}

	for _, data := range coinsToBurn {
		accAddr, err := sdk.AccAddressFromBech32(data.DelegatorAddress)
		if err != nil {
			return err
		}

		err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, accAddr, types.ModuleName, data.Amount)
		if err != nil {
			return err
		}

		err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, data.Amount)
		if err != nil {
			return err
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeBurnLockedRepresentation,
				sdk.NewAttribute(sdk.AttributeKeyAmount, data.Amount.String()),
				sdk.NewAttribute(types.AttributeKeyUser, data.DelegatorAddress),
			),
		)
	}

	return nil
}
