// Code generated by MockGen. DO NOT EDIT.
// Source: ./generated/auth_grpc.pb.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	proto "github.com/go-park-mail-ru/2023_1_4from5/internal/models/proto"
	generated "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/grpc/generated"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockAuthServiceClient is a mock of AuthServiceClient interface.
type MockAuthServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceClientMockRecorder
}

// MockAuthServiceClientMockRecorder is the mock recorder for MockAuthServiceClient.
type MockAuthServiceClientMockRecorder struct {
	mock *MockAuthServiceClient
}

// NewMockAuthServiceClient creates a new mock instance.
func NewMockAuthServiceClient(ctrl *gomock.Controller) *MockAuthServiceClient {
	mock := &MockAuthServiceClient{ctrl: ctrl}
	mock.recorder = &MockAuthServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthServiceClient) EXPECT() *MockAuthServiceClientMockRecorder {
	return m.recorder
}

// CheckUser mocks base method.
func (m *MockAuthServiceClient) CheckUser(ctx context.Context, in *generated.User, opts ...grpc.CallOption) (*generated.User, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CheckUser", varargs...)
	ret0, _ := ret[0].(*generated.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUser indicates an expected call of CheckUser.
func (mr *MockAuthServiceClientMockRecorder) CheckUser(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUser", reflect.TypeOf((*MockAuthServiceClient)(nil).CheckUser), varargs...)
}

// CheckUserVersion mocks base method.
func (m *MockAuthServiceClient) CheckUserVersion(ctx context.Context, in *generated.AccessDetails, opts ...grpc.CallOption) (*generated.UserVersion, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CheckUserVersion", varargs...)
	ret0, _ := ret[0].(*generated.UserVersion)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUserVersion indicates an expected call of CheckUserVersion.
func (mr *MockAuthServiceClientMockRecorder) CheckUserVersion(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUserVersion", reflect.TypeOf((*MockAuthServiceClient)(nil).CheckUserVersion), varargs...)
}

// EncryptPwd mocks base method.
func (m *MockAuthServiceClient) EncryptPwd(ctx context.Context, in *generated.EncryptPwdMg, opts ...grpc.CallOption) (*generated.EncryptPwdMg, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "EncryptPwd", varargs...)
	ret0, _ := ret[0].(*generated.EncryptPwdMg)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EncryptPwd indicates an expected call of EncryptPwd.
func (mr *MockAuthServiceClientMockRecorder) EncryptPwd(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EncryptPwd", reflect.TypeOf((*MockAuthServiceClient)(nil).EncryptPwd), varargs...)
}

// IncUserVersion mocks base method.
func (m *MockAuthServiceClient) IncUserVersion(ctx context.Context, in *generated.AccessDetails, opts ...grpc.CallOption) (*proto.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "IncUserVersion", varargs...)
	ret0, _ := ret[0].(*proto.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IncUserVersion indicates an expected call of IncUserVersion.
func (mr *MockAuthServiceClientMockRecorder) IncUserVersion(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncUserVersion", reflect.TypeOf((*MockAuthServiceClient)(nil).IncUserVersion), varargs...)
}

// SignIn mocks base method.
func (m *MockAuthServiceClient) SignIn(ctx context.Context, in *generated.LoginUser, opts ...grpc.CallOption) (*generated.Token, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SignIn", varargs...)
	ret0, _ := ret[0].(*generated.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignIn indicates an expected call of SignIn.
func (mr *MockAuthServiceClientMockRecorder) SignIn(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIn", reflect.TypeOf((*MockAuthServiceClient)(nil).SignIn), varargs...)
}

// SignUp mocks base method.
func (m *MockAuthServiceClient) SignUp(ctx context.Context, in *generated.User, opts ...grpc.CallOption) (*generated.Token, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SignUp", varargs...)
	ret0, _ := ret[0].(*generated.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignUp indicates an expected call of SignUp.
func (mr *MockAuthServiceClientMockRecorder) SignUp(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockAuthServiceClient)(nil).SignUp), varargs...)
}

// MockAuthServiceServer is a mock of AuthServiceServer interface.
type MockAuthServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceServerMockRecorder
}

// MockAuthServiceServerMockRecorder is the mock recorder for MockAuthServiceServer.
type MockAuthServiceServerMockRecorder struct {
	mock *MockAuthServiceServer
}

// NewMockAuthServiceServer creates a new mock instance.
func NewMockAuthServiceServer(ctrl *gomock.Controller) *MockAuthServiceServer {
	mock := &MockAuthServiceServer{ctrl: ctrl}
	mock.recorder = &MockAuthServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthServiceServer) EXPECT() *MockAuthServiceServerMockRecorder {
	return m.recorder
}

// CheckUser mocks base method.
func (m *MockAuthServiceServer) CheckUser(arg0 context.Context, arg1 *generated.User) (*generated.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUser", arg0, arg1)
	ret0, _ := ret[0].(*generated.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUser indicates an expected call of CheckUser.
func (mr *MockAuthServiceServerMockRecorder) CheckUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUser", reflect.TypeOf((*MockAuthServiceServer)(nil).CheckUser), arg0, arg1)
}

// CheckUserVersion mocks base method.
func (m *MockAuthServiceServer) CheckUserVersion(arg0 context.Context, arg1 *generated.AccessDetails) (*generated.UserVersion, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUserVersion", arg0, arg1)
	ret0, _ := ret[0].(*generated.UserVersion)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUserVersion indicates an expected call of CheckUserVersion.
func (mr *MockAuthServiceServerMockRecorder) CheckUserVersion(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUserVersion", reflect.TypeOf((*MockAuthServiceServer)(nil).CheckUserVersion), arg0, arg1)
}

// EncryptPwd mocks base method.
func (m *MockAuthServiceServer) EncryptPwd(arg0 context.Context, arg1 *generated.EncryptPwdMg) (*generated.EncryptPwdMg, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EncryptPwd", arg0, arg1)
	ret0, _ := ret[0].(*generated.EncryptPwdMg)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EncryptPwd indicates an expected call of EncryptPwd.
func (mr *MockAuthServiceServerMockRecorder) EncryptPwd(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EncryptPwd", reflect.TypeOf((*MockAuthServiceServer)(nil).EncryptPwd), arg0, arg1)
}

// IncUserVersion mocks base method.
func (m *MockAuthServiceServer) IncUserVersion(arg0 context.Context, arg1 *generated.AccessDetails) (*proto.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncUserVersion", arg0, arg1)
	ret0, _ := ret[0].(*proto.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IncUserVersion indicates an expected call of IncUserVersion.
func (mr *MockAuthServiceServerMockRecorder) IncUserVersion(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncUserVersion", reflect.TypeOf((*MockAuthServiceServer)(nil).IncUserVersion), arg0, arg1)
}

// SignIn mocks base method.
func (m *MockAuthServiceServer) SignIn(arg0 context.Context, arg1 *generated.LoginUser) (*generated.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIn", arg0, arg1)
	ret0, _ := ret[0].(*generated.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignIn indicates an expected call of SignIn.
func (mr *MockAuthServiceServerMockRecorder) SignIn(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIn", reflect.TypeOf((*MockAuthServiceServer)(nil).SignIn), arg0, arg1)
}

// SignUp mocks base method.
func (m *MockAuthServiceServer) SignUp(arg0 context.Context, arg1 *generated.User) (*generated.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", arg0, arg1)
	ret0, _ := ret[0].(*generated.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignUp indicates an expected call of SignUp.
func (mr *MockAuthServiceServerMockRecorder) SignUp(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockAuthServiceServer)(nil).SignUp), arg0, arg1)
}

// mustEmbedUnimplementedAuthServiceServer mocks base method.
func (m *MockAuthServiceServer) mustEmbedUnimplementedAuthServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedAuthServiceServer")
}

// mustEmbedUnimplementedAuthServiceServer indicates an expected call of mustEmbedUnimplementedAuthServiceServer.
func (mr *MockAuthServiceServerMockRecorder) mustEmbedUnimplementedAuthServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedAuthServiceServer", reflect.TypeOf((*MockAuthServiceServer)(nil).mustEmbedUnimplementedAuthServiceServer))
}

// MockUnsafeAuthServiceServer is a mock of UnsafeAuthServiceServer interface.
type MockUnsafeAuthServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeAuthServiceServerMockRecorder
}

// MockUnsafeAuthServiceServerMockRecorder is the mock recorder for MockUnsafeAuthServiceServer.
type MockUnsafeAuthServiceServerMockRecorder struct {
	mock *MockUnsafeAuthServiceServer
}

// NewMockUnsafeAuthServiceServer creates a new mock instance.
func NewMockUnsafeAuthServiceServer(ctrl *gomock.Controller) *MockUnsafeAuthServiceServer {
	mock := &MockUnsafeAuthServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeAuthServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeAuthServiceServer) EXPECT() *MockUnsafeAuthServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedAuthServiceServer mocks base method.
func (m *MockUnsafeAuthServiceServer) mustEmbedUnimplementedAuthServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedAuthServiceServer")
}

// mustEmbedUnimplementedAuthServiceServer indicates an expected call of mustEmbedUnimplementedAuthServiceServer.
func (mr *MockUnsafeAuthServiceServerMockRecorder) mustEmbedUnimplementedAuthServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedAuthServiceServer", reflect.TypeOf((*MockUnsafeAuthServiceServer)(nil).mustEmbedUnimplementedAuthServiceServer))
}
