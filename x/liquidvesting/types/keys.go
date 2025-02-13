package types

import (
	"time"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "liquidvesting"
	StoreKey   = ModuleName
)

var (
	ParamsKey = collections.NewPrefix(0x01)

	InsuranceFundKey             = collections.NewPrefix(0x10)
	BurnCoinsQueueKey            = collections.NewPrefix(0x20)
	CoveredLockedSharesKeyPrefix = collections.NewPrefix(0x30)
)

// GetBurnCoinsQueueTimeKey creates the prefix to obtain the list of
// coins to burn for each delegator
func GetBurnCoinsQueueTimeKey(timestamp time.Time) []byte {
	bz := sdk.FormatTimeBytes(timestamp)
	return append(BurnCoinsQueueKey, bz...)
}
