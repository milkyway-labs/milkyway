package types

import (
	"encoding/binary"
)

const (
	ModuleName = "pools"
	StoreKey   = ModuleName
)

var (
	ParamsKey = []byte{0x01}

	NextPoolIDKey = []byte{0xa1}
	PoolPrefix    = []byte{0xa2}
)

// GetPoolIDBytes returns the byte representation of the pool ID
func GetPoolIDBytes(poolID uint32) (poolIDBz []byte) {
	poolIDBz = make([]byte, 4)
	binary.BigEndian.PutUint32(poolIDBz, poolID)
	return poolIDBz
}

// GetPoolIDFromBytes returns the pool ID from a byte array
func GetPoolIDFromBytes(bz []byte) (poolID uint32) {
	return binary.BigEndian.Uint32(bz)
}

// GetPoolStoreKey turns a pool ID into a key used to store a pool in the KVStore
func GetPoolStoreKey(poolID uint32) []byte {
	return append(PoolPrefix, GetPoolIDBytes(poolID)...)
}
