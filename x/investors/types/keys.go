package types

import (
	"cosmossdk.io/collections"
)

const (
	ModuleName = "investors"
	StoreKey   = ModuleName
)

var (
	InvestorsRewardRatioKey           = collections.NewPrefix(0x01)
	ValidatorsInvestorSharesKeyPrefix = collections.NewPrefix(0x02)
)
