package keepers

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/x/group"
	"github.com/cosmos/gogoproto/proto"
	marketmapkeeper "github.com/skip-mev/connect/v2/x/marketmap/keeper"
	marketmaptypes "github.com/skip-mev/connect/v2/x/marketmap/types"
	oraclekeeper "github.com/skip-mev/connect/v2/x/oracle/keeper"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"
	feemarketkeeper "github.com/skip-mev/feemarket/x/feemarket/keeper"
	feemarkettypes "github.com/skip-mev/feemarket/x/feemarket/types"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	groupkeeper "github.com/cosmos/cosmos-sdk/x/group/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	pfmroutertypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"
	ratelimittypes "github.com/cosmos/ibc-apps/modules/rate-limiting/v8/types"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ibcfee "github.com/cosmos/ibc-go/v8/modules/apps/29-fee"
	ibcfeekeeper "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/keeper"
	ibcfeetypes "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/types"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	providertypes "github.com/cosmos/interchain-security/v6/x/ccv/provider/types"

	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	pfmrouter "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward"
	pfmrouterkeeper "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/keeper"
	ratelimit "github.com/cosmos/ibc-apps/modules/rate-limiting/v8"
	ratelimitkeeper "github.com/cosmos/ibc-apps/modules/rate-limiting/v8/keeper"
	icsprovider "github.com/cosmos/interchain-security/v6/x/ccv/provider"
	icsproviderkeeper "github.com/cosmos/interchain-security/v6/x/ccv/provider/keeper"

	assetskeeper "github.com/milkyway-labs/milkyway/v6/x/assets/keeper"
	assetstypes "github.com/milkyway-labs/milkyway/v6/x/assets/types"
	"github.com/milkyway-labs/milkyway/v6/x/liquidvesting"
	liquidvestingkeeper "github.com/milkyway-labs/milkyway/v6/x/liquidvesting/keeper"
	liquidvestingtypes "github.com/milkyway-labs/milkyway/v6/x/liquidvesting/types"
	operatorskeeper "github.com/milkyway-labs/milkyway/v6/x/operators/keeper"
	operatorstypes "github.com/milkyway-labs/milkyway/v6/x/operators/types"
	poolskeeper "github.com/milkyway-labs/milkyway/v6/x/pools/keeper"
	poolstypes "github.com/milkyway-labs/milkyway/v6/x/pools/types"
	restakingkeeper "github.com/milkyway-labs/milkyway/v6/x/restaking/keeper"
	restakingtypes "github.com/milkyway-labs/milkyway/v6/x/restaking/types"
	rewardskeeper "github.com/milkyway-labs/milkyway/v6/x/rewards/keeper"
	rewardstypes "github.com/milkyway-labs/milkyway/v6/x/rewards/types"
	serviceskeeper "github.com/milkyway-labs/milkyway/v6/x/services/keeper"
	servicestypes "github.com/milkyway-labs/milkyway/v6/x/services/types"
)

type AppKeepers struct {
	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// keepers
	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            bankkeeper.BaseKeeper
	CapabilityKeeper      *capabilitykeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	MintKeeper            mintkeeper.Keeper
	DistrKeeper           distrkeeper.Keeper
	GovKeeper             *govkeeper.Keeper
	GroupKeeper           groupkeeper.Keeper
	CrisisKeeper          *crisiskeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	ParamsKeeper          paramskeeper.Keeper
	WasmKeeper            wasmkeeper.Keeper
	EvidenceKeeper        evidencekeeper.Keeper
	AuthzKeeper           authzkeeper.Keeper
	ConsensusParamsKeeper consensusparamkeeper.Keeper
	FeeGrantKeeper        feegrantkeeper.Keeper

	// Skip
	MarketMapKeeper *marketmapkeeper.Keeper
	OracleKeeper    *oraclekeeper.Keeper
	FeeMarketKeeper *feemarketkeeper.Keeper

	// IBC
	IBCKeeper       *ibckeeper.Keeper
	TransferKeeper  ibctransferkeeper.Keeper
	PFMRouterKeeper *pfmrouterkeeper.Keeper
	RateLimitKeeper *ratelimitkeeper.Keeper

	// ICS
	ProviderKeeper icsproviderkeeper.Keeper

	// Custom
	ServicesKeeper      *serviceskeeper.Keeper
	OperatorsKeeper     *operatorskeeper.Keeper
	PoolsKeeper         *poolskeeper.Keeper
	RestakingKeeper     *restakingkeeper.Keeper
	AssetsKeeper        *assetskeeper.Keeper
	RewardsKeeper       *rewardskeeper.Keeper
	LiquidVestingKeeper *liquidvestingkeeper.Keeper

	// Modules
	IBCFeeKeeper    ibcfeekeeper.Keeper
	TransferModule  transfer.AppModule
	PFMRouterModule pfmrouter.AppModule
	RateLimitModule ratelimit.AppModule
	ProviderModule  icsprovider.AppModule

	// make scoped keepers public for test purposes
	ScopedIBCKeeper         capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper    capabilitykeeper.ScopedKeeper
	ScopedICSproviderkeeper capabilitykeeper.ScopedKeeper
	scopedWasmKeeper        capabilitykeeper.ScopedKeeper
}

func NewAppKeeper(
	appCodec codec.Codec,
	bApp *baseapp.BaseApp,
	legacyAmino *codec.LegacyAmino,
	maccPerms map[string][]string,
	blockedAddress map[string]bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	logger log.Logger,
	appOpts servertypes.AppOptions,
	wasmOpts []wasmkeeper.Option,
) AppKeepers {
	appKeepers := AppKeepers{}

	// Create codecs
	addressCodec := address.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())

	// Set keys KVStoreKey, TransientStoreKey, MemoryStoreKey
	appKeepers.GenerateKeys()

	if err := bApp.RegisterStreamingServices(appOpts, appKeepers.keys); err != nil {
		logger.Error("failed to load state streaming", "err", err)
		os.Exit(1)
	}

	appKeepers.ParamsKeeper = initParamsKeeper(
		appCodec,
		legacyAmino,
		appKeepers.keys[paramstypes.StoreKey],
		appKeepers.tkeys[paramstypes.TStoreKey],
	)

	// set the BaseApp's parameter store
	appKeepers.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[consensusparamtypes.StoreKey]),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		runtime.EventService{},
	)
	bApp.SetParamStore(appKeepers.ConsensusParamsKeeper.ParamsStore)

	// add capability keeper and ScopeToModule for ibc module
	appKeepers.CapabilityKeeper = capabilitykeeper.NewKeeper(
		appCodec,
		appKeepers.keys[capabilitytypes.StoreKey],
		appKeepers.memKeys[capabilitytypes.MemStoreKey],
	)

	appKeepers.ScopedIBCKeeper = appKeepers.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	appKeepers.ScopedTransferKeeper = appKeepers.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	appKeepers.ScopedICSproviderkeeper = appKeepers.CapabilityKeeper.ScopeToModule(providertypes.ModuleName)
	appKeepers.scopedWasmKeeper = appKeepers.CapabilityKeeper.ScopeToModule(wasmtypes.ModuleName)

	// Applications that wish to enforce statically created ScopedKeepers should call `Seal` after creating
	// their scoped modules in `NewApp` with `ScopeToModule`
	appKeepers.CapabilityKeeper.Seal()

	// Add normal keepers
	appKeepers.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		maccPerms,
		addressCodec,
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[banktypes.StoreKey]),
		appKeepers.AccountKeeper,
		blockedAddress,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		logger,
	)

	appKeepers.CrisisKeeper = crisiskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[crisistypes.StoreKey]),
		invCheckPeriod,
		appKeepers.BankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		appKeepers.AccountKeeper.AddressCodec(),
	)

	appKeepers.AuthzKeeper = authzkeeper.NewKeeper(
		runtime.NewKVStoreService(appKeepers.keys[authzkeeper.StoreKey]),
		appCodec,
		bApp.MsgServiceRouter(),
		appKeepers.AccountKeeper,
	)

	appKeepers.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[feegrant.StoreKey]),
		appKeepers.AccountKeeper,
	)

	appKeepers.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[stakingtypes.StoreKey]),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	)

	appKeepers.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[distrtypes.StoreKey]),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		legacyAmino,
		runtime.NewKVStoreService(appKeepers.keys[slashingtypes.StoreKey]),
		appKeepers.StakingKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	appKeepers.StakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			appKeepers.DistrKeeper.Hooks(),
			appKeepers.SlashingKeeper.Hooks(),
			appKeepers.ProviderKeeper.Hooks(),
		),
	)

	appKeepers.FeeMarketKeeper = feemarketkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[feemarkettypes.StoreKey],
		appKeepers.AccountKeeper,
		&DefaultFeemarketDenomResolver{},
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.MarketMapKeeper = marketmapkeeper.NewKeeper(
		runtime.NewKVStoreService(appKeepers.keys[marketmaptypes.StoreKey]),
		appCodec,
		authtypes.NewModuleAddress(govtypes.ModuleName),
	)

	oracleKeeper := oraclekeeper.NewKeeper(
		runtime.NewKVStoreService(appKeepers.keys[oracletypes.StoreKey]),
		appCodec,
		appKeepers.MarketMapKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName),
	)
	appKeepers.OracleKeeper = &oracleKeeper

	// Add the oracle keeper as a hook to market map keeper so new market map entries can be created
	// and propagated to the oracle keeper.
	appKeepers.MarketMapKeeper.SetHooks(appKeepers.OracleKeeper.Hooks())

	// UpgradeKeeper must be created before IBCKeeper
	appKeepers.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		runtime.NewKVStoreService(appKeepers.keys[upgradetypes.StoreKey]),
		appCodec,
		homePath,
		bApp,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.GroupKeeper = groupkeeper.NewKeeper(
		appKeepers.keys[group.StoreKey],
		appCodec,
		bApp.MsgServiceRouter(),
		appKeepers.AccountKeeper,
		group.DefaultConfig(),
	)

	// UpgradeKeeper must be created before IBCKeeper
	appKeepers.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibcexported.StoreKey],
		appKeepers.GetSubspace(ibcexported.ModuleName),
		appKeepers.StakingKeeper,
		appKeepers.UpgradeKeeper,
		appKeepers.ScopedIBCKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.ProviderKeeper = icsproviderkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[providertypes.StoreKey],
		appKeepers.GetSubspace(providertypes.ModuleName),
		appKeepers.ScopedICSproviderkeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.PortKeeper,
		appKeepers.IBCKeeper.ConnectionKeeper,
		appKeepers.IBCKeeper.ClientKeeper,
		appKeepers.StakingKeeper,
		appKeepers.SlashingKeeper,
		appKeepers.AccountKeeper,
		appKeepers.DistrKeeper,
		appKeepers.BankKeeper,
		govkeeper.Keeper{}, // cyclic dependency between provider and governance, will be set later
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
		authtypes.FeeCollectorName,
	)

	// gov depends on provider, so needs to be set after
	govConfig := govtypes.DefaultConfig()
	// set the MaxMetadataLen for proposals to the same value as it was pre-sdk v0.47.x
	govConfig.MaxMetadataLen = 10200
	appKeepers.GovKeeper = govkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[govtypes.StoreKey]),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		// use the ProviderKeeper as StakingKeeper for gov
		// because governance should be based on the consensus-active validators
		appKeepers.ProviderKeeper,
		appKeepers.DistrKeeper,
		bApp.MsgServiceRouter(),
		govConfig,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// mint keeper must be created after provider keeper
	appKeepers.MintKeeper = mintkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[minttypes.StoreKey]),
		appKeepers.ProviderKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.ProviderKeeper.SetGovKeeper(*appKeepers.GovKeeper)

	appKeepers.ProviderModule = icsprovider.NewAppModule(
		&appKeepers.ProviderKeeper,
		appKeepers.GetSubspace(providertypes.ModuleName),
		appKeepers.keys[providertypes.StoreKey],
	)

	// Register the proposal types
	// Deprecated: Avoid adding new handlers, instead use the new proposal flow
	// by granting the governance module the right to execute the message.
	// See: https://docs.cosmos.network/main/modules/gov#proposal-messages
	govRouter := govv1beta1.NewRouter()
	govRouter.
		AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(appKeepers.ParamsKeeper))

	// Set legacy router for backwards compatibility with gov v1beta1
	appKeepers.GovKeeper.SetLegacyRouter(govRouter)

	appKeepers.GovKeeper = appKeepers.GovKeeper.SetHooks(
		govtypes.NewMultiGovHooks(
			appKeepers.ProviderKeeper.Hooks(),
		),
	)

	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[evidencetypes.StoreKey]),
		appKeepers.StakingKeeper,
		appKeepers.SlashingKeeper,
		appKeepers.AccountKeeper.AddressCodec(),
		runtime.ProvideCometInfoService(),
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	appKeepers.EvidenceKeeper = *evidenceKeeper

	appKeepers.IBCFeeKeeper = ibcfeekeeper.NewKeeper(
		appCodec, appKeepers.keys[ibcfeetypes.StoreKey],
		appKeepers.IBCKeeper.ChannelKeeper, // may be replaced with IBC middleware
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.PortKeeper, appKeepers.AccountKeeper, appKeepers.BankKeeper,
	)

	govAuthority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	// Create RateLimit keeper
	appKeepers.RateLimitKeeper = ratelimitkeeper.NewKeeper(
		appCodec, // BinaryCodec
		runtime.NewKVStoreService(appKeepers.keys[ratelimittypes.StoreKey]), // StoreKey
		appKeepers.GetSubspace(ratelimittypes.ModuleName),                   // param Subspace
		govAuthority, // authority
		appKeepers.BankKeeper,
		appKeepers.IBCKeeper.ChannelKeeper, // ChannelKeeper
		appKeepers.IBCFeeKeeper,            // ICS4Wrapper
	)

	// PFMRouterKeeper must be created before TransferKeeper
	appKeepers.PFMRouterKeeper = pfmrouterkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[pfmroutertypes.StoreKey],
		nil, // Will be zero-value here. Reference is set later on with SetTransferKeeper.
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.DistrKeeper,
		appKeepers.BankKeeper,
		appKeepers.RateLimitKeeper, // ICS4Wrapper
		govAuthority,
	)

	appKeepers.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibctransfertypes.StoreKey],
		appKeepers.GetSubspace(ibctransfertypes.ModuleName),
		appKeepers.PFMRouterKeeper, // ISC4 Wrapper: PFM Router middleware
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.ScopedTransferKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// ----------------------
	// --- Custom modules ---
	// ----------------------

	// Custom modules
	appKeepers.ServicesKeeper = serviceskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[servicestypes.StoreKey]),
		appKeepers.AccountKeeper,
		appKeepers.DistrKeeper,
		govAuthority,
	)
	appKeepers.OperatorsKeeper = operatorskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[operatorstypes.StoreKey]),
		appKeepers.AccountKeeper,
		appKeepers.DistrKeeper,
		govAuthority,
	)
	appKeepers.PoolsKeeper = poolskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[poolstypes.StoreKey]),
		appKeepers.AccountKeeper,
	)
	appKeepers.AssetsKeeper = assetskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[assetstypes.StoreKey]),
		govAuthority,
	)
	appKeepers.RestakingKeeper = restakingkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[restakingtypes.StoreKey]),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.PoolsKeeper,
		appKeepers.OperatorsKeeper,
		appKeepers.ServicesKeeper,
		appKeepers.OracleKeeper,
		appKeepers.AssetsKeeper,
		govAuthority,
	)
	appKeepers.RewardsKeeper = rewardskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[rewardstypes.StoreKey]),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.DistrKeeper,
		appKeepers.OracleKeeper,
		appKeepers.PoolsKeeper,
		appKeepers.OperatorsKeeper,
		appKeepers.ServicesKeeper,
		appKeepers.RestakingKeeper,
		appKeepers.AssetsKeeper,
		govAuthority,
	)

	// Set hooks based on the rewards keeper
	appKeepers.LiquidVestingKeeper = liquidvestingkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[liquidvestingtypes.StoreKey]),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.OperatorsKeeper,
		appKeepers.PoolsKeeper,
		appKeepers.ServicesKeeper,
		appKeepers.RestakingKeeper,
		authtypes.NewModuleAddress(liquidvestingtypes.ModuleName).String(),
		govAuthority,
	)

	// Set the restrictions on sending tokens
	appKeepers.BankKeeper.AppendSendRestriction(appKeepers.LiquidVestingKeeper.SendRestrictionFn)

	// Set the hooks up to this point
	appKeepers.PoolsKeeper.SetHooks(
		appKeepers.RewardsKeeper.PoolsHooks(),
	)
	appKeepers.OperatorsKeeper.SetHooks(operatorstypes.NewMultiOperatorsHooks(
		appKeepers.RestakingKeeper.OperatorsHooks(),
		appKeepers.RewardsKeeper.OperatorsHooks(),
	))
	appKeepers.ServicesKeeper.SetHooks(servicestypes.NewMultiServicesHooks(
		appKeepers.RestakingKeeper.ServicesHooks(),
		appKeepers.RewardsKeeper.ServicesHooks(),
	))
	appKeepers.RestakingKeeper.SetHooks(
		appKeepers.RewardsKeeper.RestakingHooks(),
	)
	appKeepers.RestakingKeeper.SetRestakeRestriction(
		appKeepers.LiquidVestingKeeper.RestakeRestrictionFn,
	)

	// Must be called on PFMRouter AFTER TransferKeeper initialized
	appKeepers.PFMRouterKeeper.SetTransferKeeper(appKeepers.TransferKeeper)

	wasmDir := homePath
	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic("error while reading wasm config: " + err.Error())
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
		Stargate: wasmkeeper.AcceptListStargateQuerier(queryAllowlist, bApp.GRPCQueryRouter(), appCodec),
	}))

	appKeepers.WasmKeeper = wasmkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[wasmtypes.StoreKey]),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		distrkeeper.NewQuerier(appKeepers.DistrKeeper),
		appKeepers.IBCFeeKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.PortKeeper,
		appKeepers.scopedWasmKeeper,
		appKeepers.TransferKeeper,
		bApp.MsgServiceRouter(),
		bApp.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		wasmkeeper.BuiltInCapabilities(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		wasmOpts...,
	)

	// Middleware Stacks
	appKeepers.TransferModule = transfer.NewAppModule(appKeepers.TransferKeeper)
	appKeepers.PFMRouterModule = pfmrouter.NewAppModule(appKeepers.PFMRouterKeeper, appKeepers.GetSubspace(pfmroutertypes.ModuleName))
	appKeepers.RateLimitModule = ratelimit.NewAppModule(appCodec, *appKeepers.RateLimitKeeper)

	// Create Transfer Stack (from bottom to top of stack)
	// - core IBC
	// - ibcfee
	// - ratelimit
	// - pfm
	// - provider
	// - transfer
	//
	// This is how transfer stack will work in the end:
	// * RecvPacket -> IBC core -> Fee -> RateLimit -> PFM -> Provider -> Transfer (AddRoute)
	// * SendPacket -> Transfer -> Provider -> PFM -> RateLimit -> Fee -> IBC core (ICS4Wrapper)

	var transferStack porttypes.IBCModule
	transferStack = transfer.NewIBCModule(appKeepers.TransferKeeper)
	transferStack = icsprovider.NewIBCMiddleware(
		transferStack,
		appKeepers.ProviderKeeper,
	)
	transferStack = pfmrouter.NewIBCMiddleware(
		transferStack,
		appKeepers.PFMRouterKeeper,
		0, // retries on timeout
		pfmrouterkeeper.DefaultForwardTransferPacketTimeoutTimestamp,
		pfmrouterkeeper.DefaultRefundTransferPacketTimeoutTimestamp,
	)
	transferStack = liquidvesting.NewIBCMiddleware(
		transferStack,
		appKeepers.LiquidVestingKeeper,
	)
	transferStack = ratelimit.NewIBCMiddleware(
		*appKeepers.RateLimitKeeper,
		transferStack,
	)
	transferStack = ibcfee.NewIBCMiddleware(
		transferStack,
		appKeepers.IBCFeeKeeper,
	)

	var wasmStack porttypes.IBCModule
	wasmStack = wasm.NewIBCHandler(appKeepers.WasmKeeper, appKeepers.IBCKeeper.ChannelKeeper, appKeepers.IBCFeeKeeper)
	wasmStack = ibcfee.NewIBCMiddleware(wasmStack, appKeepers.IBCFeeKeeper)

	// Create IBC Router & seal
	ibcRouter := porttypes.NewRouter().
		AddRoute(ibctransfertypes.ModuleName, transferStack).
		AddRoute(providertypes.ModuleName, appKeepers.ProviderModule).
		AddRoute(wasmtypes.ModuleName, wasmStack)

	appKeepers.IBCKeeper.SetRouter(ibcRouter)

	return appKeepers
}

// GetSubspace returns a param subspace for a given module name.
func (appKeepers *AppKeepers) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, ok := appKeepers.ParamsKeeper.GetSubspace(moduleName)
	if !ok {
		panic("couldn't load subspace for module: " + moduleName)
	}
	return subspace
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	// register the key tables for legacy param subspaces
	keyTable := ibcclienttypes.ParamKeyTable()
	keyTable.RegisterParamSet(&ibcconnectiontypes.Params{})
	paramsKeeper.Subspace(authtypes.ModuleName).WithKeyTable(authtypes.ParamKeyTable())         //nolint: staticcheck
	paramsKeeper.Subspace(stakingtypes.ModuleName).WithKeyTable(stakingtypes.ParamKeyTable())   //nolint: staticcheck // SA1019
	paramsKeeper.Subspace(banktypes.ModuleName).WithKeyTable(banktypes.ParamKeyTable())         //nolint: staticcheck // SA1019
	paramsKeeper.Subspace(minttypes.ModuleName).WithKeyTable(minttypes.ParamKeyTable())         //nolint: staticcheck // SA1019
	paramsKeeper.Subspace(distrtypes.ModuleName).WithKeyTable(distrtypes.ParamKeyTable())       //nolint: staticcheck // SA1019
	paramsKeeper.Subspace(slashingtypes.ModuleName).WithKeyTable(slashingtypes.ParamKeyTable()) //nolint: staticcheck // SA1019
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govv1.ParamKeyTable())              //nolint: staticcheck // SA1019
	paramsKeeper.Subspace(crisistypes.ModuleName).WithKeyTable(crisistypes.ParamKeyTable())     //nolint: staticcheck // SA1019
	paramsKeeper.Subspace(ibcexported.ModuleName).WithKeyTable(keyTable)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName).WithKeyTable(ibctransfertypes.ParamKeyTable())
	paramsKeeper.Subspace(pfmroutertypes.ModuleName).WithKeyTable(pfmroutertypes.ParamKeyTable())
	paramsKeeper.Subspace(ratelimittypes.ModuleName).WithKeyTable(ratelimittypes.ParamKeyTable())
	paramsKeeper.Subspace(providertypes.ModuleName).WithKeyTable(providertypes.ParamKeyTable())
	paramsKeeper.Subspace(wasmtypes.ModuleName)

	return paramsKeeper
}

type DefaultFeemarketDenomResolver struct{}

func (r *DefaultFeemarketDenomResolver) ConvertToDenom(_ sdk.Context, coin sdk.DecCoin, denom string) (sdk.DecCoin, error) {
	if coin.Denom == denom {
		return coin, nil
	}

	return sdk.DecCoin{}, fmt.Errorf("error resolving denom")
}

func (r *DefaultFeemarketDenomResolver) ExtraDenoms(_ sdk.Context) ([]string, error) {
	return []string{}, nil
}
