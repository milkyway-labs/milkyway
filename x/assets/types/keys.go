package types

import (
	"cosmossdk.io/collections"
)

const (
	ModuleName = "assets"
	StoreKey   = ModuleName
)

var (
	ParamsKey = collections.NewPrefix(0x01)

	AssetKeyPrefix       = collections.NewPrefix(0x11)
	TickerIndexKeyPrefix = collections.NewPrefix(0x12)
)
