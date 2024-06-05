package types

import (
	"encoding/binary"
)

const (
	ModuleName = "operators"

	DoNotModify = "[do-not-modify]"
)

var (
	ParamsKey = []byte{0x01}

	NextOperatorIDKey = []byte{0xa1}
	OperatorPrefix    = []byte{0xa2}
)

// GetOperatorIDBytes returns the byte representation of the operator ID
func GetOperatorIDBytes(operatorID uint32) (operatorIDBz []byte) {
	operatorIDBz = make([]byte, 4)
	binary.BigEndian.PutUint32(operatorIDBz, operatorID)
	return operatorIDBz
}

// GetOperatorIDFromBytes returns the operator ID from a byte array
func GetOperatorIDFromBytes(bz []byte) (operatorID uint32) {
	return binary.BigEndian.Uint32(bz)
}

// OperatorStoreKey turns a operator ID into a key used to store a operator in the KVStore
func OperatorStoreKey(operatorID uint32) []byte {
	return append(OperatorPrefix, GetOperatorIDBytes(operatorID)...)
}
