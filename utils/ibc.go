package utils

import (
	"encoding/json"
	"strings"

	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
)

// IsIcs20Packet checks if the data is a ICS20 packet.
func IsIcs20Packet(packetData []byte) (isIcs20 bool, ics20data transfertypes.FungibleTokenPacketData) {
	var data transfertypes.FungibleTokenPacketData
	decoder := json.NewDecoder(strings.NewReader(string(packetData)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&data); err != nil {
		return false, data
	}
	return true, data
}
