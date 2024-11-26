package keeper_test

import (
	"testing"

	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	"github.com/cosmos/cosmos-sdk/runtime"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	milkyway "github.com/milkyway-labs/milkyway/v2/app"
	bankkeeper "github.com/milkyway-labs/milkyway/v2/x/bank/keeper"
	"github.com/milkyway-labs/milkyway/v2/x/pools/keeper"
	"github.com/milkyway-labs/milkyway/v2/x/pools/types"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	storetypes "cosmossdk.io/store/types"
	db "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

type KeeperTestSuite struct {
	suite.Suite

	cdc            codec.Codec
	legacyAminoCdc *codec.LegacyAmino
	ctx            sdk.Context

	storeService corestoretypes.KVStoreService

	ak authkeeper.AccountKeeper
	bk bankkeeper.Keeper
	k  *keeper.Keeper
}

func (suite *KeeperTestSuite) SetupTest() {
	// Define store keys
	keys := storetypes.NewKVStoreKeys(types.StoreKey, authtypes.StoreKey, banktypes.StoreKey)

	// Create logger
	logger := log.NewNopLogger()

	// Create an in-memory db
	memDB := db.NewMemDB()
	ms := store.NewCommitMultiStore(memDB, logger, metrics.NewNoOpMetrics())
	for _, key := range keys {
		ms.MountStoreWithDB(key, storetypes.StoreTypeIAVL, memDB)
	}

	if err := ms.LoadLatestVersion(); err != nil {
		panic(err)
	}

	suite.ctx = sdk.NewContext(ms, tmproto.Header{ChainID: "test-chain"}, false, log.NewNopLogger())
	suite.cdc, suite.legacyAminoCdc = milkyway.MakeCodecs()

	// Authority address
	authorityAddr := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	// Build keepers
	suite.ak = authkeeper.NewAccountKeeper(
		suite.cdc,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		milkyway.MaccPerms,
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		authorityAddr,
	)
	suite.bk = bankkeeper.NewKeeper(
		suite.cdc,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		suite.ak,
		nil,
		authorityAddr,
		logger,
	)

	suite.storeService = runtime.NewKVStoreService(keys[types.StoreKey])
	suite.k = keeper.NewKeeper(
		suite.cdc,
		suite.storeService,
		suite.ak,
	)
}
