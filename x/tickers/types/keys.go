package types

import (
	"cosmossdk.io/collections"
)

const (
	ModuleName = "tickers"
	StoreKey   = ModuleName
)

var (
	ParamsKey = collections.NewPrefix(0x01)

	TickerKeyPrefix      = collections.NewPrefix(0x11)
	TickerIndexKeyPrefix = collections.NewPrefix(0x12)
)
