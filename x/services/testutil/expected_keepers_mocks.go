// Code generated by MockGen. DO NOT EDIT.
// Source: ./x/services/keeper/expected_keepers.go
//
// Generated by this command:
//
//	mockgen -source ./x/services/keeper/expected_keepers.go -package testutil -destination ./x/services/testutil/expected_keepers_mocks.go
//

// Package testutil is a generated GoMock package.
package testutil

import (
	context "context"
	reflect "reflect"

	types "github.com/cosmos/cosmos-sdk/types"
	gomock "go.uber.org/mock/gomock"
)

// MockCommunityPoolKeeper is a mock of CommunityPoolKeeper interface.
type MockCommunityPoolKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockCommunityPoolKeeperMockRecorder
}

// MockCommunityPoolKeeperMockRecorder is the mock recorder for MockCommunityPoolKeeper.
type MockCommunityPoolKeeperMockRecorder struct {
	mock *MockCommunityPoolKeeper
}

// NewMockCommunityPoolKeeper creates a new mock instance.
func NewMockCommunityPoolKeeper(ctrl *gomock.Controller) *MockCommunityPoolKeeper {
	mock := &MockCommunityPoolKeeper{ctrl: ctrl}
	mock.recorder = &MockCommunityPoolKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommunityPoolKeeper) EXPECT() *MockCommunityPoolKeeperMockRecorder {
	return m.recorder
}

// FundCommunityPool mocks base method.
func (m *MockCommunityPoolKeeper) FundCommunityPool(ctx context.Context, amount types.Coins, sender types.AccAddress) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FundCommunityPool", ctx, amount, sender)
	ret0, _ := ret[0].(error)
	return ret0
}

// FundCommunityPool indicates an expected call of FundCommunityPool.
func (mr *MockCommunityPoolKeeperMockRecorder) FundCommunityPool(ctx, amount, sender any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FundCommunityPool", reflect.TypeOf((*MockCommunityPoolKeeper)(nil).FundCommunityPool), ctx, amount, sender)
}