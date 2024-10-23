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

	servicesAllowedOperators, err := k.GetAllServicesAllowedOperators(ctx)
	if err != nil {
		panic(err)
	}

	servicesSecuringPools, err := k.GetAllServicesSecuringPools(ctx)
	if err != nil {
		panic(err)
	}

	return types.NewGenesis(
		operatorsJoinedServices,
		servicesAllowedOperators,
		servicesSecuringPools,
		k.GetAllDelegations(ctx),
		k.GetAllUnbondingDelegations(ctx),
		k.GetParams(ctx),
	)
}

// InitGenesis initializes the genesis store using the provided data
func (k *Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) {
	// Store the services joined by the operators
	for _, record := range data.OperatorsJoinedServices {
		err := k.SaveOperatorJoinedServices(ctx, record.OperatorID, record.JoinedServices)
		if err != nil {
			panic(err)
		}
	}

	// Store the whitelisted operators for each service
	for _, record := range data.ServicesAllowedOperators {
		for _, operatorID := range record.OperatorIDs {
			err := k.AddOperatorToServiceAllowList(ctx, record.ServiceID, operatorID)
			if err != nil {
				panic(err)
			}
		}
	}

	// Store the whitelisted pools for each service
	for _, record := range data.ServicesSecuringPools {
		for _, poolID := range record.PoolIDs {
			err := k.AddPoolToServiceSecuringPools(ctx, record.ServiceID, poolID)
			if err != nil {
				panic(err)
			}
		}
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
