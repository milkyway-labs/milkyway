package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

// RegisterInvariants registers all operators module invariants
func RegisterInvariants(ir sdk.InvariantRegistry, keeper *Keeper) {
	ir.RegisterRoute(types.ModuleName, "valid-operators",
		ValidOperatorsInvariant(keeper))
}

// --------------------------------------------------------------------------------------------------------------------

// ValidOperatorsInvariant checks that all the operators are valid
func ValidOperatorsInvariant(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (message string, broken bool) {
		var invalidOperators []types.Operator
		k.IterateOperators(ctx, func(operator types.Operator) (stop bool) {
			invalid := false

			// Make sure the operator ID is never greater or equal to the next operator ID
			nextOperatorID, err := k.GetNextOperatorID(ctx)
			if err != nil {
				invalid = true
			}

			if operator.ID >= nextOperatorID {
				invalid = true
			}

			// Make sure the operator is valid
			err = operator.Validate()
			if err != nil {
				invalid = true
			}

			if invalid {
				invalidOperators = append(invalidOperators, operator)
			}

			return false
		})

		return sdk.FormatInvariant(types.ModuleName, "invalid operators",
			fmt.Sprintf("the following operators are invalid:\n %s", formatOutputOperators(invalidOperators)),
		), invalidOperators != nil
	}
}

// formatOutputOperators concatenates the given operators information into a string
func formatOutputOperators(operators []types.Operator) (output string) {
	for _, operator := range operators {
		output += operator.String() + "\n"
	}
	return output
}