package keepers

import (
	"cosmossdk.io/x/evidence"
	feegrantmodule "cosmossdk.io/x/feegrant/module"
	"cosmossdk.io/x/upgrade"
	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	sdkparams "github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	pfmrouter "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward"
	ratelimit "github.com/cosmos/ibc-apps/modules/rate-limiting/v8"
	"github.com/cosmos/ibc-go/modules/capability"
	ibcfee "github.com/cosmos/ibc-go/v8/modules/apps/29-fee"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibc "github.com/cosmos/ibc-go/v8/modules/core"
	ibctm "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	no_valupdates_staking "github.com/cosmos/interchain-security/v6/x/ccv/no_valupdates_staking"
	icsprovider "github.com/cosmos/interchain-security/v6/x/ccv/provider"
	"github.com/skip-mev/connect/v2/x/marketmap"
	"github.com/skip-mev/connect/v2/x/oracle"
	"github.com/skip-mev/feemarket/x/feemarket"

	"github.com/milkyway-labs/milkyway/v3/x/assets"
	"github.com/milkyway-labs/milkyway/v3/x/liquidvesting"
	"github.com/milkyway-labs/milkyway/v3/x/operators"
	"github.com/milkyway-labs/milkyway/v3/x/pools"
	"github.com/milkyway-labs/milkyway/v3/x/restaking"
	"github.com/milkyway-labs/milkyway/v3/x/rewards"
	"github.com/milkyway-labs/milkyway/v3/x/services"
)

var AppModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	vesting.AppModuleBasic{},
	bank.AppModuleBasic{},
	capability.AppModuleBasic{},
	crisis.AppModuleBasic{},
	gov.AppModuleBasic{},
	mint.AppModuleBasic{},
	slashing.AppModuleBasic{},
	distr.AppModuleBasic{},
	no_valupdates_staking.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	feegrantmodule.AppModuleBasic{},
	authzmodule.AppModuleBasic{},
	ibc.AppModuleBasic{},
	ibctm.AppModuleBasic{},
	sdkparams.AppModuleBasic{},
	consensus.AppModuleBasic{},
	wasm.AppModuleBasic{},

	// Skip modules
	feemarket.AppModuleBasic{},
	oracle.AppModuleBasic{},
	marketmap.AppModuleBasic{},

	// IBC Modules
	ibcfee.AppModuleBasic{},
	transfer.AppModuleBasic{},
	pfmrouter.AppModuleBasic{},
	ratelimit.AppModuleBasic{},
	icsprovider.AppModuleBasic{},

	// MilkyWay modules
	services.AppModuleBasic{},
	operators.AppModuleBasic{},
	pools.AppModuleBasic{},
	restaking.AppModuleBasic{},
	assets.AppModuleBasic{},
	rewards.AppModuleBasic{},
	liquidvesting.AppModuleBasic{},
)
