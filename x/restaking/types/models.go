package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewPoolDelegation creates a new Delegation instance representing a pool delegation
func NewPoolDelegation(poolID uint32, userAddress string, shares sdk.DecCoins) Delegation {
	return Delegation{
		Type:        DELEGATION_TYPE_POOL,
		UserAddress: userAddress,
		TargetID:    poolID,
		Shares:      shares,
	}
}

// NewOperatorDelegation creates a new Delegation instance representing a delegation to an operator
func NewOperatorDelegation(operatorID uint32, userAddress string, shares sdk.DecCoins) Delegation {
	return Delegation{
		Type:        DELEGATION_TYPE_OPERATOR,
		TargetID:    operatorID,
		UserAddress: userAddress,
		Shares:      shares,
	}
}

// NewServiceDelegation creates a new Delegation instance representing a delegation to a service
func NewServiceDelegation(serviceID uint32, userAddress string, shares sdk.DecCoins) Delegation {
	return Delegation{
		Type:        DELEGATION_TYPE_SERVICE,
		TargetID:    serviceID,
		UserAddress: userAddress,
		Shares:      shares,
	}
}

// Validate validates the delegation
func (d Delegation) Validate() error {
	if d.Type == DELEGATION_TYPE_UNSPECIFIED {
		return fmt.Errorf("invalid delegation type")
	}

	if d.TargetID == 0 {
		return fmt.Errorf("invalid target id")
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

// MustMarshalDelegation marshals the given pool delegation using the provided codec
func MustMarshalDelegation(cdc codec.BinaryCodec, delegation Delegation) []byte {
	bz, err := cdc.Marshal(&delegation)
	if err != nil {
		panic(err)
	}
	return bz
}

// UnmarshalDelegation unmarshals a pool delegation from the given bytes using the provided codec
func UnmarshalDelegation(cdc codec.BinaryCodec, bz []byte) (Delegation, error) {
	var delegation Delegation
	err := cdc.Unmarshal(bz, &delegation)
	if err != nil {
		return Delegation{}, err
	}
	return delegation, nil
}

// MustUnmarshalDelegation unmarshals a pool delegation from the given bytes using the provided codec
func MustUnmarshalDelegation(cdc codec.BinaryCodec, bz []byte) Delegation {
	delegation, err := UnmarshalDelegation(cdc, bz)
	if err != nil {
		panic(err)
	}
	return delegation
}

// NewDelegationResponse creates a new DelegationResponse instance
func NewDelegationResponse(delegation Delegation, balance sdk.Coins) DelegationResponse {
	return DelegationResponse{
		Delegation: delegation,
		Balance:    balance,
	}
}
