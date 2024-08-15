package utils

import (
	"bytes"
	"encoding/binary"
)

func CompositeKey(parts ...[]byte) []byte {
	return bytes.Join(parts, nil)
}

// Uint32ToBytes marshals uint32 to a bigendian byte slice so it can be sorted
func Uint32ToBytes(i uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, i)
	return b
}

// Uint32FromBytes returns an uint32 from big endian encoded bytes. If encoding
// is empty, zero is returned.
func Uint32FromBytes(bz []byte) uint32 {
	if len(bz) == 0 {
		return 0
	}
	return binary.BigEndian.Uint32(bz)
}
