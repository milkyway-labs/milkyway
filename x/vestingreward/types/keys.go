package types

import (
	"cosmossdk.io/collections"
)

const (
	ModuleName = "vestingreward"
	StoreKey   = ModuleName
)

var (
	VestingAccountsRewardRatioKey           = collections.NewPrefix(0x01)
	ValidatorsVestingAccountSharesKeyPrefix = collections.NewPrefix(0x02)
)
