package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	operatorstypes "github.com/milkyway-labs/milkyway/v4/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/v4/x/pools/types"
	servicestypes "github.com/milkyway-labs/milkyway/v4/x/services/types"
)

const (
	ModuleName = "restaking"
	StoreKey   = ModuleName
)

var (
	ParamsKey         = []byte{0x01}
	UnbondingIDKey    = []byte{0x02}
	UnbondingIndexKey = []byte{0x03}
	UnbondingTypeKey  = []byte{0x04}

	OperatorJoinedServicesPrefix       = []byte{0x13}
	ServiceOperatorsAllowListPrefix    = []byte{0x14}
	ServiceSecuringPoolsPrefix         = []byte{0x15}
	ServiceJoinedByOperatorIndexPrefix = []byte{0x16}

	PoolDelegationPrefix          = []byte{0xa1}
	PoolDelegationsByPoolIDPrefix = []byte{0xa2}
	PoolUnbondingDelegationPrefix = []byte{0xa3}

	OperatorDelegationPrefix          = []byte{0xb1}
	OperatorDelegationByOperatorID    = []byte{0xb2}
	OperatorUnbondingDelegationPrefix = []byte{0xb3}

	ServiceDelegationPrefix            = []byte{0xc1}
	ServiceDelegationByServiceIDPrefix = []byte{0xc2}
	ServiceUnbondingDelegationPrefix   = []byte{0xc3}

	UnbondingQueueKey = []byte{0xd1}

	UserPreferencesPrefix = []byte{0xe1}
)

// --------------------------------------------------------------------------------------------------------------------

// GetUnbondingIndexKey returns a key for the index for looking up UnbondingDelegations by the UnbondingDelegationEntries they contain
func GetUnbondingIndexKey(id uint64) []byte {
	return append(UnbondingIndexKey, sdk.Uint64ToBigEndian(id)...)
}

// GetUnbondingTypeKey returns a key for an index containing the type of unbonding operations
func GetUnbondingTypeKey(id uint64) []byte {
	return append(UnbondingTypeKey, sdk.Uint64ToBigEndian(id)...)
}

// GetUnbondingDelegationTimeKey creates the prefix for all unbonding delegations from a delegator
func GetUnbondingDelegationTimeKey(timestamp time.Time) []byte {
	bz := sdk.FormatTimeBytes(timestamp)
	return append(UnbondingQueueKey, bz...)
}

// --------------------------------------------------------------------------------------------------------------------

type DelegationKeyBuilder func(delegatorAddress string, targetID uint32) []byte

type DelegationByTargetIDBuilder func(targetID uint32, delegationAddress string) []byte

type UnbondingDelegationKeyBuilder func(delegatorAddress string, targetID uint32) []byte

// --------------------------------------------------------------------------------------------------------------------

// UserPoolDelegationsStorePrefix returns the prefix used to store all the delegations to a given pool
func UserPoolDelegationsStorePrefix(userAddress string) []byte {
	return append(PoolDelegationPrefix, []byte(userAddress)...)
}

// UserPoolDelegationStoreKey returns the key used to store the user -> pool delegation association
func UserPoolDelegationStoreKey(delegator string, poolID uint32) []byte {
	return append(UserPoolDelegationsStorePrefix(delegator), poolstypes.GetPoolIDBytes(poolID)...)
}

// DelegationsByPoolIDStorePrefix returns the prefix used to store the delegations to a given pool
func DelegationsByPoolIDStorePrefix(poolID uint32) []byte {
	return append(PoolDelegationsByPoolIDPrefix, poolstypes.GetPoolIDBytes(poolID)...)
}

// DelegationByPoolIDStoreKey returns the key used to store the pool -> user delegation association
func DelegationByPoolIDStoreKey(poolID uint32, delegatorAddress string) []byte {
	return append(DelegationsByPoolIDStorePrefix(poolID), []byte(delegatorAddress)...)
}

// PoolUnbondingDelegationsStorePrefix returns the prefix used to store all the unbonding delegations to a given pool
func PoolUnbondingDelegationsStorePrefix(delegatorAddress string) []byte {
	return append(PoolUnbondingDelegationPrefix, []byte(delegatorAddress)...)
}

// UserPoolUnbondingDelegationKey returns the key used to store the unbonding delegation for the given pool and delegator
func UserPoolUnbondingDelegationKey(delegatorAddress string, poolID uint32) []byte {
	return append(PoolUnbondingDelegationsStorePrefix(delegatorAddress), poolstypes.GetPoolIDBytes(poolID)...)
}

// --------------------------------------------------------------------------------------------------------------------

// UserOperatorDelegationsStorePrefix returns the prefix used to store all the delegations to a given operator
func UserOperatorDelegationsStorePrefix(userAddress string) []byte {
	return append(OperatorDelegationPrefix, []byte(userAddress)...)
}

// UserOperatorDelegationStoreKey returns the key used to store the user -> operator delegation association
func UserOperatorDelegationStoreKey(delegator string, operatorID uint32) []byte {
	return append(UserOperatorDelegationsStorePrefix(delegator), operatorstypes.GetOperatorIDBytes(operatorID)...)
}

// DelegationsByOperatorIDStorePrefix returns the prefix used to store the delegations to a given operator
func DelegationsByOperatorIDStorePrefix(operatorID uint32) []byte {
	return append(OperatorDelegationByOperatorID, operatorstypes.GetOperatorIDBytes(operatorID)...)
}

// DelegationByOperatorIDStoreKey returns the key used to store the operator -> user delegation association
func DelegationByOperatorIDStoreKey(operatorID uint32, delegatorAddress string) []byte {
	return append(DelegationsByOperatorIDStorePrefix(operatorID), []byte(delegatorAddress)...)
}

// OperatorUnbondingDelegationsStorePrefix returns the prefix used to store all the unbonding delegations to a given pool
func OperatorUnbondingDelegationsStorePrefix(delegatorAddress string) []byte {
	return append(OperatorUnbondingDelegationPrefix, []byte(delegatorAddress)...)
}

// UserOperatorUnbondingDelegationKey returns the key used to store the unbonding delegation for the given pool and delegator
func UserOperatorUnbondingDelegationKey(delegatorAddress string, operatorID uint32) []byte {
	return append(OperatorUnbondingDelegationsStorePrefix(delegatorAddress), operatorstypes.GetOperatorIDBytes(operatorID)...)
}

// --------------------------------------------------------------------------------------------------------------------

// UserServiceDelegationsStorePrefix returns the prefix used to store all the delegations to a given service
func UserServiceDelegationsStorePrefix(userAddress string) []byte {
	return append(ServiceDelegationPrefix, []byte(userAddress)...)
}

// UserServiceDelegationStoreKey returns the key used to store the user -> service delegation association
func UserServiceDelegationStoreKey(delegator string, serviceID uint32) []byte {
	return append(UserServiceDelegationsStorePrefix(delegator), servicestypes.GetServiceIDBytes(serviceID)...)
}

// DelegationsByServiceIDStorePrefix returns the prefix used to store the delegations to a given service
func DelegationsByServiceIDStorePrefix(serviceID uint32) []byte {
	return append(ServiceDelegationByServiceIDPrefix, servicestypes.GetServiceIDBytes(serviceID)...)
}

// DelegationByServiceIDStoreKey returns the key used to store the service -> user delegation association
func DelegationByServiceIDStoreKey(serviceID uint32, delegatorAddress string) []byte {
	return append(DelegationsByServiceIDStorePrefix(serviceID), []byte(delegatorAddress)...)
}

// ServiceUnbondingDelegationsStorePrefix returns the prefix used to store all the unbonding delegations to a given pool
func ServiceUnbondingDelegationsStorePrefix(delegatorAddress string) []byte {
	return append(ServiceUnbondingDelegationPrefix, []byte(delegatorAddress)...)
}

// UserServiceUnbondingDelegationKey returns the key used to store the unbonding delegation for the given pool and delegator
func UserServiceUnbondingDelegationKey(delegatorAddress string, serviceID uint32) []byte {
	return append(ServiceUnbondingDelegationsStorePrefix(delegatorAddress), servicestypes.GetServiceIDBytes(serviceID)...)
}
