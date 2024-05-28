package move

import (
	"bytes"
	"encoding/hex"
	"errors"
	"strings"
)

var StdAddress AccountAddress

// initialize StdAddress
func init() {
	var err error
	StdAddress, err = NewAccountAddress("0x1")
	if err != nil {
		panic(err)
	}
}

// NewAccountAddressFromBytes return AccountAddress from the bytes
func NewAccountAddressFromBytes(bz []byte) (AccountAddress, error) {
	lengthDiff := len(AccountAddress{}) - len(bz)
	if lengthDiff > 0 {
		bz = append(bytes.Repeat([]byte{0}, lengthDiff), bz...)
	} else if lengthDiff < 0 {
		return AccountAddress{}, errors.New("invalid length of address")
	}

	return BcsDeserializeAccountAddress(bz)
}

// NewAccountAddress return AccountAddress from the hex string
func NewAccountAddress(hexAddr string) (AccountAddress, error) {
	hexStr := strings.TrimPrefix(hexAddr, "0x")
	if len(hexStr)%2 == 1 {
		hexStr = "0" + hexStr
	}

	bz, err := hex.DecodeString(hexStr)
	if err != nil {
		return AccountAddress{}, errors.New("invalid hex address")
	}

	accountAddress, err := NewAccountAddressFromBytes(bz)
	return accountAddress, err
}

func (addr AccountAddress) Bytes() []byte {
	outBz := make([]byte, len(addr))
	copy(outBz, addr[:])
	return outBz
}
