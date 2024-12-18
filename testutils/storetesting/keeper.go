package storetesting

import (
	"slices"
	"testing"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	milkyway "github.com/milkyway-labs/milkyway/v6/app"
)

type BaseKeeperTestData struct {
	Keys    map[string]*storetypes.KVStoreKey
	Context sdk.Context

	Cdc         codec.Codec
	LegacyAmino *codec.LegacyAmino

	AuthorityAddress string

	AccountKeeper      authkeeper.AccountKeeper
	BankKeeper         bankkeeper.BaseKeeper
	DistributionKeeper distrkeeper.Keeper
}

// NewBaseKeeperTestData returns a new BaseKeeperTestData
func NewBaseKeeperTestData(t *testing.T, keys []string) BaseKeeperTestData {
	t.Helper()

	// Set the Cosmos SDK configuration to use another Bech32 prefix
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("cosmos", "cosmospub")
	config.SetBech32PrefixForValidator("cosmosvaloper", "cosmosvaloperpub")
	config.SetBech32PrefixForConsensusNode("cosmosvalcons", "cosmosvalconspub")

	var data BaseKeeperTestData

	// Define store keys
	keys = append(keys, []string{authtypes.StoreKey, banktypes.StoreKey, distrtypes.StoreKey}...)
	slices.Sort(keys)
	keys = slices.Compact(keys)
	data.Keys = storetypes.NewKVStoreKeys(keys...)

	// Setup the context
	data.Context = BuildContext(data.Keys, nil, nil)

	// Setup the codecs
	data.Cdc, data.LegacyAmino = milkyway.MakeCodecs()

	// Authority address
	data.AuthorityAddress = authtypes.NewModuleAddress(govtypes.ModuleName).String()

	// Build keepers
	data.AccountKeeper = authkeeper.NewAccountKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.Keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		milkyway.MaccPerms,
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		data.AuthorityAddress,
	)
	data.BankKeeper = bankkeeper.NewBaseKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.Keys[banktypes.StoreKey]),
		data.AccountKeeper,
		milkyway.BlockedModuleAccountAddrs(milkyway.ModuleAccountAddrs()),
		data.AuthorityAddress,
		log.NewNopLogger(),
	)
	data.DistributionKeeper = distrkeeper.NewKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.Keys[distrtypes.StoreKey]),
		data.AccountKeeper,
		data.BankKeeper,
		nil,
		authtypes.FeeCollectorName,
		data.AuthorityAddress,
	)

	// Init the module's genesis state as the default ones
	data.AccountKeeper.InitGenesis(data.Context, *authtypes.DefaultGenesisState())
	data.BankKeeper.InitGenesis(data.Context, banktypes.DefaultGenesisState())
	data.DistributionKeeper.InitGenesis(data.Context, *distrtypes.DefaultGenesisState())

	return data
}
