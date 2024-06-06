package services

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/operators/keeper"
	"github.com/milkyway-labs/milkyway/x/operators/types"
)

// EndBlocker is called at the end of every block.
//
// It iterates over all the operators that are set to be inactivated by the current block time
// and updates their status to inactive.
func EndBlocker(ctx sdk.Context, keeper *keeper.Keeper) {
	// Iterate over all the active polls that have been ended by the current block time
	keeper.IterateInactivatingOperatorQueue(ctx, ctx.BlockTime(), func(operator types.Operator) (stop bool) {

		// Update the operator status
		operator.Status = types.OPERATOR_STATUS_INACTIVE
		err := keeper.UpdateOperator(ctx, operator)
		if err != nil {
			panic(fmt.Sprintf("error while updating operator: %s", err))
		}

		// Emit an event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompletedOperatorInactivation,
				sdk.NewAttribute(types.AttributeKeyOperatorID, fmt.Sprintf("%d", operator.ID)),
			),
		)

		return false
	})
}
