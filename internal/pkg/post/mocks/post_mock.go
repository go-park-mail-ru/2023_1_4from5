// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	models "github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockPostUsecase is a mock of PostUsecase interface.
type MockPostUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockPostUsecaseMockRecorder
}

// MockPostUsecaseMockRecorder is the mock recorder for MockPostUsecase.
type MockPostUsecaseMockRecorder struct {
	mock *MockPostUsecase
}

// NewMockPostUsecase creates a new mock instance.
func NewMockPostUsecase(ctrl *gomock.Controller) *MockPostUsecase {
	mock := &MockPostUsecase{ctrl: ctrl}
	mock.recorder = &MockPostUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPostUsecase) EXPECT() *MockPostUsecaseMockRecorder {
	return m.recorder
}

// AddLike mocks base method.
func (m *MockPostUsecase) AddLike(ctx context.Context, userID, postID uuid.UUID) (models.Like, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddLike", ctx, userID, postID)
	ret0, _ := ret[0].(models.Like)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddLike indicates an expected call of AddLike.
func (mr *MockPostUsecaseMockRecorder) AddLike(ctx, userID, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddLike", reflect.TypeOf((*MockPostUsecase)(nil).AddLike), ctx, userID, postID)
}

// CreatePost mocks base method.
func (m *MockPostUsecase) CreatePost(ctx context.Context, postData models.PostCreationData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePost", ctx, postData)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreatePost indicates an expected call of CreatePost.
func (mr *MockPostUsecaseMockRecorder) CreatePost(ctx, postData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePost", reflect.TypeOf((*MockPostUsecase)(nil).CreatePost), ctx, postData)
}

// DeletePost mocks base method.
func (m *MockPostUsecase) DeletePost(ctx context.Context, postID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePost", ctx, postID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePost indicates an expected call of DeletePost.
func (mr *MockPostUsecaseMockRecorder) DeletePost(ctx, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePost", reflect.TypeOf((*MockPostUsecase)(nil).DeletePost), ctx, postID)
}

// EditPost mocks base method.
func (m *MockPostUsecase) EditPost(ctx context.Context, postData models.PostEditData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditPost", ctx, postData)
	ret0, _ := ret[0].(error)
	return ret0
}

// EditPost indicates an expected call of EditPost.
func (mr *MockPostUsecaseMockRecorder) EditPost(ctx, postData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditPost", reflect.TypeOf((*MockPostUsecase)(nil).EditPost), ctx, postData)
}

// GetPost mocks base method.
func (m *MockPostUsecase) GetPost(ctx context.Context, postID, userID uuid.UUID) (models.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPost", ctx, postID, userID)
	ret0, _ := ret[0].(models.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPost indicates an expected call of GetPost.
func (mr *MockPostUsecaseMockRecorder) GetPost(ctx, postID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPost", reflect.TypeOf((*MockPostUsecase)(nil).GetPost), ctx, postID, userID)
}

// IsCreator mocks base method.
func (m *MockPostUsecase) IsCreator(ctx context.Context, userID, creatorID uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsCreator", ctx, userID, creatorID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsCreator indicates an expected call of IsCreator.
func (mr *MockPostUsecaseMockRecorder) IsCreator(ctx, userID, creatorID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsCreator", reflect.TypeOf((*MockPostUsecase)(nil).IsCreator), ctx, userID, creatorID)
}

// IsPostOwner mocks base method.
func (m *MockPostUsecase) IsPostOwner(ctx context.Context, userId, postId uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsPostOwner", ctx, userId, postId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsPostOwner indicates an expected call of IsPostOwner.
func (mr *MockPostUsecaseMockRecorder) IsPostOwner(ctx, userId, postId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsPostOwner", reflect.TypeOf((*MockPostUsecase)(nil).IsPostOwner), ctx, userId, postId)
}

// RemoveLike mocks base method.
func (m *MockPostUsecase) RemoveLike(ctx context.Context, userID, postID uuid.UUID) (models.Like, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveLike", ctx, userID, postID)
	ret0, _ := ret[0].(models.Like)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoveLike indicates an expected call of RemoveLike.
func (mr *MockPostUsecaseMockRecorder) RemoveLike(ctx, userID, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveLike", reflect.TypeOf((*MockPostUsecase)(nil).RemoveLike), ctx, userID, postID)
}

// MockPostRepo is a mock of PostRepo interface.
type MockPostRepo struct {
	ctrl     *gomock.Controller
	recorder *MockPostRepoMockRecorder
}

// MockPostRepoMockRecorder is the mock recorder for MockPostRepo.
type MockPostRepoMockRecorder struct {
	mock *MockPostRepo
}

// NewMockPostRepo creates a new mock instance.
func NewMockPostRepo(ctrl *gomock.Controller) *MockPostRepo {
	mock := &MockPostRepo{ctrl: ctrl}
	mock.recorder = &MockPostRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPostRepo) EXPECT() *MockPostRepoMockRecorder {
	return m.recorder
}

// AddLike mocks base method.
func (m *MockPostRepo) AddLike(ctx context.Context, userID, postID uuid.UUID) (models.Like, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddLike", ctx, userID, postID)
	ret0, _ := ret[0].(models.Like)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddLike indicates an expected call of AddLike.
func (mr *MockPostRepoMockRecorder) AddLike(ctx, userID, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddLike", reflect.TypeOf((*MockPostRepo)(nil).AddLike), ctx, userID, postID)
}

// CreatePost mocks base method.
func (m *MockPostRepo) CreatePost(ctx context.Context, postData models.PostCreationData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePost", ctx, postData)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreatePost indicates an expected call of CreatePost.
func (mr *MockPostRepoMockRecorder) CreatePost(ctx, postData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePost", reflect.TypeOf((*MockPostRepo)(nil).CreatePost), ctx, postData)
}

// DeletePost mocks base method.
func (m *MockPostRepo) DeletePost(ctx context.Context, postID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePost", ctx, postID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePost indicates an expected call of DeletePost.
func (mr *MockPostRepoMockRecorder) DeletePost(ctx, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePost", reflect.TypeOf((*MockPostRepo)(nil).DeletePost), ctx, postID)
}

// EditPost mocks base method.
func (m *MockPostRepo) EditPost(ctx context.Context, postData models.PostEditData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditPost", ctx, postData)
	ret0, _ := ret[0].(error)
	return ret0
}

// EditPost indicates an expected call of EditPost.
func (mr *MockPostRepoMockRecorder) EditPost(ctx, postData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditPost", reflect.TypeOf((*MockPostRepo)(nil).EditPost), ctx, postData)
}

// GetPost mocks base method.
func (m *MockPostRepo) GetPost(ctx context.Context, postID, userID uuid.UUID) (models.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPost", ctx, postID, userID)
	ret0, _ := ret[0].(models.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPost indicates an expected call of GetPost.
func (mr *MockPostRepoMockRecorder) GetPost(ctx, postID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPost", reflect.TypeOf((*MockPostRepo)(nil).GetPost), ctx, postID, userID)
}

// GetSubsByID mocks base method.
func (m *MockPostRepo) GetSubsByID(ctx context.Context, subsIDs ...uuid.UUID) ([]models.Subscription, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range subsIDs {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSubsByID", varargs...)
	ret0, _ := ret[0].([]models.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubsByID indicates an expected call of GetSubsByID.
func (mr *MockPostRepoMockRecorder) GetSubsByID(ctx interface{}, subsIDs ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, subsIDs...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubsByID", reflect.TypeOf((*MockPostRepo)(nil).GetSubsByID), varargs...)
}

// IsCreator mocks base method.
func (m *MockPostRepo) IsCreator(ctx context.Context, userID, creatorID uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsCreator", ctx, userID, creatorID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsCreator indicates an expected call of IsCreator.
func (mr *MockPostRepoMockRecorder) IsCreator(ctx, userID, creatorID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsCreator", reflect.TypeOf((*MockPostRepo)(nil).IsCreator), ctx, userID, creatorID)
}

// IsPostAvailable mocks base method.
func (m *MockPostRepo) IsPostAvailable(ctx context.Context, userID, postID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsPostAvailable", ctx, userID, postID)
	ret0, _ := ret[0].(error)
	return ret0
}

// IsPostAvailable indicates an expected call of IsPostAvailable.
func (mr *MockPostRepoMockRecorder) IsPostAvailable(ctx, userID, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsPostAvailable", reflect.TypeOf((*MockPostRepo)(nil).IsPostAvailable), ctx, userID, postID)
}

// IsPostOwner mocks base method.
func (m *MockPostRepo) IsPostOwner(ctx context.Context, userId, postId uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsPostOwner", ctx, userId, postId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsPostOwner indicates an expected call of IsPostOwner.
func (mr *MockPostRepoMockRecorder) IsPostOwner(ctx, userId, postId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsPostOwner", reflect.TypeOf((*MockPostRepo)(nil).IsPostOwner), ctx, userId, postId)
}

// RemoveLike mocks base method.
func (m *MockPostRepo) RemoveLike(ctx context.Context, userID, postID uuid.UUID) (models.Like, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveLike", ctx, userID, postID)
	ret0, _ := ret[0].(models.Like)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoveLike indicates an expected call of RemoveLike.
func (mr *MockPostRepoMockRecorder) RemoveLike(ctx, userID, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveLike", reflect.TypeOf((*MockPostRepo)(nil).RemoveLike), ctx, userID, postID)
}
