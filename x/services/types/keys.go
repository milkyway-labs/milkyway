package types

import (
	"encoding/binary"

	"cosmossdk.io/collections"
)

const (
	ModuleName = "services"
	StoreKey   = ModuleName

	DoNotModify = "[do-not-modify]"
)

var (
	ParamsKey = []byte{0x01}

	NextServiceIDKey        = []byte{0xa1}
	ServicePrefix           = []byte{0xa2}
	ServiceAddressSetPrefix = collections.NewPrefix(0xa3)
	ServiceParamsPrefix     = collections.NewPrefix(0xa4)
)

// GetServiceIDBytes returns the byte representation of the service ID
func GetServiceIDBytes(serviceID uint32) (serviceIDBz []byte) {
	serviceIDBz = make([]byte, 4)
	binary.BigEndian.PutUint32(serviceIDBz, serviceID)
	return serviceIDBz
}

// GetServiceIDFromBytes returns the service ID from a byte array
func GetServiceIDFromBytes(bz []byte) (serviceID uint32) {
	return binary.BigEndian.Uint32(bz)
}

// ServiceStoreKey turns a service ID into a key used to store a service in the KVStore
func ServiceStoreKey(serviceID uint32) []byte {
	return append(ServicePrefix, GetServiceIDBytes(serviceID)...)
}
