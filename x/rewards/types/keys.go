package types

import (
	"cosmossdk.io/collections"
)

const (
	ModuleName = "rewards"
	StoreKey   = ModuleName

	// RewardsPoolName is the account name of global rewards pool where rewards
	// are moved from each rewards plan's rewards pool and distributed to
	// delegators later.
	RewardsPoolName = "rewards_pool"
)

var (
	ParamsKey = collections.NewPrefix(0x01)

	NextRewardsPlanIDKey           = collections.NewPrefix(0xa1)
	RewardsPlanKeyPrefix           = collections.NewPrefix(0xa2)
	LastRewardsAllocationTimeKey   = collections.NewPrefix(0xa3)
	DelegatorWithdrawAddrKeyPrefix = collections.NewPrefix(0xa4)

	PoolDelegatorStartingInfoKeyPrefix       = collections.NewPrefix(0xb1)
	PoolHistoricalRewardsKeyPrefix           = collections.NewPrefix(0xb2)
	PoolCurrentRewardsKeyPrefix              = collections.NewPrefix(0xb3)
	PoolOutstandingRewardsKeyPrefix          = collections.NewPrefix(0xb4)
	PoolServiceTotalDelegatorSharesKeyPrefix = collections.NewPrefix(0xb5)

	OperatorAccumulatedCommissionKeyPrefix = collections.NewPrefix(0xc1)
	OperatorDelegatorStartingInfoKeyPrefix = collections.NewPrefix(0xc2)
	OperatorHistoricalRewardsKeyPrefix     = collections.NewPrefix(0xc3)
	OperatorCurrentRewardsKeyPrefix        = collections.NewPrefix(0xc4)
	OperatorOutstandingRewardsKeyPrefix    = collections.NewPrefix(0xc5)

	ServiceDelegatorStartingInfoKeyPrefix = collections.NewPrefix(0xd1)
	ServiceHistoricalRewardsKeyPrefix     = collections.NewPrefix(0xd2)
	ServiceCurrentRewardsKeyPrefix        = collections.NewPrefix(0xd3)
	ServiceOutstandingRewardsKeyPrefix    = collections.NewPrefix(0xd4)
)
