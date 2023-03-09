// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	models "github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockCreatorUsecase is a mock of CreatorUsecase interface.
type MockCreatorUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockCreatorUsecaseMockRecorder
}

// MockCreatorUsecaseMockRecorder is the mock recorder for MockCreatorUsecase.
type MockCreatorUsecaseMockRecorder struct {
	mock *MockCreatorUsecase
}

// NewMockCreatorUsecase creates a new mock instance.
func NewMockCreatorUsecase(ctrl *gomock.Controller) *MockCreatorUsecase {
	mock := &MockCreatorUsecase{ctrl: ctrl}
	mock.recorder = &MockCreatorUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCreatorUsecase) EXPECT() *MockCreatorUsecaseMockRecorder {
	return m.recorder
}

// GetPage mocks base method.
func (m *MockCreatorUsecase) GetPage(details *models.AccessDetails, creatorUUID string) (models.CreatorPage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPage", details, creatorUUID)
	ret0, _ := ret[0].(models.CreatorPage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPage indicates an expected call of GetPage.
func (mr *MockCreatorUsecaseMockRecorder) GetPage(details, creatorUUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPage", reflect.TypeOf((*MockCreatorUsecase)(nil).GetPage), details, creatorUUID)
}

// MockCreatorRepo is a mock of CreatorRepo interface.
type MockCreatorRepo struct {
	ctrl     *gomock.Controller
	recorder *MockCreatorRepoMockRecorder
}

// MockCreatorRepoMockRecorder is the mock recorder for MockCreatorRepo.
type MockCreatorRepoMockRecorder struct {
	mock *MockCreatorRepo
}

// NewMockCreatorRepo creates a new mock instance.
func NewMockCreatorRepo(ctrl *gomock.Controller) *MockCreatorRepo {
	mock := &MockCreatorRepo{ctrl: ctrl}
	mock.recorder = &MockCreatorRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCreatorRepo) EXPECT() *MockCreatorRepoMockRecorder {
	return m.recorder
}

// GetPage mocks base method.
func (m *MockCreatorRepo) GetPage(userId, creatorId uuid.UUID) (models.CreatorPage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPage", userId, creatorId)
	ret0, _ := ret[0].(models.CreatorPage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPage indicates an expected call of GetPage.
func (mr *MockCreatorRepoMockRecorder) GetPage(userId, creatorId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPage", reflect.TypeOf((*MockCreatorRepo)(nil).GetPage), userId, creatorId)
}
