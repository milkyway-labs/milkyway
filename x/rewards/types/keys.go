package types

import (
	"cosmossdk.io/collections"
)

const (
	ModuleName = "rewards"
	StoreKey   = ModuleName
)

var (
	ParamsKey = collections.NewPrefix(0x01)

	NextRewardsPlanIDKey         = collections.NewPrefix(0x11)
	RewardsPlanKeyPrefix         = collections.NewPrefix(0x12)
	LastRewardsAllocationTimeKey = collections.NewPrefix(0x13)

	PoolDelegatorStartingInfoKeyPrefix = collections.NewPrefix(0x21)
	PoolHistoricalRewardsKeyPrefix     = collections.NewPrefix(0x22)
	PoolCurrentRewardsKeyPrefix        = collections.NewPrefix(0x23)
	PoolOutstandingRewardsKeyPrefix    = collections.NewPrefix(0x24)

	OperatorAccumulatedCommissionKeyPrefix = collections.NewPrefix(0x31)
	OperatorDelegatorStartingInfoKeyPrefix = collections.NewPrefix(0x32)
	OperatorHistoricalRewardsKeyPrefix     = collections.NewPrefix(0x33)
	OperatorCurrentRewardsKeyPrefix        = collections.NewPrefix(0x34)
	OperatorOutstandingRewardsKeyPrefix    = collections.NewPrefix(0x35)

	ServiceDelegatorStartingInfoKeyPrefix = collections.NewPrefix(0x41)
	ServiceHistoricalRewardsKeyPrefix     = collections.NewPrefix(0x42)
	ServiceCurrentRewardsKeyPrefix        = collections.NewPrefix(0x43)
	ServiceOutstandingRewardsKeyPrefix    = collections.NewPrefix(0x44)
)
