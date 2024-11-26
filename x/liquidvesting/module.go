package liquidvesting

import (
	"context"
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/milkyway-labs/milkyway/v2/x/liquidvesting/keeper"
	"github.com/milkyway-labs/milkyway/v2/x/liquidvesting/types"
)

const (
	consensusVersion = 1
)

var (
	_ appmodule.AppModule   = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// ----------------------------------------------------------------------------
// AppModuleBasic
// ----------------------------------------------------------------------------

// AppModuleBasic implements the AppModuleBasic interface for the
// liquidvesting module.
type AppModuleBasic struct {
	cdc codec.BinaryCodec
}

func NewAppModuleBasic(cdc codec.BinaryCodec) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

// Name implements module.AppModuleBasic.
func (a AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterInterfaces implements module.AppModuleBasic.
func (a AppModuleBasic) RegisterInterfaces(reg cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(reg)
}

// RegisterLegacyAminoCodec implements module.AppModuleBasic.
func (a AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

// DefaultGenesis returns the module's default genesis state.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return genState.Validate()
}

// RegisterGRPCGatewayRoutes implements module.AppModuleBasic.
func (a AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
	if err != nil {
		panic(err)
	}
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule implements the AppModule interface for the
// liquidvesting module.
type AppModule struct {
	AppModuleBasic

	keeper *keeper.Keeper
}

func NewAppModule(cdc codec.Codec, keeper *keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(cdc),
		keeper:         keeper,
	}
}

// Name returns the module's name.
func (a AppModule) Name() string {
	return a.AppModuleBasic.Name()
}

// RegisterServices registers a GRPC query service to respond to the module-specific GRPC queries.
func (a AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServer(a.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQuerier(a.keeper))
}

// InitGenesis performs the module's genesis initialization It returns no validator updates.
func (a AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) {
	var genState types.GenesisState
	cdc.MustUnmarshalJSON(gs, &genState)

	err := a.keeper.InitGenesis(ctx, &genState)
	if err != nil {
		panic(err)
	}
}

// ExportGenesis returns the module's exported genesis state as raw JSON bytes.
func (a AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genState, err := a.keeper.ExportGenesis(ctx)
	if err != nil {
		panic(err)
	}

	return cdc.MustMarshalJSON(genState)
}

// ConsensusVersion implements ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return consensusVersion }

// IsAppModule implements appmodule.AppModule.
func (a AppModule) IsAppModule() {}

// IsOnePerModuleType implements appmodule.AppModule.
func (a AppModule) IsOnePerModuleType() {}

// EndBlock executes all ABCI EndBlock logic respective to the restaking module.
func (a AppModule) EndBlock(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return EndBlocker(sdkCtx, a.keeper)
}
