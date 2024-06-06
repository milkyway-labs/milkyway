package keeper_test

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"go.uber.org/mock/gomock"

	"github.com/milkyway-labs/milkyway/app"
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
	k        *keeper.Keeper

	ctrl       *gomock.Controller
	poolKeeper *testutil.MockCommunityPoolKeeper
}

func (suite *KeeperTestSuite) SetupTest() {
	// Define store keys
	keys := storetypes.NewMemoryStoreKeys(types.StoreKey)
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
	suite.k = keeper.NewKeeper(
		suite.cdc,
		suite.storeKey,
		suite.poolKeeper,
		authtypes.NewModuleAddress("gov").String(),
	)
}

func (suite *KeeperTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}
