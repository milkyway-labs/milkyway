package types

import (
	"encoding/binary"
	"fmt"
	"time"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "operators"
	StoreKey   = ModuleName

	DoNotModify = "[do-not-modify]"
)

var (
	ParamsKey = []byte{0x01}

	NextOperatorIDKey               = []byte{0xa1}
	OperatorPrefix                  = []byte{0xa2}
	InactivatingOperatorQueuePrefix = []byte{0xa3}
	OperatorAddressSetPrefix        = collections.NewPrefix(0xa4)
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

// OperatorStoreKey returns a operator ID into a key used to store a operator in the KVStore
func OperatorStoreKey(operatorID uint32) []byte {
	return append(OperatorPrefix, GetOperatorIDBytes(operatorID)...)
}

var (
	lenTime   = len(sdk.FormatTimeBytes(time.Now()))
	lenPrefix = len(InactivatingOperatorQueuePrefix)
)

// InactivatingOperatorByTime returns the key for all inactivating operators that expire at the given time
func InactivatingOperatorByTime(endTime time.Time) []byte {
	return append(InactivatingOperatorQueuePrefix, sdk.FormatTimeBytes(endTime)...)
}

// InactivatingOperatorQueueKey returns the key for an inactivating operator in the queue
func InactivatingOperatorQueueKey(operatorID uint32, endTime time.Time) []byte {
	return append(InactivatingOperatorByTime(endTime), GetOperatorIDBytes(operatorID)...)
}

// SplitInactivatingOperatorQueueKey split the inactivating operator queue key into operatorID and endTime
func SplitInactivatingOperatorQueueKey(key []byte) (operatorID uint32, endTime time.Time) {
	if len(key[lenPrefix:]) != 4+lenTime {
		panic(fmt.Errorf("unexpected key length (%d ≠ %d)", len(key[1:]), lenTime+4))
	}

	endTime, err := sdk.ParseTimeBytes(key[lenPrefix : lenPrefix+lenTime])
	if err != nil {
		panic(err)
	}

	operatorID = GetOperatorIDFromBytes(key[lenPrefix+lenTime:])
	return operatorID, endTime
}
