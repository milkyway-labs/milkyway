package types

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/utils"
	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

// NewServiceParams creates a new ServiceParams instance
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

// DefaultServiceParams returns the default service params
func DefaultServiceParams() ServiceParams {
	return NewServiceParams(math.LegacyZeroDec(), nil, nil)
}

// Validate validates the service params
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

// GetDelegationTypeFromTarget returns the unbonding delegation type based on the target
func GetDelegationTypeFromTarget(target DelegationTarget) (DelegationType, error) {
	switch target.(type) {
	case *poolstypes.Pool:
		return DELEGATION_TYPE_POOL, nil
	case *operatorstypes.Operator:
		return DELEGATION_TYPE_OPERATOR, nil
	case *servicestypes.Service:
		return DELEGATION_TYPE_SERVICE, nil
	default:
		return DELEGATION_TYPE_UNSPECIFIED, fmt.Errorf("invalid unbonding target type")
	}
}

// NewUnbondingDelegationEntry creates a new UnbondingDelegationEntry instance
func NewUnbondingDelegationEntry(creationHeight int64, completionTime time.Time, balance sdk.Coins, unbondingID uint64) UnbondingDelegationEntry {
	return UnbondingDelegationEntry{
		CreationHeight: creationHeight,
		CompletionTime: completionTime,
		InitialBalance: balance,
		Balance:        balance,
		UnbondingID:    unbondingID,
	}
}

// IsMature tells whether is the current entry mature
func (e UnbondingDelegationEntry) IsMature(currentTime time.Time) bool {
	return !e.CompletionTime.After(currentTime)
}

// Validate validates the unbonding delegation entry
func (e UnbondingDelegationEntry) Validate() error {
	if e.UnbondingID == 0 {
		return fmt.Errorf("invalid unbonding id")
	}

	if e.CreationHeight == 0 {
		return fmt.Errorf("invalid creation height")
	}

	if e.CompletionTime.IsZero() {
		return fmt.Errorf("invalid completion time")
	}

	if !e.InitialBalance.IsValid() {
		return fmt.Errorf("invalid initial balance")
	}

	if !e.Balance.IsValid() {
		return fmt.Errorf("invalid balance")
	}

	return nil
}

// NewPoolUnbondingDelegation creates a new UnbondingDelegation instance representing an
// unbonding delegation to a pool
func NewPoolUnbondingDelegation(
	delegatorAddress string, poolID uint32,
	creationHeight int64, completionTime time.Time, balance sdk.Coins, id uint64,
) UnbondingDelegation {
	return UnbondingDelegation{
		Type:             DELEGATION_TYPE_POOL,
		DelegatorAddress: delegatorAddress,
		TargetID:         poolID,
		Entries: []UnbondingDelegationEntry{
			NewUnbondingDelegationEntry(creationHeight, completionTime, balance, id),
		},
	}
}

// NewOperatorUnbondingDelegation creates a new UnbondingDelegation instance representing an
// unbonding delegation to an operator
func NewOperatorUnbondingDelegation(
	delegatorAddress string, operatorID uint32,
	creationHeight int64, completionTime time.Time, balance sdk.Coins, id uint64,
) UnbondingDelegation {
	return UnbondingDelegation{
		Type:             DELEGATION_TYPE_OPERATOR,
		DelegatorAddress: delegatorAddress,
		TargetID:         operatorID,
		Entries: []UnbondingDelegationEntry{
			NewUnbondingDelegationEntry(creationHeight, completionTime, balance, id),
		},
	}
}

// NewServiceUnbondingDelegation creates a new UnbondingDelegation instance representing an
// unbonding delegation to a service
func NewServiceUnbondingDelegation(
	delegatorAddress string, serviceID uint32,
	creationHeight int64, completionTime time.Time, balance sdk.Coins, id uint64,
) UnbondingDelegation {
	return UnbondingDelegation{
		Type:             DELEGATION_TYPE_SERVICE,
		DelegatorAddress: delegatorAddress,
		TargetID:         serviceID,
		Entries: []UnbondingDelegationEntry{
			NewUnbondingDelegationEntry(creationHeight, completionTime, balance, id),
		},
	}
}

// AddEntry allows to append a new entry to the unbonding delegation
func (ubd *UnbondingDelegation) AddEntry(creationHeight int64, completionTime time.Time, balance sdk.Coins, unbondingID uint64) bool {
	// Check the entries exists with creation_height and complete_time
	entryIndex := sort.Search(len(ubd.Entries), func(i int) bool {
		return ubd.Entries[i].CreationHeight == creationHeight && ubd.Entries[i].CompletionTime.Equal(completionTime)
	})

	// entryIndex exists
	if entryIndex != len(ubd.Entries) {
		ubdEntry := ubd.Entries[entryIndex]
		ubdEntry.Balance = ubdEntry.Balance.Add(balance...)
		ubdEntry.InitialBalance = ubdEntry.InitialBalance.Add(balance...)

		// Update the entry
		ubd.Entries[entryIndex] = ubdEntry
		return false
	}

	// Append the new unbonding delegation entry
	entry := NewUnbondingDelegationEntry(creationHeight, completionTime, balance, unbondingID)
	ubd.Entries = append(ubd.Entries, entry)
	return true
}

// RemoveEntry removes the entry at index i from the unbonding delegation
func (ubd *UnbondingDelegation) RemoveEntry(i int64) {
	ubd.Entries = append(ubd.Entries[:i], ubd.Entries[i+1:]...)
}

// Validate validates the unbonding delegation
func (ubd UnbondingDelegation) Validate() error {
	if ubd.Type == DELEGATION_TYPE_UNSPECIFIED {
		return fmt.Errorf("invalid unbonding delegation type")
	}

	if ubd.TargetID == 0 {
		return fmt.Errorf("invalid target id")
	}

	_, err := sdk.AccAddressFromBech32(ubd.DelegatorAddress)
	if err != nil {
		return fmt.Errorf("invalid delegator address")
	}

	if len(ubd.Entries) == 0 {
		return fmt.Errorf("empty entries")
	}

	for _, entry := range ubd.Entries {
		err = entry.Validate()
		if err != nil {
			return err
		}
	}

	return nil
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

// --------------------------------------------------------------------------------------------------------------------

// NewEmptyOperatorJoinedServices creates a new empty OperatorJoinedServices
// instance.
func NewEmptyOperatorJoinedServices() OperatorJoinedServices {
	return OperatorJoinedServices{
		ServiceIDs: nil,
	}
}

// NewOperatorJoinedServices creates a new OperatorJoinedServices instance.
func NewOperatorJoinedServices(serviceIDs []uint32) OperatorJoinedServices {
	return OperatorJoinedServices{
		ServiceIDs: serviceIDs,
	}
}

// Validate checks if the OperatorJoinedServices instance is valid.
func (o *OperatorJoinedServices) Validate() error {
	duplicate := utils.FindDuplicate(o.ServiceIDs)
	if duplicate != nil {
		return fmt.Errorf("duplicated service id: %d", *duplicate)
	}
	if utils.Contains(o.ServiceIDs, 0) {
		return fmt.Errorf("the service id cannot be 0]")
	}

	return nil
}

// UnsafeAdd adds the serviceID to the list of joined services without
// performing any check on the validity of the provided value and on the
// resulting validity of OperatorJoinedServices.
func (o *OperatorJoinedServices) UnsafeAdd(serviceID uint32) {
	o.ServiceIDs = append(o.ServiceIDs, serviceID)
}

// Contains returns true if the serviceID is already present.
func (o *OperatorJoinedServices) Contains(serviceID uint32) bool {
	return utils.Contains(o.ServiceIDs, serviceID)
}

// Add adds the serviceID to the list of joined services and returns
// true if the serviceID was added, false if was already in the list.
func (o *OperatorJoinedServices) Add(serviceID uint32) error {
	if serviceID == 0 {
		return fmt.Errorf("service id cannot be 0")
	}
	if o.Contains(serviceID) {
		return fmt.Errorf("service id %d already present", serviceID)
	}
	o.UnsafeAdd(serviceID)
	return nil
}

// ParseOperatorID tries parsing the given value as an service id
func ParseServiceID(value string) (uint32, error) {
	operatorID, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid service ID: %s", value)
	}
	return uint32(operatorID), nil
}
