package testutils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

var _ types.OperatorsHooks = &MockHooks{}

type MockHooks struct {
	CalledMap map[string]bool
}

func NewMockHooks() *MockHooks {
	return &MockHooks{CalledMap: make(map[string]bool)}
}

func (m MockHooks) AfterOperatorRegistered(ctx sdk.Context, operatorID uint32) error {
	m.CalledMap["AfterOperatorRegistered"] = true
	return nil
}

func (m MockHooks) AfterOperatorInactivatingStarted(ctx sdk.Context, operatorID uint32) error {
	m.CalledMap["AfterOperatorInactivatingStarted"] = true
	return nil
}

func (m MockHooks) AfterOperatorInactivatingCompleted(ctx sdk.Context, operatorID uint32) error {
	m.CalledMap["AfterOperatorInactivatingCompleted"] = true
	return nil
}

func (m MockHooks) AfterOperatorDeleted(ctx sdk.Context, operatorID uint32) error {
	m.CalledMap["AfterOperatorDeleted"] = true
	return nil
}

func (m MockHooks) AfterOperatorReactivated(ctx sdk.Context, operatorID uint32) error {
	m.CalledMap["AfterOperatorReactivated"] = true
	return nil
}
