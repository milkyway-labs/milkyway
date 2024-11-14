package utils_test

import (
	"encoding/json"
	"testing"

	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	nfttransfertypes "github.com/initia-labs/initia/x/ibc/nft-transfer/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/utils"
)

func Test_DeserializeFungibleTokenPacketData(t *testing.T) {
	expected := transfertypes.NewFungibleTokenPacketData("denom", "1000000", "0x1", "0x2", "memo")
	bz, err := json.Marshal(expected)
	require.NoError(t, err)

	msg, ok := utils.DeserializeFungibleTokenPacketData(bz)
	require.True(t, ok)
	require.Equal(t, expected, msg)

	nftTransferMsg := nfttransfertypes.NewNonFungibleTokenPacketData("class_id", "uri", "data", []string{"1", "2", "3"}, []string{"uri1", "uri2", "uri3"}, []string{"data1", "data2", "data3"}, "sender", "receiver", "memo")
	bz, err = json.Marshal(nftTransferMsg)
	require.NoError(t, err)

	_, ok = utils.DeserializeFungibleTokenPacketData(bz)
	require.False(t, ok)
}
