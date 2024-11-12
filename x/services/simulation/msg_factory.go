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
	DefaultWeightMsgCreateService int = 100

	OpWeightMsgCreateService = "op_weight_msg_create_service"
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
	var weightMsgCreateService int

	appParams.GetOrGenerate(OpWeightMsgCreateService, &weightMsgCreateService, nil, func(_ *rand.Rand) {
		weightMsgCreateService = DefaultWeightMsgCreateService
	})

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(weightMsgCreateService, SimulateMsgCreateService(txGen, ak, bk, k)),
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
