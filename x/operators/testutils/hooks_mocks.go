package testutils

import (
	"context"

	"github.com/milkyway-labs/milkyway/v2/x/operators/types"
)

var _ types.OperatorsHooks = &MockHooks{}

type MockHooks struct {
	CalledMap map[string]bool
}

func NewMockHooks() *MockHooks {
	return &MockHooks{CalledMap: make(map[string]bool)}
}

func (m MockHooks) AfterOperatorRegistered(ctx context.Context, operatorID uint32) error {
	m.CalledMap["AfterOperatorRegistered"] = true
	return nil
}

func (m MockHooks) AfterOperatorInactivatingStarted(ctx context.Context, operatorID uint32) error {
	m.CalledMap["AfterOperatorInactivatingStarted"] = true
	return nil
}

func (m MockHooks) AfterOperatorInactivatingCompleted(ctx context.Context, operatorID uint32) error {
	m.CalledMap["AfterOperatorInactivatingCompleted"] = true
	return nil
}

func (m MockHooks) BeforeOperatorDeleted(ctx context.Context, operatorID uint32) error {
	m.CalledMap["BeforeOperatorDeleted"] = true
	return nil
}

func (m MockHooks) AfterOperatorReactivated(ctx context.Context, operatorID uint32) error {
	m.CalledMap["AfterOperatorReactivated"] = true
	return nil
}
