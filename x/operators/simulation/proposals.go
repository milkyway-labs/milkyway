package simulation

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	"github.com/milkyway-labs/milkyway/v4/x/operators/types"
)

const (
	DefaultWeightMsgUpdateParams = 30

	OperationWeightMsgUpdateParams = "op_weight_msg_update_params"
)

func ProposalMsgs(stakingKeeper *stakingkeeper.Keeper) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			OperationWeightMsgUpdateParams,
			DefaultWeightMsgUpdateParams,
			SimulateMsgUpdateParams(stakingKeeper),
		),
	}
}

// SimulateMsgUpdateParams returns a random MsgUpdateParams
func SimulateMsgUpdateParams(stakingKeeper *stakingkeeper.Keeper) func(r *rand.Rand, _ sdk.Context, _ []simtypes.Account) sdk.Msg {
	return func(r *rand.Rand, ctx sdk.Context, _ []simtypes.Account) sdk.Msg {
		// Use the default gov module account address as authority
		var authority sdk.AccAddress = address.Module(govtypes.ModuleName)

		// Get the stake denom
		bondDenom, err := stakingKeeper.BondDenom(ctx)
		if err != nil {
			panic("failed to get bond denom")
		}

		// Generate the new random params
		params := RandomParams(r, bondDenom)

		// Return the message
		return types.NewMsgUpdateParams(params, authority.String())
	}
}
