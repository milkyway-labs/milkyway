package liquidvesting

import (
	"encoding/json"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"

	"github.com/milkyway-labs/milkyway/v6/utils"
	"github.com/milkyway-labs/milkyway/v6/x/liquidvesting/keeper"
	"github.com/milkyway-labs/milkyway/v6/x/liquidvesting/types"
)

type IBCModule struct {
	keeper *keeper.Keeper
	app    porttypes.IBCModule
}

// NewIBCMiddleware creates a new IBCModule given the keeper
func NewIBCMiddleware(app porttypes.IBCModule, k *keeper.Keeper) IBCModule {
	return IBCModule{
		keeper: k,
		app:    app,
	}
}

// OnChanOpenInit implements the IBCModule interface
func (im IBCModule) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	channelCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) (string, error) {
	return im.app.OnChanOpenInit(
		ctx,
		order,
		connectionHops,
		portID,
		channelID,
		channelCap,
		counterparty,
		version,
	)
}

// OnChanOpenTry implements the IBCModule interface.
func (im IBCModule) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	counterpartyVersion string,
) (string, error) {
	return im.app.OnChanOpenTry(
		ctx,
		order,
		connectionHops,
		portID,
		channelID,
		chanCap,
		counterparty,
		counterpartyVersion,
	)
}

// OnChanOpenAck implements the IBCModule interface
func (im IBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	counterpartyChannelID string,
	counterpartyVersion string,
) error {
	// core/04-channel/types contains a helper function to split middleware and underlying app version
	// _, _ := channeltypes.SplitChannelVersion(counterpartyVersion)
	// doCustomLogic()
	// call the underlying applications OnChanOpenTry callback
	return im.app.OnChanOpenAck(ctx, portID, channelID, counterpartyChannelID, counterpartyVersion)
}

// OnChanOpenConfirm implements the IBCModule interface
func (im IBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// doCustomLogic()
	return im.app.OnChanOpenConfirm(ctx, portID, channelID)
}

// OnChanCloseInit implements the IBCModule interface
func (im IBCModule) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// doCustomLogic()
	return im.app.OnChanCloseInit(ctx, portID, channelID)
}

// OnChanCloseConfirm implements the IBCModule interface
func (im IBCModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// doCustomLogic()
	return im.app.OnChanCloseConfirm(ctx, portID, channelID)
}

// OnRecvPacket implements the IBCModule interface. A successful acknowledgement
// is returned if the packet data is successfully decoded and the receive application
// logic returns without error.
func (im IBCModule) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) ibcexported.Acknowledgement {
	// Decode the packet data
	data, ok := utils.DeserializeFungibleTokenPacketData(packet.GetData())
	if !ok {
		return channeltypes.NewErrorAcknowledgement(errors.New("invalid packet data"))
	}

	// Check if the packet contains a MsgDepositInsurance
	msgDepositInsurance, found, err := im.containsMsgDepositInsurance(data)
	if !found {
		// If the packet does not contain a MsgDepositInsurance, just pass it to the application
		return im.app.OnRecvPacket(ctx, packet, relayer)
	}

	if err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}

	// Handle the MsgDepositInsurance
	if err := im.keeper.OnRecvPacket(ctx, packet, data, msgDepositInsurance); err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}

	return im.app.OnRecvPacket(ctx, packet, relayer)
}

// containsMsgDepositInsurance checks if the packet data contains a MsgDepositInsurance.
// If the packet data contains a MsgDepositInsurance, it returns the message and true.
// Otherwise, it returns an empty MsgDepositInsurance and false.
func (im IBCModule) containsMsgDepositInsurance(data ibctransfertypes.FungibleTokenPacketData) (types.MsgDepositInsurance, bool, error) {
	objFound, object := utils.JSONStringHasKey(data.GetMemo(), types.ModuleName)
	if !objFound {
		return types.MsgDepositInsurance{}, false, nil
	}

	// Parse the message from the memo
	bytes, err := json.Marshal(object[types.ModuleName])
	if err != nil {
		return types.MsgDepositInsurance{}, true, err
	}

	var depositMsg types.MsgDepositInsurance
	if err = json.Unmarshal(bytes, &depositMsg); err != nil {
		return types.MsgDepositInsurance{}, true, err
	}

	return depositMsg, true, nil
}

// OnAcknowledgementPacket implements the IBCModule interface
func (im IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	return im.app.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
}

// OnTimeoutPacket implements the IBCModule interface
func (im IBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	return im.app.OnTimeoutPacket(ctx, packet, relayer)
}

// GetAppVersion implements the IBCModule interface
func (im IBCModule) GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool) {
	return ibctransfertypes.Version, true
}
