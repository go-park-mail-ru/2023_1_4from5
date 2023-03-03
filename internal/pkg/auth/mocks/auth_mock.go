// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	models "github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockAuthUsecase is a mock of AuthUsecase interface.
type MockAuthUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockAuthUsecaseMockRecorder
}

// MockAuthUsecaseMockRecorder is the mock recorder for MockAuthUsecase.
type MockAuthUsecaseMockRecorder struct {
	mock *MockAuthUsecase
}

// NewMockAuthUsecase creates a new mock instance.
func NewMockAuthUsecase(ctrl *gomock.Controller) *MockAuthUsecase {
	mock := &MockAuthUsecase{ctrl: ctrl}
	mock.recorder = &MockAuthUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthUsecase) EXPECT() *MockAuthUsecaseMockRecorder {
	return m.recorder
}

// SignIn mocks base method.
func (m *MockAuthUsecase) SignIn(user models.LoginUser) (string, int) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIn", user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(int)
	return ret0, ret1
}

// SignIn indicates an expected call of SignIn.
func (mr *MockAuthUsecaseMockRecorder) SignIn(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIn", reflect.TypeOf((*MockAuthUsecase)(nil).SignIn), user)
}

// SignUp mocks base method.
func (m *MockAuthUsecase) SignUp(user models.User) (string, int) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(int)
	return ret0, ret1
}

// SignUp indicates an expected call of SignUp.
func (mr *MockAuthUsecaseMockRecorder) SignUp(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockAuthUsecase)(nil).SignUp), user)
}

// MockAuthRepo is a mock of AuthRepo interface.
type MockAuthRepo struct {
	ctrl     *gomock.Controller
	recorder *MockAuthRepoMockRecorder
}

// MockAuthRepoMockRecorder is the mock recorder for MockAuthRepo.
type MockAuthRepoMockRecorder struct {
	mock *MockAuthRepo
}

// NewMockAuthRepo creates a new mock instance.
func NewMockAuthRepo(ctrl *gomock.Controller) *MockAuthRepo {
	mock := &MockAuthRepo{ctrl: ctrl}
	mock.recorder = &MockAuthRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthRepo) EXPECT() *MockAuthRepoMockRecorder {
	return m.recorder
}

// CheckUser mocks base method.
func (m *MockAuthRepo) CheckUser(user models.User) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUser", user)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUser indicates an expected call of CheckUser.
func (mr *MockAuthRepoMockRecorder) CheckUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUser", reflect.TypeOf((*MockAuthRepo)(nil).CheckUser), user)
}

// CreateUser mocks base method.
func (m *MockAuthRepo) CreateUser(user models.User) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", user)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockAuthRepoMockRecorder) CreateUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockAuthRepo)(nil).CreateUser), user)
}

// MockTokenGenerator is a mock of TokenGenerator interface.
type MockTokenGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockTokenGeneratorMockRecorder
}

// MockTokenGeneratorMockRecorder is the mock recorder for MockTokenGenerator.
type MockTokenGeneratorMockRecorder struct {
	mock *MockTokenGenerator
}

// NewMockTokenGenerator creates a new mock instance.
func NewMockTokenGenerator(ctrl *gomock.Controller) *MockTokenGenerator {
	mock := &MockTokenGenerator{ctrl: ctrl}
	mock.recorder = &MockTokenGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenGenerator) EXPECT() *MockTokenGeneratorMockRecorder {
	return m.recorder
}

// GetToken mocks base method.
func (m *MockTokenGenerator) GetToken(user models.User) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetToken", user)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetToken indicates an expected call of GetToken.
func (mr *MockTokenGeneratorMockRecorder) GetToken(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetToken", reflect.TypeOf((*MockTokenGenerator)(nil).GetToken), user)
}

// MockEncrypter is a mock of Encrypter interface.
type MockEncrypter struct {
	ctrl     *gomock.Controller
	recorder *MockEncrypterMockRecorder
}

// MockEncrypterMockRecorder is the mock recorder for MockEncrypter.
type MockEncrypterMockRecorder struct {
	mock *MockEncrypter
}

// NewMockEncrypter creates a new mock instance.
func NewMockEncrypter(ctrl *gomock.Controller) *MockEncrypter {
	mock := &MockEncrypter{ctrl: ctrl}
	mock.recorder = &MockEncrypterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEncrypter) EXPECT() *MockEncrypterMockRecorder {
	return m.recorder
}

// EncryptPswd mocks base method.
func (m *MockEncrypter) EncryptPswd(pswd string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EncryptPswd", pswd)
	ret0, _ := ret[0].(string)
	return ret0
}

// EncryptPswd indicates an expected call of EncryptPswd.
func (mr *MockEncrypterMockRecorder) EncryptPswd(pswd interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EncryptPswd", reflect.TypeOf((*MockEncrypter)(nil).EncryptPswd), pswd)
}
