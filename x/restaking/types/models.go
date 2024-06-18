package types

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
)

// NewPoolDelegation creates a new PoolDelegation instance
func NewPoolDelegation(poolID uint32, userAddress string, shares sdkmath.LegacyDec) PoolDelegation {
	return PoolDelegation{
		PoolID:      poolID,
		UserAddress: userAddress,
		Shares:      shares,
	}
}

// --------------------------------------------------------------------------------------------------------------------

// MustMarshalPoolDelegation marshals the given pool delegation using the provided codec
func MustMarshalPoolDelegation(cdc codec.BinaryCodec, delegation PoolDelegation) []byte {
	bz, err := cdc.Marshal(&delegation)
	if err != nil {
		panic(err)
	}
	return bz
}

// UnmarshalPoolDelegation unmarshals a pool delegation from the given bytes using the provided codec
func UnmarshalPoolDelegation(cdc codec.BinaryCodec, bz []byte) (PoolDelegation, error) {
	var delegation PoolDelegation
	err := cdc.Unmarshal(bz, &delegation)
	if err != nil {
		return PoolDelegation{}, err
	}
	return delegation, nil
}

// MustUnmarshalPoolDelegation unmarshals a pool delegation from the given bytes using the provided codec
func MustUnmarshalPoolDelegation(cdc codec.BinaryCodec, bz []byte) PoolDelegation {
	delegation, err := UnmarshalPoolDelegation(cdc, bz)
	if err != nil {
		panic(err)
	}
	return delegation
}
