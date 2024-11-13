package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/milkyway-labs/milkyway/testutils/simtesting"
	"github.com/milkyway-labs/milkyway/x/services/keeper"
	"github.com/milkyway-labs/milkyway/x/services/types"
)

// Simulation operation weights constants
const (
	DefaultWeightMsgCreateService            int = 100
	DefaultWeightMsgUpdateService            int = 100
	DefaultWeightMsgActivateService          int = 100
	DefaultWeightMsgDeactivateService        int = 100
	DefaultWeightMsgTransferServiceOwnership int = 100
	DefaultWeightMsgDeleteService            int = 100
	DefaultWeightMsgSetServiceParams         int = 100

	OpWeightMsgCreateService            = "op_weight_msg_create_service"
	OpWeightMsgUpdateService            = "op_weight_msg_update_service"
	OpWeightMsgActivateService          = "op_weight_msg_activate_service"
	OpWeightMsgDeactivateService        = "op_weight_msg_deactivate_service"
	OpWeightMsgTransferServiceOwnership = "op_weight_msg_transfer_service_ownership"
	OpWeightMsgDeleteService            = "op_weight_msg_delete_service"
	OpWeightMsgSetServiceParams         = "op_weight_msg_set_service_params"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams,
	cdc codec.JSONCodec,
	txGen client.TxConfig,
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
	appParams.GetOrGenerate(OpWeightMsgCreateService, &weightMsgCreateService, nil, func(_ *rand.Rand) {
		weightMsgCreateService = DefaultWeightMsgCreateService
	})

	appParams.GetOrGenerate(OpWeightMsgUpdateService, &weightMsgUpdateService, nil, func(_ *rand.Rand) {
		weightMsgUpdateService = DefaultWeightMsgUpdateService
	})

	appParams.GetOrGenerate(OpWeightMsgActivateService, &weightMsgActivateService, nil, func(_ *rand.Rand) {
		weightMsgActivateService = DefaultWeightMsgActivateService
	})

	appParams.GetOrGenerate(OpWeightMsgDeactivateService, &weightMsgDeactivateService, nil, func(_ *rand.Rand) {
		weightMsgDeactivateService = DefaultWeightMsgDeactivateService
	})

	appParams.GetOrGenerate(OpWeightMsgTransferServiceOwnership, &weightMsgTransferServiceOwnership, nil, func(_ *rand.Rand) {
		weightMsgTransferServiceOwnership = DefaultWeightMsgTransferServiceOwnership
	})

	appParams.GetOrGenerate(OpWeightMsgDeleteService, &weightMsgDeleteService, nil, func(_ *rand.Rand) {
		weightMsgDeleteService = DefaultWeightMsgDeleteService
	})

	appParams.GetOrGenerate(OpWeightMsgSetServiceParams, &weightMsgSetServiceParams, nil, func(_ *rand.Rand) {
		weightMsgSetServiceParams = DefaultWeightMsgSetServiceParams
	})

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(weightMsgCreateService, SimulateMsgCreateService(txGen, ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgUpdateService, SimulateMsgUpdateService(txGen, ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgActivateService, SimulateMsgActivateService(txGen, ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgDeactivateService, SimulateMsgDeactivateService(txGen, ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgTransferServiceOwnership, SimulateMsgTransferServiceOwnership(txGen, ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgDeleteService, SimulateMsgDeleteService(txGen, ak, bk, k)),
		simulation.NewWeightedOperation(weightMsgSetServiceParams, SimulateMsgSetServiceParams(txGen, ak, bk, k)),
	}
}

func SimulateMsgCreateService(
	txGen client.TxConfig,
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, "MsgCreateService", "skip"), nil, nil
		}

		signer, _ := simtypes.RandomAcc(r, accs)
		service := RandomService(r, 1, signer.Address.String())
		msg := types.NewMsgCreateService(
			service.Name,
			service.Description,
			service.Website,
			service.PictureURL,
			service.Admin,
		)

		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, signer)
	}
}

func SimulateMsgUpdateService(
	txGen client.TxConfig,
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, "MsgUpdateService", "skip"), nil, nil
		}

		// Get a random service to update
		service, found := GetRandomExistingService(r, ctx, k, nil)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, "MsgUpdateService", "skip"), nil, nil
		}

		// Get the service admin sim account
		adminAddr := sdk.MustAccAddressFromBech32(service.Admin)
		simAccount, found := simtesting.GetSimAccount(adminAddr, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, "service admin not found", "skip"), nil, nil
		}

		// Generate the new service fields
		newService := RandomService(r, service.ID, service.Admin)
		// Create the msg
		msg := types.NewMsgUpdateService(
			service.ID,
			newService.Name,
			newService.Description,
			newService.Website,
			newService.PictureURL,
			simAccount.Address.String(),
		)

		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, simAccount)
	}
}

func SimulateMsgActivateService(
	txGen client.TxConfig,
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, "MsgActivateService", "skip"), nil, nil
		}

		// Get a random service to activate
		service, found := GetRandomExistingService(r, ctx, k, func(s types.Service) bool {
			return s.Status == types.SERVICE_STATUS_CREATED || s.Status == types.SERVICE_STATUS_INACTIVE
		})
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, "MsgActivateService", "skip"), nil, nil
		}

		// Get the service admin sim account
		adminAddr := sdk.MustAccAddressFromBech32(service.Admin)
		simAccount, found := simtesting.GetSimAccount(adminAddr, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, "service admin not found", "skip"), nil, nil
		}

		// Create the msg
		msg := types.NewMsgActivateService(service.ID, simAccount.Address.String())

		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, simAccount)
	}
}

func SimulateMsgDeactivateService(
	txGen client.TxConfig,
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, "MsgDeactivateService", "skip"), nil, nil
		}

		// Get a random service
		service, found := GetRandomExistingService(r, ctx, k, func(s types.Service) bool {
			return s.Status == types.SERVICE_STATUS_ACTIVE
		})
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, "MsgDeactivateService", "skip"), nil, nil
		}

		// Get the service admin sim account
		adminAddr := sdk.MustAccAddressFromBech32(service.Admin)
		simAccount, found := simtesting.GetSimAccount(adminAddr, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, "service admin not found", "skip"), nil, nil
		}

		// Create the msg
		msg := types.NewMsgDeactivateService(service.ID, simAccount.Address.String())

		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, simAccount)
	}
}

func SimulateMsgTransferServiceOwnership(
	txGen client.TxConfig,
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, "MsgTransferServiceOwnership", "skip"), nil, nil
		}

		// Get a random service
		service, found := GetRandomExistingService(r, ctx, k, nil)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, "MsgTransferServiceOwnership", "skip"), nil, nil
		}

		// Get the service admin sim account
		adminAddr := sdk.MustAccAddressFromBech32(service.Admin)
		simAccount, found := simtesting.GetSimAccount(adminAddr, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, "service admin not found", "skip"), nil, nil
		}

		// Get a new admin
		newAdminAccount, _ := simtypes.RandomAcc(r, accs)

		// Create the msg
		msg := types.NewMsgTransferServiceOwnership(service.ID, newAdminAccount.Address.String(), simAccount.Address.String())

		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, simAccount)
	}
}

func SimulateMsgDeleteService(
	txGen client.TxConfig,
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, "MsgDeleteService", "skip"), nil, nil
		}

		// Get a random service
		service, found := GetRandomExistingService(r, ctx, k, func(s types.Service) bool {
			return s.Status == types.SERVICE_STATUS_INACTIVE
		})
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, "MsgDeleteService", "skip"), nil, nil
		}

		// Get the service admin sim account
		adminAddr := sdk.MustAccAddressFromBech32(service.Admin)
		simAccount, found := simtesting.GetSimAccount(adminAddr, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, "service admin not found", "skip"), nil, nil
		}

		// Create the msg
		msg := types.NewMsgDeleteService(service.ID, simAccount.Address.String())

		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, simAccount)
	}
}

func SimulateMsgSetServiceParams(
	txGen client.TxConfig,
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	k *keeper.Keeper,
) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// No account skipping
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, "MsgSetServiceParams", "skip"), nil, nil
		}

		// Get a random service
		service, found := GetRandomExistingService(r, ctx, k, nil)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, "MsgSetServiceParams", "skip"), nil, nil
		}

		// Get the service admin sim account
		adminAddr := sdk.MustAccAddressFromBech32(service.Admin)
		simAccount, found := simtesting.GetSimAccount(adminAddr, accs)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, "service admin not found", "skip"), nil, nil
		}

		serviceParams := types.DefaultServiceParams()

		// Create the msg
		msg := types.NewMsgSetServiceParams(service.ID, serviceParams, service.Admin)

		return simtesting.SendMsg(r, types.ModuleName, app, ak, bk, msg, ctx, simAccount)
	}
}
