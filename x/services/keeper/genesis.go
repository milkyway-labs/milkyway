package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v4/x/services/types"
)

// ExportGenesis returns the GenesisState associated with the given context
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	nextServiceID, err := k.GetNextServiceID(ctx)
	if err != nil {
		panic(err)
	}

	services, err := k.GetServices(ctx)
	if err != nil {
		panic(err)
	}

	servicesParams, err := k.GetAllServicesParams(ctx)
	if err != nil {
		panic(err)
	}

	params, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}

	return types.NewGenesisState(
		nextServiceID,
		services,
		servicesParams,
		params,
	)
}

// --------------------------------------------------------------------------------------------------------------------

// InitGenesis initializes the state from a GenesisState
func (k *Keeper) InitGenesis(ctx sdk.Context, state *types.GenesisState) error {
	// Set the next service ID
	err := k.SetNextServiceID(ctx, state.NextServiceID)
	if err != nil {
		return err
	}

	// Store the services
	for _, service := range state.Services {
		err = k.CreateService(ctx, service)
		if err != nil {
			return err
		}
	}

	for _, serviceParams := range state.ServicesParams {
		err = k.SetServiceParams(ctx, serviceParams.ServiceID, serviceParams.Params)
		if err != nil {
			return err
		}
	}

	// Store params
	err = k.SetParams(ctx, state.Params)
	if err != nil {
		return err
	}

	return nil
}
