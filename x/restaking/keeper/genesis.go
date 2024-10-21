package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

// ExportGenesis returns a new GenesisState instance containing the information currently present inside the store
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	operatorsJoinedServices, err := k.GetAllOperatorsJoinedServices(ctx)
	if err != nil {
		panic(err)
	}

	return types.NewGenesis(
		operatorsJoinedServices,
		k.GetAllServicesParams(ctx),
		k.GetAllDelegations(ctx),
		k.GetAllUnbondingDelegations(ctx),
		k.GetParams(ctx),
	)
}

// InitGenesis initializes the genesis store using the provided data
func (k *Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) {
	// Store the services joined by the operators
	for _, record := range data.OperatorsJoinedServices {
		err := k.SetOperatorJoinedServices(ctx, record.OperatorID, record.JoinedServices)
		if err != nil {
			panic(err)
		}
	}

	for _, record := range data.ServicesParams {
		k.SaveServiceParams(ctx, record.ServiceID, record.Params)
	}

	// Store the delegations
	for _, delegation := range data.Delegations {
		err := k.SetDelegation(ctx, delegation)
		if err != nil {
			panic(err)
		}
	}

	// Store the unbonding delegations
	for _, ubd := range data.UnbondingDelegations {
		_, err := k.SetUnbondingDelegation(ctx, ubd)
		if err != nil {
			panic(err)
		}

		for _, entry := range ubd.Entries {
			k.InsertUBDQueue(ctx, ubd, entry.CompletionTime)
		}
	}

	// Store the params
	k.SetParams(ctx, data.Params)
}
