package keeper_test

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	"github.com/stretchr/testify/suite"

	"github.com/milkyway-labs/milkyway/v8/x/liquidvesting/keeper"
	"github.com/milkyway-labs/milkyway/v8/x/liquidvesting/testutils"
	"github.com/milkyway-labs/milkyway/v8/x/liquidvesting/types"
	operatorskeeper "github.com/milkyway-labs/milkyway/v8/x/operators/keeper"
	operatorstypes "github.com/milkyway-labs/milkyway/v8/x/operators/types"
	poolskeeper "github.com/milkyway-labs/milkyway/v8/x/pools/keeper"
	poolstypes "github.com/milkyway-labs/milkyway/v8/x/pools/types"
	restakingkeeper "github.com/milkyway-labs/milkyway/v8/x/restaking/keeper"
	serviceskeeper "github.com/milkyway-labs/milkyway/v8/x/services/keeper"
	servicestypes "github.com/milkyway-labs/milkyway/v8/x/services/types"
)

const (
	IBCDenom       = "ibc/D79E7D83AB399BFFF93433E54FAA480C191248FC556924A2A8351AE2638B3877"
	LockedIBCDenom = "locked/" + IBCDenom
)

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

type KeeperTestSuite struct {
	suite.Suite

	cdc codec.Codec
	ctx sdk.Context

	liquidVestingModuleAddress sdk.AccAddress

	ak authkeeper.AccountKeeper
	bk bankkeeper.BaseKeeper
	ok *operatorskeeper.Keeper
	pk *poolskeeper.Keeper
	sk *serviceskeeper.Keeper
	rk *restakingkeeper.Keeper

	k *keeper.Keeper

	ibcm porttypes.IBCModule
}

func (suite *KeeperTestSuite) SetupTest() {
	data := testutils.NewKeeperTestData(suite.T())
	suite.ctx = data.Context
	suite.cdc = data.Cdc

	suite.liquidVestingModuleAddress = authtypes.NewModuleAddress(types.ModuleName)

	suite.ak = data.AccountKeeper
	suite.bk = data.BankKeeper
	suite.ibcm = data.IBCMiddleware
	suite.pk = data.PoolsKeeper

	suite.ok = data.OperatorsKeeper
	suite.sk = data.ServicesKeeper
	suite.rk = data.RestakingKeeper
	suite.k = data.Keeper

	// Setup IBC
	suite.ibcm = data.IBCMiddleware
}

// --------------------------------------------------------------------------------------------------------------------

// fundAccount add the given amount of coins to the account's balance
func (suite *KeeperTestSuite) fundAccount(ctx sdk.Context, address string, amount sdk.Coins) {
	// Mint the tokens in the insurance fund.
	err := suite.bk.MintCoins(ctx, types.ModuleName, amount)
	suite.Assert().NoError(err)

	err = suite.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.MustAccAddressFromBech32(address), amount)
	suite.Assert().NoError(err)
}

// mintLockedRepresentation mints the locked representation of the provided amount to
// the user balance
func (suite *KeeperTestSuite) mintLockedRepresentation(ctx sdk.Context, address string, amount sdk.Coins) {
	accAddress, err := sdk.AccAddressFromBech32(address)
	suite.Assert().NoError(err)

	_, err = suite.k.MintLockedRepresentation(ctx, accAddress, amount)
	suite.Assert().NoError(err)
}

// fundAccountInsuranceFund add the given amount of coins to the account's insurance fund
func (suite *KeeperTestSuite) fundAccountInsuranceFund(ctx sdk.Context, address string, amount sdk.Coins) {
	// Mint the tokens in the insurance fund.
	err := suite.bk.MintCoins(ctx, types.ModuleName, amount)
	suite.Assert().NoError(err)

	// Assign those tokens to the user insurance fund
	err = suite.k.AddToUserInsuranceFund(ctx, address, amount)
	suite.Assert().NoError(err)
}

// createPool creates a test pool with the given id and denom
func (suite *KeeperTestSuite) createPool(ctx sdk.Context, id uint32, denom string) {
	err := suite.pk.SavePool(ctx, poolstypes.Pool{
		ID:              id,
		Denom:           denom,
		Address:         poolstypes.GetPoolAddress(id).String(),
		Tokens:          sdkmath.NewInt(0),
		DelegatorShares: sdkmath.LegacyNewDec(0),
	})
	suite.Assert().NoError(err)
}

// createService creates a test service with the provided id
func (suite *KeeperTestSuite) createService(ctx sdk.Context, id uint32) {
	err := suite.sk.CreateService(ctx, servicestypes.NewService(
		id,
		servicestypes.SERVICE_STATUS_ACTIVE,
		fmt.Sprintf("test %d", id),
		fmt.Sprintf("test service %d", id),
		"",
		"",
		fmt.Sprintf("service-%d-admin", id),
		false,
	))
	suite.Assert().NoError(err)
}

func (suite *KeeperTestSuite) createOperator(ctx sdk.Context, id uint32) {
	err := suite.ok.CreateOperator(ctx, operatorstypes.NewOperator(
		id,
		operatorstypes.OPERATOR_STATUS_ACTIVE,
		fmt.Sprintf("operator-%d", id),
		"",
		"",
		fmt.Sprintf("operator-%d-admin", id),
	))
	suite.Assert().NoError(err)
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
