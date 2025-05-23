package milkyway

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"
	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/log"
	"cosmossdk.io/x/tx/signing"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	abci "github.com/cometbft/cometbft/abci/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	tmos "github.com/cometbft/cometbft/libs/os"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	sigtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	txmodule "github.com/cosmos/cosmos-sdk/x/auth/tx/config"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/gogoproto/proto"
	ibctesting "github.com/cosmos/ibc-go/v8/testing"
	providertypes "github.com/cosmos/interchain-security/v6/x/ccv/provider/types"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	oracleconfig "github.com/skip-mev/connect/v2/oracle/config"
	oracleclient "github.com/skip-mev/connect/v2/service/clients/oracle"
	"github.com/spf13/cast"

	milkywayante "github.com/milkyway-labs/milkyway/v12/ante"
	milkywayabci "github.com/milkyway-labs/milkyway/v12/app/abci"
	"github.com/milkyway-labs/milkyway/v12/app/forks"
	v7 "github.com/milkyway-labs/milkyway/v12/app/forks/v7"
	v8 "github.com/milkyway-labs/milkyway/v12/app/forks/v8"
	"github.com/milkyway-labs/milkyway/v12/app/keepers"
	"github.com/milkyway-labs/milkyway/v12/app/upgrades"
	v10 "github.com/milkyway-labs/milkyway/v12/app/upgrades/v10"
	v11 "github.com/milkyway-labs/milkyway/v12/app/upgrades/v11"
	v11warpfix "github.com/milkyway-labs/milkyway/v12/app/upgrades/v11-warp-fix"
	v12 "github.com/milkyway-labs/milkyway/v12/app/upgrades/v12"
	v12commissionfix "github.com/milkyway-labs/milkyway/v12/app/upgrades/v12-commission-fix"
	v6 "github.com/milkyway-labs/milkyway/v12/app/upgrades/v6"
	v9 "github.com/milkyway-labs/milkyway/v12/app/upgrades/v9"
	_ "github.com/milkyway-labs/milkyway/v12/client/docs/statik"
	liquidvestingtypes "github.com/milkyway-labs/milkyway/v12/x/liquidvesting/types"
)

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string

	Forks = []forks.Fork{
		v7.Fork,
		v8.Fork,
	}

	Upgrades = []upgrades.Upgrade{
		v6.Upgrade,
		v9.Upgrade,
		v10.Upgrade,
		v11.Upgrade,
		v11warpfix.Upgrade,
		v12.Upgrade,
		v12commissionfix.Upgrade,
	}
)

var (
	_ runtime.AppI            = (*MilkyWayApp)(nil)
	_ servertypes.Application = (*MilkyWayApp)(nil)
	_ ibctesting.TestingApp   = (*MilkyWayApp)(nil)
)

// MilkyWayApp extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type MilkyWayApp struct {
	*baseapp.BaseApp
	keepers.AppKeepers

	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// external fields
	oracleClient oracleclient.OracleClient

	// the module manager
	mm           *module.Manager
	ModuleBasics module.BasicManager

	// the configurator
	sm           *module.SimulationManager
	configurator module.Configurator
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, "."+appName)
}

// NewMilkyWayApp returns a reference to an initialized MilkyWayApp.
func NewMilkyWayApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	oracleConfig oracleconfig.AppConfig,
	appOpts servertypes.AppOptions,
	wasmOpts []wasmkeeper.Option,
	baseAppOptions ...func(*baseapp.BaseApp),
) *MilkyWayApp {
	legacyAmino := codec.NewLegacyAmino()
	interfaceRegistry, err := types.NewInterfaceRegistryWithOptions(types.InterfaceRegistryOptions{
		ProtoFiles: proto.HybridResolver,
		SigningOptions: signing.Options{
			AddressCodec: address.Bech32Codec{
				Bech32Prefix: sdk.GetConfig().GetBech32AccountAddrPrefix(),
			},
			ValidatorAddressCodec: address.Bech32Codec{
				Bech32Prefix: sdk.GetConfig().GetBech32ValidatorAddrPrefix(),
			},
		},
	})
	if err != nil {
		panic(err)
	}

	encodingCfg := MakeEncodingConfig()
	appCodec := encodingCfg.Marshaler
	txConfig := encodingCfg.TxConfig

	std.RegisterLegacyAminoCodec(legacyAmino)
	std.RegisterInterfaces(interfaceRegistry)

	// App Opts
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))
	invCheckPeriod := cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod))

	bApp := baseapp.NewBaseApp(
		appName,
		logger,
		db,
		txConfig.TxDecoder(),
		baseAppOptions...,
	)

	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)
	bApp.SetTxEncoder(txConfig.TxEncoder())

	app := &MilkyWayApp{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		txConfig:          txConfig,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
	}

	moduleAccountAddresses := ModuleAccountAddrs()

	// Setup keepers
	app.AppKeepers = keepers.NewAppKeeper(
		appCodec,
		bApp,
		legacyAmino,
		MaccPerms,
		BlockedModuleAccountAddrs(moduleAccountAddresses),
		skipUpgradeHeights,
		homePath,
		invCheckPeriod,
		logger,
		appOpts,
		wasmOpts,
	)

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(appModules(app, appCodec, txConfig, skipGenesisInvariants)...)
	app.ModuleBasics = newBasicManagerFromManager(app)

	enabledSignModes := append([]sigtypes.SignMode(nil), authtx.DefaultSignModes...)
	enabledSignModes = append(enabledSignModes, sigtypes.SignMode_SIGN_MODE_TEXTUAL)

	txConfigOpts := authtx.ConfigOptions{
		EnabledSignModes:           enabledSignModes,
		TextualCoinMetadataQueryFn: txmodule.NewBankKeeperCoinMetadataQueryFn(app.BankKeeper),
	}
	txConfig, err = authtx.NewTxConfigWithOptions(
		appCodec,
		txConfigOpts,
	)
	if err != nil {
		panic(err)
	}
	app.txConfig = txConfig

	// NOTE: upgrade module is required to be prioritized
	app.mm.SetOrderPreBlockers(
		upgradetypes.ModuleName,
	)
	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	// NOTE: capability module's beginblocker must come before any modules using capabilities (e.g. IBC)
	// Tell the app's module manager how to set the order of BeginBlockers, which are run at the beginning of every block.
	app.mm.SetOrderBeginBlockers(orderBeginBlockers()...)

	app.mm.SetOrderEndBlockers(orderEndBlockers()...)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: The genutils module must also occur after auth so that it can access the params from auth.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.mm.SetOrderInitGenesis(orderInitBlockers()...)

	// Uncomment if you want to set a custom migration order here.
	// app.mm.SetOrderMigrations(custom order)

	app.mm.RegisterInvariants(app.CrisisKeeper)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	err = app.mm.RegisterServices(app.configurator)
	if err != nil {
		panic(err)
	}

	autocliv1.RegisterQueryServer(app.GRPCQueryRouter(), runtimeservices.NewAutoCLIQueryService(app.mm.Modules))

	reflectionSvc, err := runtimeservices.NewReflectionService()
	if err != nil {
		panic(err)
	}
	reflectionv1.RegisterReflectionServiceServer(app.GRPCQueryRouter(), reflectionSvc)

	// add test gRPC service for testing gRPC queries in isolation
	testdata.RegisterQueryServer(app.GRPCQueryRouter(), testdata.QueryImpl{})

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	app.sm = module.NewSimulationManager(simulationModules(app, appCodec, skipGenesisInvariants)...)

	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(app.GetKVStoreKey())
	app.MountTransientStores(app.GetTransientStoreKey())
	app.MountMemoryStores(app.GetMemoryStoreKey())

	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic("error while reading wasm config: " + err.Error())
	}

	anteHandler, err := milkywayante.NewAnteHandler(
		milkywayante.HandlerOptions{
			AccountKeeper:   app.AccountKeeper,
			BankKeeper:      app.BankKeeper,
			FeegrantKeeper:  app.FeeGrantKeeper,
			SignModeHandler: txConfig.SignModeHandler(),
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,

			Codec:                 appCodec,
			IBCKeeper:             app.IBCKeeper,
			StakingKeeper:         app.StakingKeeper,
			WasmConfig:            &wasmConfig,
			TXCounterStoreService: runtime.NewKVStoreService(app.AppKeepers.GetKey(wasmtypes.StoreKey)),
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to create AnteHandler: %s", err))
	}

	// set ante and post handlers
	app.SetAnteHandler(anteHandler)

	// Initialize the ABCI extensions
	data := milkywayabci.InitializeOracleABCIExtensions(milkywayabci.SetupData{
		ChainID:       app.ChainID(),
		Logger:        app.Logger(),
		Keepers:       app.AppKeepers,
		ModuleManager: app.mm,
		OracleConfig:  oracleConfig,
	})

	// Set the oracle client
	app.SetOracleClient(data.OracleClient)

	// Set the standard callbacks
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	// Set the callbacks used by ABCI vote extensions
	app.SetPrepareProposal(data.PrepareProposalHandler)
	app.SetProcessProposal(data.ProcessProposalHandler)
	app.SetPreBlocker(data.PreBlockHandler)
	app.SetExtendVoteHandler(data.ExtendVoteHandler)
	app.SetVerifyVoteExtensionHandler(data.VerifyVoteExtensionHandler)

	if manager := app.SnapshotManager(); manager != nil {
		err = manager.RegisterExtensions(wasmkeeper.NewWasmSnapshotter(app.CommitMultiStore(), &app.AppKeepers.WasmKeeper))
		if err != nil {
			panic("failed to register snapshot extension: " + err.Error())
		}
	}

	app.setupUpgradeHandlers()
	app.setupUpgradeStoreLoaders()
	app.setupForksPreBlockers()

	// At startup, after all modules have been registered, check that all prot
	// annotations are correct.
	protoFiles, err := proto.MergedRegistry()
	if err != nil {
		panic(err)
	}
	err = msgservice.ValidateProtoAnnotations(protoFiles)
	if err != nil {
		// Once we switch to using protoreflect-based antehandlers, we might
		// want to panic here instead of logging a warning.
		fmt.Fprintln(os.Stderr, err.Error())
	}

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(fmt.Sprintf("failed to load latest version: %s", err))
		}

		ctx := app.BaseApp.NewUncachedContext(true, tmproto.Header{})

		if err := app.AppKeepers.WasmKeeper.InitializePinnedCodes(ctx); err != nil {
			tmos.Exit(fmt.Sprintf("WasmKeeper failed initialize pinned codes %s", err))
		}
	}

	return app
}

// SetOracleClient sets the oracle client
func (app *MilkyWayApp) SetOracleClient(oracleClient oracleclient.OracleClient) {
	app.oracleClient = oracleClient
}

// Name returns the name of the App
func (app *MilkyWayApp) Name() string { return app.BaseApp.Name() }

// BeginBlocker application updates every begin block
func (app *MilkyWayApp) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	return app.mm.BeginBlock(ctx)
}

// EndBlocker application updates every end block
func (app *MilkyWayApp) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	return app.mm.EndBlock(ctx)
}

// InitChainer application update at chain initialization
func (app *MilkyWayApp) InitChainer(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	if err := app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap()); err != nil {
		panic(err)
	}

	response, err := app.mm.InitGenesis(ctx, app.appCodec, genesisState)
	if err != nil {
		panic(err)
	}

	return response, nil
}

// LoadHeight loads a particular height
func (app *MilkyWayApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range MaccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlockedModuleAccountAddrs returns all the app's blocked module account
// addresses.
func BlockedModuleAccountAddrs(modAccAddrs map[string]bool) map[string]bool {
	// remove module accounts that are ALLOWED to received funds
	delete(modAccAddrs, authtypes.NewModuleAddress(govtypes.ModuleName).String())
	delete(modAccAddrs, authtypes.NewModuleAddress(liquidvestingtypes.ModuleName).String())

	// Remove the ConsumerRewardsPool from the group of blocked recipient addresses in bank
	delete(modAccAddrs, authtypes.NewModuleAddress(providertypes.ConsumerRewardsPool).String())

	return modAccAddrs
}

// LegacyAmino returns MilkyWayApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *MilkyWayApp) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// AppCodec returns MilkyWay's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *MilkyWayApp) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns MilkyWay's InterfaceRegistry
func (app *MilkyWayApp) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// SimulationManager implements the SimulationApp interface
func (app *MilkyWayApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *MilkyWayApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	cmtservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	app.ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register nodeservice grpc-gateway routes.
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if err := server.RegisterSwaggerAPI(apiSvr.ClientCtx, apiSvr.Router, apiConfig.Swagger); err != nil {
		panic(err)
	}
}

// RegisterNodeService allows query minimum-gas-prices in app.toml
func (app *MilkyWayApp) RegisterNodeService(clientCtx client.Context, cfg config.Config) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter(), cfg)
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *MilkyWayApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *MilkyWayApp) RegisterTendermintService(clientCtx client.Context) {
	cmtservice.RegisterTendermintService(
		clientCtx,
		app.BaseApp.GRPCQueryRouter(),
		app.interfaceRegistry,
		app.Query,
	)
}

// configure store loader that checks if version == upgradeHeight and applies store upgrades
func (app *MilkyWayApp) setupUpgradeStoreLoaders() {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	for _, upgrade := range Upgrades {
		if upgradeInfo.Name == upgrade.UpgradeName {
			storeUpgrades := upgrade.StoreUpgrades
			app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
		}
	}
}

func (app *MilkyWayApp) setupUpgradeHandlers() {
	for _, upgrade := range Upgrades {
		app.UpgradeKeeper.SetUpgradeHandler(
			upgrade.UpgradeName,
			upgrade.CreateUpgradeHandler(
				app.mm,
				app.configurator,
				&app.AppKeepers,
			),
		)
	}
}

func (app *MilkyWayApp) setupForksPreBlockers() {
	for _, fork := range Forks {
		currentPreBlocker := app.PreBlocker()
		app.SetPreBlocker(func(ctx sdk.Context, req *abci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
			if req.Height == fork.ForkHeight {
				fork.BeginForkLogic(ctx, app.mm, app.configurator, &app.AppKeepers)
			}
			return currentPreBlocker(ctx, req)
		})
	}
}

// RegisterSwaggerAPI registers swagger route with API Server
func RegisterSwaggerAPI(rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

func (app *MilkyWayApp) OnTxSucceeded(_ sdk.Context, _, _ string, _ []byte, _ []byte) {
}

func (app *MilkyWayApp) OnTxFailed(_ sdk.Context, _, _ string, _ []byte, _ []byte) {
}

// AutoCliOpts returns the autocli options for the app.
func (app *MilkyWayApp) AutoCliOpts() autocli.AppOptions {
	modules := make(map[string]appmodule.AppModule, 0)
	for _, m := range app.mm.Modules {
		if moduleWithName, ok := m.(module.HasName); ok {
			moduleName := moduleWithName.Name()
			if appModule, ok := moduleWithName.(appmodule.AppModule); ok {
				modules[moduleName] = appModule
			}
		}
	}

	return autocli.AppOptions{
		Modules:               modules,
		AddressCodec:          authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		ValidatorAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		ConsensusAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	}
}

// TestingApp functions

// GetBaseApp implements the TestingApp interface.
func (app *MilkyWayApp) GetBaseApp() *baseapp.BaseApp {
	return app.BaseApp
}

// GetTxConfig implements the TestingApp interface.
func (app *MilkyWayApp) GetTxConfig() client.TxConfig {
	return app.txConfig
}

// GetTestGovKeeper implements the TestingApp interface.
func (app *MilkyWayApp) GetTestGovKeeper() *govkeeper.Keeper {
	return app.AppKeepers.GovKeeper
}

// EmptyAppOptions is a stub implementing AppOptions
type EmptyAppOptions struct{}

// EmptyWasmOptions is a stub implementing Wasmkeeper Option
var EmptyWasmOptions []wasmkeeper.Option

// Get implements AppOptions
func (ao EmptyAppOptions) Get(_ string) interface{} {
	return nil
}

// Close closes the underlying baseapp, the oracle service, and the prometheus server if required.
// This method blocks on the closure of both the prometheus server, and the oracle-service
func (app *MilkyWayApp) Close() error {
	if err := app.BaseApp.Close(); err != nil {
		return err
	}

	// Close the oracle service
	if app.oracleClient != nil {
		if err := app.oracleClient.Stop(); err != nil {
			return err
		}
	}

	return nil
}

// StartOracleClient starts the oracle client
func (app *MilkyWayApp) StartOracleClient(ctx context.Context) error {
	if app.oracleClient != nil {
		return app.oracleClient.Start(ctx)
	}

	return nil
}
