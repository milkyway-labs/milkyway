package testutils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
)

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
