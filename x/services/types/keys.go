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
	ParamsKey = collections.NewPrefix(0x01)

	NextServiceIDKey        = collections.NewPrefix(0xa1)
	ServicePrefix           = collections.NewPrefix(0xa2)
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
