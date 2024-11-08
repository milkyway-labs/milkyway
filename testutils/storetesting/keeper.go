package storetesting

import (
	"testing"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/milkyway-labs/milkyway/app"
	bankkeeper "github.com/milkyway-labs/milkyway/x/bank/keeper"
)

type BaseKeeperTestData struct {
	Keys    map[string]*storetypes.KVStoreKey
	Context sdk.Context

	Cdc         codec.Codec
	LegacyAmino *codec.LegacyAmino

	AuthorityAddress string

	AccountKeeper authkeeper.AccountKeeper
	BankKeeper    bankkeeper.Keeper
}

// NewBaseKeeperTestData returns a new BaseKeeperTestData
func NewBaseKeeperTestData(t *testing.T, keys []string) BaseKeeperTestData {
	t.Helper()

	var data BaseKeeperTestData

	// Define store keys
	data.Keys = storetypes.NewKVStoreKeys(keys...)

	// Setup the context
	data.Context = BuildContext(data.Keys, nil, nil)

	// Setup the codecs
	data.Cdc, data.LegacyAmino = app.MakeCodecs()

	// Authority address
	data.AuthorityAddress = authtypes.NewModuleAddress(govtypes.ModuleName).String()

	// Build keepers
	data.AccountKeeper = authkeeper.NewAccountKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.Keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		app.GetMaccPerms(),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		data.AuthorityAddress,
	)
	data.BankKeeper = bankkeeper.NewKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.Keys[banktypes.StoreKey]),
		data.AccountKeeper,
		nil,
		data.AuthorityAddress,
		log.NewNopLogger(),
	)

	return data
}
