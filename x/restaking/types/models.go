package types

import (
	"fmt"
	"time"

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

// --------------------------------------------------------------------------------------------------------------------

func NewUnbondingDelegationEntry(creationHeight int64, completionTime time.Time, balance sdk.Coins, unbondingID uint64) UnbondingDelegationEntry {
	return UnbondingDelegationEntry{
		CreationHeight: creationHeight,
		CompletionTime: completionTime,
		InitialBalance: balance,
		Balance:        balance,
		UnbondingId:    unbondingID,
	}
}

func NewPoolUnbondingDelegation(
	delegatorAddress string, poolID uint32,
	creationHeight int64, minTime time.Time, balance sdk.Coins, id uint64,
) UnbondingDelegation {
	return UnbondingDelegation{
		Type:             UNBONDING_DELEGATION_TYPE_POOL,
		DelegatorAddress: delegatorAddress,
		TargetID:         poolID,
		Entries: []UnbondingDelegationEntry{
			NewUnbondingDelegationEntry(creationHeight, minTime, balance, id),
		},
	}
}

// AddEntry - append entry to the unbonding delegation
func (ubd *UnbondingDelegation) AddEntry(creationHeight int64, minTime time.Time, balance sdk.Coins, unbondingID uint64) {
	// Check the entries exists with creation_height and complete_time
	entryIndex := -1
	for index, ubdEntry := range ubd.Entries {
		if ubdEntry.CreationHeight == creationHeight && ubdEntry.CompletionTime.Equal(minTime) {
			entryIndex = index
			break
		}
	}

	// entryIndex exists
	if entryIndex != -1 {
		ubdEntry := ubd.Entries[entryIndex]
		ubdEntry.Balance = ubdEntry.Balance.Add(balance...)
		ubdEntry.InitialBalance = ubdEntry.InitialBalance.Add(balance...)

		// Update the entry
		ubd.Entries[entryIndex] = ubdEntry
	} else {
		// Append the new unbonding delegation entry
		entry := NewUnbondingDelegationEntry(creationHeight, minTime, balance, unbondingID)
		ubd.Entries = append(ubd.Entries, entry)
	}
}

// MarshalUnbondingDelegation marshals the unbonding delegation using the provided codec
func MarshalUnbondingDelegation(cdc codec.BinaryCodec, unbondingDelegation UnbondingDelegation) []byte {
	bz, err := cdc.Marshal(&unbondingDelegation)
	if err != nil {
		panic(err)
	}
	return bz
}

// MustMarshalUnbondingDelegation marshals the unbonding delegation using the provided codec
func MustMarshalUnbondingDelegation(cdc codec.BinaryCodec, unbondingDelegation UnbondingDelegation) []byte {
	return MarshalUnbondingDelegation(cdc, unbondingDelegation)
}

// UnmarshalUnbondingDelegation unmarshals the unbonding delegation from the given bytes using the provided codec
func UnmarshalUnbondingDelegation(cdc codec.BinaryCodec, bz []byte) (UnbondingDelegation, error) {
	var unbondingDelegation UnbondingDelegation
	err := cdc.Unmarshal(bz, &unbondingDelegation)
	if err != nil {
		return UnbondingDelegation{}, err
	}
	return unbondingDelegation, nil
}

// MustUnmarshalUnbondingDelegation unmarshals the unbonding delegation from the given bytes using the provided codec
func MustUnmarshalUnbondingDelegation(cdc codec.BinaryCodec, bz []byte) UnbondingDelegation {
	unbondingDelegation, err := UnmarshalUnbondingDelegation(cdc, bz)
	if err != nil {
		panic(err)
	}
	return unbondingDelegation
}
