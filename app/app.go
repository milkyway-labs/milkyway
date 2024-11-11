package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gorilla/mux"
	"github.com/initia-labs/OPinit/x/opchild"
	opchildtypes "github.com/initia-labs/OPinit/x/opchild/types"
	ibchooks "github.com/initia-labs/initia/x/ibc-hooks"
	ibchookstypes "github.com/initia-labs/initia/x/ibc-hooks/types"
	"github.com/rakyll/statik/fs"
	"github.com/skip-mev/block-sdk/v2/block"
	"github.com/skip-mev/block-sdk/v2/x/auction"
	auctionante "github.com/skip-mev/block-sdk/v2/x/auction/ante"
	auctionkeeper "github.com/skip-mev/block-sdk/v2/x/auction/keeper"
	auctiontypes "github.com/skip-mev/block-sdk/v2/x/auction/types"
	"github.com/skip-mev/connect/v2/x/marketmap"
	marketmaptypes "github.com/skip-mev/connect/v2/x/marketmap/types"
	"github.com/skip-mev/connect/v2/x/oracle"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"
	"github.com/spf13/cast"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"
	"cosmossdk.io/core/address"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/feegrant"
	feegrantmodule "cosmossdk.io/x/feegrant/module"
	"cosmossdk.io/x/upgrade"
	upgradetypes "cosmossdk.io/x/upgrade/types"
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
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	cosmosante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/posthandler"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	groupmodule "github.com/cosmos/cosmos-sdk/x/group/module"
	"github.com/cosmos/gogoproto/proto"

	// ibc imports
	"github.com/Stride-Labs/ibc-rate-limiting/ratelimit"
	ratelimittypes "github.com/Stride-Labs/ibc-rate-limiting/ratelimit/types"
	"github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward"
	packetforwardkeeper "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/keeper"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"
	"github.com/cosmos/ibc-go/modules/capability"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ica "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts"
	icacontroller "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
	icahost "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	ibcfee "github.com/cosmos/ibc-go/v8/modules/apps/29-fee"
	ibcfeetypes "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/types"
	ibctransfer "github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v8/modules/core"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	solomachine "github.com/cosmos/ibc-go/v8/modules/light-clients/06-solomachine"
	ibctm "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"

	// initia imports

	initialanes "github.com/initia-labs/initia/app/lanes"
	"github.com/initia-labs/initia/app/params"
	ibctestingtypes "github.com/initia-labs/initia/x/ibc/testing/types"
	icaauth "github.com/initia-labs/initia/x/intertx"
	icaauthkeeper "github.com/initia-labs/initia/x/intertx/keeper"
	icaauthtypes "github.com/initia-labs/initia/x/intertx/types"

	opchildlanes "github.com/initia-labs/OPinit/x/opchild/lanes"
	// skip imports
	mevabci "github.com/skip-mev/block-sdk/v2/abci"
	blockchecktx "github.com/skip-mev/block-sdk/v2/abci/checktx"
	signer_extraction "github.com/skip-mev/block-sdk/v2/adapters/signer_extraction_adapter"
	blockbase "github.com/skip-mev/block-sdk/v2/block/base"
	mevlane "github.com/skip-mev/block-sdk/v2/lanes/mev"
	// CosmWasm imports
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	// local imports
	appante "github.com/milkyway-labs/milkyway/app/ante"
	ibcwasmhooks "github.com/milkyway-labs/milkyway/app/ibc-hooks"
	"github.com/milkyway-labs/milkyway/app/upgrades"
	"github.com/milkyway-labs/milkyway/utils"
	"github.com/milkyway-labs/milkyway/x/assets"
	assetstypes "github.com/milkyway-labs/milkyway/x/assets/types"
	"github.com/milkyway-labs/milkyway/x/bank"
	"github.com/milkyway-labs/milkyway/x/epochs"
	epochstypes "github.com/milkyway-labs/milkyway/x/epochs/types"
	"github.com/milkyway-labs/milkyway/x/icacallbacks"
	icacallbackstypes "github.com/milkyway-labs/milkyway/x/icacallbacks/types"
	"github.com/milkyway-labs/milkyway/x/interchainquery"
	icqtypes "github.com/milkyway-labs/milkyway/x/interchainquery/types"
	"github.com/milkyway-labs/milkyway/x/liquidvesting"
	liquidvestinghooks "github.com/milkyway-labs/milkyway/x/liquidvesting/hooks"
	liquidvestingtypes "github.com/milkyway-labs/milkyway/x/liquidvesting/types"
	"github.com/milkyway-labs/milkyway/x/operators"
	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	"github.com/milkyway-labs/milkyway/x/pools"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	"github.com/milkyway-labs/milkyway/x/records"
	recordskeeper "github.com/milkyway-labs/milkyway/x/records/keeper"
	recordstypes "github.com/milkyway-labs/milkyway/x/records/types"
	"github.com/milkyway-labs/milkyway/x/restaking"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	"github.com/milkyway-labs/milkyway/x/rewards"
	rewardstypes "github.com/milkyway-labs/milkyway/x/rewards/types"
	"github.com/milkyway-labs/milkyway/x/services"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
	"github.com/milkyway-labs/milkyway/x/stakeibc"
	stakeibckeeper "github.com/milkyway-labs/milkyway/x/stakeibc/keeper"
	stakeibctypes "github.com/milkyway-labs/milkyway/x/stakeibc/types"
	"github.com/milkyway-labs/milkyway/x/tokenfactory"
	tokenfactorykeeper "github.com/milkyway-labs/milkyway/x/tokenfactory/keeper"
	tokenfactorytypes "github.com/milkyway-labs/milkyway/x/tokenfactory/types"

	// noble forwarding keeper
	"github.com/noble-assets/forwarding/v2/x/forwarding"
	forwardingtypes "github.com/noble-assets/forwarding/v2/x/forwarding/types"

	// kvindexer
	indexer "github.com/initia-labs/kvindexer"
	indexerconfig "github.com/initia-labs/kvindexer/config"
	blocksubmodule "github.com/initia-labs/kvindexer/submodules/block"
	"github.com/initia-labs/kvindexer/submodules/tx"
	nft "github.com/initia-labs/kvindexer/submodules/wasm-nft"
	pair "github.com/initia-labs/kvindexer/submodules/wasm-pair"
	indexermodule "github.com/initia-labs/kvindexer/x/kvindexer"
	indexerkeeper "github.com/initia-labs/kvindexer/x/kvindexer/keeper"

	// unnamed import of statik for swagger UI support
	_ "github.com/milkyway-labs/milkyway/client/docs/statik"
)

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string

	ModuleBasics = module.NewBasicManager(
		auth.AppModule{},
		bank.AppModule{},
		crisis.AppModule{},
		opchild.AppModule{},
		capability.AppModule{},
		feegrantmodule.AppModule{},
		upgrade.AppModule{},
		authzmodule.AppModule{},
		groupmodule.AppModule{},
		consensus.AppModule{},
		wasm.AppModule{},
		auction.AppModule{},
		tokenfactory.AppModule{},

		// ibc modules
		ibc.AppModule{},
		ibctransfer.AppModule{},
		ica.AppModule{},
		icaauth.AppModule{},
		ibcfee.AppModule{},
		ibctm.AppModule{},
		solomachine.AppModule{},
		packetforward.AppModule{},
		ibchooks.AppModule{},
		forwarding.AppModule{},
		ratelimit.AppModule{},

		// connect modules
		oracle.AppModule{},
		marketmap.AppModule{},

		// liquid staking modules
		stakeibc.AppModule{},
		epochs.AppModule{},
		interchainquery.AppModule{},
		records.AppModule{},
		icacallbacks.AppModule{},

		// custom modules
		services.AppModule{},
		operators.AppModule{},
		pools.AppModule{},
		restaking.AppModule{},
		assets.AppModule{},
		rewards.AppModule{},
		liquidvesting.AppModule{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:  nil,
		icatypes.ModuleName:         nil,
		ibcfeetypes.ModuleName:      nil,
		ibctransfertypes.ModuleName: {authtypes.Minter, authtypes.Burner},
		// x/auction's module account must be instantiated upon genesis to accrue auction rewards not
		// distributed to proposers
		auctiontypes.ModuleName:           nil,
		opchildtypes.ModuleName:           {authtypes.Minter, authtypes.Burner},
		tokenfactorytypes.ModuleName:      {authtypes.Minter, authtypes.Burner},
		icqtypes.ModuleName:               nil,
		stakeibctypes.ModuleName:          {authtypes.Minter, authtypes.Burner, authtypes.Staking},
		stakeibctypes.RewardCollectorName: nil,
		rewardstypes.RewardsPoolName:      nil,

		// connect oracle permissions
		oracletypes.ModuleName: nil,

		// MilkyWay permissions
		liquidvestingtypes.ModuleName: {authtypes.Minter, authtypes.Burner},

		// this is only for testing
		authtypes.Minter: {authtypes.Minter},
	}
)

var _ servertypes.Application = (*MilkyWayApp)(nil)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, "."+AppName)
}

// MilkyWayApp extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type MilkyWayApp struct {
	*baseapp.BaseApp

	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry types.InterfaceRegistry

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// keepers
	// TODO: add gov keeper

	// the module manager
	ModuleManager      *module.Manager
	BasicModuleManager module.BasicManager

	// the configurator
	configurator module.Configurator

	// Override of BaseApp's CheckTx
	checkTxHandler blockchecktx.CheckTx

	// kvindexer
	indexerKeeper *indexerkeeper.Keeper
	indexerModule indexermodule.AppModuleBasic
}

// NewMilkyWayApp returns a reference to an initialized Initia.
func NewMilkyWayApp(
	logger log.Logger,
	db dbm.DB,
	kvindexerDB dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	wasmOpts []wasmkeeper.Option,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *MilkyWayApp {
	// load the configs
	mempoolTxs := cast.ToInt(appOpts.Get(server.FlagMempoolMaxTxs))
	queryGasLimit := cast.ToInt(appOpts.Get(server.FlagQueryGasLimit))

	logger.Info("mempool max txs", "max_txs", mempoolTxs)
	logger.Info("query gas limit", "gas_limit", queryGasLimit)

	encodingConfig := params.MakeEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	appCodec := encodingConfig.Codec
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry
	txConfig := encodingConfig.TxConfig

	bApp := baseapp.NewBaseApp(AppName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)
	bApp.SetTxEncoder(txConfig.TxEncoder())

	keys := storetypes.NewKVStoreKeys()
	tkeys := storetypes.NewTransientStoreKeys(forwardingtypes.TransientStoreKey)
	memKeys := storetypes.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	// register streaming services
	if err := bApp.RegisterStreamingServices(appOpts, keys); err != nil {
		panic(err)
	}

	app := &MilkyWayApp{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		txConfig:          txConfig,
		interfaceRegistry: interfaceRegistry,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	ac := authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	vc := authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix())
	cc := authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix())

	authorityAccAddr := authtypes.NewModuleAddress(opchildtypes.ModuleName)
	authorityAddr, err := ac.BytesToString(authorityAccAddr)
	if err != nil {
		panic(err)
	}

	hooksICS4Wrapper := ibchooks.NewICS4Middleware(
		app.RateLimitKeeper,
		ibcwasmhooks.NewWasmHooks(appCodec, ac, app.WasmKeeper),
	)

	hooksICS4LiquidVesting := ibchooks.NewICS4Middleware(
		hooksICS4Wrapper,
		liquidvestinghooks.NewIBCHooks(app.LiquidVestingKeeper),
	)

	////////////////////////////
	// Transfer configuration //
	////////////////////////////
	// Send   : transfer -> packet forward -> wasm   -> fee            -> channel
	// Receive: channel  -> fee            -> wasm   -> packet forward -> forwarding -> transfer

	var transferStack porttypes.IBCModule
	{
		packetForwardKeeper := &packetforwardkeeper.Keeper{}

		// Create Transfer Keepers
		transferKeeper := ibctransferkeeper.NewKeeper(
			appCodec,
			keys[ibctransfertypes.StoreKey],
			nil, // we don't need migration
			// ics4wrapper: transfer -> packet forward
			packetForwardKeeper,
			app.IBCKeeper.ChannelKeeper,
			app.IBCKeeper.PortKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.ScopedTransferKeeper,
			authorityAddr,
		)
		app.TransferKeeper = &transferKeeper
		transferStack = ibctransfer.NewIBCModule(*app.TransferKeeper)

		// forwarding middleware
		transferStack = forwarding.NewMiddleware(
			// receive: forwarding -> transfer
			transferStack,
			app.AccountKeeper,
			app.ForwardingKeeper,
		)

		// create packet forward middleware
		*packetForwardKeeper = *packetforwardkeeper.NewKeeper(
			appCodec,
			keys[packetforwardtypes.StoreKey],
			app.TransferKeeper,
			app.IBCKeeper.ChannelKeeper,
			communityPoolKeeper,
			app.BankKeeper,
			// ics4wrapper: transfer -> packet forward -> fee
			hooksICS4Wrapper, // TODO: is it correct?
			authorityAddr,
		)
		app.PacketForwardKeeper = packetForwardKeeper
		transferStack = packetforward.NewIBCMiddleware(
			// receive: packet forward -> forwarding -> transfer
			transferStack,
			app.PacketForwardKeeper,
			0,
			packetforwardkeeper.DefaultForwardTransferPacketTimeoutTimestamp,
			packetforwardkeeper.DefaultRefundTransferPacketTimeoutTimestamp,
		)

		// create wasm middleware for transfer
		transferStack = ibchooks.NewIBCMiddleware(
			// receive: wasm -> packet forward -> forwarding -> transfer
			transferStack,
			hooksICS4Wrapper,
			app.IBCHooksKeeper,
		)

		transferStack = ibchooks.NewIBCMiddleware(
			// receive: liquidvesting -> wasm -> packet forward -> forwarding -> transfer
			transferStack,
			hooksICS4LiquidVesting,
			app.IBCHooksKeeper,
		)

		transferStack = ratelimit.NewIBCMiddleware(app.RateLimitKeeper, transferStack)

		app.RecordsKeeper = *recordskeeper.NewKeeper(
			appCodec,
			keys[recordstypes.StoreKey],
			keys[recordstypes.MemStoreKey],
			app.AccountKeeper,
			*app.TransferKeeper,
			*app.IBCKeeper,
			app.ICACallbacksKeeper,
		)
		transferStack = records.NewIBCModule(app.RecordsKeeper, transferStack)

		// create ibcfee middleware for transfer
		transferStack = ibcfee.NewIBCMiddleware(
			// receive: fee -> wasm -> packet forward -> forwarding -> transfer
			transferStack,
			// ics4wrapper: transfer -> packet forward -> wasm -> fee -> channel
			*app.IBCFeeKeeper,
		)
	}

	///////////////////////
	// ICA configuration //
	///////////////////////

	var icaHostStack porttypes.IBCModule
	var icaControllerStack porttypes.IBCModule
	var icaCallbacksStack porttypes.IBCModule
	{
		icaHostKeeper := icahostkeeper.NewKeeper(
			appCodec, keys[icahosttypes.StoreKey],
			nil, // we don't need migration
			app.IBCFeeKeeper,
			app.IBCKeeper.ChannelKeeper,
			app.IBCKeeper.PortKeeper,
			app.AccountKeeper,
			app.ScopedICAHostKeeper,
			app.MsgServiceRouter(),
			authorityAddr,
		)
		icaHostKeeper.WithQueryRouter(bApp.GRPCQueryRouter())
		app.ICAHostKeeper = &icaHostKeeper

		icaControllerKeeper := icacontrollerkeeper.NewKeeper(
			appCodec, keys[icacontrollertypes.StoreKey],
			nil, // we don't need migration
			app.IBCFeeKeeper,
			app.IBCKeeper.ChannelKeeper,
			app.IBCKeeper.PortKeeper,
			app.ScopedICAControllerKeeper,
			app.MsgServiceRouter(),
			authorityAddr,
		)
		app.ICAControllerKeeper = &icaControllerKeeper

		icaAuthKeeper := icaauthkeeper.NewKeeper(
			appCodec,
			*app.ICAControllerKeeper,
			app.ScopedICAAuthKeeper,
			ac,
		)
		app.ICAAuthKeeper = &icaAuthKeeper

		icaHostIBCModule := icahost.NewIBCModule(*app.ICAHostKeeper)
		icaHostStack = ibcfee.NewIBCMiddleware(icaHostIBCModule, *app.IBCFeeKeeper)

		icaAuthIBCModule := icaauth.NewIBCModule(*app.ICAAuthKeeper)
		icaControllerIBCModule := icacontroller.NewIBCMiddleware(icaAuthIBCModule, *app.ICAControllerKeeper)
		icaControllerStack = ibcfee.NewIBCMiddleware(icaControllerIBCModule, *app.IBCFeeKeeper)

		icaCallbacksStack = icacallbacks.NewIBCModule(app.ICACallbacksKeeper)

		app.StakeIBCKeeper = stakeibckeeper.NewKeeper(
			appCodec,
			keys[stakeibctypes.StoreKey],
			keys[stakeibctypes.MemStoreKey],
			runtime.NewKVStoreService(keys[stakeibctypes.StoreKey]),
			authorityAddr,
			app.AccountKeeper,
			app.BankKeeper,
			*app.ICAControllerKeeper,
			*app.IBCKeeper,
			app.InterchainQueryKeeper,
			app.RecordsKeeper,
			app.ICACallbacksKeeper,
			app.RateLimitKeeper,
			app.OPChildKeeper,
		)
		app.StakeIBCKeeper.SetHooks(stakeibctypes.NewMultiStakeIBCHooks())
		icaCallbacksStack = stakeibc.NewIBCMiddleware(icaCallbacksStack, app.StakeIBCKeeper)
		icaCallbacksStack = icacontroller.NewIBCMiddleware(icaCallbacksStack, *app.ICAControllerKeeper)
		icaCallbacksStack = ibcfee.NewIBCMiddleware(icaCallbacksStack, *app.IBCFeeKeeper)
	}

	//////////////////////////////
	// Wasm IBC Configuration   //
	//////////////////////////////

	var wasmIBCStack porttypes.IBCModule
	{
		wasmIBCModule := wasm.NewIBCHandler(
			app.WasmKeeper,
			app.IBCKeeper.ChannelKeeper,
			// ics4wrapper: wasm -> fee
			app.IBCFeeKeeper,
		)

		// create wasm middleware for wasm IBC stack
		hookMiddleware := ibchooks.NewIBCMiddleware(
			// receive: hook -> wasm
			wasmIBCModule,
			hooksICS4Wrapper, // TODO: is it correct?
			app.IBCHooksKeeper,
		)

		wasmIBCStack = ibcfee.NewIBCMiddleware(
			// receive: fee -> hook -> wasm
			hookMiddleware,
			*app.IBCFeeKeeper,
		)
	}

	//////////////////////////////
	// IBC router Configuration //
	//////////////////////////////

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferStack).
		// TODO: add ica callbacks stack
		AddRoute(icahosttypes.SubModuleName, icaHostStack).
		AddRoute(icacontrollertypes.SubModuleName, icaCallbacksStack).
		AddRoute(icaauthtypes.ModuleName, icaControllerStack).
		AddRoute(wasmtypes.ModuleName, wasmIBCStack)

	app.IBCKeeper.SetRouter(ibcRouter)

	//////////////////////////////
	// WasmKeeper Configuration //
	//////////////////////////////
	wasmDir := filepath.Join(homePath, "wasm")
	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic(fmt.Sprintf("error while reading wasm config: %s", err))
	}

	// allow connect queries
	queryAllowlist := make(map[string]proto.Message)
	queryAllowlist["/connect.oracle.v2.Query/GetAllCurrencyPairs"] = &oracletypes.GetAllCurrencyPairsResponse{}
	queryAllowlist["/connect.oracle.v2.Query/GetPrice"] = &oracletypes.GetPriceResponse{}
	queryAllowlist["/connect.oracle.v2.Query/GetPrices"] = &oracletypes.GetPricesResponse{}
	queryAllowlist["/milkyway.operators.v1.Query/Operator"] = &operatorstypes.QueryOperatorResponse{}
	queryAllowlist["/milkyway.restaking.v1.Query/ServiceOperators"] = &restakingtypes.QueryServiceOperatorsResponse{}

	// use accept list stargate querier
	wasmOpts = append(wasmOpts, wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Stargate: wasmkeeper.AcceptListStargateQuerier(queryAllowlist, app.GRPCQueryRouter(), appCodec),
	}))

	// The last arguments can contain custom message handlers, and custom query handlers,
	// if we want to allow any custom callbacks
	*app.WasmKeeper = wasmkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[wasmtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		// we do not support staking feature, so don't need to provide these keepers
		nil,
		nil,
		app.IBCFeeKeeper, // ISC4 Wrapper: fee IBC middleware
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.PortKeeper,
		app.ScopedWasmKeeper,
		app.TransferKeeper,
		app.MsgServiceRouter(),
		app.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		slices.DeleteFunc(wasmkeeper.BuiltInCapabilities(), func(s string) bool {
			return s == "staking"
		}),
		authorityAddr,
		wasmOpts...,
	)

	// x/auction module keeper initialization

	// initialize the keeper
	auctionKeeper := auctionkeeper.NewKeeperWithRewardsAddressProvider(
		app.appCodec,
		app.keys[auctiontypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		opchildlanes.NewRewardsAddressProvider(authtypes.FeeCollectorName),
		authorityAddr,
	)
	app.AuctionKeeper = &auctionKeeper

	contractKeeper := wasmkeeper.NewDefaultPermissionKeeper(app.WasmKeeper)

	tokenfactoryKeeper := tokenfactorykeeper.NewKeeper(
		ac,
		appCodec,
		runtime.NewKVStoreService(keys[tokenfactorytypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		communityPoolKeeper,
		authorityAddr,
	)
	app.TokenFactoryKeeper = &tokenfactoryKeeper
	app.TokenFactoryKeeper.SetContractKeeper(contractKeeper)

	app.BankKeeper.SetHooks(app.TokenFactoryKeeper.Hooks())

	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.

	app.ModuleManager = module.NewManager(
		auth.NewAppModule(appCodec, *app.AccountKeeper, nil, nil),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, nil),
		crisis.NewAppModule(app.CrisisKeeper, skipGenesisInvariants, nil),
		opchild.NewAppModule(appCodec, app.OPChildKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, *app.FeeGrantKeeper, app.interfaceRegistry),
		upgrade.NewAppModule(app.UpgradeKeeper, ac),
		authzmodule.NewAppModule(appCodec, *app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		groupmodule.NewAppModule(appCodec, *app.GroupKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		consensus.NewAppModule(appCodec, *app.ConsensusParamsKeeper),
		wasm.NewAppModule(appCodec, app.WasmKeeper, nil /* unused */, app.AccountKeeper, app.BankKeeper, app.MsgServiceRouter(), nil),
		auction.NewAppModule(app.appCodec, *app.AuctionKeeper),
		tokenfactory.NewAppModule(appCodec, app.TokenFactoryKeeper, *app.AccountKeeper, *app.BankKeeper),
		// ibc modules
		ibc.NewAppModule(app.IBCKeeper),
		ibctransfer.NewAppModule(*app.TransferKeeper),
		ica.NewAppModule(app.ICAControllerKeeper, app.ICAHostKeeper),
		icaauth.NewAppModule(appCodec, *app.ICAAuthKeeper),
		ibcfee.NewAppModule(*app.IBCFeeKeeper),
		ibctm.NewAppModule(),
		solomachine.NewAppModule(),
		packetforward.NewAppModule(app.PacketForwardKeeper, nil),
		ibchooks.NewAppModule(appCodec, *app.IBCHooksKeeper),
		forwarding.NewAppModule(app.ForwardingKeeper),
		ratelimit.NewAppModule(appCodec, app.RateLimitKeeper),
		// connect modules
		oracle.NewAppModule(appCodec, *app.OracleKeeper),
		marketmap.NewAppModule(appCodec, app.MarketMapKeeper),
		// liquid staking modules
		stakeibc.NewAppModule(appCodec, app.StakeIBCKeeper, app.AccountKeeper, app.BankKeeper),
		epochs.NewAppModule(appCodec, app.EpochsKeeper),
		interchainquery.NewAppModule(appCodec, app.InterchainQueryKeeper),
		records.NewAppModule(appCodec, app.RecordsKeeper, app.AccountKeeper, app.BankKeeper),
		icacallbacks.NewAppModule(appCodec, app.ICACallbacksKeeper, app.AccountKeeper, app.BankKeeper),
		// custom modules
		services.NewAppModule(appCodec, app.ServicesKeeper, app.PoolsKeeper),
		operators.NewAppModule(appCodec, app.OperatorsKeeper),
		pools.NewAppModule(appCodec, app.PoolsKeeper),
		restaking.NewAppModule(appCodec, app.RestakingKeeper),
		assets.NewAppModule(appCodec, app.AssetsKeeper),
		rewards.NewAppModule(appCodec, app.RewardsKeeper),
		liquidvesting.NewAppModule(appCodec, app.LiquidVestingKeeper),
	)

	if err := app.setupIndexer(kvindexerDB, appOpts, ac, vc, appCodec); err != nil {
		panic(err)
	}

	// BasicModuleManager defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration and genesis verification.
	// By default it is composed of all the module from the module manager.
	// Additionally, app module basics can be overwritten by passing them as argument.
	app.BasicModuleManager = module.NewBasicManagerFromManager(
		app.ModuleManager,
		map[string]module.AppModuleBasic{
			genutiltypes.ModuleName: genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
		})
	app.BasicModuleManager.RegisterLegacyAminoCodec(legacyAmino)
	app.BasicModuleManager.RegisterInterfaces(interfaceRegistry)

	// NOTE: upgrade module is required to be prioritized
	app.ModuleManager.SetOrderPreBlockers(
		upgradetypes.ModuleName,
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.ModuleManager.SetOrderBeginBlockers(
		capabilitytypes.ModuleName,
		opchildtypes.ModuleName,
		authz.ModuleName,
		ibcexported.ModuleName,
		oracletypes.ModuleName,
		marketmaptypes.ModuleName,
		stakeibctypes.ModuleName,
		epochstypes.ModuleName,
		ratelimittypes.ModuleName,

		rewardstypes.ModuleName,
		servicestypes.ModuleName,
		operatorstypes.ModuleName,
		poolstypes.ModuleName,
		restakingtypes.ModuleName,
	)

	app.ModuleManager.SetOrderEndBlockers(
		crisistypes.ModuleName,
		opchildtypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		group.ModuleName,
		oracletypes.ModuleName,
		marketmaptypes.ModuleName,
		forwardingtypes.ModuleName,
		stakeibctypes.ModuleName,
		icqtypes.ModuleName,
		ratelimittypes.ModuleName,

		servicestypes.ModuleName,
		operatorstypes.ModuleName,
		poolstypes.ModuleName,
		restakingtypes.ModuleName,
		liquidvestingtypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	genesisModuleOrder := []string{
		capabilitytypes.ModuleName, authtypes.ModuleName, banktypes.ModuleName,
		opchildtypes.ModuleName, genutiltypes.ModuleName, authz.ModuleName, group.ModuleName,
		upgradetypes.ModuleName, feegrant.ModuleName, consensusparamtypes.ModuleName,
		ibcexported.ModuleName, ibctransfertypes.ModuleName, icatypes.ModuleName,
		icaauthtypes.ModuleName, ibcfeetypes.ModuleName, auctiontypes.ModuleName,
		wasmtypes.ModuleName, oracletypes.ModuleName, marketmaptypes.ModuleName,
		packetforwardtypes.ModuleName, tokenfactorytypes.ModuleName,
		ibchookstypes.ModuleName, forwardingtypes.ModuleName,
		stakeibctypes.ModuleName, epochstypes.ModuleName, icqtypes.ModuleName,
		recordstypes.ModuleName, ratelimittypes.ModuleName, icacallbackstypes.ModuleName,

		servicestypes.ModuleName, operatorstypes.ModuleName, poolstypes.ModuleName, restakingtypes.ModuleName,
		assetstypes.ModuleName, rewardstypes.ModuleName, liquidvestingtypes.ModuleName,
		crisistypes.ModuleName,
	}

	app.ModuleManager.SetOrderInitGenesis(genesisModuleOrder...)
	app.ModuleManager.SetOrderExportGenesis(genesisModuleOrder...)

	app.ModuleManager.RegisterInvariants(app.CrisisKeeper)

	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	err = app.ModuleManager.RegisterServices(app.configurator)
	if err != nil {
		panic(err)
	}

	app.indexerModule.RegisterServices(app.configurator)

	// register upgrade handler for later use
	app.RegisterUpgradeHandlers()

	autocliv1.RegisterQueryServer(app.GRPCQueryRouter(), runtimeservices.NewAutoCLIQueryService(app.ModuleManager.Modules))

	reflectionSvc, err := runtimeservices.NewReflectionService()
	if err != nil {
		panic(err)
	}
	reflectionv1.RegisterReflectionServiceServer(app.GRPCQueryRouter(), reflectionSvc)

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetPreBlocker(app.PreBlocker)
	app.SetBeginBlocker(app.BeginBlocker)
	app.setPostHandler()
	app.SetEndBlocker(app.EndBlocker)

	// initialize and set the InitiaApp mempool. The current mempool will be the
	// x/auction module's mempool which will extract the top bid from the current block's auction
	// and insert the txs at the top of the block spots.
	signerExtractor := signer_extraction.NewDefaultAdapter()

	systemLane := initialanes.NewSystemLane(blockbase.LaneConfig{
		Logger:          app.Logger(),
		TxEncoder:       app.txConfig.TxEncoder(),
		TxDecoder:       app.txConfig.TxDecoder(),
		MaxBlockSpace:   math.LegacyMustNewDecFromStr("0.01"),
		MaxTxs:          1,
		SignerExtractor: signerExtractor,
	}, opchildlanes.SystemLaneMatchHandler())

	factory := mevlane.NewDefaultAuctionFactory(app.txConfig.TxDecoder(), signerExtractor)
	mevLane := mevlane.NewMEVLane(blockbase.LaneConfig{
		Logger:          app.Logger(),
		TxEncoder:       app.txConfig.TxEncoder(),
		TxDecoder:       app.txConfig.TxDecoder(),
		MaxBlockSpace:   math.LegacyMustNewDecFromStr("0.09"),
		MaxTxs:          100,
		SignerExtractor: signerExtractor,
	}, factory, factory.MatchHandler())

	freeLane := initialanes.NewFreeLane(blockbase.LaneConfig{
		Logger:          app.Logger(),
		TxEncoder:       app.txConfig.TxEncoder(),
		TxDecoder:       app.txConfig.TxDecoder(),
		MaxBlockSpace:   math.LegacyMustNewDecFromStr("0.1"),
		MaxTxs:          100,
		SignerExtractor: signerExtractor,
	}, opchildlanes.NewFreeLaneMatchHandler(ac, app.OPChildKeeper).MatchHandler())

	defaultLane := initialanes.NewDefaultLane(blockbase.LaneConfig{
		Logger:          app.Logger(),
		TxEncoder:       app.txConfig.TxEncoder(),
		TxDecoder:       app.txConfig.TxDecoder(),
		MaxBlockSpace:   math.LegacyMustNewDecFromStr("0.8"),
		MaxTxs:          mempoolTxs,
		SignerExtractor: signerExtractor,
	})

	lanes := []block.Lane{systemLane, mevLane, freeLane, defaultLane}
	mempool, err := block.NewLanedMempool(app.Logger(), lanes)
	if err != nil {
		panic(err)
	}

	app.SetMempool(mempool)
	anteHandler := app.setAnteHandler(mevLane, freeLane, wasmConfig, keys[wasmtypes.StoreKey])

	// set the ante handler for each lane
	//
	opt := []blockbase.LaneOption{
		blockbase.WithAnteHandler(anteHandler),
	}
	systemLane.(*blockbase.BaseLane).WithOptions(
		opt...,
	)
	mevLane.WithOptions(
		opt...,
	)
	freeLane.(*blockbase.BaseLane).WithOptions(
		opt...,
	)
	defaultLane.(*blockbase.BaseLane).WithOptions(
		opt...,
	)

	// override the base-app's ABCI methods (CheckTx, PrepareProposal, ProcessProposal)
	proposalHandlers := mevabci.NewProposalHandler(
		app.Logger(),
		app.txConfig.TxDecoder(),
		app.txConfig.TxEncoder(),
		mempool,
	)

	// override base-app's ProcessProposal + PrepareProposal
	app.SetPrepareProposal(proposalHandlers.PrepareProposalHandler())
	app.SetProcessProposal(proposalHandlers.ProcessProposalHandler())

	// overrde base-app's CheckTx
	mevCheckTx := blockchecktx.NewMEVCheckTxHandler(
		app.BaseApp,
		app.txConfig.TxDecoder(),
		mevLane,
		anteHandler,
		app.BaseApp.CheckTx,
	)
	checkTxHandler := blockchecktx.NewMempoolParityCheckTx(
		app.Logger(), mempool,
		app.txConfig.TxDecoder(), mevCheckTx.CheckTx(),
	)
	app.SetCheckTx(checkTxHandler.CheckTx())

	////////////////
	/// lane end ///
	////////////////

	// At startup, after all modules have been registered, check that all proto
	// annotations are correct.
	protoFiles, err := proto.MergedRegistry()
	if err != nil {
		panic(err)
	}
	err = msgservice.ValidateProtoAnnotations(protoFiles)
	if err != nil {
		errMsg := ""

		// ignore injective and ibc-rate-limiting proto annotations comes from github.com/cosoms/relayer
		for _, s := range strings.Split(err.Error(), "\n") {
			if strings.Contains(s, "injective") || strings.Contains(s, "ratelimit") {
				continue
			}

			errMsg += s + "\n"
		}

		if errMsg != "" {
			// Once we switch to using protoreflect-based antehandlers, we might
			// want to panic here instead of logging a warning.
			fmt.Fprintln(os.Stderr, errMsg)
		}
	}

	// must be before Loading version
	// requires the snapshot store to be created and registered as a BaseAppOption
	// see cmd/wasmd/root.go: 206 - 214 approx
	if manager := app.SnapshotManager(); manager != nil {
		err := manager.RegisterExtensions(
			wasmkeeper.NewWasmSnapshotter(app.CommitMultiStore(), app.WasmKeeper),
		)
		if err != nil {
			panic(fmt.Errorf("failed to register snapshot extension: %s", err))
		}
	}

	// Load the latest state from disk if necessary, and initialize the base-app. From this point on
	// no more modifications to the base-app can be made
	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}

		ctx := app.BaseApp.NewUncachedContext(true, tmproto.Header{})

		// Initialize pinned codes in wasmvm as they are not persisted there
		if err := app.WasmKeeper.InitializePinnedCodes(ctx); err != nil {
			tmos.Exit(fmt.Sprintf("failed initialize pinned codes %s", err))
		}
	}

	return app
}

// CheckTx will check the transaction with the provided checkTxHandler. We override the default
// handler so that we can verify bid transactions before they are inserted into the mempool.
// With the POB CheckTx, we can verify the bid transaction and all of the bundled transactions
// before inserting the bid transaction into the mempool.
func (app *MilkyWayApp) CheckTx(req *abci.RequestCheckTx) (*abci.ResponseCheckTx, error) {
	return app.checkTxHandler(req)
}

// SetCheckTx sets the checkTxHandler for the app.
func (app *MilkyWayApp) SetCheckTx(handler blockchecktx.CheckTx) {
	app.checkTxHandler = handler
}

func (app *MilkyWayApp) setAnteHandler(
	mevLane auctionante.MEVLane,
	freeLane block.Lane,
	wasmConfig wasmtypes.WasmConfig,
	txCounterStoreKey *storetypes.KVStoreKey,
) sdk.AnteHandler {
	anteHandler, err := appante.NewAnteHandler(
		appante.HandlerOptions{
			HandlerOptions: cosmosante.HandlerOptions{
				AccountKeeper:   app.AccountKeeper,
				BankKeeper:      app.BankKeeper,
				FeegrantKeeper:  app.FeeGrantKeeper,
				SignModeHandler: app.txConfig.SignModeHandler(),
			},
			IBCkeeper:             app.IBCKeeper,
			Codec:                 app.appCodec,
			OPChildKeeper:         app.OPChildKeeper,
			TxEncoder:             app.txConfig.TxEncoder(),
			AuctionKeeper:         *app.AuctionKeeper,
			MevLane:               mevLane,
			FreeLane:              freeLane,
			WasmKeeper:            app.WasmKeeper,
			WasmConfig:            &wasmConfig,
			TXCounterStoreService: runtime.NewKVStoreService(txCounterStoreKey),
		},
	)
	if err != nil {
		panic(err)
	}

	app.SetAnteHandler(anteHandler)
	return anteHandler
}

func (app *MilkyWayApp) setPostHandler() {
	postHandler, err := posthandler.NewPostHandler(
		posthandler.HandlerOptions{},
	)
	if err != nil {
		panic(err)
	}

	app.SetPostHandler(postHandler)
}

// Name returns the name of the App
func (app *MilkyWayApp) Name() string { return app.BaseApp.Name() }

// PreBlocker application updates every pre block
func (app *MilkyWayApp) PreBlocker(ctx sdk.Context, _ *abci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
	return app.ModuleManager.PreBlock(ctx)
}

// BeginBlocker application updates every begin block
func (app *MilkyWayApp) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	return app.ModuleManager.BeginBlock(ctx)
}

// EndBlocker application updates every end block
func (app *MilkyWayApp) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	return app.ModuleManager.EndBlock(ctx)
}

// InitChainer application update at chain initialization
func (app *MilkyWayApp) InitChainer(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	if err := app.UpgradeKeeper.SetModuleVersionMap(ctx, app.ModuleManager.GetVersionMap()); err != nil {
		panic(err)
	}
	return app.ModuleManager.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads a particular height
func (app *MilkyWayApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *MilkyWayApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *MilkyWayApp) BlacklistedModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	// DO NOT REMOVE: StringMapKeys fixes non-deterministic map iteration
	for _, acc := range utils.StringMapKeys(maccPerms) {
		// don't blacklist stakeibc module account, so that it can ibc transfer tokens
		if acc == stakeibctypes.ModuleName ||
			acc == stakeibctypes.RewardCollectorName ||
			// don't blacklist liquidvesting module account, so that it can receive ibc transfers
			acc == liquidvestingtypes.ModuleName {
			continue
		}
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}
	return modAccAddrs
}

// LegacyAmino returns SimApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *MilkyWayApp) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// AppCodec returns Initia's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *MilkyWayApp) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns Initia's InterfaceRegistry
func (app *MilkyWayApp) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *MilkyWayApp) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *MilkyWayApp) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *MilkyWayApp) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *MilkyWayApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx

	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register new tendermint queries routes from grpc-gateway.
	cmtservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register node gRPC service for grpc-gateway.
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register grpc-gateway routes for all modules.
	app.BasicModuleManager.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register grpc-gateway routes for indexer module.
	app.indexerModule.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(apiSvr.Router)
	}
}

// Simulate customize gas simulation to add fee deduction gas amount.
func (app *MilkyWayApp) Simulate(txBytes []byte) (sdk.GasInfo, *sdk.Result, error) {
	gasInfo, result, err := app.BaseApp.Simulate(txBytes)
	gasInfo.GasUsed += FeeDeductionGasAmount
	return gasInfo, result, err
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *MilkyWayApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(
		app.BaseApp.GRPCQueryRouter(), clientCtx,
		app.Simulate, app.interfaceRegistry,
	)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *MilkyWayApp) RegisterTendermintService(clientCtx client.Context) {
	cmtservice.RegisterTendermintService(
		clientCtx,
		app.BaseApp.GRPCQueryRouter(),
		app.interfaceRegistry, app.Query,
	)
}

func (app *MilkyWayApp) RegisterNodeService(clientCtx client.Context, cfg config.Config) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter(), cfg)
}

// Configurator returns the app's configurator.
func (app *MilkyWayApp) Configurator() module.Configurator {
	return app.configurator
}

// registerUpgrade registers the given upgrade to be supported by the app
func (app *MilkyWayApp) registerUpgrade(upgrade upgrades.Upgrade) {
	app.UpgradeKeeper.SetUpgradeHandler(upgrade.Name(), upgrade.Handler())

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == upgrade.Name() && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		// Configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, upgrade.StoreUpgrades()))
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

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}

//////////////////////////////////////
// TestingApp functions

// GetBaseApp implements the TestingApp interface.
func (app *MilkyWayApp) GetBaseApp() *baseapp.BaseApp {
	return app.BaseApp
}

// GetAccountKeeper implements the TestingApp interface.
func (app *MilkyWayApp) GetAccountKeeper() *authkeeper.AccountKeeper {
	return app.AccountKeeper
}

// GetStakingKeeper implements the TestingApp interface.
// It returns opchild instead of original staking keeper.
func (app *MilkyWayApp) GetStakingKeeper() ibctestingtypes.StakingKeeper {
	return app.OPChildKeeper
}

// GetIBCKeeper implements the TestingApp interface.
func (app *MilkyWayApp) GetIBCKeeper() *ibckeeper.Keeper {
	return app.IBCKeeper
}

// GetICAControllerKeeper implements the TestingApp interface.
func (app *MilkyWayApp) GetICAControllerKeeper() *icacontrollerkeeper.Keeper {
	return app.ICAControllerKeeper
}

// GetICAAuthKeeper implements the TestingApp interface.
func (app *MilkyWayApp) GetICAAuthKeeper() *icaauthkeeper.Keeper {
	return app.ICAAuthKeeper
}

// GetScopedIBCKeeper implements the TestingApp interface.
func (app *MilkyWayApp) GetScopedIBCKeeper() capabilitykeeper.ScopedKeeper {
	return app.ScopedIBCKeeper
}

// TxConfig implements the TestingApp interface.
func (app *MilkyWayApp) TxConfig() client.TxConfig {
	return app.txConfig
}

func (app *MilkyWayApp) setupIndexer(db dbm.DB, appOpts servertypes.AppOptions, ac, vc address.Codec, appCodec codec.Codec) error {
	// initialize the indexer fake-keeper
	indexerConfig, err := indexerconfig.NewConfig(appOpts)
	if err != nil {
		panic(err)
	}
	app.indexerKeeper = indexerkeeper.NewKeeper(
		appCodec,
		"wasm",
		db,
		indexerConfig,
		ac,
		vc,
	)
	smBlock, err := blocksubmodule.NewBlockSubmodule(appCodec, app.indexerKeeper, app.OPChildKeeper)
	if err != nil {
		panic(err)
	}
	smTx, err := tx.NewTxSubmodule(appCodec, app.indexerKeeper)
	if err != nil {
		panic(err)
	}
	smPair, err := pair.NewPairSubmodule(appCodec, app.indexerKeeper, app.IBCKeeper.ChannelKeeper, app.TransferKeeper)
	if err != nil {
		panic(err)
	}
	smNft, err := nft.NewWasmNFTSubmodule(ac, appCodec, app.indexerKeeper, app.WasmKeeper, smPair)
	if err != nil {
		panic(err)
	}
	err = app.indexerKeeper.RegisterSubmodules(smBlock, smTx, smPair, smNft)
	if err != nil {
		panic(err)
	}

	app.indexerModule = indexermodule.NewAppModuleBasic(app.indexerKeeper)
	// Add your implementation here

	indexer, err := indexer.NewIndexer(app.GetBaseApp().Logger(), app.indexerKeeper)
	if err != nil || indexer == nil {
		return nil
	}

	if err = indexer.Validate(); err != nil {
		return err
	}

	if err = indexer.Prepare(nil); err != nil {
		return err
	}

	if err = app.indexerKeeper.Seal(); err != nil {
		return err
	}

	if err = indexer.Start(nil); err != nil {
		return err
	}

	streamingManager := storetypes.StreamingManager{
		ABCIListeners: []storetypes.ABCIListener{indexer},
		StopNodeOnErr: true,
	}
	app.SetStreamingManager(streamingManager)

	return nil
}

// Close closes the underlying baseapp, the oracle service, and the prometheus server if required.
// This method blocks on the closure of both the prometheus server, and the oracle-service
func (app *MilkyWayApp) Close() error {
	if app.indexerKeeper != nil {
		if err := app.indexerKeeper.Close(); err != nil {
			return err
		}
	}

	if err := app.BaseApp.Close(); err != nil {
		return err
	}

	return nil
}
