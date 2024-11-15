package types

import (
	"encoding/binary"

	"cosmossdk.io/collections"
)

const (
	ModuleName = "pools"
	StoreKey   = ModuleName
)

var (
	ParamsKey = collections.NewPrefix(0x01)

	NextPoolIDKey        = collections.NewPrefix(0xa1)
	PoolPrefix           = collections.NewPrefix(0xa2)
	PoolAddressSetPrefix = collections.NewPrefix(0xa3)
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
