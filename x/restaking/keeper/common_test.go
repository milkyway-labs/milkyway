package keeper_test

import (
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

	"github.com/milkyway-labs/milkyway/app"
	appkeepers "github.com/milkyway-labs/milkyway/app/keepers"
	bankkeeper "github.com/milkyway-labs/milkyway/x/bank/keeper"
	operatorskeeper "github.com/milkyway-labs/milkyway/x/operators/keeper"
	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolskeeper "github.com/milkyway-labs/milkyway/x/pools/keeper"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	"github.com/milkyway-labs/milkyway/x/restaking/keeper"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
	serviceskeeper "github.com/milkyway-labs/milkyway/x/services/keeper"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"

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

	ak authkeeper.AccountKeeper
	bk bankkeeper.Keeper
	pk *poolskeeper.Keeper
	ok *operatorskeeper.Keeper
	sk *serviceskeeper.Keeper
	k  *keeper.Keeper
}

func (suite *KeeperTestSuite) SetupTest() {
	// Define store keys
	keys := storetypes.NewKVStoreKeys(
		types.StoreKey,
		authtypes.StoreKey, banktypes.StoreKey, poolstypes.StoreKey, operatorstypes.StoreKey, servicestypes.StoreKey,
	)
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
	communityPoolKeeper := appkeepers.NewCommunityPoolKeeper(
		suite.bk,
		authtypes.FeeCollectorName,
	)
	suite.pk = poolskeeper.NewKeeper(
		suite.cdc,
		keys[poolstypes.StoreKey],
		suite.ak,
	)
	suite.ok = operatorskeeper.NewKeeper(
		suite.cdc,
		keys[operatorstypes.StoreKey],
		suite.ak,
		communityPoolKeeper,
		authorityAddr,
	)
	suite.sk = serviceskeeper.NewKeeper(
		suite.cdc,
		keys[servicestypes.StoreKey],
		suite.ak,
		communityPoolKeeper,
		authorityAddr,
	)
	suite.k = keeper.NewKeeper(
		suite.cdc,
		suite.storeKey,
		suite.ak,
		suite.bk,
		suite.pk,
		suite.ok,
		suite.sk,
		authorityAddr,
	).SetHooks(newMockHooks())
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

var _ types.RestakingHooks = &mockHooks{}

type mockHooks struct {
	CalledMap map[string]bool
}

func newMockHooks() *mockHooks {
	return &mockHooks{CalledMap: make(map[string]bool)}
}

func (m mockHooks) BeforePoolDelegationCreated(sdk.Context, uint32, string) error {
	m.CalledMap["BeforePoolDelegationCreated"] = true
	return nil
}

func (m mockHooks) BeforePoolDelegationSharesModified(sdk.Context, uint32, string) error {
	m.CalledMap["BeforePoolDelegationSharesModified"] = true
	return nil
}

func (m mockHooks) AfterPoolDelegationModified(sdk.Context, uint32, string) error {
	m.CalledMap["AfterPoolDelegationModified"] = true
	return nil
}

func (m mockHooks) BeforeOperatorDelegationCreated(sdk.Context, uint32, string) error {
	m.CalledMap["BeforeOperatorDelegationCreated"] = true
	return nil
}

func (m mockHooks) BeforeOperatorDelegationSharesModified(sdk.Context, uint32, string) error {
	m.CalledMap["BeforeOperatorDelegationSharesModified"] = true
	return nil
}

func (m mockHooks) AfterOperatorDelegationModified(sdk.Context, uint32, string) error {
	m.CalledMap["AfterOperatorDelegationModified"] = true
	return nil
}

func (m mockHooks) BeforeServiceDelegationCreated(sdk.Context, uint32, string) error {
	m.CalledMap["BeforeServiceDelegationCreated"] = true
	return nil
}

func (m mockHooks) BeforeServiceDelegationSharesModified(sdk.Context, uint32, string) error {
	m.CalledMap["BeforeServiceDelegationSharesModified"] = true
	return nil
}

func (m mockHooks) AfterServiceDelegationModified(sdk.Context, uint32, string) error {
	m.CalledMap["AfterServiceDelegationModified"] = true
	return nil
}
