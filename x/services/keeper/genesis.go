package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/services/types"
)

// ExportGenesis returns the GenesisState associated with the given context
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	servicesParams, err := k.GetAllServicesParams(ctx)
	if err != nil {
		panic(err)
	}

	return types.NewGenesisState(
		k.exportNextServiceID(ctx),
		k.GetServices(ctx),
		servicesParams,
		k.GetParams(ctx),
	)
}

// exportNextServiceID returns the next Service ID stored in the KVStore
func (k *Keeper) exportNextServiceID(ctx sdk.Context) uint32 {
	nextServiceID, err := k.GetNextServiceID(ctx)
	if err != nil {
		panic(err)
	}
	return nextServiceID
}

// --------------------------------------------------------------------------------------------------------------------

// InitGenesis initializes the state from a GenesisState
func (k *Keeper) InitGenesis(ctx sdk.Context, state *types.GenesisState) error {
	// Set the next service ID
	k.SetNextServiceID(ctx, state.NextServiceID)

	// Store the services
	for _, service := range state.Services {
		if err := k.CreateService(ctx, service); err != nil {
			return err
		}
	}

	for _, serviceParams := range state.ServicesParams {
		err := k.SetServiceParams(ctx, serviceParams.ServiceID, serviceParams.Params)
		if err != nil {
			return err
		}
	}

	// Store params
	k.SetParams(ctx, state.Params)

	return nil
}
