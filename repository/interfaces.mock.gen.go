// Code generated by MockGen. DO NOT EDIT.
// Source: repository/interfaces.go

// Package repository is a generated GoMock package.
package repository

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRepositoryInterface is a mock of RepositoryInterface interface.
type MockRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryInterfaceMockRecorder
}

// MockRepositoryInterfaceMockRecorder is the mock recorder for MockRepositoryInterface.
type MockRepositoryInterfaceMockRecorder struct {
	mock *MockRepositoryInterface
}

// NewMockRepositoryInterface creates a new mock instance.
func NewMockRepositoryInterface(ctrl *gomock.Controller) *MockRepositoryInterface {
	mock := &MockRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepositoryInterface) EXPECT() *MockRepositoryInterfaceMockRecorder {
	return m.recorder
}

// GetMyProfile mocks base method.
func (m *MockRepositoryInterface) GetMyProfile(ctx context.Context, input GetMyProfileInput) (GetMyProfileOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMyProfile", ctx, input)
	ret0, _ := ret[0].(GetMyProfileOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMyProfile indicates an expected call of GetMyProfile.
func (mr *MockRepositoryInterfaceMockRecorder) GetMyProfile(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMyProfile", reflect.TypeOf((*MockRepositoryInterface)(nil).GetMyProfile), ctx, input)
}

// Login mocks base method.
func (m *MockRepositoryInterface) Login(ctx context.Context, input LoginInput) (LoginOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, input)
	ret0, _ := ret[0].(LoginOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockRepositoryInterfaceMockRecorder) Login(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockRepositoryInterface)(nil).Login), ctx, input)
}

// Registration mocks base method.
func (m *MockRepositoryInterface) Registration(ctx context.Context, input RegistrationInput) (RegistrationOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Registration", ctx, input)
	ret0, _ := ret[0].(RegistrationOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Registration indicates an expected call of Registration.
func (mr *MockRepositoryInterfaceMockRecorder) Registration(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Registration", reflect.TypeOf((*MockRepositoryInterface)(nil).Registration), ctx, input)
}

// UpdateMyProfile mocks base method.
func (m *MockRepositoryInterface) UpdateMyProfile(ctx context.Context, input UpdateMyProfileInput) (UpdateMyProfileOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMyProfile", ctx, input)
	ret0, _ := ret[0].(UpdateMyProfileOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateMyProfile indicates an expected call of UpdateMyProfile.
func (mr *MockRepositoryInterfaceMockRecorder) UpdateMyProfile(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMyProfile", reflect.TypeOf((*MockRepositoryInterface)(nil).UpdateMyProfile), ctx, input)
}
