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

// isDelegation implements Delegation
func (d PoolDelegation) isDelegation() {}

// Validate validates the pool delegation
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

// NewPoolDelegationResponse creates a new PoolDelegationResponse instance
func NewPoolDelegationResponse(poolID uint32, userAddress string, shares sdkmath.LegacyDec, balance sdk.Coin) PoolDelegationResponse {
	return PoolDelegationResponse{
		Delegation: PoolDelegation{
			UserAddress: userAddress,
			PoolID:      poolID,
			Shares:      shares,
		},
		Balance: balance,
	}

}

// --------------------------------------------------------------------------------------------------------------------

func NewOperatorDelegation(operatorID uint32, userAddress string, shares sdk.DecCoins) OperatorDelegation {
	return OperatorDelegation{
		OperatorID:  operatorID,
		UserAddress: userAddress,
		Shares:      shares,
	}
}

func (d OperatorDelegation) isDelegation() {}

func (d OperatorDelegation) Validate() error {
	if d.OperatorID == 0 {
		return fmt.Errorf("invalid operator id")
	}

	_, err := sdk.AccAddressFromBech32(d.UserAddress)
	if err != nil {
		return fmt.Errorf("invalid user address: %s", d.UserAddress)
	}

	if d.Shares.IsAnyNegative() {
		return ErrInvalidShares
	}

	return nil
}

// MustMarshalOperatorDelegation marshals the given operator delegation using the provided codec
func MustMarshalOperatorDelegation(cdc codec.BinaryCodec, delegation OperatorDelegation) []byte {
	bz, err := cdc.Marshal(&delegation)
	if err != nil {
		panic(err)
	}
	return bz
}

// UnmarshalOperatorDelegation unmarshals an operator delegation from the given bytes using the provided codec
func UnmarshalOperatorDelegation(cdc codec.BinaryCodec, bz []byte) (OperatorDelegation, error) {
	var delegation OperatorDelegation
	err := cdc.Unmarshal(bz, &delegation)
	if err != nil {
		return OperatorDelegation{}, err
	}
	return delegation, nil
}

// MustUnmarshalOperatorDelegation unmarshals an operator delegation from the given bytes using the provided codec
func MustUnmarshalOperatorDelegation(cdc codec.BinaryCodec, bz []byte) OperatorDelegation {
	delegation, err := UnmarshalOperatorDelegation(cdc, bz)
	if err != nil {
		panic(err)
	}
	return delegation
}

// NewOperatorDelegationResponse creates a new OperatorDelegationResponse instance
func NewOperatorDelegationResponse(operatorID uint32, userAddress string, shares sdk.DecCoins, balance sdk.Coins) OperatorDelegationResponse {
	return OperatorDelegationResponse{
		Delegation: OperatorDelegation{
			UserAddress: userAddress,
			OperatorID:  operatorID,
			Shares:      shares,
		},
		Balance: balance,
	}
}

// --------------------------------------------------------------------------------------------------------------------

func NewServiceDelegation(serviceID uint32, userAddress string, shares sdk.DecCoins) ServiceDelegation {
	return ServiceDelegation{
		ServiceID:   serviceID,
		UserAddress: userAddress,
		Shares:      shares,
	}
}

// isDelegation implements Delegation
func (d ServiceDelegation) isDelegation() {}

// Validate validates the service delegation
func (d ServiceDelegation) Validate() error {
	if d.ServiceID == 0 {
		return fmt.Errorf("invalid service id")
	}

	_, err := sdk.AccAddressFromBech32(d.UserAddress)
	if err != nil {
		return fmt.Errorf("invalid user address: %s", d.UserAddress)
	}

	if d.Shares.IsAnyNegative() {
		return ErrInvalidShares
	}

	return nil
}

// MustMarshalServiceDelegation marshals the given service delegation using the provided codec
func MustMarshalServiceDelegation(cdc codec.BinaryCodec, delegation ServiceDelegation) []byte {
	bz, err := cdc.Marshal(&delegation)
	if err != nil {
		panic(err)
	}
	return bz
}

// UnmarshalServiceDelegation unmarshals a service delegation from the given bytes using the provided codec
func UnmarshalServiceDelegation(cdc codec.BinaryCodec, bz []byte) (ServiceDelegation, error) {
	var delegation ServiceDelegation
	err := cdc.Unmarshal(bz, &delegation)
	if err != nil {
		return ServiceDelegation{}, err
	}
	return delegation, nil
}

// MustUnmarshalServiceDelegation unmarshals a service delegation from the given bytes using the provided codec
func MustUnmarshalServiceDelegation(cdc codec.BinaryCodec, bz []byte) ServiceDelegation {
	delegation, err := UnmarshalServiceDelegation(cdc, bz)
	if err != nil {
		panic(err)
	}
	return delegation
}

// NewServiceDelegationResponse creates a new ServiceDelegationResponse instance
func NewServiceDelegationResponse(serviceID uint32, userAddress string, shares sdk.DecCoins, balance sdk.Coins) ServiceDelegationResponse {
	return ServiceDelegationResponse{
		Delegation: ServiceDelegation{
			UserAddress: userAddress,
			ServiceID:   serviceID,
			Shares:      shares,
		},
		Balance: balance,
	}
}
