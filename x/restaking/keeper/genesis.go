package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

// ExportGenesis returns a new GenesisState instance containing the information currently present inside the store
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	var operatorParamsRecords []types.OperatorParamsRecord
	k.IterateAllOperatorParams(ctx, func(operatorID uint32, params types.OperatorParams) (stop bool) {
		operatorParamsRecords = append(operatorParamsRecords, types.OperatorParamsRecord{
			OperatorID: operatorID,
			Params:     params,
		})
		return false
	})

	var serviceParamsRecords []types.ServiceParamsRecord
	k.IterateAllServiceParams(ctx, func(serviceID uint32, params types.ServiceParams) (stop bool) {
		serviceParamsRecords = append(serviceParamsRecords, types.ServiceParamsRecord{
			ServiceID: serviceID,
			Params:    params,
		})
		return false
	})

	return types.NewGenesis(
		operatorParamsRecords,
		serviceParamsRecords,
		k.GetAllDelegations(ctx),
		k.GetParams(ctx),
	)
}

// InitGenesis initializes the genesis store using the provided data
func (k *Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) {
	for _, record := range data.OperatorsParams {
		k.SaveOperatorParams(ctx, record.OperatorID, record.Params)
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

	// Store the params
	k.SetParams(ctx, data.Params)
}
