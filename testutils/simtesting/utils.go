package simtesting

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// SendMsg sends a transaction with the specified message.
func SendMsg(
	r *rand.Rand, moduleName string, app *baseapp.BaseApp, ak authkeeper.AccountKeeper, bk bankkeeper.Keeper,
	msg sdk.Msg, ctx sdk.Context,
	simAccount simtypes.Account,
) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
	deposit := sdk.Coins{}
	spendableCoins := bk.SpendableCoins(ctx, simAccount.Address)
	for _, v := range spendableCoins {
		if bk.IsSendEnabledCoin(ctx, v) {
			deposit = deposit.Add(simtypes.RandSubsetCoins(r, sdk.NewCoins(v))...)
		}
	}

	if deposit.IsZero() {
		msgType := sdk.MsgTypeURL(msg)
		return simtypes.NoOpMsg(moduleName, msgType, "skip because of broke account"), nil, nil
	}

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	txConfig := tx.NewTxConfig(codec.NewProtoCodec(interfaceRegistry), tx.DefaultSignModes)
	txCtx := simulation.OperationInput{
		R:               r,
		App:             app,
		TxGen:           txConfig,
		Cdc:             nil,
		Msg:             msg,
		Context:         ctx,
		SimAccount:      simAccount,
		AccountKeeper:   ak,
		Bankkeeper:      bk,
		ModuleName:      moduleName,
		CoinsSpentInMsg: deposit,
	}
	return simulation.GenAndDeliverTxWithRandFees(txCtx)
}

// GetSimAccount gets the Account with the given address
func GetSimAccount(address sdk.Address, accs []simtypes.Account) (simtypes.Account, bool) {
	for _, acc := range accs {
		if acc.Address.Equals(address) {
			return acc, true
		}
	}
	return simtypes.Account{}, false
}
