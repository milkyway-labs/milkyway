package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
)

const (
	ModuleName = "restaking"
	StoreKey   = ModuleName
)

var (
	ParamsKey = []byte{0x01}

	PoolDelegationPrefix          = []byte{0xa1}
	UnbondingPoolDelegationPrefix = []byte{0xa2}

	ServiceDelegationPrefix          = []byte{0xb1}
	UnbondingServiceDelegationPrefix = []byte{0xb2}

	OperatorDelegationPrefix          = []byte{0xc1}
	UnbondingOperatorDelegationPrefix = []byte{0xc2}
)

// PoolDelegationsStorePrefix returns the prefix used to store all the delegations to a given pool
func PoolDelegationsStorePrefix(poolID uint32) []byte {
	return append(PoolDelegationPrefix, poolstypes.GetPoolIDBytes(poolID)...)
}

// UserPoolDelegationStoreKey returns the key used to store the delegation of a user to a given pool
func UserPoolDelegationStoreKey(poolID uint32, delegator string) []byte {
	return append(PoolDelegationsStorePrefix(poolID), []byte(delegator)...)
}

// UnbondingDelegationsByTimeStorePrefix returns the prefix used to store all the unbonding delegations
// that expire at the given time
func UnbondingDelegationsByTimeStorePrefix(endTime time.Time) []byte {
	return append(UnbondingPoolDelegationPrefix, sdk.FormatTimeBytes(endTime)...)
}

// UnbondingPoolDelegationsStorePrefix returns the prefix used to store all the unbonding delegations
// to a given pool that expire at the given time
func UnbondingPoolDelegationsStorePrefix(poolID uint32, endTime time.Time) []byte {
	return append(UnbondingDelegationsByTimeStorePrefix(endTime), poolstypes.GetPoolIDBytes(poolID)...)
}

// UserUnbondingPoolDelegationStoreKey returns the key used to store the unbonding delegation of a user
// to a given pool that expires at the given time
func UserUnbondingPoolDelegationStoreKey(poolID uint32, delegator string, endTime time.Time) []byte {
	return append(UnbondingPoolDelegationsStorePrefix(poolID, endTime), []byte(delegator)...)
}
