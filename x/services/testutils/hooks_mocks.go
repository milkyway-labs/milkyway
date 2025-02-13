package testutils

import (
	"context"

	"github.com/milkyway-labs/milkyway/v9/x/services/types"
)

var _ types.ServicesHooks = &MockHooks{}

type MockHooks struct {
	CalledMap map[string]bool
}

func NewMockHooks() *MockHooks {
	return &MockHooks{CalledMap: make(map[string]bool)}
}

func (m MockHooks) AfterServiceCreated(_ context.Context, _ uint32) error {
	m.CalledMap["AfterServiceCreated"] = true
	return nil
}

func (m MockHooks) AfterServiceActivated(_ context.Context, _ uint32) error {
	m.CalledMap["AfterServiceActivated"] = true
	return nil
}

func (m MockHooks) AfterServiceDeactivated(_ context.Context, _ uint32) error {
	m.CalledMap["AfterServiceDeactivated"] = true
	return nil
}

func (m MockHooks) BeforeServiceDeleted(_ context.Context, _ uint32) error {
	m.CalledMap["BeforeServiceDeleted"] = true
	return nil
}

func (m MockHooks) AfterServiceAccreditationModified(_ context.Context, _ uint32) error {
	m.CalledMap["AfterServiceAccreditationModified"] = true
	return nil
}
