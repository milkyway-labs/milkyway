package utils

import (
	"bytes"
	"encoding/binary"
)

func CompositeKey(parts ...[]byte) []byte {
	return bytes.Join(parts, nil)
}

// Uint32ToBigEndian marshals uint32 to a bigendian byte slice so it can be sorted
func Uint32ToBigEndian(i uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, i)
	return b
}
