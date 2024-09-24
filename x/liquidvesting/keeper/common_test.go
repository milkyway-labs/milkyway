package keeper_test

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibchooks "github.com/initia-labs/initia/x/ibc-hooks"
	ibchookskeeper "github.com/initia-labs/initia/x/ibc-hooks/keeper"
	ibchookstypes "github.com/initia-labs/initia/x/ibc-hooks/types"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	db "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/suite"

	"github.com/milkyway-labs/milkyway/app"
	appkeepers "github.com/milkyway-labs/milkyway/app/keepers"
	bankkeeper "github.com/milkyway-labs/milkyway/x/bank/keeper"
	"github.com/milkyway-labs/milkyway/x/liquidvesting/keeper"
	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
	operatorskeeper "github.com/milkyway-labs/milkyway/x/operators/keeper"
	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolskeeper "github.com/milkyway-labs/milkyway/x/pools/keeper"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	restakingkeeper "github.com/milkyway-labs/milkyway/x/restaking/keeper"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	serviceskeeper "github.com/milkyway-labs/milkyway/x/services/keeper"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

type KeeperTestSuite struct {
	suite.Suite

	cdc            codec.Codec
	legacyAminoCdc *codec.LegacyAmino
	ctx            sdk.Context

	ak   authkeeper.AccountKeeper
	bk   *bankkeeper.Keeper
	ibck *ibchookskeeper.Keeper
	ibcm ibchooks.IBCMiddleware
	ok   *operatorskeeper.Keeper
	pk   *poolskeeper.Keeper
	sk   *serviceskeeper.Keeper
	rk   *restakingkeeper.Keeper

	k *keeper.Keeper
}

func (suite *KeeperTestSuite) SetupTest() {
	// Define store keys
	keys := storetypes.NewKVStoreKeys(
		types.StoreKey, authtypes.StoreKey, banktypes.StoreKey,
		operatorstypes.StoreKey, poolstypes.StoreKey, servicestypes.StoreKey,
		restakingtypes.StoreKey, stakingtypes.StoreKey,
		distributiontypes.StoreKey, ibchookstypes.StoreKey,
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
	bk := bankkeeper.NewKeeper(
		suite.cdc,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		suite.ak,
		nil,
		authorityAddr,
		logger,
	)
	suite.bk = &bk
	suite.ibck = ibchookskeeper.NewKeeper(
		suite.cdc,
		runtime.NewKVStoreService(keys[ibchookstypes.StoreKey]),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		ac,
	)
	if err := suite.ibck.Params.Set(suite.ctx, ibchookstypes.DefaultParams()); err != nil {
		panic(err)
	}

	suite.pk = poolskeeper.NewKeeper(
		suite.cdc,
		keys[poolstypes.StoreKey],
		suite.ak,
	)
	communityPoolKeeper := appkeepers.NewCommunityPoolKeeper(
		suite.bk,
		authtypes.FeeCollectorName,
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
	suite.rk = restakingkeeper.NewKeeper(
		suite.cdc,
		keys[restakingtypes.StoreKey],
		suite.ak,
		suite.bk,
		suite.pk,
		suite.ok,
		suite.sk,
		authorityAddr,
	)
	suite.k = keeper.NewKeeper(
		suite.cdc,
		keys[types.StoreKey],
		runtime.NewKVStoreService(keys[types.StoreKey]),
		suite.bk,
		suite.ok,
		suite.pk,
		suite.sk,
		suite.rk,
		authtypes.NewModuleAddress(types.ModuleName).String(),
		authorityAddr,
	)
	// Set bank hooks
	suite.bk.AppendSendRestriction(suite.k.SendRestrictionFn)
	// Set ibc hooks

	mockIBCMiddleware := mockIBCMiddleware{}
	middleware := ibchooks.NewICS4Middleware(mockIBCMiddleware, suite.k.IBCHooks())
	suite.ibcm = ibchooks.NewIBCMiddleware(mockIBCMiddleware, middleware, suite.ibck)

	account := suite.ak.GetModuleAccount(suite.ctx, types.ModuleName)
	suite.Assert().NotNil(account)
	suite.ak.SetModuleAccount(suite.ctx, account)
}

// --------------------------------------------------------------------------------------------------------------------

// fundAccount add the given amount of coins to the account's balance
func (suite *KeeperTestSuite) fundAccount(ctx sdk.Context, address string, amount sdk.Coins) {
	// Mint the tokens in the insurance fund.
	suite.Assert().NoError(suite.bk.MintCoins(ctx, types.ModuleName, amount))

	suite.Assert().NoError(suite.bk.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, sdk.MustAccAddressFromBech32(address), amount))
}

// mintVestedRepresentation mints the vested representation of the provided amount to
// the user balance
func (suite *KeeperTestSuite) mintVestedRepresentation(address string, amount sdk.Coins) {
	accAddress, err := sdk.AccAddressFromBech32(address)
	suite.Assert().NoError(err)

	_, err = suite.k.MintVestedRepresentation(
		suite.ctx, accAddress, amount,
	)
	suite.Assert().NoError(err)
}

// fundAccountInsuranceFund add the given amount of coins to the account's insurance fund
func (suite *KeeperTestSuite) fundAccountInsuranceFund(ctx sdk.Context, address string, amount sdk.Coins) {
	// Mint the tokens in the insurance fund.
	suite.Assert().NoError(suite.bk.MintCoins(suite.ctx, types.ModuleName, amount))

	// Assign those tokens to the user insurance fund
	userAddress, err := sdk.AccAddressFromBech32(address)
	suite.Assert().NoError(err)
	suite.Assert().NoError(suite.k.AddToUserInsuranceFund(
		ctx,
		userAddress,
		amount,
	))
}

// createPool creates a test pool with the given id and denom
func (suite *KeeperTestSuite) createPool(id uint32, denom string) {
	err := suite.pk.SavePool(suite.ctx, poolstypes.Pool{
		ID:              id,
		Denom:           denom,
		Address:         poolstypes.GetPoolAddress(id).String(),
		Tokens:          sdkmath.NewInt(0),
		DelegatorShares: sdkmath.LegacyNewDec(0),
	})
	suite.Assert().NoError(err)
}

// createService creates a test service with the provided id
func (suite *KeeperTestSuite) createService(id uint32) {
	err := suite.sk.CreateService(suite.ctx, servicestypes.NewService(
		id,
		servicestypes.SERVICE_STATUS_ACTIVE,
		fmt.Sprintf("test %d", id),
		fmt.Sprintf("test service %d", id),
		"",
		"",
		fmt.Sprintf("service-%d-admin", id),
	))
	suite.Assert().NoError(err)
}

func (suite *KeeperTestSuite) createOperator(id uint32) {
	suite.Assert().NoError(suite.ok.RegisterOperator(suite.ctx, operatorstypes.NewOperator(
		id,
		operatorstypes.OPERATOR_STATUS_ACTIVE,
		fmt.Sprintf("operator-%d", id),
		"",
		"",
		fmt.Sprintf("operator-%d-admin", id))))
}

// ---------------------------------------------
// ------------ IBC Mocks -----------------------
// ---------------------------------------------
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
