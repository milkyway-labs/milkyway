package move

import (
	"encoding/hex"
	"strings"

	"golang.org/x/crypto/sha3"
)

const (
	DenomTraceDenomPrefixMove = "move/"
)

// Generate named object address from the seed (address + name + 0xFE)
func NamedObjectAddress(source AccountAddress, name string) AccountAddress {
	// 0xFE is the suffix of named object address, which is
	// defined in object.move as `OBJECT_FROM_SEED_ADDRESS_SCHEME`.
	hasher := sha3.New256()
	hasher.Write(append(append(source[:], []byte(name)...), 0xFE))
	bz := hasher.Sum(nil)

	addr, err := NewAccountAddressFromBytes(bz[:])
	if err != nil {
		panic(err)
	}

	return addr
}

func UserDerivedObjectAddress(source AccountAddress, deriveFrom AccountAddress) AccountAddress {
	hasher := sha3.New256()
	hasher.Write(append(append(source[:], deriveFrom[:]...), 0xFC))
	bz := hasher.Sum(nil)

	addr, err := NewAccountAddressFromBytes(bz[:])
	if err != nil {
		panic(err)
	}

	return addr
}

// Extract metadata address from a denom
func MetadataAddressFromDenom(denom string) (AccountAddress, error) {
	if strings.HasPrefix(denom, DenomTraceDenomPrefixMove) {
		addrBz, err := hex.DecodeString(strings.TrimPrefix(denom, DenomTraceDenomPrefixMove))
		if err != nil {
			return AccountAddress{}, err
		}

		return NewAccountAddressFromBytes(addrBz)
	}

	// non move coins are generated from 0x1.
	return NamedObjectAddress(StdAddress, denom), nil
}
