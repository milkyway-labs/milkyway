package testutil

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

// TestAddress returns an account address used in test codes derived from the
// given integer n, so that it's easier to write tests with multiple accounts.
func TestAddress(n uint64) sdk.AccAddress {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	return address.Hash("test", b)
}
