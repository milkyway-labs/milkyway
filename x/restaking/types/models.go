package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewPoolDelegation creates a new PoolDelegation instance
func NewPoolDelegation(poolID uint32, userAddress string, shares sdkmath.LegacyDec) PoolDelegation {
	return PoolDelegation{
		PoolID:      poolID,
		UserAddress: userAddress,
		Shares:      shares,
	}
}

func (d PoolDelegation) Validate() error {
	if d.PoolID == 0 {
		return fmt.Errorf("invalid pool id")
	}

	_, err := sdk.AccAddressFromBech32(d.UserAddress)
	if err != nil {
		return fmt.Errorf("invalid user address: %s", d.UserAddress)
	}

	if d.Shares.IsNegative() {
		return ErrInvalidShares
	}

	return nil
}

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

// --------------------------------------------------------------------------------------------------------------------

func NewOperatorDelegation(operatorID uint32, userAddress string, shares sdkmath.LegacyDec) OperatorDelegation {
	return OperatorDelegation{
		OperatorID:  operatorID,
		UserAddress: userAddress,
		Shares:      shares,
	}
}

func (d OperatorDelegation) Validate() error {
	if d.OperatorID == 0 {
		return fmt.Errorf("invalid operator id")
	}

	_, err := sdk.AccAddressFromBech32(d.UserAddress)
	if err != nil {
		return fmt.Errorf("invalid user address: %s", d.UserAddress)
	}

	if d.Shares.IsNegative() {
		return ErrInvalidShares
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

func NewServiceDelegation(serviceID uint32, userAddress string, shares sdkmath.LegacyDec) ServiceDelegation {
	return ServiceDelegation{
		ServiceID:   serviceID,
		UserAddress: userAddress,
		Shares:      shares,
	}
}

func (d ServiceDelegation) Validate() error {
	if d.ServiceID == 0 {
		return fmt.Errorf("invalid service id")
	}

	_, err := sdk.AccAddressFromBech32(d.UserAddress)
	if err != nil {
		return fmt.Errorf("invalid user address: %s", d.UserAddress)
	}

	if d.Shares.IsNegative() {
		return ErrInvalidShares
	}

	return nil
}
