package testutils

import (
	"testing"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"go.uber.org/mock/gomock"

	milkyway "github.com/milkyway-labs/milkyway/app"
	"github.com/milkyway-labs/milkyway/testutils/storetesting"
	bankkeeper "github.com/milkyway-labs/milkyway/x/bank/keeper"
	poolskeeper "github.com/milkyway-labs/milkyway/x/pools/keeper"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	"github.com/milkyway-labs/milkyway/x/services/keeper"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

type KeeperTestData struct {
	storetesting.BaseKeeperTestData

	StoreKey storetypes.StoreKey
	MockCtrl *gomock.Controller

	PoolKeeper  *MockCommunityPoolKeeper
	PoolsKeeper *poolskeeper.Keeper

	Keeper *keeper.Keeper
	Hooks  *MockHooks
}

func NewKeeperTestData(t *testing.T) KeeperTestData {
	var data = KeeperTestData{
		BaseKeeperTestData: storetesting.NewBaseKeeperTestData(t, []string{
			authtypes.StoreKey,
			banktypes.StoreKey,
			servicestypes.StoreKey,
			poolstypes.StoreKey,
		}),
	}

	// Define store keys
	data.StoreKey = data.Keys[servicestypes.StoreKey]

	// Setup the codecs
	encodingConfig := milkyway.MakeEncodingConfig()
	data.Cdc, data.LegacyAmino = encodingConfig.Marshaler, encodingConfig.Amino

	// Mocks initializations
	data.MockCtrl = gomock.NewController(t)
	data.PoolKeeper = NewMockCommunityPoolKeeper(data.MockCtrl)

	// Authority address
	authorityAddr := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	// Build keepers
	data.AccountKeeper = authkeeper.NewAccountKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.Keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		milkyway.MaccPerms,
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		authorityAddr,
	)
	data.BankKeeper = bankkeeper.NewKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.Keys[banktypes.StoreKey]),
		data.AccountKeeper,
		nil,
		authorityAddr,
		log.NewNopLogger(),
	)
	data.PoolsKeeper = poolskeeper.NewKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.Keys[poolstypes.StoreKey]),
		data.AccountKeeper,
	)
	data.Keeper = keeper.NewKeeper(
		data.Cdc,
		runtime.NewKVStoreService(data.Keys[servicestypes.StoreKey]),
		data.AccountKeeper,
		data.DistributionKeeper,
		authorityAddr,
	)

	// Set hooks
	data.Hooks = NewMockHooks()
	data.Keeper = data.Keeper.SetHooks(data.Hooks)

	return data
}