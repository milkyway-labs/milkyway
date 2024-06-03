package types

import (
	"encoding/binary"
)

const (
	ModuleName   = "avs"
	RouterKey    = ModuleName
	StoreKey     = ModuleName
	QuerierRoute = ModuleName
)

var (
	AVSPrefix = []byte{0x02}
)

// NextAVSIDKey returns the key for the next AVS ID
func NextAVSIDKey() []byte {
	return []byte{0x01}
}

// GetAVSIDBytes returns the byte representation of the AVS ID
func GetAVSIDBytes(avsID uint32) (avsIDBz []byte) {
	avsIDBz = make([]byte, 4)
	binary.BigEndian.PutUint32(avsIDBz, avsID)
	return avsIDBz
}

// GetAVSIDFromBytes returns the AVS ID from a byte array
func GetAVSIDFromBytes(bz []byte) (avsID uint32) {
	return binary.BigEndian.Uint32(bz)
}

// AVSStoreKey turns an AVS ID into a key used to store an AVS in the KVStore
func AVSStoreKey(avsID uint32) []byte {
	return append(AVSPrefix, GetAVSIDBytes(avsID)...)
}
