package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/pools/types"
)

// RegisterInvariants registers all pools module invariants
func RegisterInvariants(ir sdk.InvariantRegistry, keeper *Keeper) {
	ir.RegisterRoute(types.ModuleName, "valid-pools",
		ValidPoolsInvariant(keeper))
}

// --------------------------------------------------------------------------------------------------------------------

// ValidPoolsInvariant checks that all the pools are valid
func ValidPoolsInvariant(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (message string, broken bool) {

		// Get the next pool id.
		nextPoolID, err := k.GetNextPoolID(ctx)
		if err != nil {
			return sdk.FormatInvariant(types.ModuleName, "invalid pools", "unable to get the next pool ID"), true
		}

		var invalidPools []types.Pool
		k.IteratePools(ctx, func(pool types.Pool) (stop bool) {
			invalid := false

			// Make sure the pool ID is never greater or equal to the next pool ID
			if pool.ID >= nextPoolID {
				invalid = true
			}

			// Make sure the pool is valid
			err = pool.Validate()
			if err != nil {
				invalid = true
			}

			if invalid {
				invalidPools = append(invalidPools, pool)
			}

			return false
		})

		return sdk.FormatInvariant(types.ModuleName, "invalid pools",
			fmt.Sprintf("the following pools are invalid:\n %s", formatOutputPools(invalidPools)),
		), invalidPools != nil
	}
}

// formatOutputPools concatenates the given pools information into a string
func formatOutputPools(pools []types.Pool) (output string) {
	for _, pool := range pools {
		output += fmt.Sprintf("%d\n", pool.ID)
	}
	return output
}
