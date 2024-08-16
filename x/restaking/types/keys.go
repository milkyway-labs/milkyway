package types

import (
	"bytes"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/utils"
	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
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

	OperatorParamsPrefix = []byte{0x11}
	ServiceParamsPrefix  = []byte{0x12}

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
)

// OperatorParamsStoreKey returns the key used to store the operator params
func OperatorParamsStoreKey(operatorID uint32) []byte {
	return utils.CompositeKey(OperatorParamsPrefix, operatorstypes.GetOperatorIDBytes(operatorID))
}

// ServiceParamsStoreKey returns the key used to store the service params
func ServiceParamsStoreKey(serviceID uint32) []byte {
	return utils.CompositeKey(ServiceParamsPrefix, servicestypes.GetServiceIDBytes(serviceID))
}

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

// ParseDelegationsByPoolIDKey parses the pool ID and delegator address from the given key
func ParseDelegationsByPoolIDKey(bz []byte) (poolID uint32, delegatorAddress string, err error) {
	prefixLength := len(PoolDelegationsByPoolIDPrefix)
	if prefix := bz[:prefixLength]; !bytes.Equal(prefix, PoolDelegationsByPoolIDPrefix) {
		return 0, "", fmt.Errorf("invalid prefix; expected: %X, got: %x", PoolDelegationsByPoolIDPrefix, prefix)
	}

	// Remove the prefix
	bz = bz[prefixLength:]

	// Read the pool ID
	poolID = poolstypes.GetPoolIDFromBytes(bz[:4])
	bz = bz[4:]

	// Read the delegator address
	delegatorAddress = string(bz)

	return poolID, delegatorAddress, nil
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

// ParseDelegationsByOperatorIDKey parses the operator ID and delegator address from the given key
func ParseDelegationsByOperatorIDKey(bz []byte) (operatorID uint32, delegatorAddress string, err error) {
	prefixLength := len(OperatorDelegationPrefix)
	if prefix := bz[:prefixLength]; !bytes.Equal(prefix, OperatorDelegationPrefix) {
		return 0, "", fmt.Errorf("invalid prefix; expected: %X, got: %x", OperatorDelegationPrefix, prefix)
	}

	// Remove the prefix
	bz = bz[prefixLength:]

	// Read the operator ID
	operatorID = operatorstypes.GetOperatorIDFromBytes(bz[:4])
	bz = bz[4:]

	// Read the delegator address
	delegatorAddress = string(bz)

	return operatorID, delegatorAddress, nil
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

// ParseDelegationsByServiceIDKey parses the service ID and delegator address from the given key
func ParseDelegationsByServiceIDKey(bz []byte) (serviceID uint32, delegatorAddress string, err error) {
	prefixLength := len(ServiceDelegationPrefix)
	if prefix := bz[:prefixLength]; !bytes.Equal(prefix, ServiceDelegationPrefix) {
		return 0, "", fmt.Errorf("invalid prefix; expected: %X, got: %x", ServiceDelegationPrefix, prefix)
	}

	// Remove the prefix
	bz = bz[prefixLength:]

	// Read the service ID
	serviceID = servicestypes.GetServiceIDFromBytes(bz[:4])
	bz = bz[4:]

	// Read the delegator address
	delegatorAddress = string(bz)

	return serviceID, delegatorAddress, nil
}

// ServiceUnbondingDelegationsStorePrefix returns the prefix used to store all the unbonding delegations to a given pool
func ServiceUnbondingDelegationsStorePrefix(delegatorAddress string) []byte {
	return append(ServiceUnbondingDelegationPrefix, []byte(delegatorAddress)...)
}

// UserServiceUnbondingDelegationKey returns the key used to store the unbonding delegation for the given pool and delegator
func UserServiceUnbondingDelegationKey(delegatorAddress string, serviceID uint32) []byte {
	return append(ServiceUnbondingDelegationsStorePrefix(delegatorAddress), servicestypes.GetServiceIDBytes(serviceID)...)
}
