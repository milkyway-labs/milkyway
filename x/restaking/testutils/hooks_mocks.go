package testutils

import (
	"context"

	"github.com/milkyway-labs/milkyway/v9/x/restaking/types"
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
func (m MockHooks) BeforePoolDelegationCreated(context.Context, uint32, string) error {
	m.CalledMap["BeforePoolDelegationCreated"] = true
	return nil
}

// BeforePoolDelegationSharesModified implements restakingtypes.Hooks
func (m MockHooks) BeforePoolDelegationSharesModified(context.Context, uint32, string) error {
	m.CalledMap["BeforePoolDelegationSharesModified"] = true
	return nil
}

// AfterPoolDelegationModified implements restakingtypes.Hooks
func (m MockHooks) AfterPoolDelegationModified(context.Context, uint32, string) error {
	m.CalledMap["AfterPoolDelegationModified"] = true
	return nil
}

// BeforeOperatorDelegationCreated implements restakingtypes.Hooks
func (m MockHooks) BeforeOperatorDelegationCreated(context.Context, uint32, string) error {
	m.CalledMap["BeforeOperatorDelegationCreated"] = true
	return nil
}

// BeforeOperatorDelegationSharesModified implements restakingtypes.Hooks
func (m MockHooks) BeforeOperatorDelegationSharesModified(context.Context, uint32, string) error {
	m.CalledMap["BeforeOperatorDelegationSharesModified"] = true
	return nil
}

// AfterOperatorDelegationModified implements restakingtypes.Hooks
func (m MockHooks) AfterOperatorDelegationModified(context.Context, uint32, string) error {
	m.CalledMap["AfterOperatorDelegationModified"] = true
	return nil
}

// BeforeServiceDelegationCreated implements restakingtypes.Hooks
func (m MockHooks) BeforeServiceDelegationCreated(context.Context, uint32, string) error {
	m.CalledMap["BeforeServiceDelegationCreated"] = true
	return nil
}

// BeforeServiceDelegationSharesModified implements restakingtypes.Hooks
func (m MockHooks) BeforeServiceDelegationSharesModified(context.Context, uint32, string) error {
	m.CalledMap["BeforeServiceDelegationSharesModified"] = true
	return nil
}

// AfterServiceDelegationModified implements restakingtypes.Hooks
func (m MockHooks) AfterServiceDelegationModified(context.Context, uint32, string) error {
	m.CalledMap["AfterServiceDelegationModified"] = true
	return nil
}

// BeforePoolDelegationRemoved implements restakingtypes.Hooks
func (m MockHooks) BeforePoolDelegationRemoved(ctx context.Context, poolID uint32, delegator string) error {
	m.CalledMap["BeforePoolDelegationRemoved"] = true
	return nil
}

// BeforeOperatorDelegationRemoved implements restakingtypes.Hooks
func (m MockHooks) BeforeOperatorDelegationRemoved(ctx context.Context, operatorID uint32, delegator string) error {
	m.CalledMap["BeforeOperatorDelegationRemoved"] = true
	return nil
}

// BeforeServiceDelegationRemoved implements restakingtypes.Hooks
func (m MockHooks) BeforeServiceDelegationRemoved(ctx context.Context, serviceID uint32, delegator string) error {
	m.CalledMap["BeforeServiceDelegationRemoved"] = true
	return nil
}

// AfterUnbondingInitiated implements restakingtypes.Hooks
func (m MockHooks) AfterUnbondingInitiated(ctx context.Context, unbondingDelegationID uint64) error {
	m.CalledMap["AfterUnbondingInitiated"] = true
	return nil
}

// AfterUserPreferencesModified implements restakingtypes.Hooks
func (m MockHooks) AfterUserPreferencesModified(ctx context.Context, userAddress string, oldPreferences, newPreferences types.UserPreferences) error {
	m.CalledMap["AfterUserPreferencesModified"] = true
	return nil
}
