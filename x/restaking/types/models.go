package types

import (
	"fmt"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/utils"
)

func NewOperatorParams(commissionRate math.LegacyDec, joinedServicesIDs []uint32) OperatorParams {
	return OperatorParams{
		CommissionRate:    commissionRate,
		JoinedServicesIDs: joinedServicesIDs,
	}
}

func DefaultOperatorParams() OperatorParams {
	return NewOperatorParams(math.LegacyZeroDec(), nil)
}

func (p *OperatorParams) Validate() error {
	if p.CommissionRate.IsNegative() || p.CommissionRate.GT(math.LegacyOneDec()) {
		return fmt.Errorf("invalid commission rate: %s", p.CommissionRate.String())
	}

	if duplicate := utils.FindDuplicate(p.JoinedServicesIDs); duplicate != nil {
		return fmt.Errorf("duplicated joined service id: %v", duplicate)
	}

	for _, serviceID := range p.JoinedServicesIDs {
		if serviceID == 0 {
			return fmt.Errorf("invalid joined service id: %d", serviceID)
		}
	}
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

func NewServiceParams(
	slashFraction math.LegacyDec,
	whitelistedPoolsIDs []uint32,
	whitelistedOperatorsIDs []uint32,
) ServiceParams {
	return ServiceParams{
		SlashFraction:           slashFraction,
		WhitelistedPoolsIDs:     whitelistedPoolsIDs,
		WhitelistedOperatorsIDs: whitelistedOperatorsIDs,
	}
}

func DefaultServiceParams() ServiceParams {
	return NewServiceParams(math.LegacyZeroDec(), nil, nil)
}

func (p *ServiceParams) Validate() error {
	if p.SlashFraction.IsNegative() || p.SlashFraction.GT(math.LegacyOneDec()) {
		return fmt.Errorf("invalid slash fraction: %s", p.SlashFraction.String())
	}

	if duplicate := utils.FindDuplicate(p.WhitelistedPoolsIDs); duplicate != nil {
		return fmt.Errorf("duplicated whitelisted pool id: %v", duplicate)
	}

	if duplicate := utils.FindDuplicate(p.WhitelistedOperatorsIDs); duplicate != nil {
		return fmt.Errorf("duplicated whitelisted operator id: %v", duplicate)
	}

	for _, poolID := range p.WhitelistedPoolsIDs {
		if poolID == 0 {
			return fmt.Errorf("invalid whitelisted pool id: %d", poolID)
		}
	}

	for _, operatorID := range p.WhitelistedOperatorsIDs {
		if operatorID == 0 {
			return fmt.Errorf("invalid whitelisted operator id: %d", operatorID)
		}
	}
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

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
