package ibc_hooks_test

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	db "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibchooks "github.com/initia-labs/initia/x/ibc-hooks"
	ibchookskeeper "github.com/initia-labs/initia/x/ibc-hooks/keeper"
	ibchookstypes "github.com/initia-labs/initia/x/ibc-hooks/types"
	"github.com/stretchr/testify/suite"

	"github.com/milkyway-labs/milkyway/app"
	liquidvestinghooks "github.com/milkyway-labs/milkyway/x/liquidvesting/ibc-hooks"
	liquidvestingkeeper "github.com/milkyway-labs/milkyway/x/liquidvesting/keeper"
	liquidvestingtypes "github.com/milkyway-labs/milkyway/x/liquidvesting/types"
)

func TestHooksTestSuite(t *testing.T) {
	suite.Run(t, new(HooksTestSuite))
}

type HooksTestSuite struct {
	suite.Suite

	cdc            codec.Codec
	legacyAminoCdc *codec.LegacyAmino
	ctx            sdk.Context

	ak                 authkeeper.AccountKeeper
	bk                 bankkeeper.Keeper
	IBCHooksKeeper     *ibchookskeeper.Keeper
	IBCHooksMiddleware ibchooks.IBCMiddleware
	k                  *liquidvestingkeeper.Keeper
}

func (suite *HooksTestSuite) SetupTest() {
	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		distributiontypes.StoreKey, ibchookstypes.StoreKey,
		liquidvestingtypes.StoreKey,
	)

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

	ac := authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())

	accountKeeper := authkeeper.NewAccountKeeper(
		suite.cdc,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]), // target store
		authtypes.ProtoBaseAccount,                          // prototype
		map[string][]string{},
		ac,
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	bankKeeper := bankkeeper.NewBaseKeeper(
		suite.cdc,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		accountKeeper,
		map[string]bool{},
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		logger,
	)
	if err := bankKeeper.SetParams(suite.ctx, banktypes.DefaultParams()); err != nil {
		panic(err)
	}

	ibcHooksKeeper := ibchookskeeper.NewKeeper(
		suite.cdc,
		runtime.NewKVStoreService(keys[ibchookstypes.StoreKey]),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		ac,
	)
	if err := ibcHooksKeeper.Params.Set(suite.ctx, ibchookstypes.DefaultParams()); err != nil {
		panic(err)
	}

	// ibc middleware setup
	mockIBCMiddleware := mockIBCMiddleware{}
	liquidVestingKeeper := liquidvestingkeeper.NewKeeper(
		suite.cdc,
		runtime.NewKVStoreService(keys[liquidvestingtypes.StoreKey]),
		bankKeeper,
		nil,
		nil,
		nil,
		nil,
		authtypes.NewModuleAddress(liquidvestingtypes.ModuleName).String(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	liquidVestingHooks := liquidvestinghooks.NewHooks(liquidVestingKeeper)
	middleware := ibchooks.NewICS4Middleware(mockIBCMiddleware, liquidVestingHooks)
	ibcHookMiddleware := ibchooks.NewIBCMiddleware(mockIBCMiddleware, middleware, ibcHooksKeeper)

	suite.ak = accountKeeper
	suite.bk = bankKeeper
	suite.IBCHooksKeeper = ibcHooksKeeper
	suite.IBCHooksMiddleware = ibcHookMiddleware
	suite.k = liquidVestingKeeper
}

func (suite *HooksTestSuite) TearDownTest() {}

func (suite *HooksTestSuite) userAddress(name string) sdk.AccAddress {
	seed := []byte(name)
	key := ed25519.GenPrivKeyFromSecret(seed)
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return addr
}

// do nothing ibc middleware
var (
	_ porttypes.IBCModule   = mockIBCMiddleware{}
	_ porttypes.ICS4Wrapper = mockIBCMiddleware{}
)

type mockIBCMiddleware struct{}

// GetAppVersion implements types.ICS4Wrapper.
func (m mockIBCMiddleware) GetAppVersion(ctx sdk.Context, portID string, channelID string) (string, bool) {
	return "", false
}

// SendPacket implements types.ICS4Wrapper.
func (m mockIBCMiddleware) SendPacket(ctx sdk.Context, chanCap *capabilitytypes.Capability, sourcePort string, sourceChannel string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64, data []byte) (sequence uint64, err error) {
	return 0, nil
}

// WriteAcknowledgement implements types.ICS4Wrapper.
func (m mockIBCMiddleware) WriteAcknowledgement(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet ibcexported.PacketI, ack ibcexported.Acknowledgement) error {
	return nil
}

// OnAcknowledgementPacket implements types.IBCModule.
func (m mockIBCMiddleware) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress) error {
	return nil
}

// OnChanCloseConfirm implements types.IBCModule.
func (m mockIBCMiddleware) OnChanCloseConfirm(ctx sdk.Context, portID string, channelID string) error {
	return nil
}

// OnChanCloseInit implements types.IBCModule.
func (m mockIBCMiddleware) OnChanCloseInit(ctx sdk.Context, portID string, channelID string) error {
	return nil
}

// OnChanOpenAck implements types.IBCModule.
func (m mockIBCMiddleware) OnChanOpenAck(ctx sdk.Context, portID string, channelID string, counterpartyChannelID string, counterpartyVersion string) error {
	return nil
}

// OnChanOpenConfirm implements types.IBCModule.
func (m mockIBCMiddleware) OnChanOpenConfirm(ctx sdk.Context, portID string, channelID string) error {
	return nil
}

// OnChanOpenInit implements types.IBCModule.
func (m mockIBCMiddleware) OnChanOpenInit(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID string, channelID string, channelCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, version string) (string, error) {
	return "", nil
}

// OnChanOpenTry implements types.IBCModule.
func (m mockIBCMiddleware) OnChanOpenTry(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID string, channelID string, channelCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, counterpartyVersion string) (version string, err error) {
	return "", nil
}

// OnRecvPacket implements types.IBCModule.
func (m mockIBCMiddleware) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) ibcexported.Acknowledgement {
	return channeltypes.NewResultAcknowledgement([]byte{byte(1)})
}

// OnTimeoutPacket implements types.IBCModule.
func (m mockIBCMiddleware) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) error {
	return nil
}
