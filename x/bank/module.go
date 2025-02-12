package bank

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/bank/exported"
	"github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/milkyway-labs/milkyway/v7/x/bank/keeper"
)

const ConsensusVersion = 1

// AppModule implements an application module for the bank module.
type AppModule struct {
	keeper keeper.Keeper
	bank.AppModule
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper, accountKeeper types.AccountKeeper, ss exported.Subspace) AppModule {
	return AppModule{
		keeper:    keeper,
		AppModule: bank.NewAppModule(cdc, keeper.BaseKeeper, accountKeeper, ss),
	}
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}

// ConsensusVersion implements ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return ConsensusVersion }
