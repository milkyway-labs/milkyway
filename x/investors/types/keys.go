package types

import (
	"cosmossdk.io/collections"
)

const (
	ModuleName = "investors"
	StoreKey   = ModuleName
)

var (
	InvestorsRewardRatioKey         = collections.NewPrefix(0x01)
	InvestorsVestingQueueKeyPrefix  = collections.NewPrefix(0x02)
	VestingInvestorsKeyPrefix       = collections.NewPrefix(0x03)
	VestingInvestorRewardsKeyPrefix = collections.NewPrefix(0x04)
	CurrentDelegatorKey             = collections.NewPrefix(0x05)
)
