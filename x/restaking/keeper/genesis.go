package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

// ExportGenesis returns a new GenesisState instance containing the information currently present inside the store
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return types.NewGenesis(
		k.GetAllPoolDelegations(ctx),
		k.GetAllServiceDelegations(ctx),
		k.GetAllOperatorDelegations(ctx),
		k.GetParams(ctx),
	)
}

// InitGenesis initializes the genesis store using the provided data
func (k *Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) {
	// Store the pools delegations
	for _, delegation := range data.PoolsDelegations {
		k.SavePoolDelegation(ctx, delegation)
	}

	// Store the services delegations
	for _, delegation := range data.ServicesDelegations {
		k.SaveServiceDelegation(ctx, delegation)
	}

	// Store the operators delegations
	for _, delegation := range data.OperatorsDelegations {
		k.SaveOperatorDelegation(ctx, delegation)
	}

	// Store the params
	k.SetParams(ctx, data.Params)
}
