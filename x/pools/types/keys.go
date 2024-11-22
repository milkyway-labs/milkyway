package types

import (
	"encoding/binary"
)

const (
	ModuleName = "pools"
	StoreKey   = ModuleName
)

var (
	NextPoolIDKey        = []byte{0xa1}
	PoolPrefix           = []byte{0xa2}
	PoolAddressSetPrefix = []byte{0xa3}
)

// GetPoolIDBytes returns the byte representation of the pool ID
func GetPoolIDBytes(poolID uint32) (poolIDBz []byte) {
	poolIDBz = make([]byte, 4)
	binary.BigEndian.PutUint32(poolIDBz, poolID)
	return poolIDBz
}
