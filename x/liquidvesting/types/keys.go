package types

import (
	"cosmossdk.io/collections"
)

const (
	ModuleName = "liquidvesting"
	StoreKey   = ModuleName
)

var (
	ParamsKey = collections.NewPrefix(0x01)

	InsuranceFundKey = collections.NewPrefix(0x10)
)
