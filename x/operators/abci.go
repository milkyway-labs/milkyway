package operators

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v11/x/operators/keeper"
	"github.com/milkyway-labs/milkyway/v11/x/operators/types"
)

// BeginBlocker is called at the beginning of every block.
//
// It iterates over all the operators that are set to be inactivated by the current block time
// and updates their status to inactive.
func BeginBlocker(ctx sdk.Context, keeper *keeper.Keeper) error {
	// Iterate over all the operators that are set to be inactivated by the current block time
	return keeper.IterateInactivatingOperatorQueue(ctx, ctx.BlockTime(), func(operator types.Operator) (stop bool, err error) {
		// Complete the operator inactivation process
		err = keeper.CompleteOperatorInactivation(ctx, operator)
		if err != nil {
			return true, err
		}

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
