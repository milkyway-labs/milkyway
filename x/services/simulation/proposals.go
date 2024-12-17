package simulation

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/milkyway-labs/milkyway/v4/x/services/keeper"
	"github.com/milkyway-labs/milkyway/v4/x/services/types"
)

// Simulation operation weights constants
//
//nolint:gosec // The followings are not sensitive information
const (
	DefaultWeightMsgUpdateParams               int = 50
	DefaultWeightMsgAccreditService            int = 50
	DefaultWeightMsgRevokeServiceAccreditation int = 50

	OperationWeightMsgUpdateParams               = "op_weight_msg_update_params"
	OperationWeightMsgAccreditService            = "op_weight_msg_accredit_service"
	OperationWeightMsgRevokeServiceAccreditation = "op_weight_msg_revoke_service_accreditation"
)

// ProposalMsgs defines the module weighted proposals' contents
func ProposalMsgs(keeper *keeper.Keeper) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			OperationWeightMsgUpdateParams,
			DefaultWeightMsgUpdateParams,
			SimulateMsgUpdateParams,
		),
		simulation.NewWeightedProposalMsg(
			OperationWeightMsgAccreditService,
			DefaultWeightMsgAccreditService,
			SimulateMsgAccreditService(keeper),
		),
		simulation.NewWeightedProposalMsg(
			OperationWeightMsgRevokeServiceAccreditation,
			DefaultWeightMsgRevokeServiceAccreditation,
			SimulateMsgRevokeServiceAccreditation(keeper),
		),
	}
}

// SimulateMsgUpdateParams returns a random MsgUpdateParams
func SimulateMsgUpdateParams(r *rand.Rand, _ sdk.Context, _ []simtypes.Account) sdk.Msg {
	// use the default gov module account address as authority
	var authority sdk.AccAddress = address.Module("gov")

	params := types.DefaultParams()
	params.ServiceRegistrationFee = sdk.NewCoins(sdk.NewInt64Coin("umilk", int64(r.Intn(10_000_000)+1)))

	return &types.MsgUpdateParams{
		Authority: authority.String(),
		Params:    params,
	}
}

func SimulateMsgAccreditService(
	k *keeper.Keeper,
) simtypes.MsgSimulatorFn {
	return func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
		// Get a random service
		service, found := GetRandomExistingService(r, ctx, k, func(s types.Service) bool {
			return !s.Accredited
		})
		if !found {
			return nil
		}

		// use the default gov module account address as authority
		var authority sdk.AccAddress = address.Module(govtypes.ModuleName)

		// Create the msg
		return types.NewMsgAccreditService(service.ID, authority.String())
	}
}

func SimulateMsgRevokeServiceAccreditation(
	k *keeper.Keeper,
) simtypes.MsgSimulatorFn {
	return func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
		// Get a random service
		service, found := GetRandomExistingService(r, ctx, k, func(s types.Service) bool {
			return s.Accredited
		})
		if !found {
			return nil
		}

		// use the default gov module account address as authority
		var authority sdk.AccAddress = address.Module(govtypes.ModuleName)

		// Create the msg
		return types.NewMsgRevokeServiceAccreditation(service.ID, authority.String())
	}
}
