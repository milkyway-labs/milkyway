package milkyway

import (
	hyperlanetypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	warptypes "github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/skip-mev/connect/v2/x/marketmap"
	marketmaptypes "github.com/skip-mev/connect/v2/x/marketmap/types"
	"github.com/skip-mev/connect/v2/x/oracle"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"

	hyperlane "github.com/bcp-innovations/hyperlane-cosmos/x/core"
	"github.com/bcp-innovations/hyperlane-cosmos/x/warp"
	pfmroutertypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"
	ratelimittypes "github.com/cosmos/ibc-apps/modules/rate-limiting/v8/types"
	"github.com/cosmos/ibc-go/modules/capability"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ibcfee "github.com/cosmos/ibc-go/v8/modules/apps/29-fee"
	ibcfeetypes "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v8/modules/core"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	no_valupdates_genutil "github.com/cosmos/interchain-security/v6/x/ccv/no_valupdates_genutil"
	no_valupdates_staking "github.com/cosmos/interchain-security/v6/x/ccv/no_valupdates_staking"
	providertypes "github.com/cosmos/interchain-security/v6/x/ccv/provider/types"

	"cosmossdk.io/x/evidence"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	feegrantmodule "cosmossdk.io/x/feegrant/module"
	"cosmossdk.io/x/upgrade"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	sdkparams "github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	"github.com/milkyway-labs/milkyway/v12/x/assets"
	assetstypes "github.com/milkyway-labs/milkyway/v12/x/assets/types"
	"github.com/milkyway-labs/milkyway/v12/x/bank"
	ibchooks "github.com/milkyway-labs/milkyway/v12/x/ibc-hooks"
	ibchookstypes "github.com/milkyway-labs/milkyway/v12/x/ibc-hooks/types"
	"github.com/milkyway-labs/milkyway/v12/x/investors"
	investorstypes "github.com/milkyway-labs/milkyway/v12/x/investors/types"
	"github.com/milkyway-labs/milkyway/v12/x/liquidvesting"
	liquidvestingtypes "github.com/milkyway-labs/milkyway/v12/x/liquidvesting/types"
	"github.com/milkyway-labs/milkyway/v12/x/operators"
	operatorstypes "github.com/milkyway-labs/milkyway/v12/x/operators/types"
	"github.com/milkyway-labs/milkyway/v12/x/pools"
	poolstypes "github.com/milkyway-labs/milkyway/v12/x/pools/types"
	"github.com/milkyway-labs/milkyway/v12/x/restaking"
	restakingtypes "github.com/milkyway-labs/milkyway/v12/x/restaking/types"
	"github.com/milkyway-labs/milkyway/v12/x/rewards"
	rewardstypes "github.com/milkyway-labs/milkyway/v12/x/rewards/types"
	"github.com/milkyway-labs/milkyway/v12/x/services"
	servicestypes "github.com/milkyway-labs/milkyway/v12/x/services/types"
	"github.com/milkyway-labs/milkyway/v12/x/tokenfactory"
	tokenfactorytypes "github.com/milkyway-labs/milkyway/v12/x/tokenfactory/types"
)

var MaccPerms = map[string][]string{
	authtypes.FeeCollectorName:        nil,
	distrtypes.ModuleName:             nil,
	minttypes.ModuleName:              {authtypes.Minter},
	stakingtypes.BondedPoolName:       {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName:    {authtypes.Burner, authtypes.Staking},
	govtypes.ModuleName:               {authtypes.Burner},
	ibctransfertypes.ModuleName:       {authtypes.Minter, authtypes.Burner},
	ibcfeetypes.ModuleName:            nil,
	providertypes.ConsumerRewardsPool: nil,
	wasmtypes.ModuleName:              {authtypes.Burner},
	tokenfactorytypes.ModuleName:      {authtypes.Minter, authtypes.Burner},

	// Skip
	oracletypes.ModuleName: nil,

	// Osmosis IBC hooks
	ibchookstypes.ModuleName: nil,

	// MilkyWay permissions
	rewardstypes.RewardsPoolName:  nil,
	liquidvestingtypes.ModuleName: {authtypes.Minter, authtypes.Burner},
	investorstypes.ModuleName:     nil,

	// Warp module
	warptypes.ModuleName:      {authtypes.Minter, authtypes.Burner},
	hyperlanetypes.ModuleName: nil,
}

func appModules(
	app *MilkyWayApp,
	appCodec codec.Codec,
	txConfig client.TxEncodingConfig,
	skipGenesisInvariants bool,
) []module.AppModule {
	return []module.AppModule{
		no_valupdates_genutil.NewAppModule(
			app.AccountKeeper,
			app.StakingKeeper,
			app,
			txConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, nil, app.GetSubspace(authtypes.ModuleName)),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		crisis.NewAppModule(app.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisistypes.ModuleName)),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper, nil, app.GetSubspace(minttypes.ModuleName)),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName), app.interfaceRegistry),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		no_valupdates_staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		upgrade.NewAppModule(app.UpgradeKeeper, app.AccountKeeper.AddressCodec()),
		evidence.NewAppModule(app.EvidenceKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		ibc.NewAppModule(app.IBCKeeper),
		ibctm.NewAppModule(),
		sdkparams.NewAppModule(app.ParamsKeeper),
		consensus.NewAppModule(appCodec, app.ConsensusParamsKeeper),
		wasm.NewAppModule(appCodec, &app.AppKeepers.WasmKeeper, app.AppKeepers.StakingKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.MsgServiceRouter(), app.GetSubspace(wasmtypes.ModuleName)),
		tokenfactory.NewAppModule(*app.TokenFactoryKeeper, app.AccountKeeper, app.BankKeeper),

		// Skip modules
		oracle.NewAppModule(appCodec, *app.OracleKeeper),
		marketmap.NewAppModule(appCodec, app.MarketMapKeeper),

		// Hyperlane modules
		hyperlane.NewAppModule(appCodec, app.HyperlaneKeeper),
		warp.NewAppModule(appCodec, app.WarpKeeper),

		// IBC Modules
		ibcfee.NewAppModule(app.IBCFeeKeeper),
		app.TransferModule,
		app.PFMRouterModule,
		app.RateLimitModule,
		app.ProviderModule,
		ibchooks.NewAppModule(app.AccountKeeper, *app.IBCHooksKeeper),

		// MilkyWay modules
		services.NewAppModule(appCodec, app.ServicesKeeper, app.AccountKeeper, app.BankKeeper),
		operators.NewAppModule(appCodec, app.OperatorsKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		pools.NewAppModule(appCodec, app.PoolsKeeper),
		restaking.NewAppModule(appCodec, app.RestakingKeeper, app.AccountKeeper, app.BankKeeper, app.PoolsKeeper, app.OperatorsKeeper, app.ServicesKeeper),
		assets.NewAppModule(appCodec, app.AssetsKeeper),
		rewards.NewAppModule(appCodec, app.RewardsKeeper, app.AccountKeeper, app.BankKeeper, app.PoolsKeeper, app.OperatorsKeeper, app.ServicesKeeper),
		liquidvesting.NewAppModule(appCodec, app.LiquidVestingKeeper),
		investors.NewAppModule(appCodec, app.InvestorsKeeper),
	}
}

// ModuleBasics defines the module BasicManager that is in charge of setting up basic,
// non-dependant module elements, such as codec registration
// and genesis verification.
func newBasicManagerFromManager(app *MilkyWayApp) module.BasicManager {
	basicManager := module.NewBasicManagerFromManager(
		app.mm,
		map[string]module.AppModuleBasic{
			genutiltypes.ModuleName: genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
			govtypes.ModuleName: gov.NewAppModuleBasic(
				[]govclient.ProposalHandler{
					paramsclient.ProposalHandler,
				},
			),
		})
	basicManager.RegisterLegacyAminoCodec(app.legacyAmino)
	basicManager.RegisterInterfaces(app.interfaceRegistry)
	return basicManager
}

// simulationModules returns modules for simulation manager
// define the order of the modules for deterministic simulations
func simulationModules(
	app *MilkyWayApp,
	appCodec codec.Codec,
	_ bool,
) []module.AppModuleSimulation {
	return []module.AppModuleSimulation{
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper, nil, app.GetSubspace(minttypes.ModuleName)),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName), app.interfaceRegistry),
		sdkparams.NewAppModule(app.ParamsKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		wasm.NewAppModule(appCodec, &app.AppKeepers.WasmKeeper, app.AppKeepers.StakingKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.MsgServiceRouter(), app.GetSubspace(wasmtypes.ModuleName)),
		ibc.NewAppModule(app.IBCKeeper),
		app.TransferModule,

		// MilkyWay modules
		services.NewAppModule(appCodec, app.ServicesKeeper, app.AccountKeeper, app.BankKeeper),
		operators.NewAppModule(appCodec, app.OperatorsKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		restaking.NewAppModule(appCodec, app.RestakingKeeper, app.AccountKeeper, app.BankKeeper, app.PoolsKeeper, app.OperatorsKeeper, app.ServicesKeeper),
		rewards.NewAppModule(appCodec, app.RewardsKeeper, app.AccountKeeper, app.BankKeeper, app.PoolsKeeper, app.OperatorsKeeper, app.ServicesKeeper),
	}
}

/*
orderBeginBlockers tells the app's module manager how to set the order of
BeginBlockers, which are run at the beginning of every block.

Interchain Security Requirements:
During begin block slashing happens after distr.BeginBlocker so that
there is nothing left over in the validator fee pool, so as to keep the
CanWithdrawInvariant invariant.
NOTE: staking module is required if HistoricalEntries param > 0
NOTE: capability module's beginblocker must come before any modules using capabilities (e.g. IBC)
*/
func orderBeginBlockers() []string {
	return []string{
		capabilitytypes.ModuleName,
		minttypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		pfmroutertypes.ModuleName,
		ratelimittypes.ModuleName,
		ibcfeetypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		providertypes.ModuleName,
		consensusparamtypes.ModuleName,
		wasmtypes.ModuleName,
		ratelimittypes.ModuleName,

		// Skip modules
		oracletypes.ModuleName,
		marketmaptypes.ModuleName,

		// MilkyWay modules
		rewardstypes.ModuleName,
		servicestypes.ModuleName,
		operatorstypes.ModuleName,
		poolstypes.ModuleName,
		restakingtypes.ModuleName,
		investorstypes.ModuleName,
	}
}

/*
Interchain Security Requirements:
- provider.EndBlock gets validator updates from the staking module;
thus, staking.EndBlock must be executed before provider.EndBlock;
- creating a new consumer chain requires the following order,
CreateChildClient(), staking.EndBlock, provider.EndBlock;
thus, gov.EndBlock must be executed before staking.EndBlock
*/
func orderEndBlockers() []string {
	return []string{
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		pfmroutertypes.ModuleName,
		ratelimittypes.ModuleName,
		capabilitytypes.ModuleName,
		ibcfeetypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		providertypes.ModuleName,
		consensusparamtypes.ModuleName,
		wasmtypes.ModuleName,

		// Skip modules
		marketmaptypes.ModuleName,
		oracletypes.ModuleName,

		// MilkyWay modules
		servicestypes.ModuleName,
		operatorstypes.ModuleName,
		poolstypes.ModuleName,
		restakingtypes.ModuleName,
		liquidvestingtypes.ModuleName,
	}
}

/*
NOTE: The genutils module must occur after staking so that pools are
properly initialized with tokens from genesis accounts.
NOTE: The genutils module must also occur after auth so that it can access the params from auth.
NOTE: Capability module must occur first so that it can initialize any capabilities
so that other modules that want to create or claim capabilities afterwards in InitChain
can do so safely.
*/
func orderInitBlockers() []string {
	return []string{
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibcexported.ModuleName,
		ibcfeetypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		pfmroutertypes.ModuleName,
		ratelimittypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		tokenfactorytypes.ModuleName,

		// Skip modules
		oracletypes.ModuleName,
		marketmaptypes.ModuleName,

		// Hyperlane modules
		hyperlanetypes.ModuleName,
		warptypes.ModuleName,

		// MilkyWay modules
		servicestypes.ModuleName,
		operatorstypes.ModuleName,
		poolstypes.ModuleName,
		restakingtypes.ModuleName,
		assetstypes.ModuleName,
		rewardstypes.ModuleName,
		liquidvestingtypes.ModuleName,
		investorstypes.ModuleName,

		providertypes.ModuleName,
		consensusparamtypes.ModuleName,
		wasmtypes.ModuleName,
		ibchookstypes.ModuleName,

		// crisis needs to be last so that the genesis state is consistent
		// when it checks invariants
		crisistypes.ModuleName,
	}
}
