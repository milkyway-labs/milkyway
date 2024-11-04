package testutils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

var _ types.RestakingHooks = &MockHooks{}

// MockHooks is a mock implementation of the RestakingHooks interface
type MockHooks struct {
	CalledMap map[string]bool
}

// NewMockHooks returns a new MockHooks
func NewMockHooks() *MockHooks {
	return &MockHooks{CalledMap: make(map[string]bool)}
}

// BeforePoolDelegationCreated implements restakingtypes.Hooks
func (m MockHooks) BeforePoolDelegationCreated(sdk.Context, uint32, string) error {
	m.CalledMap["BeforePoolDelegationCreated"] = true
	return nil
}

// BeforePoolDelegationSharesModified implements restakingtypes.Hooks
func (m MockHooks) BeforePoolDelegationSharesModified(sdk.Context, uint32, string) error {
	m.CalledMap["BeforePoolDelegationSharesModified"] = true
	return nil
}

// AfterPoolDelegationModified implements restakingtypes.Hooks
func (m MockHooks) AfterPoolDelegationModified(sdk.Context, uint32, string) error {
	m.CalledMap["AfterPoolDelegationModified"] = true
	return nil
}

// BeforeOperatorDelegationCreated implements restakingtypes.Hooks
func (m MockHooks) BeforeOperatorDelegationCreated(sdk.Context, uint32, string) error {
	m.CalledMap["BeforeOperatorDelegationCreated"] = true
	return nil
}

// BeforeOperatorDelegationSharesModified implements restakingtypes.Hooks
func (m MockHooks) BeforeOperatorDelegationSharesModified(sdk.Context, uint32, string) error {
	m.CalledMap["BeforeOperatorDelegationSharesModified"] = true
	return nil
}

// AfterOperatorDelegationModified implements restakingtypes.Hooks
func (m MockHooks) AfterOperatorDelegationModified(sdk.Context, uint32, string) error {
	m.CalledMap["AfterOperatorDelegationModified"] = true
	return nil
}

// BeforeServiceDelegationCreated implements restakingtypes.Hooks
func (m MockHooks) BeforeServiceDelegationCreated(sdk.Context, uint32, string) error {
	m.CalledMap["BeforeServiceDelegationCreated"] = true
	return nil
}

// BeforeServiceDelegationSharesModified implements restakingtypes.Hooks
func (m MockHooks) BeforeServiceDelegationSharesModified(sdk.Context, uint32, string) error {
	m.CalledMap["BeforeServiceDelegationSharesModified"] = true
	return nil
}

// AfterServiceDelegationModified implements restakingtypes.Hooks
func (m MockHooks) AfterServiceDelegationModified(sdk.Context, uint32, string) error {
	m.CalledMap["AfterServiceDelegationModified"] = true
	return nil
}

// BeforePoolDelegationRemoved implements restakingtypes.Hooks
func (m MockHooks) BeforePoolDelegationRemoved(ctx sdk.Context, poolID uint32, delegator string) error {
	m.CalledMap["BeforePoolDelegationRemoved"] = true
	return nil
}

// BeforeOperatorDelegationRemoved implements restakingtypes.Hooks
func (m MockHooks) BeforeOperatorDelegationRemoved(ctx sdk.Context, operatorID uint32, delegator string) error {
	m.CalledMap["BeforeOperatorDelegationRemoved"] = true
	return nil
}

// BeforeServiceDelegationRemoved implements restakingtypes.Hooks
func (m MockHooks) BeforeServiceDelegationRemoved(ctx sdk.Context, serviceID uint32, delegator string) error {
	m.CalledMap["BeforeServiceDelegationRemoved"] = true
	return nil
}

// AfterUnbondingInitiated implements restakingtypes.Hooks
func (m MockHooks) AfterUnbondingInitiated(ctx sdk.Context, unbondingDelegationID uint64) error {
	m.CalledMap["AfterUnbondingInitiated"] = true
	return nil
}
