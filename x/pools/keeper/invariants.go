package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v3/utils"
	"github.com/milkyway-labs/milkyway/v3/x/pools/types"
)

// RegisterInvariants registers all pools module invariants
func RegisterInvariants(ir sdk.InvariantRegistry, keeper *Keeper) {
	ir.RegisterRoute(types.ModuleName, "valid-pools",
		ValidPoolsInvariant(keeper))
	ir.RegisterRoute(types.ModuleName, "unique-pools",
		UniquePoolsInvariant(keeper))
}

// --------------------------------------------------------------------------------------------------------------------

// formatOutputPools concatenates the given pools information into a string
func formatOutputPools(pools []types.Pool) (output string) {
	// Get the unique IDs
	uniquePoolIDs := utils.RemoveDuplicates(utils.Map(pools, func(pool types.Pool) uint32 {
		return pool.ID
	}))

	// Create the message string
	for _, poolID := range uniquePoolIDs {
		output += fmt.Sprintf("%d\n", poolID)
	}

	return output
}

// ValidPoolsInvariant checks that all the pools are valid
func ValidPoolsInvariant(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (message string, broken bool) {

		// Get the next pool id.
		nextPoolID, err := k.GetNextPoolID(ctx)
		if err != nil {
			return sdk.FormatInvariant(types.ModuleName, "invalid pools", "unable to get the next pool ID"), true
		}

		var invalidPools []types.Pool
		err = k.IteratePools(ctx, func(pool types.Pool) (stop bool, err error) {
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

			return false, nil
		})
		if err != nil {
			panic(err)
		}

		return sdk.FormatInvariant(types.ModuleName, "invalid pools",
			fmt.Sprintf("the following pools are invalid:\n %s", formatOutputPools(invalidPools)),
		), invalidPools != nil
	}
}

// UniquePoolsInvariant checks that there are no duplicated pools for the same denom
func UniquePoolsInvariant(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (message string, broken bool) {

		var invalidPools []types.Pool
		err := k.IteratePools(ctx, func(pool types.Pool) (stop bool, err error) {
			otherPool, found, err := k.GetPoolByDenom(ctx, pool.Denom)
			if err != nil {
				return false, err
			}

			if found && otherPool.ID != pool.ID {
				invalidPools = append(invalidPools, pool)
			}
			return false, nil
		})
		if err != nil {
			panic(err)
		}

		return sdk.FormatInvariant(types.ModuleName, "invalid pools",
			fmt.Sprintf("the following pools have the same denoms:\n %s", formatOutputPools(invalidPools)),
		), invalidPools != nil
	}
}
