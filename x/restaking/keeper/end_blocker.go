package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v7/x/restaking/types"
)

// CompleteMatureUnbondingDelegations runs the endblocker logic for delegations
func (k *Keeper) CompleteMatureUnbondingDelegations(ctx sdk.Context) error {
	// Remove all mature unbonding delegations from the ubd queue.
	matureUnbonds, err := k.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
	if err != nil {
		return err
	}

	for _, data := range matureUnbonds {

		balances, err := k.CompleteUnbonding(ctx, data)
		if err != nil {
			return err
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteUnbonding,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeUnbondingDelegationType, data.UnbondingDelegationType.String()),
				sdk.NewAttribute(types.AttributeTargetID, fmt.Sprintf("%d", data.TargetID)),
				sdk.NewAttribute(types.AttributeKeyDelegator, data.DelegatorAddress),
			),
		)
	}

	return nil
}
