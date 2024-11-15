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
	ParamsKey = []byte{0x01}

	NextServiceIDKey        = []byte{0xa1}
	ServicePrefix           = []byte{0xa2}
	ServiceAddressSetPrefix = []byte{0xa3}
	ServiceParamsPrefix     = []byte{0xa4}
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
