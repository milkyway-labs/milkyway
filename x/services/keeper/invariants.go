package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/services/types"
)

// RegisterInvariants registers all services module invariants
func RegisterInvariants(ir sdk.InvariantRegistry, keeper *Keeper) {
	ir.RegisterRoute(types.ModuleName, "valid-services",
		ValidServicesInvariant(keeper))
}

// --------------------------------------------------------------------------------------------------------------------

// ValidServicesInvariant checks that all the services are valid
func ValidServicesInvariant(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (message string, broken bool) {
		var invalidServices []types.Service
		k.IterateServices(ctx, func(service types.Service) (stop bool) {
			invalid := false

			// Make sure the service ID is never greater or equal to the next service ID
			nextServiceID, err := k.GetNextServiceID(ctx)
			if err != nil {
				invalid = true
			}

			if service.ID >= nextServiceID {
				invalid = true
			}

			// Make sure the service is valid
			err = service.Validate()
			if err != nil {
				invalid = true
			}

			if invalid {
				invalidServices = append(invalidServices, service)
			}

			return false
		})

		return sdk.FormatInvariant(types.ModuleName, "invalid services",
			fmt.Sprintf("the following services are invalid:\n %s", formatOutputServices(invalidServices)),
		), invalidServices != nil
	}
}

// formatOutputServices concatenates the given services information into a string
func formatOutputServices(services []types.Service) (output string) {
	for _, service := range services {
		output += fmt.Sprintf("%d\n", service.ID)
	}
	return output
}