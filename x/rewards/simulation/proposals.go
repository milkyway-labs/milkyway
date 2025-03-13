package simulation

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/milkyway-labs/milkyway/v10/utils"
	"github.com/milkyway-labs/milkyway/v10/x/rewards/types"
)

const (
	DefaultWeightMsgUpdateParams = 30

	OperationWeightMsgUpdateParams = "op_weight_msg_update_params"
)

func ProposalMsgs(bankKeeper bankkeeper.Keeper) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			OperationWeightMsgUpdateParams,
			DefaultWeightMsgUpdateParams,
			SimulateMsgUpdateParams(bankKeeper),
		),
	}
}

// SimulateMsgUpdateParams returns a random MsgUpdateParams
func SimulateMsgUpdateParams(bankKeeper bankkeeper.Keeper) func(r *rand.Rand, _ sdk.Context, _ []simtypes.Account) sdk.Msg {
	return func(r *rand.Rand, ctx sdk.Context, _ []simtypes.Account) sdk.Msg {
		// Use the default gov module account address as authority
		var authority sdk.AccAddress = address.Module(govtypes.ModuleName)

		// Get all the denoms
		metadata := bankKeeper.GetAllDenomMetaData(ctx)
		denoms := utils.Map(metadata, func(md banktypes.Metadata) string {
			return md.Base
		})

		// Generate the new random params
		params := RandomParams(r, denoms)

		// Return the message
		return types.NewMsgUpdateParams(params, authority.String())
	}
}
