package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/milkyway-labs/milkyway/v3/x/rewards/keeper"
	"github.com/milkyway-labs/milkyway/v3/x/rewards/types"
)

// Simulation operation weights constants
const (
	DefaultWeightMsgCreateRewardsPlan          int = 80
	DefaultWeightMsgEditRewardsPlan            int = 40
	DefaultWeightMsgSetWithdrawAddress         int = 20
	DefaultWeightMsgWithdrawDelegatorReward    int = 30
	DefaultWeightMsgWithdrawOperatorCommission int = 10

	OperationWeightMsgCreateRewardsPlan          = "op_weight_msg_create_rewards_plan"
	OperationWeightMsgEditRewardsPlan            = "op_weight_msg_edit_rewards_plan"
	OperationWeightMsgSetWithdrawAddress         = "op_weight_msg_set_withdraw_address"
	OperationWeightMsgWithdrawDelegatorReward    = "op_weight_msg_withdraw_delegator_reward"
	OperationWeightMsgWithdrawOperatorCommission = "op_weight_msg_withdraw_operator_commission"
)

func WeightedOperations(
	appParams simtypes.AppParams,
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	k *keeper.Keeper,
) simulation.WeightedOperations {
	var (
		weightMsgCreateRewardsPlan          int
		weightMsgEditRewardsPlan            int
		weightMsgSetWithdrawAddress         int
		weightMsgWithdrawDelegatorReward    int
		weightMsgWithdrawOperatorCommission int
	)

	// Generate the weights for the messages
	appParams.GetOrGenerate(OperationWeightMsgCreateRewardsPlan, &weightMsgCreateRewardsPlan, nil, func(_ *rand.Rand) {
		weightMsgCreateRewardsPlan = DefaultWeightMsgCreateRewardsPlan
	})

	appParams.GetOrGenerate(OperationWeightMsgEditRewardsPlan, &weightMsgEditRewardsPlan, nil, func(_ *rand.Rand) {
		weightMsgEditRewardsPlan = DefaultWeightMsgEditRewardsPlan
	})

	appParams.GetOrGenerate(OperationWeightMsgSetWithdrawAddress, &weightMsgSetWithdrawAddress, nil, func(_ *rand.Rand) {
		weightMsgSetWithdrawAddress = DefaultWeightMsgEditRewardsPlan
	})

	appParams.GetOrGenerate(OperationWeightMsgWithdrawDelegatorReward, &weightMsgWithdrawDelegatorReward, nil, func(_ *rand.Rand) {
		weightMsgWithdrawDelegatorReward = DefaultWeightMsgWithdrawDelegatorReward
	})

	appParams.GetOrGenerate(OperationWeightMsgWithdrawOperatorCommission, &weightMsgWithdrawOperatorCommission, nil, func(_ *rand.Rand) {
		weightMsgWithdrawOperatorCommission = DefaultWeightMsgWithdrawOperatorCommission
	})

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(weightMsgCreateRewardsPlan, SimulateMsgCreateRewardsPlan(ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgEditRewardsPlan, SimulateMsgEditRewardsPlan(ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgSetWithdrawAddress, SimulateMsgSetWithdrawAddress(ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgWithdrawDelegatorReward, SimulateMsgWithdrawDelegatorReward(ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgWithdrawOperatorCommission, SimulateMsgWithdrawOperatorCommission(ak, bk, k)),
	}
}

// --------------------------------------------------------------------------------------------------------------------

func SimulateMsgCreateRewardsPlan(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k *keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgCreateRewardsPlan{}
		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "TODO"), nil, nil
	}
}

func SimulateMsgEditRewardsPlan(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k *keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgEditRewardsPlan{}
		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "TODO"), nil, nil
	}
}

func SimulateMsgSetWithdrawAddress(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k *keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgSetWithdrawAddress{}
		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "TODO"), nil, nil
	}
}

func SimulateMsgWithdrawDelegatorReward(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k *keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgWithdrawDelegatorReward{}
		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "TODO"), nil, nil
	}
}

func SimulateMsgWithdrawOperatorCommission(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k *keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgWithdrawOperatorCommission{}
		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "TODO"), nil, nil
	}
}
