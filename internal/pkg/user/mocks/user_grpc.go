// Code generated by MockGen. DO NOT EDIT.
// Source: ./generated/user_grpc.pb.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	proto "github.com/go-park-mail-ru/2023_1_4from5/internal/models/proto"
	generated "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/delivery/grpc/generated"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockUserServiceClient is a mock of UserServiceClient interface.
type MockUserServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceClientMockRecorder
}

// MockUserServiceClientMockRecorder is the mock recorder for MockUserServiceClient.
type MockUserServiceClientMockRecorder struct {
	mock *MockUserServiceClient
}

// NewMockUserServiceClient creates a new mock instance.
func NewMockUserServiceClient(ctrl *gomock.Controller) *MockUserServiceClient {
	mock := &MockUserServiceClient{ctrl: ctrl}
	mock.recorder = &MockUserServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserServiceClient) EXPECT() *MockUserServiceClientMockRecorder {
	return m.recorder
}

// BecomeCreator mocks base method.
func (m *MockUserServiceClient) BecomeCreator(ctx context.Context, in *generated.BecameCreatorInfoMessage, opts ...grpc.CallOption) (*proto.UUIDResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "BecomeCreator", varargs...)
	ret0, _ := ret[0].(*proto.UUIDResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BecomeCreator indicates an expected call of BecomeCreator.
func (mr *MockUserServiceClientMockRecorder) BecomeCreator(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BecomeCreator", reflect.TypeOf((*MockUserServiceClient)(nil).BecomeCreator), varargs...)
}

// CheckIfCreator mocks base method.
func (m *MockUserServiceClient) CheckIfCreator(ctx context.Context, in *proto.UUIDMessage, opts ...grpc.CallOption) (*generated.CheckCreatorMessage, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CheckIfCreator", varargs...)
	ret0, _ := ret[0].(*generated.CheckCreatorMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckIfCreator indicates an expected call of CheckIfCreator.
func (mr *MockUserServiceClientMockRecorder) CheckIfCreator(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckIfCreator", reflect.TypeOf((*MockUserServiceClient)(nil).CheckIfCreator), varargs...)
}

// DeletePhoto mocks base method.
func (m *MockUserServiceClient) DeletePhoto(ctx context.Context, in *proto.UUIDMessage, opts ...grpc.CallOption) (*proto.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeletePhoto", varargs...)
	ret0, _ := ret[0].(*proto.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeletePhoto indicates an expected call of DeletePhoto.
func (mr *MockUserServiceClientMockRecorder) DeletePhoto(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePhoto", reflect.TypeOf((*MockUserServiceClient)(nil).DeletePhoto), varargs...)
}

// Donate mocks base method.
func (m *MockUserServiceClient) Donate(ctx context.Context, in *generated.DonateMessage, opts ...grpc.CallOption) (*generated.DonateResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Donate", varargs...)
	ret0, _ := ret[0].(*generated.DonateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Donate indicates an expected call of Donate.
func (mr *MockUserServiceClientMockRecorder) Donate(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Donate", reflect.TypeOf((*MockUserServiceClient)(nil).Donate), varargs...)
}

// Follow mocks base method.
func (m *MockUserServiceClient) Follow(ctx context.Context, in *generated.FollowMessage, opts ...grpc.CallOption) (*proto.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Follow", varargs...)
	ret0, _ := ret[0].(*proto.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Follow indicates an expected call of Follow.
func (mr *MockUserServiceClientMockRecorder) Follow(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Follow", reflect.TypeOf((*MockUserServiceClient)(nil).Follow), varargs...)
}

// GetProfile mocks base method.
func (m *MockUserServiceClient) GetProfile(ctx context.Context, in *proto.UUIDMessage, opts ...grpc.CallOption) (*generated.UserProfile, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetProfile", varargs...)
	ret0, _ := ret[0].(*generated.UserProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfile indicates an expected call of GetProfile.
func (mr *MockUserServiceClientMockRecorder) GetProfile(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfile", reflect.TypeOf((*MockUserServiceClient)(nil).GetProfile), varargs...)
}

// Subscribe mocks base method.
func (m *MockUserServiceClient) Subscribe(ctx context.Context, in *generated.SubscriptionDetails, opts ...grpc.CallOption) (*proto.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Subscribe", varargs...)
	ret0, _ := ret[0].(*proto.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockUserServiceClientMockRecorder) Subscribe(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockUserServiceClient)(nil).Subscribe), varargs...)
}

// Unfollow mocks base method.
func (m *MockUserServiceClient) Unfollow(ctx context.Context, in *generated.FollowMessage, opts ...grpc.CallOption) (*proto.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Unfollow", varargs...)
	ret0, _ := ret[0].(*proto.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Unfollow indicates an expected call of Unfollow.
func (mr *MockUserServiceClientMockRecorder) Unfollow(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unfollow", reflect.TypeOf((*MockUserServiceClient)(nil).Unfollow), varargs...)
}

// UpdatePassword mocks base method.
func (m *MockUserServiceClient) UpdatePassword(ctx context.Context, in *generated.UpdatePasswordMessage, opts ...grpc.CallOption) (*proto.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdatePassword", varargs...)
	ret0, _ := ret[0].(*proto.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePassword indicates an expected call of UpdatePassword.
func (mr *MockUserServiceClientMockRecorder) UpdatePassword(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePassword", reflect.TypeOf((*MockUserServiceClient)(nil).UpdatePassword), varargs...)
}

// UpdatePhoto mocks base method.
func (m *MockUserServiceClient) UpdatePhoto(ctx context.Context, in *proto.UUIDMessage, opts ...grpc.CallOption) (*generated.ImageID, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdatePhoto", varargs...)
	ret0, _ := ret[0].(*generated.ImageID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePhoto indicates an expected call of UpdatePhoto.
func (mr *MockUserServiceClientMockRecorder) UpdatePhoto(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePhoto", reflect.TypeOf((*MockUserServiceClient)(nil).UpdatePhoto), varargs...)
}

// UpdateProfileInfo mocks base method.
func (m *MockUserServiceClient) UpdateProfileInfo(ctx context.Context, in *generated.UpdateProfileInfoMessage, opts ...grpc.CallOption) (*proto.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateProfileInfo", varargs...)
	ret0, _ := ret[0].(*proto.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateProfileInfo indicates an expected call of UpdateProfileInfo.
func (mr *MockUserServiceClientMockRecorder) UpdateProfileInfo(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProfileInfo", reflect.TypeOf((*MockUserServiceClient)(nil).UpdateProfileInfo), varargs...)
}

// UserFollows mocks base method.
func (m *MockUserServiceClient) UserFollows(ctx context.Context, in *proto.UUIDMessage, opts ...grpc.CallOption) (*generated.FollowsMessage, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UserFollows", varargs...)
	ret0, _ := ret[0].(*generated.FollowsMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserFollows indicates an expected call of UserFollows.
func (mr *MockUserServiceClientMockRecorder) UserFollows(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserFollows", reflect.TypeOf((*MockUserServiceClient)(nil).UserFollows), varargs...)
}

// UserSubscriptions mocks base method.
func (m *MockUserServiceClient) UserSubscriptions(ctx context.Context, in *proto.UUIDMessage, opts ...grpc.CallOption) (*generated.SubscriptionsMessage, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UserSubscriptions", varargs...)
	ret0, _ := ret[0].(*generated.SubscriptionsMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserSubscriptions indicates an expected call of UserSubscriptions.
func (mr *MockUserServiceClientMockRecorder) UserSubscriptions(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserSubscriptions", reflect.TypeOf((*MockUserServiceClient)(nil).UserSubscriptions), varargs...)
}

// MockUserServiceServer is a mock of UserServiceServer interface.
type MockUserServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceServerMockRecorder
}

// MockUserServiceServerMockRecorder is the mock recorder for MockUserServiceServer.
type MockUserServiceServerMockRecorder struct {
	mock *MockUserServiceServer
}

// NewMockUserServiceServer creates a new mock instance.
func NewMockUserServiceServer(ctrl *gomock.Controller) *MockUserServiceServer {
	mock := &MockUserServiceServer{ctrl: ctrl}
	mock.recorder = &MockUserServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserServiceServer) EXPECT() *MockUserServiceServerMockRecorder {
	return m.recorder
}

// BecomeCreator mocks base method.
func (m *MockUserServiceServer) BecomeCreator(arg0 context.Context, arg1 *generated.BecameCreatorInfoMessage) (*proto.UUIDResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BecomeCreator", arg0, arg1)
	ret0, _ := ret[0].(*proto.UUIDResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BecomeCreator indicates an expected call of BecomeCreator.
func (mr *MockUserServiceServerMockRecorder) BecomeCreator(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BecomeCreator", reflect.TypeOf((*MockUserServiceServer)(nil).BecomeCreator), arg0, arg1)
}

// CheckIfCreator mocks base method.
func (m *MockUserServiceServer) CheckIfCreator(arg0 context.Context, arg1 *proto.UUIDMessage) (*generated.CheckCreatorMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckIfCreator", arg0, arg1)
	ret0, _ := ret[0].(*generated.CheckCreatorMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckIfCreator indicates an expected call of CheckIfCreator.
func (mr *MockUserServiceServerMockRecorder) CheckIfCreator(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckIfCreator", reflect.TypeOf((*MockUserServiceServer)(nil).CheckIfCreator), arg0, arg1)
}

// DeletePhoto mocks base method.
func (m *MockUserServiceServer) DeletePhoto(arg0 context.Context, arg1 *proto.UUIDMessage) (*proto.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePhoto", arg0, arg1)
	ret0, _ := ret[0].(*proto.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeletePhoto indicates an expected call of DeletePhoto.
func (mr *MockUserServiceServerMockRecorder) DeletePhoto(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePhoto", reflect.TypeOf((*MockUserServiceServer)(nil).DeletePhoto), arg0, arg1)
}

// Donate mocks base method.
func (m *MockUserServiceServer) Donate(arg0 context.Context, arg1 *generated.DonateMessage) (*generated.DonateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Donate", arg0, arg1)
	ret0, _ := ret[0].(*generated.DonateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Donate indicates an expected call of Donate.
func (mr *MockUserServiceServerMockRecorder) Donate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Donate", reflect.TypeOf((*MockUserServiceServer)(nil).Donate), arg0, arg1)
}

// Follow mocks base method.
func (m *MockUserServiceServer) Follow(arg0 context.Context, arg1 *generated.FollowMessage) (*proto.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Follow", arg0, arg1)
	ret0, _ := ret[0].(*proto.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Follow indicates an expected call of Follow.
func (mr *MockUserServiceServerMockRecorder) Follow(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Follow", reflect.TypeOf((*MockUserServiceServer)(nil).Follow), arg0, arg1)
}

// GetProfile mocks base method.
func (m *MockUserServiceServer) GetProfile(arg0 context.Context, arg1 *proto.UUIDMessage) (*generated.UserProfile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfile", arg0, arg1)
	ret0, _ := ret[0].(*generated.UserProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfile indicates an expected call of GetProfile.
func (mr *MockUserServiceServerMockRecorder) GetProfile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfile", reflect.TypeOf((*MockUserServiceServer)(nil).GetProfile), arg0, arg1)
}

// Subscribe mocks base method.
func (m *MockUserServiceServer) Subscribe(arg0 context.Context, arg1 *generated.SubscriptionDetails) (*proto.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", arg0, arg1)
	ret0, _ := ret[0].(*proto.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockUserServiceServerMockRecorder) Subscribe(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockUserServiceServer)(nil).Subscribe), arg0, arg1)
}

// Unfollow mocks base method.
func (m *MockUserServiceServer) Unfollow(arg0 context.Context, arg1 *generated.FollowMessage) (*proto.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unfollow", arg0, arg1)
	ret0, _ := ret[0].(*proto.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Unfollow indicates an expected call of Unfollow.
func (mr *MockUserServiceServerMockRecorder) Unfollow(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unfollow", reflect.TypeOf((*MockUserServiceServer)(nil).Unfollow), arg0, arg1)
}

// UpdatePassword mocks base method.
func (m *MockUserServiceServer) UpdatePassword(arg0 context.Context, arg1 *generated.UpdatePasswordMessage) (*proto.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePassword", arg0, arg1)
	ret0, _ := ret[0].(*proto.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePassword indicates an expected call of UpdatePassword.
func (mr *MockUserServiceServerMockRecorder) UpdatePassword(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePassword", reflect.TypeOf((*MockUserServiceServer)(nil).UpdatePassword), arg0, arg1)
}

// UpdatePhoto mocks base method.
func (m *MockUserServiceServer) UpdatePhoto(arg0 context.Context, arg1 *proto.UUIDMessage) (*generated.ImageID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePhoto", arg0, arg1)
	ret0, _ := ret[0].(*generated.ImageID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePhoto indicates an expected call of UpdatePhoto.
func (mr *MockUserServiceServerMockRecorder) UpdatePhoto(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePhoto", reflect.TypeOf((*MockUserServiceServer)(nil).UpdatePhoto), arg0, arg1)
}

// UpdateProfileInfo mocks base method.
func (m *MockUserServiceServer) UpdateProfileInfo(arg0 context.Context, arg1 *generated.UpdateProfileInfoMessage) (*proto.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProfileInfo", arg0, arg1)
	ret0, _ := ret[0].(*proto.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateProfileInfo indicates an expected call of UpdateProfileInfo.
func (mr *MockUserServiceServerMockRecorder) UpdateProfileInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProfileInfo", reflect.TypeOf((*MockUserServiceServer)(nil).UpdateProfileInfo), arg0, arg1)
}

// UserFollows mocks base method.
func (m *MockUserServiceServer) UserFollows(arg0 context.Context, arg1 *proto.UUIDMessage) (*generated.FollowsMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserFollows", arg0, arg1)
	ret0, _ := ret[0].(*generated.FollowsMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserFollows indicates an expected call of UserFollows.
func (mr *MockUserServiceServerMockRecorder) UserFollows(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserFollows", reflect.TypeOf((*MockUserServiceServer)(nil).UserFollows), arg0, arg1)
}

// UserSubscriptions mocks base method.
func (m *MockUserServiceServer) UserSubscriptions(arg0 context.Context, arg1 *proto.UUIDMessage) (*generated.SubscriptionsMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserSubscriptions", arg0, arg1)
	ret0, _ := ret[0].(*generated.SubscriptionsMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserSubscriptions indicates an expected call of UserSubscriptions.
func (mr *MockUserServiceServerMockRecorder) UserSubscriptions(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserSubscriptions", reflect.TypeOf((*MockUserServiceServer)(nil).UserSubscriptions), arg0, arg1)
}

// mustEmbedUnimplementedUserServiceServer mocks base method.
func (m *MockUserServiceServer) mustEmbedUnimplementedUserServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedUserServiceServer")
}

// mustEmbedUnimplementedUserServiceServer indicates an expected call of mustEmbedUnimplementedUserServiceServer.
func (mr *MockUserServiceServerMockRecorder) mustEmbedUnimplementedUserServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedUserServiceServer", reflect.TypeOf((*MockUserServiceServer)(nil).mustEmbedUnimplementedUserServiceServer))
}

// MockUnsafeUserServiceServer is a mock of UnsafeUserServiceServer interface.
type MockUnsafeUserServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeUserServiceServerMockRecorder
}

// MockUnsafeUserServiceServerMockRecorder is the mock recorder for MockUnsafeUserServiceServer.
type MockUnsafeUserServiceServerMockRecorder struct {
	mock *MockUnsafeUserServiceServer
}

// NewMockUnsafeUserServiceServer creates a new mock instance.
func NewMockUnsafeUserServiceServer(ctrl *gomock.Controller) *MockUnsafeUserServiceServer {
	mock := &MockUnsafeUserServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeUserServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeUserServiceServer) EXPECT() *MockUnsafeUserServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedUserServiceServer mocks base method.
func (m *MockUnsafeUserServiceServer) mustEmbedUnimplementedUserServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedUserServiceServer")
}

// mustEmbedUnimplementedUserServiceServer indicates an expected call of mustEmbedUnimplementedUserServiceServer.
func (mr *MockUnsafeUserServiceServerMockRecorder) mustEmbedUnimplementedUserServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedUserServiceServer", reflect.TypeOf((*MockUnsafeUserServiceServer)(nil).mustEmbedUnimplementedUserServiceServer))
}