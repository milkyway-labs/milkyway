package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

// ExportGenesis returns a new GenesisState instance containing the information currently present inside the store
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return types.NewGenesis(
		k.GetAllDelegations(ctx),
		k.GetParams(ctx),
	)
}

// InitGenesis initializes the genesis store using the provided data
func (k *Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) {
	// Store the delegations
	for _, delegation := range data.Delegations {
		switch delegation.Type {
		case types.DELEGATION_TYPE_POOL:
			k.SavePoolDelegation(ctx, delegation)
		case types.DELEGATION_TYPE_OPERATOR:
			k.SaveOperatorDelegation(ctx, delegation)
		case types.DELEGATION_TYPE_SERVICE:
			k.SaveServiceDelegation(ctx, delegation)
		}
	}

	// Store the params
	k.SetParams(ctx, data.Params)
}
