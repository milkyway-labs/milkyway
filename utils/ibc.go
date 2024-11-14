package utils

import (
	"encoding/json"
	"strings"

	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
)

// DeserializeFungibleTokenPacketData deserializes the packet data and returns the FungibleTokenPacketData
func DeserializeFungibleTokenPacketData(packetData []byte) (ics20data transfertypes.FungibleTokenPacketData, isIcs20 bool) {
	var data transfertypes.FungibleTokenPacketData
	decoder := json.NewDecoder(strings.NewReader(string(packetData)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&data); err != nil {
		return data, false
	}
	return data, true
}
