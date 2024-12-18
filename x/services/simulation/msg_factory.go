package simulation

import (
	"bytes"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/milkyway-labs/milkyway/v6/testutils/simtesting"
	"github.com/milkyway-labs/milkyway/v6/x/services/keeper"
	"github.com/milkyway-labs/milkyway/v6/x/services/types"
)

// Simulation operation weights constants
const (
	DefaultWeightMsgCreateService            int = 80
	DefaultWeightMsgUpdateService            int = 30
	DefaultWeightMsgActivateService          int = 60
	DefaultWeightMsgDeactivateService        int = 20
	DefaultWeightMsgTransferServiceOwnership int = 15
	DefaultWeightMsgDeleteService            int = 10
	DefaultWeightMsgSetServiceParams         int = 40

	OperationWeightMsgCreateService            = "op_weight_msg_create_service"
	OperationWeightMsgUpdateService            = "op_weight_msg_update_service"
	OperationWeightMsgActivateService          = "op_weight_msg_activate_service"
	OperationWeightMsgDeactivateService        = "op_weight_msg_deactivate_service"
	OperationWeightMsgTransferServiceOwnership = "op_weight_msg_transfer_service_ownership"
	OperationWeightMsgDeleteService            = "op_weight_msg_delete_service"
	OperationWeightMsgSetServiceParams         = "op_weight_msg_set_service_params"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams,
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	k *keeper.Keeper,
) simulation.WeightedOperations {
	var (
		weightMsgCreateService            int
		weightMsgUpdateService            int
		weightMsgActivateService          int
		weightMsgDeactivateService        int
		weightMsgTransferServiceOwnership int
		weightMsgDeleteService            int
		weightMsgSetServiceParams         int
	)

	// Generate the weights
	appParams.GetOrGenerate(OperationWeightMsgCreateService, &weightMsgCreateService, nil, func(_ *rand.Rand) {
		weightMsgCreateService = DefaultWeightMsgCreateService
	})

	appParams.GetOrGenerate(OperationWeightMsgUpdateService, &weightMsgUpdateService, nil, func(_ *rand.Rand) {
		weightMsgUpdateService = DefaultWeightMsgUpdateService
	})

	appParams.GetOrGenerate(OperationWeightMsgActivateService, &weightMsgActivateService, nil, func(_ *rand.Rand) {
		weightMsgActivateService = DefaultWeightMsgActivateService
	})

	appParams.GetOrGenerate(OperationWeightMsgDeactivateService, &weightMsgDeactivateService, nil, func(_ *rand.Rand) {
		weightMsgDeactivateService = DefaultWeightMsgDeactivateService
	})

	appParams.GetOrGenerate(OperationWeightMsgTransferServiceOwnership, &weightMsgTransferServiceOwnership, nil, func(_ *rand.Rand) {
		weightMsgTransferServiceOwnership = DefaultWeightMsgTransferServiceOwnership
	})

	appParams.GetOrGenerate(OperationWeightMsgDeleteService, &weightMsgDeleteService, nil, func(_ *rand.Rand) {
		weightMsgDeleteService = DefaultWeightMsgDeleteService
	})

	appParams.GetOrGenerate(OperationWeightMsgSetServiceParams, &weightMsgSetServiceParams, nil, func(_ *rand.Rand) {
		weightMsgSetServiceParams = DefaultWeightMsgSetServiceParams
	})

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(weightMsgCreateService, SimulateMsgCreateService(ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgUpdateService, SimulateMsgUpdateService(ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgActivateService, SimulateMsgActivateService(ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgDeactivateService, SimulateMsgDeactivateService(ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgTransferServiceOwnership, SimulateMsgTransferServiceOwnership(ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgDeleteService, SimulateMsgDeleteService(ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgSetServiceParams, SimulateMsgSetServiceParams(ak, bk, k)),
	}
}

func SimulateMsgCreateService(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k *keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var msg = &types.MsgCreateService{}

		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		// Create a random service
		service := RandomService(r, accs)

		// Get the admin account that should sign the transaction
		adminAddress, err := sdk.AccAddressFromBech32(service.Admin)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "error while parsing admin address"), nil, nil
		}

		// Make sure the admin has enough funds to pay for the creation fees
		params, err := k.GetParams(ctx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "could not get params"), nil, nil
		}

		// Check the fees that should be paid
		feesAmount := sdk.NewCoins()
		for _, feeCoin := range params.ServiceRegistrationFee {
			if bk.GetBalance(ctx, adminAddress, feeCoin.Denom).IsGTE(feeCoin) {
				feesAmount = feesAmount.Add(feeCoin)
				break
			}
		}

		if !params.ServiceRegistrationFee.IsZero() && feesAmount.IsZero() {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "insufficient funds"), nil, nil
		}

		// Get the account that will sign the transaction
		signer, found := simtesting.GetSimAccount(adminAddress, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "admin account not found"), nil, nil
		}

		// Create the message
		msg = types.NewMsgCreateService(
			service.Name,
			service.Description,
			service.Website,
			service.PictureURL,
			feesAmount,
			service.Admin,
		)
		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, signer)
	}
}

func SimulateMsgUpdateService(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k *keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var msg = &types.MsgUpdateService{}

		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		// Get a random service to update
		service, found := GetRandomExistingService(r, ctx, k, nil)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no service found"), nil, nil
		}

		// Get the service admin sim account
		adminAddr, err := sdk.AccAddressFromBech32(service.Admin)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "error while parsing admin address"), nil, nil
		}

		// Get the params to be updates
		updatedName := types.DoNotModify
		if r.Intn(100) < 50 {
			// 50% chance of updating the name
			updatedName = simtypes.RandStringOfLength(r, 24)
		}

		updatedDescription := types.DoNotModify
		if r.Intn(100) < 50 {
			// 50% chance of updating the description
			updatedDescription = simtypes.RandStringOfLength(r, 24)
		}

		updatedWebsite := types.DoNotModify
		if r.Intn(100) < 50 {
			// 50% chance of updating the website
			updatedWebsite = simtypes.RandStringOfLength(r, 24)
		}

		updatedPictureURL := types.DoNotModify
		if r.Intn(100) < 50 {
			// 50% chance of updating the picture URL
			updatedPictureURL = simtypes.RandStringOfLength(r, 24)
		}

		signer, found := simtesting.GetSimAccount(adminAddr, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, "service admin not found", "skip"), nil, nil
		}

		// Create the msg
		msg = types.NewMsgUpdateService(
			service.ID,
			updatedName,
			updatedDescription,
			updatedWebsite,
			updatedPictureURL,
			signer.Address.String(),
		)
		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, signer)
	}
}

func SimulateMsgActivateService(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k *keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var msg = &types.MsgActivateService{}

		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "skip"), nil, nil
		}

		// Get a random service to activate
		service, found := GetRandomExistingService(r, ctx, k, func(s types.Service) bool {
			return s.Status == types.SERVICE_STATUS_CREATED || s.Status == types.SERVICE_STATUS_INACTIVE
		})
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "skip"), nil, nil
		}

		// Get the service admin sim account
		adminAddr, err := sdk.AccAddressFromBech32(service.Admin)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "error while parsing admin address"), nil, nil
		}

		simAccount, found := simtesting.GetSimAccount(adminAddr, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "service admin not found"), nil, nil
		}

		// Create the msg
		msg = types.NewMsgActivateService(service.ID, simAccount.Address.String())
		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, simAccount)
	}
}

func SimulateMsgDeactivateService(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k *keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var msg = &types.MsgDeactivateService{}

		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		// Get a random service
		service, found := GetRandomExistingService(r, ctx, k, func(s types.Service) bool {
			return s.Status == types.SERVICE_STATUS_ACTIVE
		})
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "service not found"), nil, nil
		}

		// Get the service admin sim account
		adminAddr, err := sdk.AccAddressFromBech32(service.Admin)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "error while parsing admin address"), nil, nil
		}

		simAccount, found := simtesting.GetSimAccount(adminAddr, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "admin account not found"), nil, nil
		}

		// Create the msg
		msg = types.NewMsgDeactivateService(service.ID, simAccount.Address.String())
		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, simAccount)
	}
}

func SimulateMsgTransferServiceOwnership(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k *keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var msg = &types.MsgTransferServiceOwnership{}

		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		// Get a random service
		service, found := GetRandomExistingService(r, ctx, k, nil)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "service not found"), nil, nil
		}

		// Get the service admin sim account
		adminAddr, err := sdk.AccAddressFromBech32(service.Admin)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "error while parsing admin address"), nil, nil
		}

		// Get a new admin
		newAdminAccount, _ := simtypes.RandomAcc(r, accs)
		if bytes.Equal(newAdminAccount.Address, adminAddr) {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "new admin is the same as the current one"), nil, nil
		}

		// Get the signer
		signer, found := simtesting.GetSimAccount(adminAddr, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "service admin not found"), nil, nil
		}

		// Create the msg
		msg = types.NewMsgTransferServiceOwnership(service.ID, newAdminAccount.Address.String(), signer.Address.String())
		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, signer)
	}
}

func SimulateMsgDeleteService(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k *keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var msg = &types.MsgDeleteService{}

		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		// Get a random service
		service, found := GetRandomExistingService(r, ctx, k, func(s types.Service) bool {
			return s.Status == types.SERVICE_STATUS_INACTIVE
		})
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "service not found"), nil, nil
		}

		// Get the service admin sim account
		adminAddr, err := sdk.AccAddressFromBech32(service.Admin)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "error while parsing admin address"), nil, nil
		}

		signer, found := simtesting.GetSimAccount(adminAddr, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "service admin not found"), nil, nil
		}

		// Create the msg
		msg = types.NewMsgDeleteService(service.ID, signer.Address.String())
		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, signer)
	}
}

func SimulateMsgSetServiceParams(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k *keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var msg = &types.MsgSetServiceParams{}

		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "no accounts"), nil, nil
		}

		// Get a random service
		service, found := GetRandomExistingService(r, ctx, k, nil)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "service not found"), nil, nil
		}

		// Get new random params
		serviceParams := RandomServiceParams(r)

		// Get the service admin sim account
		adminAddr, err := sdk.AccAddressFromBech32(service.Admin)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "error while parsing admin address"), nil, nil
		}

		signer, found := simtesting.GetSimAccount(adminAddr, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "service admin not found"), nil, nil
		}

		// Create the msg
		msg = types.NewMsgSetServiceParams(service.ID, serviceParams, service.Admin)
		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, signer)
	}
}
