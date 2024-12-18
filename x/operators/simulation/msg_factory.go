package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/milkyway-labs/milkyway/v6/testutils/simtesting"
	"github.com/milkyway-labs/milkyway/v6/x/operators/keeper"
	"github.com/milkyway-labs/milkyway/v6/x/operators/types"
)

// Simulation operation weights constants
const (
	DefaultWeightMsgRegisterOperator          int = 80
	DefaultWeightMsgUpdateOperator            int = 40
	DefaultWeightMsgDeactivateOperator        int = 20
	DefaultWeightMsgReactivateOperator        int = 30
	DefaultWeightMsgTransferOperatorOwnership int = 10
	DefaultWeightMsgDeleteOperator            int = 25
	DefaultWeightMsgSetOperatorParams         int = 10

	OperationWeightMsgRegisterOperator          = "op_weight_msg_register_operator"
	OperationWeightMsgDeactivateOperator        = "op_weight_msg_deactivate_operator"
	OperationWeightMsgUpdateOperator            = "op_weight_msg_update_operator"
	OperationWeightMsgReactivateOperator        = "op_weight_msg_reactivate_operator"
	OperationWeightMsgTransferOperatorOwnership = "op_weight_msg_transfer_operator_ownership"
	OperationWeightMsgDeleteOperator            = "op_weight_msg_delete_operator"
	OperationWeightMsgSetOperatorParams         = "op_weight_msg_set_operator_params"
)

func WeightedOperations(
	appParams simtypes.AppParams,
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	k *keeper.Keeper,
) simulation.WeightedOperations {
	var (
		weightMsgRegisterOperator          int
		weightMsgUpdateOperator            int
		weightMsgDeactivateOperator        int
		weightMsgReactivateOperator        int
		weightMsgTransferOperatorOwnership int
		weightMsgDeleteOperator            int
		weightMsgSetOperatorParams         int
	)

	// Generate the weights for the messages
	appParams.GetOrGenerate(OperationWeightMsgRegisterOperator, &weightMsgRegisterOperator, nil, func(_ *rand.Rand) {
		weightMsgRegisterOperator = DefaultWeightMsgRegisterOperator
	})

	appParams.GetOrGenerate(OperationWeightMsgUpdateOperator, &weightMsgUpdateOperator, nil, func(_ *rand.Rand) {
		weightMsgUpdateOperator = DefaultWeightMsgUpdateOperator
	})

	appParams.GetOrGenerate(OperationWeightMsgDeactivateOperator, &weightMsgDeactivateOperator, nil, func(_ *rand.Rand) {
		weightMsgDeactivateOperator = DefaultWeightMsgDeactivateOperator
	})

	appParams.GetOrGenerate(OperationWeightMsgReactivateOperator, &weightMsgReactivateOperator, nil, func(_ *rand.Rand) {
		weightMsgReactivateOperator = DefaultWeightMsgReactivateOperator
	})

	appParams.GetOrGenerate(OperationWeightMsgTransferOperatorOwnership, &weightMsgTransferOperatorOwnership, nil, func(_ *rand.Rand) {
		weightMsgTransferOperatorOwnership = DefaultWeightMsgTransferOperatorOwnership
	})

	appParams.GetOrGenerate(OperationWeightMsgDeleteOperator, &weightMsgDeleteOperator, nil, func(_ *rand.Rand) {
		weightMsgDeleteOperator = DefaultWeightMsgDeleteOperator
	})

	appParams.GetOrGenerate(OperationWeightMsgSetOperatorParams, &weightMsgSetOperatorParams, nil, func(_ *rand.Rand) {
		weightMsgSetOperatorParams = DefaultWeightMsgSetOperatorParams
	})

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(weightMsgRegisterOperator, SimulateMsgRegisterOperator(ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgUpdateOperator, SimulateUpdateOperator(ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgDeactivateOperator, SimulateDeactivateOperator(ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgReactivateOperator, SimulateReactivateOperator(ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgTransferOperatorOwnership, SimulateTransferOperatorOwnership(ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgDeleteOperator, SimulateDeleteOperator(ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgSetOperatorParams, SimulateSetOperatorParams(ak, bk, k)),
	}
}

// --------------------------------------------------------------------------------------------------------------------

func SimulateMsgRegisterOperator(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k *keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var msg = &types.MsgRegisterOperator{}

		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		// Get a random operator
		operator := RandomOperator(r, accs)

		// Get the admin account that should sign the transaction
		adminAddress, err := sdk.AccAddressFromBech32(operator.Admin)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "invalid admin address"), nil, nil
		}

		// Make sure the admin has enough funds to pay for the creation fees
		params, err := k.GetParams(ctx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get params"), nil, nil
		}

		// Check the fees that should be paid
		feesAmount := sdk.NewCoins()
		for _, feeCoin := range params.OperatorRegistrationFee {
			if bk.GetBalance(ctx, adminAddress, feeCoin.Denom).IsGTE(feeCoin) {
				feesAmount = feesAmount.Add(feeCoin)
				break
			}
		}

		if !params.OperatorRegistrationFee.IsZero() && feesAmount.IsZero() {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "insufficient funds"), nil, nil
		}

		signer, found := simtesting.GetSimAccount(adminAddress, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "admin account not found"), nil, nil
		}

		// Create the message
		msg = types.NewMsgRegisterOperator(
			operator.Moniker,
			operator.Website,
			operator.PictureURL,
			feesAmount,
			operator.Admin,
		)
		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, signer)
	}
}

func SimulateUpdateOperator(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k *keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var msg = &types.MsgUpdateOperator{}

		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		// Get a random operator from the one existing
		operator, found := GetRandomExistingOperator(r, ctx, k, nil)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get operator"), nil, nil
		}

		// Generate a new data
		updatedMoniker := types.DoNotModify
		if r.Intn(100) < 50 {
			// 50% chance of changing the moniker
			updatedMoniker = simtypes.RandStringOfLength(r, 20)
		}

		updatedWebsite := types.DoNotModify
		if r.Intn(100) < 50 {
			// 50% chance of changing the website
			updatedWebsite = simtypes.RandStringOfLength(r, 20)
		}

		updatedPictureURL := types.DoNotModify
		if r.Intn(100) < 50 {
			// 50% chance of changing the picture URL
			updatedPictureURL = simtypes.RandStringOfLength(r, 20)
		}

		// Get the admin account that should sign the transaction
		adminAddress, err := sdk.AccAddressFromBech32(operator.Admin)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "invalid admin address"), nil, nil
		}

		signer, found := simtesting.GetSimAccount(adminAddress, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "admin account not found"), nil, nil
		}

		// Create the message
		msg = types.NewMsgUpdateOperator(
			operator.ID,
			updatedMoniker,
			updatedWebsite,
			updatedPictureURL,
			operator.Admin,
		)
		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, signer)
	}
}

func SimulateDeactivateOperator(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k *keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var msg = &types.MsgDeactivateOperator{}

		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		// Get a random operator from the one existing
		operator, found := GetRandomExistingOperator(r, ctx, k, func(operator types.Operator) bool {
			return operator.Status == types.OPERATOR_STATUS_ACTIVE
		})
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get operator"), nil, nil
		}

		// Get the admin account that should sign the transaction
		adminAddress, err := sdk.AccAddressFromBech32(operator.Admin)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "invalid admin address"), nil, nil
		}

		signer, found := simtesting.GetSimAccount(adminAddress, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "admin account not found"), nil, nil
		}

		// Create the message
		msg = types.NewMsgDeactivateOperator(operator.ID, operator.Admin)
		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, signer)
	}
}

func SimulateReactivateOperator(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k *keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var msg = &types.MsgReactivateOperator{}

		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		// Get a random operator from the one existing
		operator, found := GetRandomExistingOperator(r, ctx, k, func(operator types.Operator) bool {
			return operator.Status == types.OPERATOR_STATUS_INACTIVE
		})
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get inactive operator"), nil, nil
		}

		// Get the admin account that should sign the transaction
		adminAddress, err := sdk.AccAddressFromBech32(operator.Admin)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "invalid admin address"), nil, nil
		}

		signer, found := simtesting.GetSimAccount(adminAddress, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "admin account not found"), nil, nil
		}

		// Create the message
		msg = types.NewMsgReactivateOperator(operator.ID, operator.Admin)
		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, signer)
	}
}

func SimulateTransferOperatorOwnership(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k *keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var msg = &types.MsgTransferOperatorOwnership{}

		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		// Get a random operator from the one existing
		operator, found := GetRandomExistingOperator(r, ctx, k, nil)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get operator"), nil, nil
		}

		// Get a random new admin
		newAdminAccount, _ := simtypes.RandomAcc(r, accs)

		// Get the current admin account that should sign the transaction
		adminAddress, err := sdk.AccAddressFromBech32(operator.Admin)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "invalid admin address"), nil, nil
		}

		signer, found := simtesting.GetSimAccount(adminAddress, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "admin account not found"), nil, nil
		}

		// Create the message
		msg = types.NewMsgTransferOperatorOwnership(operator.ID, newAdminAccount.Address.String(), operator.Admin)
		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, signer)
	}
}

func SimulateDeleteOperator(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k *keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var msg = &types.MsgDeleteOperator{}

		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		// Get a random operator from the one existing
		operator, found := GetRandomExistingOperator(r, ctx, k, func(operator types.Operator) bool {
			return operator.Status == types.OPERATOR_STATUS_INACTIVE
		})
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get operator"), nil, nil
		}

		// Get the admin account that should sign the transaction
		adminAddress, err := sdk.AccAddressFromBech32(operator.Admin)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "invalid admin address"), nil, nil
		}

		signer, found := simtesting.GetSimAccount(adminAddress, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "admin account not found"), nil, nil
		}

		// Create the message
		msg = types.NewMsgDeleteOperator(operator.ID, operator.Admin)
		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, signer)
	}
}

func SimulateSetOperatorParams(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k *keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var msg = &types.MsgSetOperatorParams{}

		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		// Get a random operator from the one existing
		operator, found := GetRandomExistingOperator(r, ctx, k, nil)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get operator"), nil, nil
		}

		// Generate a new data
		newParams := RandomOperatorParams(r)

		// Get the admin account that should sign the transaction
		adminAddress, err := sdk.AccAddressFromBech32(operator.Admin)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "invalid admin address"), nil, nil
		}

		signer, found := simtesting.GetSimAccount(adminAddress, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "admin account not found"), nil, nil
		}

		// Create the message
		msg = types.NewMsgSetOperatorParams(operator.ID, newParams, operator.Admin)
		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, signer)
	}
}
