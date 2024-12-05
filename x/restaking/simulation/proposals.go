package simulation

import (
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/milkyway-labs/milkyway/v2/x/restaking/keeper"
	"github.com/milkyway-labs/milkyway/v2/x/restaking/types"
)

// Simulation operation weights constants
const (
	DefaultWeightMsgUpdateParams int = 50

	OperationWeightMsgUpdateParams = "op_weight_msg_update_params"
)

// ProposalMsgs defines the module weighted proposals' contents
func ProposalMsgs(keeper *keeper.Keeper) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			OperationWeightMsgUpdateParams,
			DefaultWeightMsgUpdateParams,
			SimulateMsgUpdateParams,
		),
	}
}

// SimulateMsgUpdateParams returns a random MsgUpdateParams
func SimulateMsgUpdateParams(r *rand.Rand, _ sdk.Context, _ []simtypes.Account) sdk.Msg {
	// use the default gov module account address as authority
	var authority sdk.AccAddress = address.Module("gov")

	params := types.DefaultParams()
	unbondingDays := time.Duration(r.Intn(7) + 1)
	params.UnbondingTime = time.Hour * 24 * unbondingDays

	return &types.MsgUpdateParams{
		Authority: authority.String(),
		Params:    params,
	}
}
