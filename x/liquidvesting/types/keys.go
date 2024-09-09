package types

import (
	"cosmossdk.io/collections"
)

const (
	ModuleName = "liquidvesting"
)

var (
	ParamsKey = collections.NewPrefix(0x01)

	InsuranceFundKey = collections.NewPrefix(0x10)
)
