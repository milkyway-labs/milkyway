package distribution

import (
	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/distribution/exported"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/milkyway-labs/milkyway/v7/x/distribution/keeper"
)

const ConsensusVersion = 1

var (
	_ module.AppModuleBasic = AppModule{}
	_ appmodule.AppModule   = AppModule{}
)

type AppModule struct {
	distribution.AppModule
}

func NewAppModule(
	cdc codec.Codec, keeper keeper.Keeper, accountKeeper distrtypes.AccountKeeper,
	bankKeeper distrtypes.BankKeeper, stakingKeeper distrtypes.StakingKeeper, ss exported.Subspace,
) AppModule {
	return AppModule{
		AppModule: distribution.NewAppModule(cdc, keeper.Keeper, accountKeeper, bankKeeper, stakingKeeper, ss),
	}
}

func (am AppModule) ConsensusVersion() uint64 { return ConsensusVersion }
