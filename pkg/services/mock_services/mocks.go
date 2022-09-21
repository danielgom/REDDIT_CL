// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/services/services.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	internal "RD-Clone-API/pkg/internal"
	errors "RD-Clone-API/pkg/util/errors"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserService is a mock of UserService interface.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// Login mocks base method.
func (m *MockUserService) Login(arg0 context.Context, arg1 *internal.LoginRequest) (*internal.LoginResponse, errors.CommonError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0, arg1)
	ret0, _ := ret[0].(*internal.LoginResponse)
	ret1, _ := ret[1].(errors.CommonError)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockUserServiceMockRecorder) Login(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserService)(nil).Login), arg0, arg1)
}

// SignUp mocks base method.
func (m *MockUserService) SignUp(arg0 context.Context, arg1 *internal.RegisterRequest) (*internal.RegisterResponse, errors.CommonError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", arg0, arg1)
	ret0, _ := ret[0].(*internal.RegisterResponse)
	ret1, _ := ret[1].(errors.CommonError)
	return ret0, ret1
}

// SignUp indicates an expected call of SignUp.
func (mr *MockUserServiceMockRecorder) SignUp(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockUserService)(nil).SignUp), arg0, arg1)
}

// VerifyAccount mocks base method.
func (m *MockUserService) VerifyAccount(arg0 context.Context, arg1 string) errors.CommonError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyAccount", arg0, arg1)
	ret0, _ := ret[0].(errors.CommonError)
	return ret0
}

// VerifyAccount indicates an expected call of VerifyAccount.
func (mr *MockUserServiceMockRecorder) VerifyAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyAccount", reflect.TypeOf((*MockUserService)(nil).VerifyAccount), arg0, arg1)
}

// MockRefreshTokenService is a mock of RefreshTokenService interface.
type MockRefreshTokenService struct {
	ctrl     *gomock.Controller
	recorder *MockRefreshTokenServiceMockRecorder
}

// MockRefreshTokenServiceMockRecorder is the mock recorder for MockRefreshTokenService.
type MockRefreshTokenServiceMockRecorder struct {
	mock *MockRefreshTokenService
}

// NewMockRefreshTokenService creates a new mock instance.
func NewMockRefreshTokenService(ctrl *gomock.Controller) *MockRefreshTokenService {
	mock := &MockRefreshTokenService{ctrl: ctrl}
	mock.recorder = &MockRefreshTokenServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRefreshTokenService) EXPECT() *MockRefreshTokenServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRefreshTokenService) Create(arg0 context.Context) (string, errors.CommonError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(errors.CommonError)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRefreshTokenServiceMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRefreshTokenService)(nil).Create), arg0)
}

// Validate mocks base method.
func (m *MockRefreshTokenService) Validate(arg0 context.Context, arg1 string) errors.CommonError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", arg0, arg1)
	ret0, _ := ret[0].(errors.CommonError)
	return ret0
}

// Validate indicates an expected call of Validate.
func (mr *MockRefreshTokenServiceMockRecorder) Validate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockRefreshTokenService)(nil).Validate), arg0, arg1)
}
