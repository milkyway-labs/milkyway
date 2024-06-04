package types

import (
	"encoding/binary"
)

const (
	ModuleName = "services"
	StoreKey   = ModuleName

	DoNotModify = "[do-not-modify]"
)

var (
	AVSPrefix = []byte{0x01}

	ParamsKey = []byte{0x10}
)

// NextAVSIDKey returns the key for the next Service ID
func NextAVSIDKey() []byte {
	return []byte{0x01}
}

// GetAVSIDBytes returns the byte representation of the Service ID
func GetAVSIDBytes(avsID uint32) (avsIDBz []byte) {
	avsIDBz = make([]byte, 4)
	binary.BigEndian.PutUint32(avsIDBz, avsID)
	return avsIDBz
}

// GetAVSIDFromBytes returns the Service ID from a byte array
func GetAVSIDFromBytes(bz []byte) (avsID uint32) {
	return binary.BigEndian.Uint32(bz)
}

// AVSStoreKey turns an Service ID into a key used to store an Service in the KVStore
func AVSStoreKey(avsID uint32) []byte {
	return append(AVSPrefix, GetAVSIDBytes(avsID)...)
}
