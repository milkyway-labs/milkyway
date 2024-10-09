package utils_test

import (
	"encoding/json"
	"testing"

	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	nfttransfertypes "github.com/initia-labs/initia/x/ibc/nft-transfer/types"
	"github.com/stretchr/testify/require"

	utils "github.com/milkyway-labs/milkyway/utils"
)

func Test_isIcs20Packet(t *testing.T) {
	transferMsg := transfertypes.NewFungibleTokenPacketData("denom", "1000000", "0x1", "0x2", "memo")
	bz, err := json.Marshal(transferMsg)
	require.NoError(t, err)

	ok, _transferMsg := utils.IsIcs20Packet(bz)
	require.True(t, ok)
	require.Equal(t, transferMsg, _transferMsg)

	nftTransferMsg := nfttransfertypes.NewNonFungibleTokenPacketData("class_id", "uri", "data", []string{"1", "2", "3"}, []string{"uri1", "uri2", "uri3"}, []string{"data1", "data2", "data3"}, "sender", "receiver", "memo")
	bz, err = json.Marshal(nftTransferMsg)
	require.NoError(t, err)

	ok, _ = utils.IsIcs20Packet(bz)
	require.False(t, ok)
}
