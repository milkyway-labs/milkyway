package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v11/x/restaking/types"
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

	delegations, err := k.GetAllDelegations(ctx)
	if err != nil {
		panic(err)
	}

	unbondingDelegations, err := k.GetAllUnbondingDelegations(ctx)
	if err != nil {
		panic(err)
	}

	preferences, err := k.GetUserPreferencesEntries(ctx)
	if err != nil {
		panic(err)
	}

	params, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}

	return types.NewGenesis(
		operatorsJoinedServices,
		servicesAllowedOperators,
		servicesSecuringPools,
		delegations,
		unbondingDelegations,
		preferences,
		params,
	)
}

// InitGenesis initializes the genesis store using the provided data
func (k *Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) error {
	// Store the services joined by the operators
	for _, record := range data.OperatorsJoinedServices {
		for _, serviceID := range record.ServiceIDs {
			err := k.AddServiceToOperatorJoinedServices(ctx, record.OperatorID, serviceID)
			if err != nil {
				return err
			}
		}
	}

	// Store the whitelisted operators for each service
	for _, record := range data.ServicesAllowedOperators {
		for _, operatorID := range record.OperatorIDs {
			err := k.AddOperatorToServiceAllowList(ctx, record.ServiceID, operatorID)
			if err != nil {
				return err
			}
		}
	}

	// Store the whitelisted pools for each service
	for _, record := range data.ServicesSecuringPools {
		for _, poolID := range record.PoolIDs {
			err := k.AddPoolToServiceSecuringPools(ctx, record.ServiceID, poolID)
			if err != nil {
				return err
			}
		}
	}

	// Store the delegations
	for _, delegation := range data.Delegations {
		err := k.SetDelegation(ctx, delegation)
		if err != nil {
			return err
		}
	}

	// Store the unbonding delegations
	for _, ubd := range data.UnbondingDelegations {
		_, err := k.SetUnbondingDelegation(ctx, ubd)
		if err != nil {
			return err
		}

		for _, entry := range ubd.Entries {
			err = k.InsertUBDQueue(ctx, ubd, entry.CompletionTime)
			if err != nil {
				return err
			}
		}
	}

	// Store the user preferences
	for _, entry := range data.UsersPreferences {
		err := k.SetUserPreferences(ctx, entry.UserAddress, entry.Preferences)
		if err != nil {
			return err
		}
	}

	// Store the params
	return k.SetParams(ctx, data.Params)
}
