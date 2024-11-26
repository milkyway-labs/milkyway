package utils_test

import (
	"encoding/json"
	"testing"

	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/v2/utils"
)

func Test_DeserializeFungibleTokenPacketData(t *testing.T) {
	expected := transfertypes.NewFungibleTokenPacketData("denom", "1000000", "0x1", "0x2", "memo")
	bz, err := json.Marshal(expected)
	require.NoError(t, err)

	msg, ok := utils.DeserializeFungibleTokenPacketData(bz)
	require.True(t, ok)
	require.Equal(t, expected, msg)
}
