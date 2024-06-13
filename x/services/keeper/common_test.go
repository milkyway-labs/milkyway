package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	"github.com/cosmos/cosmos-sdk/runtime"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"go.uber.org/mock/gomock"

	"github.com/milkyway-labs/milkyway/app"
	"github.com/milkyway-labs/milkyway/app/keepers"
	bankkeeper "github.com/milkyway-labs/milkyway/x/bank/keeper"
	"github.com/milkyway-labs/milkyway/x/services/keeper"
	"github.com/milkyway-labs/milkyway/x/services/testutil"
	"github.com/milkyway-labs/milkyway/x/services/types"

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

	storeKey storetypes.StoreKey

	ak    authkeeper.AccountKeeper
	bk    bankkeeper.Keeper
	k     *keeper.Keeper
	hooks *mockHooks

	ctrl       *gomock.Controller
	poolKeeper *testutil.MockCommunityPoolKeeper
}

func (suite *KeeperTestSuite) SetupTest() {
	// Define store keys
	keys := storetypes.NewKVStoreKeys(types.StoreKey, authtypes.StoreKey, banktypes.StoreKey)
	suite.storeKey = keys[types.StoreKey]

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
	suite.cdc, suite.legacyAminoCdc = app.MakeCodecs()

	// Mocks initializations
	suite.ctrl = gomock.NewController(suite.T())
	suite.poolKeeper = testutil.NewMockCommunityPoolKeeper(suite.ctrl)

	// Authority address
	authorityAddr := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	// Build keepers
	suite.ak = authkeeper.NewAccountKeeper(
		suite.cdc,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		app.GetMaccPerms(),
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
	suite.k = keeper.NewKeeper(
		suite.cdc,
		suite.storeKey,
		suite.ak,
		keepers.NewCommunityPoolKeeper(suite.bk, authtypes.FeeCollectorName),
		authorityAddr,
	)

	// Set hooks
	suite.hooks = newMockHooks()
	suite.k = suite.k.SetHooks(suite.hooks)
}

func (suite *KeeperTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

// --------------------------------------------------------------------------------------------------------------------

// fundAccount adds the given amount of coins to the account with the given address
func (suite *KeeperTestSuite) fundAccount(ctx sdk.Context, address string, amount sdk.Coins) {
	// Mint the coins
	moduleAcc := suite.ak.GetModuleAccount(ctx, authtypes.Minter)

	err := suite.bk.MintCoins(ctx, moduleAcc.GetName(), amount)
	suite.Require().NoError(err)

	// Get the amount to the user
	userAddress, err := sdk.AccAddressFromBech32(address)
	suite.Require().NoError(err)
	err = suite.bk.SendCoinsFromModuleToAccount(ctx, moduleAcc.GetName(), userAddress, amount)
	suite.Require().NoError(err)
}

// --------------------------------------------------------------------------------------------------------------------

var _ types.ServicesHooks = &mockHooks{}

type mockHooks struct {
	CalledMap map[string]bool
}

func newMockHooks() *mockHooks {
	return &mockHooks{CalledMap: make(map[string]bool)}
}

func (m mockHooks) AfterServiceCreated(_ context.Context, _ uint32) {
	m.CalledMap["AfterServiceCreated"] = true
}

func (m mockHooks) AfterServiceActivated(_ context.Context, _ uint32) {
	m.CalledMap["AfterServiceActivated"] = true
}

func (m mockHooks) AfterServiceDeactivated(_ context.Context, _ uint32) {
	m.CalledMap["AfterServiceDeactivated"] = true
}
