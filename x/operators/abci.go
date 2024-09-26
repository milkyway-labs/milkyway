package operators

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/operators/keeper"
	"github.com/milkyway-labs/milkyway/x/operators/types"
)

// BeginBlocker is called at the beginning of every block.
//
// It iterates over all the operators that are set to be inactivated by the current block time
// and updates their status to inactive.
func BeginBlocker(ctx sdk.Context, keeper *keeper.Keeper) error {
	// Iterate over all the operators that are set to be inactivated by the current block time
	return keeper.IterateInactivatingOperatorQueue(ctx, ctx.BlockTime(), func(operator types.Operator) (stop bool, err error) {
		// Complete the operator inactivation process
		keeper.CompleteOperatorInactivation(ctx, operator)

		// Emit an event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteOperatorInactivation,
				sdk.NewAttribute(types.AttributeKeyOperatorID, fmt.Sprintf("%d", operator.ID)),
			),
		)

		return false, nil
	})
}
