package testutils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/services/types"
)

var _ types.ServicesHooks = &MockHooks{}

type MockHooks struct {
	CalledMap map[string]bool
}

func NewMockHooks() *MockHooks {
	return &MockHooks{CalledMap: make(map[string]bool)}
}

func (m MockHooks) AfterServiceCreated(_ sdk.Context, _ uint32) error {
	m.CalledMap["AfterServiceCreated"] = true
	return nil
}

func (m MockHooks) AfterServiceActivated(_ sdk.Context, _ uint32) error {
	m.CalledMap["AfterServiceActivated"] = true
	return nil
}

func (m MockHooks) AfterServiceDeactivated(_ sdk.Context, _ uint32) error {
	m.CalledMap["AfterServiceDeactivated"] = true
	return nil
}

func (m MockHooks) AfterServiceDeleted(_ sdk.Context, _ uint32) error {
	m.CalledMap["AfterServiceDeleted"] = true
	return nil
}

func (m MockHooks) AfterServiceAccreditationModified(_ sdk.Context, _ uint32, _ bool) error {
	m.CalledMap["AfterServiceAccreditationModified"] = true
	return nil
}
