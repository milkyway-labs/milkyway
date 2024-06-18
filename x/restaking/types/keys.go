package types

import (
	"bytes"
	"fmt"

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
	PoolDelegationsByPoolIDPrefix = []byte{0x71}

	ServiceDelegationPrefix          = []byte{0xb1}
	UnbondingServiceDelegationPrefix = []byte{0xb2}

	OperatorDelegationPrefix          = []byte{0xc1}
	UnbondingOperatorDelegationPrefix = []byte{0xc2}
)

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
