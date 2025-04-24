package types

import (
	"fmt"
	"slices"
	"sort"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v11/utils"
	operatorstypes "github.com/milkyway-labs/milkyway/v11/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/v11/x/pools/types"
	servicestypes "github.com/milkyway-labs/milkyway/v11/x/services/types"
)

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

// GetDelegationTypeFromTarget returns the delegation type based on the target
func GetDelegationTypeFromTarget(target DelegationTarget) (DelegationType, error) {
	switch target.(type) {
	case poolstypes.Pool:
		return DELEGATION_TYPE_POOL, nil
	case operatorstypes.Operator:
		return DELEGATION_TYPE_OPERATOR, nil
	case servicestypes.Service:
		return DELEGATION_TYPE_SERVICE, nil
	default:
		return DELEGATION_TYPE_UNSPECIFIED, fmt.Errorf("invalid delegation target type : %T", target)
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

// NewUserPreferences creates a new UserPreferences instance
func NewUserPreferences(trustedServices []TrustedServiceEntry) UserPreferences {
	return UserPreferences{
		TrustedServices: trustedServices,
	}
}

// DefaultUserPreferences returns the default user preferences
func DefaultUserPreferences() UserPreferences {
	return NewUserPreferences(nil)
}

// Validate validates the user preferences
func (p UserPreferences) Validate() error {
	var trustedServicesIDs []uint32
	for _, entry := range p.TrustedServices {
		// Validate the entry
		if err := entry.Validate(); err != nil {
			return err
		}

		// Make sure there are no other entries for the same service id
		if slices.Contains(trustedServicesIDs, entry.ServiceID) {
			return fmt.Errorf("duplicate service id: %d", entry.ServiceID)
		}

		trustedServicesIDs = append(trustedServicesIDs, entry.ServiceID)
	}

	return nil
}

// GetTrustedServiceEntry returns the trusted service entry for the given service id
func (p UserPreferences) GetTrustedServiceEntry(serviceID uint32) TrustedServiceEntry {
	for _, entry := range p.TrustedServices {
		if entry.ServiceID == serviceID {
			return entry
		}
	}
	return TrustedServiceEntry{}
}

// TrustedServicesIDs returns the list of service ids that the user trusts
func (p UserPreferences) TrustedServicesIDs() []uint32 {
	return utils.Map(p.TrustedServices, func(entry TrustedServiceEntry) uint32 { return entry.ServiceID })
}

// IsServiceTrusted tells whether the user trusts the given service
func (p UserPreferences) IsServiceTrusted(serviceID uint32) bool {
	return slices.Contains(p.TrustedServicesIDs(), serviceID)
}

// IsServiceTrustedWithPool tells whether the user trusts the given service with the given pool
func (p UserPreferences) IsServiceTrustedWithPool(serviceID, poolID uint32) bool {
	if len(p.TrustedServices) == 0 {
		return false
	}
	entry := p.GetTrustedServiceEntry(serviceID)
	return entry.ServiceID != 0 && (len(entry.PoolsIDs) == 0 || slices.Contains(entry.PoolsIDs, poolID))
}

// ComputeChangedServicesIDs returns the list of service IDs that have changed
// between the two user preferences.
// A service is considered changed in the following cases:
// - It was trusted in the old preferences but not in the new preferences
// - It was not trusted in the old preferences but it is in the new preferences
// - It was trusted in both preferences but the pools have changed
//
// The result will contain the all services ids that have changed
// (either removed,added or modified).
// CONTRACT: this function is never triggered in Begin/End block or any proposal execution
func ComputeChangedServicesIDs(before UserPreferences, after UserPreferences) (changed []uint32) {
	beforeTrusted := before.TrustedServicesIDs()
	afterTrusted := after.TrustedServicesIDs()

	// Get the deleted entries (A - B)
	changed = append(changed, utils.Difference(beforeTrusted, afterTrusted)...)

	// Get the added entries (B - A)
	changed = append(changed, utils.Difference(afterTrusted, beforeTrusted)...)

	// Find the list of services that have changed pools (A ∩ B)
	for _, serviceID := range utils.Intersect(beforeTrusted, afterTrusted) {
		beforePools := before.GetTrustedServiceEntry(serviceID).PoolsIDs
		afterPools := after.GetTrustedServiceEntry(serviceID).PoolsIDs

		if !slices.Equal(beforePools, afterPools) {
			changed = append(changed, serviceID)
		}
	}

	return changed
}

// NewTrustedServiceEntry creates a new TrustedServiceEntry instance
func NewTrustedServiceEntry(serviceID uint32, poolsIDs []uint32) TrustedServiceEntry {
	return TrustedServiceEntry{
		ServiceID: serviceID,
		PoolsIDs:  poolsIDs,
	}
}

// Validate validates the trusted service entry
func (e TrustedServiceEntry) Validate() error {
	if e.ServiceID == 0 {
		return fmt.Errorf("invalid service id")
	}

	for _, poolID := range e.PoolsIDs {
		if poolID == 0 {
			return fmt.Errorf("invalid pool id")
		}
	}

	duplicated := utils.FindDuplicate(e.PoolsIDs)
	if duplicated != nil {
		return fmt.Errorf("duplicated pool id: %d", *duplicated)
	}

	return nil
}
